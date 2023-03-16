defmodule SupaManager do
  @moduledoc """
  Manage self-hosted Supabase instances using the Supabase studio.
  """

  # A Project by Harry Bairstow

  def router do
    quote do
      use Phoenix.Router, helpers: false
      import Plug.Conn
      import Phoenix.Controller
    end
  end

  def controller do
    quote do
      use Phoenix.Controller,
        namespace: SupaManager,
        formats: [:json]

      import Plug.Conn
    end
  end

  defmacro __using__(which) when is_atom(which) do
    apply(__MODULE__, which, [])
  end
end
