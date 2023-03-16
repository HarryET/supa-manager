defmodule SupaManager.Controllers.Notification do
  use SupaManager, :controller

  def list(conn, _params) do
    conn
    |> put_status(200)
    |> json([])
  end
end
