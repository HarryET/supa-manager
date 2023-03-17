defmodule SupaManager.Repo.Migrations.UpdateProjectsAddDatabase do
  use Ecto.Migration

  def change do
    alter table(:projects) do
      add :db_password, :text, null: false
      add :db_username, :text, null: false
      add :db_host, :text, null: false
      add :db_port, :text, null: false
      add :db_status, :string, null: false
    end
  end
end
