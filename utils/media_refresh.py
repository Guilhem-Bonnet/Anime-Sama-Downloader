from __future__ import annotations

import threading
import time
from dataclasses import dataclass

import requests
from requests import Response

from utils.config import load_ini_config


def _env_get(key: str) -> str | None:
    import os

    v = os.environ.get(key)
    if v is None:
        return None
    v = str(v).strip()
    return v or None


def _parse_bool(v: str | None) -> bool | None:
    if v is None:
        return None
    s = str(v).strip().lower()
    if s in {"1", "true", "yes", "y", "on"}:
        return True
    if s in {"0", "false", "no", "n", "off"}:
        return False
    return None


def _cfg_get(cfg, section: str, key: str) -> str | None:
    try:
        if cfg.has_option(section, key):
            v = cfg.get(section, key)
            v = str(v).strip()
            return v or None
    except Exception:
        return None
    return None


@dataclass(frozen=True)
class JellyfinConfig:
    url: str
    api_key: str


@dataclass(frozen=True)
class PlexConfig:
    url: str
    token: str
    section_id: str


def _load_configs() -> tuple[bool, float, JellyfinConfig | None, PlexConfig | None]:
    cfg = load_ini_config()

    enabled_raw = (
        _env_get("ASD_MEDIA_REFRESH_ENABLED")
        or _env_get("ASD_MEDIA_REFRESH")
        or _cfg_get(cfg, "MEDIA_REFRESH", "enabled")
    )
    enabled = _parse_bool(enabled_raw)

    debounce_raw = (
        _env_get("ASD_MEDIA_REFRESH_DEBOUNCE_SECONDS")
        or _env_get("ASD_MEDIA_REFRESH_DEBOUNCE")
        or _cfg_get(cfg, "MEDIA_REFRESH", "debounce_seconds")
    )
    debounce_s = 20.0
    if debounce_raw:
        try:
            debounce_s = max(0.0, float(str(debounce_raw).strip()))
        except Exception:
            debounce_s = 20.0

    jellyfin_url = _env_get("ASD_JELLYFIN_URL") or _cfg_get(cfg, "JELLYFIN", "url")
    jellyfin_key = _env_get("ASD_JELLYFIN_API_KEY") or _cfg_get(cfg, "JELLYFIN", "api_key")
    jellyfin = None
    if jellyfin_url and jellyfin_key:
        jellyfin = JellyfinConfig(url=jellyfin_url.rstrip("/"), api_key=jellyfin_key)

    plex_url = _env_get("ASD_PLEX_URL") or _cfg_get(cfg, "PLEX", "url")
    plex_token = _env_get("ASD_PLEX_TOKEN") or _cfg_get(cfg, "PLEX", "token")
    plex_section = _env_get("ASD_PLEX_SECTION_ID") or _cfg_get(cfg, "PLEX", "section_id")
    plex = None
    if plex_url and plex_token and plex_section:
        plex = PlexConfig(url=plex_url.rstrip("/"), token=plex_token, section_id=str(plex_section).strip())

    # Default behavior: auto-enable only when a target is configured.
    if enabled is None:
        enabled = bool(jellyfin or plex)

    return bool(enabled), float(debounce_s), jellyfin, plex


def _refresh_jellyfin(conf: JellyfinConfig) -> None:
    # Jellyfin supports X-Emby-Token header.
    url = f"{conf.url}/Library/Refresh"
    requests.post(url, headers={"X-Emby-Token": conf.api_key}, timeout=10).raise_for_status()


def _refresh_plex(conf: PlexConfig) -> None:
    url = f"{conf.url}/library/sections/{conf.section_id}/refresh"
    requests.get(url, params={"X-Plex-Token": conf.token}, timeout=10).raise_for_status()


def get_media_refresh_status() -> dict:
    """Return a public (non-secret) status snapshot for UI/logging."""
    enabled, debounce_s, jellyfin, plex = _load_configs()
    return {
        "enabled": bool(enabled),
        "debounce_seconds": float(debounce_s),
        "jellyfin": {
            "configured": bool(jellyfin is not None),
            "url": (jellyfin.url if jellyfin is not None else None),
        },
        "plex": {
            "configured": bool(plex is not None),
            "url": (plex.url if plex is not None else None),
        },
    }


def _resp_to_status(r: Response) -> tuple[bool, int]:
    code = int(getattr(r, "status_code", 0) or 0)
    return (200 <= code < 300, code)


def test_media_apis(timeout_s: float = 5.0) -> dict:
    """Test Jellyfin/Plex API connectivity using configured credentials.

    Note: this may trigger a library refresh.
    Returns a JSON-serializable dict and never raises.
    """
    enabled, _debounce_s, jellyfin, plex = _load_configs()
    timeout = max(1.0, float(timeout_s or 5.0))

    out: dict = {
        "enabled": bool(enabled),
        "timeout_s": float(timeout),
        "jellyfin": {
            "configured": bool(jellyfin is not None),
            "url": (jellyfin.url if jellyfin is not None else None),
            "ok": False,
            "status_code": None,
            "error": None,
        },
        "plex": {
            "configured": bool(plex is not None),
            "url": (plex.url if plex is not None else None),
            "ok": False,
            "status_code": None,
            "error": None,
        },
    }

    if jellyfin is not None:
        try:
            url = f"{jellyfin.url}/Library/Refresh"
            r = requests.post(url, headers={"X-Emby-Token": jellyfin.api_key}, timeout=timeout)
            ok, code = _resp_to_status(r)
            out["jellyfin"]["ok"] = ok
            out["jellyfin"]["status_code"] = code
        except Exception as e:
            out["jellyfin"]["ok"] = False
            out["jellyfin"]["error"] = str(e)

    if plex is not None:
        try:
            url = f"{plex.url}/library/sections/{plex.section_id}/refresh"
            r = requests.get(url, params={"X-Plex-Token": plex.token}, timeout=timeout)
            ok, code = _resp_to_status(r)
            out["plex"]["ok"] = ok
            out["plex"]["status_code"] = code
        except Exception as e:
            out["plex"]["ok"] = False
            out["plex"]["error"] = str(e)

    return out


class _MediaRefreshScheduler:
    def __init__(self) -> None:
        self._enabled, self._debounce_s, self._jellyfin, self._plex = _load_configs()
        self._cv = threading.Condition()
        self._stop = False
        self._pending = False
        self._next_run_ts: float | None = None
        self._inflight = False
        self._thread = threading.Thread(target=self._run_loop, name="media-refresh", daemon=True)
        self._thread.start()

    @property
    def enabled(self) -> bool:
        return bool(self._enabled)

    def trigger(self) -> bool:
        if not self._enabled:
            return False
        with self._cv:
            self._pending = True
            self._next_run_ts = time.time() + float(self._debounce_s)
            self._cv.notify_all()
        return True

    def flush(self, timeout_s: float = 15.0) -> bool:
        if not self._enabled:
            return False
        deadline = time.time() + float(timeout_s)
        with self._cv:
            self._pending = True
            self._next_run_ts = time.time()
            self._cv.notify_all()
            while (self._pending or self._inflight) and time.time() < deadline:
                remaining = max(0.0, deadline - time.time())
                self._cv.wait(timeout=remaining)
        return True

    def shutdown(self, timeout_s: float = 2.0) -> None:
        with self._cv:
            self._stop = True
            self._cv.notify_all()
        try:
            self._thread.join(timeout=float(timeout_s))
        except Exception:
            pass

    def _do_refresh(self) -> None:
        if not self._enabled:
            return
        try:
            if self._jellyfin is not None:
                _refresh_jellyfin(self._jellyfin)
        except Exception:
            # Best-effort; never fail downloads.
            pass
        try:
            if self._plex is not None:
                _refresh_plex(self._plex)
        except Exception:
            pass

    def _run_loop(self) -> None:
        while True:
            with self._cv:
                while not self._stop and not self._pending:
                    self._cv.wait(timeout=1.0)

                if self._stop:
                    return

                assert self._pending
                run_at = self._next_run_ts or time.time()
                delay = max(0.0, run_at - time.time())
                if delay > 0:
                    self._cv.wait(timeout=delay)
                    continue

                # Run now.
                self._pending = False
                self._next_run_ts = None
                self._inflight = True

            try:
                self._do_refresh()
            finally:
                with self._cv:
                    self._inflight = False
                    self._cv.notify_all()


_SCHEDULER: _MediaRefreshScheduler | None = None
_SCHEDULER_LOCK = threading.Lock()


def _get_scheduler() -> _MediaRefreshScheduler:
    global _SCHEDULER
    with _SCHEDULER_LOCK:
        if _SCHEDULER is None:
            _SCHEDULER = _MediaRefreshScheduler()
        return _SCHEDULER


def schedule_media_refresh() -> bool:
    """Schedule a debounced refresh (Jellyfin/Plex) if configured."""
    return _get_scheduler().trigger()


def flush_media_refresh(timeout_s: float = 15.0) -> bool:
    """Force a refresh now and wait up to timeout_s (best-effort)."""
    return _get_scheduler().flush(timeout_s=timeout_s)


def shutdown_media_refresh(timeout_s: float = 2.0) -> None:
    """Stop the refresh background thread (best-effort)."""
    global _SCHEDULER
    with _SCHEDULER_LOCK:
        sched = _SCHEDULER
        _SCHEDULER = None
    if sched is not None:
        sched.shutdown(timeout_s=timeout_s)
