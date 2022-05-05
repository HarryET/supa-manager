use rocket_sync_db_pools::{diesel, database};

#[database("supa_admin_db")]
pub struct PostgresDbConn(diesel::PgConnection);