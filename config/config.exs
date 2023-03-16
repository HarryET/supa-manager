import Config

config :supa_manager,
  ecto_repos: [SupaManager.Repo]

# Phoenix Endpoint
config :supa_manager, SupaManager.Endpoint,
  url: [host: "localhost"],
  server: true

# CORS
config :cors_plug,
  origin: ["http://localhost:3000"],
  max_age: 86400,
  methods: ["*"],
  headers: ["*"]

# Oban
config :supa_manager, Oban,
  repo: SupaManager.Repo,
  plugins: [],
  queues: [default: 10]

# Logger
config :logger, :console,
  format: "$time $metadata[$level] $message\n",
  metadata: [:request_id]

# Import environment specific config. This must remain at the bottom
# of this file so it overrides the configuration defined above.
import_config "#{config_env()}.exs"
