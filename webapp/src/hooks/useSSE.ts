import { useEffect, useRef, useCallback } from 'react';

export type SSEMessageHandler = (data: any) => void;

const MAX_RETRY_DELAY = 30_000;
const INITIAL_RETRY_DELAY = 1_000;

export function useSSE(url: string, onMessage?: SSEMessageHandler, eventName: string | string[] = 'message') {
  const eventSourceRef = useRef<EventSource | null>(null);
  const callbackRef = useRef(onMessage);
  const retryDelayRef = useRef(INITIAL_RETRY_DELAY);
  const retryTimerRef = useRef<ReturnType<typeof setTimeout> | null>(null);

  // Normalize to array for uniform handling
  const eventNames = Array.isArray(eventName) ? eventName : [eventName];
  // Stable dep key derived from the sorted list
  const eventNamesKey = [...eventNames].sort().join(',');

  useEffect(() => {
    callbackRef.current = onMessage;
  }, [onMessage]);

  useEffect(() => {
    if (!url) {
      return;
    }

    let cancelled = false;

    function connect() {
      if (cancelled) return;

      const eventSource = new EventSource(url);

      const handleMessage = (event: MessageEvent) => {
        try {
          const data = JSON.parse(event.data);
          callbackRef.current?.(data);
        } catch {
          // Ignore non-JSON messages (heartbeats, etc.)
        }
      };

      eventSource.onopen = () => {
        // Reset retry delay on successful connection
        retryDelayRef.current = INITIAL_RETRY_DELAY;
      };

      // Register a listener for every requested event name
      for (const name of eventNames) {
        if (name === 'message') {
          eventSource.onmessage = handleMessage;
        } else {
          eventSource.addEventListener(name, handleMessage as EventListener);
        }
      }

      eventSource.onerror = () => {
        eventSource.close();
        if (!cancelled) {
          // Exponential backoff reconnection
          const delay = retryDelayRef.current;
          retryDelayRef.current = Math.min(delay * 2, MAX_RETRY_DELAY);
          retryTimerRef.current = setTimeout(connect, delay);
        }
      };

      eventSourceRef.current = eventSource;
    }

    connect();

    return () => {
      cancelled = true;
      if (retryTimerRef.current) clearTimeout(retryTimerRef.current);
      if (eventSourceRef.current) {
        eventSourceRef.current.close();
      }
    };
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [url, eventNamesKey]);

  return {
    close: useCallback(() => {
      if (retryTimerRef.current) clearTimeout(retryTimerRef.current);
      eventSourceRef.current?.close();
    }, []),
  };
}
