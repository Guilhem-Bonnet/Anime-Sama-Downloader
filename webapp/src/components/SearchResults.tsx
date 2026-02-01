import React, { useEffect } from 'react';
import { useSearchStore } from '../stores/search.store';
import { useJobsStore } from '../stores/jobs.store';
import { apiClient } from '../utils/api';
import { Download } from 'lucide-react';
import { EmptySearchIllustration } from './illustrations/SakuraIllustrations';

export const SearchResultsGrid: React.FC = () => {
  const { results, isSearching } = useSearchStore();
  const { addJob } = useJobsStore();
  const [showStamp, setShowStamp] = React.useState(false);

  // TEMPORARY: Show mock data if no results
  const mockResults = [
    {anime_id: 'mushishi', title: 'Mushishi', episodes: 26, source: 'AnimeSama', image_url: '/assets/cover-placeholder.svg'},
    {anime_id: 'mononoke', title: 'Mononoke', episodes: 12, source: 'AnimeSama', image_url: '/assets/cover-placeholder.svg'},
    {anime_id: 'natsume-yuujinchou', title: 'Natsume Yuujinchou', episodes: 13, source: 'AnimeSama', image_url: '/assets/cover-placeholder.svg'},
    {anime_id: 'samurai-champloo', title: 'Samurai Champloo', episodes: 26, source: 'AnimeSama', image_url: '/assets/cover-placeholder.svg'},
  ];
  const displayResults = results.length > 0 ? results : mockResults;

  const handleDownload = async (animeId: string, title: string) => {
    try {
      setShowStamp(true);
      setTimeout(() => setShowStamp(false), 1400);
      const download = await apiClient.createDownload(animeId, 1);
      if (download && download.downloadId) {
        addJob({
          id: download.downloadId,
          status: 'pending',
          progress: 0,
          downloadId: download.downloadId,
          animeId,
          episodeNumber: 1,
        });
      }
    } catch (error) {
      console.error('Failed to create download:', error);
    }
  };

  if (isSearching) {
    return (
      <div className="text-center py-20">
        <div className="inline-flex items-center justify-center gap-4">
          <div className="relative">
            <div className="w-16 h-16 border-4 border-gray-200 dark:border-gray-700 rounded-full"></div>
            <div className="absolute top-0 left-0 w-16 h-16 border-4 border-transparent border-t-cyan-500 border-r-magenta-500 rounded-full animate-spin"></div>
          </div>
        </div>
        <p className="text-gray-600 dark:text-gray-400 mt-6 text-lg font-medium animate-pulse">
          Recherche en cours...
        </p>
      </div>
    );
  }

  if (displayResults.length === 0 && !isSearching) {
    return (
      <div className="text-center py-20">
        <EmptySearchIllustration />
        <p className="text-gray-600 dark:text-gray-400 text-lg mt-4">
          Aucun résultat trouvé. Essayez une autre recherche !
        </p>
      </div>
    );
  }

  return (
    <div className="space-y-8">
      {/* Résultats */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
        <div className={`ink-stamp ${showStamp ? 'is-visible' : ''}`} aria-hidden="true" />
        {displayResults.map((result, index) => (
          <div
            key={result.anime_id}
            className="group relative frame-ornate overflow-hidden shadow-lg hover:shadow-2xl transition-all duration-300 hover:scale-[1.02] animate-fadeInUp"
            style={{ 
              animationDelay: `${index * 100}ms`,
              background: 'linear-gradient(135deg, rgba(143,106,61,0.1), rgba(125,114,103,0.05))',
            }}
          >
            {/* Image avec overlay gradient */}
            {result.image_url && (
              <div className="relative h-64 overflow-hidden">
                <img
                  src={result.image_url}
                  alt={result.title}
                  className="w-full h-full object-cover transition-transform duration-300 group-hover:scale-110"
                />
                <div className="absolute inset-0 bg-gradient-to-t from-black/80 via-black/40 to-transparent" />
                {/* Badge épisodes avec encre brune */}
                <div style={{
                  position: 'absolute',
                  top: '12px',
                  right: '12px',
                  background: 'rgba(20,17,22,0.85)',
                  backdropFilter: 'blur(8px)',
                  padding: '6px 12px',
                  borderRadius: '999px',
                  fontSize: '12px',
                  fontWeight: 600,
                  color: 'var(--sakura-text-light)',
                  border: '1px solid rgba(255,255,255,0.15)',
                }}>
                  {result.episodes} épisodes
                </div>
              </div>
            )}
            
            {/* Contenu */}
            <div style={{ padding: '16px' }}>
              <h3 style={{
                fontWeight: 700,
                fontSize: '16px',
                color: 'var(--sakura-text-primary)',
                marginBottom: '12px',
                lineHeight: '1.3',
                display: '-webkit-box',
                WebkitLineClamp: 2,
                WebkitBoxOrient: 'vertical',
                overflow: 'hidden',
              }}>
                {result.title}
              </h3>
              <p style={{
                fontSize: '13px',
                color: 'var(--sakura-text-secondary)',
                marginBottom: '16px',
                display: 'flex',
                alignItems: 'center',
                gap: '8px',
              }}>
                <span style={{
                  display: 'inline-block',
                  width: '6px',
                  height: '6px',
                  borderRadius: '50%',
                  background: 'var(--sakura-accent-brown-500)',
                }}></span>
                {result.source}
              </p>
              
              {/* Bouton Download avec encre brune */}
              <button
                onClick={() => handleDownload(result.anime_id, result.title)}
                style={{
                  width: '100%',
                  padding: '10px 12px',
                  borderRadius: '12px',
                  border: '1.5px solid var(--sakura-accent-brown-500)',
                  background: 'var(--sakura-accent-brown-500)',
                  color: 'white',
                  fontWeight: 600,
                  fontSize: '13px',
                  cursor: 'pointer',
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                  gap: '8px',
                  transition: 'all 200ms ease',
                }}
                onMouseEnter={(e) => {
                  e.currentTarget.style.background = 'var(--sakura-accent-brown-600)';
                  e.currentTarget.style.boxShadow = '0 4px 12px rgba(143, 106, 61, 0.3)';
                }}
                onMouseLeave={(e) => {
                  e.currentTarget.style.background = 'var(--sakura-accent-brown-500)';
                  e.currentTarget.style.boxShadow = 'none';
                }}
              >
                <Download className="w-4 h-4" />
                <span>Télécharger</span>
              </button>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};
