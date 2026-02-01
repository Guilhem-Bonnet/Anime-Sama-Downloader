import React, { useState, useEffect, useRef } from 'react';
import { Search, TrendingUp, History, Tag } from 'lucide-react';

export interface Suggestion {
  query: string;
  category: 'recent' | 'popular' | 'trending' | 'genre';
  score: number;
  metadata?: Record<string, any>;
}

interface SuggestionsDropdownProps {
  query: string;
  onSelectSuggestion: (suggestion: string) => void;
  isOpen: boolean;
  onOpenChange: (open: boolean) => void;
}

export const SuggestionsDropdown: React.FC<SuggestionsDropdownProps> = ({
  query,
  onSelectSuggestion,
  isOpen,
  onOpenChange,
}) => {
  const [suggestions, setSuggestions] = useState<Suggestion[]>([]);
  const [selectedIndex, setSelectedIndex] = useState(-1);
  const [loading, setLoading] = useState(false);
  const dropdownRef = useRef<HTMLDivElement>(null);

  // Fetch suggestions on query change
  useEffect(() => {
    if (!query.trim()) {
      // Fetch trending if empty query
      fetchSuggestions('');
      return;
    }

    // Debounce suggestion fetching
    const timer = setTimeout(() => {
      fetchSuggestions(query);
    }, 300);

    return () => clearTimeout(timer);
  }, [query]);

  // Close dropdown when clicking outside
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target as Node)) {
        onOpenChange(false);
      }
    };

    if (isOpen) {
      document.addEventListener('mousedown', handleClickOutside);
    }

    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
    };
  }, [isOpen, onOpenChange]);

  // Keyboard navigation
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if (!isOpen || suggestions.length === 0) return;

      switch (e.key) {
        case 'ArrowDown':
          e.preventDefault();
          setSelectedIndex((prev) => Math.min(prev + 1, suggestions.length - 1));
          break;
        case 'ArrowUp':
          e.preventDefault();
          setSelectedIndex((prev) => Math.max(prev - 1, -1));
          break;
        case 'Enter':
          e.preventDefault();
          if (selectedIndex >= 0) {
            handleSelect(suggestions[selectedIndex]);
          }
          break;
        case 'Escape':
          e.preventDefault();
          onOpenChange(false);
          break;
      }
    };

    if (isOpen) {
      document.addEventListener('keydown', handleKeyDown);
    }

    return () => {
      document.removeEventListener('keydown', handleKeyDown);
    };
  }, [isOpen, suggestions, selectedIndex, onOpenChange]);

  const fetchSuggestions = async (q: string) => {
    setLoading(true);
    try {
      const params = new URLSearchParams();
      if (q) {
        params.append('q', q);
      }
      params.append('limit', '8');

      // Fallback to search endpoint (works with mock server)
      const response = await fetch(`/api/search?${params.toString()}`);
      if (!response.ok) {
        setSuggestions([]);
        setSelectedIndex(-1);
        return;
      }

      const data = await response.json();
      const results = data.results || [];
      const mapped: Suggestion[] = results.map((item: any) => ({
        query: item.title || item.anime_id || 'Unknown',
        category: q ? 'popular' : 'trending',
        score: 1,
        metadata: {
          episodes: item.episodes,
        },
      }));

      setSuggestions(mapped);
      setSelectedIndex(-1);
    } catch (err) {
      console.error('Failed to fetch suggestions:', err);
      setSuggestions([]);
    } finally {
      setLoading(false);
    }
  };

  const handleSelect = (suggestion: Suggestion) => {
    onSelectSuggestion(suggestion.query);
    onOpenChange(false);
  };

  if (!isOpen || (suggestions.length === 0 && !loading)) {
    return null;
  }

  const getCategoryIcon = (category: Suggestion['category']) => {
    switch (category) {
      case 'recent':
        return <History className="w-4 h-4 text-blue-500" />;
      case 'popular':
        return <Search className="w-4 h-4 text-green-500" />;
      case 'trending':
        return <TrendingUp className="w-4 h-4 text-red-500" />;
      case 'genre':
        return <Tag className="w-4 h-4 text-purple-500" />;
      default:
        return null;
    }
  };

  const getCategoryLabel = (category: Suggestion['category']) => {
    switch (category) {
      case 'recent':
        return 'Recent';
      case 'popular':
        return 'Popular';
      case 'trending':
        return 'Trending';
      case 'genre':
        return 'Genre';
      default:
        return '';
    }
  };

  return (
    <div ref={dropdownRef} className="absolute top-full left-0 right-0 mt-1 bg-white dark:bg-gray-800 rounded-lg shadow-lg border border-gray-200 dark:border-gray-700 z-50 max-h-96 overflow-y-auto">
      {loading && (
        <div className="px-4 py-3 text-center text-sm text-gray-500 dark:text-gray-400">
          Loading suggestions...
        </div>
      )}

      {suggestions.length === 0 && !loading && (
        <div className="px-4 py-3 text-center text-sm text-gray-500 dark:text-gray-400">
          No suggestions found
        </div>
      )}

      {suggestions.map((suggestion, index) => (
        <div
          key={`${suggestion.category}-${suggestion.query}`}
          onClick={() => handleSelect(suggestion)}
          className={`flex items-center gap-3 px-4 py-3 cursor-pointer transition-colors ${
            index === selectedIndex
              ? 'bg-cyan-100 dark:bg-cyan-900'
              : 'hover:bg-gray-50 dark:hover:bg-gray-700'
          }`}
        >
          {getCategoryIcon(suggestion.category)}

          <div className="flex-1 min-w-0">
            <div className="text-sm font-medium text-gray-900 dark:text-white truncate">
              {suggestion.query}
            </div>
            {suggestion.metadata && (
              <div className="text-xs text-gray-500 dark:text-gray-400">
                {suggestion.metadata.anime_count && `${suggestion.metadata.anime_count} anime`}
                {suggestion.metadata.episodes && `${suggestion.metadata.episodes} episodes`}
              </div>
            )}
          </div>

          <div className="flex items-center gap-1">
            <span className="text-xs bg-gray-200 dark:bg-gray-700 text-gray-700 dark:text-gray-300 px-2 py-1 rounded">
              {getCategoryLabel(suggestion.category)}
            </span>
          </div>
        </div>
      ))}
    </div>
  );
};
