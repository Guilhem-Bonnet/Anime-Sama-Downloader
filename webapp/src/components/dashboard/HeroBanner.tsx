import React from 'react';

interface HeroBannerProps {
  jobs?: any[];
  subscriptions?: any[];
  trendingAnime?: any;
}

const HeroBanner: React.FC<HeroBannerProps> = ({
  jobs = [],
  subscriptions = [],
  trendingAnime,
}) => {
  // Fallback strategy
  let anime = trendingAnime;
  let ctaText = 'Découvrir';
  let coverUrl = trendingAnime?.coverImage?.large;

  if (jobs.length > 0 && jobs[0]?.animeTitle) {
    anime = jobs[0];
    ctaText = 'Reprendre';
    coverUrl = undefined;
  } else {
    // Check subscriptions with new episodes available
    const subWithNew = subscriptions.find((s: any) => {
      const avail = s.lastAvailableEpisode || 0;
      const downloaded = s.lastDownloadedEpisode || 0;
      return avail > downloaded;
    });
    if (subWithNew) {
      const newCount = (subWithNew.lastAvailableEpisode || 0) - (subWithNew.lastDownloadedEpisode || 0);
      anime = subWithNew;
      ctaText = `Voir ${newCount} nouveaux`;
      coverUrl = undefined;
    }
  }

  const title = anime?.title?.romaji || anime?.animeTitle || anime?.label || 'Bienvenue';
  const season = anime?.season || '';
  const year = anime?.seasonYear || '';
  const language = anime?.language || '';

  return (
    <div 
      className="hero-banner"
      style={{
        backgroundImage: coverUrl ? `linear-gradient(135deg, rgba(10, 14, 26, 0.7) 0%, rgba(26, 31, 46, 0.7) 100%), url('${coverUrl}')` : 'linear-gradient(135deg, var(--night-bg-base), var(--night-bg-elevated))',
        backgroundSize: 'cover',
        backgroundPosition: 'center',
      }}
    >
      <div className="hero-banner__content">
        <h2 className="hero-banner__title">{title}</h2>
        {(season || year) && (
          <div className="hero-banner__meta">
            {season && <span className="hero-banner__badge">{season}</span>}
            {year && <span className="hero-banner__badge">{year}</span>}
            {language && <span className="hero-banner__badge">{language}</span>}
          </div>
        )}
        <button className="hero-banner__cta">{ctaText}</button>
      </div>
    </div>
  );
};

HeroBanner.displayName = 'HeroBanner';

export default HeroBanner;
