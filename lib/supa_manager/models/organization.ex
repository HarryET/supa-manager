defmodule SupaManager.Models.Organization do
  use SupaManager.Models.Schema

  schema "organizations" do
    field(:slug, :string)
    field(:name, :string)

    timestamps()
  end
end
