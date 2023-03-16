defmodule SupaManager.Controllers.Billing do
  use SupaManager, :controller

  @moduledoc """
  Billing controller for SupaManager.

  Currently we discard all billing however, this could be used in future.
  """

  def overdue_invoices(conn, _params) do
    conn
    |> put_status(200)
    |> json([])
  end
end
