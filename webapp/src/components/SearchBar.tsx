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
    <div className="w-full max-w-2xl mx-auto">
      <form onSubmit={handleSubmit} className="mb-4 relative">
        <div className="relative">
          <input
            type="text"
            value={localQuery}
            onChange={handleChange}
            onFocus={() => setSuggestionsOpen(true)}
            placeholder="Search for anime..."
            className="w-full px-4 py-3 rounded-lg bg-gray-200 dark:bg-gray-700 text-gray-900 dark:text-white placeholder-gray-500 dark:placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-cyan-500"
            disabled={isSearching}
            autoComplete="off"
          />
          {isSearching && <div className="absolute right-4 top-3 animate-spin">⌛</div>}

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
