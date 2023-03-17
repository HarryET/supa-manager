defmodule SupaManager.MixProject do
  use Mix.Project

  def project do
    [
      app: :supa_manager,
      version: "0.1.0",
      elixir: "~> 1.14-rc",
      start_permanent: Mix.env() == :prod,
      deps: deps(),
      aliases: aliases()
    ]
  end

  # Run "mix help compile.app" to learn about applications.
  def application do
    [
      extra_applications: [:logger],
      mod: {SupaManager.Application, []}
    ]
  end

  # Run "mix help deps" to learn about dependencies.
  defp deps do
    [
      {:plug_cowboy, "~> 2.6"},
      {:jason, "~> 1.3"},
      {:cors_plug, "~> 3.0"},
      {:phoenix, "~> 1.7"},
      # Database
      {:ecto_sql, "~> 3.0"},
      {:postgrex, ">= 0.0.0"},
      {:uuid, "~> 1.1"},
      # K8S
      {:kazan, git: "https://github.com/HarryET/kazan.git"},
      # Queues
      {:oban, "~> 2.14"},
      # Security
      {:argon2_elixir, "~> 3.0"},
      {:joken, "~> 2.6"},
      {:pbcs, "~> 0.1"}
    ]
  end

  defp aliases do
    [
      setup: ["deps.get", "ecto.setup"],
      "ecto.setup": ["ecto.create", "ecto.migrate", "run priv/repo/seeds.exs"],
      "ecto.reset": ["ecto.drop", "ecto.setup"]
    ]
  end
end
