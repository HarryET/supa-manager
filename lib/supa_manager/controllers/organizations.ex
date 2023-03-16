defmodule SupaManager.Controllers.Organizations do
  use SupaManager, :controller

  # get "/organizations" do
  #   import Ecto.Query

  #   orgs = SupaManager.Repo.all(from(o in SupaManager.Models.Organization))

  #   conn
  #   |> put_resp_content_type("application/json")
  #   |> send_resp(
  #     200,
  #     Jason.encode!(
  #       Enum.map(orgs, fn o ->
  #         %{
  #           "id" => o.id,
  #           "slug" => o.slug,
  #           "name" => o.name,
  #           "billing_email" => "admin@supamanager.io"
  #         }
  #       end)
  #     )
  #   )
  # end

  def list(conn, _params) do
    send_resp(conn, 500, "todo: implement")
  end
end
