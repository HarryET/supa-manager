defmodule SupaManager.Endpoint do
  use Phoenix.Endpoint, otp_app: :supa_manager

  plug(Plug.Logger)

  plug(CORSPlug)

  plug(Plug.Parsers,
    parsers: [:json],
    pass: ["application/json"],
    json_decoder: Jason
  )

  plug(Plug.RequestId)
  plug(Plug.Telemetry, event_prefix: [:phoenix, :endpoint])

  plug(SupaManager.Router)
end
