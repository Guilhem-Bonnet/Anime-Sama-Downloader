const API = '/api/v1';

async function parseError(r: Response): Promise<string> {
  try {
    const j = await r.json();
    if (j && typeof j.error === 'string') return j.error;
    return JSON.stringify(j);
  } catch {
    return await r.text();
  }
}

async function fetchJson<T>(path: string, init?: RequestInit): Promise<T> {
  const r = await fetch(`${API}${path}`, init);
  if (!r.ok) throw new Error(await parseError(r));
  return (await r.json()) as T;
}

export type Settings = {
  destination: string;
  outputNamingMode: 'legacy' | 'media-server';
  separateLang: boolean;
  maxWorkers: number;
  maxConcurrentDownloads: number;
  jellyfinUrl?: string;
  jellyfinApiKey?: string;
  plexUrl?: string;
  plexToken?: string;
  plexSectionId?: string;
  anilistToken?: string;
};

export type Subscription = {
  id: string;
  baseUrl: string;
  label: string;
  player: string;
  lastScheduledEpisode: number;
  lastDownloadedEpisode: number;
  lastAvailableEpisode: number;
  nextCheckAt: string;
  lastCheckedAt: string;
  createdAt: string;
  updatedAt: string;
};

export type SyncResult = {
  subscription: Subscription;
  selectedPlayer: string;
  maxAvailableEpisode: number;
  enqueuedEpisodes: number[];
  enqueuedJobIDs: string[];
  message: string;
};

export type JobState = 'queued' | 'running' | 'muxing' | 'completed' | 'failed' | 'canceled';

export type Job = {
  id: string;
  type: string;
  state: JobState;
  progress: number;
  createdAt: string;
  updatedAt: string;
  params?: any;
  result?: any;
  errorCode?: string;
  error?: string;
};

export type AniListViewer = { id: number; name: string };

export type AniListAiringScheduleEntry = {
  id: number;
  airingAt: number;
  episode: number;
  media: {
    id: number;
    title: { romaji: string; english: string; native: string };
  };
};

export type AnimeSamaResolveRequest = {
  titles: string[];
  season?: number;
  lang?: string;
  maxCandidates?: number;
};

export type AnimeSamaResolvedCandidate = {
  catalogueUrl: string;
  baseUrl: string;
  slug: string;
  matchedTitle: string;
  score: number;
};

export type AnimeSamaResolveResponse = {
  candidates: AnimeSamaResolvedCandidate[];
};

export type AniListImportAutoRequest = {
  statuses: string[];
  season: number;
  lang: string;
  maxCandidates: number;
  minScore: number;
};

export type AniListImportAutoResponse = {
  created: Subscription[];
  skipped: Array<{ anilistMediaId: number; title: string; reason: string; baseUrl?: string; topScore?: number }>;
  errors: Array<{ baseUrl: string; error: string }>;
};

export async function apiGetSettings(): Promise<Settings> {
  return await fetchJson<Settings>('/settings');
}

export async function apiPutSettings(s: Settings): Promise<Settings> {
  return await fetchJson<Settings>('/settings', {
    method: 'PUT',
    headers: { 'content-type': 'application/json' },
    body: JSON.stringify(s),
  });
}

export async function apiListSubscriptions(limit = 200): Promise<Subscription[]> {
  return await fetchJson<Subscription[]>(`/subscriptions?limit=${encodeURIComponent(String(limit))}`);
}

export async function apiCreateSubscription(params: { baseUrl: string; label: string; player?: string }): Promise<Subscription> {
  return await fetchJson<Subscription>('/subscriptions/', {
    method: 'POST',
    headers: { 'content-type': 'application/json' },
    body: JSON.stringify(params),
  });
}

export async function apiDeleteSubscription(id: string): Promise<void> {
  await fetchJson<void>(`/subscriptions/${encodeURIComponent(id)}`, { method: 'DELETE' });
}

export async function apiSyncSubscription(id: string, enqueue = true): Promise<SyncResult> {
  const qs = enqueue ? '' : '?enqueue=false';
  return await fetchJson<SyncResult>(`/subscriptions/${encodeURIComponent(id)}/sync${qs}`, { method: 'POST' });
}

export async function apiSyncAll(params: { enqueue?: boolean; dueOnly?: boolean; limit?: number } = {}): Promise<{ results: SyncResult[]; errors: Array<{ id: string; error: string }> }> {
  const sp = new URLSearchParams();
  if (params.enqueue === false) sp.set('enqueue', 'false');
  if (params.dueOnly) sp.set('dueOnly', 'true');
  if (typeof params.limit === 'number' && params.limit > 0) sp.set('limit', String(params.limit));
  const qs = sp.toString() ? `?${sp.toString()}` : '';
  return await fetchJson(`/subscriptions/sync-all${qs}`, { method: 'POST' });
}

export async function apiListJobs(limit = 200): Promise<Job[]> {
  return await fetchJson<Job[]>(`/jobs?limit=${encodeURIComponent(String(limit))}`);
}

export async function apiCancelJob(id: string): Promise<Job> {
  return await fetchJson<Job>(`/jobs/${encodeURIComponent(id)}/cancel`, { method: 'POST' });
}

export async function apiAniListViewer(): Promise<AniListViewer> {
  return await fetchJson<AniListViewer>('/anilist/viewer');
}

export async function apiAniListAiring(days = 7, limit = 50): Promise<AniListAiringScheduleEntry[]> {
  const qs = new URLSearchParams({ days: String(days), limit: String(limit) });
  return await fetchJson<AniListAiringScheduleEntry[]>(`/anilist/airing?${qs.toString()}`);
}

export async function apiAnimeSamaResolve(req: AnimeSamaResolveRequest): Promise<AnimeSamaResolveResponse> {
  return await fetchJson<AnimeSamaResolveResponse>('/animesama/resolve', {
    method: 'POST',
    headers: { 'content-type': 'application/json' },
    body: JSON.stringify(req),
  });
}

export async function apiAniListImportAuto(req: AniListImportAutoRequest): Promise<AniListImportAutoResponse> {
  return await fetchJson<AniListImportAutoResponse>('/import/anilist/auto', {
    method: 'POST',
    headers: { 'content-type': 'application/json' },
    body: JSON.stringify(req),
  });
}
