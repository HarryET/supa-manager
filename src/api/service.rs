use bollard::Docker;
use rocket::{Route, State};

pub fn routes() -> Vec<Route> {
    routes![get_service, start_service, stop_service, restart_service]
}

#[get("/<id>/services/<service>")]
async fn get_service(id: &str, service: &str, _docker: &State<Docker>) -> String {
    format!("Get Service, {} for instance, {}", service, id)
}

#[post("/<id>/services/<service>/start")]
async fn start_service(id: &str, service: &str, _docker: &State<Docker>) -> String {
    format!("Start Service, {} for instance, {}", service, id)
}

#[post("/<id>/services/<service>/stop")]
async fn stop_service(id: &str, service: &str, _docker: &State<Docker>) -> String {
    format!("Stop Service, {} for instance, {}", service, id)
}

#[post("/<id>/services/<service>/restart")]
async fn restart_service(id: &str, service: &str, _docker: &State<Docker>) -> String {
    format!("Restart Service, {} for instance, {}", service, id)
}