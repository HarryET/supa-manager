defmodule SupaManager.Controllers.Organizations do
  use SupaManager, :controller

  alias SupaManager.Models.Organization
  alias SupaManager.Models.Links.OrganizationMember
  alias SupaManager.Repo

  def list(conn, _params) do
    user = Repo.preload(conn.assigns[:current_user], :organizations)

    conn
    |> put_status(200)
    |> json(
      Enum.map(user.organizations, fn o ->
        %{
          "id" => o.id,
          "slug" => o.slug,
          "name" => o.name,
          "billing_email" => o.billing_email
        }
      end)
    )
  end

  def create(conn, %{"kind" => kind, "name" => name} = _params) do
    # TODO: Store kind
    # TODO: Generate slug
    with {:ok, %Organization{} = org} <-
           Repo.insert(
             %Organization{}
             |> Organization.changeset(%{
               "name" => name,
               "slug" => String.downcase(String.split(name, " ") |> Enum.join("-")),
               "owner_id" => conn.assigns[:current_user].id,
               "billing_email" => conn.assigns[:current_user].email,
               "kind" => kind
             })
           ),
         {:ok, _member} <-
           Repo.insert(
             %OrganizationMember{}
             |> OrganizationMember.changeset(%{
               "organization_id" => org.id,
               "user_id" => conn.assigns[:current_user].id
             })
           ) do
      conn
      |> put_status(201)
      |> json(%{
        "id" => org.id,
        "slug" => org.slug,
        "name" => org.name,
        "billing_email" => org.billing_email
      })
    else
      e ->
        IO.inspect(e)

        conn
        |> put_status(500)
        |> json(%{error: "Internal Server Error"})
    end
  end

  def free_tier_limit(conn, _params) do
    conn
    |> put_status(200)
    |> json([])
  end
end
