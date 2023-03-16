defmodule SupaManager.Repo.Migrations.ProjectsOrgsUsersLink do
  use Ecto.Migration

  def change do
    create table(:organization_members, primary_key: false) do
      add :id, :uuid, primary_key: true

      add :user_id, references(:users, type: :uuid)
      add :organization_id, references(:organizations, type: :uuid)

      # TODO roles

      timestamps()
    end

    create unique_index(:organization_members, [:user_id, :organization_id])
  end
end
