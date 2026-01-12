import React, { useEffect, useMemo, useRef, useState } from 'react';
import {
  apiCancelAll,
  apiCancelJob,
  apiClearFinished,
  apiClearPending,
  apiEnqueue,
  apiJobs,
  apiRetryJob,
  apiSearch,
  apiSeasons,
  type JobItem,
  type JobsSnapshot,
} from './api';

function statusBadge(status: JobItem['status']) {
  const cls =
    status === 'SUCCESS'
      ? 'ok'
      : status === 'RUNNING'
      ? 'run'
      : status === 'FAILED'
      ? 'fail'
      : status === 'CANCELLED'
      ? 'cancel'
      : '';
  return <span className={`badge ${cls}`}>{status}</span>;
}

function fmtBytes(n: number | null | undefined) {
  if (n == null || !Number.isFinite(n)) return '—';
  const units = ['B', 'KB', 'MB', 'GB', 'TB'];
  let v = Math.max(0, n);
  let i = 0;
  while (v >= 1024 && i < units.length - 1) {
    v /= 1024;
    i++;
  }
  return `${v.toFixed(i === 0 ? 0 : 1)} ${units[i]}`;
}

function fmtSpeed(bps: number | null | undefined) {
  if (bps == null || !Number.isFinite(bps) || bps <= 0) return '—';
  return `${fmtBytes(bps)}/s`;
}

function fmtEta(sec: number | null | undefined) {
  if (sec == null || !Number.isFinite(sec) || sec < 0) return '—';
  const s = Math.round(sec);
  const m = Math.floor(s / 60);
  const r = s % 60;
  return m > 0 ? `${m}m${String(r).padStart(2, '0')}` : `${r}s`;
}

export function App() {
  const [query, setQuery] = useState('');
  const [lang, setLang] = useState('vostfr');
  const [baseUrl, setBaseUrl] = useState<string | null>(null);
  const [seasons, setSeasons] = useState<number[]>([]);
  const [season, setSeason] = useState<number>(1);

  const [selectionMode, setSelectionMode] = useState<'simple' | 'advanced'>('simple');
  const [selAll, setSelAll] = useState(true);
  const [selFrom, setSelFrom] = useState(1);
  const [selTo, setSelTo] = useState(1);
  const [advancedSel, setAdvancedSel] = useState('ALL');

  const [destRoot, setDestRoot] = useState('');

  const [jobs, setJobs] = useState<JobsSnapshot>({ pending: 0, running: 0, total: 0, jobs: [] });
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const [logLines, setLogLines] = useState<string[]>([]);
  const logRef = useRef<HTMLDivElement | null>(null);

  const selectionText = useMemo(() => {
    if (selectionMode === 'advanced') return advancedSel.trim() || 'ALL';
    if (selAll) return 'ALL';
    const a = Math.max(1, selFrom);
    const b = Math.max(1, selTo);
    return a === b ? String(a) : `${Math.min(a, b)}-${Math.max(a, b)}`;
  }, [selectionMode, advancedSel, selAll, selFrom, selTo]);

  async function refreshJobs() {
    const snap = await apiJobs();
    setJobs(snap);
  }

  useEffect(() => {
    refreshJobs().catch(() => void 0);

    const es = new EventSource('/api/events');
    es.onmessage = (ev) => {
      try {
        const msg = JSON.parse(ev.data);
        if (msg.type === 'job') {
          refreshJobs().catch(() => void 0);
        }
        if (msg.type === 'progress') {
          const jid = String(msg.job_id || '');
          const p = msg.progress || {};
          setJobs((prev) => {
            if (!prev.jobs?.length) return prev;
            const jobs = prev.jobs.map((j) =>
              j.job_id === jid
                ? {
                    ...j,
                    progress_percent: p.percent ?? j.progress_percent,
                    progress_downloaded: p.downloaded ?? j.progress_downloaded,
                    progress_total: p.total ?? j.progress_total,
                    progress_speed_bps: p.speed_bps ?? j.progress_speed_bps,
                    progress_eta_seconds: p.eta_seconds ?? j.progress_eta_seconds,
                    progress_stage: p.stage ?? j.progress_stage,
                    progress_message: p.message ?? j.progress_message,
                  }
                : j
            );
            return { ...prev, jobs };
          });
        }
        if (msg.type === 'log') {
          setLogLines((prev) => {
            const next = [...prev, msg.msg as string];
            return next.slice(-400);
          });
        }
      } catch {
        // ignore
      }
    };
    es.onerror = () => {
      // keep trying; browser reconnects SSE
    };

    return () => {
      es.close();
    };
  }, []);

  useEffect(() => {
    // auto scroll logs
    const el = logRef.current;
    if (!el) return;
    el.scrollTop = el.scrollHeight;
  }, [logLines]);

  async function onSearch() {
    setError(null);
    setBusy(true);
    try {
      const r = await apiSearch(query.trim());
      if (!r.base_url) {
        setError('Aucun résultat.');
        return;
      }
      setBaseUrl(r.base_url);
      const s = await apiSeasons(r.base_url, lang);
      setSeasons(s.seasons || []);
      if (s.seasons?.length) setSeason(s.seasons[0]);
      setLogLines((p) => [...p, `[INFO] Trouvé: ${r.base_url}`, `[INFO] Saisons: ${(s.seasons || []).map((x) => 'S' + x).join(', ') || '—'}`].slice(-400));
    } catch (e: any) {
      setError(e?.message || String(e));
    } finally {
      setBusy(false);
    }
  }

  async function onEnqueue() {
    if (!baseUrl) {
      setError('Fais une recherche d’abord.');
      return;
    }
    setError(null);
    setBusy(true);
    try {
      const r = await apiEnqueue({
        base_url: baseUrl,
        lang,
        season,
        selection: selectionText,
        dest_root: destRoot,
      });
      if (r.error) {
        setError(r.error);
        setLogLines((p) => [...p, `[ERROR] ${r.error}`].slice(-400));
        return;
      }
      setLogLines((p) => [...p, `[INFO] Ajouté: ${r.enqueued} épisode(s)`].slice(-400));
      await refreshJobs();
    } catch (e: any) {
      setError(e?.message || String(e));
    } finally {
      setBusy(false);
    }
  }

  return (
    <div className="container">
      <div className="h1">Anime-Sama Downloader</div>
      <div className="muted small">UI Web branchée sur le core (queue + API).</div>

      <div className="grid" style={{ marginTop: 14 }}>
        <div className="card">
          <div className="row">
            <input className="input" value={query} onChange={(e) => setQuery(e.target.value)} placeholder="Nom de l’anime" />
            <select className="select" value={lang} onChange={(e) => setLang(e.target.value)}>
              <option value="vostfr">vostfr</option>
              <option value="vf">vf</option>
              <option value="vo">vo</option>
            </select>
            <button className="btn primary" onClick={onSearch} disabled={busy || !query.trim()}>
              Chercher
            </button>
          </div>

          <div className="row" style={{ marginTop: 10 }}>
            <div className="kpi">URL: <span className="muted">{baseUrl || '—'}</span></div>
          </div>

          <div className="row" style={{ marginTop: 10 }}>
            <select className="select" value={String(season)} onChange={(e) => setSeason(Number(e.target.value))}>
              {(seasons.length ? seasons : [1]).map((s) => (
                <option key={s} value={String(s)}>
                  S{s}
                </option>
              ))}
            </select>

            <div className="row" style={{ flex: 1 }}>
              <select className="select" value={selectionMode} onChange={(e) => setSelectionMode(e.target.value as any)}>
                <option value="simple">Sélection simple</option>
                <option value="advanced">Sélection avancée</option>
              </select>

              {selectionMode === 'simple' ? (
                <>
                  <label className="row muted small" style={{ gap: 6 }}>
                    <input type="checkbox" checked={selAll} onChange={(e) => setSelAll(e.target.checked)} />
                    Tous les épisodes
                  </label>
                  {!selAll ? (
                    <>
                      <input className="input" style={{ maxWidth: 130, minWidth: 100 }} value={selFrom} onChange={(e) => setSelFrom(Number(e.target.value || '1'))} type="number" min={1} />
                      <span className="muted">→</span>
                      <input className="input" style={{ maxWidth: 130, minWidth: 100 }} value={selTo} onChange={(e) => setSelTo(Number(e.target.value || '1'))} type="number" min={1} />
                    </>
                  ) : null}
                </>
              ) : (
                <input className="input" value={advancedSel} onChange={(e) => setAdvancedSel(e.target.value)} placeholder="ALL | 1-6 | S1E1-6 | ALLSEASONS" />
              )}
            </div>
          </div>

          <div className="row" style={{ marginTop: 10 }}>
            <input className="input" value={destRoot} onChange={(e) => setDestRoot(e.target.value)} placeholder="Dossier destination (optionnel)" />
            <button className="btn primary" onClick={onEnqueue} disabled={busy || !baseUrl}>
              Ajouter à la file
            </button>
          </div>

          <div className="row" style={{ marginTop: 10 }}>
            <button
              className="btn danger"
              onClick={() => {
                apiCancelAll().then(() => setLogLines((p) => [...p, '[INFO] Annuler tout demandé'].slice(-400)));
              }}
              disabled={busy}
            >
              Annuler tout
            </button>
            <button
              className="btn"
              onClick={() => {
                apiClearPending().then((r) => setLogLines((p) => [...p, `[INFO] File vidée: ${r.cleared}`].slice(-400)));
              }}
              disabled={busy}
            >
              Vider file
            </button>
            <button className="btn" onClick={refreshJobs} disabled={busy}>
              Rafraîchir
            </button>
            <button
              className="btn"
              onClick={() => {
                apiClearFinished().then((r) => setLogLines((p) => [...p, `[INFO] Nettoyé: ${r.cleared} terminé(s)`].slice(-400)));
                refreshJobs().catch(() => void 0);
              }}
              disabled={busy}
            >
              Nettoyer terminés
            </button>
            <div className="kpi">Sélection: <span className="muted">{selectionText}</span></div>
          </div>

          {error ? <div style={{ marginTop: 10, color: '#ffb3ad' }}>{error}</div> : null}
        </div>

        <div className="card">
          <div className="row" style={{ justifyContent: 'space-between' }}>
            <div>
              <div className="h1" style={{ margin: 0, fontSize: 16 }}>Téléchargements</div>
              <div className="kpi">File: {jobs.pending} en attente | {jobs.running} en cours | total {jobs.total}</div>
            </div>
          </div>

          <div style={{ marginTop: 10, overflow: 'auto' }}>
            <table className="table">
              <thead>
                <tr>
                  <th>Épisode</th>
                  <th>Progress</th>
                  <th>Statut</th>
                  <th>Action</th>
                </tr>
              </thead>
              <tbody>
                {jobs.jobs.map((j) => (
                  <tr key={j.job_id}>
                    <td>
                      <div>{j.label}</div>
                      {j.result_path ? <div className="muted small">{j.result_path}</div> : null}
                      {j.error ? <div style={{ color: '#ffb3ad' }} className="small">{j.error}</div> : null}
                    </td>
                    <td>
                      <div className="small muted">{j.progress_stage || '—'}</div>
                      <div className="small">
                        {j.progress_percent != null ? `${Math.min(100, Math.max(0, j.progress_percent)).toFixed(1)}%` : '—'}
                        {' · '}
                        {fmtBytes(j.progress_downloaded)} / {fmtBytes(j.progress_total)}
                      </div>
                      <div className="small muted">{fmtSpeed(j.progress_speed_bps)} · ETA {fmtEta(j.progress_eta_seconds)}</div>
                    </td>
                    <td>{statusBadge(j.status)}</td>
                    <td>
                      <button
                        className="btn"
                        disabled={j.status !== 'RUNNING' && j.status !== 'PENDING'}
                        onClick={() => apiCancelJob(j.job_id).then(() => refreshJobs())}
                      >
                        Annuler
                      </button>
                      <button
                        className="btn"
                        style={{ marginLeft: 8 }}
                        disabled={j.status !== 'FAILED' && j.status !== 'CANCELLED'}
                        onClick={() => apiRetryJob(j.job_id).then(() => refreshJobs())}
                      >
                        Retry
                      </button>
                    </td>
                  </tr>
                ))}
                {!jobs.jobs.length ? (
                  <tr>
                    <td colSpan={4} className="muted">Aucun job.</td>
                  </tr>
                ) : null}
              </tbody>
            </table>
          </div>
        </div>
      </div>

      <div className="card" style={{ marginTop: 16 }}>
        <div className="row" style={{ justifyContent: 'space-between' }}>
          <div>
            <div className="h1" style={{ margin: 0, fontSize: 16 }}>Logs</div>
            <div className="muted small">Flux temps-réel (SSE). Idéal pour debug.</div>
          </div>
          <button className="btn" onClick={() => setLogLines([])}>Vider</button>
        </div>
        <div ref={logRef} className="log" style={{ marginTop: 10 }}>
          {logLines.map((l, i) => (
            <div key={i} className="logline">{l}</div>
          ))}
        </div>
      </div>
    </div>
  );
}
