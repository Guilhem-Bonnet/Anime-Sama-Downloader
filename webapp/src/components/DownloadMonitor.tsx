import React, { useEffect, useState } from 'react';
import { useJobsStore } from '../stores/jobs.store';
import { useSSE } from '../hooks/useSSE';
import { StatusBadge } from './custom/StatusBadge';
import { DownloadProgress } from './custom/DownloadProgress';
import { EmptyDownloadsIllustration } from './illustrations/SakuraIllustrations';

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
      <div className="text-center py-16">
        <EmptyDownloadsIllustration />
        <p className="text-gray-500 dark:text-gray-400 mt-4">
          Aucun téléchargement en cours.
        </p>
      </div>
    );
  }

  return (
    <div className="space-y-8">
      {/* Active Downloads */}
      {jobs.filter((j) => j.status === 'running' || j.status === 'pending').length > 0 && (
        <div>
          <div style={{ display: 'flex', alignItems: 'center', gap: '12px', marginBottom: '16px' }}>
            <h3 style={{ fontSize: '18px', fontWeight: 600, margin: 0, color: 'var(--sakura-text-primary)' }}>
              Téléchargements Actifs ({jobs.filter((j) => j.status === 'running').length})
            </h3>
            <span style={{ fontSize: '12px', padding: '4px 12px', background: 'var(--sakura-accent-brown-500)', borderRadius: '12px', color: 'white', fontWeight: 500 }}>
              En cours
            </span>
          </div>
          <div className="space-y-3">
            {jobs
              .filter((j) => j.status === 'running' || j.status === 'pending')
              .map((job) => (
                <div
                  key={job.id}
                  className="frame-ornate"
                  style={{
                    padding: '16px',
                    background: 'linear-gradient(135deg, rgba(143,106,61,0.05), rgba(125,114,103,0.05))',
                  }}
                >
                  <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', marginBottom: '12px' }}>
                    <div>
                      <h4 style={{ fontWeight: 600, color: 'var(--sakura-text-primary)', margin: 0, marginBottom: '4px' }}>
                        {job.animeId} • Épisode {job.episodeNumber}
                      </h4>
                      <p style={{ fontSize: '12px', color: 'var(--sakura-text-secondary)', margin: 0 }}>
                        ID: {job.id.substring(0, 12)}...
                      </p>
                    </div>
                    <StatusBadge status={job.status} size="sm" />
                  </div>
                  <DownloadProgress progress={job.progress} />
                </div>
              ))}
          </div>
        </div>
      )}

      {/* Completed Downloads */}
      {jobs.filter((j) => j.status === 'completed').length > 0 && (
        <div>
          <div style={{ display: 'flex', alignItems: 'center', gap: '12px', marginBottom: '16px' }}>
            <h3 style={{ fontSize: '18px', fontWeight: 600, margin: 0, color: 'var(--sakura-text-primary)' }}>
              Téléchargements Terminés ({jobs.filter((j) => j.status === 'completed').length})
            </h3>
            <span style={{ fontSize: '12px', padding: '4px 12px', background: '#059669', borderRadius: '12px', color: 'white', fontWeight: 500 }}>
              Complété
            </span>
          </div>
          <div className="space-y-2">
            {jobs
              .filter((j) => j.status === 'completed')
              .map((job) => (
                <div
                  key={job.id}
                  className="frame-ornate"
                  style={{
                    padding: '12px 16px',
                    background: 'linear-gradient(135deg, rgba(5,150,105,0.05), rgba(34,197,94,0.05))',
                    display: 'flex',
                    justifyContent: 'space-between',
                    alignItems: 'center',
                  }}
                >
                  <span style={{ fontSize: '14px', color: 'var(--sakura-text-primary)', fontWeight: 500 }}>
                    {job.animeId} • Ep {job.episodeNumber}
                  </span>
                  <StatusBadge status="completed" size="sm" />
                </div>
              ))}
          </div>
        </div>
      )}

      {/* Failed Downloads */}
      {jobs.some((j) => j.status === 'failed') && (
        <div>
          <div style={{ display: 'flex', alignItems: 'center', gap: '12px', marginBottom: '16px' }}>
            <h3 style={{ fontSize: '18px', fontWeight: 600, margin: 0, color: 'var(--sakura-text-primary)' }}>
              Erreurs ({jobs.filter((j) => j.status === 'failed').length})
            </h3>
            <span style={{ fontSize: '12px', padding: '4px 12px', background: '#dc2626', borderRadius: '12px', color: 'white', fontWeight: 500 }}>
              Échoué
            </span>
          </div>
          <div className="space-y-2">
            {jobs
              .filter((j) => j.status === 'failed')
              .map((job) => (
                <div
                  key={job.id}
                  className="frame-ornate"
                  style={{
                    padding: '12px 16px',
                    background: 'linear-gradient(135deg, rgba(220,38,38,0.05), rgba(239,68,68,0.05))',
                  }}
                >
                  <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start' }}>
                    <div>
                      <p style={{ fontSize: '14px', fontWeight: 500, color: 'var(--sakura-text-primary)', margin: 0, marginBottom: '4px' }}>
                        {job.animeId} • Ep {job.episodeNumber}
                      </p>
                      {job.errorMessage && (
                        <p style={{ fontSize: '12px', color: '#dc2626', margin: 0 }}>
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
