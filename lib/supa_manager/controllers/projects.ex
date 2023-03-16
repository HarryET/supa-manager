defmodule SupaManager.Controllers.Projects do
  use SupaManager, :controller

  alias SupaManager.Repo

  def list(conn, _params) do
    user = Repo.preload(conn.assigns[:current_user], organizations: [:projects])

    projects = Enum.flat_map(user.organizations, fn o -> o.projects end)

    conn
    |> put_status(200)
    |> json(
      Enum.map(projects, fn p ->
        %{
          id: p.id,
          ref: p.id,
          name: p.name,
          organization_id: p.organization_id,
          cloud_provider: "aws",
          status: p.status,
          region: "eu-central-1"
        }
      end)
    )
  end
end
