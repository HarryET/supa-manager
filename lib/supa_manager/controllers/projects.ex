defmodule SupaManager.Controllers.Projects do
  use SupaManager, :controller

  # get "/projects" do
  #   conn
  #   |> put_resp_content_type("application/json")
  #   |> send_resp(
  #     200,
  #     Jason.encode!([
  #       %{
  #         "id" => 1,
  #         "ref" => "mng",
  #         "name" => "Managed Project",
  #         "organization_id" => 1,
  #         "cloud_provider" => "aws",
  #         "status" => "UNKNOWN",
  #         "region" => "eu-central-1"
  #       }
  #     ])
  #   )
  # end

  def list(conn, _params) do
    send_resp(conn, 500, "todo: implement")
  end
end
