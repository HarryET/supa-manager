mod api;
mod models;
mod schema;
mod database;
mod constants;
mod utils;
mod config;
mod ui;

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
use std::path::{Path, PathBuf};
use bollard::Docker;
use diesel_migrations::embed_migrations;
use rocket::fs::NamedFile;
use crate::constants::{DB_IMAGE, DB_TAG, GOTRUE_IMAGE, GOTRUE_TAG, META_IMAGE, META_TAG, POSTGREST_IMAGE, POSTGREST_TAG, REALTIME_IMAGE, REALTIME_TAG, STUDIO_IMAGE, STUDIO_TAG};
use crate::database::PostgresDbConn;
use rocket_dyn_templates::Template;

embed_migrations!();

#[get("/<file..>")]
async fn static_files(file: PathBuf) -> Option<NamedFile> {
    NamedFile::open(Path::new("public/").join(file)).await.ok()
}

#[rocket::main]
async fn main() {
    dotenv().ok();

    let _ = migrate();

    let docker = Docker::connect_with_local_defaults().expect("Failed to connect to docker.");

    let _ = rocket::build()
        .manage(docker)
        .mount("/", ui::routes())
        .mount("/public", routes![static_files])
        // TODO Can dynamic route parts be here?
        .mount("/instances", api::instances::routes())
        .mount("/instances", api::instance::routes())
        .mount("/instances", api::service::routes())
        .attach(PostgresDbConn::fairing())
        .attach(Template::fairing())
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