import Config

# API
config :supa_manager,
  port: 4000

# Oban
config :my_app, Oban, testing: :inline

# Ecto
config :supa_manager,
  ecto_repos: [SupaManager.Repo]

config :supa_manager, SupaManager.Repo,
  database: "supa_manager",
  username: "postgres",
  password: "postgres",
  hostname: "localhost"
