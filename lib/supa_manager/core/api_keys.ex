defmodule SupaManager.Core.ApiKeys do
  def service_key(proj, jwt_secret), do: gen_token(proj, jwt_secret, "service_role")
  def service_key!(proj, jwt_secret), do: gen_token!(proj, jwt_secret, "service_role")

  def anon_key(proj, jwt_secret), do: gen_token(proj, jwt_secret, "anon")
  def anon_key!(proj, jwt_secret), do: gen_token!(proj, jwt_secret, "anon")

  @spec gen_token(any(), any(), any()) :: {:ok, binary()} | {:error, any()}
  def gen_token(proj, jwt_secret, role) do
    token_cfg =
      Joken.Config.default_claims(
        iss: "supabase",
        aud: "supamanager.io",
        # Note: exp doesn't seam important for anon tokens
        default_exp: 24 * 60 * 60
      )

    case Joken.generate_and_sign(
           token_cfg,
           %{
             "role" => role,
             "ref" => proj.ref,
             "iat" => DateTime.to_unix(DateTime.from_naive!(proj.inserted_at, "Etc/UTC"))
           },
           Joken.Signer.create("HS256", jwt_secret)
         ) do
      {:ok, token, _claims} -> {:ok, token}
      {:error, err} -> {:error, err}
    end
  end

  @spec gen_token!(any(), any(), any()) :: binary()
  def gen_token!(proj, jwt_secret, role) do
    case gen_token(proj, jwt_secret, role) do
      {:ok, token} -> token
      {:error, err} -> raise err
    end
  end
end
