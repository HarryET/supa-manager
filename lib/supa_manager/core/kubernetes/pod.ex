defmodule SupaManager.Core.Kubernetes.Pod do
  @spec new(Kazan.Apis.Core.V1.Pod.t()) ::
          {:ok, Kazan.Apis.Core.V1.Pod.t()} | {:error, any}
  def new(pod) do
    case Kazan.Apis.Core.V1.create_namespaced_pod(
           pod,
           SupaManager.Core.Kubernetes.namespace()
         ) do
      {:ok, req} ->
        case Kazan.run(req) do
          {:ok, %Kazan.Apis.Core.V1.Pod{} = pod} ->
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
