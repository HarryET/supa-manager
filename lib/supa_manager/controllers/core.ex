defmodule SupaManager.Controllers.Core do
  use SupaManager, :controller

  def index(conn, _params) do
    conn
    |> put_status(200)
    |> json(%{status: "OK"})
  end

  def not_found(conn, _params) do
    conn
    |> put_status(404)
    |> json(%{error: "Not Found"})
  end
end
