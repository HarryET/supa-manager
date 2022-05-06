use std::collections::HashMap;
use serde::{Serialize};
use uuid::Uuid;
use crate::constants::DOMAIN;
use crate::schema::instances;

#[derive(Queryable, Identifiable, Insertable, Debug, PartialEq, Serialize, Clone)]
#[diesel(table_name = crate::models::Instance)]
pub struct Instance {
    pub id: Uuid,
    pub nickname: Option<String>,
    pub hostname: String,

    pub studio_container_id: Option<String>,
    pub studio_enabled: bool,

    pub database_container_id: String,
    pub gotrue_container_id: String,
    pub realtime_container_id: String,
    pub postgrest_container_id: String,
    pub postgres_meta_container_id: String,
}

impl Instance {
    pub fn to_url(&self) -> String {
        format!("http://{}.{}", self.hostname, DOMAIN)
    }

    pub fn get_docker_project_name(&self) -> String {
        let mut n = "".to_string();
        if self.nickname.is_some() {
            n = format!("_{}", self.nickname.as_ref().unwrap())
        }
        format!("{}{}", self.hostname.to_string(), n)
    }

    pub fn to_private_docker_labels(&self) -> HashMap<String, String> {
        let mut map = HashMap::new();

        for (k, v) in vec![
            ("supa-manager.instance", self.id.to_string().as_str()),
            ("traefik.enable", "false"),
            /*
                NOTE (HarryET): This is a hack to make docker desktop group the containers together,
                                I found what label to use by looking at the docker compose v2 source and
                                the label is in this file:
                                https://github.com/docker/compose/blob/v2/pkg/api/labels.go#L29
            */
            ("com.docker.compose.project", self.get_docker_project_name().as_str()),
        ] {
            map.insert(k.to_string(), v.to_string());
        }

        return map;
    }

    pub fn to_docker_labels(&self, service: String, path: Option<String>, entrypoint: Option<String>, port: Option<String>) -> HashMap<String, String> {
        let mut map = HashMap::new();

        let mut p = "".to_string();

        if path.is_some() {
            p = format!("&& Path(`{}`)", path.unwrap())
        }

        for (k, v) in vec![
            ("supa-manager.instance", self.id.to_string().as_str()),
            /*
            NOTE (HarryET): This is a hack to make docker desktop group the containers together,
                            I found what label to use by looking at the docker compose v2 source and
                            the label is in this file:
                            https://github.com/docker/compose/blob/v2/pkg/api/labels.go#L29
             */
            ("com.docker.compose.project", self.get_docker_project_name().as_str()),
            ("traefik.enable", "true"),
            (format!("traefik.http.routers.{}-{}.rule", self.hostname, service).as_str(), format!("Host(`{}.{}`){}", self.hostname, DOMAIN, p).as_str()),
            (format!("traefik.http.routers.{}-{}.entrypoints", self.hostname, service).as_str(), entrypoint.unwrap_or("web".to_string()).as_str()),
            (format!("traefik.http.routers.{}-{}.service", self.hostname, service).as_str(), format!("{}-{}-svc", self.hostname, service).as_str()),
            (format!("traefik.http.services.{}-{}-svc.loadbalancer.server.port", self.hostname, service).as_str(), port.unwrap_or("80".to_string()).as_str()),
            //(format!("traefik.http.routers.{}-{}.certresolver", self.hostname, service).as_str(), "supa-manager-certs"),
        ] {
            map.insert(k.to_string(), v.to_string());
        }

        return map;
    }
}

pub fn new_blank_instance() -> Instance {
    Instance {
        id: Uuid::new_v4(),
        nickname: None,
        hostname: cuid::cuid().unwrap(),
        studio_container_id: None,
        studio_enabled: false,
        database_container_id: "".to_string(),
        gotrue_container_id: "".to_string(),
        realtime_container_id: "".to_string(),
        postgrest_container_id: "".to_string(),
        postgres_meta_container_id: "".to_string()
    }
}