from __future__ import annotations

import os

from utils.config import get_default_download_path


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
    safe_slug = (anime_slug or "anime").strip().strip("/") or "anime"
    safe_lang = (lang or "vostfr").strip().strip("/") or "vostfr"

    dest_root = (dest_root or get_default_download_path()).strip()
    abs_root = os.path.abspath(os.path.expanduser(dest_root))
    dest_dir = os.path.join(abs_root, safe_slug, f"Saison {int(season)}", safe_lang)
    filename = f"{safe_slug}-S{int(season)}E{int(episode)}.{ext.lstrip('.')}"
    return dest_dir, os.path.join(dest_dir, filename)
