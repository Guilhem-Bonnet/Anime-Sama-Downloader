import React, { useRef, useEffect, useState } from 'react';
import { Clock, X } from 'lucide-react';
import { useRecentSearches, RecentSearch } from '../../hooks/useRecentSearches';

interface RecentSearchesDropdownProps {
  onSelectSearch: (query: string) => void;
  currentQuery: string;
}

export const RecentSearchesDropdown: React.FC<RecentSearchesDropdownProps> = ({
  onSelectSearch,
  currentQuery,
}) => {
  const { recentSearches, clearRecentSearches, removeRecentSearch } = useRecentSearches();
  const [isOpen, setIsOpen] = useState(false);
  const dropdownRef = useRef<HTMLDivElement>(null);

  // Close dropdown when clicking outside
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target as Node)) {
        setIsOpen(false);
      }
    };

    if (isOpen) {
      document.addEventListener('mousedown', handleClickOutside);
    }

    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
    };
  }, [isOpen]);

  const handleSelectSearch = (search: RecentSearch) => {
    onSelectSearch(search.query);
    setIsOpen(false);
  };

  const handleRemove = (e: React.MouseEvent, query: string) => {
    e.stopPropagation();
    removeRecentSearch(query);
  };

  const handleClearAll = (e: React.MouseEvent) => {
    e.stopPropagation();
    clearRecentSearches();
  };

  // Filter out current query from recent searches
  const filteredSearches = recentSearches.filter((s) => s.query !== currentQuery);

  if (filteredSearches.length === 0) {
    return null;
  }

  return (
    <div ref={dropdownRef} className="relative">
      <button
        onClick={() => setIsOpen(!isOpen)}
        className="flex items-center gap-2 px-3 py-2 text-sm text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-gray-200 transition-colors"
        title="Recent searches"
      >
        <Clock className="w-4 h-4" />
        <span>Recent ({filteredSearches.length})</span>
      </button>

      {isOpen && (
        <div className="absolute top-full left-0 mt-2 w-80 bg-white dark:bg-gray-800 rounded-lg shadow-lg border border-gray-200 dark:border-gray-700 z-50">
          <div className="flex items-center justify-between px-4 py-3 border-b border-gray-200 dark:border-gray-700">
            <h3 className="text-sm font-medium text-gray-700 dark:text-gray-300">
              Recent Searches
            </h3>
            <button
              onClick={handleClearAll}
              className="text-xs text-red-500 hover:text-red-600 dark:text-red-400 dark:hover:text-red-300"
            >
              Clear All
            </button>
          </div>

          <div className="max-h-96 overflow-y-auto">
            {filteredSearches.map((search) => (
              <button
                key={search.query}
                onClick={() => handleSelectSearch(search)}
                className="flex items-center justify-between w-full px-4 py-3 hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors group"
              >
                <div className="flex items-center gap-3 flex-1 min-w-0">
                  <Clock className="w-4 h-4 text-gray-400 flex-shrink-0" />
                  <span className="text-sm text-gray-700 dark:text-gray-300 truncate">
                    {search.query}
                  </span>
                </div>

                <button
                  onClick={(e) => handleRemove(e, search.query)}
                  className="ml-2 p-1 opacity-0 group-hover:opacity-100 hover:bg-gray-200 dark:hover:bg-gray-600 rounded transition-all"
                  title="Remove from recent searches"
                >
                  <X className="w-3 h-3 text-gray-500 dark:text-gray-400" />
                </button>
              </button>
            ))}
          </div>
        </div>
      )}
    </div>
  );
};
