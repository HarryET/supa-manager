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
        post("/password-check", Profile, :password_check)
      end

      scope "/organizations" do
        get("/", Organizations, :list)
        post("/", Organizations, :create)

        scope "/:slug" do
          scope "/members" do
            get("/reached-free-project-limit", Organizations, :free_tier_limit)
          end

          get "/payments", Billing, :org_payments
        end
      end

      scope "/projects" do
        get("/", Projects, :list)
        post("/", Projects, :create)

        scope "/:id" do
          get("/", Projects, :get)
          get("/usage", Usage, :for_project)
          get("/subscription", Billing, :project_subscription)
        end
      end

      get "/stripe/invoices/overdue", Billing, :overdue_invoices
    end

    scope "/telemetry" do
      match(:*, "/*any", Telemetry, :discard)
    end
  end
end
