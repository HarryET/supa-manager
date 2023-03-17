defmodule SupaManager.Core.ObanLogger do
  require Logger

  def handle_event(
        [:oban, :job, :exception],
        %{duration: duration, queue_time: _queue_time},
        %{queue: queue, worker: worker, stacktrace: stacktrace} = _meta,
        nil
      ) do
    Logger.warn("[#{queue}] #{worker} failed in #{duration}")
    Logger.error(Exception.format_stacktrace(stacktrace))
  end
end
