import Config

# Rest API
config :supa_manager,
  port: 80

# Database
config :supa_manager,
  ecto_repos: [SupaManager.Repo]
