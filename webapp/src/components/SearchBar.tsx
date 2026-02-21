import React, { useState, useEffect, useRef, useCallback } from 'react';
import { useSearchStore } from '../stores/search.store';
import { useDebounce } from '../hooks/useDebounce';
import { useRecentSearches } from '../hooks/useRecentSearches';
import { FilterPanel } from './search/FilterPanel';
import { RecentSearchesDropdown } from './search/RecentSearchesDropdown';
import { SuggestionsDropdown } from './search/SuggestionsDropdown';
import { Search, X, Loader2 } from 'lucide-react';

export const SearchBar: React.FC = () => {
  const { query, filters, setQuery, setFilters, performSearch, isSearching, error } =
    useSearchStore();
  const [localQuery, setLocalQuery] = useState(query);
  const [suggestionsOpen, setSuggestionsOpen] = useState(false);
  const [isFocused, setIsFocused] = useState(false);
  const debouncedQuery = useDebounce(localQuery, 500);
  const { addRecentSearch } = useRecentSearches();
  const filtersRef = useRef(filters);
  filtersRef.current = filters;

  // Fire search only when debouncedQuery changes — no unstable deps
  useEffect(() => {
    if (debouncedQuery.trim()) {
      performSearch(debouncedQuery, filtersRef.current);
      addRecentSearch(debouncedQuery);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [debouncedQuery]);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    setLocalQuery(value);
    setSuggestionsOpen(true);
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (localQuery.trim()) {
      setQuery(localQuery);
      performSearch(localQuery, filters);
      addRecentSearch(localQuery);
      setSuggestionsOpen(false);
    }
  };

  const handleFiltersChange = (newFilters: typeof filters) => {
    setFilters(newFilters);
    // Re-trigger search with new filters (no double search — effect only watches debouncedQuery)
    if (localQuery.trim()) {
      performSearch(localQuery, newFilters);
    }
  };

  const handleSelectRecentSearch = (recentQuery: string) => {
    setLocalQuery(recentQuery);
    setQuery(recentQuery);
    performSearch(recentQuery, filters);
    addRecentSearch(recentQuery);
    setSuggestionsOpen(false);
  };

  const handleSelectSuggestion = (suggestion: string) => {
    setLocalQuery(suggestion);
    setQuery(suggestion);
    performSearch(suggestion, filters);
    addRecentSearch(suggestion);
    setSuggestionsOpen(false);
  };

  return (
    <div className="w-full max-w-3xl mx-auto">
      <form onSubmit={handleSubmit} className="mb-6 relative">
        <div className={`relative group search-focus-ink ${isFocused ? 'is-focused' : ''}`}>
          {/* Search Icon */}
          <div className="absolute left-5 top-1/2 -translate-y-1/2 text-gray-400 group-focus-within:text-gray-200 transition-colors">
            <Search className="w-5 h-5" />
          </div>
          
          <input
            type="text"
            value={localQuery}
            onChange={handleChange}
            onFocus={() => {
              setIsFocused(true);
              setSuggestionsOpen(true);
            }}
            onBlur={() => setIsFocused(false)}
            placeholder="Rechercher un anime... (ex: Attack on Titan)"
            className="input w-full pl-14 pr-14 py-4 rounded-2xl text-lg shadow-lg hover:shadow-xl transition-all duration-300"
            autoComplete="off"
          />
          
          {/* Loading spinner */}
          {isSearching && (
            <div className="absolute right-5 top-1/2 -translate-y-1/2">
              <Loader2 className="w-5 h-5 animate-spin text-gray-400" />
            </div>
          )}
          
          {/* Clear button */}
          {!isSearching && localQuery && (
            <button
              type="button"
              onClick={() => {
                setLocalQuery('');
                setQuery('');
              }}
              className="absolute right-5 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600 dark:hover:text-gray-300 transition-colors"
            >
              <X className="w-5 h-5" />
            </button>
          )}

          <SuggestionsDropdown
            query={debouncedQuery}
            onSelectSuggestion={handleSelectSuggestion}
            isOpen={suggestionsOpen}
            onOpenChange={setSuggestionsOpen}
          />
        </div>
        {error && <p className="text-red-500 text-sm mt-2">{error}</p>}
      </form>

      <div className="flex items-center gap-4 mb-4">
        <RecentSearchesDropdown
          onSelectSearch={handleSelectRecentSearch}
          currentQuery={localQuery}
        />
      </div>

      <FilterPanel filters={filters} onFiltersChange={handleFiltersChange} />
    </div>
  );
};
