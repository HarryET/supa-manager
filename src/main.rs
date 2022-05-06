mod api;
mod models;
mod schema;
mod database;
mod constants;
mod utils;
mod config;

#[macro_use]
extern crate rocket;
#[macro_use]
extern crate diesel;
#[macro_use]
extern crate diesel_migrations;
extern crate dotenv;

use diesel::prelude::*;
use diesel::pg::PgConnection;
use dotenv::dotenv;
use std::env;
use bollard::Docker;
use diesel_migrations::embed_migrations;
use crate::constants::{DB_IMAGE, DB_TAG, GOTRUE_IMAGE, GOTRUE_TAG, META_IMAGE, META_TAG, POSTGREST_IMAGE, POSTGREST_TAG, REALTIME_IMAGE, REALTIME_TAG, STUDIO_IMAGE, STUDIO_TAG};
use crate::database::PostgresDbConn;

embed_migrations!();

#[get("/")]
fn index() -> &'static str {
    "SupaManager, a project by Harry Bairstow;\nManage self-hosted Supabase instances with an easy to use API & Web Portal (soon)"
}

#[rocket::main]
async fn main() {
    dotenv().ok();

    let _ = migrate();

    let docker = Docker::connect_with_local_defaults().expect("Failed to connect to docker.");

    let _ = rocket::build()
        .manage(docker)
        .mount("/", routes![index])
        // TODO Can dynamic route parts be here?
        .mount("/instances", crate::api::instances::routes())
        .mount("/instances", crate::api::instance::routes())
        .mount("/instances", crate::api::service::routes())
        .attach(PostgresDbConn::fairing())
        .launch()
        .await;
}

fn migrate() -> Result<(), diesel_migrations::RunMigrationsError> {
    let database_url = env::var("DATABASE_URL")
        .expect("DATABASE_URL must be set");
    let db_conn = PgConnection::establish(&database_url)
        .expect(&format!("Error connecting to {}", database_url));

    embedded_migrations::run(&db_conn)
}