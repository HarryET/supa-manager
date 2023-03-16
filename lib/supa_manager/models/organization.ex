defmodule SupaManager.Models.Organization do
  use SupaManager.Models.Schema

  schema "organizations" do
    field(:slug, :string)
    field(:name, :string)
    field(:kind, :string)
    field(:billing_email, :string)

    belongs_to :owner, SupaManager.Models.User

    has_many :projects, SupaManager.Models.Project

    many_to_many :users, SupaManager.Models.User,
      join_through: SupaManager.Models.Links.OrganizationMember

    timestamps()
  end

  def changeset(org, attrs) do
    org
    |> cast(attrs, [:slug, :name, :kind, :billing_email, :owner_id])
    |> validate_required([:slug, :name, :kind, :billing_email, :owner_id])
    |> cast_assoc(:owner)
  end
end
