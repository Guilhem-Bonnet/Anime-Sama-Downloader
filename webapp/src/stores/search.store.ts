import { create } from 'zustand';
import { persist } from 'zustand/middleware';

const API = '/api/v1';

// Shared AbortController — cancels previous in-flight search
let searchAbortController: AbortController | null = null;

export interface SearchResult {
  id: string;
  title: string;
  thumbnail_url?: string;
  year?: number;
  status?: string;
  episode_count?: number;
  genres?: string[];
}

export interface SearchFilters {
  genres: string[];
  status: string;
  yearMin: number;
  yearMax: number;
}

export interface SearchState {
  query: string;
  filters: SearchFilters;
  results: SearchResult[];
  isSearching: boolean;
  error?: string;
  lastSearchTime?: number;
}

export interface SearchActions {
  setQuery: (query: string) => void;
  setFilters: (filters: SearchFilters) => void;
  setResults: (results: SearchResult[]) => void;
  setIsSearching: (searching: boolean) => void;
  setError: (error?: string) => void;
  clearResults: () => void;
  performSearch: (query?: string, filters?: SearchFilters) => Promise<void>;
}

export const useSearchStore = create<SearchState & SearchActions>()(
  persist(
    (set, get) => ({
      // State
      query: '',
      filters: {
        genres: [],
        status: '',
        yearMin: 0,
        yearMax: 0,
      },
      results: [],
      isSearching: false,

      // Actions
      setQuery: (query) => set({ query }),
      setFilters: (filters) => set({ filters }),
      setResults: (results) => set({ results }),
      setIsSearching: (isSearching) => set({ isSearching }),
      setError: (error) => set({ error }),
      clearResults: () =>
        set({
          results: [],
          query: '',
          filters: { genres: [], status: '', yearMin: 0, yearMax: 0 },
          error: undefined,
        }),

      performSearch: async (queryArg?: string, filtersArg?: SearchFilters) => {
        const { query, filters } = get();

        const searchQuery = queryArg !== undefined ? queryArg : query;
        const searchFilters = filtersArg !== undefined ? filtersArg : filters;

        if (!searchQuery.trim()) {
          set({ results: [] });
          return;
        }

        // Abort previous in-flight request
        if (searchAbortController) {
          searchAbortController.abort();
        }
        const controller = new AbortController();
        searchAbortController = controller;

        set({ isSearching: true, error: undefined });

        try {
          const params = new URLSearchParams({ q: searchQuery });
          if (searchFilters.genres.length > 0) params.set('genres', searchFilters.genres.join(','));
          if (searchFilters.status) params.set('status', searchFilters.status);
          if (searchFilters.yearMin > 0) params.set('year_min', String(searchFilters.yearMin));
          if (searchFilters.yearMax > 0) params.set('year_max', String(searchFilters.yearMax));

          const response = await fetch(`${API}/search?${params.toString()}`, {
            signal: controller.signal,
          });
          if (!response.ok) throw new Error(`Recherche échouée: ${response.statusText}`);
          const data = await response.json();
          const results: SearchResult[] = Array.isArray(data) ? data : [];
          set({ results, lastSearchTime: Date.now() });
        } catch (err) {
          // Don't treat abort as an error
          if (err instanceof DOMException && err.name === 'AbortError') return;
          set({ error: err instanceof Error ? err.message : 'Erreur de recherche', results: [] });
        } finally {
          // Only clear isSearching if this is still the current request
          if (searchAbortController === controller) {
            set({ isSearching: false });
          }
        }
      },
    }),
    {
      name: 'search-filters-storage',
      partialize: (state) => ({ filters: state.filters }),
    }
  )
);
