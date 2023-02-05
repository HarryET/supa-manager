defmodule SupaManager.Models.Utils do
  @spec gen_id :: String.t()
  def gen_id do
    UUID.uuid4(:hex)
  end
end
