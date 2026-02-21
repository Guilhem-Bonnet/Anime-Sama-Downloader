import React from 'react';

export interface JobRowProps {
  id: string;
  animeTitle: string;
  episode: number;
  progress: number; // 0-100
  eta?: string;
  speed?: string;
  status: 'queued' | 'downloading' | 'paused' | 'completed' | 'failed';
  onPause?: () => void;
  onResume?: () => void;
  onCancel: () => void;
  onRetry?: () => void;
}

const JobRow: React.FC<JobRowProps> = ({
  id,
  animeTitle,
  episode,
  progress,
  eta,
  speed,
  status,
  onPause,
  onResume,
  onCancel,
  onRetry,
}) => {
  const statusColors: Record<string, string> = {
    queued: 'var(--night-text-secondary)',
    downloading: 'var(--night-bg-warning)',
    paused: 'var(--night-bg-info)',
    completed: 'var(--night-bg-success)',
    failed: 'var(--night-bg-error)',
  };

  const statusLabels: Record<string, string> = {
    queued: 'En attente',
    downloading: 'Téléchargement',
    paused: 'Mis en pause',
    completed: 'Terminé',
    failed: 'Erreur',
  };

  const showActions = status !== 'queued' && status !== 'completed';

  return (
    <div
      className={`job-row job-row--${status}`}
      data-testid={`job-row-${id}`}
      role="region"
      aria-label={`Tâche: ${animeTitle} Ép. ${episode}`}
    >
      {/* Title + Episode */}
      <div className="job-row__info">
        <h4 className="job-row__title">{animeTitle}</h4>
        <span className="job-row__episode">Ép. {episode}</span>
      </div>

      {/* Progress Section */}
      <div className="job-row__progress-section">
        <div className="job-row__progress-bar">
          <div
            className="job-row__progress-fill"
            style={{ width: `${progress}%` }}
            role="progressbar"
            aria-valuenow={progress}
            aria-valuemin={0}
            aria-valuemax={100}
            aria-label={`Progression: ${progress}%`}
          ></div>
        </div>
        <span className="job-row__progress-text">{progress}%</span>
      </div>

      {/* Meta: ETA + Speed */}
      <div className="job-row__meta">
        {eta && (
          <span className="job-row__eta" title={`Temps estimé: ${eta}`}>
            {eta}
          </span>
        )}
        {speed && (
          <span className="job-row__speed" title={`Vitesse: ${speed}`}>
            {speed}
          </span>
        )}
      </div>

      {/* Status Badge */}
      <div className="job-row__status">
        <span
          className="job-row__badge"
          style={{ backgroundColor: statusColors[status] }}
        >
          {statusLabels[status]}
        </span>
      </div>

      {/* Actions (conditional by status) */}
      {showActions && (
        <div className="job-row__actions">
          {status === 'downloading' && onPause && (
            <button
              className="job-row__action-btn job-row__action-btn--pause"
              onClick={onPause}
              aria-label={`Mettre en pause ${animeTitle} Ép. ${episode}`}
            >
              ⏸
            </button>
          )}
          {status === 'paused' && onResume && (
            <button
              className="job-row__action-btn job-row__action-btn--resume"
              onClick={onResume}
              aria-label={`Reprendre ${animeTitle} Ép. ${episode}`}
            >
              ▶
            </button>
          )}
          {status === 'downloading' && (
            <button
              className="job-row__action-btn job-row__action-btn--cancel"
              onClick={onCancel}
              aria-label={`Annuler ${animeTitle} Ép. ${episode}`}
            >
              ✕
            </button>
          )}
          {status === 'paused' && (
            <button
              className="job-row__action-btn job-row__action-btn--cancel"
              onClick={onCancel}
              aria-label={`Annuler ${animeTitle} Ép. ${episode}`}
            >
              ✕
            </button>
          )}
          {status === 'failed' && onRetry && (
            <button
              className="job-row__action-btn job-row__action-btn--retry"
              onClick={onRetry}
              aria-label={`Réessayer ${animeTitle} Ép. ${episode}`}
            >
              ↻
            </button>
          )}
          {status === 'failed' && (
            <button
              className="job-row__action-btn job-row__action-btn--cancel"
              onClick={onCancel}
              aria-label={`Supprimer ${animeTitle} Ép. ${episode}`}
            >
              ✕
            </button>
          )}
        </div>
      )}

      {/* Live region for status updates */}
      <div aria-live="polite" aria-atomic="true" className="sr-only">
        {status === 'downloading' && `${progress}% téléchargé`}
        {status === 'completed' && `Téléchargement terminé`}
        {status === 'failed' && `Erreur lors du téléchargement`}
      </div>
    </div>
  );
};

JobRow.displayName = 'JobRow';

export default JobRow;
