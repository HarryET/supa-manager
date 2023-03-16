defmodule SupaManager.Models.User do
  use SupaManager.Models.Schema

  schema "users" do
    field(:email, :string)
    field(:password, :string)

    field(:username, :string)
    field(:first_name, :string)
    field(:last_name, :string)

    many_to_many :organizations, SupaManager.Models.Organization,
      join_through: SupaManager.Models.Links.OrganizationMember

    timestamps()
  end

  def changeset(user, attrs) do
    user
    |> cast(attrs, [:email, :password, :username, :first_name, :last_name])
    |> validate_required([:email, :password])
  end
end
