from __future__ import annotations

import os
from dataclasses import dataclass

from textual.app import App, ComposeResult
from textual.containers import Horizontal, Vertical
from textual.widgets import Button, Footer, Header, Input, Label, Select, Static

from utils.search import resolve_anime_sama_base_url
from utils.fetch import fetch_episodes, select_best_player, fetch_video_source
from utils.downloaders.downloader import download_video


@dataclass
class DownloadPlan:
    base_url: str | None = None
    season: int = 1
    lang: str = "vostfr"
    player: str | None = None
    episodes: str = "1"
    directory: str = "./videos"


class StatusBox(Static):
    def set(self, text: str) -> None:
        self.update(text)


class AnimeSamaTUI(App):
    CSS = """
    Screen { align: center middle; }
    #root { width: 100%; height: 100%; }
    #panel { width: 96%; height: auto; }
    .row { height: auto; }
    #status { height: 4; }
    """

    BINDINGS = [("q", "quit", "Quit")]

    def __init__(self):
        super().__init__()
        self.plan = DownloadPlan()
        self.episodes_data: dict[str, list[str]] | None = None

    def compose(self) -> ComposeResult:
        yield Header(show_clock=True)

        with Vertical(id="root"):
            with Vertical(id="panel"):
                yield Label("Recherche anime → saison/langue → épisodes → téléchargement", classes="row")

                with Horizontal(classes="row"):
                    yield Input(placeholder="Nom de l’anime (ex: Kaiju No 8)", id="query")
                    yield Button("Chercher", id="search")

                with Horizontal(classes="row"):
                    yield Label("Saison:")
                    yield Select([(str(i), str(i)) for i in range(1, 21)], value="1", id="season")
                    yield Label("Langue:")
                    yield Select([(v, v) for v in ("vostfr", "vf", "vo")], value="vostfr", id="lang")

                with Horizontal(classes="row"):
                    yield Label("Player:")
                    yield Select([], id="player")

                with Horizontal(classes="row"):
                    yield Input(placeholder="Épisodes (ex: 1-5, 1,2,3 ou all)", value="1", id="episodes")
                    yield Input(placeholder="Dossier (ex: ./videos)", value="./videos", id="directory")

                with Horizontal(classes="row"):
                    yield Button("Charger épisodes", id="load")
                    yield Button("Télécharger", id="download", variant="success")

                yield StatusBox("Prêt.", id="status")

        yield Footer()

    def _set_status(self, msg: str) -> None:
        self.query_one("#status", StatusBox).set(msg)

    async def on_button_pressed(self, event: Button.Pressed) -> None:
        bid = event.button.id
        if bid == "search":
            await self._do_search()
        elif bid == "load":
            await self._do_load_episodes()
        elif bid == "download":
            await self._do_download()

    async def _do_search(self) -> None:
        query = self.query_one("#query", Input).value.strip()
        if not query:
            self._set_status("Entre un nom d’anime.")
            return

        self._set_status("Recherche (AniList + résolution URL)…")

        try:
            base = resolve_anime_sama_base_url(query, provider="anilist")
        except Exception as e:
            self._set_status(f"Erreur recherche: {e}")
            return

        if not base:
            self._set_status("Aucun match trouvé. Essaie un nom plus précis.")
            return

        self.plan.base_url = base
        self._set_status(f"Trouvé: {base}")

    async def _do_load_episodes(self) -> None:
        if not self.plan.base_url:
            self._set_status("Fais une recherche d’abord.")
            return

        season = int(self.query_one("#season", Select).value or "1")
        lang = str(self.query_one("#lang", Select).value or "vostfr")
        self.plan.season = season
        self.plan.lang = lang

        full_url = f"{self.plan.base_url.rstrip('/')}/saison{season}/{lang}/"

        self._set_status("Chargement des épisodes…")
        try:
            episodes = fetch_episodes(full_url)
        except Exception as e:
            self._set_status(f"Erreur episodes: {e}")
            return

        if not episodes:
            self._set_status("Impossible de récupérer les épisodes (URL invalide ou blocage).")
            return

        self.episodes_data = episodes

        # Populate player select
        options = [(k, k) for k in episodes.keys()]
        player_select = self.query_one("#player", Select)
        player_select.set_options(options)

        best_player, _ = select_best_player(episodes)
        if best_player:
            player_select.value = best_player
            self.plan.player = best_player

        self._set_status(f"Épisodes chargés. Players: {', '.join(episodes.keys())}")

    async def _do_download(self) -> None:
        if not self.plan.base_url or not self.episodes_data:
            self._set_status("Charge les épisodes avant de télécharger.")
            return

        player = self.query_one("#player", Select).value
        if not player:
            self._set_status("Sélectionne un player.")
            return

        episodes_input = self.query_one("#episodes", Input).value.strip().lower() or "1"
        directory = self.query_one("#directory", Input).value.strip() or "./videos"
        directory = os.path.abspath(os.path.expanduser(directory))
        os.makedirs(directory, exist_ok=True)

        urls = self.episodes_data.get(player) or []
        if not urls:
            self._set_status("Player vide.")
            return

        # Parse episodes selection
        if episodes_input == "all":
            indices = list(range(len(urls)))
        else:
            indices: list[int] = []
            for part in episodes_input.split(","):
                part = part.strip()
                if not part:
                    continue
                if "-" in part:
                    a, b = part.split("-", 1)
                    start = int(a.strip())
                    end = int(b.strip())
                    if start > end:
                        start, end = end, start
                    indices.extend(range(start - 1, end))
                else:
                    indices.append(int(part) - 1)
            indices = sorted({i for i in indices if 0 <= i < len(urls)})

        if not indices:
            self._set_status("Aucun épisode valide.")
            return

        selected_urls = [urls[i] for i in indices]
        episode_numbers = [i + 1 for i in indices]

        self._set_status("Résolution des sources vidéo…")
        try:
            sources = fetch_video_source(selected_urls)
        except Exception as e:
            self._set_status(f"Erreur source: {e}")
            return

        if isinstance(sources, str):
            sources = [sources]

        if not sources or len(sources) != len(selected_urls):
            self._set_status("Impossible de résoudre toutes les sources.")
            return

        # Download sequentially (simple MVP)
        for ep_num, page_url, src in zip(episode_numbers, selected_urls, sources):
            safe_anime = (self.plan.base_url.rstrip("/").split("/")[-1] or "episode")
            save_path = os.path.join(directory, f"{safe_anime}_episode_{ep_num}.mp4")
            self._set_status(f"Téléchargement épisode {ep_num}…")
            ok, _ = download_video(src, save_path, use_ts_threading=True, url=page_url, automatic_mp4=False)
            if not ok:
                self._set_status(f"Échec épisode {ep_num}.")
                return

        self._set_status("Terminé.")


def run_tui() -> int:
    app = AnimeSamaTUI()
    app.run()
    return 0
