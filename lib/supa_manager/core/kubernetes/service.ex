defmodule SupaManager.Core.Kubernetes.Service do
  @spec new(Kazan.Apis.Core.V1.Service.t()) ::
          {:ok, Kazan.Apis.Core.V1.Service.t()} | {:error, any}
  def new(service) do
    case Kazan.Apis.Core.V1.create_namespaced_service(
           service,
           SupaManager.Core.Kubernetes.namespace()
         ) do
      {:ok, req} ->
        case Kazan.run(req) do
          {:ok, %Kazan.Apis.Core.V1.Service{} = pod} ->
            {:ok, pod}

          {:ok, res} ->
            {:error, res}

          {:error, err} ->
            {:error, err}
        end

      {:error, err} ->
        {:error, err}
    end
  end
end
