import React, { useEffect, useState } from 'react';
import { useJobsStore } from '../stores/jobs.store';
import { useSSE } from '../hooks/useSSE';
import { StatusBadge } from './custom/StatusBadge';
import { DownloadProgress } from './custom/DownloadProgress';

export const DownloadMonitor: React.FC = () => {
  const { jobs, updateJobProgress, updateJobStatus } = useJobsStore();
  const [sseLogs, setSSELogs] = useState<string[]>([]);

  // Subscribe to SSE updates for first active job
  const activeJob = jobs.find((j) => j.status === 'running');

  const { close } = useSSE(
    activeJob ? `/api/jobs/${activeJob.id}/progress` : '',
    activeJob
      ? (data) => {
          updateJobProgress(activeJob.id, data.progress || 0);
          setSSELogs((prev) => [
            ...prev,
            `[${new Date().toLocaleTimeString()}] Job ${activeJob.id}: ${data.progress}% complete`,
          ]);

          if (data.progress >= 100) {
            updateJobStatus(activeJob.id, 'completed');
          }
        }
      : undefined
  );

  if (jobs.length === 0) {
    return (
      <div className="text-center py-12">
        <p className="text-gray-500 dark:text-gray-400">No downloads</p>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Active Downloads */}
      <div>
        <h3 className="text-lg font-semibold mb-4 text-gray-900 dark:text-white">
          Active Downloads ({jobs.filter((j) => j.status === 'running').length})
        </h3>
        <div className="space-y-3">
          {jobs
            .filter((j) => j.status === 'running' || j.status === 'pending')
            .map((job) => (
              <div
                key={job.id}
                className="p-4 bg-gray-100 dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700"
              >
                <div className="flex justify-between items-start mb-2">
                  <div>
                    <h4 className="font-semibold text-gray-900 dark:text-white">
                      {job.animeId} - Ep {job.episodeNumber}
                    </h4>
                    <p className="text-xs text-gray-500">{job.id}</p>
                  </div>
                  <StatusBadge status={job.status} size="sm" />
                </div>
                <DownloadProgress progress={job.progress} />
              </div>
            ))}
        </div>
      </div>

      {/* Completed Downloads */}
      <div>
        <h3 className="text-lg font-semibold mb-4 text-gray-900 dark:text-white">
          Completed ({jobs.filter((j) => j.status === 'completed').length})
        </h3>
        <div className="space-y-2">
          {jobs
            .filter((j) => j.status === 'completed')
            .map((job) => (
              <div
                key={job.id}
                className="p-3 bg-green-50 dark:bg-green-900/20 rounded border border-green-200 dark:border-green-700 flex justify-between items-center"
              >
                <span className="text-sm text-gray-900 dark:text-white">
                  {job.animeId} - Ep {job.episodeNumber}
                </span>
                <StatusBadge status="completed" size="sm" />
              </div>
            ))}
        </div>
      </div>

      {/* Failed Downloads */}
      {jobs.some((j) => j.status === 'failed') && (
        <div>
          <h3 className="text-lg font-semibold mb-4 text-gray-900 dark:text-white">
            Failed ({jobs.filter((j) => j.status === 'failed').length})
          </h3>
          <div className="space-y-2">
            {jobs
              .filter((j) => j.status === 'failed')
              .map((job) => (
                <div
                  key={job.id}
                  className="p-3 bg-red-50 dark:bg-red-900/20 rounded border border-red-200 dark:border-red-700"
                >
                  <div className="flex justify-between items-start">
                    <div>
                      <p className="text-sm font-medium text-gray-900 dark:text-white">
                        {job.animeId} - Ep {job.episodeNumber}
                      </p>
                      {job.errorMessage && (
                        <p className="text-xs text-red-600 dark:text-red-400 mt-1">
                          {job.errorMessage}
                        </p>
                      )}
                    </div>
                    <StatusBadge status="failed" size="sm" />
                  </div>
                </div>
              ))}
          </div>
        </div>
      )}
    </div>
  );
};
