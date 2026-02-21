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
  if (!params) return `Téléchargement #${job.id.substring(0, 8)}`;
  // Prefer the formatted label (e.g. "Naruto - S01E03")
  if (params.label) return String(params.label);
  // Fallback: build from parts
  const title = params.animeTitle || params.animeId || 'Téléchargement';
  const season = params.seasonNumber;
  const ep = params.episodeNumber;
  if (season && ep) {
    const s = String(season).padStart(2, '0');
    const e = String(ep).padStart(2, '0');
    return `${title} - S${s}E${e}`;
  }
  if (ep) return `${title} • Épisode ${ep}`;
  return title;
}

/** Translate backend error messages into user-friendly French */
function friendlyError(job: Job): string {
  const msg = job.error || '';
  if (msg.includes('missing params.url')) return 'URL de téléchargement non disponible (source anime-sama requise)';
  if (msg.includes('invalid params.url')) return 'URL de téléchargement invalide';
  if (msg.includes('missing_ffmpeg')) return 'ffmpeg non installé sur le serveur';
  if (job.errorCode === 'worker_canceled') return 'Worker arrêté pendant l\'attente';
  return msg || 'Erreur inconnue';
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
      <div style={{ textAlign: 'center', padding: '48px 16px' }}>
        <EmptyDownloadsIllustration />
        <p style={{ color: 'var(--night-text-secondary, #999)', marginTop: '16px' }}>
          Aucun téléchargement. Recherchez un anime et lancez un téléchargement.
        </p>
      </div>
    );
  }

  const activeJobs = jobs.filter((j) => j.state === 'running' || j.state === 'queued' || j.state === 'muxing');
  const completedJobs = jobs.filter((j) => j.state === 'completed');
  const failedJobs = jobs.filter((j) => j.state === 'failed' || j.state === 'canceled');

  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: '24px' }}>
      {/* Summary bar */}
      <div
        style={{
          display: 'flex',
          gap: '16px',
          padding: '12px 16px',
          background: 'var(--night-bg-secondary, #1e1e1e)',
          borderRadius: '8px',
          fontSize: '13px',
          color: 'var(--night-text-secondary, #999)',
          flexWrap: 'wrap',
        }}
      >
        <span>{jobs.length} téléchargement{jobs.length > 1 ? 's' : ''} au total</span>
        {activeJobs.length > 0 && (
          <span style={{ color: 'var(--night-accent-brown-400, #b8860b)' }}>
            ● {activeJobs.length} en cours
          </span>
        )}
        {completedJobs.length > 0 && (
          <span style={{ color: '#059669' }}>✓ {completedJobs.length} terminé{completedJobs.length > 1 ? 's' : ''}</span>
        )}
        {failedJobs.length > 0 && (
          <span style={{ color: '#dc2626' }}>✕ {failedJobs.length} échoué{failedJobs.length > 1 ? 's' : ''}</span>
        )}
      </div>
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
          <div style={{ display: 'flex', flexDirection: 'column', gap: '12px' }}>
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
                    <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
                      <button
                        onClick={() => cancelJob(job.id)}
                        style={{
                          background: 'rgba(220,38,38,0.15)',
                          border: '1px solid rgba(220,38,38,0.3)',
                          borderRadius: '6px',
                          padding: '4px 10px',
                          fontSize: '12px',
                          color: '#dc2626',
                          cursor: 'pointer',
                          fontWeight: 500,
                        }}
                        title="Annuler ce téléchargement"
                      >
                        ✕ Annuler
                      </button>
                      <StatusBadge status={displayStatus(job.state)} size="sm" />
                    </div>
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
          <div style={{ display: 'flex', flexDirection: 'column', gap: '8px' }}>
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
          <div style={{ display: 'flex', flexDirection: 'column', gap: '8px' }}>
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
                      <p style={{ fontSize: '12px', color: '#dc2626', margin: 0 }}>
                        {friendlyError(job)}
                      </p>
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
