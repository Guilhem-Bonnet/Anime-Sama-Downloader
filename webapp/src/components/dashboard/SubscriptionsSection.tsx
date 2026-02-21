import React from 'react';
import { useNavigate } from 'react-router-dom';
import { SubscriptionCard } from '../anime/index';

interface SubscriptionsSectionProps {
  subscriptions: any[];
}

const SubscriptionsSection: React.FC<SubscriptionsSectionProps> = ({ subscriptions }) => {
  const navigate = useNavigate();

  if (!subscriptions || subscriptions.length === 0) {
    return (
      <div className="empty-state">
        <div className="empty-state__icon">📺</div>
        <h3 className="empty-state__title">Aucun abonnement actif</h3>
        <p className="empty-state__description">
          Abonnez-vous à des séries pour récupérer automatiquement les nouveaux épisodes
        </p>
        <button 
          className="empty-state__cta"
          onClick={() => navigate('/search')}
        >
          Commencer maintenant
        </button>
      </div>
    );
  }

  return (
    <div className="subscriptions-grid">
      {subscriptions.map((sub) => (
        <div key={sub.id} className="subscription-card-wrapper">
          {sub.lastAvailableEpisode > sub.lastDownloadedEpisode && (
            <div className="subscription-badge">
              {sub.lastAvailableEpisode - sub.lastDownloadedEpisode} nouveau
            </div>
          )}
          <SubscriptionCard
            id={sub.id}
            animeTitle={sub.label || sub.baseUrl}
            season={1}
            coverUrl=""
            lastEpisode={sub.lastDownloadedEpisode || 0}
            nextEpisode={sub.lastAvailableEpisode || undefined}
          />
        </div>
      ))}
    </div>
  );
};

SubscriptionsSection.displayName = 'SubscriptionsSection';

export default SubscriptionsSection;
