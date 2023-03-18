defmodule SupaManager.Controllers.Hooks.Dev do
  use SupaManager, :controller

  require Logger

  def handle_hook(conn, %{"project_id" => project_id} = _params) do
    Logger.info("[Hooks_Dev] Received webhook for project #{project_id}")

    conn
    |> put_status(200)
    |> json(%{status: "OK"})
  end
end
