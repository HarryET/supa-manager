defmodule SupaManager.Repo do
  use Ecto.Repo,
    otp_app: :supa_manager,
    adapter: Ecto.Adapters.Postgres
end
