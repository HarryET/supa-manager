defmodule SupaManager.Core.Versions do
  @moduledoc """
  Module for providing images and versions for services
  """

  @spec get_image(SupaManager.Core.Kubernetes.service()) :: String.t()
  def get_image(service) do
    case service do
      :postgres -> "supabase/postgres"
      _ -> throw("Unknown service")
    end
  end

  @spec get_version(SupaManager.Core.Kubernetes.service()) :: String.t()
  def get_version(service) do
    case service do
      :postgres -> "15.1.0.55"
      _ -> throw("Unknown service")
    end
  end
end
