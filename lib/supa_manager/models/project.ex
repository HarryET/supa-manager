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

    field :jwt_secret, :string

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
      :region,
      :jwt_secret
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
      :region,
      :jwt_secret
    ])
    |> cast_assoc(:organization)
  end

  def update_status_changeset(project, attrs) do
    project
    |> cast(attrs, [:status])
    |> validate_required([:status])
  end

  def update_db_status_changeset(project, attrs) do
    project
    |> cast(attrs, [:db_status])
    |> validate_required([:db_status])
  end
end
