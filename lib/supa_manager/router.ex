defmodule SupaManager.Router do
  use SupaManager, :router

  pipeline :api do
    plug(:accepts, ["json"])
  end

  scope "/", SupaManager.Controllers do
    pipe_through(:api)

    get("/", Core, :index)

    post("/signup", Users, :signup)

    scope "/telemetry" do
      match(:*, "/*any", Telemetry, :discard)
    end
  end
end
