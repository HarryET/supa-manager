use std::collections::HashMap;
use chrono::prelude::*;
use rocket::Route;
use rocket_dyn_templates::Template;

pub fn routes() -> Vec<Route> {
    routes![index, login]
}

#[get("/")]
fn index() -> Template {
    let args: HashMap<String, String> = HashMap::new();
    Template::render("index", &args)
}

#[get("/login")]
fn login() -> Template {
    let mut args: HashMap<&str, String> = HashMap::new();
    args.insert("version", "0.1.0".to_string()); // TODO make dynamic
    args.insert("year", Utc::now().year().to_string());
    Template::render("login", &args)
}