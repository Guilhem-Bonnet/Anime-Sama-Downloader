import React, { useEffect } from 'react';
import { useSearchStore } from '../../stores/search.store';
import { useJobsStore } from '../../stores/jobs.store';
import { apiClient } from '../../utils/api';

export const SearchResultsGrid: React.FC = () => {
  const { results, isSearching } = useSearchStore();
  const { addJob } = useJobsStore();

  const handleDownload = async (animeId: string, title: string) => {
    try {
      const download = await apiClient.createDownload(animeId, 1);
      if (download) {
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
      <div className="text-center py-12">
        <div className="inline-block animate-spin text-2xl">⌛</div>
        <p className="text-gray-500 dark:text-gray-400 mt-4">Searching...</p>
      </div>
    );
  }

  if (results.length === 0) {
    return (
      <div className="text-center py-12">
        <p className="text-gray-500 dark:text-gray-400">No results found. Try another search.</p>
      </div>
    );
  }

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
      {results.map((result) => (
        <div
          key={result.animeId}
          className="p-4 bg-gray-100 dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 hover:shadow-lg transition-shadow"
        >
          {result.imageUrl && (
            <img
              src={result.imageUrl}
              alt={result.title}
              className="w-full h-32 object-cover rounded mb-3"
            />
          )}
          <h3 className="font-semibold text-gray-900 dark:text-white">{result.title}</h3>
          <p className="text-sm text-gray-600 dark:text-gray-400 mt-1">
            {result.episodes} episodes • {result.source}
          </p>
          <button
            onClick={() => handleDownload(result.animeId, result.title)}
            className="w-full mt-3 px-4 py-2 bg-magenta-600 hover:bg-magenta-700 text-white rounded transition-colors"
          >
            Download
          </button>
        </div>
      ))}
    </div>
  );
};
