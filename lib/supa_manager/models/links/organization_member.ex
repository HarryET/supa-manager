defmodule SupaManager.Models.Links.OrganizationMember do
  use SupaManager.Models.Schema

  schema "organization_members" do
    belongs_to :user, SupaManager.Models.User
    belongs_to :organization, SupaManager.Models.Organization

    timestamps()
  end

  def changeset(member, attrs) do
    member
    |> cast(attrs, [:user_id, :organization_id])
    |> validate_required([:user_id, :organization_id])
    |> cast_assoc(:user)
    |> cast_assoc(:organization)
  end
end
