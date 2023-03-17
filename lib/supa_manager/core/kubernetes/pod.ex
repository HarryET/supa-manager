defmodule SupaManager.Core.Kubernetes.Pod do
  @spec new(any, SupaManager.Core.Kubernetes.service(), list(Kazan.Apis.Core.V1.EnvVar)) ::
          {:ok, Kazan.Apis.Core.V1.Pod.t()} | {:error, any}
  def new(project_id, service, env) do
    image = SupaManager.Core.Versions.get_image(service)
    version = SupaManager.Core.Versions.get_version(service)

    {:ok, req} =
      Kazan.Apis.Core.V1.create_namespaced_pod(
        %Kazan.Apis.Core.V1.Pod{
          metadata: %Kazan.Models.Apimachinery.Meta.V1.ObjectMeta{
            name: "#{project_id}-#{service}",
            namespace: "default",
            labels: %{
              "supamanager.io/managed" => "true",
              "supamanager.io/project" => project_id
            }
          },
          spec: %Kazan.Apis.Core.V1.PodSpec{
            containers: [
              %Kazan.Apis.Core.V1.Container{
                name: "#{service}",
                image: "#{image}:#{version}",
                env: env
              }
            ]
          }
        },
        SupaManager.Core.Kubernetes.namespace()
      )

    case Kazan.run(req) do
      {:ok, %Kazan.Apis.Core.V1.Pod{} = pod} ->
        {:ok, pod}

      {:ok, res} ->
        {:error, res}

      {:error, err} ->
        {:error, err}
    end
  end
end
