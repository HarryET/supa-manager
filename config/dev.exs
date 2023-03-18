import Config

# Phoenix Endpoint
config :supa_manager, SupaManager.Endpoint,
  http: [ip: {0, 0, 0, 0}, port: 4000],
  check_origin: false,
  code_reloader: true,
  debug_errors: true,
  secret_key_base: "tevXr3e7YAOW1z7V7N6Qi8xZo7xL9+FXyhGh5lPLYjSxGUq+QBKyxx4zpHT38vgm"

# Ecto
config :supa_manager, SupaManager.Repo,
  database: "supa_manager",
  username: "postgres",
  password: "postgres",
  hostname: "localhost"

# Encryption
config :supa_manager, SupaManager.Core.Encryption, password: "password"

# Logger
config :logger, :console, format: "[$level] $message\n"

# Webhooks
config :supa_manager, SupaManager.Core.Hooks,
  url: "http://localhost:4000/dev/hooks",
  secret: "hook-secret"
