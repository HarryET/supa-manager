defmodule SupaManager.Repo.Migrations.CreateOrganizationsTable do
  use Ecto.Migration

  def change do
    create table(:organizations, primary_key: false) do
      add(:id, :uuid, primary_key: true)
      add(:slug, :string)
      add(:name, :string)

      timestamps()
    end

    create(unique_index(:organizations, :name))
  end
end
