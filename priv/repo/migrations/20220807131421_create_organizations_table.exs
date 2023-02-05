defmodule SupaManager.Repo.Migrations.CreateOrganizationsTable do
  use Ecto.Migration

  def change do
    create table :organizations, primary_key: false do
      add :slug, :string, primary_key: true
      add :name, :string

      timestamps()
    end

    create unique_index :organizations, :name
  end
end
