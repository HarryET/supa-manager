defmodule SupaManager.Repo.Migrations.UpdateOrgs do
  use Ecto.Migration

  def change do
    alter table(:organizations) do
      add :kind, :text, null: false
      add :owner_id, references(:users, type: :uuid), null: false
      add :billing_email, :text, null: false
    end
  end
end
