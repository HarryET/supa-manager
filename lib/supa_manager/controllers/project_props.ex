defmodule SupaManager.Controllers.ProjectProps do
  use SupaManager, :controller

  import Ecto.Query
  alias SupaManager.Repo
  alias SupaManager.Models.Project

  def jwt_secret_update_status(conn, %{"ref" => ref} = _params) do
    with %Project{} = _proj <- Repo.one(from p in Project, where: p.ref == ^ref) do
      conn
      |> put_status(200)
      |> json(%{
        jwtSecretUpdateStatus: nil
      })
    else
      _ ->
        conn
        |> put_status(404)
        |> json(%{
          message: "Project not found"
        })
    end
  end

  def settings(conn, %{"ref" => ref} = _params) do
    with %Project{} = proj <- Repo.one(from p in Project, where: p.ref == ^ref),
         {:ok, jwt_secret} = SupaManager.Core.Encryption.decrypt(proj.jwt_secret) do
      conn
      |> put_status(200)
      |> json(%{
        project: %{
          id: proj.id,
          ref: proj.ref,
          name: proj.name,
          status: proj.status,
          cloud_provider: proj.cloud_provider,
          region: proj.region,
          inserted_at: proj.inserted_at,
          updated_at: proj.updated_at,
          db_dns_name: "",
          db_port: proj.db_port,
          db_host: proj.db_host,
          db_name: proj.db_username,
          db_ssl: false,
          db_user: proj.db_username,
          jwt_secret: jwt_secret
        },
        services: [
          %{
            id: proj.id,
            name: "Default API",
            app_config: %{
              endpoint: SupaManager.Core.Domains.get(proj.ref),
              # TODO store in db
              db_schema: "public",
              # TODO investigate
              realtime_multitenant_enabled: false
            },
            app: %{
              id: 1,
              name: "Auto API"
            },
            service_api_keys: [
              %{
                tags: "anon",
                name: "anon key",
                api_key: SupaManager.Core.ApiKeys.anon_key(proj, jwt_secret)
              },
              %{
                tags: "service_role",
                name: "service_role key",
                api_key: SupaManager.Core.ApiKeys.service_key(proj, jwt_secret)
              }
            ]
          }
        ]
      })
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
