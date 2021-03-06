use std::collections::HashMap;
use bollard::container::{CreateContainerOptions, NetworkingConfig};
use bollard::{container, Docker};
use bollard::errors::Error;
use bollard::models::{ContainerCreateResponse, EndpointSettings};
use bollard::network::ConnectNetworkOptions;

// pub fn download_image(docker: &Docker, image_name: &str, tag: &str, repo: Option<&str>) -> impl Stream<Item = Result<CreateImageInfo, Error>> {
//     &docker
//         .create_image(
//             Some(CreateImageOptions {
//                 from_image: format!("{}:{}", image_name, tag),
//                 repo: repo.unwrap_or("registry.hub.docker.com").to_string(),
//                 ..Default::default()
//             }),
//             None,
//             None,
//         )
// }

pub async fn create_container(docker: &Docker, name: String, cfg: container::Config<String>) -> Result<ContainerCreateResponse, Error> {
    docker.create_container(Some(CreateContainerOptions {
        name
    }), cfg).await
}

pub async fn start_container(docker: &Docker, id: &String) -> Result<(), Error> {
    docker.start_container::<String>(&*id, None).await
}

pub async fn start_containers(docker: &Docker, ids: Vec<String>) -> Result<(), Error> {
    for id in ids {
        let r = start_container(docker, &id).await;
        if r.is_err() {
            return r;
        }
    }

    Ok(())
}

pub fn new_net_config(net_name: String, alias: Vec<String>) -> NetworkingConfig<String> {
    let mut cfg: HashMap<String, EndpointSettings> = HashMap::new();

    cfg.insert(net_name, EndpointSettings {
        aliases: Some(alias.clone()),
        ..Default::default()
    });

    NetworkingConfig {
        endpoints_config: cfg
    }
}

pub async fn attach_to_traefik(docker: &Docker, container_id: String, alias: Vec<String>) -> Result<(), Error> {
    docker.connect_network("traefik", ConnectNetworkOptions {
        container: container_id,
        endpoint_config: EndpointSettings {
            aliases: Some(alias.clone()),
            ..Default::default()
        }
    }).await
}