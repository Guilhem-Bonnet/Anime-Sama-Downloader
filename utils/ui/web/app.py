from __future__ import annotations

import asyncio
import json
import os
from pathlib import Path
import threading
import time
from typing import Any, AsyncIterator

from fastapi import FastAPI, HTTPException
from fastapi.responses import HTMLResponse, StreamingResponse
from fastapi.staticfiles import StaticFiles

from utils.download_manager import DownloadJob, DownloadManager
from utils.config import get_default_download_path, get_max_concurrent_downloads
from utils.fetch import fetch_episodes, fetch_video_source, rank_players
from utils.output_paths import build_episode_output_path
from utils.search import resolve_anime_sama_base_url
from utils.selection import parse_episode_expr, parse_tracker_selection
from utils.downloaders.downloader import download_video


def create_app() -> FastAPI:
    app = FastAPI(title="Anime-Sama Downloader", version="0.1")

    # Event bus for job updates (SSE)
    subscribers: set[asyncio.Queue[str]] = set()
    subscribers_lock = threading.Lock()

    # Simple in-memory cache for seasons scans (per base_url+lang)
    seasons_cache: dict[tuple[str, str], tuple[float, list[int]]] = {}
    seasons_cache_lock = threading.Lock()
    seasons_cache_ttl_s = 10 * 60

    def publish(event: dict[str, Any]) -> None:
        data = json.dumps(event, ensure_ascii=False)
        with subscribers_lock:
            for q in list(subscribers):
                try:
                    q.put_nowait(data)
                except Exception:
                    pass

    mgr = DownloadManager(
        max_parallel=get_max_concurrent_downloads(default=10),
        executor_name="web-dl",
        on_event=lambda job, evt: publish({
            "type": "job",
            "event": evt,
            "job": {
                "job_id": job.job_id,
                "label": job.label,
                "status": job.status,
                "result_path": job.result_path,
                "error": job.error,
                "created_at": job.created_at,
                "started_at": job.started_at,
                "finished_at": job.finished_at,
            },
            "ts": time.time(),
        }),
    )

    def _minimal_index_html() -> str:
        return """<!doctype html>
<html lang=\"fr\">
<head>
  <meta charset=\"utf-8\" />
  <meta name=\"viewport\" content=\"width=device-width, initial-scale=1\" />
  <title>Anime-Sama Downloader</title>
  <style>
    body{font-family:system-ui,Segoe UI,Roboto,Arial,sans-serif;max-width:980px;margin:24px auto;padding:0 16px;}
    .row{display:flex;gap:12px;flex-wrap:wrap;align-items:center}
    input,select,button{padding:10px 12px;font-size:14px}
    button{cursor:pointer}
    pre{background:#0b1020;color:#e8e8e8;padding:12px;border-radius:8px;overflow:auto}
    .card{border:1px solid #ddd;border-radius:10px;padding:12px;margin:12px 0}
    .muted{color:#666}
        .hint{background:#fff7e6;border:1px solid #f2d3a1;color:#6b4b16;border-radius:10px;padding:10px 12px;margin:12px 0}
  </style>
</head>
<body>
  <h1>Anime-Sama Downloader (Web)</h1>
    <div class=\"hint\">
        <strong>Astuce dev</strong> : pour la “vraie app” (SPA), ouvre <a href=\"http://127.0.0.1:5173/\">http://127.0.0.1:5173/</a>.<br/>
        Cette page est un fallback (utile si tu n’as pas build le frontend).
    </div>
  <div class=\"card\">
    <div class=\"row\">
      <input id=\"q\" placeholder=\"Nom de l'anime\" style=\"flex:1;min-width:240px\" />
      <select id=\"lang\">
        <option value=\"vostfr\">vostfr</option>
        <option value=\"vf\">vf</option>
        <option value=\"vo\">vo</option>
      </select>
      <button id=\"search\">Chercher</button>
    </div>
    <div class=\"row\" style=\"margin-top:10px\">
      <span class=\"muted\">URL:</span>
      <span id=\"url\" class=\"muted\">—</span>
    </div>
    <div class=\"row\" style=\"margin-top:10px\">
      <select id=\"season\"></select>
      <input id=\"sel\" placeholder=\"ALL | 1-6 | S1E1-6 | ALLSEASONS\" style=\"flex:1;min-width:240px\" />
      <button id=\"enqueue\">Ajouter à la file</button>
    </div>
    <div class=\"row\" style=\"margin-top:10px\">
      <input id=\"dest\" placeholder=\"Dossier destination\" style=\"flex:1;min-width:240px\" />
      <button id=\"cancel\">Annuler tout</button>
      <button id=\"clear\">Vider file</button>
    </div>
  </div>

  <div class=\"card\">
    <h2>Téléchargements</h2>
    <div class=\"muted\" id=\"qs\">File: —</div>
    <pre id=\"log\"></pre>
  </div>

<script>
let baseUrl = null;
const log = (s)=>{ const el=document.getElementById('log'); el.textContent += s+"\n"; el.scrollTop=el.scrollHeight; };
const setQs = (s)=>document.getElementById('qs').textContent=s;

async function refreshJobs(){
  const r = await fetch('/api/jobs');
  const j = await r.json();
  setQs(`File: ${j.pending} en attente | ${j.running} en cours | total ${j.total}`);
}

document.getElementById('search').onclick = async ()=>{
  const q = document.getElementById('q').value.trim();
  const lang = document.getElementById('lang').value;
  if(!q) return;
  log(`[INFO] Recherche: ${q}`);
  const r = await fetch('/api/search', {method:'POST', headers:{'content-type':'application/json'}, body:JSON.stringify({query:q, lang})});
  const j = await r.json();
  if(!j.base_url){ log('[ERROR] Aucun résultat'); return; }
  baseUrl = j.base_url;
  document.getElementById('url').textContent = baseUrl;

  const s = await fetch('/api/seasons', {method:'POST', headers:{'content-type':'application/json'}, body:JSON.stringify({base_url:baseUrl, lang})});
  const sj = await s.json();
  const sel = document.getElementById('season');
  sel.innerHTML = '';
  for(const season of (sj.seasons||[])){
    const opt=document.createElement('option');
    opt.value=String(season); opt.textContent='S'+season; sel.appendChild(opt);
  }
  if(sel.options.length) sel.value = sel.options[0].value;
  log(`[INFO] Saisons: ${(sj.seasons||[]).map(x=>'S'+x).join(', ') || '—'}`);
};

document.getElementById('enqueue').onclick = async ()=>{
  if(!baseUrl){ log('[ERROR] Fais une recherche d\'abord'); return; }
  const lang = document.getElementById('lang').value;
  const season = Number(document.getElementById('season').value || '1');
  const selection = document.getElementById('sel').value;
  const dest = document.getElementById('dest').value;
  const r = await fetch('/api/enqueue', {method:'POST', headers:{'content-type':'application/json'}, body:JSON.stringify({base_url:baseUrl, lang, season, selection, dest_root:dest})});
  const j = await r.json();
  if(j.error){ log('[ERROR] '+j.error); return; }
  log(`[INFO] Ajouté: ${j.enqueued} épisode(s)`);
  await refreshJobs();
};

document.getElementById('cancel').onclick = async ()=>{
  await fetch('/api/cancel_all', {method:'POST'});
  log('[INFO] Annuler tout demandé');
};

document.getElementById('clear').onclick = async ()=>{
  const r = await fetch('/api/clear_pending', {method:'POST'});
  const j = await r.json();
  log(`[INFO] File vidée: ${j.cleared}`);
};

const es = new EventSource('/api/jobs/stream');
es.onmessage = (ev)=>{
  try{ const msg = JSON.parse(ev.data);
    if(msg.type==='job'){
      log(`[${(msg.event||'').toUpperCase()}] ${msg.job.label}` + (msg.job.error?` | ${msg.job.error}`:'') + (msg.job.result_path?` | ${msg.job.result_path}`:''));
      refreshJobs();
    }
  }catch(e){ }
};
refreshJobs();
</script>
</body>
</html>"""

    def _dist_dir() -> Path:
        # Repo layout: <root>/utils/ui/web/app.py -> <root>/webapp/dist
        return Path(__file__).resolve().parents[3] / "webapp" / "dist"

    @app.post("/api/search")
    def api_search(payload: dict[str, Any]) -> dict[str, Any]:
        query = str(payload.get("query") or "").strip()
        if not query:
            raise HTTPException(400, "query required")
        base = resolve_anime_sama_base_url(query, provider="anilist")
        return {"base_url": base}

    @app.get("/api/health")
    def api_health() -> dict[str, Any]:
        return {"status": "ok"}

    @app.get("/api/defaults")
    def api_defaults() -> dict[str, Any]:
        return {
            "download_root": get_default_download_path(),
            "max_concurrent_downloads": get_max_concurrent_downloads(default=10),
        }

    @app.post("/api/season_info")
    def api_season_info(payload: dict[str, Any]) -> dict[str, Any]:
        base_url = str(payload.get("base_url") or "").strip().rstrip("/")
        lang = str(payload.get("lang") or "vostfr").strip().lower()
        try:
            season = int(payload.get("season") or 1)
        except Exception:
            raise HTTPException(400, "season must be an integer")

        if not base_url:
            raise HTTPException(400, "base_url required")
        if lang not in {"vostfr", "vf", "vo"}:
            raise HTTPException(400, "lang must be one of: vostfr, vf, vo")
        if season <= 0:
            raise HTTPException(400, "season must be >= 1")

        full_url = f"{base_url}/saison{season}/{lang}/"
        eps = fetch_episodes(full_url, quiet=True)
        if not eps:
            return {"season": season, "max_episodes": 0, "available": []}

        max_eps = max((len(v) for v in eps.values() if v), default=0)
        available: list[int] = []
        for n in range(1, max_eps + 1):
            idx0 = n - 1
            if any(((idx0 < len(urls)) and bool(urls[idx0])) for urls in eps.values() if urls):
                available.append(n)
        return {"season": season, "max_episodes": max_eps, "available": available}

    @app.post("/api/seasons")
    def api_seasons(payload: dict[str, Any]) -> dict[str, Any]:
        base_url = str(payload.get("base_url") or "").strip().rstrip("/")
        lang = str(payload.get("lang") or "vostfr")
        if not base_url:
            raise HTTPException(400, "base_url required")

        cache_key = (base_url, lang)
        now = time.time()
        with seasons_cache_lock:
            cached = seasons_cache.get(cache_key)
            if cached and (now - cached[0]) <= seasons_cache_ttl_s:
                return {"seasons": cached[1], "cached": True}

        seasons: list[int] = []
        consecutive_misses = 0
        for season in range(1, 51):
            full_url = f"{base_url}/saison{season}/{lang}/"
            eps = fetch_episodes(full_url, quiet=True, timeout=8)
            if eps:
                seasons.append(season)
                consecutive_misses = 0
            else:
                consecutive_misses += 1
                if seasons and consecutive_misses >= 3:
                    break
                if not seasons and consecutive_misses >= 5:
                    break

        with seasons_cache_lock:
            seasons_cache[cache_key] = (now, seasons)

        return {"seasons": seasons, "cached": False}

    def _make_episode_job(
        base_catalogue_url: str,
        season: int,
        ep_num: int,
        episodes: dict[str, list[str]],
        lang: str,
        dest_root: str,
    ) -> DownloadJob:
        idx0 = ep_num - 1
        ranked = rank_players(episodes) or list(episodes.keys())
        slug = base_catalogue_url.rstrip("/").split("/")[-1] or "anime"

        dest_root = (dest_root or "").strip()
        dest_dir, save_path = build_episode_output_path(dest_root, slug, season, lang, ep_num, ext="mp4")
        os.makedirs(dest_dir, exist_ok=True)

        # With 10 jobs parallel, keep per-job fanout low.
        ts_workers = 2
        mp4_workers = 2

        label = f"{slug} S{season}E{ep_num}"

        # Create the job first so we can reference job_id in callbacks.
        job = DownloadJob(label=label, run=lambda _ev: None)

        def _progress(p: dict[str, Any]) -> None:
            mgr.update_job_progress(
                job.job_id,
                percent=p.get("percent"),
                downloaded=p.get("downloaded"),
                total=p.get("total"),
                speed_bps=p.get("speed_bps"),
                eta_seconds=p.get("eta_seconds"),
                stage=p.get("stage"),
                message=p.get("message"),
            )
            publish({
                "type": "progress",
                "job_id": job.job_id,
                "progress": {
                    "percent": p.get("percent"),
                    "downloaded": p.get("downloaded"),
                    "total": p.get("total"),
                    "speed_bps": p.get("speed_bps"),
                    "eta_seconds": p.get("eta_seconds"),
                    "stage": p.get("stage"),
                    "message": p.get("message"),
                },
                "ts": time.time(),
            })

        def _runner(cancel_event: threading.Event) -> str | None:
            for player in ranked:
                if cancel_event.is_set():
                    return None

                urls = episodes.get(player) or []
                if idx0 >= len(urls):
                    continue
                page_url = urls[idx0]
                if not page_url:
                    continue

                try:
                    src = fetch_video_source(page_url)
                except Exception as e:
                    publish({"type": "log", "level": "error", "msg": f"{label}: erreur source {player}: {e}", "ts": time.time()})
                    continue

                if not src:
                    continue

                def _cb(message: str, status_type: str = "info") -> None:
                    publish({"type": "log", "level": status_type, "msg": f"{label}: {message}", "ts": time.time()})

                ok, out = download_video(
                    src,
                    save_path,
                    use_ts_threading=False,
                    url=page_url,
                    automatic_mp4=False,
                    log_callback=_cb,
                    progress_callback=_progress,
                    use_tqdm=False,
                    use_mp4_threading=True,
                    mp4_workers=mp4_workers,
                    ts_workers=ts_workers,
                    cancel_event=cancel_event,
                )
                if ok and out:
                    return out

            return None

        job.run = _runner
        return job

    @app.post("/api/enqueue")
    def api_enqueue(payload: dict[str, Any]) -> dict[str, Any]:
        base_url = str(payload.get("base_url") or "").strip().rstrip("/")
        lang = str(payload.get("lang") or "vostfr").strip().lower()
        if lang not in {"vostfr", "vf", "vo"}:
            raise HTTPException(400, "lang must be one of: vostfr, vf, vo")

        try:
            season_default = int(payload.get("season") or 1)
        except Exception:
            raise HTTPException(400, "season must be an integer")

        if season_default <= 0:
            raise HTTPException(400, "season must be >= 1")
        selection = str(payload.get("selection") or "").strip()
        dest_root = str(payload.get("dest_root") or "").strip()

        if not base_url:
            raise HTTPException(400, "base_url required")

        selection_parsed = parse_tracker_selection(selection, default_season=season_default)
        all_seasons_requested = any(s is None for s, _ in selection_parsed)
        if all_seasons_requested:
            seasons_to_download = list(range(1, 21))
        else:
            seasons_to_download = sorted({int(s) for s, _ in selection_parsed if isinstance(s, int)})

        per_season_specs: dict[int, list[str]] = {}
        if all_seasons_requested:
            for s in seasons_to_download:
                per_season_specs.setdefault(s, []).append("ALL")
        else:
            for s, spec in selection_parsed:
                if s is None:
                    continue
                per_season_specs.setdefault(int(s), []).append(spec)

        enqueued = 0
        any_season_available = False
        any_episode_available = False
        for season in seasons_to_download:
            full_url = f"{base_url}/saison{season}/{lang}/"
            eps = fetch_episodes(full_url, quiet=True)
            if not eps:
                continue
            any_season_available = True
            max_eps = max((len(v) for v in eps.values() if v), default=0)
            if max_eps <= 0:
                continue
            any_episode_available = True

            requested: set[int] = set()
            for spec in per_season_specs.get(season, ["1"]):
                if str(spec).upper() == "ALL":
                    requested.update(range(1, max_eps + 1))
                else:
                    requested.update(parse_episode_expr(str(spec), max_episodes=max_eps))

            # keep only available
            requested = {n for n in requested if any(((n - 1) < len(urls) and urls[n - 1]) for urls in eps.values())}
            for ep in sorted(requested):
                mgr.enqueue(_make_episode_job(base_url, season, ep, eps, lang, dest_root))
                enqueued += 1

        if enqueued <= 0:
            if not any_season_available:
                return {"enqueued": 0, "error": "Impossible de récupérer la liste des épisodes (saison/lang indisponible ou blocage réseau)."}
            if not any_episode_available:
                return {"enqueued": 0, "error": "Saison détectée, mais aucun épisode exploitable n’a été trouvé."}
            return {"enqueued": 0, "error": "Aucun épisode n’a été ajouté (sélection vide ou non disponible)."}

        return {"enqueued": enqueued}

    def _job_to_dict(job: DownloadJob) -> dict[str, Any]:
        return {
            "job_id": job.job_id,
            "label": job.label,
            "status": job.status,
            "result_path": job.result_path,
            "error": job.error,
            "created_at": job.created_at,
            "started_at": job.started_at,
            "finished_at": job.finished_at,

            "progress_percent": job.progress_percent,
            "progress_downloaded": job.progress_downloaded,
            "progress_total": job.progress_total,
            "progress_speed_bps": job.progress_speed_bps,
            "progress_eta_seconds": job.progress_eta_seconds,
            "progress_stage": job.progress_stage,
            "progress_message": job.progress_message,
        }

    @app.get("/api/jobs")
    def api_jobs_compat() -> dict[str, Any]:
        # Backward compatibility for the minimal HTML UI.
        return {
            "pending": mgr.pending_count(),
            "running": mgr.running_count(),
            "total": len(mgr.list_jobs()),
        }

    @app.get("/api/jobs/list")
    def api_jobs_list() -> dict[str, Any]:
        jobs = sorted(mgr.list_jobs(), key=lambda j: j.created_at)
        return {
            "pending": mgr.pending_count(),
            "running": mgr.running_count(),
            "total": len(jobs),
            "jobs": [_job_to_dict(j) for j in jobs],
        }

    @app.post("/api/jobs/{job_id}/cancel")
    def api_cancel_job(job_id: str) -> dict[str, Any]:
        ok = mgr.cancel(str(job_id))
        if not ok:
            raise HTTPException(404, "job not found")
        return {"ok": True}

    @app.post("/api/jobs/{job_id}/retry")
    def api_retry_job(job_id: str) -> dict[str, Any]:
        new_id = mgr.retry(str(job_id))
        if not new_id:
            raise HTTPException(404, "job not found")
        return {"ok": True, "job_id": new_id}

    @app.post("/api/jobs/clear_finished")
    def api_clear_finished() -> dict[str, Any]:
        cleared = mgr.clear_finished()
        return {"cleared": cleared}

    @app.post("/api/cancel_all")
    def api_cancel_all() -> dict[str, Any]:
        mgr.cancel_all()
        return {"ok": True}

    @app.post("/api/clear_pending")
    def api_clear_pending() -> dict[str, Any]:
        cleared = mgr.clear_pending()
        return {"cleared": cleared}

    async def _sse_stream() -> StreamingResponse:
        q: asyncio.Queue[str] = asyncio.Queue(maxsize=200)
        with subscribers_lock:
            subscribers.add(q)

        async def gen() -> AsyncIterator[bytes]:
            try:
                yield b"retry: 1500\n\n"
                while True:
                    data = await q.get()
                    yield ("data: " + data + "\n\n").encode("utf-8")
            finally:
                with subscribers_lock:
                    subscribers.discard(q)

        return StreamingResponse(gen(), media_type="text/event-stream")

    @app.get("/api/jobs/stream")
    async def api_jobs_stream() -> StreamingResponse:
        # Backward compatibility endpoint.
        return await _sse_stream()

    @app.get("/api/events")
    async def api_events() -> StreamingResponse:
        # Preferred endpoint for the SPA.
        return await _sse_stream()

    dist = _dist_dir()
    if dist.is_dir() and (dist / "index.html").is_file():
        # Serve built SPA in production mode: `cd webapp && npm i && npm run build`.
        app.mount("/", StaticFiles(directory=str(dist), html=True), name="spa")
    else:
        @app.get("/", response_class=HTMLResponse)
        def index() -> str:
            # Fallback minimal UI when SPA isn't built.
            return _minimal_index_html()

    return app


def run_web(host: str = "127.0.0.1", port: int = 8000) -> int:
    # Lazy import so core can be used without uvicorn.
    import uvicorn

    try:
        from utils.config import get_web_bind

        cfg_host, cfg_port = get_web_bind()
        host = cfg_host or host
        port = int(cfg_port or port)
    except Exception:
        pass

    uvicorn.run(create_app(), host=host, port=port, log_level="info")
    return 0
