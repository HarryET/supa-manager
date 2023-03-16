defmodule SupaManager.Repo.Migrations.CreateProjectsTable do
  use Ecto.Migration

  def change do
    create table(:projects, primary_key: false) do
      add(:id, :uuid, primary_key: true)
      add(:name, :string)
      add(:organization_id, references(:organizations, on_delete: :nothing, type: :uuid))

      add(:status, :string)

      timestamps()
    end
  end
end
