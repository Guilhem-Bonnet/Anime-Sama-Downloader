import React, { useEffect, useMemo, useRef, useState } from 'react';
import {
  apiCancelAll,
  apiCancelJob,
  apiClearFinished,
  apiClearPending,
  apiDefaults,
  apiEnqueue,
  apiJobs,
  apiMediaTest,
  apiRetryJob,
  apiSearch,
  apiSeasonInfo,
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

  const [selectionMode, setSelectionMode] = useState<'simple' | 'checkboxes' | 'advanced'>('simple');
  const [selAll, setSelAll] = useState(true);
  const [selFrom, setSelFrom] = useState(1);
  const [selTo, setSelTo] = useState(1);
  const [advancedSel, setAdvancedSel] = useState('ALL');

  const [maxEpisodes, setMaxEpisodes] = useState<number>(0);
  const [availableEpisodes, setAvailableEpisodes] = useState<number[]>([]);
  const [selectedEpisodes, setSelectedEpisodes] = useState<number[]>([]);
  const [rangeFrom, setRangeFrom] = useState<number>(1);
  const [rangeTo, setRangeTo] = useState<number>(12);

  const [destRoot, setDestRoot] = useState('');
  const [defaultDestRoot, setDefaultDestRoot] = useState('');

  // En Docker, on n'autorise pas de chemin absolu "hôte" dans l'UI.
  // On propose uniquement un sous-dossier relatif (optionnel) sous /data/videos.
  const [dockerSubdir, setDockerSubdir] = useState('');

  const [isDocker, setIsDocker] = useState(false);
  const [allowedDestPrefixes, setAllowedDestPrefixes] = useState<string[]>([]);

  const [mediaTestBusy, setMediaTestBusy] = useState(false);
  const [mediaTestResult, setMediaTestResult] = useState<any | null>(null);

  const [jobs, setJobs] = useState<JobsSnapshot>({ pending: 0, running: 0, total: 0, jobs: [] });
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const [logLines, setLogLines] = useState<string[]>([]);
  const logRef = useRef<HTMLDivElement | null>(null);

  function encodeEpisodeSpecFromList(list: number[], max?: number, available?: number[]) {
    const uniq = Array.from(new Set(list.filter((n) => Number.isFinite(n) && n > 0))).sort((a, b) => a - b);
    if (!uniq.length) return 'ALL';

    if (max && available && available.length) {
      const availSet = new Set(available);
      const allAvail = available.every((n) => uniq.includes(n)) && uniq.every((n) => availSet.has(n));
      if (allAvail) return 'ALL';
    }

    const parts: string[] = [];
    let i = 0;
    while (i < uniq.length) {
      const start = uniq[i];
      let end = start;
      while (i + 1 < uniq.length && uniq[i + 1] === end + 1) {
        i++;
        end = uniq[i];
      }
      parts.push(start === end ? String(start) : `${start}-${end}`);
      i++;
    }
    return parts.join(',');
  }

  function buildRangeList(a: number, b: number, available: number[]) {
    const start = Math.max(1, Math.min(a, b));
    const end = Math.max(1, Math.max(a, b));
    const avail = new Set(available);
    const out: number[] = [];
    for (let n = start; n <= end; n++) {
      if (avail.has(n)) out.push(n);
    }
    return out;
  }

  const selectionText = useMemo(() => {
    if (selectionMode === 'advanced') return advancedSel.trim() || 'ALL';
    if (selectionMode === 'checkboxes') {
      return encodeEpisodeSpecFromList(selectedEpisodes, maxEpisodes, availableEpisodes);
    }
    if (selAll) return 'ALL';
    const a = Math.max(1, selFrom);
    const b = Math.max(1, selTo);
    return a === b ? String(a) : `${Math.min(a, b)}-${Math.max(a, b)}`;
  }, [selectionMode, advancedSel, selAll, selFrom, selTo, selectedEpisodes, maxEpisodes, availableEpisodes]);

  async function refreshJobs() {
    const snap = await apiJobs();
    setJobs(snap);
  }

  useEffect(() => {
    refreshJobs().catch(() => void 0);

    // Defaults (download root, etc.)
    apiDefaults()
      .then((d) => {
        setDefaultDestRoot(d.download_root || '');
        if (!destRoot.trim()) setDestRoot(d.download_root || '');
        setIsDocker(Boolean(d.is_docker));
        setAllowedDestPrefixes(Array.isArray(d.allowed_dest_prefixes) ? (d.allowed_dest_prefixes as string[]) : []);
      })
      .catch(() => void 0);

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

  const destProblem = useMemo(() => {
    if (!isDocker) return null;
    const d = destRoot.trim();
    if (!d) return null;
    // Les chemins relatifs sont résolus côté backend sous le download_root.
    if (!d.startsWith('/')) return null;
    if (!allowedDestPrefixes.length) return null;
      const ok = allowedDestPrefixes.some((p) => d === p || d.startsWith(p.replace(/\/+$/, '') + '/'));
    if (ok) return null;
    return `En Docker, la destination doit être sous: ${allowedDestPrefixes.join(', ')} (sinon ça écrit dans le FS du conteneur).`;
  }, [isDocker, destRoot, allowedDestPrefixes]);

  useEffect(() => {
    if (!baseUrl) return;
    apiSeasonInfo(baseUrl, lang, season)
      .then((info) => {
        const max = Number(info.max_episodes || 0);
        const avail = (info.available || []).map((n) => Number(n)).filter((n) => Number.isFinite(n) && n > 0);
        setMaxEpisodes(max);
        setAvailableEpisodes(avail);

        // Default selection = all available in checkbox mode.
        setSelectedEpisodes((prev) => {
          if (selectionMode !== 'checkboxes') return prev;
          return avail.slice();
        });
      })
      .catch(() => {
        setMaxEpisodes(0);
        setAvailableEpisodes([]);
      });
  }, [baseUrl, lang, season]);

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

  async function onTestMediaApi() {
    setError(null);
    setMediaTestBusy(true);
    try {
      const r = await apiMediaTest(5);
      setMediaTestResult(r);

      const jf = r?.jellyfin;
      const px = r?.plex;
      const jfTxt = jf?.configured ? (jf?.ok ? `OK (${jf?.status_code ?? '—'})` : `KO (${jf?.status_code ?? '—'})`) : 'non configuré';
      const pxTxt = px?.configured ? (px?.ok ? `OK (${px?.status_code ?? '—'})` : `KO (${px?.status_code ?? '—'})`) : 'non configuré';

      setLogLines((p) => [...p, `[INFO] Test API Jellyfin: ${jfTxt}`, `[INFO] Test API Plex: ${pxTxt}`].slice(-400));
    } catch (e: any) {
      setError(`Test API impossible: ${String(e?.message || e)}`);
    } finally {
      setMediaTestBusy(false);
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
      const dockerDest = dockerSubdir.trim();
      const r = await apiEnqueue({
        base_url: baseUrl,
        lang,
        season,
        selection: selectionText,
        dest_root: isDocker ? dockerDest : destRoot,
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
                <option value="checkboxes">Sélection par épisodes</option>
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
                selectionMode === 'checkboxes' ? (
                  <div style={{ flex: 1, minWidth: 240 }}>
                    <div className="row" style={{ justifyContent: 'space-between' }}>
                      <div className="muted small">
                        Dispo: {availableEpisodes.length ? `${availableEpisodes.length}/${maxEpisodes || '—'}` : '—'}
                      </div>
                      <div className="row">
                        <button
                          className="btn sm"
                          type="button"
                          onClick={() => setSelectedEpisodes(availableEpisodes.slice())}
                          disabled={!availableEpisodes.length}
                        >
                          Tout dispo
                        </button>
                        <button className="btn sm" type="button" onClick={() => setSelectedEpisodes([])} disabled={!availableEpisodes.length}>
                          Aucun
                        </button>
                      </div>
                    </div>

                    <div className="row" style={{ marginTop: 8 }}>
                      <span className="muted small">Range:</span>
                      <input
                        className="input"
                        style={{ maxWidth: 110, minWidth: 90 }}
                        value={rangeFrom}
                        onChange={(e) => setRangeFrom(Number(e.target.value || '1'))}
                        type="number"
                        min={1}
                      />
                      <span className="muted small">→</span>
                      <input
                        className="input"
                        style={{ maxWidth: 110, minWidth: 90 }}
                        value={rangeTo}
                        onChange={(e) => setRangeTo(Number(e.target.value || '1'))}
                        type="number"
                        min={1}
                      />
                      <button
                        className="btn sm"
                        type="button"
                        disabled={!availableEpisodes.length}
                        onClick={() => {
                          const list = buildRangeList(rangeFrom, rangeTo, availableEpisodes);
                          setSelectedEpisodes(list);
                        }}
                      >
                        Remplacer
                      </button>
                      <button
                        className="btn sm"
                        type="button"
                        disabled={!availableEpisodes.length}
                        onClick={() => {
                          const list = buildRangeList(rangeFrom, rangeTo, availableEpisodes);
                          setSelectedEpisodes((prev) => Array.from(new Set([...prev, ...list])).sort((a, b) => a - b));
                        }}
                      >
                        Ajouter
                      </button>
                      <button
                        className="btn sm"
                        type="button"
                        disabled={!availableEpisodes.length}
                        onClick={() => {
                          const list = buildRangeList(rangeFrom, rangeTo, availableEpisodes);
                          const toRemove = new Set(list);
                          setSelectedEpisodes((prev) => prev.filter((n) => !toRemove.has(n)));
                        }}
                      >
                        Retirer
                      </button>
                    </div>

                    <div className="epgrid">
                      {Array.from({ length: Math.max(0, maxEpisodes || 0) }).slice(0, 400).map((_, i) => {
                        const ep = i + 1;
                        const isAvailable = availableEpisodes.includes(ep);
                        const checked = selectedEpisodes.includes(ep);
                        return (
                          <label key={ep} className={`ep ${isAvailable ? '' : 'off'}`}>
                            <input
                              type="checkbox"
                              disabled={!isAvailable}
                              checked={checked}
                              onChange={(e) => {
                                const on = e.target.checked;
                                setSelectedEpisodes((prev) => {
                                  const set = new Set(prev);
                                  if (on) set.add(ep);
                                  else set.delete(ep);
                                  return Array.from(set).sort((a, b) => a - b);
                                });
                              }}
                            />
                            <span>EP {ep}</span>
                          </label>
                        );
                      })}
                    </div>

                    <div className="muted small" style={{ marginTop: 6 }}>
                      Envoi: <span className="badge">{selectionText}</span>
                    </div>
                  </div>
                ) : (
                  <input className="input" value={advancedSel} onChange={(e) => setAdvancedSel(e.target.value)} placeholder="ALL | 1-6 | S1E1-6 | ALLSEASONS" />
                )
              )}
            </div>
          </div>

          <div className="row" style={{ marginTop: 10 }}>
            {isDocker ? (
              <input
                className="input"
                value={dockerSubdir}
                onChange={(e) => setDockerSubdir(e.target.value)}
                placeholder="Sous-dossier (optionnel) sous /data/videos"
              />
            ) : (
              <input className="input" value={destRoot} onChange={(e) => setDestRoot(e.target.value)} placeholder="Dossier destination (optionnel)" />
            )}
            <button className="btn primary" onClick={onEnqueue} disabled={busy || !baseUrl || (!isDocker && Boolean(destProblem))}>
              Ajouter à la file
            </button>
          </div>

          {isDocker ? (
            <div className="muted small" style={{ marginTop: 6 }}>
              Docker: sortie = <span className="badge">/data/videos</span> (monté sur l’hôte). Pour changer le dossier hôte, définis <span className="badge">ASD_HOST_DOWNLOAD_ROOT</span> dans <span className="badge">.env</span>.
            </div>
          ) : null}

          {destProblem ? (
            <div style={{ marginTop: 8 }}>
              <div style={{ color: '#ffb3ad' }} className="small">{destProblem}</div>
              {defaultDestRoot ? (
                <div className="row" style={{ marginTop: 6 }}>
                  <button className="btn sm" type="button" onClick={() => setDestRoot(defaultDestRoot)}>
                    Réinitialiser au défaut
                  </button>
                </div>
              ) : null}
            </div>
          ) : null}

          <div className="row" style={{ marginTop: 10, alignItems: 'center' }}>
            <button className="btn" type="button" onClick={onTestMediaApi} disabled={busy || mediaTestBusy}>
              Tester l’API (Jellyfin / Plex)
            </button>
            <div className="muted small">
              {mediaTestResult ? (
                <>
                  Jellyfin: <span className="badge">{mediaTestResult?.jellyfin?.configured ? (mediaTestResult?.jellyfin?.ok ? 'OK' : 'KO') : 'N/A'}</span>{' '}
                  Plex: <span className="badge">{mediaTestResult?.plex?.configured ? (mediaTestResult?.plex?.ok ? 'OK' : 'KO') : 'N/A'}</span>
                </>
              ) : (
                <>Vérifie la connexion + credentials configurés (env/config.ini).</>
              )}
            </div>
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
