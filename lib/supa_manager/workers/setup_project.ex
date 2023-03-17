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

      case SupaManager.Core.Kubernetes.Pod.new(proj.id, :postgres) do
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
