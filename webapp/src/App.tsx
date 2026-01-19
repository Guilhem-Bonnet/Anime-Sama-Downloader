import React, { useEffect, useMemo, useState } from 'react';
import {
  apiAniListAiring,
  apiAniListImportAuto,
  apiAniListViewer,
  apiAnimeSamaResolve,
  apiCancelJob,
  apiCreateSubscription,
  apiDeleteSubscription,
  apiGetSettings,
  apiListJobs,
  apiListSubscriptions,
  apiPutSettings,
  apiSyncAll,
  apiSyncSubscription,
  type AniListAiringScheduleEntry,
  type AniListImportAutoRequest,
  type AniListImportAutoResponse,
  type AniListViewer,
  type AnimeSamaResolvedCandidate,
  type Job,
  type Settings,
  type Subscription,
} from './api';

type Tab = 'abonnements' | 'calendrier' | 'anilist' | 'jobs' | 'settings';

type SubSort = 'label' | 'nextCheckAt';
type CalendarMode = 'checks' | 'anilist';

function badgeForState(state: Job['state']) {
  const cls =
    state === 'completed'
      ? 'ok'
      : state === 'running' || state === 'muxing'
      ? 'run'
      : state === 'failed'
      ? 'fail'
      : state === 'canceled'
      ? 'cancel'
      : '';
  return <span className={`badge ${cls}`}>{state}</span>;
}

function fmtWhen(iso: string | null | undefined) {
  if (!iso) return '—';
  const d = new Date(iso);
  if (!Number.isFinite(d.getTime())) return '—';
  return d.toLocaleString();
}

function fmtUnix(ts: number | null | undefined) {
  if (!ts || !Number.isFinite(ts)) return '—';
  return new Date(ts * 1000).toLocaleString();
}

function fmtRelative(iso: string | null | undefined) {
  if (!iso) return '—';
  const d = new Date(iso);
  if (!Number.isFinite(d.getTime())) return '—';
  const ms = d.getTime() - Date.now();
  const abs = Math.abs(ms);
  const sign = ms < 0 ? '-' : '+';
  const mins = Math.round(abs / 60000);
  if (mins < 60) return `${sign}${mins} min`;
  const hrs = Math.round(mins / 60);
  if (hrs < 48) return `${sign}${hrs} h`;
  const days = Math.round(hrs / 24);
  return `${sign}${days} j`;
}

function isDue(iso: string | null | undefined) {
  if (!iso) return false;
  const d = new Date(iso);
  if (!Number.isFinite(d.getTime())) return false;
  return d.getTime() <= Date.now();
}

function pct(x: number | null | undefined) {
  if (!Number.isFinite(x as any)) return 0;
  return Math.max(0, Math.min(100, Math.round((x as number) * 100)));
}

function pretty(v: any) {
  if (v === null || v === undefined) return '';
  try {
    return JSON.stringify(v, null, 2);
  } catch {
    return String(v);
  }
}

function titleOf(entry: AniListAiringScheduleEntry) {
  return entry.media.title.romaji || entry.media.title.english || entry.media.title.native || `#${entry.media.id}`;
}

function titlesOf(entry: AniListAiringScheduleEntry) {
  const t = entry.media.title;
  const vals = [t.romaji, t.english, t.native]
    .filter((s) => typeof s === 'string')
    .map((s) => String(s).trim())
    .filter(Boolean);
  const seen = new Set<string>();
  const out: string[] = [];
  for (const s of vals) {
    const k = s.toLowerCase();
    if (seen.has(k)) continue;
    seen.add(k);
    out.push(s);
  }
  return out;
}

function jobDownloadInfo(j: Job) {
  const p = (j.params || {}) as any;
  const r = (j.result || {}) as any;
  const url = (r.resolvedUrl || r.url || p.resolvedUrl || p.url) as string | undefined;
  const sourceUrl = (r.sourceUrl || p.sourceUrl) as string | undefined;
  const path = (r.path || p.path) as string | undefined;
  const mode = (r.mode || p.mode) as string | undefined;
  const bytes = (r.bytes ?? r.sizeBytes) as number | undefined;
  return {
    url: typeof url === 'string' ? url : undefined,
    sourceUrl: typeof sourceUrl === 'string' ? sourceUrl : undefined,
    path: typeof path === 'string' ? path : undefined,
    mode: typeof mode === 'string' ? mode : undefined,
    bytes: typeof bytes === 'number' ? bytes : undefined,
  };
}

function fmtBytes(n: number | null | undefined) {
  if (!Number.isFinite(n as any) || (n as number) < 0) return '—';
  const v = n as number;
  if (v < 1024) return `${v} B`;
  const kb = v / 1024;
  if (kb < 1024) return `${kb.toFixed(1)} KB`;
  const mb = kb / 1024;
  if (mb < 1024) return `${mb.toFixed(1)} MB`;
  const gb = mb / 1024;
  return `${gb.toFixed(2)} GB`;
}

export function App() {
  const [tab, setTab] = useState<Tab>('abonnements');

  const [busy, setBusy] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [events, setEvents] = useState<string[]>([]);

  const [settings, setSettings] = useState<Settings | null>(null);
  const [settingsDraft, setSettingsDraft] = useState<Settings | null>(null);

  const [subs, setSubs] = useState<Subscription[]>([]);
  const [jobs, setJobs] = useState<Job[]>([]);

  const [subsFilter, setSubsFilter] = useState('');
  const [subsSort, setSubsSort] = useState<SubSort>('nextCheckAt');

  const [calendarMode, setCalendarMode] = useState<CalendarMode>('checks');

  const [jobDetailsID, setJobDetailsID] = useState<string | null>(null);
  const [showEvents, setShowEvents] = useState(true);

  const [newBaseUrl, setNewBaseUrl] = useState('');
  const [newLabel, setNewLabel] = useState('');
  const [newPlayer, setNewPlayer] = useState('auto');

  const [viewer, setViewer] = useState<AniListViewer | null>(null);
  const [airingDays, setAiringDays] = useState(7);
  const [airing, setAiring] = useState<AniListAiringScheduleEntry[]>([]);

  const [resolveSeason, setResolveSeason] = useState(1);
  const [resolveLang, setResolveLang] = useState('vostfr');
  const [resolveMaxCandidates, setResolveMaxCandidates] = useState(5);
  const [resolved, setResolved] = useState<{
    entryId: number;
    title: string;
    candidates: AnimeSamaResolvedCandidate[];
  } | null>(null);

  const [importReq, setImportReq] = useState<AniListImportAutoRequest>({
    statuses: ['CURRENT', 'PLANNING'],
    season: 1,
    lang: 'vostfr',
    maxCandidates: 5,
    minScore: 0.65,
  });
  const [importRes, setImportRes] = useState<AniListImportAutoResponse | null>(null);

  async function refreshSettings() {
    const s = await apiGetSettings();
    setSettings(s);
    setSettingsDraft(s);
  }

  async function refreshSubs() {
    setSubs(await apiListSubscriptions(500));
  }

  async function refreshJobs() {
    setJobs(await apiListJobs(300));
  }

  async function refreshViewer() {
    setViewer(await apiAniListViewer());
  }

  async function refreshAiring() {
    setAiring(await apiAniListAiring(airingDays, 60));
  }

  useEffect(() => {
    refreshSettings().catch(() => void 0);
    refreshSubs().catch(() => void 0);
    refreshJobs().catch(() => void 0);

    const es = new EventSource('/api/v1/events');
    const topics = [
      'hello',
      'ping',
      'job.created',
      'job.started',
      'job.progress',
      'job.result',
      'job.completed',
      'job.failed',
      'job.canceled',
      'subscription.created',
      'subscription.updated',
      'subscription.synced',
    ];

    for (const t of topics) {
      es.addEventListener(t, (ev) => {
        try {
          const data = (ev as MessageEvent).data;
          setEvents((prev) => [...prev, `${new Date().toLocaleTimeString()} ${t} ${String(data).slice(0, 220)}`].slice(-200));
        } catch {
          // ignore
        }
        if (t.startsWith('job.')) refreshJobs().catch(() => void 0);
        if (t.startsWith('subscription.')) refreshSubs().catch(() => void 0);
      });
    }

    return () => {
      es.close();
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  useEffect(() => {
    if (tab === 'calendrier') refreshAiring().catch(() => void 0);
    if (tab === 'anilist') refreshViewer().catch(() => void 0);
  }, [tab]);

  useEffect(() => {
    if (tab === 'calendrier') refreshSubs().catch(() => void 0);
  }, [tab]);

  const subsSorted = useMemo(() => {
    const q = subsFilter.trim().toLowerCase();
    const filtered = !q
      ? [...subs]
      : subs.filter((s) => (s.label + ' ' + s.baseUrl).toLowerCase().includes(q));

    const byNext = (a: Subscription, b: Subscription) => {
      const da = new Date(a.nextCheckAt).getTime();
      const db = new Date(b.nextCheckAt).getTime();
      if (Number.isFinite(da) && Number.isFinite(db)) return da - db;
      if (Number.isFinite(da)) return -1;
      if (Number.isFinite(db)) return 1;
      return a.label.localeCompare(b.label);
    };

    const byLabel = (a: Subscription, b: Subscription) => a.label.localeCompare(b.label);

    return filtered.sort((a, b) => {
      const da = isDue(a.nextCheckAt);
      const db = isDue(b.nextCheckAt);
      if (da !== db) return da ? -1 : 1;
      return subsSort === 'nextCheckAt' ? byNext(a, b) : byLabel(a, b);
    });
  }, [subs, subsFilter, subsSort]);

  const airingByDay = useMemo(() => {
    const m = new Map<string, AniListAiringScheduleEntry[]>();
    for (const e of airing) {
      const d = new Date(e.airingAt * 1000);
      const key = `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`;
      if (!m.has(key)) m.set(key, []);
      m.get(key)!.push(e);
    }
    return Array.from(m.entries()).sort((a, b) => a[0].localeCompare(b[0]));
  }, [airing]);

  const subsByCheck = useMemo(() => {
    const now = Date.now();
    const due: Subscription[] = [];
    const soon: Subscription[] = [];
    const later: Subscription[] = [];
    for (const s of subs) {
      const t = new Date(s.nextCheckAt).getTime();
      if (!Number.isFinite(t)) {
        later.push(s);
        continue;
      }
      if (t <= now) {
        due.push(s);
        continue;
      }
      if (t <= now + 24 * 3600 * 1000) {
        soon.push(s);
        continue;
      }
      later.push(s);
    }
    const byTime = (a: Subscription, b: Subscription) => new Date(a.nextCheckAt).getTime() - new Date(b.nextCheckAt).getTime();
    due.sort(byTime);
    soon.sort(byTime);
    later.sort(byTime);
    return { due, soon, later };
  }, [subs]);

  async function wrap<T>(fn: () => Promise<T>) {
    setError(null);
    setBusy(true);
    try {
      return await fn();
    } catch (e: any) {
      setError(e?.message || String(e));
      throw e;
    } finally {
      setBusy(false);
    }
  }

  return (
    <div className="container">
      <div className="row" style={{ justifyContent: 'space-between' }}>
        <div>
          <div className="h1">Anime-Sama Downloader — UI</div>
          <div className="muted small">Backend Go: abonnements, AniList, jobs, calendrier</div>
        </div>
        <div className="row">
          {(['abonnements', 'calendrier', 'anilist', 'jobs', 'settings'] as Tab[]).map((t) => (
            <button key={t} className={`btn sm ${tab === t ? 'primary' : ''}`} onClick={() => setTab(t)}>
              {t}
            </button>
          ))}
        </div>
      </div>

      {error && (
        <div className="card" style={{ borderColor: 'rgba(231,76,60,.55)', marginTop: 12 }}>
          <div className="row" style={{ justifyContent: 'space-between' }}>
            <div>
              <div style={{ color: 'var(--err)', fontWeight: 600 }}>Erreur</div>
              <div className="small muted" style={{ whiteSpace: 'pre-wrap' }}>
                {error}
              </div>
            </div>
            <button className="btn sm" onClick={() => setError(null)}>
              Fermer
            </button>
          </div>
        </div>
      )}

      <div className="grid" style={{ marginTop: 12 }}>
        <div className="card">
          {tab === 'abonnements' && (
            <>
              <div className="row" style={{ justifyContent: 'space-between' }}>
                <div>
                  <div style={{ fontWeight: 600 }}>Abonnements</div>
                  <div className="muted small">BaseURL = saison + langue (ex: /catalogue/.../saison1/vostfr/)</div>
                </div>
                <div className="row">
                  <button className="btn sm" disabled={busy} onClick={() => wrap(refreshSubs)}>
                    Rafraîchir
                  </button>
                  <button
                    className="btn sm primary"
                    disabled={busy}
                    onClick={() =>
                      wrap(async () => {
                        await apiSyncAll({ enqueue: true, dueOnly: true, limit: 500 });
                        await refreshSubs();
                        await refreshJobs();
                      })
                    }
                  >
                    Sync due
                  </button>
                </div>
              </div>

              <div className="row" style={{ marginTop: 10 }}>
                <input className="input" placeholder="Base URL (Anime-Sama)" value={newBaseUrl} onChange={(e) => setNewBaseUrl(e.target.value)} />
                <input className="input" placeholder="Label" value={newLabel} onChange={(e) => setNewLabel(e.target.value)} />
                <select className="select" value={newPlayer} onChange={(e) => setNewPlayer(e.target.value)}>
                  <option value="auto">auto</option>
                  <option value="Player 1">Player 1</option>
                  <option value="Player 2">Player 2</option>
                </select>
                <button
                  className="btn primary"
                  disabled={busy}
                  onClick={() =>
                    wrap(async () => {
                      await apiCreateSubscription({ baseUrl: newBaseUrl, label: newLabel, player: newPlayer });
                      setNewBaseUrl('');
                      setNewLabel('');
                      setNewPlayer('auto');
                      await refreshSubs();
                    })
                  }
                >
                  Ajouter
                </button>
              </div>

              <div className="row" style={{ marginTop: 10, justifyContent: 'space-between' }}>
                <input
                  className="input"
                  placeholder="Filtrer (label ou URL)"
                  value={subsFilter}
                  onChange={(e) => setSubsFilter(e.target.value)}
                />
                <select className="select" value={subsSort} onChange={(e) => setSubsSort(e.target.value as SubSort)}>
                  <option value="nextCheckAt">tri: prochain check</option>
                  <option value="label">tri: label</option>
                </select>
              </div>

              <table className="table" style={{ marginTop: 10 }}>
                <thead>
                  <tr>
                    <th>Label</th>
                    <th>Épisodes</th>
                    <th>Prochain check</th>
                    <th></th>
                  </tr>
                </thead>
                <tbody>
                  {subsSorted.map((s) => (
                    <tr key={s.id}>
                      <td>
                        <div style={{ fontWeight: 600 }}>{s.label}</div>
                        <div className="muted small" style={{ wordBreak: 'break-all' }}>
                          {s.baseUrl}
                        </div>
                        <div className="row" style={{ gap: 8 }}>
                          <div className="muted small">player: {s.player}</div>
                          {isDue(s.nextCheckAt) && <span className="badge warn">due</span>}
                        </div>
                      </td>
                      <td>
                        <div className="small">dispo: {s.lastAvailableEpisode || 0}</div>
                        <div className="small">téléchargé: {s.lastDownloadedEpisode || 0}</div>
                        <div className="bar" title="téléchargé / disponible">
                          <div
                            className="barfill"
                            style={{
                              width: `${
                                s.lastAvailableEpisode > 0
                                  ? Math.min(100, Math.round((100 * (s.lastDownloadedEpisode || 0)) / s.lastAvailableEpisode))
                                  : 0
                              }%`,
                            }}
                          />
                        </div>
                        <div className="muted small">check: {fmtWhen(s.lastCheckedAt)}</div>
                      </td>
                      <td className="small">{fmtWhen(s.nextCheckAt)}</td>
                      <td>
                        <div className="row" style={{ justifyContent: 'flex-end' }}>
                          <a className="btn sm" href={s.baseUrl} target="_blank" rel="noreferrer">
                            Ouvrir
                          </a>
                          <button
                            className="btn sm"
                            disabled={busy}
                            onClick={() =>
                              wrap(async () => {
                                await apiSyncSubscription(s.id, true);
                                await refreshSubs();
                                await refreshJobs();
                              })
                            }
                          >
                            Sync + enqueue
                          </button>
                          <button
                            className="btn sm danger"
                            disabled={busy}
                            onClick={() =>
                              wrap(async () => {
                                await apiDeleteSubscription(s.id);
                                await refreshSubs();
                              })
                            }
                          >
                            Suppr
                          </button>
                        </div>
                      </td>
                    </tr>
                  ))}
                  {!subsSorted.length && (
                    <tr>
                      <td colSpan={4} className="muted">
                        Aucun abonnement.
                      </td>
                    </tr>
                  )}
                </tbody>
              </table>
            </>
          )}

          {tab === 'calendrier' && (
            <>
              <div className="row" style={{ justifyContent: 'space-between' }}>
                <div>
                  <div style={{ fontWeight: 600 }}>Calendrier</div>
                  <div className="muted small">Checks abonnements + schedule AniList</div>
                </div>
                <div className="row">
                  <button
                    className={`btn sm ${calendarMode === 'checks' ? 'primary' : ''}`}
                    onClick={() => setCalendarMode('checks')}
                  >
                    Abonnements
                  </button>
                  <button
                    className={`btn sm ${calendarMode === 'anilist' ? 'primary' : ''}`}
                    onClick={() => setCalendarMode('anilist')}
                  >
                    AniList
                  </button>
                </div>
              </div>

              {calendarMode === 'checks' && (
                <>
                  <div className="row" style={{ marginTop: 10, justifyContent: 'space-between' }}>
                    <div className="row">
                      <span className="badge warn">due {subsByCheck.due.length}</span>
                      <span className="badge"><span className="muted">next:</span> {subsByCheck.soon[0] ? fmtRelative(subsByCheck.soon[0].nextCheckAt) : '—'}</span>
                    </div>
                    <div className="row">
                      <button className="btn sm" disabled={busy} onClick={() => wrap(refreshSubs)}>
                        Rafraîchir
                      </button>
                      <button
                        className="btn sm primary"
                        disabled={busy}
                        onClick={() =>
                          wrap(async () => {
                            await apiSyncAll({ enqueue: true, dueOnly: true, limit: 500 });
                            await refreshSubs();
                            await refreshJobs();
                          })
                        }
                      >
                        Sync due
                      </button>
                    </div>
                  </div>

                  {subsByCheck.due.length > 0 && (
                    <div style={{ marginTop: 12 }}>
                      <div style={{ fontWeight: 600, marginBottom: 6 }}>En retard</div>
                      <table className="table">
                        <thead>
                          <tr>
                            <th>Label</th>
                            <th>Quand</th>
                            <th></th>
                          </tr>
                        </thead>
                        <tbody>
                          {subsByCheck.due.slice(0, 40).map((s) => (
                            <tr key={s.id}>
                              <td>
                                <div style={{ fontWeight: 600 }}>{s.label}</div>
                                <div className="muted small" style={{ wordBreak: 'break-all' }}>{s.baseUrl}</div>
                              </td>
                              <td className="small">
                                {fmtWhen(s.nextCheckAt)}
                                <div className="muted small">{fmtRelative(s.nextCheckAt)}</div>
                              </td>
                              <td>
                                <div className="row" style={{ justifyContent: 'flex-end' }}>
                                  <a className="btn sm" href={s.baseUrl} target="_blank" rel="noreferrer">Ouvrir</a>
                                  <button
                                    className="btn sm"
                                    disabled={busy}
                                    onClick={() =>
                                      wrap(async () => {
                                        await apiSyncSubscription(s.id, true);
                                        await refreshSubs();
                                        await refreshJobs();
                                      })
                                    }
                                  >
                                    Sync
                                  </button>
                                </div>
                              </td>
                            </tr>
                          ))}
                        </tbody>
                      </table>
                    </div>
                  )}

                  <div style={{ marginTop: 12 }}>
                    <div style={{ fontWeight: 600, marginBottom: 6 }}>Prochaines 24h</div>
                    <table className="table">
                      <thead>
                        <tr>
                          <th>Label</th>
                          <th>Quand</th>
                          <th>Épisodes</th>
                        </tr>
                      </thead>
                      <tbody>
                        {subsByCheck.soon.slice(0, 60).map((s) => (
                          <tr key={s.id}>
                            <td>
                              <div style={{ fontWeight: 600 }}>{s.label}</div>
                              <div className="muted small">{fmtRelative(s.nextCheckAt)}</div>
                            </td>
                            <td className="small">{fmtWhen(s.nextCheckAt)}</td>
                            <td className="small">
                              dispo: {s.lastAvailableEpisode || 0} / dl: {s.lastDownloadedEpisode || 0}
                            </td>
                          </tr>
                        ))}
                        {!subsByCheck.soon.length && (
                          <tr>
                            <td colSpan={3} className="muted">Aucun check planifié (24h).</td>
                          </tr>
                        )}
                      </tbody>
                    </table>
                  </div>
                </>
              )}

              {calendarMode === 'anilist' && (
                <>
                  <div className="row" style={{ marginTop: 10, justifyContent: 'space-between' }}>
                    <div className="row">
                      <select className="select" value={String(airingDays)} onChange={(e) => setAiringDays(Number(e.target.value) || 7)}>
                        <option value="3">3 jours</option>
                        <option value="7">7 jours</option>
                        <option value="14">14 jours</option>
                      </select>
                      <button className="btn sm" disabled={busy} onClick={() => wrap(refreshAiring)}>
                        Rafraîchir
                      </button>

                      <label className="small muted" style={{ marginLeft: 10 }}>lang</label>
                      <select className="select" value={resolveLang} onChange={(e) => setResolveLang(e.target.value)}>
                        <option value="vostfr">vostfr</option>
                        <option value="vf">vf</option>
                      </select>

                      <label className="small muted">saison</label>
                      <input className="select" type="number" min={1} value={resolveSeason} onChange={(e) => setResolveSeason(Number(e.target.value) || 1)} />

                      <label className="small muted">max</label>
                      <input className="select" type="number" min={1} max={10} value={resolveMaxCandidates} onChange={(e) => setResolveMaxCandidates(Number(e.target.value) || 5)} />
                    </div>
                    <button className="btn sm" onClick={() => setTab('anilist')}>
                      Import watchlist
                    </button>
                  </div>

                  {airingByDay.map(([day, list]) => (
                    <div key={day} style={{ marginTop: 12 }}>
                      <div style={{ fontWeight: 600, marginBottom: 6 }}>{day}</div>
                      <table className="table">
                        <thead>
                          <tr>
                            <th>Heure</th>
                            <th>Titre</th>
                            <th>Ep</th>
                            <th></th>
                          </tr>
                        </thead>
                        <tbody>
                          {list
                            .slice()
                            .sort((a, b) => a.airingAt - b.airingAt)
                            .map((e) => (
                              <React.Fragment key={e.id}>
                                <tr>
                                  <td className="small">{fmtUnix(e.airingAt)}</td>
                                  <td>{titleOf(e)}</td>
                                  <td className="small">{e.episode}</td>
                                  <td>
                                    <div className="row" style={{ justifyContent: 'flex-end' }}>
                                      <button
                                        className="btn sm"
                                        disabled={busy}
                                        onClick={() =>
                                          wrap(async () => {
                                            const titles = titlesOf(e);
                                            const res = await apiAnimeSamaResolve({
                                              titles,
                                              season: resolveSeason,
                                              lang: resolveLang,
                                              maxCandidates: resolveMaxCandidates,
                                            });
                                            setResolved({ entryId: e.id, title: titleOf(e), candidates: res.candidates || [] });
                                          })
                                        }
                                      >
                                        Résoudre
                                      </button>
                                      {resolved?.entryId === e.id && (
                                        <button className="btn sm" disabled={busy} onClick={() => setResolved(null)}>
                                          Fermer
                                        </button>
                                      )}
                                    </div>
                                  </td>
                                </tr>

                                {resolved?.entryId === e.id && (
                                  <tr>
                                    <td colSpan={4}>
                                      {!resolved.candidates.length ? (
                                        <div className="muted small">Aucun candidat trouvé.</div>
                                      ) : (
                                        <table className="table" style={{ marginTop: 6 }}>
                                          <thead>
                                            <tr>
                                              <th>Score</th>
                                              <th>Match</th>
                                              <th>Base URL</th>
                                              <th></th>
                                            </tr>
                                          </thead>
                                          <tbody>
                                            {resolved.candidates.map((c) => (
                                              <tr key={c.baseUrl + '|' + c.matchedTitle}>
                                                <td className="small">{Math.round((c.score || 0) * 100)}%</td>
                                                <td className="small">{c.matchedTitle || c.slug}</td>
                                                <td className="small" style={{ wordBreak: 'break-all' }}>{c.baseUrl}</td>
                                                <td>
                                                  <div className="row" style={{ justifyContent: 'flex-end' }}>
                                                    <a className="btn sm" href={c.catalogueUrl} target="_blank" rel="noreferrer">
                                                      Catalogue
                                                    </a>
                                                    <a className="btn sm" href={c.baseUrl} target="_blank" rel="noreferrer">
                                                      Saison
                                                    </a>
                                                    <button
                                                      className="btn sm primary"
                                                      disabled={busy}
                                                      onClick={() =>
                                                        wrap(async () => {
                                                          const label = `${resolved.title} (S${resolveSeason} ${resolveLang})`;
                                                          await apiCreateSubscription({ baseUrl: c.baseUrl, label, player: newPlayer });
                                                          setResolved(null);
                                                          await refreshSubs();
                                                          setTab('abonnements');
                                                        })
                                                      }
                                                    >
                                                      Créer
                                                    </button>
                                                  </div>
                                                </td>
                                              </tr>
                                            ))}
                                          </tbody>
                                        </table>
                                      )}
                                    </td>
                                  </tr>
                                )}
                              </React.Fragment>
                            ))}
                        </tbody>
                      </table>
                    </div>
                  ))}
                  {!airingByDay.length && <div className="muted" style={{ marginTop: 12 }}>Aucune donnée.</div>}
                </>
              )}
            </>
          )}

          {tab === 'anilist' && (
            <>
              <div className="row" style={{ justifyContent: 'space-between' }}>
                <div>
                  <div style={{ fontWeight: 600 }}>AniList</div>
                  <div className="muted small">Connexion via settings.anilistToken</div>
                </div>
                <button className="btn sm" disabled={busy} onClick={() => wrap(refreshViewer)}>
                  Viewer
                </button>
              </div>

              <div className="card" style={{ marginTop: 10 }}>
                <div className="row" style={{ justifyContent: 'space-between' }}>
                  <div>
                    <div style={{ fontWeight: 600 }}>Import auto watchlist → abonnements</div>
                    <div className="muted small">Crée des subscriptions (saison/langue) depuis ta watchlist.</div>
                  </div>
                  <button
                    className="btn sm primary"
                    disabled={busy}
                    onClick={() =>
                      wrap(async () => {
                        const res = await apiAniListImportAuto(importReq);
                        setImportRes(res);
                        await refreshSubs();
                      })
                    }
                  >
                    Lancer
                  </button>
                </div>

                <div className="row" style={{ marginTop: 10 }}>
                  <label className="small muted">statuses</label>
                  {['CURRENT', 'PLANNING', 'PAUSED', 'COMPLETED'].map((st) => (
                    <label key={st} className="small" style={{ display: 'flex', gap: 6, alignItems: 'center' }}>
                      <input
                        type="checkbox"
                        checked={importReq.statuses.includes(st)}
                        onChange={(e) => {
                          setImportReq((prev) => {
                            const next = new Set(prev.statuses);
                            if (e.target.checked) next.add(st);
                            else next.delete(st);
                            return { ...prev, statuses: Array.from(next) };
                          });
                        }}
                      />
                      {st}
                    </label>
                  ))}
                </div>

                <div className="row" style={{ marginTop: 10 }}>
                  <label className="small muted">lang</label>
                  <select className="select" value={importReq.lang} onChange={(e) => setImportReq((p) => ({ ...p, lang: e.target.value }))}>
                    <option value="vostfr">vostfr</option>
                    <option value="vf">vf</option>
                  </select>
                  <label className="small muted">saison</label>
                  <input className="select" type="number" min={1} value={importReq.season} onChange={(e) => setImportReq((p) => ({ ...p, season: Number(e.target.value) || 1 }))} />
                  <label className="small muted">maxCandidates</label>
                  <input className="select" type="number" min={1} value={importReq.maxCandidates} onChange={(e) => setImportReq((p) => ({ ...p, maxCandidates: Number(e.target.value) || 5 }))} />
                  <label className="small muted">minScore</label>
                  <input
                    className="select"
                    type="number"
                    step={0.05}
                    min={0}
                    max={1}
                    value={importReq.minScore}
                    onChange={(e) => setImportReq((p) => ({ ...p, minScore: Number(e.target.value) || 0 }))}
                  />
                </div>

                {viewer && (
                  <div className="muted small" style={{ marginTop: 10 }}>
                    Viewer: {viewer.name} (#{viewer.id})
                  </div>
                )}
              </div>

              {importRes && (
                <div style={{ marginTop: 12 }}>
                  <div className="row" style={{ justifyContent: 'space-between' }}>
                    <div style={{ fontWeight: 600 }}>Résultat</div>
                    <button className="btn sm" onClick={() => setImportRes(null)}>
                      Clear
                    </button>
                  </div>
                  <div className="row" style={{ marginTop: 8 }}>
                    <span className="badge ok">created {importRes.created?.length || 0}</span>
                    <span className="badge">skipped {importRes.skipped?.length || 0}</span>
                    <span className="badge fail">errors {importRes.errors?.length || 0}</span>
                  </div>
                  {(importRes.skipped?.length || 0) > 0 && (
                    <div style={{ marginTop: 10 }}>
                      <div className="muted small" style={{ marginBottom: 6 }}>
                        Skipped
                      </div>
                      <table className="table">
                        <thead>
                          <tr>
                            <th>Titre</th>
                            <th>Raison</th>
                          </tr>
                        </thead>
                        <tbody>
                          {importRes.skipped.slice(0, 30).map((s, idx) => (
                            <tr key={idx}>
                              <td>{s.title}</td>
                              <td className="muted small">{s.reason}</td>
                            </tr>
                          ))}
                        </tbody>
                      </table>
                    </div>
                  )}
                </div>
              )}
            </>
          )}
        </div>

        <div className="card">
          {tab === 'jobs' && (
            <>
              <div className="row" style={{ justifyContent: 'space-between' }}>
                <div>
                  <div style={{ fontWeight: 600 }}>Jobs</div>
                  <div className="muted small">Queue + progression + résultats</div>
                </div>
                <div className="row">
                  <button className="btn sm" disabled={busy} onClick={() => wrap(refreshJobs)}>
                    Rafraîchir
                  </button>
                  <button className="btn sm" onClick={() => setShowEvents((v) => !v)}>
                    {showEvents ? 'Masquer SSE' : 'Voir SSE'}
                  </button>
                </div>
              </div>

              <table className="table" style={{ marginTop: 10 }}>
                <thead>
                  <tr>
                    <th>État</th>
                    <th>Type</th>
                    <th>Progress</th>
                    <th></th>
                  </tr>
                </thead>
                <tbody>
                  {jobs
                    .slice()
                    .sort((a, b) => (a.createdAt < b.createdAt ? 1 : -1))
                    .slice(0, 80)
                    .map((j) => (
                      <React.Fragment key={j.id}>
                        <tr>
                          <td>
                            {badgeForState(j.state)}
                            <div className="muted small">{fmtWhen(j.createdAt)}</div>
                          </td>
                          <td>
                            <div className="row" style={{ justifyContent: 'space-between' }}>
                              <div style={{ fontWeight: 600 }}>{j.type}</div>
                              <button
                                className="btn sm"
                                onClick={() => setJobDetailsID((cur) => (cur === j.id ? null : j.id))}
                              >
                                {jobDetailsID === j.id ? 'Fermer' : 'Détails'}
                              </button>
                            </div>

                            {j.type === 'download' && (() => {
                              const info = jobDownloadInfo(j);
                              return (
                                <div className="muted small" style={{ marginTop: 6 }}>
                                  <div className="row" style={{ gap: 8 }}>
                                    {info.mode && <span className="badge">{info.mode}</span>}
                                    {Number.isFinite(info.bytes as any) && <span className="badge ok">{fmtBytes(info.bytes)}</span>}
                                    {info.path && <span className="badge">{info.path.split('/').slice(-2).join('/')}</span>}
                                  </div>
                                  {info.url && (
                                    <div className="row" style={{ marginTop: 6, justifyContent: 'space-between' }}>
                                      <div className="muted small" style={{ wordBreak: 'break-all' }}>{info.url}</div>
                                      <a className="btn sm" href={info.url} target="_blank" rel="noreferrer">Ouvrir</a>
                                    </div>
                                  )}
                                  {!info.url && info.sourceUrl && (
                                    <div className="muted small" style={{ wordBreak: 'break-all', marginTop: 6 }}>{info.sourceUrl}</div>
                                  )}
                                </div>
                              );
                            })()}

                            {j.error && (
                              <div className="muted small" style={{ color: 'var(--err)', whiteSpace: 'pre-wrap' }}>
                                {j.error}
                              </div>
                            )}
                            {j.errorCode && <div className="muted small">code: {j.errorCode}</div>}
                          </td>
                          <td className="small">
                            {pct(j.progress)}%
                            <div className="bar" title="progression">
                              <div className="barfill" style={{ width: `${pct(j.progress)}%` }} />
                            </div>
                          </td>
                          <td>
                            <div className="row" style={{ justifyContent: 'flex-end' }}>
                              {(j.state === 'queued' || j.state === 'running' || j.state === 'muxing') && (
                                <button
                                  className="btn sm"
                                  disabled={busy}
                                  onClick={() =>
                                    wrap(async () => {
                                      await apiCancelJob(j.id);
                                      await refreshJobs();
                                    })
                                  }
                                >
                                  Cancel
                                </button>
                              )}
                            </div>
                          </td>
                        </tr>

                        {jobDetailsID === j.id && (
                          <tr>
                            <td colSpan={4}>
                              <div className="row" style={{ justifyContent: 'space-between', marginBottom: 8 }}>
                                <div className="muted small">id: {j.id}</div>
                                <div className="row">
                                  <button className="btn sm" onClick={() => navigator.clipboard?.writeText(j.id)}>
                                    Copier id
                                  </button>
                                  {j.type === 'download' && (() => {
                                    const info = jobDownloadInfo(j);
                                    return (
                                      <>
                                        {info.path && (
                                          <button className="btn sm" onClick={() => navigator.clipboard?.writeText(info.path!)}>
                                            Copier path
                                          </button>
                                        )}
                                        {info.url && (
                                          <button className="btn sm" onClick={() => navigator.clipboard?.writeText(info.url!)}>
                                            Copier url
                                          </button>
                                        )}
                                      </>
                                    );
                                  })()}
                                </div>
                              </div>
                              <div className="grid2">
                                <div>
                                  <div className="muted small" style={{ marginBottom: 6 }}>
                                    Params
                                  </div>
                                  <pre className="codeblock">{pretty(j.params) || '—'}</pre>
                                </div>
                                <div>
                                  <div className="muted small" style={{ marginBottom: 6 }}>
                                    Result
                                  </div>
                                  <pre className="codeblock">{pretty(j.result) || '—'}</pre>
                                </div>
                              </div>
                            </td>
                          </tr>
                        )}
                      </React.Fragment>
                    ))}
                  {!jobs.length && (
                    <tr>
                      <td colSpan={4} className="muted">
                        Aucun job.
                      </td>
                    </tr>
                  )}
                </tbody>
              </table>

              {showEvents && (
                <div style={{ marginTop: 12 }}>
                  <div className="row" style={{ justifyContent: 'space-between', marginBottom: 6 }}>
                    <div className="muted small">Évènements (SSE)</div>
                    <button className="btn sm" onClick={() => setEvents([])}>
                      Clear
                    </button>
                  </div>
                  <div className="log">
                    {events.map((l, idx) => (
                      <div key={idx} className="logline">
                        {l}
                      </div>
                    ))}
                  </div>
                </div>
              )}
            </>
          )}

          {tab === 'settings' && (
            <>
              <div className="row" style={{ justifyContent: 'space-between' }}>
                <div>
                  <div style={{ fontWeight: 600 }}>Settings</div>
                  <div className="muted small">Modifie destination, workers, AniList token, etc.</div>
                </div>
                <div className="row">
                  <button className="btn sm" disabled={busy} onClick={() => wrap(refreshSettings)}>
                    Recharger
                  </button>
                  <button
                    className="btn sm primary"
                    disabled={busy || !settingsDraft}
                    onClick={() =>
                      wrap(async () => {
                        const saved = await apiPutSettings(settingsDraft!);
                        setSettings(saved);
                        setSettingsDraft(saved);
                      })
                    }
                  >
                    Sauver
                  </button>
                </div>
              </div>

              {!settingsDraft && <div className="muted" style={{ marginTop: 10 }}>Chargement…</div>}
              {settingsDraft && (
                <>
                  <div className="row" style={{ marginTop: 10 }}>
                    <label className="small muted">destination</label>
                    <input className="input" value={settingsDraft.destination} onChange={(e) => setSettingsDraft((p) => ({ ...(p as Settings), destination: e.target.value }))} />
                  </div>

                  <div className="row" style={{ marginTop: 10 }}>
                    <label className="small muted">maxWorkers</label>
                    <input className="select" type="number" min={1} value={settingsDraft.maxWorkers} onChange={(e) => setSettingsDraft((p) => ({ ...(p as Settings), maxWorkers: Number(e.target.value) || 1 }))} />
                    <label className="small muted">maxConcurrentDownloads</label>
                    <input
                      className="select"
                      type="number"
                      min={1}
                      value={settingsDraft.maxConcurrentDownloads}
                      onChange={(e) => setSettingsDraft((p) => ({ ...(p as Settings), maxConcurrentDownloads: Number(e.target.value) || 1 }))}
                    />
                    <label className="small muted">naming</label>
                    <select
                      className="select"
                      value={settingsDraft.outputNamingMode}
                      onChange={(e) => setSettingsDraft((p) => ({ ...(p as Settings), outputNamingMode: e.target.value as Settings['outputNamingMode'] }))}
                    >
                      <option value="legacy">legacy</option>
                      <option value="media-server">media-server</option>
                    </select>
                    <label className="small" style={{ display: 'flex', gap: 6, alignItems: 'center' }}>
                      <input type="checkbox" checked={settingsDraft.separateLang} onChange={(e) => setSettingsDraft((p) => ({ ...(p as Settings), separateLang: e.target.checked }))} />
                      separateLang
                    </label>
                  </div>

                  <div className="row" style={{ marginTop: 10 }}>
                    <label className="small muted">anilistToken</label>
                    <input
                      className="input"
                      placeholder="(optionnel)"
                      value={settingsDraft.anilistToken || ''}
                      onChange={(e) => setSettingsDraft((p) => ({ ...(p as Settings), anilistToken: e.target.value }))}
                    />
                  </div>

                  {settings && (
                    <div className="muted small" style={{ marginTop: 10 }}>
                      Actuel: destination={settings.destination}, workers={settings.maxWorkers}, downloads={settings.maxConcurrentDownloads}
                    </div>
                  )}
                </>
              )}
            </>
          )}

          {tab !== 'jobs' && tab !== 'settings' && (
            <>
              <div style={{ fontWeight: 600 }}>Jobs (aperçu)</div>
              <div className="muted small">Les jobs se rafraîchissent via SSE.</div>
              <table className="table" style={{ marginTop: 10 }}>
                <thead>
                  <tr>
                    <th>État</th>
                    <th>Type</th>
                    <th>Progress</th>
                  </tr>
                </thead>
                <tbody>
                  {jobs
                    .slice()
                    .sort((a, b) => (a.createdAt < b.createdAt ? 1 : -1))
                    .slice(0, 8)
                    .map((j) => (
                      <tr key={j.id}>
                        <td>{badgeForState(j.state)}</td>
                        <td>{j.type}</td>
                        <td className="small">{Math.round((j.progress || 0) * 100)}%</td>
                      </tr>
                    ))}
                  {!jobs.length && (
                    <tr>
                      <td colSpan={3} className="muted">
                        Aucun job.
                      </td>
                    </tr>
                  )}
                </tbody>
              </table>
            </>
          )}
        </div>
      </div>
    </div>
  );
}
