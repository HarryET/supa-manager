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
          ref: p.id,
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

  def get(conn, %{"id" => id} = _params) do
    with %Project{} = proj <- Repo.one(from p in Project, where: p.id == ^id) do
      conn
      |> put_status(200)
      |> json(%{
        id: proj.id,
        ref: proj.id,
        name: proj.name,
        organization_id: proj.organization_id,
        cloud_provider: "aws",
        status: proj.status,
        region: "eu-central-1",
        inserted_at: proj.inserted_at
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
          "db_region" => _region,
          "name" => name,
          "org_id" => org_id
        } = _params
      ) do
    with encrypted_pass <- SupaManager.Encryption.encrypt(db_pass),
         true <- is_binary(encrypted_pass),
         {:ok, %Project{} = proj} <-
           Repo.insert(
             %Project{}
             |> Project.changeset(%{
               name: name,
               organization_id: org_id,
               status: "UNKNOWN",
               db_password: encrypted_pass,
               db_username: "postgres",
               db_status: :pending
             })
           ) do
      %{id: proj.id}
      |> SupaManager.Workers.SetupProject.new()
      |> Oban.insert()

      conn
      |> put_status(201)
      |> json(%{
        id: proj.id,
        ref: proj.id,
        name: proj.name,
        status: proj.status,
        organization_id: proj.organization_id,
        cloud_provider: "aws",
        region: "eu-central-1",
        inserted_at: proj.inserted_at,
        subscription_id: "sub_not_implemented",
        endpoint: "https://#{proj.id}.#{Application.get_env(:supa_manager, :supabase_host)}",
        # TODO - generate anon_key and service_key
        anon_key: "not_implemented",
        service_key: "not_implemented"
      })
    else
      _ ->
        conn
        |> put_status(500)
        |> json(%{message: "Failed to create project"})
    end
  end
end
