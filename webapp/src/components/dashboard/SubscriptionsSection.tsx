import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { SubscriptionCard } from '../anime/index';

interface SubscriptionsSectionProps {
  subscriptions: any[];
}

/**
 * Fetch cover images from AniList by anime title, returns a map title → coverUrl.
 */
async function fetchCoversFromAniList(titles: string[]): Promise<Record<string, string>> {
  const covers: Record<string, string> = {};
  if (titles.length === 0) return covers;

  // Deduplicate and limit (avoid hammering AniList)
  const uniqueTitles = [...new Set(titles)].slice(0, 20);

  await Promise.allSettled(
    uniqueTitles.map(async (title) => {
      try {
        const res = await fetch('https://graphql.anilist.co', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            query: `query($search:String){Page(page:1,perPage:1){media(search:$search,type:ANIME){title{romaji english}coverImage{large medium}}}}`,
            variables: { search: title },
          }),
        });
        const data = await res.json();
        const media = data?.data?.Page?.media?.[0];
        if (media) {
          const url = media.coverImage?.large || media.coverImage?.medium || '';
          if (url) covers[title] = url;
        }
      } catch {
        /* best-effort */
      }
    })
  );

  return covers;
}

const SubscriptionsSection: React.FC<SubscriptionsSectionProps> = ({ subscriptions }) => {
  const navigate = useNavigate();
  const [coverMap, setCoverMap] = useState<Record<string, string>>({});

  useEffect(() => {
    const titles = subscriptions
      .map((s) => s.label || '')
      .filter(Boolean);
    if (titles.length > 0) {
      fetchCoversFromAniList(titles).then(setCoverMap);
    }
  }, [subscriptions]);

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
            coverUrl={coverMap[sub.label] || ''}
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
