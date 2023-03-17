defmodule SupaManager.Core.Kubernetes do
  @type service :: :postgres

  @namespace "default"

  @spec namespace :: String.t()
  def namespace, do: @namespace
end
