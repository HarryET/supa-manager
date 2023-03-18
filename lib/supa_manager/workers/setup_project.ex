defmodule SupaManager.Workers.SetupProject do
  use Oban.Worker, queue: :default

  require Logger

  import Ecto.Query
  alias SupaManager.Models.Project
  alias SupaManager.Repo

  @impl Oban.Worker
  def perform(%Oban.Job{args: %{"id" => id} = _args}) do
    with %Project{} = proj <- Repo.one(from(p in Project, where: p.id == ^id)) do
      Logger.info("[#{id}] Starting services")

      # Update DB to COMING_UP
      proj = Repo.update!(Project.update_status_changeset(proj, %{status: "COMING_UP"}))

      SupaManager.Core.Hooks.update_dns(%{
        project_id: proj.id,
        dns_name: SupaManager.Core.Domains.get_db(proj.ref),
        action: :create,
        reason: :new_project
      })

      SupaManager.Core.Hooks.update_dns(%{
        project_id: proj.id,
        dns_name: SupaManager.Core.Domains.get(proj.ref),
        action: :create,
        reason: :new_project
      })

      http_service(proj)
      postgres_pod(proj)
      postgres_service(proj)
    else
      _ -> {:error, "Couldn't find project"}
    end
  end

  @spec http_service(any) :: :ok | {:error, binary()}
  def http_service(proj) do
    case SupaManager.Core.Kubernetes.Service.new(%Kazan.Apis.Core.V1.Service{
           metadata: %Kazan.Models.Apimachinery.Meta.V1.ObjectMeta{
             name: "p-#{proj.ref}-http-service",
             namespace: "default",
             labels: %{
               "supamanager.io/managed" => "true",
               "supamanager.io/service" => "http-service",
               "supamanager.io/project" => proj.id,
               "supamanager.io/project-ref" => proj.ref,
               "supamanager.io/region" => proj.region
             }
           },
           spec: %Kazan.Apis.Core.V1.ServiceSpec{
             selector: %{
               "supamanager.io/service" => "http-proxy"
             },
             ports: [
               %Kazan.Apis.Core.V1.ServicePort{
                 name: "http",
                 port: 80,
                 target_port: 80
               }
             ]
           }
         }) do
      {:ok, _service} ->
        Logger.info("[#{proj.id}] Created HTTP service")
        :ok

      {:error, error} ->
        Logger.error("Failed create HTTP service")
        {:error, error}
    end
  end

  @spec postgres_service(any) :: :ok | {:error, binary()}
  def postgres_service(proj) do
    case SupaManager.Core.Kubernetes.Service.new(%Kazan.Apis.Core.V1.Service{
           metadata: %Kazan.Models.Apimachinery.Meta.V1.ObjectMeta{
             name: "p-#{proj.ref}-pg-service",
             namespace: "default",
             labels: %{
               "supamanager.io/managed" => "true",
               "supamanager.io/service" => "pg-service",
               "supamanager.io/project" => proj.id,
               "supamanager.io/project-ref" => proj.ref,
               "supamanager.io/region" => proj.region
             }
           },
           spec: %Kazan.Apis.Core.V1.ServiceSpec{
             selector: %{
               "supamanager.io/service" => "postgres"
             },
             ports: [
               %Kazan.Apis.Core.V1.ServicePort{
                 name: "pg",
                 port: 5432,
                 target_port: 5432
               }
             ]
           }
         }) do
      {:ok, _service} ->
        Logger.info("[#{proj.id}] Created Postgres service")
        :ok

      {:error, error} ->
        Logger.error("Failed create Postgres service")
        {:error, error}
    end
  end

  @spec postgres_pod(any) :: :ok | {:error, binary()}
  def postgres_pod(proj) do
    {:ok, db_pass} = SupaManager.Core.Encryption.decrypt(proj.db_password)

    env = [
      %Kazan.Apis.Core.V1.EnvVar{
        name: "POSTGRES_USER",
        value: proj.db_username
      },
      # TODO make a secret
      %Kazan.Apis.Core.V1.EnvVar{
        name: "POSTGRES_PASSWORD",
        value: db_pass
      }
    ]

    case SupaManager.Core.Kubernetes.Pod.new(%Kazan.Apis.Core.V1.Pod{
           metadata: %Kazan.Models.Apimachinery.Meta.V1.ObjectMeta{
             name: "#{proj.id}-postgres",
             namespace: "default",
             labels: %{
               "supamanager.io/managed" => "true",
               "supamanager.io/service" => "postgres",
               "supamanager.io/project" => proj.id,
               "supamanager.io/project-ref" => proj.ref,
               "supamanager.io/region" => proj.region
             }
           },
           spec: %Kazan.Apis.Core.V1.PodSpec{
             containers: [
               %Kazan.Apis.Core.V1.Container{
                 name: "postgres",
                 image:
                   "#{SupaManager.Core.Versions.get_image(:postgres)}:#{SupaManager.Core.Versions.get_version(:postgres)}",
                 env: env
               }
             ]
           }
         }) do
      {:ok, _pod} ->
        Logger.info("[#{proj.id}] Started postgres")
        Repo.update!(Project.update_db_status_changeset(proj, %{db_status: :ready}))
        :ok

      {:error, error} ->
        Logger.error("Failed to start postgres")
        {:error, error}

      _ ->
        throw("Unknown error creating postgres pod")
    end
  end
end
