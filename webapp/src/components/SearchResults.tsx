import React, { useEffect } from 'react';
import { useSearchStore } from '../stores/search.store';
import { useJobsStore } from '../stores/jobs.store';
import { apiClient } from '../utils/api';
import { Download, Play } from 'lucide-react';

export const SearchResultsGrid: React.FC = () => {
  const { results, isSearching } = useSearchStore();
  const { addJob } = useJobsStore();

  // TEMPORARY: Show mock data if no results
  const mockResults = [
    {anime_id: 'test-1', title: 'Attack on Titan', episodes: 75, source: 'AnimeSama', image_url: 'https://cdn.myanimelist.net/images/anime/10/47347.jpg'},
    {anime_id: 'test-2', title: 'Demon Slayer', episodes: 26, source: 'AnimeSama', image_url: 'https://cdn.myanimelist.net/images/anime/1286/99889.jpg'},
  ];
  const displayResults = results.length > 0 ? results : mockResults;

  const handleDownload = async (animeId: string, title: string) => {
    try {
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
        <div className="flex justify-center mb-4">
          <div className="p-6 bg-gray-100 dark:bg-gray-800 rounded-full">
            <Play className="w-12 h-12 text-gray-400" />
          </div>
        </div>
        <p className="text-gray-600 dark:text-gray-400 text-lg">
          Aucun résultat trouvé. Essayez une autre recherche !
        </p>
      </div>
    );
  }

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
      {displayResults.map((result, index) => (
        <div
          key={result.anime_id}
          className="group relative bg-gradient-to-br from-white to-gray-50 dark:from-gray-800 dark:to-gray-900 rounded-2xl overflow-hidden shadow-lg hover:shadow-2xl transition-all duration-300 border border-gray-200/50 dark:border-gray-700/50 hover:scale-[1.02] hover:border-cyan-500/50 animate-fadeInUp"
          style={{ animationDelay: `${index * 100}ms` }}
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
              {/* Badge épisodes */}
              <div className="absolute top-3 right-3 bg-black/70 backdrop-blur-sm px-3 py-1 rounded-full text-xs font-medium text-white border border-white/20">
                {result.episodes} eps
              </div>
            </div>
          )}
          
          {/* Contenu */}
          <div className="p-5">
            <h3 className="font-bold text-lg text-gray-900 dark:text-white mb-2 line-clamp-2 group-hover:text-transparent group-hover:bg-clip-text group-hover:bg-gradient-to-r group-hover:from-magenta-500 group-hover:to-cyan-500 transition-all">
              {result.title}
            </h3>
            <p className="text-sm text-gray-600 dark:text-gray-400 mb-4 flex items-center gap-2">
              <span className="inline-block w-2 h-2 rounded-full bg-cyan-500"></span>
              {result.source}
            </p>
            
            {/* Bouton Download avec gradient */}
            <button
              onClick={() => handleDownload(result.anime_id, result.title)}
              className="w-full py-3 px-4 bg-gradient-to-r from-magenta-600 to-pink-600 hover:from-magenta-700 hover:to-pink-700 text-white font-semibold rounded-xl shadow-lg shadow-magenta-500/30 hover:shadow-xl hover:shadow-magenta-500/50 transition-all duration-300 flex items-center justify-center gap-2"
            >
              <Download className="w-5 h-5" />
              <span>Download</span>
            </button>
          </div>
        </div>
      ))}
    </div>
  );
};
