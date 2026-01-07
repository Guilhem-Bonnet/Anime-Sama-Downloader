"""Anime metadata providers (external databases).

Currently implemented:
- AniList (GraphQL): no API key required for basic search.

This module is used to enrich search queries (titles, synonyms) to better
resolve Anime-Sama catalogue URLs.
"""

from __future__ import annotations

from dataclasses import dataclass
import time
from typing import Iterable

from utils.http_pool import get_http_pool
from utils.config import load_config, save_config


ANILIST_ENDPOINT = "https://graphql.anilist.co"


@dataclass(frozen=True)
class AnimeTitles:
    titles: tuple[str, ...]

    def as_list(self) -> list[str]:
        return list(self.titles)


def _now() -> float:
    return time.time()


def _config_get(config: dict, *path: str, default=None):
    cur = config
    for key in path:
        if not isinstance(cur, dict) or key not in cur:
            return default
        cur = cur[key]
    return cur


def _config_set(config: dict, value, *path: str) -> None:
    cur = config
    for key in path[:-1]:
        cur = cur.setdefault(key, {})
    cur[path[-1]] = value


def _unique_nonempty(strings: Iterable[str]) -> list[str]:
    seen: set[str] = set()
    out: list[str] = []
    for s in strings:
        if not s:
            continue
        s2 = str(s).strip()
        if not s2:
            continue
        if s2.lower() in seen:
            continue
        seen.add(s2.lower())
        out.append(s2)
    return out


def anilist_search_titles(query: str, limit: int = 5, cache_ttl_seconds: int = 7 * 24 * 3600) -> AnimeTitles:
    """Return candidate titles/synonyms for an anime query using AniList.

    No API key is required for basic search.

    Caches results in the existing config file (~/.anime-sama-downloader.json).
    """

    q = (query or "").strip()
    if not q:
        return AnimeTitles(titles=())

    cache_key = q.lower()
    config = load_config()

    cached = _config_get(config, "anime_db", "anilist", cache_key)
    if isinstance(cached, dict):
        ts = cached.get("ts")
        titles = cached.get("titles")
        if isinstance(ts, (int, float)) and isinstance(titles, list):
            if _now() - float(ts) <= cache_ttl_seconds:
                return AnimeTitles(titles=tuple(_unique_nonempty(titles)))

    gql = """
    query ($search: String, $page: Int, $perPage: Int) {
      Page(page: $page, perPage: $perPage) {
        media(search: $search, type: ANIME, sort: POPULARITY_DESC) {
          title {
            romaji
            english
            native
            userPreferred
          }
          synonyms
        }
      }
    }
    """.strip()

    payload = {
        "query": gql,
        "variables": {"search": q, "page": 1, "perPage": max(1, min(limit, 10))},
    }

    headers = {
        "Content-Type": "application/json",
        "Accept": "application/json",
        "User-Agent": "Anime-Sama-Downloader/2.x (Python; AniList)",
    }

    pool = get_http_pool()
    resp = pool.post(ANILIST_ENDPOINT, json=payload, headers=headers, timeout=15)
    resp.raise_for_status()

    data = resp.json() if resp.content else {}
    media = (
        data.get("data", {})
        .get("Page", {})
        .get("media", [])
    )

    candidates: list[str] = []
    for item in media:
        title_obj = item.get("title") or {}
        candidates.extend(
            [
                title_obj.get("userPreferred"),
                title_obj.get("romaji"),
                title_obj.get("english"),
                title_obj.get("native"),
            ]
        )
        syns = item.get("synonyms")
        if isinstance(syns, list):
            candidates.extend(syns)

    titles = _unique_nonempty(candidates)

    _config_set(config, {"ts": _now(), "titles": titles}, "anime_db", "anilist", cache_key)
    save_config(config)

    return AnimeTitles(titles=tuple(titles))
