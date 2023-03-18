defmodule SupaManager.Jwt do
  use Joken.Config

  @aud "supamanager.io"

  @impl true
  def token_config do
    default_claims(iss: @aud, aud: @aud)
  end
end
