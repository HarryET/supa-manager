defmodule SupaManager.Repo.Migrations.CreateProjectsTable do
  use Ecto.Migration

  def change do
    create table :projects, primary_key: false do
      add :id, :string, primary_key: true
      add :name, :string
      add :organization_id, references(:organizations, on_delete: :nothing)

      add :status, :string

      timestamps()
    end

    create unique_index :organizations, :name
  end
end
