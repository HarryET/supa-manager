use bollard::container::{RemoveContainerOptions, RestartContainerOptions, StartContainerOptions, StopContainerOptions};
use bollard::Docker;
use rocket::{Route, State};
use rocket::serde::json::Json;
use crate::api::errors::ApiError;
use crate::models::Instance;
use crate::PostgresDbConn;
use crate::schema::*;
use diesel::prelude::*;
use rocket::http::Status;
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
async fn delete_instance(id: &str, docker: &State<Docker>, db: PostgresDbConn) -> Result<Status, ApiError> {
    let uuid = Uuid::parse_str(id)?;
    let instance: Instance = db.run(move |conn| {
        instances::table
            .filter(instances::id.eq(uuid))
            .select(instances::all_columns)
            .first(conn)
    }).await?;

    let rm_opts = RemoveContainerOptions {
        v: true,
        ..Default::default()
    };

    let stop_opts = StopContainerOptions {
        ..Default::default()
    };

    if instance.studio_enabled {
        docker.stop_container(instance.studio_container_id.as_ref().unwrap().as_str(), Some(stop_opts)).await?;
        docker.remove_container(instance.studio_container_id.unwrap().as_str(), Some(rm_opts)).await?;
    }

    docker.stop_container(instance.postgres_meta_container_id.clone().as_str(), Some(stop_opts)).await?;
    docker.remove_container(instance.postgres_meta_container_id.as_str(), Some(rm_opts)).await?;

    docker.stop_container(instance.realtime_container_id.clone().as_str(), Some(stop_opts)).await?;
    docker.remove_container(instance.realtime_container_id.as_str(), Some(rm_opts)).await?;

    docker.stop_container(instance.gotrue_container_id.clone().as_str(), Some(stop_opts)).await?;
    docker.remove_container(instance.gotrue_container_id.as_str(), Some(rm_opts)).await?;

    docker.stop_container(instance.postgrest_container_id.clone().as_str(), Some(stop_opts)).await?;
    docker.remove_container(instance.postgrest_container_id.as_str(), Some(rm_opts)).await?;

    docker.stop_container(instance.database_container_id.clone().as_str(), Some(stop_opts)).await?;
    docker.remove_container(instance.database_container_id.as_str(), Some(rm_opts)).await?;

    Ok(Status::Ok)
}

#[post("/<id>/start")]
async fn start_instance(id: &str, docker: &State<Docker>, db: PostgresDbConn) -> Result<Status, ApiError> {
    let uuid = Uuid::parse_str(id)?;
    let instance: Instance = db.run(move |conn| {
        instances::table
            .filter(instances::id.eq(uuid))
            .select(instances::all_columns)
            .first(conn)
    }).await?;

    let start_opts = StartContainerOptions::<String> {
        ..Default::default()
    };

    let g_start_opts = || { start_opts.clone() };

    if instance.studio_enabled {
        docker.start_container(instance.studio_container_id.as_ref().unwrap().as_str(), Some(g_start_opts())).await?;
    }

    docker.start_container(instance.postgres_meta_container_id.clone().as_str(), Some(g_start_opts())).await?;
    docker.start_container(instance.realtime_container_id.clone().as_str(), Some(g_start_opts())).await?;
    docker.start_container(instance.gotrue_container_id.clone().as_str(), Some(g_start_opts())).await?;
    docker.start_container(instance.postgrest_container_id.clone().as_str(), Some(g_start_opts())).await?;
    docker.start_container(instance.database_container_id.clone().as_str(), Some(g_start_opts())).await?;

    Ok(Status::Ok)
}

#[post("/<id>/stop")]
async fn stop_instance(id: &str, docker: &State<Docker>, db: PostgresDbConn) -> Result<Status, ApiError> {
    let uuid = Uuid::parse_str(id)?;
    let instance: Instance = db.run(move |conn| {
        instances::table
            .filter(instances::id.eq(uuid))
            .select(instances::all_columns)
            .first(conn)
    }).await?;

    let stop_opts = StopContainerOptions {
        ..Default::default()
    };

    if instance.studio_enabled {
        docker.stop_container(instance.studio_container_id.as_ref().unwrap().as_str(), Some(stop_opts)).await?;
    }

    docker.stop_container(instance.postgres_meta_container_id.clone().as_str(), Some(stop_opts)).await?;
    docker.stop_container(instance.realtime_container_id.clone().as_str(), Some(stop_opts)).await?;
    docker.stop_container(instance.gotrue_container_id.clone().as_str(), Some(stop_opts)).await?;
    docker.stop_container(instance.postgrest_container_id.clone().as_str(), Some(stop_opts)).await?;
    docker.stop_container(instance.database_container_id.clone().as_str(), Some(stop_opts)).await?;

    Ok(Status::Ok)
}

#[post("/<id>/restart")]
async fn restart_instance(id: &str, docker: &State<Docker>, db: PostgresDbConn) -> Result<Status, ApiError> {
    let uuid = Uuid::parse_str(id)?;
    let instance: Instance = db.run(move |conn| {
        instances::table
            .filter(instances::id.eq(uuid))
            .select(instances::all_columns)
            .first(conn)
    }).await?;

    let restart_opts = RestartContainerOptions {
        ..Default::default()
    };

    if instance.studio_enabled {
        docker.restart_container(instance.studio_container_id.unwrap().as_str(), Some(restart_opts)).await?;
    }

    docker.restart_container(instance.postgres_meta_container_id.as_str(), Some(restart_opts)).await?;
    docker.restart_container(instance.realtime_container_id.as_str(), Some(restart_opts)).await?;
    docker.restart_container(instance.gotrue_container_id.as_str(), Some(restart_opts)).await?;
    docker.restart_container(instance.postgrest_container_id.as_str(), Some(restart_opts)).await?;
    // No reason to restart database!
    //docker.restart_container(instance.database_container_id.as_str(), Some(restart_opts)).await?;

    Ok(Status::Ok)
}