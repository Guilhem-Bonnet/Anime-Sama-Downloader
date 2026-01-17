export type JobStatus = 'PENDING' | 'RUNNING' | 'SUCCESS' | 'FAILED' | 'CANCELLED';

export type JobItem = {
  job_id: string;
  label: string;
  status: JobStatus;
  result_path?: string | null;
  error?: string | null;
  created_at?: number | null;
  started_at?: number | null;
  finished_at?: number | null;

  progress_percent?: number | null;
  progress_downloaded?: number | null;
  progress_total?: number | null;
  progress_speed_bps?: number | null;
  progress_eta_seconds?: number | null;
  progress_stage?: string | null;
  progress_message?: string | null;
};

export type JobsSnapshot = {
  pending: number;
  running: number;
  total: number;
  jobs: JobItem[];
};

export async function apiSearch(query: string): Promise<{ base_url: string | null }> {
  const r = await fetch('/api/search', {
    method: 'POST',
    headers: { 'content-type': 'application/json' },
    body: JSON.stringify({ query }),
  });
  if (!r.ok) throw new Error(await r.text());
  return await r.json();
}

export async function apiSeasons(base_url: string, lang: string): Promise<{ seasons: number[] }> {
  const r = await fetch('/api/seasons', {
    method: 'POST',
    headers: { 'content-type': 'application/json' },
    body: JSON.stringify({ base_url, lang }),
  });
  if (!r.ok) throw new Error(await r.text());
  return await r.json();
}

export async function apiDefaults(): Promise<{
  download_root: string;
  max_concurrent_downloads: number;
  is_docker?: boolean;
  allowed_dest_prefixes?: string[];
}> {
  const r = await fetch('/api/defaults');
  if (!r.ok) throw new Error(await r.text());
  return await r.json();
}

export async function apiSeasonInfo(base_url: string, lang: string, season: number): Promise<{ season: number; max_episodes: number; available: number[] }> {
  const r = await fetch('/api/season_info', {
    method: 'POST',
    headers: { 'content-type': 'application/json' },
    body: JSON.stringify({ base_url, lang, season }),
  });
  if (!r.ok) throw new Error(await r.text());
  return await r.json();
}

export async function apiEnqueue(params: {
  base_url: string;
  lang: string;
  season: number;
  selection: string;
  dest_root: string;
}): Promise<{ enqueued: number; error?: string }> {
  const r = await fetch('/api/enqueue', {
    method: 'POST',
    headers: { 'content-type': 'application/json' },
    body: JSON.stringify(params),
  });
  if (!r.ok) throw new Error(await r.text());
  return await r.json();
}

export async function apiJobs(): Promise<JobsSnapshot> {
  const r = await fetch('/api/jobs/list');
  if (!r.ok) throw new Error(await r.text());
  return await r.json();
}

export async function apiCancelAll(): Promise<void> {
  const r = await fetch('/api/cancel_all', { method: 'POST' });
  if (!r.ok) throw new Error(await r.text());
}

export async function apiClearPending(): Promise<{ cleared: number }> {
  const r = await fetch('/api/clear_pending', { method: 'POST' });
  if (!r.ok) throw new Error(await r.text());
  return await r.json();
}

export async function apiCancelJob(job_id: string): Promise<{ ok: boolean }> {
  const r = await fetch(`/api/jobs/${encodeURIComponent(job_id)}/cancel`, { method: 'POST' });
  if (!r.ok) throw new Error(await r.text());
  return await r.json();
}

export async function apiRetryJob(job_id: string): Promise<{ ok: boolean; job_id: string }> {
  const r = await fetch(`/api/jobs/${encodeURIComponent(job_id)}/retry`, { method: 'POST' });
  if (!r.ok) throw new Error(await r.text());
  return await r.json();
}

export async function apiClearFinished(): Promise<{ cleared: number }> {
  const r = await fetch('/api/jobs/clear_finished', { method: 'POST' });
  if (!r.ok) throw new Error(await r.text());
  return await r.json();
}
