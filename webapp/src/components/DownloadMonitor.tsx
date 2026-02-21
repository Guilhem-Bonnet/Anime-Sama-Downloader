import React, { useEffect } from 'react';
import { useJobsStore, progressPercent } from '../stores/jobs.store';
import type { Job } from '../api';
import { useSSE } from '../hooks/useSSE';
import { StatusBadge } from './custom/StatusBadge';
import { DownloadProgress } from './custom/DownloadProgress';
import { EmptyDownloadsIllustration } from './illustrations/NocturneIllustrations';

/** Map backend job state to display status */
function displayStatus(state: Job['state']): string {
  switch (state) {
    case 'queued': return 'pending';
    case 'running':
    case 'muxing': return 'running';
    case 'completed': return 'completed';
    case 'failed':
    case 'canceled': return 'failed';
    default: return 'pending';
  }
}

function jobLabel(job: Job): string {
  const params = job.params as Record<string, any> | undefined;
  const title = params?.label || params?.animeTitle || params?.animeId || 'Téléchargement';
  const ep = params?.episodeNumber;
  return ep ? `${title} • Épisode ${ep}` : title;
}

export const DownloadMonitor: React.FC = () => {
  const { jobs, loadJobs, updateJobFromSSE, cancelJob } = useJobsStore();

  // Load jobs on mount
  useEffect(() => {
    loadJobs();
  }, []);

  // Single SSE connection for all job events
  useSSE('/api/v1/events', (data: any) => {
    if (data && data.id) {
      updateJobFromSSE(data);
    }
  });

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

  const activeJobs = jobs.filter((j) => j.state === 'running' || j.state === 'queued' || j.state === 'muxing');
  const completedJobs = jobs.filter((j) => j.state === 'completed');
  const failedJobs = jobs.filter((j) => j.state === 'failed' || j.state === 'canceled');

  return (
    <div className="space-y-8">
      {/* Active Downloads */}
      {activeJobs.length > 0 && (
        <div>
          <div style={{ display: 'flex', alignItems: 'center', gap: '12px', marginBottom: '16px' }}>
            <h3 style={{ fontSize: '18px', fontWeight: 600, margin: 0, color: 'var(--night-text-primary)' }}>
              Téléchargements Actifs ({activeJobs.length})
            </h3>
            <span style={{ fontSize: '12px', padding: '4px 12px', background: 'var(--night-accent-brown-500)', borderRadius: '12px', color: 'white', fontWeight: 500 }}>
              En cours
            </span>
          </div>
          <div className="space-y-3">
            {activeJobs.map((job) => (
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
                      <h4 style={{ fontWeight: 600, color: 'var(--night-text-primary)', margin: 0, marginBottom: '4px' }}>
                        {jobLabel(job)}
                      </h4>
                      <p style={{ fontSize: '12px', color: 'var(--night-text-secondary)', margin: 0 }}>
                        ID: {job.id.substring(0, 12)}... · {job.state}
                      </p>
                    </div>
                    <StatusBadge status={displayStatus(job.state)} size="sm" />
                  </div>
                  <DownloadProgress progress={progressPercent(job)} />
                </div>
              ))}
          </div>
        </div>
      )}

      {/* Completed Downloads */}
      {completedJobs.length > 0 && (
        <div>
          <div style={{ display: 'flex', alignItems: 'center', gap: '12px', marginBottom: '16px' }}>
            <h3 style={{ fontSize: '18px', fontWeight: 600, margin: 0, color: 'var(--night-text-primary)' }}>
              Téléchargements Terminés ({completedJobs.length})
            </h3>
            <span style={{ fontSize: '12px', padding: '4px 12px', background: '#059669', borderRadius: '12px', color: 'white', fontWeight: 500 }}>
              Complété
            </span>
          </div>
          <div className="space-y-2">
            {completedJobs.map((job) => (
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
                  <span style={{ fontSize: '14px', color: 'var(--night-text-primary)', fontWeight: 500 }}>
                    {jobLabel(job)}
                  </span>
                  <StatusBadge status="completed" size="sm" />
                </div>
              ))}
          </div>
        </div>
      )}

      {/* Failed Downloads */}
      {failedJobs.length > 0 && (
        <div>
          <div style={{ display: 'flex', alignItems: 'center', gap: '12px', marginBottom: '16px' }}>
            <h3 style={{ fontSize: '18px', fontWeight: 600, margin: 0, color: 'var(--night-text-primary)' }}>
              Erreurs ({failedJobs.length})
            </h3>
            <span style={{ fontSize: '12px', padding: '4px 12px', background: '#dc2626', borderRadius: '12px', color: 'white', fontWeight: 500 }}>
              Échoué
            </span>
          </div>
          <div className="space-y-2">
            {failedJobs.map((job) => (
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
                      <p style={{ fontSize: '14px', fontWeight: 500, color: 'var(--night-text-primary)', margin: 0, marginBottom: '4px' }}>
                        {jobLabel(job)}
                      </p>
                      {job.error && (
                        <p style={{ fontSize: '12px', color: '#dc2626', margin: 0 }}>
                          {job.error}
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
