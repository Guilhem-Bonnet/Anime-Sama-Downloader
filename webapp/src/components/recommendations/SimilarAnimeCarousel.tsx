import React, { useState, useEffect } from 'react';
import { Star, Zap } from 'lucide-react';

export interface RecommendationScore {
  anime_id: string;
  title: string;
  score: number;
  reason?: string;
}

interface SimilarAnimeCarouselProps {
  animeId: string;
  title: string;
  limit?: number;
  onSelectAnime?: (animeId: string) => void;
}

export const SimilarAnimeCarousel: React.FC<SimilarAnimeCarouselProps> = ({
  animeId,
  title,
  limit = 10,
  onSelectAnime,
}) => {
  const [recommendations, setRecommendations] = useState<RecommendationScore[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetchRecommendations();
  }, [animeId]);

  const fetchRecommendations = async () => {
    setLoading(true);
    setError(null);

    try {
      const params = new URLSearchParams();
      params.append('anime_id', animeId);
      params.append('limit', limit.toString());

      const response = await fetch(`/api/v1/recommendations/similar?${params.toString()}`);
      if (!response.ok) throw new Error('Failed to fetch recommendations');

      const data = await response.json();
      setRecommendations(data.recommendations || []);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Unknown error');
      setRecommendations([]);
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <div className="py-6 px-4 bg-gray-50 dark:bg-gray-900 rounded-lg">
        <h2 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">
          Similar Anime
        </h2>
        <div className="text-center text-gray-500 dark:text-gray-400">Loading...</div>
      </div>
    );
  }

  if (error || recommendations.length === 0) {
    return null;
  }

  return (
    <div className="py-6 px-4 bg-gray-50 dark:bg-gray-900 rounded-lg">
      <h2 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">
        Similar Anime
      </h2>

      <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-5 gap-4">
        {recommendations.map((rec) => (
          <button
            key={rec.anime_id}
            onClick={() => onSelectAnime?.(rec.anime_id)}
            className="p-3 bg-white dark:bg-gray-800 rounded-lg hover:shadow-lg transition-shadow border border-gray-200 dark:border-gray-700 text-left"
          >
            <div className="text-sm font-medium text-gray-900 dark:text-white truncate mb-2">
              {rec.title}
            </div>

            <div className="flex items-center gap-1 mb-2">
              <Zap className="w-3 h-3 text-yellow-500" />
              <span className="text-xs text-gray-600 dark:text-gray-400">
                {rec.score.toFixed(1)}
              </span>
            </div>

            {rec.reason && (
              <div className="text-xs text-gray-500 dark:text-gray-500">
                {rec.reason}
              </div>
            )}
          </button>
        ))}
      </div>
    </div>
  );
};
