defmodule SupaManager.Models.Organization do
  use Ecto.Schema
  @primary_key {:slug, :string, autogenerate: {SupaManager.Models.Utils, :gen_id, []}}

  schema "organizations" do
    field :name, :string

    timestamps()
  end
end
