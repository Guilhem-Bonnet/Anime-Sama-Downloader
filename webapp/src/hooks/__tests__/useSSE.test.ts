import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { renderHook, act } from '@testing-library/react';
import { useSSE } from '../useSSE';

// Mock EventSource
class MockEventSource {
  static instances: MockEventSource[] = [];
  url: string;
  onopen: (() => void) | null = null;
  onmessage: ((event: MessageEvent) => void) | null = null;
  onerror: (() => void) | null = null;
  readyState = 0;
  private listeners = new Map<string, EventListener[]>();

  constructor(url: string) {
    this.url = url;
    MockEventSource.instances.push(this);
  }

  addEventListener(event: string, listener: EventListener) {
    const existing = this.listeners.get(event) || [];
    this.listeners.set(event, [...existing, listener]);
  }

  removeEventListener(event: string, listener: EventListener) {
    const existing = this.listeners.get(event) || [];
    this.listeners.set(event, existing.filter((l) => l !== listener));
  }

  close = vi.fn();

  // Test helpers
  simulateOpen() {
    this.readyState = 1;
    this.onopen?.();
  }

  simulateMessage(data: string, event = 'message') {
    const msg = new MessageEvent(event, { data });
    if (event === 'message' && this.onmessage) {
      this.onmessage(msg);
    } else {
      const listeners = this.listeners.get(event) || [];
      listeners.forEach((l) => l(msg));
    }
  }

  simulateError() {
    this.onerror?.();
  }
}

describe('useSSE', () => {
  beforeEach(() => {
    MockEventSource.instances = [];
    vi.useFakeTimers();
    (globalThis as any).EventSource = MockEventSource;
  });

  afterEach(() => {
    vi.useRealTimers();
    delete (globalThis as any).EventSource;
  });

  it('crée une connexion EventSource', () => {
    renderHook(() => useSSE('/api/v1/events'));

    expect(MockEventSource.instances).toHaveLength(1);
    expect(MockEventSource.instances[0].url).toBe('/api/v1/events');
  });

  it('ne crée pas de connexion si l\'URL est vide', () => {
    renderHook(() => useSSE(''));

    expect(MockEventSource.instances).toHaveLength(0);
  });

  it('appelle onMessage avec les données JSON parsées', () => {
    const handler = vi.fn();
    renderHook(() => useSSE('/api/v1/events', handler));

    const es = MockEventSource.instances[0];
    es.simulateOpen();
    es.simulateMessage('{"type":"update","id":"j-1"}');

    expect(handler).toHaveBeenCalledWith({ type: 'update', id: 'j-1' });
  });

  it('ignore les messages non-JSON (heartbeats)', () => {
    const handler = vi.fn();
    renderHook(() => useSSE('/api/v1/events', handler));

    const es = MockEventSource.instances[0];
    es.simulateOpen();
    es.simulateMessage(':heartbeat');

    expect(handler).not.toHaveBeenCalled();
  });

  it('écoute un event personnalisé', () => {
    const handler = vi.fn();
    renderHook(() => useSSE('/api/v1/events', handler, 'job-update'));

    const es = MockEventSource.instances[0];
    es.simulateOpen();
    es.simulateMessage('{"progress":0.5}', 'job-update');

    expect(handler).toHaveBeenCalledWith({ progress: 0.5 });
  });

  it('ferme la connexion au démontage', () => {
    const { unmount } = renderHook(() => useSSE('/api/v1/events'));

    const es = MockEventSource.instances[0];
    unmount();

    expect(es.close).toHaveBeenCalled();
  });

  it('reconnecter avec backoff exponentiel après erreur', () => {
    renderHook(() => useSSE('/api/v1/events'));

    const es = MockEventSource.instances[0];
    // Simulate error — should trigger reconnect
    es.simulateError();

    expect(MockEventSource.instances).toHaveLength(1); // Still 1

    // Advance past initial retry delay (1000ms)
    act(() => {
      vi.advanceTimersByTime(1000);
    });

    expect(MockEventSource.instances).toHaveLength(2); // Reconnected

    // Simulate another error — retry should be 2000ms
    MockEventSource.instances[1].simulateError();

    act(() => {
      vi.advanceTimersByTime(1000); // Too early
    });
    expect(MockEventSource.instances).toHaveLength(2);

    act(() => {
      vi.advanceTimersByTime(1000); // Now at 2000ms total
    });
    expect(MockEventSource.instances).toHaveLength(3);
  });

  it('réinitialise le délai après connexion réussie', () => {
    renderHook(() => useSSE('/api/v1/events'));

    const es1 = MockEventSource.instances[0];
    es1.simulateError();

    // First retry at 1000ms
    act(() => { vi.advanceTimersByTime(1000); });
    expect(MockEventSource.instances).toHaveLength(2);

    // Connect success then error again
    MockEventSource.instances[1].simulateOpen(); // Resets delay
    MockEventSource.instances[1].simulateError();

    // Should retry again at 1000ms (reset), not 2000ms
    act(() => { vi.advanceTimersByTime(1000); });
    expect(MockEventSource.instances).toHaveLength(3);
  });

  it('close() ferme manuellement la connexion', () => {
    const { result } = renderHook(() => useSSE('/api/v1/events'));

    const es = MockEventSource.instances[0];
    act(() => {
      result.current.close();
    });

    expect(es.close).toHaveBeenCalled();
  });
});
