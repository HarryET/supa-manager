defmodule SupaManager.Models.Project do
  use SupaManager.Models.Schema

  schema "projects" do
    field(:name, :string)
    field(:status, :string)

    field :db_password, :string
    field :db_username, :string
    field :db_host, :string
    field :db_port, :string
    field :db_status, Ecto.Enum, values: [:pending, :ready, :error]

    belongs_to(:organization, SupaManager.Models.Organization)

    timestamps()
  end

  def changeset(project, attrs) do
    project
    |> cast(attrs, [
      :name,
      :status,
      :organization_id,
      :db_password,
      :db_username,
      :db_host,
      :db_port,
      :db_status
    ])
    |> validate_required([
      :name,
      :status,
      :organization_id,
      :db_password,
      :db_username,
      :db_status
    ])
    |> cast_assoc(:organization)
  end
end
