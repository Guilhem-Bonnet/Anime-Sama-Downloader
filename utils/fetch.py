"""Episode list + player extraction utilities.

This module provides:
- fetch_episodes(base_url): downloads and parses episodes.js
- select_best_player(episodes): auto-selects a usable player
- fetch_video_source(url|list[url]): resolves host pages into direct video URLs (m3u8 or mp4)
"""

from __future__ import annotations

import re
import time
from typing import Dict, List, Optional, Tuple, Union
from urllib.parse import parse_qs, urlparse

import requests

from utils.downloaders.extractor import (
    extract_movearnpre_video_source,
    extract_oneupload_video_source,
    extract_sendvid_video_source,
    extract_sibnet_video_source,
    extract_vidmoly_video_source,
)
from utils.http_pool import cached_get, get_http_pool
from utils.parsers import parse_m3u8_content
from utils.var import print_status


http_pool = get_http_pool()


def _canonicalize_anime_sama_domain(url: str) -> str:
    # anime-sama.tv redirige actuellement vers .si sans conserver le path.
    return re.sub(r"^https://anime-sama\.(?:tv|org|fr|si)/", "https://anime-sama.si/", (url or ""))


def _player_score(player_name: str, urls: List[str]) -> int:
    source_quality = {
        "sendvid": 10,
        "sibnet": 9,
        "m3u8": 8,
        "vidmoly": 7,
        "oneupload": 6,
        "movearnpre": 5,
        "smoothpre": 5,
        "mivalyo": 4,
    }

    available_count = sum(1 for u in urls if u and u.strip())

    quality_score = 0
    for u in urls:
        if not u:
            continue
        u_lower = u.lower()
        for source, score in source_quality.items():
            if source in u_lower:
                quality_score += score
                break
        else:
            quality_score += 3

    return (available_count * 100) + quality_score


def rank_players(episodes: Dict[str, List[str]]) -> List[str]:
    """Return players ordered by heuristics (availability + likely quality).

    This is used for automatic selection and fallback when a host fails.
    """
    if not episodes:
        return []

    scored: List[Tuple[int, str]] = []
    for player_name, urls in episodes.items():
        if not urls:
            continue
        scored.append((_player_score(player_name, urls), player_name))

    scored.sort(reverse=True)
    return [name for _, name in scored]


def select_best_player(episodes: Dict[str, List[str]]) -> Tuple[Optional[str], Optional[str]]:
    """Automatically select the best player based on availability and quality.

    Priority: player with most episodes available, prefer known reliable sources.
    Returns: (player_name, player_key) (same value for backward compatibility).
    """
    if not episodes:
        return None, None

    ranked = rank_players(episodes)
    best_player = ranked[0] if ranked else None

    return best_player, best_player


def fetch_episodes(base_url: str, quiet: bool = False, timeout: int = 20) -> Optional[Dict[str, List[str]]]:
    """Fetch and parse episodes.js for a season/language URL.

    Args:
        base_url: season/language URL
        quiet: suppress stdout logs
        timeout: requests timeout in seconds
    """
    base_url = _canonicalize_anime_sama_domain(base_url)
    js_url = base_url.rstrip("/") + "/episodes.js"
    if not quiet:
        print_status("Fetching episode list...", "loading")

    try:
        headers = {
            "User-Agent": "Mozilla/5.0 (X11; Linux x86_64) Gecko/20100101 Firefox/120.0",
            "Accept": "text/javascript,*/*;q=0.1",
            "Referer": base_url,
        }
        response = cached_get(js_url, headers=headers, timeout=timeout, use_cache=False)
        response.raise_for_status()
    except Exception as e:
        if not quiet:
            print_status(f"Failed to fetch episodes.js: {e}", "error")
        return None

    js_content = response.text
    pattern = re.compile(r"var\s+(eps\d+)\s*=\s*\[([^\]]*)\];", re.MULTILINE)
    matches = pattern.findall(js_content)

    def _is_plausible_episode_url(u: str) -> bool:
        u = (u or "").strip()
        if not (u.startswith("http://") or u.startswith("https://")):
            return False

        # Common placeholder forms when a season exists structurally but has no real videos.
        if u.endswith("="):
            return False
        if u.endswith("/embed/"):
            return False

        # Anime-Sama sometimes exposes placeholder VK embeds for non-existing seasons.
        # Example: https://vk.com/video_ext.php?oid=&hd=3
        if "vk.com/video_ext.php" in u:
            try:
                q = parse_qs(urlparse(u).query)
                oid = (q.get("oid") or [""])[0]
                vid = (q.get("id") or [""])[0]
            except Exception:
                return False
            if not re.fullmatch(r"-?\d+", oid or ""):
                return False
            if vid and not re.fullmatch(r"\d+", vid):
                return False

        # Sibnet placeholder: missing videoid
        if "video.sibnet.ru" in u and "shell.php" in u:
            try:
                q = parse_qs(urlparse(u).query)
                videoid = (q.get("videoid") or [""])[0]
            except Exception:
                return False
            if not re.fullmatch(r"\d+", videoid or ""):
                return False

        # Vidmoly placeholder: embed-.html or empty id
        if "vidmoly" in u and "/embed-" in u:
            p = urlparse(u).path
            m = re.match(r"^/embed-([A-Za-z0-9]+)\.html$", p or "")
            if not m:
                return False

        # SendVid placeholder: embed/ with empty id
        if "sendvid.com" in u and "/embed/" in u:
            p = urlparse(u).path or ""
            # /embed/<id>
            parts = [x for x in p.split("/") if x]
            if len(parts) < 2 or parts[0] != "embed" or not parts[1].strip():
                return False

        return True

    episodes: Dict[str, List[str]] = {}
    for name, content in matches:
        player_num_match = re.search(r"\d+", name)
        if not player_num_match:
            continue
        player_num = player_num_match.group()
        player_name = f"Player {player_num}"
        # Preserve indexes: episodes.js may contain empty strings for missing episodes.
        raw_items = re.findall(r"'([^']*)'", content)
        urls: List[str] = []
        for item in raw_items:
            item = (item or "").strip()
            if _is_plausible_episode_url(item):
                urls.append(item)
            else:
                urls.append("")
        episodes[player_name] = urls

    # Some seasons exist as empty arrays (no episodes). Treat as missing.
    has_any_url = any(_is_plausible_episode_url(u) for urls in episodes.values() for u in (urls or []))
    if not has_any_url:
        if not quiet:
            print_status("No episodes found in episodes.js", "error")
        return None

    if not quiet:
        print_status(f"Found {len(episodes)} players with episodes!", "success")
    return episodes


def get_sibnet_redirect_location(video_url: str) -> Optional[str]:
    headers = {
        "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Gecko/20100101 Firefox/108.0",
        "Accept": "video/webm,video/mp4,video/*;q=0.9,*/*;q=0.8",
        "Accept-Language": "en-US,en;q=0.5",
        "Referer": "https://video.sibnet.ru/",
    }
    try:
        response = http_pool.get(video_url, headers=headers, allow_redirects=False, timeout=15)
        if response.status_code == 302:
            redirect_url = response.headers.get("location")
            if not redirect_url:
                return None
            if redirect_url.startswith("//"):
                return f"https:{redirect_url}"
            return redirect_url
        print_status(f"Expected redirect (302), got {response.status_code}", "warning")
        return None
    except requests.RequestException as e:
        print_status(f"Failed to get redirect location: {e}", "error")
        return None


def fetch_page_content(url: str) -> Optional[str]:
    headers = {
        "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Gecko/20100101 Firefox/108.0",
        "Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8",
        "Accept-Language": "en-US,en;q=0.5",
        "Referer": "https://sendvid.com/" if "sendvid.com" in url else "https://oneupload.net/",
    }
    try:
        print_status("Connecting to server...", "loading")
        response = http_pool.get(url, headers=headers, timeout=20)
        response.raise_for_status()
        return response.text
    except requests.RequestException as e:
        print_status(f"Failed to connect to {url}: {e}", "error")
        return None


def _best_stream_from_m3u8(m3u8_url: str, referer: str) -> Optional[str]:
    headers = {
        "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Gecko/20100101 Firefox/108.0",
        "Referer": referer,
    }
    response = http_pool.get(m3u8_url, headers=headers, timeout=20)
    response.raise_for_status()

    streams = parse_m3u8_content(response.text, m3u8_url)
    if not streams:
        return None
    best_stream = max(streams, key=lambda x: int(x.get("BANDWIDTH", 0)))
    return best_stream.get("url")


def fetch_video_source(url: Union[str, List[str]]):
    """Resolve one or many episode host URLs into direct playable URLs."""

    def process_single_url(single_url: str) -> Optional[str]:
        if not single_url:
            return None

        print_status(f"Processing video URL: {single_url[:50]}...", "loading")

        # Domain normalizations
        if "vidmoly.to" in single_url:
            single_url = single_url.replace("vidmoly.to", "vidmoly.net")
            print_status("Converted vidmoly.to to vidmoly.net", "info")
        if "oneupload.to" in single_url:
            single_url = single_url.replace("oneupload.to", "oneupload.net")

        # SENDVID
        if "sendvid.com" in single_url:
            html_content = fetch_page_content(single_url)
            return extract_sendvid_video_source(html_content)

        # SIBNET
        if "video.sibnet.ru" in single_url:
            html_content = fetch_page_content(single_url)
            video_source = extract_sibnet_video_source(html_content)
            if video_source:
                print_status("Getting direct download link...", "loading")
                return get_sibnet_redirect_location(video_source)
            return None

        # ONEUPLOAD
        if "oneupload.net" in single_url:
            html_content = fetch_page_content(single_url)
            m3u8_url = extract_oneupload_video_source(html_content)
            if not m3u8_url:
                return None
            try:
                best = _best_stream_from_m3u8(m3u8_url, referer="https://oneupload.net/")
                if not best:
                    print_status("No video streams found in M3U8 playlist", "error")
                return best
            except requests.RequestException as e:
                print_status(f"Failed to fetch M3U8 playlist: {e}", "error")
                return None

        # VIDMOLY
        if "vidmoly.net" in single_url:
            html_content = fetch_page_content(single_url)
            m3u8_url = extract_vidmoly_video_source(html_content)
            if not m3u8_url:
                return None
            try:
                best = _best_stream_from_m3u8(m3u8_url, referer="https://vidmoly.net/")
                if not best:
                    print_status("No video streams found in M3U8 playlist", "error")
                return best
            except requests.RequestException as e:
                print_status(f"Failed to fetch M3U8 playlist: {e}", "error")
                return None

        # MOVEARNPRE / MIVALYO / SMOOTHPRE (packed JS -> HLS)
        if any(h in single_url.lower() for h in ["movearnpre.com", "mivalyo.com", "smoothpre.com"]):
            m3u8_url = extract_movearnpre_video_source(single_url)
            if not m3u8_url:
                return None

            referer = "https://movearnpre.com/"
            if "mivalyo.com" in single_url.lower():
                referer = "https://mivalyo.com/"
            if "smoothpre.com" in single_url.lower():
                referer = "https://smoothpre.com/"

            max_retries = 5
            for attempt in range(max_retries):
                try:
                    print_status("Parsing M3U8 content...", "loading")
                    best = _best_stream_from_m3u8(m3u8_url, referer=referer)
                    if not best:
                        print_status("No video streams found in M3U8 playlist", "error")
                        return None
                    return best
                except requests.RequestException as e:
                    if attempt < max_retries - 1:
                        print_status(f"Attempt {attempt + 1} failed: {e}. Retrying...", "warning")
                        time.sleep(1)
                        continue
                    print_status(
                        "Service not responding or extraction failed (retry exhausted). Try again later.",
                        "error",
                    )
                    return None

        print_status(
            "Unsupported video source. Only some hosts are supported (see README).",
            "error",
        )
        return None

    if isinstance(url, str):
        return process_single_url(url)
    if isinstance(url, list):
        return [process_single_url(u) for u in url]

    print_status("Invalid input: URL must be a string or a list of strings.", "error")
    return None