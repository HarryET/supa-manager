defmodule SupaManager.Core.Domains do
  def base, do: Application.get_env(:supa_manager, :supabase_host)

  @spec get(String.t()) :: String.t()
  def get(ref) do
    "#{ref}.#{base()}"
  end

  @spec get_db(String.t()) :: String.t()
  def get_db(ref) do
    "db.#{get(ref)}"
  end
end
