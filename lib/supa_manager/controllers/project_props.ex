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

  #   {
  #   "services": [
  #     {
  #         "id": 365303,
  #         "name": "Default API",
  #         "app_config": {
  #             "endpoint": "pxtlksmmatueuqcktmkr.supabase.co",
  #             "db_schema": "public",
  #             "realtime_multitenant_enabled": true
  #         },
  #         "app": {
  #             "id": 1,
  #             "name": "Auto API"
  #         },
  #         "service_api_keys": [
  #             {
  #                 "tags": "anon",
  #                 "name": "anon key",
  #                 "api_key": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6InB4dGxrc21tYXR1ZXVxY2t0bWtyIiwicm9sZSI6ImFub24iLCJpYXQiOjE2NzkwNzQ5OTcsImV4cCI6MTk5NDY1MDk5N30.D3QPEyK2WpgkD5FwT1JYER1iaWvLnYEJlsawv6BJvC4"
  #             },
  #             {
  #                 "tags": "service_role",
  #                 "name": "service_role key",
  #                 "api_key": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6InB4dGxrc21tYXR1ZXVxY2t0bWtyIiwicm9sZSI6InNlcnZpY2Vfcm9sZSIsImlhdCI6MTY3OTA3NDk5NywiZXhwIjoxOTk0NjUwOTk3fQ.skHN78BgUc8KlN2uPlvajjPe3SZvd4ow-SGBXiKj894"
  #             }
  #         ]
  #     }
  # ]
  # }

  def settings(conn, %{"ref" => ref} = _params) do
    with %Project{} = proj <- Repo.one(from p in Project, where: p.ref == ^ref) do
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
          db_dns_name: "",
          db_port: proj.db_port,
          db_host: proj.db_host,
          db_name: proj.db_username,
          db_ssl: false,
          db_user: proj.db_username,
          jwt_secret: "encrypted rubbish"
        },
        # TODO services
        services: [
          %{
            id: proj.id,
            name: "Default API",
            app_config: %{
              endpoint: SupaManager.Core.Domains.get(proj.ref),
              db_schema: "public",
              realtime_multitenant_enabled: true
            },
            app: %{
              id: 1,
              name: "Auto API"
            },
            service_api_keys: [
              %{
                tags: "anon",
                name: "anon key",
                api_key:
                  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6InB4dGxrc21tYXR1ZXVxY2t0bWtyIiwicm9sZSI6ImFub24iLCJpYXQiOjE2NzkwNzQ5OTcsImV4cCI6MTk5NDY1MDk5N30.D3QPEyK2WpgkD5FwT1JYER1iaWvLnYEJlsawv6BJvC4"
              },
              %{
                tags: "service_role",
                name: "service_role key",
                api_key:
                  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6InB4dGxrc21tYXR1ZXVxY2t0bWtyIiwicm9sZSI6InNlcnZpY2Vfcm9sZSIsImlhdCI6MTY3OTA3NDk5NywiZXhwIjoxOTk0NjUwOTk3fQ.skHN78BgUc8KlN2uPlvajjPe3SZvd4ow-SGBXiKj894"
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
