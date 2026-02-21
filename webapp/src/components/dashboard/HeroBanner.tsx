import React, { useEffect, useState } from 'react';

interface HeroBannerProps {
  jobs?: any[];
  subscriptions?: any[];
  trendingAnime?: any;
}

/**
 * Resolve a cover image URL from AniList by anime title (best-effort).
 */
async function resolveCoverByTitle(title: string): Promise<string | undefined> {
  try {
    const res = await fetch('https://graphql.anilist.co', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        query: `query($search:String){Page(page:1,perPage:1){media(search:$search,type:ANIME){coverImage{large}}}}`,
        variables: { search: title },
      }),
    });
    const data = await res.json();
    return data?.data?.Page?.media?.[0]?.coverImage?.large || undefined;
  } catch {
    return undefined;
  }
}

const HeroBanner: React.FC<HeroBannerProps> = ({
  jobs = [],
  subscriptions = [],
  trendingAnime,
}) => {
  const [resolvedCover, setResolvedCover] = useState<string | undefined>(undefined);

  // Fallback strategy for displayed anime info
  let anime = trendingAnime;
  let ctaText = 'Découvrir';
  let titleToResolve: string | undefined;

  if (jobs.length > 0 && jobs[0]?.animeTitle) {
    anime = jobs[0];
    ctaText = 'Reprendre';
    titleToResolve = jobs[0].animeTitle;
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
      titleToResolve = subWithNew.label;
    }
  }

  // Resolve cover image from AniList when we have a title but no trending image
  useEffect(() => {
    if (titleToResolve) {
      resolveCoverByTitle(titleToResolve).then(setResolvedCover);
    }
  }, [titleToResolve]);

  // Cover priority: resolved cover for current anime → trending anime cover → gradient
  const coverUrl = titleToResolve
    ? resolvedCover
    : trendingAnime?.coverImage?.large;

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
