import { useState, useEffect, useCallback } from 'react';

const MAX_RECENT_SEARCHES = 10;
const STORAGE_KEY = 'recent-searches';

export interface RecentSearch {
  query: string;
  timestamp: number;
}

export const useRecentSearches = () => {
  const [recentSearches, setRecentSearches] = useState<RecentSearch[]>([]);

  // Load from localStorage on mount
  useEffect(() => {
    const stored = localStorage.getItem(STORAGE_KEY);
    if (stored) {
      try {
        const parsed = JSON.parse(stored);
        setRecentSearches(parsed);
      } catch (err) {
        console.error('Failed to parse recent searches:', err);
        localStorage.removeItem(STORAGE_KEY);
      }
    }
  }, []);

  // Stable reference — won't cause re-render loops in effect deps
  const addRecentSearch = useCallback((query: string) => {
    if (!query.trim()) return;

    setRecentSearches((prev) => {
      const filtered = prev.filter((s) => s.query !== query);
      const updated = [{ query, timestamp: Date.now() }, ...filtered];
      const limited = updated.slice(0, MAX_RECENT_SEARCHES);
      localStorage.setItem(STORAGE_KEY, JSON.stringify(limited));
      return limited;
    });
  }, []);

  // Clear all recent searches
  const clearRecentSearches = () => {
    setRecentSearches([]);
    localStorage.removeItem(STORAGE_KEY);
  };

  // Remove a single recent search
  const removeRecentSearch = (query: string) => {
    setRecentSearches((prev) => {
      const updated = prev.filter((s) => s.query !== query);
      localStorage.setItem(STORAGE_KEY, JSON.stringify(updated));
      return updated;
    });
  };

  return {
    recentSearches,
    addRecentSearch,
    clearRecentSearches,
    removeRecentSearch,
  };
};
