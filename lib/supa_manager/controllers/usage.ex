defmodule SupaManager.Controllers.Usage do
  use SupaManager, :controller

  import Ecto.Query
  alias SupaManager.Repo
  alias SupaManager.Models.Project

  def for_project(conn, %{"id" => id}) do
    # TODO load services
    with %Project{} = proj <- Repo.one(from p in Project, where: p.id == ^id) do
      if proj.status == "COMING_UP" do
        conn
        |> put_status(400)
        |> json(%{
          message: "Failed to query project's current usage"
        })
      else
        conn
        |> put_status(200)
        |> json(%{
          db_size: %{
            usage: 0,
            limit: 536_870_912,
            cost: 0,
            available_in_plan: true
          },
          db_egress: %{
            usage: 0,
            limit: 2_147_483_648,
            cost: 0,
            available_in_plan: true
          },
          storage_size: %{
            usage: 0,
            limit: 1_073_741_824,
            cost: 0,
            available_in_plan: true
          },
          storage_egress: %{
            usage: 0,
            limit: 2_147_483_648,
            cost: 0,
            available_in_plan: true
          },
          storage_image_render_count: %{
            usage: 0,
            limit: -1,
            cost: 0,
            available_in_plan: false
          },
          monthly_active_users: %{
            usage: 0,
            limit: 50000,
            cost: 0,
            available_in_plan: true
          },
          func_invocations: %{
            usage: 0,
            limit: 500_000,
            cost: 0,
            available_in_plan: true
          },
          func_count: %{
            usage: 2,
            limit: 10,
            cost: 0,
            available_in_plan: true
          },
          disk_volume_size_gb: 8,
          realtime_message_count: %{
            usage: 0,
            limit: 2_000_000,
            cost: 0,
            available_in_plan: true
          },
          realtime_peak_connection: %{
            usage: 0,
            limit: 200,
            cost: 0,
            available_in_plan: true
          }
        })
      end
    else
      _ ->
        conn
        |> put_status(404)
        |> json(%{
          message: "Project not found"
        })
    end
  end
end
