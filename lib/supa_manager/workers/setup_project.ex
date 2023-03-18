defmodule SupaManager.Workers.SetupProject do
  use Oban.Worker, queue: :default

  require Logger

  import Ecto.Query
  alias SupaManager.Models.Project
  alias SupaManager.Repo

  @impl Oban.Worker
  def perform(%Oban.Job{args: %{"id" => id} = _args}) do
    with %Project{} = proj <- Repo.one(from(p in Project, where: p.id == ^id)) do
      Logger.info("[#{id}] Starting services")

      # Update DB to COMING_UP
      proj = Repo.update!(Project.update_status_changeset(proj, %{status: "COMING_UP"}))

      {:ok, db_pass} = SupaManager.Core.Encryption.decrypt(proj.db_password)

      env = [
        %Kazan.Apis.Core.V1.EnvVar{
          name: "POSTGRES_USER",
          value: proj.db_username
        },
        # TODO make a secret
        %Kazan.Apis.Core.V1.EnvVar{
          name: "POSTGRES_PASSWORD",
          value: db_pass
        }
      ]

      case SupaManager.Core.Kubernetes.Pod.new(proj.id, :postgres, env) do
        {:ok, _pod} ->
          Logger.info("[#{id}] Started postgres")
          # TODO update project status
          proj = Repo.update!(Project.update_db_status_changeset(proj, %{db_status: :ready}))

          SupaManager.Core.Hooks.update_dns(%{
            project_id: proj.id,
            dns_name: SupaManager.Core.Domains.get_db(proj.ref),
            action: :create,
            reason: :new_project
          })

        {:error, error} ->
          Logger.error("Failed to start postgres")
          {:error, error}
      end
    else
      _ -> {:error, "Couldn't find project"}
    end
  end
end
