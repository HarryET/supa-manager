defmodule SupaManager.Controllers.Telemetry do
  use SupaManager, :controller

  @moduledoc """
  Telemetry controller for SupaManager.

  Currently we discard all telemetry however, this could be stored in future behind an optional flag.
  """

  def discard(conn, _params) do
    send_resp(conn, 200, "telemetry discarded")
  end
end
