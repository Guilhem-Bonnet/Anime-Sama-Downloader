import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { useSearchStore } from '../search.store';

describe('useSearchStore', () => {
  let fetchSpy: ReturnType<typeof vi.fn>;

  beforeEach(() => {
    useSearchStore.setState({
      query: '',
      filters: { genres: [], status: '', yearMin: 0, yearMax: 0 },
      results: [],
      isSearching: false,
      error: undefined,
      lastSearchTime: undefined,
    });
    fetchSpy = vi.fn();
    globalThis.fetch = fetchSpy as typeof fetch;
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('a un état initial correct', () => {
    const state = useSearchStore.getState();
    expect(state.query).toBe('');
    expect(state.results).toEqual([]);
    expect(state.isSearching).toBe(false);
  });

  describe('setters simples', () => {
    it('setQuery met à jour la requête', () => {
      useSearchStore.getState().setQuery('naruto');
      expect(useSearchStore.getState().query).toBe('naruto');
    });

    it('setFilters met à jour les filtres', () => {
      const filters = { genres: ['action'], status: 'airing', yearMin: 2020, yearMax: 2024 };
      useSearchStore.getState().setFilters(filters);
      expect(useSearchStore.getState().filters).toEqual(filters);
    });

    it('clearResults réinitialise tout', () => {
      useSearchStore.setState({ query: 'x', results: [{ id: '1', title: 'x' }], error: 'e' });
      useSearchStore.getState().clearResults();
      const state = useSearchStore.getState();
      expect(state.query).toBe('');
      expect(state.results).toEqual([]);
      expect(state.error).toBeUndefined();
    });
  });

  describe('performSearch', () => {
    it('ne fait rien si la requête est vide', async () => {
      useSearchStore.setState({ query: '' });
      await useSearchStore.getState().performSearch();
      expect(fetchSpy).not.toHaveBeenCalled();
    });

    it('effectue la recherche et stocke les résultats', async () => {
      const mockResults = [{ id: '1', title: 'Naruto' }];
      fetchSpy.mockResolvedValue({
        ok: true,
        json: () => Promise.resolve(mockResults),
      });

      useSearchStore.setState({ query: 'naruto' });
      await useSearchStore.getState().performSearch();

      expect(fetchSpy).toHaveBeenCalledWith(
        expect.stringContaining('/api/v1/search?q=naruto'),
        expect.objectContaining({ signal: expect.any(AbortSignal) })
      );
      expect(useSearchStore.getState().results).toEqual(mockResults);
      expect(useSearchStore.getState().isSearching).toBe(false);
      expect(useSearchStore.getState().lastSearchTime).toBeDefined();
    });

    it('accepte une requête en argument', async () => {
      fetchSpy.mockResolvedValue({ ok: true, json: () => Promise.resolve([]) });

      await useSearchStore.getState().performSearch('bleach');

      expect(fetchSpy).toHaveBeenCalledWith(
        expect.stringContaining('q=bleach'),
        expect.anything()
      );
    });

    it('inclut les filtres dans les paramètres', async () => {
      fetchSpy.mockResolvedValue({ ok: true, json: () => Promise.resolve([]) });

      const filters = { genres: ['action', 'comedy'], status: 'completed', yearMin: 2020, yearMax: 2024 };
      useSearchStore.setState({ query: 'test', filters });
      await useSearchStore.getState().performSearch();

      const url = fetchSpy.mock.calls[0][0] as string;
      expect(url).toContain('genres=action%2Ccomedy');
      expect(url).toContain('status=completed');
      expect(url).toContain('year_min=2020');
      expect(url).toContain('year_max=2024');
    });

    it('gère les erreurs HTTP', async () => {
      fetchSpy.mockResolvedValue({ ok: false, statusText: 'Internal Server Error' });

      useSearchStore.setState({ query: 'test' });
      await useSearchStore.getState().performSearch();

      expect(useSearchStore.getState().error).toContain('Internal Server Error');
      expect(useSearchStore.getState().results).toEqual([]);
    });

    it('ignore les erreurs AbortError', async () => {
      const abortError = new DOMException('aborted', 'AbortError');
      fetchSpy.mockRejectedValue(abortError);

      useSearchStore.setState({ query: 'test' });
      await useSearchStore.getState().performSearch();

      // No error should be set for abort
      expect(useSearchStore.getState().error).toBeUndefined();
    });

    it('gère les réponses non-tableau', async () => {
      fetchSpy.mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({ data: 'not an array' }),
      });

      useSearchStore.setState({ query: 'test' });
      await useSearchStore.getState().performSearch();

      expect(useSearchStore.getState().results).toEqual([]);
    });
  });
});
