from __future__ import annotations

import asyncio
import os
import re
from dataclasses import dataclass
import threading

from textual.app import App, ComposeResult
from textual.containers import Horizontal, Vertical
from textual.widgets import Button, Footer, Header, Input, Label, RichLog, Select, Static, TabbedContent, TabPane

from utils.config import (
    get_default_download_path,
    get_max_concurrent_downloads,
    get_preferred_download_root,
    set_last_used_path,
    set_preferred_download_root,
)
from utils.downloaders.downloader import download_video
from utils.fetch import fetch_episodes, fetch_video_source, rank_players
from utils.search import resolve_anime_sama_base_url
from utils.download_manager import DownloadJob, DownloadManager
from utils.output_paths import build_episode_output_path
from utils.selection import parse_episode_expr, parse_tracker_selection


@dataclass
class DownloadPlan:
    base_url: str | None = None
    lang: str = "vostfr"


class StatusBox(Static):
    def set(self, text: str) -> None:
        self.update(text)


class AnimeSamaTUI(App):
    CSS = """
    Screen { align: center top; }
    #root { width: 100%; height: 100%; }
    #panel { width: 96%; height: 100%; overflow-y: auto; }
    .row { height: auto; }
    #status { height: 5; }
    #search { width: 20; }
    #detect_seasons { width: 22; }
    #load { width: 24; }
    #download { width: 22; }
    #availability { height: 3; }
    #selection { width: 1fr; }
    #log { height: 12; }
    """

    BINDINGS = [("q", "quit", "Quit")]

    def __init__(self):
        super().__init__()
        self.plan = DownloadPlan()
        # Episodes cache: {season: {player: [urls...]}}
        self.episodes_by_season: dict[int, dict[str, list[str]]] = {}
        self.detected_seasons: list[int] = []
        self.base_download_root: str = os.path.abspath(os.path.expanduser(get_default_download_path()))
        self._found_url: str | None = None
        self._seasons_summary: str | None = None
        self._status_msg: str = "Prêt."
        self._is_downloading: bool = False
        self._is_detecting_seasons: bool = False
        self._dl_manager: DownloadManager | None = None

    def compose(self) -> ComposeResult:
        yield Header(show_clock=True)

        with Vertical(id="root"):
            with Vertical(id="panel"):
                yield Label("Recherche → saisons dispo → sélection (anime-tracker) → téléchargement", classes="row")

                yield StatusBox("Prêt.", id="status")

                with TabbedContent(id="tabs"):
                    with TabPane("Télécharger", id="tab-download"):
                        with Vertical(classes="row"):
                            yield Input(placeholder="Nom de l’anime (ex: One Piece)", id="query")
                            with Horizontal(classes="row"):
                                yield Button("Chercher", id="search")
                                yield Button("Rafraîchir saisons", id="detect_seasons")

                        with Horizontal(classes="row"):
                            yield Label("Langue:")
                            yield Select([(v, v) for v in ("vostfr", "vf", "vo")], value="vostfr", id="lang")
                            yield Label("Saison:")
                            yield Select([], id="season")

                        with Horizontal(classes="row"):
                            yield Label("Sélection:")
                            yield Input(
                                placeholder="S2E14 | S1E2,7 | S1E1-6 | S1EALL | ALL (= saison sélectionnée) | ALLSEASONS",
                                value="",
                                id="selection",
                            )

                        yield Static("Disponibilité: —", id="availability", classes="row")
                        yield Static("Destination: — (onglet Options)", id="destination", classes="row")

                        yield RichLog(id="log", highlight=True, wrap=True, markup=False, classes="row")

                        with Horizontal(classes="row"):
                            yield Button("Charger dispo", id="load")
                            yield Button("Télécharger", id="download", variant="success")

                        with Horizontal(classes="row"):
                            yield Button("Annuler tout", id="cancel_all")
                            yield Button("Vider file", id="clear_queue")

                    with TabPane("Téléchargements", id="tab-queue"):
                        yield Static("File: —", id="queue_status", classes="row")
                        yield RichLog(id="queue_log", highlight=True, wrap=True, markup=False, classes="row")

                    with TabPane("Options", id="tab-options"):
                        yield Static(
                            "Destination par défaut: disque monté avec le plus d’espace libre (Linux), sinon ~/Téléchargements ou ~/Downloads.",
                            classes="row",
                        )
                        yield Input(
                            placeholder="Dossier racine (ex: /mnt/HDD/Anime-Sama-Downloader ou ~/Téléchargements)",
                            value=self.base_download_root,
                            id="download_root",
                        )
                        with Horizontal(classes="row"):
                            yield Button("Utiliser défaut smart", id="use_default")
                            yield Button("Sauver options", id="save_options", variant="success")

        yield Footer()

    def on_mount(self) -> None:
        self._focus_download_tab()
        self._update_season_options()
        self._update_destination_label()
        self._render_status()
        try:
            self.query_one("#query", Input).focus()
        except Exception:
            pass

        try:
            self.query_one("#log", RichLog).write("Prêt.")
        except Exception:
            pass

        self._ensure_download_manager()
        self._refresh_queue_status()

    def _ensure_download_manager(self) -> None:
        if self._dl_manager is not None:
            return

        def _on_event(job: DownloadJob, event: str) -> None:
            msg = f"[{event.upper()}] {job.label}"
            if job.status == "FAILED" and job.error:
                msg += f" | {job.error}"
            if job.status == "SUCCESS" and job.result_path:
                msg += f" | {job.result_path}"

            try:
                self.call_from_thread(self._append_queue_log, msg)
                self.call_from_thread(self._refresh_queue_status)
            except Exception:
                pass

        self._dl_manager = DownloadManager(
            max_parallel=get_max_concurrent_downloads(default=10),
            on_event=_on_event,
            executor_name="tui-dl",
        )

    def _append_queue_log(self, msg: str) -> None:
        try:
            self.query_one("#queue_log", RichLog).write(msg)
        except Exception:
            self._append_log(msg)

    def _refresh_queue_status(self) -> None:
        mgr = self._dl_manager
        if mgr is None:
            return
        try:
            pending = mgr.pending_count()
            running = mgr.running_count()
            total = len(mgr.list_jobs())
            self.query_one("#queue_status", Static).update(
                f"File: {pending} en attente | {running} en cours | total {total}"
            )
        except Exception:
            pass

    def _focus_download_tab(self) -> None:
        try:
            tabs = self.query_one("#tabs", TabbedContent)
            tabs.active = "tab-download"
        except Exception:
            pass

    def _set_status(self, msg: str) -> None:
        self._status_msg = msg
        self._render_status()

    def _append_log(self, msg: str) -> None:
        try:
            self.query_one("#log", RichLog).write(msg)
        except Exception:
            pass

    def _log(self, msg: str) -> None:
        self._set_status(msg)
        self._append_log(msg)

    def _thread_log(self, msg: str) -> None:
        try:
            self.call_from_thread(self._append_log, msg)
            self.call_from_thread(self._set_status, msg)
        except Exception:
            pass

    def _set_download_enabled(self, enabled: bool) -> None:
        # Backward compatible name: actually toggles all main controls.
        self._set_controls_enabled(enabled)

    def _set_controls_enabled(self, enabled: bool) -> None:
        for btn_id in ("search", "detect_seasons", "load", "download"):
            try:
                self.query_one(f"#{btn_id}", Button).disabled = not enabled
            except Exception:
                pass

        for widget_id in ("query", "lang", "season", "selection"):
            try:
                self.query_one(f"#{widget_id}").disabled = not enabled
            except Exception:
                pass

    def _render_status(self) -> None:
        url = self._found_url or "—"
        seasons = self._seasons_summary or "—"
        dest = self.base_download_root or "—"
        text = f"URL: {url}\nSaisons: {seasons}\nDestination: {dest}\nÉtat: {self._status_msg}"
        self.query_one("#status", StatusBox).set(text)

    def _update_destination_label(self) -> None:
        dest = self.base_download_root
        self.query_one("#destination", Static).update(f"Destination: {dest} (onglet Options pour changer)")

    def _update_season_options(self) -> None:
        season_select = self.query_one("#season", Select)
        if self.detected_seasons:
            season_select.set_options([(f"S{s}", str(s)) for s in self.detected_seasons])
            if season_select.value is None or str(season_select.value) not in {str(s) for s in self.detected_seasons}:
                season_select.value = str(self.detected_seasons[0])
        else:
            season_select.set_options([(f"S{i}", str(i)) for i in range(1, 21)])
            if not season_select.value:
                season_select.value = "1"

    def _availability_summary(self, episodes: dict[str, list[str]]) -> str:
        if not episodes:
            return "Disponibilité: —"

        max_len = max((len(v) for v in episodes.values() if v), default=0)
        if max_len <= 0:
            return "Disponibilité: 0 épisode"

        available: list[int] = []
        missing: list[int] = []
        for i in range(max_len):
            ok = any((i < len(urls) and urls[i] and urls[i].strip()) for urls in episodes.values())
            if ok:
                available.append(i + 1)
            else:
                missing.append(i + 1)

        if not available:
            return f"Disponibilité: 0/{max_len}"

        avail_str = f"{available[0]}-{available[-1]} ({len(available)}/{max_len})"
        if missing:
            shown = ",".join(map(str, missing[:12]))
            more = "…" if len(missing) > 12 else ""
            return f"Disponibilité: {avail_str} | manquants: {shown}{more}"
        return f"Disponibilité: {avail_str}"

    async def on_button_pressed(self, event: Button.Pressed) -> None:
        bid = event.button.id
        if bid == "search":
            await self._do_search()
        elif bid == "detect_seasons":
            await self._do_detect_seasons()
        elif bid == "load":
            await self._do_load_episodes()
        elif bid == "download":
            await self._do_download()
        elif bid == "cancel_all":
            self._do_cancel_all()
        elif bid == "clear_queue":
            self._do_clear_queue()
        elif bid == "use_default":
            self._use_default_root()
        elif bid == "save_options":
            self._save_options()

    def _do_cancel_all(self) -> None:
        self._ensure_download_manager()
        if self._dl_manager is None:
            return
        self._dl_manager.cancel_all()
        self._append_queue_log("[INFO] Annulation demandée pour tous les jobs.")
        self._refresh_queue_status()

    def _do_clear_queue(self) -> None:
        self._ensure_download_manager()
        if self._dl_manager is None:
            return
        cleared = self._dl_manager.clear_pending()
        self._append_queue_log(f"[INFO] File vidée: {cleared} job(s) annulé(s).")
        self._refresh_queue_status()

    async def _do_search(self) -> None:
        self._focus_download_tab()
        query = self.query_one("#query", Input).value.strip()
        if not query:
            self._set_status("Entre un nom d’anime.")
            return

        self._log("Recherche (AniList + résolution URL)…")

        try:
            base = await asyncio.to_thread(resolve_anime_sama_base_url, query, provider="anilist")
        except Exception as e:
            self._log(f"Erreur recherche: {e}")
            return

        if not base:
            self._log("Aucun match trouvé. Essaie un nom plus précis.")
            return

        self.plan.base_url = base
        self.episodes_by_season = {}
        self.detected_seasons = []
        self._update_destination_label()
        self._found_url = base
        self._seasons_summary = "(détection…)"
        self._log(f"Trouvé: {base}")

        # UX: on enchaîne directement sur la détection des saisons.
        self._set_status("Détection des saisons…")
        try:
            asyncio.create_task(self._do_detect_seasons())
        except Exception:
            # Fallback si create_task échoue (rare): l’utilisateur peut cliquer sur Rafraîchir.
            self._set_status("Clique sur “Rafraîchir saisons”.")

    def _make_episode_job(
        self,
        base_catalogue_url: str,
        season: int,
        ep_num: int,
        episodes: dict[str, list[str]],
        lang: str,
        dest_root: str,
    ) -> DownloadJob:
        idx0 = ep_num - 1
        ranked = rank_players(episodes) or list(episodes.keys())

        slug = self._anime_slug() or "anime"
        dest_dir, save_path = build_episode_output_path(dest_root, slug, season, lang, ep_num, ext="mp4")
        os.makedirs(dest_dir, exist_ok=True)

        # Avec 10 jobs en parallèle, on réduit la fan-out interne par job.
        ts_workers = 2
        mp4_workers = 2
        use_ts_threading = False

        label = f"{slug} S{season}E{ep_num}"

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

                self._thread_log(f"{label}: source via {player}…")
                try:
                    src = fetch_video_source(page_url)
                except Exception as e:
                    self._thread_log(f"{label}: erreur source ({player}): {e}")
                    continue

                if not src:
                    continue

                self._thread_log(f"{label}: téléchargement ({player})…")

                def _cb(message: str, status_type: str = "info") -> None:
                    self._thread_log(f"{label}: {message}")

                ok, out = download_video(
                    src,
                    save_path,
                    use_ts_threading=use_ts_threading,
                    url=page_url,
                    automatic_mp4=False,
                    log_callback=_cb,
                    use_tqdm=False,
                    use_mp4_threading=True,
                    mp4_workers=mp4_workers,
                    ts_workers=ts_workers,
                    cancel_event=cancel_event,
                )
                if ok and out:
                    return out

            return None

        return DownloadJob(label=label, run=_runner)

    def _use_default_root(self) -> None:
        self.base_download_root = os.path.abspath(os.path.expanduser(get_default_download_path()))
        self.query_one("#download_root", Input).value = self.base_download_root
        self._update_destination_label()
        self._set_status("Dossier défaut appliqué.")

    def _save_options(self) -> None:
        root = self.query_one("#download_root", Input).value.strip() or get_default_download_path()
        root = os.path.abspath(os.path.expanduser(root))
        self.base_download_root = root
        set_preferred_download_root(root)
        self._update_destination_label()
        self._set_status("Options sauvegardées.")

    async def _do_detect_seasons(self) -> None:
        self._focus_download_tab()
        if not self.plan.base_url:
            self._set_status("Fais une recherche d’abord.")
            return

        if self._is_detecting_seasons:
            self._set_status("Détection des saisons déjà en cours…")
            return

        self._is_detecting_seasons = True

        lang = str(self.query_one("#lang", Select).value or "vostfr")
        self.plan.lang = lang
        self._log("Détection des saisons disponibles…")

        base_url = self.plan.base_url.rstrip("/")

        def _probe() -> tuple[list[int], dict[int, dict[str, list[str]]]]:
            seasons: list[int] = []
            episodes_by_season: dict[int, dict[str, list[str]]] = {}
            consecutive_misses = 0
            max_seasons_to_probe = 50
            max_consecutive_misses = 3
            max_consecutive_misses_start = 5

            for season in range(1, max_seasons_to_probe + 1):
                self._thread_log(f"Sondage saisons: S{season}/{max_seasons_to_probe}…")
                full_url = f"{base_url}/saison{season}/{lang}/"
                # Sondage plus rapide: timeout réduit.
                eps = fetch_episodes(full_url, quiet=True, timeout=8)
                if eps:
                    seasons.append(season)
                    episodes_by_season[season] = eps
                    consecutive_misses = 0
                else:
                    consecutive_misses += 1
                    if seasons and consecutive_misses >= max_consecutive_misses:
                        break
                    if not seasons and consecutive_misses >= max_consecutive_misses_start:
                        break

            return seasons, episodes_by_season

        try:
            seasons, episodes_by_season = await asyncio.to_thread(_probe)
        finally:
            self._is_detecting_seasons = False

        for season, eps in episodes_by_season.items():
            self.episodes_by_season[season] = eps

        self.detected_seasons = seasons
        self._update_season_options()
        self._update_destination_label()

        if seasons:
            self._seasons_summary = ", ".join("S" + str(s) for s in seasons)
            self._log(f"Saisons détectées: {', '.join('S'+str(s) for s in seasons)}")
            await self._do_load_episodes()
        else:
            self._seasons_summary = "—"
            self._log("Impossible de détecter les saisons (blocage réseau ?). Fallback: S1..S20")

    async def _do_load_episodes(self) -> None:
        self._focus_download_tab()
        if not self.plan.base_url:
            self._set_status("Fais une recherche d’abord.")
            return

        season = int(self.query_one("#season", Select).value or "1")
        lang = str(self.query_one("#lang", Select).value or "vostfr")
        self.plan.lang = lang

        episodes = self.episodes_by_season.get(season)
        if not episodes:
            full_url = f"{self.plan.base_url.rstrip('/')}/saison{season}/{lang}/"
            self._log("Chargement de la disponibilité…")
            try:
                episodes = await asyncio.to_thread(fetch_episodes, full_url, quiet=True)
            except Exception as e:
                self._log(f"Erreur episodes: {e}")
                return
            if not episodes:
                self._log("Impossible de récupérer les épisodes (URL invalide ou blocage).")
                return
            self.episodes_by_season[season] = episodes

        self.query_one("#availability", Static).update(self._availability_summary(episodes))
        self._log("Disponibilité chargée.")

    # parse_tracker_selection + parse_episode_expr moved to utils.selection

    def _anime_slug(self) -> str:
        if not self.plan.base_url:
            return "anime"
        slug = self.plan.base_url.rstrip("/").split("/")[-1]
        return slug or "anime"

    def _season_dir(self, season: int, lang: str) -> str:
        root = os.path.abspath(os.path.expanduser(self.base_download_root))
        return os.path.join(root, self._anime_slug(), f"S{season}", lang)

    def _max_episode_count(self, episodes: dict[str, list[str]]) -> int:
        return max((len(v) for v in episodes.values() if v), default=0)

    def _episode_is_available(self, episodes: dict[str, list[str]], idx0: int) -> bool:
        return any((idx0 < len(urls) and urls[idx0] and urls[idx0].strip()) for urls in episodes.values())

    async def _download_one_episode_with_fallback(
        self,
        season: int,
        ep_num: int,
        episodes: dict[str, list[str]],
        lang: str,
    ) -> str | None:
        idx0 = ep_num - 1
        ranked = rank_players(episodes) or list(episodes.keys())
        dest_dir = self._season_dir(season, lang)
        os.makedirs(dest_dir, exist_ok=True)
        save_path = os.path.join(dest_dir, f"episode_{ep_num}.mp4")

        for player in ranked:
            urls = episodes.get(player) or []
            if idx0 >= len(urls):
                continue
            page_url = urls[idx0]
            if not page_url:
                continue

            self._log(f"S{season}E{ep_num}: source via {player}…")

            def _resolve_source() -> str | None:
                try:
                    return fetch_video_source(page_url)
                except Exception as e:  # pragma: no cover
                    self._thread_log(f"S{season}E{ep_num}: erreur source ({player}): {e}")
                    return None

            src = await asyncio.to_thread(_resolve_source)

            if not src:
                continue

            self._log(f"S{season}E{ep_num}: téléchargement ({player})…")

            def _download() -> bool:
                try:
                    def _cb(message: str, status_type: str = "info") -> None:
                        self._thread_log(message)

                    ok, _ = download_video(
                        src,
                        save_path,
                        use_ts_threading=True,
                        url=page_url,
                        automatic_mp4=False,
                        log_callback=_cb,
                        use_tqdm=False,
                        use_mp4_threading=True,
                    )
                    return bool(ok)
                except Exception as e:  # pragma: no cover
                    self._thread_log(f"S{season}E{ep_num}: erreur téléchargement ({player}): {e}")
                    return False

            ok = await asyncio.to_thread(_download)
            if ok:
                self._log(f"S{season}E{ep_num}: OK → {save_path}")
                return save_path

        return None

    async def _do_download(self) -> None:
        self._ensure_download_manager()
        if self._dl_manager is None:
            self._log("Erreur: gestionnaire de téléchargement indisponible.")
            return

        if not self.plan.base_url:
            self._log("Fais une recherche d’abord.")
            return

        lang = str(self.query_one("#lang", Select).value or "vostfr")
        self.plan.lang = lang

        default_season = int(self.query_one("#season", Select).value or "1")
        selection_text = self.query_one("#selection", Input).value
        selection = parse_tracker_selection(selection_text, default_season=default_season)

        # Determine seasons to download
        all_seasons_requested = any(season is None for season, _ in selection)
        if all_seasons_requested:
            seasons_to_download = self.detected_seasons or list(range(1, 21))
        else:
            seasons_to_download = sorted({int(season) for season, _ in selection if isinstance(season, int)})

        if not seasons_to_download:
            self._log("Aucune saison à télécharger.")
            return

        # Apply options
        root_override = self.query_one("#download_root", Input).value.strip()
        if root_override:
            self.base_download_root = os.path.abspath(os.path.expanduser(root_override))
        self._update_destination_label()

        if get_preferred_download_root() is None:
            set_preferred_download_root(self.base_download_root)

        # Load episodes for needed seasons (snapshot)
        base_url = self.plan.base_url.rstrip("/")
        for season in seasons_to_download:
            if season in self.episodes_by_season:
                continue
            full_url = f"{base_url}/saison{season}/{lang}/"
            self._log(f"Chargement S{season}…")
            eps = await asyncio.to_thread(fetch_episodes, full_url, quiet=True)
            if not eps:
                # Ne pas bloquer tout le batch si une saison est vide/inexistante.
                self._log(f"S{season}: pas d'épisodes (saison inexistante ou vide), ignorée.")
                continue
            self.episodes_by_season[season] = eps

        # Build per-season specs
        per_season_specs: dict[int, list[str]] = {}
        if all_seasons_requested:
            for season in seasons_to_download:
                per_season_specs.setdefault(season, []).append("ALL")
        else:
            for season, spec in selection:
                if season is None:
                    continue
                per_season_specs.setdefault(int(season), []).append(spec)

        enqueued = 0
        for season in seasons_to_download:
            episodes = self.episodes_by_season.get(season)
            if not episodes:
                continue

            max_eps = self._max_episode_count(episodes)
            if max_eps <= 0:
                continue

            specs = per_season_specs.get(season, ["1"])
            requested: set[int] = set()
            for spec in specs:
                if str(spec).upper() == "ALL":
                    requested.update(range(1, max_eps + 1))
                else:
                    requested.update(parse_episode_expr(str(spec), max_episodes=max_eps))

            requested = {n for n in requested if self._episode_is_available(episodes, n - 1)}
            if not requested:
                self._log(f"S{season}: aucun épisode disponible pour la sélection.")
                continue

            for ep in sorted(requested):
                job = self._make_episode_job(base_url, season, ep, episodes, lang, self.base_download_root)
                self._dl_manager.enqueue(job)
                enqueued += 1

        if enqueued:
            try:
                set_last_used_path(self.base_download_root)
            except Exception:
                pass
            self._log(f"Ajouté à la file: {enqueued} épisode(s). Tu peux continuer à chercher/ajouter.")
            self._refresh_queue_status()
        else:
            self._log("Rien à ajouter à la file (sélection vide ou indisponible).")


def run_tui() -> int:
    app = AnimeSamaTUI()
    app.run()
    return 0
