import React, { useState, useEffect, useRef, KeyboardEvent } from 'react';
import { useNavigate } from 'react-router-dom';
import { useDebounce } from '../../hooks/useDebounce';

export interface AutocompleteSuggestion {
  id: string;
  title: string;
  thumbnail_url: string;
  year: number;
}

interface AutocompleteSuggestionsProps {
  query: string;
  onSelect: (suggestion: AutocompleteSuggestion) => void;
  onClose: () => void;
  apiEndpoint?: string;
}

export function AutocompleteSuggestions({
  query,
  onSelect,
  onClose,
  apiEndpoint = '/api/v1/search/autocomplete',
}: AutocompleteSuggestionsProps) {
  const navigate = useNavigate();
  const [suggestions, setSuggestions] = useState<AutocompleteSuggestion[]>([]);
  const [selectedIndex, setSelectedIndex] = useState(-1);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const containerRef = useRef<HTMLDivElement>(null);
  const abortControllerRef = useRef<AbortController | null>(null);

  const debouncedQuery = useDebounce(query, 300);

  // Handle selection: navigate to anime detail page
  const handleSelect = (suggestion: AutocompleteSuggestion) => {
    onSelect(suggestion);
    navigate(`/anime/${suggestion.id}`);
  };

  // Fetch autocomplete suggestions
  useEffect(() => {
    // Reset state when query is too short
    if (debouncedQuery.length < 2) {
      setSuggestions([]);
      setSelectedIndex(-1);
      setError(null);
      return;
    }

    // Cancel previous request
    if (abortControllerRef.current) {
      abortControllerRef.current.abort();
    }

    // Create new abort controller
    const controller = new AbortController();
    abortControllerRef.current = controller;

    setLoading(true);
    setError(null);

    fetch(`${apiEndpoint}?q=${encodeURIComponent(debouncedQuery)}`, {
      signal: controller.signal,
    })
      .then((res) => {
        if (!res.ok) {
          throw new Error(`HTTP ${res.status}`);
        }
        return res.json();
      })
      .then((data: AutocompleteSuggestion[]) => {
        setSuggestions(data);
        setSelectedIndex(-1);
        setLoading(false);
      })
      .catch((err) => {
        if (err.name === 'AbortError') {
          // Request was cancelled, ignore
          return;
        }
        setError('Failed to load suggestions');
        setSuggestions([]);
        setLoading(false);
      });

    // Cleanup on unmount
    return () => {
      controller.abort();
    };
  }, [debouncedQuery, apiEndpoint]);

  // Handle click outside to close
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (
        containerRef.current &&
        !containerRef.current.contains(event.target as Node)
      ) {
        onClose();
      }
    };

    document.addEventListener('mousedown', handleClickOutside);
    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
    };
  }, [onClose]);

  // Handle keyboard navigation
  const handleKeyDown = (e: KeyboardEvent) => {
    if (suggestions.length === 0) return;

    switch (e.key) {
      case 'ArrowDown':
        e.preventDefault();
        setSelectedIndex((prev) =>
          prev < suggestions.length - 1 ? prev + 1 : 0
        );
        break;
      case 'ArrowUp':
        e.preventDefault();
        setSelectedIndex((prev) =>
          prev > 0 ? prev - 1 : suggestions.length - 1
        );
        break;
      case 'Enter':
        e.preventDefault();
        if (selectedIndex >= 0 && selectedIndex < suggestions.length) {
          onSelect(suggestions[selectedIndex]);
        }
        break;
      case 'Escape':
        e.preventDefault();
        onClose();
        break;
    }
  };

  // Scroll selected item into view
  useEffect(() => {
    if (selectedIndex >= 0 && containerRef.current) {
      const selectedElement = containerRef.current.querySelector(
        `[data-index="${selectedIndex}"]`
      );
      if (selectedElement) {
        selectedElement.scrollIntoView({
          block: 'nearest',
          behavior: 'smooth',
        });
      }
    }
  }, [selectedIndex]);

  // Don't render if query is too short
  if (query.length < 2) {
    return null;
  }

  return (
    <div
      ref={containerRef}
      className="autocomplete-container"
      onKeyDown={handleKeyDown}
      tabIndex={-1}
    >
      {loading && (
        <div className="autocomplete-loading">
          <span>Loading suggestions...</span>
        </div>
      )}

      {error && (
        <div className="autocomplete-error">
          <span>{error}</span>
        </div>
      )}

      {!loading && !error && suggestions.length === 0 && (
        <div className="autocomplete-no-results">
          <span>No results found</span>
        </div>
      )}

      {!loading && !error && suggestions.length > 0 && (
        <ul className="autocomplete-list">
          {suggestions.map((suggestion, index) => (
            <li
              key={suggestion.id}
              data-index={index}
              className={`autocomplete-item ${
                index === selectedIndex ? 'autocomplete-item-selected' : ''
              }`}
              onClick={() => handleSelect(suggestion)}
              onMouseEnter={() => setSelectedIndex(index)}
            >
              <img
                src={suggestion.thumbnail_url}
                alt={suggestion.title}
                className="autocomplete-thumbnail"
              />
              <div className="autocomplete-text">
                <span className="autocomplete-title">{suggestion.title}</span>
                <span className="autocomplete-year">{suggestion.year}</span>
              </div>
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}
