defmodule SupaManager.Controllers.Profile do
  use SupaManager, :controller

  alias SupaManager.Repo

  def index(conn, _params) do
    user = conn.assigns[:current_user]

    conn
    |> put_status(200)
    |> json(%{
      id: user.id,
      primary_email: user.email,
      username: user.username,
      first_name: user.first_name,
      last_name: user.last_name
    })
  end

  def password_check(conn, %{"password" => _password} = _params) do
    # TODO: score password

    conn
    |> put_status(200)
    |> json(%{
      result: %{
        score: 4
      }
    })
  end

  # TODO actually use permissions
  def permissions(conn, _params) do
    user = Repo.preload(conn.assigns[:current_user], :organizations)

    conn
    |> put_status(200)
    |> json(
      Enum.map(user.organizations, fn o ->
        %{
          organization_id: o.id,
          resources: [
            "%"
          ],
          actions: [
            "billing:Write",
            "infra:Execute",
            "read:Read",
            "write:Create",
            "write:Update"
          ],
          condition: nil
        }
      end)
    )
  end
end
