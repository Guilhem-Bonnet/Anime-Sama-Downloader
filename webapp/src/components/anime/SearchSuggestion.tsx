import React from 'react';

export interface SearchSuggestionProps {
  title: string;
  season?: string;
  language: 'VOSTFR' | 'VF';
  query: string;
  onSelect: () => void;
}

const SearchSuggestion: React.FC<SearchSuggestionProps> = ({
  title,
  season,
  language,
  query,
  onSelect,
}) => {
  // Highlight query in title
  const renderHighlightedTitle = (text: string, highlight: string) => {
    if (!highlight || highlight.length === 0) {
      return text;
    }

    const parts = text.split(new RegExp(`(${highlight})`, 'gi'));

    return (
      <>
        {parts.map((part, i) => 
          part.toLowerCase() === highlight.toLowerCase() ? (
            <mark key={i} className="search-suggestion__highlight">
              {part}
            </mark>
          ) : (
            <span key={i}>{part}</span>
          )
        )}
      </>
    );
  };

  return (
    <div
      className="search-suggestion"
      role="option"
      aria-selected="false"
      tabIndex={0}
      onClick={onSelect}
      onKeyDown={(e) => {
        if (e.key === 'Enter' || e.key === ' ') {
          e.preventDefault();
          onSelect();
        }
      }}
      data-testid={`search-suggestion-${title}`}
    >
      {/* Search Icon */}
      <span className="search-suggestion__icon" aria-hidden="true">
        🔍
      </span>

      {/* Title with Highlight */}
      <div className="search-suggestion__content">
        <div className="search-suggestion__title">
          {renderHighlightedTitle(title, query)}
        </div>

        {/* Metadata: Season + Language */}
        <div className="search-suggestion__metadata">
          {season && (
            <span className="search-suggestion__season">{season}</span>
          )}
          <span className="search-suggestion__language">{language}</span>
        </div>
      </div>
    </div>
  );
};

SearchSuggestion.displayName = 'SearchSuggestion';

export default SearchSuggestion;
