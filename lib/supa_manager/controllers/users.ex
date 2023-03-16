defmodule SupaManager.Controllers.Users do
  use SupaManager, :controller

  def signup(conn, _params) do
    send_resp(conn, 500, "todo: implement")
  end

  # get "/profile" do
  #   conn
  #   |> put_resp_content_type("application/json")
  #   |> send_resp(
  #     200,
  #     Jason.encode!(%{
  #       "id" => 1,
  #       "primary_email" => "admin@supamanager.io",
  #       "username" => "admin",
  #       "first_name" => "Supa",
  #       "last_name" => "Manager",
  #       "organizations" => []
  #     })
  #   )
  # end

  def profile(conn, _params) do
    send_resp(conn, 500, "todo: implement")
  end
end
