from __future__ import annotations

from typing import Callable


UIRunner = Callable[[], int]


def _textual_runner() -> int:
    # Lazy import so the core app can run without Textual installed.
    from utils.tui import run_tui

    return int(run_tui())


def _web_runner() -> int:
    # Lazy import so the core app can run without FastAPI/Uvicorn installed.
    from utils.ui.web import run_web

    return int(run_web())


_UI_RUNNERS: dict[str, UIRunner] = {
    "textual": _textual_runner,
    "web": _web_runner,
}


def list_ui_plugins() -> list[str]:
    return sorted(_UI_RUNNERS.keys())


def run_ui(name: str) -> int:
    key = (name or "").strip().lower()
    if not key:
        raise ValueError("UI plugin name is required")

    runner = _UI_RUNNERS.get(key)
    if runner is None:
        raise ValueError(f"Unknown UI plugin: {name}. Available: {', '.join(list_ui_plugins())}")

    return runner()
