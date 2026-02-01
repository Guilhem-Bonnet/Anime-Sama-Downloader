import React, { useState, useEffect } from 'react';
import { useSearchStore } from '../stores/search.store';
import { useDebounce } from '../hooks/useDebounce';
import { useRecentSearches } from '../hooks/useRecentSearches';
import { FilterPanel } from './search/FilterPanel';
import { RecentSearchesDropdown } from './search/RecentSearchesDropdown';
import { SuggestionsDropdown } from './search/SuggestionsDropdown';

export const SearchBar: React.FC = () => {
  const { query, filters, setQuery, setFilters, performSearch, isSearching, error } =
    useSearchStore();
  const [localQuery, setLocalQuery] = useState(query);
  const [suggestionsOpen, setSuggestionsOpen] = useState(false);
  const debouncedQuery = useDebounce(localQuery, 500);
  const { addRecentSearch } = useRecentSearches();

  useEffect(() => {
    if (debouncedQuery) {
      performSearch(debouncedQuery, filters);
      addRecentSearch(debouncedQuery); // Save to recent searches
    }
  }, [debouncedQuery, performSearch, filters, addRecentSearch]);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    setLocalQuery(value);
    setQuery(value);
    setSuggestionsOpen(true); // Open suggestions on input change
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (localQuery.trim()) {
      performSearch(localQuery, filters);
      addRecentSearch(localQuery); // Save to recent searches
      setSuggestionsOpen(false);
    }
  };

  const handleFiltersChange = (newFilters: typeof filters) => {
    setFilters(newFilters);
    // Re-trigger search with new filters if we have a query
    if (localQuery.trim()) {
      performSearch(localQuery, newFilters);
    }
  };

  const handleSelectRecentSearch = (recentQuery: string) => {
    setLocalQuery(recentQuery);
    setQuery(recentQuery);
    performSearch(recentQuery, filters);
    addRecentSearch(recentQuery); // Update recent searches timestamp
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
        <div className="relative group">
          {/* Search Icon */}
          <div className="absolute left-5 top-1/2 -translate-y-1/2 text-gray-400 group-focus-within:text-cyan-500 transition-colors">
            🔍
          </div>
          
          <input
            type="text"
            value={localQuery}
            onChange={handleChange}
            onFocus={() => setSuggestionsOpen(true)}
            placeholder="Rechercher un anime... (ex: Attack on Titan)"
            className="w-full pl-14 pr-14 py-4 rounded-2xl bg-white dark:bg-gray-800 text-gray-900 dark:text-white placeholder-gray-400 dark:placeholder-gray-500 border-2 border-gray-200 dark:border-gray-700 focus:outline-none focus:border-transparent focus:ring-4 focus:ring-cyan-500/30 shadow-lg hover:shadow-xl transition-all duration-300 text-lg"
            disabled={isSearching}
            autoComplete="off"
          />
          
          {/* Loading spinner */}
          {isSearching && (
            <div className="absolute right-5 top-1/2 -translate-y-1/2 animate-spin text-2xl">
              ⌛
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
              className="absolute right-5 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600 dark:hover:text-gray-300 transition-colors text-xl"
            >
              ✕
            </button>
          )}

          <SuggestionsDropdown
            query={localQuery}
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
