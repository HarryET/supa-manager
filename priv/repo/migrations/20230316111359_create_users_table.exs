defmodule SupaManager.Repo.Migrations.CreateUsersTable do
  use Ecto.Migration

  def change do
    create table(:users, primary_key: false) do
      add(:id, :uuid, primary_key: true)

      add(:email, :text, null: false)
      add(:password, :text, null: false)

      add(:username, :string)
      add(:first_name, :string)
      add(:last_name, :string)

      timestamps()
    end
  end
end
