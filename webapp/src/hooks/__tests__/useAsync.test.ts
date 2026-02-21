import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook, waitFor, act } from '@testing-library/react';
import { useAsync } from '../useAsync';

describe('useAsync', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('exécute immédiatement par défaut', async () => {
    const asyncFn = vi.fn().mockResolvedValue('data');

    const { result } = renderHook(() => useAsync(asyncFn));

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(asyncFn).toHaveBeenCalledOnce();
    expect(result.current.data).toBe('data');
    expect(result.current.error).toBeNull();
  });

  it('n\'exécute pas immédiatement si immediate=false', async () => {
    const asyncFn = vi.fn().mockResolvedValue('data');

    const { result } = renderHook(() => useAsync(asyncFn, false));

    // Give it a tick to ensure it wasn't called
    await new Promise((resolve) => setTimeout(resolve, 50));

    expect(asyncFn).not.toHaveBeenCalled();
    expect(result.current.data).toBeNull();
    expect(result.current.loading).toBe(false);
  });

  it('gère les erreurs', async () => {
    const asyncFn = vi.fn().mockRejectedValue(new Error('échec'));

    const { result } = renderHook(() => useAsync(asyncFn));

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.data).toBeNull();
    expect(result.current.error).toBeInstanceOf(Error);
    expect(result.current.error?.message).toBe('échec');
  });

  it('convertit les erreurs non-Error en Error', async () => {
    const asyncFn = vi.fn().mockRejectedValue('string error');

    const { result } = renderHook(() => useAsync(asyncFn));

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.error?.message).toBe('string error');
  });

  it('execute() peut être appelé manuellement', async () => {
    const asyncFn = vi.fn().mockResolvedValue('manual');

    const { result } = renderHook(() => useAsync(asyncFn, false));

    expect(result.current.data).toBeNull();

    await act(async () => {
      await result.current.execute();
    });

    expect(result.current.data).toBe('manual');
    expect(asyncFn).toHaveBeenCalledOnce();
  });

  it('est en loading pendant l\'exécution', async () => {
    let resolve!: (value: string) => void;
    const asyncFn = vi.fn().mockImplementation(
      () => new Promise<string>((r) => { resolve = r; })
    );

    const { result } = renderHook(() => useAsync(asyncFn));

    // Wait for loading to be true
    await waitFor(() => {
      expect(result.current.loading).toBe(true);
    });

    await act(async () => {
      resolve('done');
    });

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
      expect(result.current.data).toBe('done');
    });
  });
});
