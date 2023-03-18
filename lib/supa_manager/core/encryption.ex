defmodule SupaManager.Core.Encryption do
  @tag "SUPAMANAGER"

  @spec encrypt(binary) :: binary | {:error, binary}
  def encrypt(plain_text) do
    protected = %{
      alg: "PBES2-HS512",
      enc: "A256GCM",
      p2c: 4096,
      p2s: :crypto.strong_rand_bytes(32)
    }

    PBCS.encrypt({@tag, plain_text}, protected, password: password())
  end

  @spec decrypt(binary) :: {:ok, binary} | {:error, String.t()} | :error
  def decrypt(cipher_text) do
    PBCS.decrypt({@tag, cipher_text}, password: password())
  end

  defp password do
    Application.get_env(:supa_manager, SupaManager.Core.Encryption)[:password]
  end
end
