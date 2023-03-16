defmodule SupaManager.Router do
  use SupaManager, :router

  import SupaManager.UserAuth

  pipeline :api do
    plug(:accepts, ["json"])
  end

  pipeline :authenticated do
    plug :ensure_user
  end

  scope "/", SupaManager.Controllers do
    pipe_through(:api)

    get("/", Core, :index)

    post("/signup", Users, :signup)

    scope "/auth" do
      # Mimics GoTrue to remove additional dependency
      post("/token", Users, :login)
    end

    scope "/" do
      pipe_through(:authenticated)

      get "/notifications", Notification, :list

      scope "/profile" do
        get("/", Profile, :index)
        get("/permissions", Profile, :permissions)
      end

      get "/stripe/invoices/overdue", Billing, :overdue_invoices
    end

    scope "/telemetry" do
      match(:*, "/*any", Telemetry, :discard)
    end
  end
end
