defmodule SupaManager.Workers.SetupProject do
  use Oban.Worker, queue: :default

  require Logger

  import Ecto.Query
  alias SupaManager.Models.Project
  alias SupaManager.Repo

  @impl Oban.Worker
  def perform(%Oban.Job{args: %{"id" => id} = _args}) do
    with %Project{} <- Repo.one(from p in Project, where: p.id == ^id) do
      # TODO setup all services
      Logger.info("Setup Project (#{id})")
      :ok
    else
      _ -> {:error, "Couldn't find project"}
    end
  end
end
