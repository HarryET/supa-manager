use std::borrow::Cow;
use bollard::{container, Docker};
use bollard::network::CreateNetworkOptions;
use diesel::QueryDsl;
use rocket::{Route, State};
use crate::database::PostgresDbConn;
use rocket::serde::{json::Json};
use crate::models::Instance;
use crate::schema::*;
use diesel::prelude::*;
use passwords::PasswordGenerator;
use crate::api::errors::ApiError;
use serde::Deserialize;
use crate::{DB_IMAGE, DB_TAG, GOTRUE_IMAGE, GOTRUE_TAG, META_IMAGE, META_TAG, POSTGREST_IMAGE, POSTGREST_TAG, REALTIME_IMAGE, REALTIME_TAG, STUDIO_IMAGE, STUDIO_TAG};
use crate::constants::{DB_MIGRATIONS_DIR, DOMAIN};
use crate::utils::{attach_to_traefik, create_container, new_net_config, start_container, start_containers};
use std::default::Default;
use bollard::models::{HostConfig, RestartPolicy, RestartPolicyNameEnum};
use rocket::response::status;

pub fn routes() -> Vec<Route> {
    routes![get_instances, new_instance]
}

#[get("/")]
async fn get_instances(db: PostgresDbConn) -> Result<Json<Vec<Instance>>, ApiError> {
    let instances: Vec<Instance> = db.run(move |conn| {
        instances::table
            .select(instances::all_columns)
            .load(conn)
    }).await?;

    Ok(Json(instances))
}

#[derive(Deserialize)]
struct NewInstance {
    nickname: Option<String>,
    enable_studio: Option<bool>,
}

struct PostgresInstanceData {
    username: String,
    password: String,
    database: String,
}

impl PostgresInstanceData {
    pub fn to_url(&self, container_name: String) -> String {
        format!("postgres://{}:{}@{}:5432/{}", self.username, self.password, container_name, self.database)
    }

    pub fn to_env(&self) -> Vec<String> {
        vec![format!("POSTGRES_USER={}", self.username), format!("POSTGRES_PASSWORD={}", self.password), format!("POSTGRES_DB={}", self.database)]
    }

    pub fn to_meta_env(&self) -> Vec<String> {
        vec![format!("PG_META_DB_USER={}", self.username), format!("PG_META_DB_PASSWORD={}", self.password), format!("PG_META_DB_NAME={}", self.database)]
    }

    pub fn to_rest_env(&self, container_name: String) -> Vec<String> {
        vec![format!("PGRST_DB_URI=postgres://{}:{}@{}:5432/{}", self.username, self.password, container_name, self.database), "PGRST_DB_SCHEMA=public".to_string(), "PGRST_DB_ANON_ROLE=anon".to_string()]
    }

    pub fn to_realtime_env(&self) -> Vec<String> {
        vec![format!("DB_USER={}", self.username), format!("DB_PASSWORD={}", self.password), format!("DB_NAME={}", self.database)]
    }
}

fn gen_pg_instance_data() -> PostgresInstanceData {
    PostgresInstanceData {
        username: "postgres".to_string(),
        password: PasswordGenerator {
            length: 16,
            numbers: true,
            lowercase_letters: true,
            uppercase_letters: true,
            symbols: false,
            spaces: false,
            exclude_similar_characters: true,
            strict: true,
        }.generate_one().unwrap(),
        database: "postgres".to_string(),
    }
}

// TODO optimise!
#[post("/", data = "<body>")]
async fn new_instance(body: Json<NewInstance>, docker: &State<Docker>, db: PostgresDbConn) -> Result<status::Created<Json<Instance>>, ApiError> {
    let mut instance = crate::models::new_blank_instance();
    instance.nickname = body.nickname.clone();
    instance.studio_enabled = body.enable_studio.unwrap_or(false);
    let instance_id = instance.id.clone();

    let host_config = HostConfig {
        restart_policy: Some(RestartPolicy {
            name: Some(RestartPolicyNameEnum::UNLESS_STOPPED),
            ..Default::default()
        }),
        ..Default::default()
    };

    let network_name = format!("{}_net", instance.hostname);
    let _ = docker.create_network(CreateNetworkOptions {
        name: network_name.clone(),
        ..Default::default()
    }).await?;

    let mut containers: Vec<String> = vec![];

    let pg_instance_data = gen_pg_instance_data();
    let mut pg_host_config = host_config.clone();
    pg_host_config.binds = Some(vec![format!("{}:/docker-entrypoint-initdb.d", DB_MIGRATIONS_DIR)]);
    let db_container = create_container(&docker, format!("{}_postgres", instance.hostname), container::Config {
        image: Some(format!("{}:{}", DB_IMAGE, DB_TAG)),
        env: Some(pg_instance_data.to_env()),
        // TODO fix so it uses tcp not HTTP
        labels: Some(instance.to_docker_labels("postgres".to_string(), None, Some("pg".to_string()), Some("5432".to_string()))),
        networking_config: Some(new_net_config(network_name.clone(), vec!["postgres".to_string()])),
        host_config: Some(pg_host_config),
        ..Default::default()
    }).await?;
    instance.database_container_id = db_container.id;
    containers.push(instance.database_container_id.clone());

    let mut db_meta_container_env = vec!["PG_META_PORT=8080".to_string(), "PG_META_DB_PORT=5432".to_string(), format!("PG_META_DB_HOST={}_postgres", instance.hostname)];
    db_meta_container_env.append(&mut pg_instance_data.to_meta_env());
    let db_meta_container = create_container(&docker, format!("{}_postgres_meta", instance.hostname), container::Config {
        image: Some(format!("{}:{}", META_IMAGE, META_TAG)),
        env: Some(db_meta_container_env),
        labels: Some(instance.to_private_docker_labels()),
        networking_config: Some(new_net_config(network_name.clone(), vec!["postgres_meta".to_string()])),
        host_config: Some(host_config.clone()),
        ..Default::default()
    }).await?;
    instance.postgres_meta_container_id = db_meta_container.id;
    containers.push(instance.postgres_meta_container_id.clone());

    let jwt_secret = uuid::Uuid::new_v4().to_string();
    let mut postgrest_container_env = vec![format!("PGRST_JWT_SECRET={}", jwt_secret)];
    postgrest_container_env.append(&mut pg_instance_data.to_rest_env(format!("{}_postgres", instance.hostname)));
    let postgrest_container = create_container(&docker, format!("{}_postgrest", instance.hostname), container::Config {
        image: Some(format!("{}:{}", POSTGREST_IMAGE, POSTGREST_TAG)),
        env: Some(postgrest_container_env),
        labels: Some(instance.to_docker_labels("postgrest".to_string(), Some("/rest/v1".to_string()), None, Some("3000".to_string()))),
        networking_config: Some(new_net_config(network_name.clone(), vec!["postgrest".to_string()])),
        host_config: Some(host_config.clone()),
        ..Default::default()
    }).await?;
    instance.postgrest_container_id = postgrest_container.id;
    containers.push(instance.postgrest_container_id.clone());
    attach_to_traefik(&docker, instance.postgrest_container_id.clone(), vec!["postgrest".to_string()]).await?;

    let gotrue_config = crate::config::gotrue::new(&instance, format!("{}?search_path=auth", pg_instance_data.to_url(format!("{}_postgres", instance.hostname))), jwt_secret.clone());
    let gotrue_container = create_container(&docker, format!("{}_gotrue", instance.hostname), container::Config {
        image: Some(format!("{}:{}", GOTRUE_IMAGE, GOTRUE_TAG)),
        env: Some(gotrue_config.to_env()),
        labels: Some(instance.to_docker_labels("gotrue".to_string(), Some("/auth/v1".to_string()), None, None)),
        networking_config: Some(new_net_config(network_name.clone(), vec!["gotrue".to_string()])),
        host_config: Some(host_config.clone()),
        ..Default::default()
    }).await?;
    instance.gotrue_container_id = gotrue_container.id;
    containers.push(instance.gotrue_container_id.clone());
    attach_to_traefik(&docker, instance.gotrue_container_id.clone(), vec!["gotrue".to_string()]).await?;

    let mut realtime_container_env = vec!["DB_PORT=5432".to_string(), format!("DB_HOST={}", format!("{}_postgres", instance.hostname)), "PORT=8080".to_string(), format!("HOSTNAME={}.{}", instance.hostname.clone(), DOMAIN), "SECURE_CHANNELS=true".to_string(), format!("JWT_SECRET={}", jwt_secret)];
    realtime_container_env.append(&mut pg_instance_data.to_realtime_env());
    let realtime_container = create_container(&docker, format!("{}_realtime", instance.hostname), container::Config {
        image: Some(format!("{}:{}", REALTIME_IMAGE, REALTIME_TAG)),
        env: Some(realtime_container_env),
        // TODO make sure traefik config is correct for realtime, https://github.com/supabase/supabase/blob/master/docker/volumes/api/kong.yml#L125
        labels: Some(instance.to_docker_labels("realtime".to_string(), Some("/realtime/v1".to_string()), None, Some("8080".to_string()))),
        networking_config: Some(new_net_config(network_name.clone(), vec!["realtime".to_string()])),
        host_config: Some(host_config.clone()),
        ..Default::default()
    }).await?;
    instance.realtime_container_id = realtime_container.id;
    containers.push(instance.realtime_container_id.clone());
    attach_to_traefik(&docker, instance.realtime_container_id.clone(), vec!["realtime".to_string()]).await?;

    if instance.studio_enabled {
        let studio_container = create_container(&docker, format!("{}_studio", instance.hostname), container::Config {
            image: Some(format!("{}:{}", STUDIO_IMAGE, STUDIO_TAG)),
            env: Some(vec![format!("POSTGRES_PASSWORD={}", pg_instance_data.password), "STUDIO_PG_META_URL=http://postgres_meta:8080".to_string(), format!("SUPABASE_URL=http://{}.{}", instance.hostname, DOMAIN), format!("SUPABASE_REST_URL=http://{}.{}/rest/v1", instance.hostname, DOMAIN)]),
            labels: Some(instance.to_docker_labels("studio".to_string(), None, None, Some("3000".to_string()))),
            networking_config: Some(new_net_config(network_name.clone(), vec!["studio".to_string()])),
            host_config: Some(host_config.clone()),
            ..Default::default()
        }).await?;
        instance.studio_container_id = Some(studio_container.id);
        attach_to_traefik(&docker, instance.studio_container_id.clone().unwrap(), vec!["studio".to_string()]).await?;
        start_container(docker, instance.studio_container_id.as_ref().unwrap()).await?;
    }

    let db_instance: Instance = db.run(move |conn| {
        diesel::insert_into(instances::table)
            .values(&instance)
            .returning(instances::all_columns)
            .get_result::<Instance>(&*conn)
    }).await?;

    let _ = start_containers(docker, containers).await;

    // TODO make dynamic!
    Ok(status::Created::<Json<Instance>>::new(Cow::from(format!("http://localhost:8000/instances/{}", instance_id))).body(Json(db_instance)))
}
