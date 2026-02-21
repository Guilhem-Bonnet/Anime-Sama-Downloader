import React, { useState, useEffect } from 'react';
import { TrendingUp } from 'lucide-react';

export interface RecommendationScore {
  anime_id: string;
  title: string;
  score: number;
  reason?: string;
}

interface YouMightAlsoLikeProps {
  query: string;
  limit?: number;
  onSelectAnime?: (animeId: string, title: string) => void;
}

export const YouMightAlsoLike: React.FC<YouMightAlsoLikeProps> = ({
  query,
  limit = 6,
  onSelectAnime,
}) => {
  const [recommendations, setRecommendations] = useState<RecommendationScore[]>([]);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    if (!query.trim()) {
      setRecommendations([]);
      return;
    }

    fetchRecommendations();
  }, [query]);

  const fetchRecommendations = async () => {
    setLoading(true);

    try {
      const params = new URLSearchParams();
      params.append('q', query);
      params.append('limit', limit.toString());

      const response = await fetch(`/api/v1/recommendations/query?${params.toString()}`);
      if (!response.ok) throw new Error('Failed to fetch recommendations');

      const data = await response.json();
      setRecommendations(data.recommendations || []);
    } catch (err) {
      console.error('Failed to fetch recommendations:', err);
      setRecommendations([]);
    } finally {
      setLoading(false);
    }
  };

  if (!query.trim() || (recommendations.length === 0 && !loading)) {
    return null;
  }

  return (
    <div className="mt-8 py-6 px-4 bg-gradient-to-r from-cyan-50 to-blue-50 dark:from-gray-800 dark:to-gray-900 rounded-lg border border-cyan-200 dark:border-gray-700">
      <div className="flex items-center gap-2 mb-4">
        <TrendingUp className="w-5 h-5 text-cyan-600 dark:text-cyan-400" />
        <h2 className="text-lg font-semibold text-gray-900 dark:text-white">
          You Might Also Like
        </h2>
      </div>

      {loading && (
        <div className="text-center text-gray-600 dark:text-gray-400 py-4">
          Finding recommendations...
        </div>
      )}

      {recommendations.length > 0 && (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-3">
          {recommendations.map((rec) => (
            <button
              key={rec.anime_id}
              onClick={() => onSelectAnime?.(rec.anime_id, rec.title)}
              className="p-3 bg-white dark:bg-gray-800 rounded-lg hover:shadow-md transition-all border border-gray-200 dark:border-gray-600 hover:border-cyan-400 dark:hover:border-cyan-500 text-left group"
            >
              <div className="flex items-start justify-between gap-2">
                <div className="flex-1 min-w-0">
                  <h3 className="text-sm font-medium text-gray-900 dark:text-white group-hover:text-cyan-600 dark:group-hover:text-cyan-400 truncate transition-colors">
                    {rec.title}
                  </h3>
                  {rec.reason && (
                    <p className="text-xs text-gray-600 dark:text-gray-400 mt-1">
                      {rec.reason}
                    </p>
                  )}
                </div>
                <div className="flex-shrink-0 text-sm font-semibold bg-cyan-100 dark:bg-cyan-900 text-cyan-700 dark:text-cyan-300 px-2 py-1 rounded">
                  {rec.score.toFixed(1)}
                </div>
              </div>
            </button>
          ))}
        </div>
      )}
    </div>
  );
};
