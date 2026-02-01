import React, { useState } from 'react';
import { Input } from './components/ui/Input';
import {
  AutocompleteSuggestions,
  type AutocompleteSuggestion,
} from './components/search/AutocompleteSuggestions';

export function SearchDemo() {
  const [query, setQuery] = useState('');
  const [showSuggestions, setShowSuggestions] = useState(false);
  const [selectedAnime, setSelectedAnime] = useState<AutocompleteSuggestion | null>(null);

  const handleSelect = (suggestion: AutocompleteSuggestion) => {
    console.log('Selected anime:', suggestion);
    setSelectedAnime(suggestion);
    setQuery(suggestion.title);
    setShowSuggestions(false);
  };

  const handleClose = () => {
    setShowSuggestions(false);
  };

  return (
    <div style={{ maxWidth: '600px', margin: '40px auto', padding: '24px' }}>
      <h1 style={{ marginBottom: '24px' }}>Search Autocomplete Demo</h1>

      <div style={{ position: 'relative' }}>
        <Input
          value={query}
          onChange={(e) => {
            setQuery(e.target.value);
            setShowSuggestions(true);
          }}
          onFocus={() => setShowSuggestions(true)}
          placeholder="Search for anime (type at least 2 characters)..."
          label="Anime Search"
          hint="Try typing 'naruto', 'one piece', or any anime title"
        />

        {showSuggestions && query.length >= 2 && (
          <AutocompleteSuggestions
            query={query}
            onSelect={handleSelect}
            onClose={handleClose}
          />
        )}
      </div>

      {selectedAnime && (
        <div
          style={{
            marginTop: '32px',
            padding: '16px',
            border: '1px solid var(--border)',
            borderRadius: '12px',
            background: 'rgba(0,0,0,.14)',
          }}
        >
          <h2 style={{ marginBottom: '12px', fontSize: '18px' }}>Selected Anime</h2>
          <div style={{ display: 'flex', gap: '16px', alignItems: 'center' }}>
            <img
              src={selectedAnime.thumbnail_url}
              alt={selectedAnime.title}
              style={{
                width: '80px',
                height: '80px',
                borderRadius: '8px',
                objectFit: 'cover',
              }}
            />
            <div>
              <div style={{ fontWeight: 'bold', fontSize: '16px' }}>
                {selectedAnime.title}
              </div>
              <div style={{ color: 'var(--muted)', fontSize: '14px', marginTop: '4px' }}>
                Year: {selectedAnime.year}
              </div>
              <div style={{ color: 'var(--muted)', fontSize: '14px' }}>
                ID: {selectedAnime.id}
              </div>
            </div>
          </div>
        </div>
      )}

      <div
        style={{
          marginTop: '32px',
          padding: '16px',
          background: 'rgba(124,92,255,.08)',
          border: '1px solid rgba(124,92,255,.3)',
          borderRadius: '12px',
        }}
      >
        <h3 style={{ marginBottom: '8px', fontSize: '14px' }}>Instructions</h3>
        <ul style={{ margin: 0, paddingLeft: '20px', fontSize: '13px', color: 'var(--muted)' }}>
          <li>Type at least 2 characters to see suggestions</li>
          <li>Suggestions appear after 300ms debounce</li>
          <li>Use ↑↓ arrow keys to navigate suggestions</li>
          <li>Press Enter to select highlighted suggestion</li>
          <li>Press Esc to close suggestions</li>
          <li>Click outside to close suggestions</li>
        </ul>
      </div>
    </div>
  );
}
