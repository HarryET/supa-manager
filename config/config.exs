import Config

# CORS
config :cors_plug,
  origin: ["http://localhost:8082"],
  max_age: 86400,
  methods: ["*"],
  headers: ["*"]

# Oban
config :my_app, Oban,
  repo: SupaManager.Repo,
  plugins: [],
  queues: [default: 10]

# Import environment specific config. This must remain at the bottom
# of this file so it overrides the configuration defined above.
import_config "#{config_env()}.exs"
