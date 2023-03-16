defmodule SupaManager.Controllers.Profile do
  use SupaManager, :controller

  def index(conn, _params) do
    user = conn.assigns[:current_user]

    conn
    |> put_status(200)
    |> json(%{
      id: user.id,
      primary_email: user.email,
      username: user.username,
      first_name: user.first_name,
      last_name: user.last_name
    })
  end

  def permissions(conn, _params) do
    conn
    |> put_status(200)
    |> json([])
  end
end
