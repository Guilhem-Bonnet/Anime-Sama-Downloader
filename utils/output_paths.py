from __future__ import annotations

import os
import re

from utils.config import get_default_download_path, get_media_separate_lang, get_output_naming_mode


_INVALID_FS_CHARS_RE = re.compile(r'[<>:"/\\|?*]', re.UNICODE)


def _sanitize_component(value: str, *, fallback: str) -> str:
    s = str(value or "").strip()
    s = s.replace(os.sep, " ")
    s = s.replace("/", " ")
    s = _INVALID_FS_CHARS_RE.sub("_", s)
    s = re.sub(r"\s+", " ", s).strip()
    return s or fallback


def _format_series_name(raw: str) -> str:
    s = _sanitize_component(raw, fallback="anime")
    # human-friendly: "sword-art-online" -> "Sword Art Online"
    s = s.replace("-", " ").replace("_", " ")
    s = re.sub(r"\s+", " ", s).strip()
    words: list[str] = []
    for w in s.split(" "):
        # keep existing all-caps tokens
        words.append(w if w.isupper() else (w[:1].upper() + w[1:]))
    return " ".join(words) or "Anime"


def _lang_tag(lang: str) -> str | None:
    l = (lang or "").strip().lower()
    if not l:
        return None
    if l == "vostfr":
        return "VOSTFR"
    if l == "vf":
        return "VF"
    if l == "vo":
        return "VO"
    return l.upper()


def build_episode_output_path(
    dest_root: str,
    anime_slug: str,
    season: int,
    lang: str,
    episode: int,
    ext: str = "mp4",
) -> tuple[str, str]:
    """Return (dest_dir, file_path) for an episode.

    Target layout:
      <dest_root>/<anime_slug>/Saison <season>/<lang>/<anime_slug>-S<season>E<episode>.<ext>

    Notes:
    - We keep season/episode unpadded to match the requested example.
    - Caller is responsible for creating dest_dir.
    """
    safe_slug = _sanitize_component((anime_slug or "anime").strip().strip("/"), fallback="anime")
    safe_lang = _sanitize_component((lang or "vostfr").strip().strip("/"), fallback="vostfr")

    dest_root = (dest_root or get_default_download_path()).strip()
    abs_root = os.path.abspath(os.path.expanduser(dest_root))

    mode = (get_output_naming_mode() or "legacy").strip().lower()
    is_media = mode in {"media", "media-server", "jellyfin", "plex"}

    if is_media:
        series_name = _format_series_name(safe_slug)
        series_dir = series_name
        if get_media_separate_lang():
            tag = _lang_tag(safe_lang)
            if tag:
                series_dir = f"{series_dir} [{tag}]"

        season_dir = f"Season {int(season):02d}"
        dest_dir = os.path.join(abs_root, series_dir, season_dir)
        filename = f"{series_name} - S{int(season):02d}E{int(episode):02d}.{ext.lstrip('.')}"
        return dest_dir, os.path.join(dest_dir, filename)

    dest_dir = os.path.join(abs_root, safe_slug, f"Saison {int(season)}", safe_lang)
    filename = f"{safe_slug}-S{int(season)}E{int(episode)}.{ext.lstrip('.')}"
    return dest_dir, os.path.join(dest_dir, filename)
