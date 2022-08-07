import Config

config :supa_manager, SupaManager.Repo,
  database: "supa_manager_repo",
  username: "user",
  password: "pass",
  hostname: "localhost"

# Import environment specific config. This must remain at the bottom
# of this file so it overrides the configuration defined above.
import_config "#{config_env()}.exs"
