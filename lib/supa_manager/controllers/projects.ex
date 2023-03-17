defmodule SupaManager.Controllers.Projects do
  use SupaManager, :controller

  import Ecto.Query
  alias SupaManager.Repo
  alias SupaManager.Models.Project

  def list(conn, _params) do
    user = Repo.preload(conn.assigns[:current_user], organizations: [:projects])

    projects = Enum.flat_map(user.organizations, fn o -> o.projects end)

    conn
    |> put_status(200)
    |> json(
      Enum.map(projects, fn p ->
        %{
          id: p.id,
          ref: p.ref,
          name: p.name,
          organization_id: p.organization_id,
          cloud_provider: "aws",
          status: p.status,
          region: "eu-central-1",
          inserted_at: p.inserted_at
        }
      end)
    )
  end

  def get(conn, %{"ref" => ref} = _params) do
    with %Project{} = proj <- Repo.one(from p in Project, where: p.ref == ^ref) do
      conn
      |> put_status(200)
      |> json(%{
        id: proj.id,
        ref: proj.ref,
        name: proj.name,
        organization_id: proj.organization_id,
        cloud_provider: proj.cloud_provider,
        status: proj.status,
        region: proj.region,
        inserted_at: proj.inserted_at,
        subscription_id: "not-implemented",
        db_host: proj.db_host,
        restUrl: SupaManager.Core.Domains.get(proj.ref) <> "/rest/v1/",
        connectionString: "encrypted rubbish",
        # TODO store db version in project row
        dbVersion: "supabase-postgres-#{SupaManager.Core.Versions.get_version(:postgres)}",
        kpsVersion: "supabase-postgres-#{SupaManager.Core.Versions.get_version(:postgres)}"
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

  def status(conn, %{"ref" => ref} = _params) do
    with %Project{} = proj <- Repo.one(from p in Project, where: p.ref == ^ref) do
      conn
      |> put_status(200)
      |> json(%{
        status: proj.status
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

  def create(
        conn,
        %{
          "cloud_provider" => _provider,
          "db_pass" => db_pass,
          "db_pricing_tier_id" => _tier,
          "db_region" => region,
          "name" => name,
          "org_id" => org_id
        } = _params
      ) do
    {:ok, ref_i} = Snowflake.next_id()
    ref = "#{ref_i}"

    with encrypted_pass <- SupaManager.Core.Encryption.encrypt(db_pass),
         true <- is_binary(encrypted_pass),
         {:ok, %Project{} = proj} <-
           Repo.insert(
             %Project{}
             |> Project.changeset(%{
               name: name,
               ref: ref,
               organization_id: org_id,
               status: "COMING_UP",
               db_password: encrypted_pass,
               db_username: "postgres",
               db_host: SupaManager.Core.Domains.get_db(ref),
               db_port: "5432",
               db_status: :pending,
               region: region
             })
           ) do
      %{id: proj.id}
      |> SupaManager.Workers.SetupProject.new()
      |> Oban.insert()

      conn
      |> put_status(201)
      |> json(%{
        id: proj.id,
        ref: proj.ref,
        name: proj.name,
        status: proj.status,
        organization_id: proj.organization_id,
        cloud_provider: "supamanager",
        region: "self-host",
        inserted_at: proj.inserted_at,
        subscription_id: "sub_not_implemented",
        endpoint: "https://#{SupaManager.Core.Domains.get(proj.ref)}",
        # TODO - generate anon_key and service_key
        anon_key: "not_implemented",
        service_key: "not_implemented"
      })
    else
      e ->
        IO.inspect(e)

        conn
        |> put_status(500)
        |> json(%{message: "Failed to create project"})
    end
  end
end
