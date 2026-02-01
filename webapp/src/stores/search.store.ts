import { create } from 'zustand';
import { persist } from 'zustand/middleware';

export interface SearchResult {
  animeId: string;
  title: string;
  episodes: number;
  source: string;
  imageUrl?: string;
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
          // Build query params
          const params = new URLSearchParams();
          params.append('q', searchQuery);

          if (searchFilters.genres.length > 0) {
            params.append('genres', searchFilters.genres.join(','));
          }
          if (searchFilters.status) {
            params.append('status', searchFilters.status);
          }
          if (searchFilters.yearMin > 0) {
            params.append('year_min', searchFilters.yearMin.toString());
          }
          if (searchFilters.yearMax > 0) {
            params.append('year_max', searchFilters.yearMax.toString());
          }

          const response = await fetch(`/api/v1/search?${params.toString()}`);
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
