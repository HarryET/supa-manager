defmodule SupaManager.UserAuth do
  import Plug.Conn
  import Ecto.Query

  alias SupaManager.Models.User
  alias SupaManager.Repo

  def ensure_user(conn, _opts) do
    with user_token = find_user_token(conn),
         true <- is_binary(user_token),
         {:ok, claims} <-
           SupaManager.Jwt.verify_and_validate(user_token),
         %User{} = user <-
           Repo.one(from u in User, where: u.id == ^claims["sub"]) do
      assign(
        conn,
        :current_user,
        user
      )
    else
      e ->
        IO.inspect(e)

        conn
        |> put_resp_content_type("application/json")
        |> send_resp(
          401,
          Jason.encode!(%{
            "error" => "Unauthorized"
          })
        )
        |> halt()
    end
  end

  def find_user_token(conn) do
    with ["Bearer " <> token] <- get_req_header(conn, "authorization"),
         true <- is_binary(token) do
      token
    else
      _ -> nil
    end
  end
end
