defmodule SupaManager.Repo.Migrations.AddJwtSecret do
  use Ecto.Migration

  def change do
    alter table(:projects) do
      add :jwt_secret, :text, null: false
    end
  end
end
