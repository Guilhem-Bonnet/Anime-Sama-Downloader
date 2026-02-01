import { useEffect, useRef, useCallback } from 'react';

export type EventListener = (data: any) => void;

export function useSSE(url: string, onMessage?: EventListener) {
  const eventSourceRef = useRef<EventSource | null>(null);
  const callbackRef = useRef(onMessage);

  useEffect(() => {
    callbackRef.current = onMessage;
  }, [onMessage]);

  useEffect(() => {
    const eventSource = new EventSource(url);

    eventSource.onmessage = (event) => {
      const data = JSON.parse(event.data);
      callbackRef.current?.(data);
    };

    eventSource.onerror = (error) => {
      console.error('SSE connection error:', error);
      eventSource.close();
    };

    eventSourceRef.current = eventSource;

    return () => {
      eventSource.close();
    };
  }, [url]);

  return {
    close: useCallback(() => {
      eventSourceRef.current?.close();
    }, []),
  };
}
