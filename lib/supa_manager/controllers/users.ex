defmodule SupaManager.Controllers.Users do
  use SupaManager, :controller

  import Ecto.Query
  alias SupaManager.Models.User
  alias SupaManager.Repo

  def signup(conn, %{"email" => email, "password" => password} = _params) do
    if Application.get_env(:supa_manager, :allow_signup) == false do
      send_resp(conn, 403, "Signup is disabled")
    else
      with nil <- Repo.one(from u in User, where: u.email == ^email) do
        with password_hash <- Argon2.hash_pwd_salt(password),
             {:ok, %User{} = _user} <-
               Repo.insert(
                 %User{}
                 |> User.changeset(%{"email" => email, "password" => password_hash})
               ) do
          send_resp(conn, 200, "User created")
        else
          _ -> send_resp(conn, 500, "Failed to create user")
        end
      else
        _ -> send_resp(conn, 400, "Email already exists")
      end
    end
  end

  def login(
        conn,
        %{"email" => email, "password" => password, "grant_type" => "password"} = _params
      ) do
    with %User{} = user <- Repo.one(from u in User, where: u.email == ^email),
         true <- Argon2.verify_pass(password, user.password) do
      with {:ok, token, _claims} <- SupaManager.Jwt.generate_and_sign(%{"sub" => user.id}),
           {:ok, refresh_token, _claims} <- SupaManager.Jwt.generate_and_sign(%{"sub" => user.id}) do
        conn
        |> put_status(200)
        |> json(%{
          access_token: token,
          token_type: "Bearer",
          expires_in: 3600,
          # TODO: Support refresh tokens
          refresh_token: refresh_token,
          # TODO: Look at GoTrue model for more fields
          # https://github.com/supabase/gotrue/blob/aaf2765d84ec4f56582f2ba550f8015c79c2bf22/internal/models/user.go#L19-L70
          user: %{
            id: user.id,
            email: user.email,
            app_metadata: %{
              provider: "email"
            }
          }
        })
      else
        _ ->
          conn
          |> put_status(500)
          |> json(%{message: "Failed to generate token"})
      end
    else
      _ ->
        conn
        |> put_status(401)
        |> json(%{message: "Bad Request: Invalid Login"})
    end
  end

  def login(conn, _params) do
    conn
    |> put_status(400)
    |> json(%{message: "Bad Request: Must be password grant login"})
  end
end
