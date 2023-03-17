defmodule SupaManager.Models.Project do
  use SupaManager.Models.Schema

  schema "projects" do
    field(:name, :string)
    field :ref, :string
    field(:status, :string)

    field :db_password, :string
    field :db_username, :string
    field :db_host, :string
    field :db_port, :string
    field :db_status, Ecto.Enum, values: [:pending, :ready, :error]

    field :cloud_provider, :string
    field :region, :string

    belongs_to(:organization, SupaManager.Models.Organization)

    timestamps()
  end

  def changeset(project, attrs) do
    project
    |> cast(attrs, [
      :name,
      :ref,
      :status,
      :organization_id,
      :db_password,
      :db_username,
      :db_host,
      :db_port,
      :db_status,
      :cloud_provider,
      :region
    ])
    |> validate_required([
      :name,
      :ref,
      :status,
      :organization_id,
      :db_password,
      :db_username,
      :db_host,
      :db_port,
      :db_status,
      :region
    ])
    |> cast_assoc(:organization)
  end
end
