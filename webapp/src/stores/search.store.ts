import { create } from 'zustand';
import { persist } from 'zustand/middleware';

export interface SearchResult {
  anime_id: string;   // from API
  title: string;
  episodes: number;
  source: string;
  image_url?: string; // from API
  description?: string;
  genres?: string[];
  year?: number;
  status?: string;
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
        const { setIsSearching, setResults, setError, query, filters } = get();

        const searchQuery = queryArg !== undefined ? queryArg : query;
        const searchFilters = filtersArg !== undefined ? filtersArg : filters;

        if (!searchQuery.trim()) {
          setResults([]);
          return;
        }

        setIsSearching(true);
        setError(undefined);

        try {
          const response = await fetch(`http://localhost:8000/api/search?q=${encodeURIComponent(searchQuery)}`);
          if (!response.ok) throw new Error('Search failed');

          const data = await response.json();
          setResults(data.results || []);
          set({ lastSearchTime: Date.now() });
        } catch (err) {
          setError(err instanceof Error ? err.message : 'Search error');
          setResults([]);
        } finally {
          setIsSearching(false);
        }
      },
    }),
    {
      name: 'search-filters-storage',
      partialize: (state) => ({ filters: state.filters }), // Only persist filters
    }
  )
);
