import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook, act } from '@testing-library/react';
import { useRecentSearches } from '../useRecentSearches';

// The vitest setup substitutes a stub localStorage (returns null for everything).
// We need a real in-memory implementation for these tests.
function createRealLocalStorage(): Storage {
  const store = new Map<string, string>();
  return {
    getItem: (key: string) => store.get(key) ?? null,
    setItem: (key: string, value: string) => { store.set(key, value); },
    removeItem: (key: string) => { store.delete(key); },
    clear: () => store.clear(),
    get length() { return store.size; },
    key: (index: number) => [...store.keys()][index] ?? null,
  };
}

describe('useRecentSearches', () => {
  let realLS: Storage;

  beforeEach(() => {
    realLS = createRealLocalStorage();
    Object.defineProperty(window, 'localStorage', { value: realLS, writable: true });
    vi.clearAllMocks();
  });

  it('retourne un tableau vide par défaut', () => {
    const { result } = renderHook(() => useRecentSearches());
    expect(result.current.recentSearches).toEqual([]);
  });

  it('charge les recherches du localStorage au montage', async () => {
    const stored = [{ query: 'naruto', timestamp: 1000 }];
    realLS.setItem('recent-searches', JSON.stringify(stored));

    const { result } = renderHook(() => useRecentSearches());

    // useEffect is async — wait a tick
    await act(async () => {});

    expect(result.current.recentSearches).toEqual(stored);
  });

  it('gère un localStorage corrompu', async () => {
    realLS.setItem('recent-searches', '{{invalid');

    const { result } = renderHook(() => useRecentSearches());
    await act(async () => {});

    expect(result.current.recentSearches).toEqual([]);
    expect(realLS.getItem('recent-searches')).toBeNull();
  });

  describe('addRecentSearch', () => {
    it('ajoute une recherche en tête de liste', () => {
      const { result } = renderHook(() => useRecentSearches());

      act(() => {
        result.current.addRecentSearch('naruto');
      });

      expect(result.current.recentSearches).toHaveLength(1);
      expect(result.current.recentSearches[0].query).toBe('naruto');
    });

    it('ignore les chaînes vides', () => {
      const { result } = renderHook(() => useRecentSearches());

      act(() => {
        result.current.addRecentSearch('');
      });
      act(() => {
        result.current.addRecentSearch('   ');
      });

      expect(result.current.recentSearches).toHaveLength(0);
    });

    it('déduplique les recherches existantes', () => {
      const { result } = renderHook(() => useRecentSearches());

      act(() => {
        result.current.addRecentSearch('naruto');
      });
      act(() => {
        result.current.addRecentSearch('bleach');
      });
      act(() => {
        result.current.addRecentSearch('naruto');
      });

      expect(result.current.recentSearches).toHaveLength(2);
      expect(result.current.recentSearches[0].query).toBe('naruto');
      expect(result.current.recentSearches[1].query).toBe('bleach');
    });

    it('limite à 10 recherches', () => {
      const { result } = renderHook(() => useRecentSearches());

      act(() => {
        for (let i = 0; i < 12; i++) {
          result.current.addRecentSearch(`search-${i}`);
        }
      });

      expect(result.current.recentSearches).toHaveLength(10);
      // Most recent should be first
      expect(result.current.recentSearches[0].query).toBe('search-11');
    });

    it('persiste dans localStorage', () => {
      const { result } = renderHook(() => useRecentSearches());

      act(() => {
        result.current.addRecentSearch('naruto');
      });

      const stored = JSON.parse(realLS.getItem('recent-searches')!);
      expect(stored).toHaveLength(1);
      expect(stored[0].query).toBe('naruto');
    });
  });

  describe('removeRecentSearch', () => {
    it('supprime une recherche spécifique', () => {
      const { result } = renderHook(() => useRecentSearches());

      act(() => {
        result.current.addRecentSearch('naruto');
        result.current.addRecentSearch('bleach');
      });

      act(() => {
        result.current.removeRecentSearch('naruto');
      });

      expect(result.current.recentSearches).toHaveLength(1);
      expect(result.current.recentSearches[0].query).toBe('bleach');
    });
  });

  describe('clearRecentSearches', () => {
    it('vide toutes les recherches', () => {
      const { result } = renderHook(() => useRecentSearches());

      act(() => {
        result.current.addRecentSearch('naruto');
        result.current.addRecentSearch('bleach');
      });

      act(() => {
        result.current.clearRecentSearches();
      });

      expect(result.current.recentSearches).toEqual([]);
      expect(realLS.getItem('recent-searches')).toBeNull();
    });
  });
});
