import React, { useEffect, useMemo, useState } from 'react';
import {
  apiAniListAiring,
  apiAniListImportAuto,
  apiAniListViewer,
  apiAnimeSamaEnqueue,
  apiAnimeSamaEpisodes,
  apiAnimeSamaResolve,
  apiAnimeSamaScan,
  apiCancelJob,
  apiCreateSubscription,
  apiDeleteSubscription,
  apiEnqueueSubscriptionEpisodes,
  apiGetSettings,
  apiGetSubscriptionEpisodes,
  apiListJobs,
  apiListSubscriptions,
  apiPutSettings,
  apiSyncAll,
  apiSyncSubscription,
  type AniListAiringScheduleEntry,
  type AniListImportAutoRequest,
  type AniListImportAutoResponse,
  type AniListViewer,
  type AnimeSamaEnqueueResponse,
  type AnimeSamaEpisodesResponse,
  type AnimeSamaResolvedCandidate,
  type AnimeSamaScanOption,
  type Job,
  type Settings,
  type SubscriptionEnqueueEpisodesResponse,
  type SubscriptionEpisodesResponse,
  type Subscription,
} from './api';

type Tab = 'recherche' | 'abonnements' | 'calendrier' | 'jobs' | 'settings';

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

function jobBucket(state: Job['state']): 'running' | 'queue' | 'done' | 'failed' {
  if (state === 'queued') return 'queue';
  if (state === 'running' || state === 'muxing') return 'running';
  if (state === 'completed' || state === 'canceled') return 'done';
  return 'failed';
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
  const [busyLabel, setBusyLabel] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [events, setEvents] = useState<string[]>([]);
  const [toast, setToast] = useState<null | { kind: 'ok' | 'info' | 'err'; message: string }>(null);

  const [settings, setSettings] = useState<Settings | null>(null);
  const [settingsDraft, setSettingsDraft] = useState<Settings | null>(null);

  const [subs, setSubs] = useState<Subscription[]>([]);
  const [jobs, setJobs] = useState<Job[]>([]);

  const [subsFilter, setSubsFilter] = useState('');
  const [subsSort, setSubsSort] = useState<SubSort>('nextCheckAt');

  const [calendarMode, setCalendarMode] = useState<CalendarMode>('checks');

  const [jobDetailsID, setJobDetailsID] = useState<string | null>(null);
  const [showEvents, setShowEvents] = useState(false);

  const [newBaseUrl, setNewBaseUrl] = useState('');
  const [searchTitle, setSearchTitle] = useState('');
  const [searchCandidates, setSearchCandidates] = useState<AnimeSamaResolvedCandidate[]>([]);
  const [searching, setSearching] = useState(false);

  const [searchSubscribe, setSearchSubscribe] = useState(false);
  const [searchDownloadNow, setSearchDownloadNow] = useState(true);

  const [episodesByBaseUrl, setEpisodesByBaseUrl] = useState<
    Record<
      string,
      {
        loading: boolean;
        data: AnimeSamaEpisodesResponse | null;
        error: string | null;
        selected: Set<number>;
        rangeFrom: number;
        rangeTo: number;
        listInput: string;
      }
    >
  >({});

  const [jobsView, setJobsView] = useState<'running' | 'queue' | 'done' | 'failed' | 'all'>('running');

  const [episodesModal, setEpisodesModal] = useState<null | {
    data: SubscriptionEpisodesResponse;
    selected: Set<number>;
  }>(null);

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

  const [scanBySlug, setScanBySlug] = useState<
    Record<
      string,
      {
        loading: boolean;
        options: AnimeSamaScanOption[];
        selectedBaseUrl: string | null;
        error: string | null;
      }
    >
  >({});

  function parseEpisodesListInput(s: string): number[] {
    const out: number[] = [];
    const seen = new Set<number>();
    const parts = String(s)
      .split(/[\s,;]+/g)
      .map((p) => p.trim())
      .filter(Boolean);
    for (const p of parts) {
      const m = p.match(/^([0-9]+)\s*-\s*([0-9]+)$/);
      if (m) {
        let a = Number(m[1]);
        let b = Number(m[2]);
        if (!Number.isFinite(a) || !Number.isFinite(b)) continue;
        if (a <= 0 || b <= 0) continue;
        if (a > b) [a, b] = [b, a];
        for (let i = a; i <= b; i++) {
          if (seen.has(i)) continue;
          seen.add(i);
          out.push(i);
        }
        continue;
      }
      const n = Number(p);
      if (!Number.isFinite(n) || n <= 0) continue;
      if (seen.has(n)) continue;
      seen.add(n);
      out.push(n);
    }
    out.sort((a, b) => a - b);
    return out;
  }

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

  function pickDefaultScanOption(options: AnimeSamaScanOption[]) {
    const preferred = options.find((o) => o.season === resolveSeason && o.lang === resolveLang);
    return preferred || options[0] || null;
  }

  async function scanCandidate(c: AnimeSamaResolvedCandidate) {
    const cur = scanBySlug[c.slug];
    if (cur?.loading) return;
    if (cur?.options?.length) return;

    setScanBySlug((m) => ({
      ...m,
      [c.slug]: { loading: true, options: m[c.slug]?.options || [], selectedBaseUrl: m[c.slug]?.selectedBaseUrl || null, error: null },
    }));

    try {
      const res = await apiAnimeSamaScan({
        catalogueUrl: c.catalogueUrl,
        maxSeason: Math.max(5, resolveSeason),
        langs: [resolveLang, 'vostfr', 'vf'],
      });
      const options = (res.options || []).slice().sort((a, b) => {
        if (a.season !== b.season) return a.season - b.season;
        return a.lang.localeCompare(b.lang);
      });
      const def = pickDefaultScanOption(options);
      setScanBySlug((m) => ({
        ...m,
        [c.slug]: {
          loading: false,
          options,
          selectedBaseUrl: def ? def.baseUrl : null,
          error: null,
        },
      }));
    } catch (e: any) {
      const msg = e?.message || String(e);
      setScanBySlug((m) => ({
        ...m,
        [c.slug]: { loading: false, options: [], selectedBaseUrl: null, error: msg },
      }));
    }
  }

  async function ensureEpisodesLoaded(baseUrl: string) {
    const cur = episodesByBaseUrl[baseUrl];
    if (cur?.loading) return;
    if (cur?.data) return;

    setEpisodesByBaseUrl((m) => ({
      ...m,
      [baseUrl]: {
        loading: true,
        data: m[baseUrl]?.data || null,
        error: null,
        selected: m[baseUrl]?.selected || new Set<number>(),
        rangeFrom: m[baseUrl]?.rangeFrom || 1,
        rangeTo: m[baseUrl]?.rangeTo || 1,
        listInput: m[baseUrl]?.listInput || '',
      },
    }));

    try {
      const data = await apiAnimeSamaEpisodes({ baseUrl });
      const selected = new Set<number>();
      for (const e of data.episodes || []) {
        if (e.available) selected.add(e.episode);
      }
      const max = data.maxAvailableEpisode || 1;
      setEpisodesByBaseUrl((m) => ({
        ...m,
        [baseUrl]: {
          loading: false,
          data,
          error: null,
          selected,
          rangeFrom: 1,
          rangeTo: max,
          listInput: '',
        },
      }));
    } catch (e: any) {
      const msg = e?.message || String(e);
      setEpisodesByBaseUrl((m) => ({
        ...m,
        [baseUrl]: {
          loading: false,
          data: null,
          error: msg,
          selected: m[baseUrl]?.selected || new Set<number>(),
          rangeFrom: m[baseUrl]?.rangeFrom || 1,
          rangeTo: m[baseUrl]?.rangeTo || 1,
          listInput: m[baseUrl]?.listInput || '',
        },
      }));
    }
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

  async function openEpisodesModal(sub: Subscription) {
    const data = await apiGetSubscriptionEpisodes(sub.id);
    const selected = new Set<number>();
    for (const e of data.episodes || []) {
      if (e.available && !e.downloaded) selected.add(e.episode);
    }
    setEpisodesModal({ data, selected });
  }

  async function runSearch() {
    await wrapAction('Recherche Anime‑Sama…', async () => {
      const q = searchTitle.trim();
      if (!q) {
        notify('info', 'Entre un nom d’anime.');
        return;
      }
      setSearching(true);
      setSearchCandidates([]);
      try {
        const res = await apiAnimeSamaResolve({ titles: [q], season: resolveSeason, lang: resolveLang, maxCandidates: resolveMaxCandidates });
        setSearchCandidates(res.candidates || []);
        notify('info', `Résultats: ${(res.candidates || []).length}`);
      } finally {
        setSearching(false);
      }
    });
  }

  // Scan automatique des candidats (onglet Recherche + résultats planning AniList).
  useEffect(() => {
    const cands: AnimeSamaResolvedCandidate[] = [];
    if (tab === 'recherche') cands.push(...(searchCandidates || []));
    if (resolved?.candidates?.length) cands.push(...resolved.candidates);
    for (const c of cands) {
      void scanCandidate(c);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [tab, searchCandidates, resolved]);

  // Une fois scanné + baseUrl choisie, on charge automatiquement la liste d'épisodes.
  useEffect(() => {
    const bases: string[] = [];
    if (tab === 'recherche') {
      for (const c of searchCandidates || []) {
        const baseUrl = scanBySlug[c.slug]?.selectedBaseUrl;
        if (baseUrl) bases.push(baseUrl);
      }
    }
    if (resolved?.candidates?.length) {
      for (const c of resolved.candidates) {
        const baseUrl = scanBySlug[c.slug]?.selectedBaseUrl;
        if (baseUrl) bases.push(baseUrl);
      }
    }
    for (const b of bases) {
      void ensureEpisodesLoaded(b);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [tab, searchCandidates, resolved, scanBySlug]);

  function notify(kind: 'ok' | 'info' | 'err', message: string) {
    setToast({ kind, message });
  }

  useEffect(() => {
    if (!toast) return;
    const t = window.setTimeout(() => setToast(null), 3200);
    return () => window.clearTimeout(t);
  }, [toast]);

  useEffect(() => {
    if (!episodesModal) return;
    const onKeyDown = (ev: KeyboardEvent) => {
      if (ev.key === 'Escape') setEpisodesModal(null);
    };
    window.addEventListener('keydown', onKeyDown);
    return () => window.removeEventListener('keydown', onKeyDown);
  }, [episodesModal]);

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
      notify('err', e?.message || String(e));
      throw e;
    } finally {
      setBusy(false);
      setBusyLabel(null);
    }
  }

  async function wrapAction<T>(label: string, fn: () => Promise<T>) {
    setBusyLabel(label);
    return await wrap(fn);
  }

  const tabLabels: Record<Tab, string> = {
    recherche: 'Recherche',
    abonnements: 'Abonnements',
    calendrier: 'Planning',
    jobs: 'Téléchargements',
    settings: 'Réglages',
  };

  return (
    <div className="container">
      <header className="topbar">
        <div className="row" style={{ justifyContent: 'space-between', alignItems: 'flex-start' }}>
          <div className="brandTitle">
            <div className="h1">Anime-Sama Downloader</div>
            <div className="muted small">Abonnements • AniList • Téléchargements • Logs</div>
          </div>
          <nav className="nav" aria-label="Navigation">
            {(['recherche', 'abonnements', 'calendrier', 'jobs', 'settings'] as Tab[]).map((t) => (
              <button
                key={t}
                className={`btn sm navbtn ${tab === t ? 'active' : ''}`}
                onClick={() => setTab(t)}
                aria-current={tab === t ? 'page' : undefined}
              >
                {tabLabels[t]}
              </button>
            ))}
          </nav>
        </div>
      </header>

      {toast && (
        <div className={`toast ${toast.kind}`} role="status" aria-live="polite" onClick={() => setToast(null)}>
          {toast.message}
        </div>
      )}

      {busy && busyLabel && (
        <div className="card soft" style={{ marginTop: 12 }}>
          <div className="row" style={{ justifyContent: 'space-between' }}>
            <div>
              <div className="sectionTitle">Action en cours</div>
              <div className="muted small" style={{ marginTop: 4 }}>{busyLabel}</div>
            </div>
            <button className="btn sm" onClick={() => notify('info', 'Patiente…')}>OK</button>
          </div>
        </div>
      )}

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

      {episodesModal && (
        <div
          className="modalBackdrop"
          onMouseDown={(ev) => {
            if (ev.target === ev.currentTarget) setEpisodesModal(null);
          }}
        >
          <div className="modal card" onMouseDown={(ev) => ev.stopPropagation()}>
          <div className="row" style={{ justifyContent: 'space-between', alignItems: 'flex-start' }}>
            <div>
              <div style={{ fontWeight: 700 }}>Épisodes</div>
              <div className="muted small" style={{ marginTop: 4 }}>
                {episodesModal.data.subscription.label} — player: {episodesModal.data.selectedPlayer || 'auto'} — dispo max: {episodesModal.data.maxAvailableEpisode}
              </div>
            </div>
            <div className="row" style={{ gap: 8 }}>
              <button className="btn sm" disabled={busy} onClick={() => setEpisodesModal(null)}>
                Fermer
              </button>
              <button
                className="btn sm"
                disabled={busy}
                onClick={() =>
                  setEpisodesModal((m) => {
                    if (!m) return m;
                    const s = new Set<number>();
                    for (const e of m.data.episodes || []) if (e.available && !e.downloaded) s.add(e.episode);
                    return { ...m, selected: s };
                  })
                }
              >
                Tout (non téléchargés)
              </button>
              <button className="btn sm" disabled={busy} onClick={() => setEpisodesModal((m) => (m ? { ...m, selected: new Set() } : m))}>
                Rien
              </button>
              <button
                className="btn sm primary"
                disabled={busy}
                onClick={() =>
                  wrapAction('Planification des téléchargements…', async () => {
                    const eps = Array.from(episodesModal.selected.values()).sort((a, b) => a - b);
                    const res: SubscriptionEnqueueEpisodesResponse = await apiEnqueueSubscriptionEpisodes(episodesModal.data.subscription.id, eps);
                    setEpisodesModal(null);
                    await refreshSubs();
                    await refreshJobs();
                    notify('ok', `Téléchargements planifiés: ${(res.enqueuedEpisodes || []).length}`);
                    if ((res.skipped || []).length) {
                      setError(`Certains épisodes ont été ignorés: ${res.skipped.map((s) => `#${s.episode} (${s.reason})`).join(', ')}`);
                    }
                  })
                }
              >
                Télécharger la sélection
              </button>
            </div>
          </div>

          <div style={{ marginTop: 10, maxHeight: 320, overflow: 'auto' }}>
            <table className="table">
              <thead>
                <tr>
                  <th></th>
                  <th>Épisode</th>
                  <th>Statut</th>
                </tr>
              </thead>
              <tbody>
                {(episodesModal.data.episodes || []).map((e) => {
                  const checked = episodesModal.selected.has(e.episode);
                  const disabled = !e.available;
                  const status = e.downloaded ? 'téléchargé' : e.scheduled ? 'planifié' : e.available ? 'disponible' : 'indispo';
                  return (
                    <tr key={e.episode}>
                      <td style={{ width: 40 }}>
                        <input
                          type="checkbox"
                          disabled={disabled || busy}
                          checked={checked}
                          onChange={(ev) => {
                            const v = ev.target.checked;
                            setEpisodesModal((m) => {
                              if (!m) return m;
                              const s = new Set(m.selected);
                              if (v) s.add(e.episode);
                              else s.delete(e.episode);
                              return { ...m, selected: s };
                            });
                          }}
                        />
                      </td>
                      <td className="small">#{e.episode}</td>
                      <td className="small">{status}</td>
                    </tr>
                  );
                })}
              </tbody>
            </table>
          </div>
          </div>
        </div>
      )}

      <div className="grid" style={{ marginTop: 12 }}>
        <div className="card">
          {tab === 'recherche' && (
            <>
              <div className="row" style={{ justifyContent: 'space-between' }}>
                <div>
                  <div className="sectionTitle">Recherche</div>
                  <div className="muted small sectionHint">Recherche → scan automatique → choix saison/lang → choix épisodes → télécharger (ou s’abonner).</div>
                </div>
                <div className="row">
                  <label className="small muted">S’abonner</label>
                  <input type="checkbox" disabled={busy} checked={searchSubscribe} onChange={(e) => setSearchSubscribe(e.target.checked)} />
                  {searchSubscribe && (
                    <>
                      <label className="small muted">Télécharger maintenant</label>
                      <input type="checkbox" disabled={busy} checked={searchDownloadNow} onChange={(e) => setSearchDownloadNow(e.target.checked)} />
                    </>
                  )}
                </div>
              </div>

              <div className="card" style={{ marginTop: 10 }}>
                <div className="cardTitle">Trouver un anime</div>
                <div className="muted small cardHint">Tape un nom. Le scan est automatique (pas de bouton). Les épisodes n’apparaissent qu’après scan.</div>

                <div className="row">
                  <input
                    className="input"
                    placeholder="Nom de l'anime (ex: Solo Leveling)"
                    value={searchTitle}
                    onChange={(e) => setSearchTitle(e.target.value)}
                    onKeyDown={(e) => {
                      if (e.key === 'Enter') runSearch().catch(() => void 0);
                    }}
                  />
                  <button className="btn sm primary" disabled={busy || searching || searchTitle.trim() === ''} onClick={() => runSearch().catch(() => void 0)}>
                    {searching ? 'Recherche…' : 'Chercher'}
                  </button>
                  <details>
                    <summary className="btn sm" style={{ display: 'inline-flex' }}>Préférences</summary>
                    <div className="row" style={{ marginTop: 10 }}>
                      <label className="small muted">Langue préférée</label>
                      <select className="select" value={resolveLang} onChange={(e) => setResolveLang(e.target.value)}>
                        <option value="vostfr">vostfr</option>
                        <option value="vf">vf</option>
                      </select>
                      <label className="small muted">Saison préférée</label>
                      <input className="select" type="number" min={1} value={resolveSeason} onChange={(e) => setResolveSeason(Number(e.target.value) || 1)} />
                      <label className="small muted">Résultats</label>
                      <input className="select" type="number" min={1} max={10} value={resolveMaxCandidates} onChange={(e) => setResolveMaxCandidates(Number(e.target.value) || 5)} />
                    </div>
                  </details>
                </div>

                {searchTitle.trim() !== '' && (
                  <div className="muted small" style={{ marginTop: 10 }}>
                    {searching ? 'Recherche en cours…' : `Résultats affichés: ${searchCandidates.length}`}
                  </div>
                )}

                {!!searchCandidates.length && (
                  <table className="table" style={{ marginTop: 10 }}>
                    <thead>
                      <tr>
                        <th>Score</th>
                        <th>Match</th>
                        <th>Saison / Lang</th>
                        <th>Épisodes</th>
                        <th></th>
                      </tr>
                    </thead>
                    <tbody>
                      {searchCandidates.map((c) => {
                        const st = scanBySlug[c.slug];
                        const baseUrl = st?.selectedBaseUrl || null;
                        const epState = baseUrl ? episodesByBaseUrl[baseUrl] : null;
                        const selectedList = baseUrl && epState ? Array.from(epState.selected).sort((a, b) => a - b) : [];

                        const canAct = !!baseUrl && !!epState?.data && selectedList.length > 0;

                        return (
                          <tr key={c.catalogueUrl + '|' + c.slug}>
                            <td className="small">{Math.round((c.score || 0) * 100)}%</td>
                            <td className="small">{c.matchedTitle || c.slug}</td>
                            <td className="small">
                              {!st ? (
                                <span className="muted">Scan en attente…</span>
                              ) : st.loading ? (
                                <span className="muted">Scan en cours…</span>
                              ) : st.error ? (
                                <span className="bad">{st.error}</span>
                              ) : !st.options.length ? (
                                <span className="muted">Aucune saison trouvée</span>
                              ) : (
                                <div className="row" style={{ gap: 8 }}>
                                  <span className="muted small">Saison</span>
                                  <select
                                    className="select"
                                    value={st.selectedBaseUrl || st.options[0].baseUrl}
                                    onChange={(e) => {
                                      const v = e.target.value;
                                      setScanBySlug((m) => ({
                                        ...m,
                                        [c.slug]: { ...m[c.slug], selectedBaseUrl: v },
                                      }));
                                    }}
                                  >
                                    {st.options.map((o) => (
                                      <option key={o.baseUrl} value={o.baseUrl}>
                                        S{o.season} {o.lang.toUpperCase()} (max ep {o.maxAvailableEpisode})
                                      </option>
                                    ))}
                                  </select>
                                </div>
                              )}
                            </td>

                            <td className="small">
                              {!baseUrl ? (
                                <span className="muted">—</span>
                              ) : !epState ? (
                                <span className="muted">Chargement…</span>
                              ) : epState.loading ? (
                                <span className="muted">Chargement…</span>
                              ) : epState.error ? (
                                <span className="bad">{epState.error}</span>
                              ) : !epState.data ? (
                                <span className="muted">—</span>
                              ) : (
                                <div style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
                                  <div className="muted small">
                                    Sélection: {selectedList.length} épisode(s) (max {epState.data.maxAvailableEpisode})
                                  </div>
                                  <div className="row" style={{ gap: 8, flexWrap: 'wrap', justifyContent: 'flex-start' }}>
                                    <button
                                      className="btn sm"
                                      disabled={busy}
                                      onClick={() => {
                                        const data = epState.data;
                                        if (!data) return;
                                        const next = new Set<number>();
                                        for (const e of data.episodes || []) if (e.available) next.add(e.episode);
                                        setEpisodesByBaseUrl((m) => ({
                                          ...m,
                                          [baseUrl]: { ...m[baseUrl], selected: next },
                                        }));
                                      }}
                                    >
                                      Tout
                                    </button>
                                    <button
                                      className="btn sm"
                                      disabled={busy}
                                      onClick={() => {
                                        setEpisodesByBaseUrl((m) => ({
                                          ...m,
                                          [baseUrl]: { ...m[baseUrl], selected: new Set<number>() },
                                        }));
                                      }}
                                    >
                                      Aucun
                                    </button>
                                    <span className="muted small">Plage</span>
                                    <input
                                      className="select"
                                      style={{ width: 80 }}
                                      type="number"
                                      min={1}
                                      value={epState.rangeFrom}
                                      onChange={(e) => {
                                        const v = Number(e.target.value) || 1;
                                        setEpisodesByBaseUrl((m) => ({
                                          ...m,
                                          [baseUrl]: { ...m[baseUrl], rangeFrom: v },
                                        }));
                                      }}
                                    />
                                    <span className="muted small">→</span>
                                    <input
                                      className="select"
                                      style={{ width: 80 }}
                                      type="number"
                                      min={1}
                                      value={epState.rangeTo}
                                      onChange={(e) => {
                                        const v = Number(e.target.value) || 1;
                                        setEpisodesByBaseUrl((m) => ({
                                          ...m,
                                          [baseUrl]: { ...m[baseUrl], rangeTo: v },
                                        }));
                                      }}
                                    />
                                    <button
                                      className="btn sm"
                                      disabled={busy}
                                      onClick={() => {
                                        const data = epState.data;
                                        if (!data) return;
                                        let a = epState.rangeFrom || 1;
                                        let b = epState.rangeTo || 1;
                                        if (a > b) [a, b] = [b, a];
                                        const avail = new Set<number>();
                                        for (const e of data.episodes || []) if (e.available) avail.add(e.episode);
                                        const next = new Set<number>();
                                        for (let i = a; i <= b; i++) if (avail.has(i)) next.add(i);
                                        setEpisodesByBaseUrl((m) => ({
                                          ...m,
                                          [baseUrl]: { ...m[baseUrl], selected: next },
                                        }));
                                      }}
                                    >
                                      Appliquer
                                    </button>
                                  </div>

                                  <div className="row" style={{ gap: 8, flexWrap: 'wrap', justifyContent: 'flex-start' }}>
                                    <span className="muted small">Liste</span>
                                    <input
                                      className="input"
                                      style={{ minWidth: 220 }}
                                      placeholder="ex: 1,2,5-8"
                                      value={epState.listInput}
                                      onChange={(e) => {
                                        const v = e.target.value;
                                        setEpisodesByBaseUrl((m) => ({
                                          ...m,
                                          [baseUrl]: { ...m[baseUrl], listInput: v },
                                        }));
                                      }}
                                    />
                                    <button
                                      className="btn sm"
                                      disabled={busy}
                                      onClick={() => {
                                        const data = epState.data;
                                        if (!data) return;
                                        const wanted = parseEpisodesListInput(epState.listInput);
                                        const avail = new Set<number>();
                                        for (const e of data.episodes || []) if (e.available) avail.add(e.episode);
                                        const next = new Set<number>();
                                        for (const n of wanted) if (avail.has(n)) next.add(n);
                                        setEpisodesByBaseUrl((m) => ({
                                          ...m,
                                          [baseUrl]: { ...m[baseUrl], selected: next },
                                        }));
                                      }}
                                    >
                                      Appliquer
                                    </button>
                                  </div>
                                </div>
                              )}
                            </td>

                            <td>
                              <div className="row" style={{ justifyContent: 'flex-end' }}>
                                <a className="btn sm" href={c.catalogueUrl} target="_blank" rel="noreferrer">
                                  Catalogue
                                </a>
                                {baseUrl && (
                                  <a className="btn sm" href={baseUrl} target="_blank" rel="noreferrer">
                                    Ouvrir saison
                                  </a>
                                )}
                                <button
                                  className="btn sm primary"
                                  disabled={busy || !canAct}
                                  onClick={() =>
                                    wrapAction(searchSubscribe ? 'Création de l’abonnement…' : 'Mise en téléchargement…', async () => {
                                      const episodes = selectedList;
                                      if (!baseUrl) throw new Error('Scan en cours: baseUrl manquante');
                                      if (!episodes.length) throw new Error('Sélectionne au moins un épisode.');

                                      if (searchSubscribe) {
                                        const sub = await apiCreateSubscription({ baseUrl });
                                        if (searchDownloadNow) {
                                          const enq: SubscriptionEnqueueEpisodesResponse = await apiEnqueueSubscriptionEpisodes(sub.id, episodes);
                                          notify('ok', `Abonné + téléchargements: ${enq.enqueuedEpisodes.length} épisode(s).`);
                                        } else {
                                          notify('ok', 'Abonnement ajouté.');
                                        }
                                        await refreshSubs();
                                        await refreshJobs();
                                        return;
                                      }

                                      const enq: AnimeSamaEnqueueResponse = await apiAnimeSamaEnqueue({ baseUrl, episodes });
                                      await refreshJobs();
                                      notify('ok', `Téléchargements: ${enq.enqueuedEpisodes.length} épisode(s) mis en file.`);
                                    })
                                  }
                                >
                                  {searchSubscribe ? 'S’abonner' : 'Télécharger'}
                                </button>
                              </div>
                            </td>
                          </tr>
                        );
                      })}
                    </tbody>
                  </table>
                )}
              </div>

              <details style={{ marginTop: 12 }}>
                <summary className="btn sm" style={{ display: 'inline-flex' }}>URL directe (avancé)</summary>
                <div className="muted small" style={{ marginTop: 8 }}>
                  Colle une URL de saison Anime‑Sama (ex: `/catalogue/.../saison1/vostfr/`). Le scan n’est pas nécessaire ici.
                </div>
                <div className="row" style={{ marginTop: 10 }}>
                  <input
                    className="input"
                    placeholder="Base URL (Anime-Sama)"
                    value={newBaseUrl}
                    onChange={(e) => setNewBaseUrl(e.target.value)}
                    onKeyDown={(e) => {
                      if (e.key === 'Enter') {
                        const v = newBaseUrl.trim();
                        if (!v) return;
                        wrapAction('Chargement des épisodes…', async () => {
                          // On charge les épisodes et pré-sélectionne tout, puis l'utilisateur clique Télécharger.
                          await ensureEpisodesLoaded(v);
                          notify('info', 'Épisodes chargés: sélectionne puis télécharge.');
                        }).catch(() => void 0);
                      }
                    }}
                  />
                </div>
              </details>
            </>
          )}

          {tab === 'abonnements' && (
            <>
              <div className="row" style={{ justifyContent: 'space-between' }}>
                <div>
                  <div className="sectionTitle">Abonnements</div>
                  <div className="muted small sectionHint">Ajoute via “Recherche”. Ici: suivi et sélection d’épisodes sur tes abonnements.</div>
                </div>
                <div className="row">
                  <button className="btn sm" disabled={busy} onClick={() => wrapAction('Rafraîchissement des abonnements…', refreshSubs)}>
                    Rafraîchir
                  </button>
                  <button
                    className="btn sm primary"
                    disabled={busy}
                    onClick={() =>
                      wrapAction('Synchronisation des abonnements “à vérifier”…', async () => {
                        await apiSyncAll({ enqueue: true, dueOnly: true, limit: 500 });
                        await refreshSubs();
                        await refreshJobs();
                        notify('ok', 'Sync lancée (à vérifier).');
                      })
                    }
                  >
                    Vérifier maintenant
                  </button>
                </div>
              </div>

              <div className="card" style={{ marginTop: 10 }}>
                <div className="cardTitle">Ajouter / télécharger</div>
                <div className="muted small cardHint">Pour ajouter un nouvel anime ou lancer un téléchargement ponctuel: passe par l’onglet “Recherche”.</div>
                <div className="row" style={{ marginTop: 10 }}>
                  <button className="btn sm primary" disabled={busy} onClick={() => setTab('recherche')}>Ouvrir Recherche</button>
                </div>
              </div>

              <div className="row" style={{ marginTop: 10, justifyContent: 'space-between' }}>
                <input
                  className="input"
                  placeholder="Rechercher dans mes abonnements (nom ou URL)"
                  value={subsFilter}
                  onChange={(e) => setSubsFilter(e.target.value)}
                />
                <select className="select" value={subsSort} onChange={(e) => setSubsSort(e.target.value as SubSort)}>
                  <option value="nextCheckAt">Tri: prochain check</option>
                  <option value="label">Tri: nom</option>
                </select>
              </div>

              <div className="cards">
                {subsSorted.map((s) => (
                  <div key={s.id} className="cardItem">
                    <div className="cardRow">
                      <div>
                        <div style={{ fontWeight: 800, letterSpacing: '.2px' }}>{s.label}</div>
                        <div className="muted small" style={{ marginTop: 4 }}>Prochaine vérif: {fmtWhen(s.nextCheckAt)} ({fmtRelative(s.nextCheckAt)})</div>
                        <div className="muted small" style={{ marginTop: 4 }}>Dernière vérif: {fmtWhen(s.lastCheckedAt)}</div>
                        <div className="airMeta">
                          {isDue(s.nextCheckAt) ? <span className="pill warn">À vérifier</span> : <span className="pill ok">OK</span>}
                          <span className="pill">Dispo {s.lastAvailableEpisode || 0}</span>
                          <span className="pill">Téléchargé {s.lastDownloadedEpisode || 0}</span>
                        </div>
                      </div>
                      <div className="row" style={{ justifyContent: 'flex-end' }}>
                        <button className="btn sm primary" disabled={busy} onClick={() => wrapAction('Chargement des épisodes…', () => openEpisodesModal(s))}>
                          Épisodes
                        </button>
                        <details>
                          <summary className="btn sm">Actions</summary>
                          <div className="row" style={{ marginTop: 10 }}>
                            <a className="btn sm" href={s.baseUrl} target="_blank" rel="noreferrer">Ouvrir Anime‑Sama</a>
                            <button
                              className="btn sm"
                              disabled={busy}
                              onClick={() =>
                                wrapAction('Vérification de l’abonnement…', async () => {
                                  await apiSyncSubscription(s.id, true);
                                  await refreshSubs();
                                  await refreshJobs();
                                  notify('ok', 'Vérification lancée.');
                                })
                              }
                            >
                              Vérifier maintenant
                            </button>
                            <button
                              className="btn sm danger"
                              disabled={busy}
                              onClick={() =>
                                wrapAction('Suppression de l’abonnement…', async () => {
                                  await apiDeleteSubscription(s.id);
                                  await refreshSubs();
                                  notify('ok', 'Abonnement supprimé.');
                                })
                              }
                            >
                              Supprimer
                            </button>
                          </div>
                        </details>
                      </div>
                    </div>
                  </div>
                ))}
                {!subsSorted.length && <div className="muted">Aucun abonnement.</div>}
              </div>
            </>
          )}

          {tab === 'calendrier' && (
            <>
              <div className="row" style={{ justifyContent: 'space-between' }}>
                <div>
                  <div className="sectionTitle">Planning</div>
                  <div className="muted small sectionHint">Deux vues: vérifications (abonnements) et sorties (AniList).</div>
                </div>
                <div className="row">
                  <button
                    className={`btn sm navbtn ${calendarMode === 'checks' ? 'active' : ''}`}
                    onClick={() => setCalendarMode('checks')}
                  >
                    Vérifications
                  </button>
                  <button
                    className={`btn sm navbtn ${calendarMode === 'anilist' ? 'active' : ''}`}
                    onClick={() => setCalendarMode('anilist')}
                  >
                    Sorties AniList
                  </button>
                </div>
              </div>

              {calendarMode === 'checks' && (
                <>
                  <div className="row" style={{ marginTop: 10, justifyContent: 'space-between' }}>
                    <div className="row">
                      <span className="badge warn">À vérifier maintenant: {subsByCheck.due.length}</span>
                      <span className="badge"><span className="muted">Prochain check:</span> {subsByCheck.soon[0] ? fmtRelative(subsByCheck.soon[0].nextCheckAt) : '—'}</span>
                    </div>
                    <div className="row">
                      <button className="btn sm" disabled={busy} onClick={() => wrapAction('Rafraîchissement des abonnements…', refreshSubs)}>
                        Rafraîchir
                      </button>
                      <button
                        className="btn sm primary"
                        disabled={busy}
                        onClick={() =>
                          wrapAction('Synchronisation des abonnements “à vérifier”…', async () => {
                            await apiSyncAll({ enqueue: true, dueOnly: true, limit: 500 });
                            await refreshSubs();
                            await refreshJobs();
                            notify('ok', 'Sync lancée (à vérifier).');
                          })
                        }
                      >
                        Vérifier maintenant
                      </button>
                    </div>
                  </div>

                  <div className="muted small" style={{ marginTop: 10 }}>
                    “À vérifier maintenant” = abonnements dont la prochaine vérification est dépassée. “Check” = on regarde si de nouveaux épisodes sont sortis.
                  </div>

                  <div className="cards">
                    <div className="cardItem">
                      <div className="cardRow">
                        <div>
                          <div className="sectionTitle">À vérifier maintenant</div>
                          <div className="muted small">Abonnements en retard de vérification.</div>
                        </div>
                        <span className="pill warn">{subsByCheck.due.length}</span>
                      </div>
                      <div style={{ marginTop: 10, display: 'flex', flexDirection: 'column', gap: 10 }}>
                        {subsByCheck.due.slice(0, 8).map((s) => (
                          <div key={s.id} className="airItem">
                            <div className="cardRow">
                              <div>
                                <div style={{ fontWeight: 700 }}>{s.label}</div>
                                <div className="muted small">{fmtRelative(s.nextCheckAt)} • {fmtWhen(s.nextCheckAt)}</div>
                              </div>
                              <div className="row" style={{ justifyContent: 'flex-end' }}>
                                <button className="btn sm" onClick={() => setTab('abonnements')}>Ouvrir</button>
                                <button
                                  className="btn sm primary"
                                  disabled={busy}
                                  onClick={() =>
                                    wrapAction('Vérification de l’abonnement…', async () => {
                                      await apiSyncSubscription(s.id, true);
                                      await refreshSubs();
                                      await refreshJobs();
                                      notify('ok', 'Vérification lancée.');
                                    })
                                  }
                                >
                                  Vérifier
                                </button>
                              </div>
                            </div>
                          </div>
                        ))}
                        {subsByCheck.due.length === 0 && <div className="muted small">Rien à vérifier.</div>}
                        {subsByCheck.due.length > 8 && <div className="muted small">+ {subsByCheck.due.length - 8} autres…</div>}
                      </div>
                    </div>

                    <div className="cardItem">
                      <div className="cardRow">
                        <div>
                          <div className="sectionTitle">Prochaines 24h</div>
                          <div className="muted small">Vérifications planifiées bientôt.</div>
                        </div>
                        <span className="pill info">{subsByCheck.soon.length}</span>
                      </div>
                      <div style={{ marginTop: 10, display: 'flex', flexDirection: 'column', gap: 10 }}>
                        {subsByCheck.soon.slice(0, 10).map((s) => (
                          <div key={s.id} className="airItem">
                            <div style={{ fontWeight: 700 }}>{s.label}</div>
                            <div className="muted small">{fmtRelative(s.nextCheckAt)} • {fmtWhen(s.nextCheckAt)}</div>
                            <div className="muted small" style={{ marginTop: 6 }}>dispo: {s.lastAvailableEpisode || 0} • téléchargé: {s.lastDownloadedEpisode || 0}</div>
                          </div>
                        ))}
                        {subsByCheck.soon.length === 0 && <div className="muted small">Aucun check planifié (24h).</div>}
                      </div>
                    </div>
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

                      <label className="small muted" style={{ marginLeft: 10 }}>Langue</label>
                      <select className="select" value={resolveLang} onChange={(e) => setResolveLang(e.target.value)}>
                        <option value="vostfr">vostfr</option>
                        <option value="vf">vf</option>
                      </select>

                      <label className="small muted">Saison</label>
                      <input className="select" type="number" min={1} value={resolveSeason} onChange={(e) => setResolveSeason(Number(e.target.value) || 1)} />

                      <label className="small muted">max</label>
                      <input className="select" type="number" min={1} max={10} value={resolveMaxCandidates} onChange={(e) => setResolveMaxCandidates(Number(e.target.value) || 5)} />
                    </div>
                    <div className="muted small">Pour importer ta watchlist: onglet “Intégrations”.</div>
                  </div>

                  <div className="dayGrid">
                    {airingByDay.map(([day, list]) => (
                      <div key={day} className="dayCol">
                        <div className="dayHead">
                          <div style={{ fontWeight: 700 }}>{day}</div>
                          <span className="pill info">{list.length}</span>
                        </div>
                        <div className="dayBody">
                          {list
                            .slice()
                            .sort((a, b) => a.airingAt - b.airingAt)
                            .slice(0, 12)
                            .map((e) => (
                              <div key={e.id} className="airItem">
                                <div style={{ fontWeight: 700 }}>{titleOf(e)}</div>
                                <div className="airMeta">
                                  <span className="pill">Ep {e.episode}</span>
                                  <span className="pill">{fmtUnix(e.airingAt)}</span>
                                </div>
                                <div className="row" style={{ justifyContent: 'flex-end', marginTop: 10 }}>
                                  <button
                                    className="btn sm"
                                    disabled={busy}
                                    onClick={() =>
                                      wrapAction('Recherche de l’URL Anime‑Sama…', async () => {
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
                                    Trouver sur Anime‑Sama
                                  </button>
                                </div>

                                {resolved?.entryId === e.id && (
                                  <div style={{ marginTop: 10 }}>
                                    {!resolved.candidates.length ? (
                                      <div className="muted small">Aucun candidat trouvé.</div>
                                    ) : (
                                      <div style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
                                        {resolved.candidates.slice(0, 3).map((c) => (
                                          <div key={c.slug + '|' + c.matchedTitle} className="cardItem" style={{ padding: 10 }}>
                                            <div className="cardRow">
                                              <div>
                                                <div style={{ fontWeight: 700 }}>{c.matchedTitle || c.slug}</div>
                                                <div className="muted small" style={{ wordBreak: 'break-all' }}>{c.catalogueUrl}</div>
                                              </div>
                                              <div className="row" style={{ justifyContent: 'flex-end' }}>
                                                <span className="pill info">{Math.round((c.score || 0) * 100)}%</span>

                                                {(() => {
                                                  const st = scanBySlug[c.slug];
                                                  if (!st || st.loading) return null;
                                                  if (st.error) return <span className="pill cancel">Erreur scan</span>;
                                                  if (!st.options.length) return null;
                                                  return (
                                                    <select
                                                      className="select"
                                                      value={st.selectedBaseUrl || st.options[0].baseUrl}
                                                      onChange={(e) => {
                                                        const v = e.target.value;
                                                        setScanBySlug((m) => ({
                                                          ...m,
                                                          [c.slug]: { ...m[c.slug], selectedBaseUrl: v },
                                                        }));
                                                      }}
                                                    >
                                                      {st.options.map((o) => (
                                                        <option key={o.baseUrl} value={o.baseUrl}>
                                                          S{o.season} {o.lang.toUpperCase()} (max ep {o.maxAvailableEpisode})
                                                        </option>
                                                      ))}
                                                    </select>
                                                  );
                                                })()}

                                                <button
                                                  className="btn sm primary"
                                                  disabled={busy || !(scanBySlug[c.slug]?.selectedBaseUrl || '')}
                                                  onClick={() =>
                                                    wrapAction('Création de l’abonnement…', async () => {
                                                      const st = scanBySlug[c.slug];
                                                      const baseUrl = st?.selectedBaseUrl;
                                                      if (!baseUrl) throw new Error('Scan en cours: choisis une saison quand disponible.');
                                                      await apiCreateSubscription({ baseUrl });
                                                      setResolved(null);
                                                      await refreshSubs();
                                                      notify('ok', 'Abonnement ajouté.');
                                                      setTab('abonnements');
                                                    })
                                                  }
                                                >
                                                  Ajouter
                                                </button>
                                              </div>
                                            </div>
                                          </div>
                                        ))}
                                        <button className="btn sm" onClick={() => setResolved(null)}>Fermer</button>
                                      </div>
                                    )}
                                  </div>
                                )}
                              </div>
                            ))}
                          {list.length > 12 && <div className="muted small">+ {list.length - 12} autres…</div>}
                        </div>
                      </div>
                    ))}
                  </div>
                  {!airingByDay.length && <div className="muted" style={{ marginTop: 12 }}>Aucune donnée.</div>}
                </>
              )}
            </>
          )}

        </div>

        <div className="card">
          {tab === 'jobs' && (
            <>
              <div className="row" style={{ justifyContent: 'space-between' }}>
                <div>
                  <div className="sectionTitle">Téléchargements</div>
                  <div className="muted small sectionHint">Progression, erreurs, détails des téléchargements.</div>
                </div>
                <div className="row">
                  <button className="btn sm" disabled={busy} onClick={() => wrapAction('Rafraîchissement des téléchargements…', refreshJobs)}>
                    Rafraîchir
                  </button>
                  <button className="btn sm" onClick={() => setShowEvents((v) => !v)}>
                    {showEvents ? 'Masquer logs temps réel' : 'Logs temps réel (debug)'}
                  </button>
                </div>
              </div>

              <div className="row" style={{ marginTop: 12, justifyContent: 'space-between' }}>
                <div className="row">
                  {(() => {
                    const counts = { running: 0, queue: 0, done: 0, failed: 0 } as any;
                    for (const j of jobs) counts[jobBucket(j.state)]++;
                    const btn = (k: typeof jobsView, label: string, pillCls?: string) => (
                      <button className={`btn sm navbtn ${jobsView === k ? 'active' : ''}`} onClick={() => setJobsView(k)}>
                        {label}{' '}
                        {pillCls ? <span className={`pill ${pillCls}`}>{k === 'all' ? jobs.length : counts[k]}</span> : <span className="pill">{k === 'all' ? jobs.length : counts[k]}</span>}
                      </button>
                    );
                    return (
                      <>
                        {btn('running', 'En cours', 'info')}
                        {btn('queue', 'En attente')}
                        {btn('failed', 'Erreurs', 'warn')}
                        {btn('done', 'Terminés', 'ok')}
                        {btn('all', 'Tout')}
                      </>
                    );
                  })()}
                </div>
                <div className="muted small">Clique une ligne pour “Détails”.</div>
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
                    .filter((j) => (jobsView === 'all' ? true : jobBucket(j.state) === jobsView))
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
                                    wrapAction('Annulation du téléchargement…', async () => {
                                      await apiCancelJob(j.id);
                                      await refreshJobs();
                                      notify('ok', 'Annulé.');
                                    })
                                  }
                                >
                                  Annuler
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
                    <div className="muted small">Logs temps réel (debug)</div>
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
                  <div className="sectionTitle">Réglages</div>
                  <div className="muted small sectionHint">Téléchargements, destination, et intégrations (AniList / Jellyfin / Plex).</div>
                </div>
                <div className="row">
                  <button className="btn sm" disabled={busy} onClick={() => wrapAction('Chargement des réglages…', refreshSettings)}>
                    Recharger
                  </button>
                  <button
                    className="btn sm primary"
                    disabled={busy || !settingsDraft}
                    onClick={() =>
                      wrapAction('Sauvegarde des réglages…', async () => {
                        const saved = await apiPutSettings(settingsDraft!);
                        setSettings(saved);
                        setSettingsDraft(saved);
                        notify('ok', 'Réglages enregistrés.');
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
                  <div className="card soft" style={{ marginTop: 12 }}>
                    <div className="cardTitle">Téléchargements</div>
                    <div className="row" style={{ marginTop: 10 }}>
                      <label className="small muted">Téléchargements en parallèle</label>
                      <input
                        className="select"
                        type="number"
                        min={1}
                        value={settingsDraft.maxConcurrentDownloads}
                        onChange={(e) => setSettingsDraft((p) => ({ ...(p as Settings), maxConcurrentDownloads: Number(e.target.value) || 1 }))}
                      />
                      <label className="small muted">Workers (traitements)</label>
                      <input
                        className="select"
                        type="number"
                        min={1}
                        value={settingsDraft.maxWorkers}
                        onChange={(e) => setSettingsDraft((p) => ({ ...(p as Settings), maxWorkers: Number(e.target.value) || 1 }))}
                      />
                    </div>
                    <div className="muted small" style={{ marginTop: 8 }}>
                      Astuce: si ta connexion est instable, baisse “Téléchargements en parallèle”.
                    </div>
                  </div>

                  <div className="card soft" style={{ marginTop: 12 }}>
                    <div className="cardTitle">Destination</div>
                    <div className="muted small cardHint">En Docker, la destination doit pointer vers un volume monté.</div>
                    <div className="row" style={{ marginTop: 10 }}>
                      <input className="input" value={settingsDraft.destination} onChange={(e) => setSettingsDraft((p) => ({ ...(p as Settings), destination: e.target.value }))} />
                    </div>
                  </div>

                  <div className="card soft" style={{ marginTop: 12 }}>
                    <div className="cardTitle">Intégrations</div>
                    <div className="muted small cardHint">Optionnel: AniList pour importer ta watchlist, Jellyfin/Plex pour bibliothèques.</div>

                    <div className="row" style={{ marginTop: 10 }}>
                      <label className="small muted">AniList token</label>
                      <input
                        className="input"
                        placeholder="(optionnel)"
                        value={settingsDraft.anilistToken || ''}
                        onChange={(e) => setSettingsDraft((p) => ({ ...(p as Settings), anilistToken: e.target.value }))}
                      />
                      <button className="btn sm" disabled={busy || !(settingsDraft.anilistToken || '').trim()} onClick={() => wrapAction('Vérification AniList…', refreshViewer)}>
                        Tester
                      </button>
                    </div>
                    {viewer && <div className="muted small" style={{ marginTop: 8 }}>Connecté AniList: {viewer.name} (#{viewer.id})</div>}

                    <details style={{ marginTop: 10 }}>
                      <summary className="btn sm" style={{ display: 'inline-flex' }}>Importer watchlist (avancé)</summary>
                      <div className="muted small" style={{ marginTop: 8 }}>Crée des abonnements automatiquement depuis ta watchlist AniList.</div>
                      <div className="row" style={{ marginTop: 10, justifyContent: 'space-between' }}>
                        <div className="muted small">Statuts: {importReq.statuses.join(', ')}</div>
                        <button
                          className="btn sm primary"
                          disabled={busy}
                          onClick={() =>
                            wrapAction('Import AniList…', async () => {
                              const res = await apiAniListImportAuto(importReq);
                              setImportRes(res);
                              await refreshSubs();
                              notify('ok', `Import terminé: created ${(res.created || []).length}`);
                            })
                          }
                        >
                          Lancer l’import
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
                        <input className="select" type="number" step={0.05} min={0} max={1} value={importReq.minScore} onChange={(e) => setImportReq((p) => ({ ...p, minScore: Number(e.target.value) || 0 }))} />
                      </div>

                      {importRes && (
                        <div style={{ marginTop: 12 }}>
                          <div className="row" style={{ justifyContent: 'space-between' }}>
                            <div className="sectionTitle">Résultat</div>
                            <button className="btn sm" onClick={() => setImportRes(null)}>Clear</button>
                          </div>
                          <div className="row" style={{ marginTop: 8 }}>
                            <span className="badge ok">created {importRes.created?.length || 0}</span>
                            <span className="badge">skipped {importRes.skipped?.length || 0}</span>
                            <span className="badge fail">errors {importRes.errors?.length || 0}</span>
                          </div>
                        </div>
                      )}
                    </details>

                    <details style={{ marginTop: 10 }}>
                      <summary className="btn sm" style={{ display: 'inline-flex' }}>Jellyfin</summary>
                      <div className="row" style={{ marginTop: 10 }}>
                        <label className="small muted">URL</label>
                        <input className="input" placeholder="http://jellyfin:8096" value={settingsDraft.jellyfinUrl || ''} onChange={(e) => setSettingsDraft((p) => ({ ...(p as Settings), jellyfinUrl: e.target.value }))} />
                      </div>
                      <div className="row" style={{ marginTop: 10 }}>
                        <label className="small muted">API key</label>
                        <input className="input" placeholder="(optionnel)" value={settingsDraft.jellyfinApiKey || ''} onChange={(e) => setSettingsDraft((p) => ({ ...(p as Settings), jellyfinApiKey: e.target.value }))} />
                      </div>
                    </details>

                    <details style={{ marginTop: 10 }}>
                      <summary className="btn sm" style={{ display: 'inline-flex' }}>Plex</summary>
                      <div className="row" style={{ marginTop: 10 }}>
                        <label className="small muted">URL</label>
                        <input className="input" placeholder="http://plex:32400" value={settingsDraft.plexUrl || ''} onChange={(e) => setSettingsDraft((p) => ({ ...(p as Settings), plexUrl: e.target.value }))} />
                      </div>
                      <div className="row" style={{ marginTop: 10 }}>
                        <label className="small muted">Token</label>
                        <input className="input" placeholder="(optionnel)" value={settingsDraft.plexToken || ''} onChange={(e) => setSettingsDraft((p) => ({ ...(p as Settings), plexToken: e.target.value }))} />
                      </div>
                      <div className="row" style={{ marginTop: 10 }}>
                        <label className="small muted">Section ID</label>
                        <input className="input" placeholder="(optionnel)" value={settingsDraft.plexSectionId || ''} onChange={(e) => setSettingsDraft((p) => ({ ...(p as Settings), plexSectionId: e.target.value }))} />
                      </div>
                    </details>
                  </div>

                  <details style={{ marginTop: 12 }}>
                    <summary className="btn sm" style={{ display: 'inline-flex' }}>Options avancées</summary>
                    <div className="row" style={{ marginTop: 10 }}>
                      <label className="small muted">Nom des fichiers</label>
                      <select className="select" value={settingsDraft.outputNamingMode} onChange={(e) => setSettingsDraft((p) => ({ ...(p as Settings), outputNamingMode: e.target.value as Settings['outputNamingMode'] }))}>
                        <option value="legacy">Simple</option>
                        <option value="media-server">Media‑server</option>
                      </select>
                      <label className="small" style={{ display: 'flex', gap: 6, alignItems: 'center' }}>
                        <input type="checkbox" checked={settingsDraft.separateLang} onChange={(e) => setSettingsDraft((p) => ({ ...(p as Settings), separateLang: e.target.checked }))} />
                        Séparer VF/VOSTFR
                      </label>
                    </div>
                  </details>

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
              <div className="sectionTitle">Raccourcis</div>
              <div className="muted small sectionHint">Va dans “Téléchargements” pour suivre la progression.</div>
              <div className="row" style={{ marginTop: 10 }}>
                <button className="btn" onClick={() => setTab('jobs')}>Ouvrir Téléchargements</button>
                <button className="btn" disabled={busy} onClick={() => wrapAction('Rafraîchissement des téléchargements…', refreshJobs)}>Rafraîchir</button>
              </div>
            </>
          )}
        </div>
      </div>
    </div>
  );
}
