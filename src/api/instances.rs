use std::collections::HashMap;
use std::fmt::format;
use bollard::container::CreateContainerOptions;
use bollard::{container, Docker};
use bollard::models::NetworkingConfig;
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
use crate::{DB_IMAGE, DB_TAG, META_IMAGE, META_TAG, STUDIO_IMAGE, STUDIO_TAG};
use crate::constants::DOMAIN;
use crate::utils::{create_container, new_net_config, start_container, start_containers};
use std::default::Default;

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

// #[derive(Deserialize)]
// struct ServiceOverrides {
//     kong: ServiceOverride,
//     studio: ServiceOverride,
//     database: ServiceOverride,
//     realtime: ServiceOverride,
//     gotrue: ServiceOverride,
//     postgrest: ServiceOverride
// }
//
// #[derive(Deserialize)]
// struct ServiceOverride {
//     image: Option<String>,
//     tag: Option<String>
// }

struct PostgresInstanceData {
    username: String,
    password: String,
    database: String,
}

impl PostgresInstanceData {
    pub fn to_env(&self) -> Vec<String> {
        vec![format!("POSTGRES_USER={}", self.username), format!("POSTGRES_PASSWORD={}", self.password), format!("POSTGRES_DB={}", self.database)]
    }

    pub fn to_meta_env(&self) -> Vec<String> {
        vec![format!("PG_META_DB_USER={}", self.username), format!("PG_META_DB_PASSWORD={}", self.password), format!("PG_META_DB_NAME={}", self.database)]
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

#[post("/", data = "<body>")]
async fn new_instance(body: Json<NewInstance>, docker: &State<Docker>) -> Result<Json<Instance>, ApiError> {
    let mut instance = crate::models::new_blank_instance();
    instance.nickname = body.nickname.clone();
    instance.studio_enabled = body.enable_studio.unwrap_or(false);

    let network_name = format!("{}_net", instance.hostname);
    let network = docker.create_network(CreateNetworkOptions {
        name: network_name.clone(),
        ..Default::default()
    }).await?;

    let mut containers: Vec<&String> = vec![];

    let pg_instance_data = gen_pg_instance_data();
    let db_container = create_container(&docker, format!("{}_postgres", instance.hostname), container::Config {
        image: Some(format!("{}:{}", DB_IMAGE, DB_TAG)),
        env: Some(pg_instance_data.to_env()),
        labels: Some(instance.to_docker_labels("postgres".to_string(), None, Some("pg".to_string()), Some("5432".to_string()))),
        networking_config: Some(new_net_config(network_name.clone(), vec!["postgres".to_string()], false)),
        ..Default::default()
    }).await?;
    instance.database_container_id = db_container.id;
    containers.push(&instance.database_container_id);

    let mut db_meta_container_env = vec!["PG_META_PORT=8080".to_string(), "PG_META_DB_PORT=5432".to_string(), format!("PG_META_DB_HOST=\"{}_postgres\"", instance.hostname)];
    db_meta_container_env.append(&mut pg_instance_data.to_meta_env());
    let db_meta_container = create_container(&docker, format!("{}_postgres_meta", instance.hostname), container::Config {
        image: Some(format!("{}:{}", META_IMAGE, META_TAG)),
        env: Some(db_meta_container_env),
        networking_config: Some(new_net_config(network_name.clone(), vec!["postgres_meta".to_string()], false)),
        ..Default::default()
    }).await?;
    instance.postgres_meta_container_id = db_meta_container.id;
    containers.push(&instance.postgres_meta_container_id);

    if instance.studio_enabled {
        let studio_container = create_container(&docker, format!("{}_studio", instance.hostname), container::Config {
            image: Some(format!("{}:{}", STUDIO_IMAGE, STUDIO_TAG)),
            env: Some(vec![format!("POSTGRES_PASSWORD={}", pg_instance_data.password), "STUDIO_PG_META_URL=http://postgres_meta:8080".to_string(), format!("SUPABASE_URL=http://{}.{}", instance.hostname, DOMAIN), format!("SUPABASE_REST_URL={}", "")]),
            labels: Some(instance.to_docker_labels("studio".to_string(), None, None, Some("3000".to_string()))),
            networking_config: Some(new_net_config(network_name.clone(), vec!["studio".to_string()], true)),
            ..Default::default()
        }).await?;
        instance.studio_container_id = Some(studio_container.id);
        start_container(docker, instance.studio_container_id.as_ref().unwrap()).await?;
    }

    start_containers(docker, containers).await?;

    Ok(Json(instance))
}
