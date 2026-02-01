import React from 'react';
import { useJobsStore } from '../stores/jobs.store';
import { StatusBadge } from './custom/StatusBadge';
import { DownloadProgress } from './custom/DownloadProgress';

export const DownloadList: React.FC = () => {
  const { jobs } = useJobsStore();

  if (jobs.length === 0) {
    return (
      <div className="text-center py-12">
        <p className="text-gray-500 dark:text-gray-400">No downloads yet</p>
      </div>
    );
  }

  return (
    <div className="space-y-4">
      {jobs.map((job) => (
        <div
          key={job.id}
          className="p-4 bg-gray-100 dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700"
        >
          <div className="flex justify-between items-start mb-2">
            <div>
              <h3 className="font-semibold text-gray-900 dark:text-white">{job.animeId}</h3>
              <p className="text-sm text-gray-600 dark:text-gray-400">
                Episode {job.episodeNumber}
              </p>
            </div>
            <StatusBadge status={job.status} />
          </div>
          <DownloadProgress progress={job.progress} />
          {job.errorMessage && (
            <p className="text-red-500 text-sm mt-2">{job.errorMessage}</p>
          )}
        </div>
      ))}
    </div>
  );
};
