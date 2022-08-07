if Mix.env() == :prod do
  # TODO convert to env vars
  config :supa_manager, SupaManager.Repo,
    database: "supa_manager",
    username: "postgres",
    password: "postgres",
    hostname: "localhost"
end
