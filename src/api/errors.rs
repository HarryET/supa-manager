use std::io::Cursor;
use bollard::errors::Error;

use rocket::request::Request;
use rocket::response::{self, Response, Responder};
use rocket::http::{ContentType, Status};
use rocket::serde::json::serde_json::json;

#[derive(Debug)]
pub enum ApiError {
    DieselError(diesel::result::Error),
    UuidError(uuid::Error),
    DockerError(bollard::errors::Error)
}

impl<'r> Responder<'r, 'r> for ApiError {
    // TODO better error handling
    fn respond_to(self, _: &Request) -> response::Result<'r> {
        match self {
            ApiError::DieselError(e) => {
                match e {
                    diesel::result::Error::NotFound => {
                        let body = "{\"message\": \"Not Found.\"}";
                        Response::build()
                            .sized_body(body.len(), Cursor::new(body))
                            .header(ContentType::JSON)
                            .status(Status::NotFound)
                            .ok()
                    },
                    _ => {
                        let body = "{\"message\": \"Internal Server Error\"}";
                        Response::build()
                            .sized_body(body.len(), Cursor::new(body))
                            .header(ContentType::JSON)
                            .status(Status::InternalServerError)
                            .ok()
                    }
                }
            },
            e => {
                dbg!(e);
                let body = "{\"message\": \"Internal Server Error\"}";
                Response::build()
                    .sized_body(body.len(), Cursor::new(body))
                    .header(ContentType::JSON)
                    .status(Status::InternalServerError)
                    .ok()
            }
        }
    }
}

impl From<diesel::result::Error> for ApiError {
    fn from(e: diesel::result::Error) -> Self {
        ApiError::DieselError(e)
    }
}

impl From<uuid::Error> for ApiError {
    fn from(e: uuid::Error) -> Self {
        ApiError::UuidError(e)
    }
}

impl From<bollard::errors::Error> for ApiError {
    fn from(e: bollard::errors::Error) -> Self {
        ApiError::DockerError(e)
    }
}