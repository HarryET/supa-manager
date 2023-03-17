defmodule SupaManager.Repo.Migrations.UpdateProjects do
  use Ecto.Migration

  def change do
    alter table(:projects) do
      add :ref, :string, null: false, unique: true
      add :cloud_provider, :string, null: false, default: "supamanager"
      add :region, :string, null: false
    end
  end
end
