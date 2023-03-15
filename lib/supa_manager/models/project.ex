defmodule SupaManager.Models.Project do
  use SupaManager.Models.Schema

  schema "projects" do
    field(:name, :string)
    field(:status, :string)

    belongs_to(:organization, SupaManager.Models.Organization)

    timestamps()
  end
end
