defmodule SupaManager.Workers.SetupProject do
  use Oban.Worker, queue: :default

  require Logger

  import Ecto.Query
  alias SupaManager.Models.Project
  alias SupaManager.Repo

  @impl Oban.Worker
  def perform(%Oban.Job{args: %{"id" => id} = _args}) do
    with %Project{} = proj <- Repo.one(from(p in Project, where: p.id == ^id)) do
      Logger.info("Setup Project (#{id})")

      {:ok, db_pass} = SupaManager.Core.Encryption.decrypt(proj.db_password)

      env = [
        %Kazan.Apis.Core.V1.EnvVar{
          name: "POSTGRES_USER",
          value: proj.db_username
        },
        %Kazan.Apis.Core.V1.EnvVar{
          name: "POSTGRES_PASSWORD",
          value: db_pass
        }
      ]

      case SupaManager.Core.Kubernetes.Pod.new(proj.id, :postgres, env) do
        {:ok, _pod} ->
          Logger.info("Started postgres pod")

        {:error, error} ->
          Logger.error("Failed to start postgres")
          {:error, error}
      end
    else
      _ -> {:error, "Couldn't find project"}
    end
  end
end
