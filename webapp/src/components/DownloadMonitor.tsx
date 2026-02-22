import React from 'react';
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
  const { jobs, isLoading, updateJobFromSSE, cancelJob } = useJobsStore();

  // Les jobs sont déjà chargés par Layout — ici on ne fait que les écouter via SSE.

  // Single SSE connection — écoute tous les événements job.* envoyés par le backend
  const JOB_EVENTS = ['job.started', 'job.progress', 'job.result', 'job.created', 'job.failed', 'job.muxing', 'job.completed', 'job.canceled'];
  useSSE('/api/v1/events', (data: any) => {
    if (data && data.id) {
      updateJobFromSSE(data);
    }
  }, JOB_EVENTS);

  // While loading and store is empty, show skeleton (not empty state)
  if (isLoading && jobs.length === 0) {
    return (
      <div style={{ display: 'flex', flexDirection: 'column', gap: '12px', padding: '16px 0' }}>
        {[1, 2, 3].map((i) => (
          <div
            key={i}
            style={{
              height: '80px',
              background: 'var(--night-bg-secondary, #1e1e1e)',
              borderRadius: '8px',
              opacity: 0.5,
              animation: 'pulse 1.5s ease-in-out infinite',
            }}
          />
        ))}
      </div>
    );
  }

  if (!isLoading && jobs.length === 0) {
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

  // Unified list — all jobs sorted newest first so recent downloads always appear at top
  const sortedJobs = [...jobs].sort(
    (a, b) => new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime()
  );

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

      {/* Unified job list — newest first */}
      <div style={{ display: 'flex', flexDirection: 'column', gap: '8px' }}>
        {sortedJobs.map((job) => {
          const isActive = job.state === 'running' || job.state === 'queued' || job.state === 'muxing';
          const isFailed = job.state === 'failed' || job.state === 'canceled';
          const isCompleted = job.state === 'completed';

          let bgColor = 'var(--night-bg-secondary, #1a1a1a)';
          if (isActive) bgColor = 'linear-gradient(135deg, rgba(143,106,61,0.08), rgba(125,114,103,0.05))';
          if (isCompleted) bgColor = 'linear-gradient(135deg, rgba(5,150,105,0.07), rgba(34,197,94,0.04))';
          if (isFailed) bgColor = 'linear-gradient(135deg, rgba(220,38,38,0.07), rgba(239,68,68,0.04))';

          return (
            <div
              key={job.id}
              className="frame-ornate"
              style={{ padding: '14px 16px', background: bgColor }}
            >
              <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', gap: '12px' }}>
                {/* Left: label + error/progress */}
                <div style={{ flex: 1, minWidth: 0 }}>
                  <div style={{ display: 'flex', alignItems: 'center', gap: '8px', marginBottom: '4px', flexWrap: 'wrap' }}>
                    <span style={{ fontWeight: 600, fontSize: '14px', color: 'var(--night-text-primary)' }}>
                      {jobLabel(job)}
                    </span>
                    <StatusBadge status={displayStatus(job.state)} size="sm" />
                  </div>
                  {isFailed && (
                    <p style={{ fontSize: '12px', color: '#dc2626', margin: 0 }}>
                      {friendlyError(job)}
                    </p>
                  )}
                  {isActive && (
                    <>
                      <DownloadProgress progress={progressPercent(job)} showLabel={false} />
                      <p style={{ fontSize: '11px', color: 'var(--night-text-secondary)', margin: '4px 0 0 0' }}>
                        {progressPercent(job)}% · {job.state}
                      </p>
                    </>
                  )}
                  {isCompleted && (
                    <p style={{ fontSize: '12px', color: '#059669', margin: 0 }}>
                      Téléchargement terminé
                    </p>
                  )}
                </div>
                {/* Right: time + cancel button for active */}
                <div style={{ display: 'flex', alignItems: 'center', gap: '8px', flexShrink: 0 }}>
                  <span style={{ fontSize: '11px', color: 'var(--night-text-secondary)' }}>
                    {new Date(job.createdAt).toLocaleTimeString('fr-FR', { hour: '2-digit', minute: '2-digit' })}
                  </span>
                  {isActive && (
                    <button
                      onClick={() => cancelJob(job.id)}
                      style={{
                        background: 'rgba(220,38,38,0.15)',
                        border: '1px solid rgba(220,38,38,0.3)',
                        borderRadius: '6px',
                        padding: '3px 8px',
                        fontSize: '11px',
                        color: '#dc2626',
                        cursor: 'pointer',
                        fontWeight: 500,
                      }}
                    >
                      ✕ Annuler
                    </button>
                  )}
                </div>
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
};
