defmodule SupaManager.Repo.Migrations.UpdateProjectsAddDatabase do
  use Ecto.Migration

  def change do
    alter table(:projects) do
      add :db_password, :text, null: false
      add :db_username, :text, null: false
      add :db_host, :text
      add :db_port, :text
      add :db_status, :string, null: false
    end
  end
end
