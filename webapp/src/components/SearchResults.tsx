import React from 'react';
import { useSearchStore } from '../stores/search.store';
import { Download } from 'lucide-react';
import { EmptySearchIllustration } from './illustrations/NocturneIllustrations';
import { useNavigate } from 'react-router-dom';

export const SearchResultsGrid: React.FC = () => {
  const { results, isSearching } = useSearchStore();
  const navigate = useNavigate();
  const [showStamp, setShowStamp] = React.useState(false);

  const displayResults = results;

  const handleViewDetail = (animeId: string) => {
    navigate(`/anime/${animeId}`);
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
            key={result.id}
            className="group relative frame-ornate overflow-hidden shadow-lg hover:shadow-2xl transition-all duration-300 hover:scale-[1.02] animate-fadeInUp"
            style={{ 
              animationDelay: `${Math.min(index * 100, 500)}ms`,
              background: 'linear-gradient(135deg, rgba(143,106,61,0.1), rgba(125,114,103,0.05))',
            }}
          >
            {/* Image avec overlay gradient */}
            {result.thumbnail_url && (
              <div className="relative h-64 overflow-hidden">
                <img
                  src={result.thumbnail_url}
                  alt={result.title}
                  className="w-full h-full object-cover transition-transform duration-300 group-hover:scale-110"
                  onError={(e) => {
                    const target = e.currentTarget;
                    if (target.src !== '/assets/cover-placeholder.svg') {
                      target.src = '/assets/cover-placeholder.svg';
                    }
                  }}
                />
                <div className="absolute inset-0 bg-gradient-to-t from-black/80 via-black/40 to-transparent" />
                {/* Badge épisodes */}
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
                  color: 'var(--night-text-light)',
                  border: '1px solid rgba(255,255,255,0.15)',
                }}>
                  {result.episode_count || '?'} épisodes
                </div>
              </div>
            )}
            
            {/* Contenu */}
            <div style={{ padding: '16px' }}>
              <h3 style={{
                fontWeight: 700,
                fontSize: '16px',
                color: 'var(--night-text-primary)',
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
                color: 'var(--night-text-secondary)',
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
                  background: 'var(--night-accent-brown-500)',
                }}></span>
                {result.year || ''} · {result.status || ''}
              </p>
              
              {/* Bouton Voir détails */}
              <button
                onClick={() => handleViewDetail(result.id)}
                style={{
                  width: '100%',
                  padding: '10px 12px',
                  borderRadius: '12px',
                  border: '1.5px solid var(--night-accent-brown-500)',
                  background: 'var(--night-accent-brown-500)',
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
                  e.currentTarget.style.background = 'var(--night-accent-brown-600)';
                  e.currentTarget.style.boxShadow = '0 4px 12px rgba(143, 106, 61, 0.3)';
                }}
                onMouseLeave={(e) => {
                  e.currentTarget.style.background = 'var(--night-accent-brown-500)';
                  e.currentTarget.style.boxShadow = 'none';
                }}
              >
                <Download className="w-4 h-4" />
                <span>Voir détails</span>
              </button>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};
