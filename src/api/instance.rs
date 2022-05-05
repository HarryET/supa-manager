use bollard::Docker;
use rocket::{Route, State};
use rocket::serde::json::Json;
use crate::api::errors::ApiError;
use crate::models::Instance;
use crate::PostgresDbConn;
use crate::schema::*;
use diesel::prelude::*;
use uuid::Uuid;

pub fn routes() -> Vec<Route> {
    routes![get_instance, delete_instance, start_instance, stop_instance, restart_instance]
}

#[get("/<id>")]
async fn get_instance(id: &str, db: PostgresDbConn) -> Result<Json<Instance>, ApiError> {
    let uuid = Uuid::parse_str(id)?;
    let instance: Instance = db.run(move |conn| {
        instances::table
            .filter(instances::id.eq(uuid))
            .select(instances::all_columns)
            .first(conn)
    }).await?;

    Ok(Json(instance))
}

#[delete("/<id>")]
async fn delete_instance(id: &str, docker: &State<Docker>) -> String {
    format!("Delete an Instance, {}", id)
}

#[post("/<id>/start")]
async fn start_instance(id: &str, docker: &State<Docker>) -> String {
    format!("Delete an Instance, {}", id)
}

#[post("/<id>/stop")]
async fn stop_instance(id: &str, docker: &State<Docker>) -> String {
    format!("Delete an Instance, {}", id)
}

#[post("/<id>/restart")]
async fn restart_instance(id: &str, docker: &State<Docker>) -> String {
    format!("Delete an Instance, {}", id)
}