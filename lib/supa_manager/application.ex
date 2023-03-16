defmodule SupaManager.Application do
  # See https://hexdocs.pm/elixir/Application.html
  # for more information on OTP Applications
  @moduledoc false

  use Application

  @impl true
  def start(_type, _args) do
    children = [
      SupaManager.Endpoint,
      SupaManager.Repo,
      {Oban, Application.fetch_env!(:supa_manager, Oban)}
    ]

    # See https://hexdocs.pm/elixir/Supervisor.html
    # for other strategies and supported options
    opts = [strategy: :one_for_one, name: SupaManager.Supervisor]
    Supervisor.start_link(children, opts)
  end
end
