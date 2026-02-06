import React, { useState, useEffect } from 'react';
import { Button } from '../ui/Button';
import { Card, CardBody } from '../ui/Card';
import { Badge } from '../ui/Badge';

export interface EpisodeSelectorProps {
  maxEpisode: number;
  selectedEpisodes: number[];
  onChange: (episodes: number[]) => void;
  previewUrl?: string;
}

export function EpisodeSelector({
  maxEpisode,
  selectedEpisodes,
  onChange,
  previewUrl,
}: EpisodeSelectorProps) {
  const [tempSelection, setTempSelection] = useState<number[]>(selectedEpisodes);
  const [rangeStart, setRangeStart] = useState<string>('');
  const [rangeEnd, setRangeEnd] = useState<string>('');
  const [rangeError, setRangeError] = useState<string>('');

  useEffect(() => {
    setTempSelection(selectedEpisodes);
  }, [selectedEpisodes]);

  const handleSelectAll = () => {
    const all = Array.from({ length: maxEpisode }, (_, i) => i + 1);
    setTempSelection(all);
    onChange(all);
  };

  const handleSelectLast = (count: number) => {
    const last = Array.from(
      { length: Math.min(count, maxEpisode) },
      (_, i) => maxEpisode - i
    ).reverse();
    setTempSelection(last);
    onChange(last);
  };

  const handleCustomRange = () => {
    const start = parseInt(rangeStart);
    const end = parseInt(rangeEnd);

    if (isNaN(start) || isNaN(end)) {
      setRangeError('Veuillez entrer des nombres valides');
      return;
    }

    if (start < 1 || end > maxEpisode) {
      setRangeError(`Les épisodes doivent être entre 1 et ${maxEpisode}`);
      return;
    }

    if (start > end) {
      setRangeError('Le début doit être inférieur ou égal à la fin');
      return;
    }

    setRangeError('');
    const range = Array.from({ length: end - start + 1 }, (_, i) => start + i);
    setTempSelection(range);
    onChange(range);
  };

  const handleToggleEpisode = (episode: number) => {
    const newSelection = tempSelection.includes(episode)
      ? tempSelection.filter((e) => e !== episode)
      : [...tempSelection, episode].sort((a, b) => a - b);
    setTempSelection(newSelection);
    onChange(newSelection);
  };

  const handleClear = () => {
    setTempSelection([]);
    onChange([]);
  };

  return (
    <Card style={{ marginTop: 'var(--space-4)' }}>
      <CardBody>
        <div style={{ display: 'flex', flexDirection: 'column', gap: 'var(--space-4)' }}>
          {/* Shortcuts */}
          <div>
            <label
              style={{
                display: 'block',
                marginBottom: 'var(--space-2)',
                fontSize: 'var(--text-label)',
                fontWeight: 600,
                color: 'var(--sakura-text-secondary)',
              }}
            >
              Raccourcis de sélection
            </label>
            <div style={{ display: 'flex', gap: 'var(--space-2)', flexWrap: 'wrap' }}>
              <Button variant="secondary" onClick={handleSelectAll}>
                Tous les épisodes
              </Button>
              <Button variant="secondary" onClick={() => handleSelectLast(5)}>
                Derniers 5
              </Button>
              <Button variant="secondary" onClick={() => handleSelectLast(10)}>
                Derniers 10
              </Button>
              <Button variant="ghost" onClick={handleClear}>
                Effacer tout
              </Button>
            </div>
          </div>

          {/* Custom Range */}
          <div>
            <label
              style={{
                display: 'block',
                marginBottom: 'var(--space-2)',
                fontSize: 'var(--text-label)',
                fontWeight: 600,
                color: 'var(--sakura-text-secondary)',
              }}
            >
              Plage personnalisée
            </label>
            <div style={{ display: 'flex', gap: 'var(--space-2)', alignItems: 'flex-start' }}>
              <input
                type="number"
                placeholder="Début"
                min="1"
                max={maxEpisode}
                value={rangeStart}
                onChange={(e) => setRangeStart(e.target.value)}
                onKeyDown={(e) => {
                  if (e.key === 'Enter') handleCustomRange();
                }}
                style={{
                  padding: '8px 12px',
                  borderRadius: 'var(--radius-md)',
                  border: `1px solid ${rangeError ? 'var(--sakura-error-border)' : 'var(--sakura-border-default)'}`,
                  background: 'var(--sakura-bg-surface)',
                  color: 'var(--sakura-text-primary)',
                  fontSize: 'var(--text-body)',
                  width: '100px',
                }}
              />
              <span style={{ color: 'var(--sakura-text-secondary)', lineHeight: '36px' }}>à</span>
              <input
                type="number"
                placeholder="Fin"
                min="1"
                max={maxEpisode}
                value={rangeEnd}
                onChange={(e) => setRangeEnd(e.target.value)}
                onKeyDown={(e) => {
                  if (e.key === 'Enter') handleCustomRange();
                }}
                style={{
                  padding: '8px 12px',
                  borderRadius: 'var(--radius-md)',
                  border: `1px solid ${rangeError ? 'var(--sakura-error-border)' : 'var(--sakura-border-default)'}`,
                  background: 'var(--sakura-bg-surface)',
                  color: 'var(--sakura-text-primary)',
                  fontSize: 'var(--text-body)',
                  width: '100px',
                }}
              />
              <Button variant="primary" onClick={handleCustomRange}>
                Appliquer
              </Button>
            </div>
            {rangeError && (
              <p
                style={{
                  marginTop: 'var(--space-2)',
                  fontSize: 'var(--text-body-sm)',
                  color: 'var(--sakura-error-text)',
                }}
              >
                {rangeError}
              </p>
            )}
          </div>

          {/* Episode Checkboxes */}
          <div>
            <div
              style={{
                display: 'flex',
                justifyContent: 'space-between',
                alignItems: 'center',
                marginBottom: 'var(--space-2)',
              }}
            >
              <label
                style={{
                  fontSize: 'var(--text-label)',
                  fontWeight: 600,
                  color: 'var(--sakura-text-secondary)',
                }}
              >
                Sélection manuelle ({tempSelection.length} sélectionné{tempSelection.length > 1 ? 's' : ''})
              </label>
              {tempSelection.length > 0 && (
                <Badge variant="info">{tempSelection.length} épisodes</Badge>
              )}
            </div>
            <div
              style={{
                maxHeight: '320px',
                overflowY: 'auto',
                border: '1px solid var(--sakura-border-default)',
                borderRadius: 'var(--radius-md)',
                padding: 'var(--space-3)',
                background: 'var(--sakura-bg-elevated)',
              }}
            >
              <div
                style={{
                  display: 'grid',
                  gridTemplateColumns: 'repeat(auto-fill, minmax(120px, 1fr))',
                  gap: 'var(--space-2)',
                }}
              >
                {Array.from({ length: maxEpisode }, (_, i) => i + 1).map((episode) => {
                  const isSelected = tempSelection.includes(episode);
                  return (
                    <label
                      key={episode}
                      style={{
                        display: 'flex',
                        alignItems: 'center',
                        gap: 'var(--space-2)',
                        padding: 'var(--space-2)',
                        borderRadius: 'var(--radius-sm)',
                        cursor: 'pointer',
                        background: isSelected
                          ? 'var(--sakura-accent-magenta-900)'
                          : 'transparent',
                        border: `1px solid ${
                          isSelected
                            ? 'var(--sakura-accent-magenta-500)'
                            : 'var(--sakura-border-subtle)'
                        }`,
                        transition: 'all var(--transition-fast)',
                      }}
                      onMouseEnter={(e) => {
                        if (!isSelected) {
                          e.currentTarget.style.background = 'var(--sakura-bg-surface)';
                        }
                      }}
                      onMouseLeave={(e) => {
                        if (!isSelected) {
                          e.currentTarget.style.background = 'transparent';
                        }
                      }}
                    >
                      <input
                        type="checkbox"
                        checked={isSelected}
                        onChange={() => handleToggleEpisode(episode)}
                        style={{
                          width: '16px',
                          height: '16px',
                          cursor: 'pointer',
                          accentColor: 'var(--sakura-accent-magenta-500)',
                        }}
                      />
                      <span
                        style={{
                          fontSize: 'var(--text-body-sm)',
                          color: isSelected
                            ? 'var(--sakura-text-primary)'
                            : 'var(--sakura-text-secondary)',
                          fontWeight: isSelected ? 600 : 400,
                        }}
                      >
                        Ép. {episode}
                      </span>
                    </label>
                  );
                })}
              </div>
            </div>
          </div>

          {/* Preview (if URL provided) */}
          {previewUrl && (
            <div>
              <label
                style={{
                  display: 'block',
                  marginBottom: 'var(--space-2)',
                  fontSize: 'var(--text-label)',
                  fontWeight: 600,
                  color: 'var(--sakura-text-secondary)',
                }}
              >
                Aperçu
              </label>
              <div
                style={{
                  aspectRatio: '16/9',
                  background: 'var(--sakura-bg-base)',
                  borderRadius: 'var(--radius-md)',
                  overflow: 'hidden',
                  border: '1px solid var(--sakura-border-default)',
                }}
              >
                <img
                  src={previewUrl}
                  alt="Preview"
                  style={{ width: '100%', height: '100%', objectFit: 'cover' }}
                />
              </div>
            </div>
          )}
        </div>
      </CardBody>
    </Card>
  );
}
