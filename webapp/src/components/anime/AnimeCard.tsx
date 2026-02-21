import React from 'react';

export interface AnimeCardProps {
  id: string;
  title: string;
  coverUrl: string;
  season?: string;
  language: 'VOSTFR' | 'VF';
  status: 'ongoing' | 'completed' | 'upcoming';
  onDetails: () => void;
  onDownload: () => void;
  variant?: 'grid' | 'list';
  disabled?: boolean;
  selected?: boolean;
}

const AnimeCard: React.FC<AnimeCardProps> = ({
  id,
  title,
  coverUrl,
  season,
  language,
  status,
  onDetails,
  onDownload,
  variant = 'grid',
  disabled = false,
  selected = false,
}) => {
  const statusColors: Record<string, string> = {
    ongoing: 'var(--night-bg-info)',
    completed: 'var(--night-bg-success)',
    upcoming: 'var(--night-bg-warning)',
  };

  const statusLabels: Record<string, string> = {
    ongoing: 'En cours',
    completed: 'Terminé',
    upcoming: 'À venir',
  };

  const styles = variant === 'list' 
    ? 'anime-card--list'
    : 'anime-card--grid';

  return (
    <article
      className={`anime-card ${styles} ${selected ? 'anime-card--selected' : ''} ${disabled ? 'anime-card--disabled' : ''}`}
      aria-label={`${title}${season ? ` - ${season}` : ''}`}
      role="article"
      data-testid={`anime-card-${id}`}
    >
      {/* Cover Image */}
      <div className="anime-card__cover-wrapper">
        <img
          src={coverUrl}
          alt={`Couverture de ${title}`}
          className="anime-card__cover"
          loading="lazy"
        />
        <div 
          className="anime-card__status-badge"
          style={{ backgroundColor: statusColors[status] }}
          aria-label={`Statut: ${statusLabels[status]}`}
        >
          {statusLabels[status]}
        </div>
      </div>

      {/* Title */}
      <h3 className="anime-card__title">{title}</h3>

      {/* Badges */}
      <div className="anime-card__badges">
        {season && (
          <span className="anime-card__badge anime-card__badge--season">
            {season}
          </span>
        )}
        <span className="anime-card__badge anime-card__badge--language">
          {language}
        </span>
      </div>

      {/* Actions */}
      <div className="anime-card__actions">
        <button
          className="anime-card__action-btn anime-card__action-btn--details"
          onClick={onDetails}
          disabled={disabled}
          aria-label={`Voir détails de ${title}`}
        >
          Détails
        </button>
        <button
          className="anime-card__action-btn anime-card__action-btn--download"
          onClick={onDownload}
          disabled={disabled}
          aria-label={`Télécharger ${title}`}
        >
          Télécharger
        </button>
      </div>
    </article>
  );
};

AnimeCard.displayName = 'AnimeCard';

export default AnimeCard;
