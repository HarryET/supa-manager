table! {
    instances (id) {
        id -> Uuid,
        nickname -> Nullable<Varchar>,
        hostname -> Varchar,
        studio_container_id -> Nullable<Text>,
        studio_enabled -> Bool,
        database_container_id -> Text,
        gotrue_container_id -> Text,
        realtime_container_id -> Text,
        postgrest_container_id -> Text,
        postgres_meta_container_id -> Text,
    }
}
