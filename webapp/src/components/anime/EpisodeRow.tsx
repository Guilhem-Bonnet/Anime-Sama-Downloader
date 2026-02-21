import React from 'react';

export interface EpisodeRowProps {
  number: number;
  title?: string;
  duration?: string;
  status: 'available' | 'downloading' | 'downloaded';
  selected: boolean;
  disabled?: boolean;
  onChange: (selected: boolean) => void;
}

const EpisodeRow: React.FC<EpisodeRowProps> = ({
  number,
  title,
  duration,
  status,
  selected,
  disabled = false,
  onChange,
}) => {
  const statusColors: Record<string, string> = {
    available: 'var(--night-text-secondary)',
    downloading: 'var(--night-bg-warning)',
    downloaded: 'var(--night-bg-success)',
  };

  const statusLabels: Record<string, string> = {
    available: 'Disponible',
    downloading: 'Téléchargement...',
    downloaded: 'Téléchargé',
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    onChange(e.target.checked);
  };

  return (
    <div
      className={`episode-row ${selected ? 'episode-row--selected' : ''} ${disabled ? 'episode-row--disabled' : ''}`}
      data-testid={`episode-row-${number}`}
    >
      {/* Checkbox */}
      <div className="episode-row__checkbox-wrapper">
        <input
          type="checkbox"
          id={`episode-${number}`}
          className="episode-row__checkbox"
          checked={selected}
          onChange={handleChange}
          disabled={disabled}
          aria-label={`Épisode ${number}`}
        />
      </div>

      {/* Episode Number */}
      <div className="episode-row__number">
        <label htmlFor={`episode-${number}`} className="episode-row__label">
          Ép. {number}
        </label>
      </div>

      {/* Episode Title (optional) */}
      {title && (
        <div className="episode-row__title">
          <span>{title}</span>
        </div>
      )}

      {/* Duration (optional) */}
      {duration && (
        <div className="episode-row__duration">
          <span>{duration}</span>
        </div>
      )}

      {/* Status Badge */}
      <div className="episode-row__status">
        <span
          className="episode-row__badge"
          style={{ backgroundColor: statusColors[status] }}
          aria-describedby={`status-${number}`}
        >
          {statusLabels[status]}
        </span>
        <span id={`status-${number}`} className="sr-only">
          {statusLabels[status]}
        </span>
      </div>

      {/* Loading indicator for downloading state */}
      {status === 'downloading' && (
        <div className="episode-row__loader" aria-hidden="true">
          <div className="episode-row__spinner"></div>
        </div>
      )}
    </div>
  );
};

EpisodeRow.displayName = 'EpisodeRow';

export default EpisodeRow;
