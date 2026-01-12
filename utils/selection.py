from __future__ import annotations

import re


def parse_tracker_selection(text: str, default_season: int) -> list[tuple[int | None, str]]:
    """Parse: S2E14, S1E2,7, S1E1-6, S1EALL, ALL.

    Returns list of (season, spec). season=None means ALL seasons.
    """
    s = (text or "").strip().upper()
    if not s:
        return [(default_season, "1")]

    # UX: most people expect "ALL" = all episodes of the selected season.
    if s in {"ALL", "*"}:
        return [(default_season, "ALL")]
    # Explicit keyword for all seasons.
    if s in {"ALLSEASONS", "ALL-SEASONS", "ALLS"}:
        return [(None, "ALL")]

    parts = [p for p in s.replace(";", " ").split() if p.strip()]
    out: list[tuple[int | None, str]] = []
    for part in parts:
        m = re.match(r"^S(?P<season>\d+)(?:E(?P<eps>.+))?$", part)
        if not m:
            # episodes-only for default season
            out.append((default_season, part.replace("E", "")))
            continue
        season = int(m.group("season"))
        eps = (m.group("eps") or "").strip()
        out.append((season, "ALL" if not eps or eps == "ALL" else eps))

    return out or [(default_season, "1")]


def parse_episode_expr(expr: str, max_episodes: int) -> list[int]:
    e = (expr or "").strip().lower()
    if not e:
        return []
    if e in {"all", "*"}:
        return list(range(1, max_episodes + 1))

    indices: list[int] = []
    for part in e.split(","):
        part = part.strip()
        if not part:
            continue
        if "-" in part:
            a, b = part.split("-", 1)
            try:
                start = int(a.strip())
                end = int(b.strip())
            except ValueError:
                continue
            if start > end:
                start, end = end, start
            indices.extend(range(start, end + 1))
        else:
            try:
                indices.append(int(part))
            except ValueError:
                continue

    return sorted({i for i in indices if 1 <= i <= max_episodes})
