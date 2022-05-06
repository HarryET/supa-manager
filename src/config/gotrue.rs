use crate::models::Instance;
use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Debug)]
pub struct GoTrueConfig {
    jwt_secret: String,
    jwt_expiry: u32,
    jwt_default_group_name: String,
    db_driver: String,
    db_namespace: String,
    api_external_url: String,
    host: String,
    port: String,

    disable_signup: String,
    site_url: String,
    mailer_auto_confirm: String,
    log_level: String,
    operator_token: String,
    database_url: String

    // TODO More Config
}

pub fn new(instance: &Instance, pg_url: String, jwt_secret: String) -> GoTrueConfig {
    GoTrueConfig {
        jwt_secret,
        jwt_expiry: 3600,
        jwt_default_group_name: "authenticated".to_string(),
        db_driver: "postgres".to_string(),
        db_namespace: "auth".to_string(),
        api_external_url: instance.to_url(),
        host: "0.0.0.0".to_string(),
        port: "8080".to_string(),
        disable_signup: "false".to_string(),
        site_url: instance.to_url(),
        mailer_auto_confirm: "true".to_string(),
        log_level: "INFO".to_string(),
        // TODO what is operator token?
        operator_token: "".to_string(),
        database_url: pg_url
    }
}

impl GoTrueConfig {
    // TODO make better with serde!
    pub fn to_env(&self) -> Vec<String> {
        vec![
            format!("GOTRUE_JWT_SECRET={}", self.jwt_secret),
            format!("GOTRUE_JWT_EXP={}", self.jwt_expiry),
            format!("GOTRUE_JWT_DEFAULT_GROUP_NAME={}", self.jwt_default_group_name),
            format!("GOTRUE_DB_DRIVER={}", self.db_driver),
            format!("DB_NAMESPACE={}", self.db_namespace),
            format!("API_EXTERNAL_URL={}", self.api_external_url),
            format!("GOTRUE_API_HOST={}", self.host),
            format!("PORT={}", self.port),
            format!("GOTRUE_DISABLE_SIGNUP={}", self.disable_signup),
            format!("GOTRUE_SITE_URL={}", self.site_url),
            format!("GOTRUE_MAILER_AUTOCONFIRM={}", self.mailer_auto_confirm),
            format!("GOTRUE_LOG_LEVEL={}", self.log_level),
            format!("GOTRUE_OPERATOR_TOKEN={}", self.operator_token),
            format!("DATABASE_URL={}", self.database_url),
        ]
    }
}