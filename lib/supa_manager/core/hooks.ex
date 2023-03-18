defmodule SupaManager.Core.Hooks do
  @moduledoc """
  Module to send Webhooks to a provided URL to allow for customization of actions e.g. update DNS
  """

  require Logger

  @type dns_req :: %{
          project_id: String.t(),
          dns_name: String.t(),
          action: :create | :update | :delete,
          reason: :new_project | :project_update | :project_delete
        }

  @spec update_dns(dns_req) :: :ok | {:error, String.t()}
  def update_dns(req) do
    Logger.debug("[Hooks] Updating DNS for #{req.project_id} on domain #{req.dns_name}")

    # TODO send webhook with secret

    :ok
  end
end
