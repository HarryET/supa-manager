defmodule SupaManager.Jwt do
  use Joken.Config

  @aud "supamanager.io"

  @impl true
  def token_config do
    default_claims()
    |> add_claim("aud", fn -> @aud end, &(&1 == @aud))
  end
end
