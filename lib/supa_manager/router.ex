defmodule SupaManager.Router do
  # Bring Plug.Router module into scope
  use Plug.Router

  # Attach the Logger to log incoming requests
  plug(Plug.Logger)

  # CORS
  plug(CORSPlug)

  # Tell Plug to match the incoming request with the defined endpoints
  plug(:match)

  # Once there is a match, parse the response body if the content-type
  # is application/json. The order is important here, as we only want to
  # parse the body if there is a matching route.(Using the Jayson parser)
  plug(Plug.Parsers,
    parsers: [:json],
    pass: ["application/json"],
    json_decoder: Jason
  )

  # Dispatch the connection to the matched handler
  plug(:dispatch)

  # Handler for GET request with "/" path
  get "/" do
    send_resp(conn, 200, "OK")
  end

  get "/profile" do
    conn
    |> put_resp_content_type("application/json")
    |> send_resp(200, Jason.encode!(%{
      "id" => 1,
      "primary_email" => "admin@supamanager.io",
      "username" => "admin",
      "first_name" => "Supa",
      "last_name" => "Manager",
      "organizations" => []
    }))
  end

  get "/projects" do
    conn
    |> put_resp_content_type("application/json")
    |> send_resp(200, Jason.encode!([
      %{
        "id" => 1,
        "ref" => "mng",
        "name" => "Managed Project",
        "organization_id" => 1,
        "cloud_provider" => "aws",
        "status" => "UNKNOWN",
        "region" => "eu-central-1"
      }
    ]))
  end

  get "/organizations" do
    import Ecto.Query

    orgs = SupaManager.Repo.all(from o in SupaManager.Models.Organization)

    conn
    |> put_resp_content_type("application/json")
    |> send_resp(200, Jason.encode!(Enum.map(orgs, fn o -> %{
      "id" => o.slug,
      "slug" => o.slug,
      "name" => o.name,
      "billing_email" => "admin@supamanager.io"
    } end)))
  end

  #? Ignore Telemetry
  post "/telemetry/:event" do
    conn
    |> send_resp(200, "OK")
  end

  # Fallback handler when there was no match
  match _ do
    send_resp(conn, 404, "Not Found")
  end
end
