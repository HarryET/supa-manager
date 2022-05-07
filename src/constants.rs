// TODO make dynamic with docker api

pub const STUDIO_IMAGE: &str = "supabase/studio";
pub const STUDIO_TAG: &str = "latest";

pub const DB_IMAGE: &str = "supabase/postgres";
pub const DB_TAG: &str = "latest";

pub const GOTRUE_IMAGE: &str = "supabase/gotrue";
pub const GOTRUE_TAG: &str = "latest";

pub const POSTGREST_IMAGE: &str = "postgrest/postgrest";
pub const POSTGREST_TAG: &str = "latest";

pub const REALTIME_IMAGE: &str = "supabase/realtime";
pub const REALTIME_TAG: &str = "latest";

pub const META_IMAGE: &str = "supabase/postgres-meta";
pub const META_TAG: &str = "latest";

// NOTE: THIS WILL ALL BE ENV CONFIG!
pub const DOMAIN: &str = "localhost";
pub const DB_MIGRATIONS_DIR: &str = "B:\\Development\\supa-manager\\_docker\\pg\\docker-entrypoint-initdb.d";
pub const PASSWORD: &str = "password";