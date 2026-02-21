import React from 'react';

export interface SubscriptionCardProps {
  id: string;
  animeTitle: string;
  season: number;
  coverUrl?: string;
  lastEpisode: number;
  nextEpisode?: number;
  language?: 'VOSTFR' | 'VF';
  quality?: string;
  interval?: string; // e.g., "Hebdomadaire"
  onClick?: () => void;
}

const SubscriptionCard: React.FC<SubscriptionCardProps> = ({
  animeTitle,
  season,
  coverUrl,
  lastEpisode,
  nextEpisode,
  language,
  quality,
  interval,
  onClick,
}) => {
  return (
    <article
      className="subscription-card"
      onClick={onClick}
      role="button"
      tabIndex={0}
      onKeyDown={(e) => {
        if ((e.key === 'Enter' || e.key === ' ') && onClick) {
          e.preventDefault();
          onClick();
        }
      }}
      aria-label={`Abonnement: ${animeTitle} Saison ${season}`}
    >
      <div className="subscription-card__cover">
        {coverUrl ? (
          <img src={coverUrl} alt={`${animeTitle} cover`} loading="lazy" />
        ) : (
          <div className="subscription-card__cover-placeholder">
            <span>📺</span>
          </div>
        )}
      </div>

      <div className="subscription-card__content">
        <h3 className="subscription-card__title">{animeTitle}</h3>
        
        <div className="subscription-card__meta">
          <span className="subscription-card__season">Saison {season}</span>
          {language && (
            <span className="subscription-card__badge">{language}</span>
          )}
        </div>

        <div className="subscription-card__progress">
          <span className="subscription-card__episode">
            Épisode {lastEpisode}
          </span>
          {nextEpisode && (
            <span className="subscription-card__next">
              Prochain: {nextEpisode}
            </span>
          )}
        </div>

        {(quality || interval) && (
          <div className="subscription-card__details">
            {quality && <span>{quality}</span>}
            {interval && <span>{interval}</span>}
          </div>
        )}
      </div>
    </article>
  );
};

SubscriptionCard.displayName = 'SubscriptionCard';

export default SubscriptionCard;
