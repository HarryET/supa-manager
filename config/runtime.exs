import Config

# GoTrue
config :supa_manager, allow_signup: System.get_env("ALLOW_SIGNUP", "true") == "true"

# Deployment
config :supa_manager, supabase_host: System.get_env("SUPABASE_HOST", "localhost:4000")

if Mix.env() == :prod do
  database_url =
    System.get_env("DATABASE_URL") ||
      raise """
      environment variable DATABASE_URL is missing.
      For example: postgresql://USER:PASS@HOST/DATABASE
      """

  maybe_ipv6 = if System.get_env("ECTO_IPV6"), do: [:inet6], else: []

  config :supa_manager, SupaManager.Repo,
    # ssl: true,
    url: database_url,
    pool_size: String.to_integer(System.get_env("ECTO_POOL_SIZE") || "10"),
    socket_options: maybe_ipv6

  # The secret key base is used to sign/encrypt cookies and other secrets.
  # A default value is used in config/dev.exs and config/test.exs but you
  # want to use a different value for prod and you most likely don't want
  # to check this value into version control, so we use an environment
  # variable instead.
  secret_key_base =
    System.get_env("SECRET_KEY_BASE") ||
      raise """
      environment variable SECRET_KEY_BASE is missing.
      You can generate one by calling: mix phx.gen.secret
      """

  host = System.get_env("HOST") || "supamanager.io"
  port = String.to_integer(System.get_env("PORT") || "4000")

  config :supa_manager, SupaManager.Endpoint,
    url: [host: host, port: 443, scheme: "https"],
    http: [
      # Enable IPv6 and bind on all interfaces.
      # Set it to  {0, 0, 0, 0, 0, 0, 0, 1} for local network only access.
      # See the documentation on https://hexdocs.pm/plug_cowboy/Plug.Cowboy.html
      # for details about using IPv6 vs IPv4 and loopback vs public addresses.
      ip: {0, 0, 0, 0, 0, 0, 0, 0},
      port: port
    ],
    secret_key_base: secret_key_base

  # JWTs
  jwt_secret =
    System.get_env("JWT_SECRET") ||
      raise """
      environment variable JWT_SECRET is missing.
      You can generate one by calling: mix phx.gen.secret
      """

  config :joken, default_signer: jwt_secret

  # Encryption
  encryption_key =
    System.get_env("ENCRYPTION_KEY") ||
      raise """
      environment variable ENCRYPTION_KEY is missing.
      You can generate one by calling: mix phx.gen.secret
      """

  config :supa_manager, SupaManager.Encryption, password: encryption_key
end
