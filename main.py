import re
import os
import sys
import subprocess
import shutil
import platform
import urllib.request
import zipfile
import tarfile
import argparse
import threading
from utils.var import Colors, print_status, print_separator, print_header, print_tutorial
from utils.download_manager import DownloadJob, DownloadManager
from utils.config import get_default_download_path, get_max_concurrent_downloads, get_site_base_url_override
from utils.output_paths import build_episode_output_path

# PLEASE DO NOT REMOVE: Original code from https://github.com/sertrafurr/Anime-Sama-Downloader

def parse_arguments():
    """Parse command line arguments for CLI usage."""
    parser = argparse.ArgumentParser(
        description='🎌 Anime-Sama Video Downloader - Download anime episodes easily',
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
Examples:
  # Interactive mode (default)
  python main.py
  
  # CLI mode with URL
        python main.py -u "https://anime-sama.si/catalogue/sword-art-online/saison1/vostfr/" -e 1-5
  
  # Download specific episodes with threading
  python main.py -u "URL" -e 3,5,7 -t -d ~/Downloads
  
  # Download all episodes
  python main.py -u "URL" -e all --auto-mp4
        """
    )
    
    parser.add_argument('-s', '--search', type=str, action='append', help='Search anime by name (repeatable)')
    parser.add_argument('--config', type=str, default=None, help='Chemin vers un config.ini (sinon auto-detect)')
    parser.add_argument('--search-provider', type=str, choices=['anilist', 'local'], default='anilist', help='Search provider for --search (default: anilist)')
    parser.add_argument('--season', type=int, help='Season number when using --search (default: 1)', default=1)
    parser.add_argument('--lang', type=str, choices=['vostfr', 'vf', 'vo'], help='Language when using --search (default: vostfr)', default='vostfr')
    parser.add_argument('-u', '--url', type=str, action='append', help='Anime-Sama URL (repeatable)')
    parser.add_argument('-e', '--episodes', type=str, help='Episodes to download (e.g., "1-5", "3,5,7", "all")')
    parser.add_argument('-p', '--player', type=int, help='Player number to use (e.g., 1, 2, 3). If omitted, auto-select best player.', default=None)
    parser.add_argument('-d', '--directory', type=str, help='Dossier racine de sortie (défaut: config.ini ou smart default)', default=None)
    parser.add_argument('-t', '--threaded', action='store_true', help='Use threaded downloads (faster)')
    parser.add_argument('--ts-threaded', action='store_true', help='Use threaded .ts segment downloads (much faster for M3U8)')
    parser.add_argument('--mp4-threaded', action='store_true', help='Use multi-part Range downloads for MP4 when supported (can be faster)')
    parser.add_argument('--auto-mp4', action='store_true', help='Automatically convert .ts files to .mp4')
    parser.add_argument('--ffmpeg', action='store_true', help='Use ffmpeg for conversion (default, faster)')
    parser.add_argument('--moviepy', action='store_true', help='Use moviepy for conversion (slower but lighter)')
    parser.add_argument('--no-tutorial', action='store_true', help='Skip tutorial prompt')
    parser.add_argument('--quick', action='store_true', help='Quick mode: use smart defaults, minimal prompts')
    parser.add_argument('-j', '--jobs', type=int, default=None, help='Max téléchargements en parallèle (1-10). Défaut: config.ini')
    parser.add_argument('-y', '--yes', action='store_true', help='Ne pas demander de confirmation (utile pour --search)')
    parser.add_argument('--tui', action='store_true', help='Launch modern terminal UI (Textual). CLI remains default.')
    parser.add_argument('--ui', type=str, default=None, help='Launch a UI plugin by name (e.g. "textual"). Overrides --tui.')
    parser.add_argument('--version', action='version', version='Anime-Sama Downloader v2.6.1')
    
    return parser.parse_args()

def check_ffmpeg_installed():
    """Check if ffmpeg is installed (silent check)."""
    try:
        subprocess.run(["ffmpeg", "-version"], check=True, capture_output=True, text=True)
        return True
    except (subprocess.CalledProcessError, FileNotFoundError):
        return False


def install_ffmpeg_with_winget():
    system = platform.system().lower()
    if system != "windows":
        raise OSError("This script only supports Windows for now. Check https://ffmpeg.org/download.html for installation instructions for your OS.")

    print("Installing FFmpeg using winget...")
    try:
        subprocess.run(["winget", "install", "ffmpeg", "--accept-source-agreements", "--accept-package-agreements"], check=True)
        print("FFmpeg installed successfully via winget.")
    except subprocess.CalledProcessError as e:
        raise RuntimeError(f"Failed to install FFmpeg via winget: {e}")
        
def package_check(ask_install=False, first_run=False):
    missing_packages = []

    try:
        import requests
    except ImportError:
        missing_packages.append("requests")

    try:
        from tqdm import tqdm
    except ImportError:
        missing_packages.append("tqdm")

    try:
        from moviepy import VideoFileClip
    except ImportError:
        missing_packages.append("moviepy")

    try:
        from bs4 import BeautifulSoup
    except ImportError:
        missing_packages.append("beautifulsoup4")

    if missing_packages and ask_install:
        print("Missing packages:", ", ".join(missing_packages))
        if not first_run:
            for package in missing_packages:
                try:
                    print_status(f"Installing {package}...", "info")
                    subprocess.check_call([sys.executable, "-m", "pip", "install", package])
                except subprocess.CalledProcessError:
                    print_status(f"Failed to install {package}.", "error")
                    return False
            missing_packages = []
            try:
                import requests
            except ImportError:
                missing_packages.append("requests")
            try:
                from tqdm import tqdm
            except ImportError:
                missing_packages.append("tqdm")
            try:
                from moviepy import VideoFileClip
            except ImportError:
                missing_packages.append("moviepy")

            try:
                from bs4 import BeautifulSoup
            except ImportError:
                missing_packages.append("beautifulsoup4")
            
            if missing_packages:
                print_status(f"Some packages still missing after installation: {', '.join(missing_packages)}", "error")
                return False
        else:
            return False
    return len(missing_packages) == 0

if not package_check(ask_install=True, first_run=True):
    print_status("Some required packages were missing. Would you like to install them now? (y/n): ", "warning")
    if not sys.stdin.isatty():
        print_status(
            "Cannot prompt for installation in non-interactive mode. Install deps and re-run: pip install -r requirements.txt",
            "error",
        )
        sys.exit(1)

    ask_user = input().strip().lower()
    if ask_user in ['y', 'yes', '1']:
        if not package_check(ask_install=True, first_run=False):
            print_status("Failed to install required packages. Please install them manually and re-run the script. pip install -r requirements.txt", "error")
            sys.exit(1)
    else:
        print_status("Cannot proceed without required packages. Exiting.", "warning")
        input("Press Enter to exit...")
        sys.exit(1)

try:
    from moviepy                  import VideoFileClip
except ImportError:
    print_status("moviepy not installed and can't seem to be installed. You should install it manually.", "error")
    VideoFileClip = None

from concurrent.futures           import ThreadPoolExecutor, as_completed

from utils.parsers                import parse_ts_segments
from utils.ts_to_mp4              import convert_ts_to_mp4
from utils.fetch                  import fetch_episodes, fetch_video_source
from utils.downloaders.downloader import download_video
from utils.stuff                  import print_episodes, get_player_choice, get_episode_choice
def extract_anime_name(base_url):
    match = re.search(r'catalogue/([^/]+)/', base_url)
    if match:
        return match.group(1)
    return "episode"

def install_package(package_name):
    try:
        subprocess.check_call([sys.executable, "-m", "pip", "install", package_name])
        return True
    except subprocess.CalledProcessError:
        return False

def check_ffmpeg_installed():
    return shutil.which("ffmpeg") is not None

def get_save_directory():
    print(f"\n{Colors.BOLD}{Colors.HEADER}📁 SAVE LOCATION{Colors.ENDC}")
    print_separator()
    
    default_dir = get_default_download_path()
    save_dir = input(f"{Colors.OKCYAN}Enter directory to save videos (default: {default_dir}): {Colors.ENDC}").strip()
    
    if not save_dir:
        save_dir = default_dir
    
    # Expand user path (~) to absolute path
    save_dir = os.path.expanduser(save_dir)
    save_dir = os.path.abspath(save_dir)
    
    try:
        os.makedirs(save_dir, exist_ok=True)
        print_status(f"Save directory set to: {save_dir}", "success")
        return save_dir
    except Exception as e:
        print_status(f"Cannot create directory {save_dir}: {str(e)}", "error")
        print_status(f"Using default directory: {default_dir}", "info")
        default_dir = os.path.abspath(default_dir)
        os.makedirs(default_dir, exist_ok=True)
        return default_dir

def validate_anime_sama_url(url):
    from urllib.parse import urlparse

    parsed = urlparse(url or "")
    if parsed.scheme not in {"http", "https"}:
        return False, "URL invalide (http/https requis)."

    # Allow default anime-sama.* domains + optional override from config/env.
    override = get_site_base_url_override() or ""
    override_netloc = urlparse(override).netloc if override else ""
    allowed = {"anime-sama.tv", "anime-sama.fr", "anime-sama.org", "anime-sama.si"}
    if override_netloc:
        allowed.add(override_netloc)

    if (parsed.netloc or "").lower() not in {d.lower() for d in allowed}:
        return False, (
            "Domaine non autorisé. Configure [SITE] base_url/domain dans config.ini si le domaine a changé."
        )

    path_ok = re.match(r"^/catalogue/[^/]+/saison\d+/(?:vostfr|vf|vo)/?$", parsed.path or "", re.IGNORECASE)
    if not path_ok:
        return False, (
            "URL invalide. Format attendu:\n"
            "  https://<domaine>/catalogue/<anime-name>/saison<NUMBER>/<language>/\n"
            "Où <language> ∈ {vostfr, vf, vo}."
        )
    return True, ""


def download_episode(
    episode_num,
    url,
    video_source,
    anime_name,
    save_dir,
    use_ts_threading=False,
    use_mp4_threading=False,
    automatic_mp4=False,
    pre_selected_tool=None,
):
    if not video_source:
        print_status(f"Could not extract video source for episode {episode_num}", "error")
        return False, None
    
    print_separator()
    print_status(f"Processing episode: {episode_num}", "info")
    print_status(f"Source: {url[:60]}...", "info")
    
    # New naming/layout: <root>/<anime>/Saison <n>/<lang>/<anime>-S<n>E<ep>.mp4
    season = 1
    lang = "vostfr"
    m = re.search(r"/saison(?P<s>\d+)/(?:vostfr|vf|vo)/", url or "")
    if m:
        try:
            season = int(m.group("s"))
        except Exception:
            season = 1
    m2 = re.search(r"/saison\d+/(?P<lang>vostfr|vf|vo)/", url or "", re.IGNORECASE)
    if m2:
        lang = str(m2.group("lang")).lower()

    dest_dir, save_path = build_episode_output_path(save_dir, anime_name, season, lang, episode_num, ext="mp4")
    os.makedirs(dest_dir, exist_ok=True)
    
    print(f"\n{Colors.BOLD}{Colors.HEADER}⬇️ DOWNLOADING EPISODE {episode_num}{Colors.ENDC}")
    print_separator()
    
    try:
        success, output_path = download_video(
            video_source,
            save_path,
            use_ts_threading=use_ts_threading,
            use_mp4_threading=use_mp4_threading,
            url=url,
            automatic_mp4=automatic_mp4,
        )
    except Exception as e:
        print_status(f"Download failed for episode {episode_num}: {str(e)}", "error")
        return False, None
    
    if not success:
        print_status(f"Failed to download episode {episode_num}", "error")
        return False, None
    
    print_separator()
    
    if 'm3u8' in video_source and output_path.endswith('.ts'):
        print_status(f"Video saved as {output_path} (MPEG-TS format, playable in VLC or similar players)", "success")
        if automatic_mp4:
            success, final_path = convert_ts_to_mp4(output_path, save_path, pre_selected_tool)
            if success:
                print_status(f"Episode {episode_num} successfully saved to: {final_path}", "success")
                try:
                    os.remove(output_path)
                    print_status(f"Removed temporary .ts file: {output_path}", "info")
                except Exception as e:
                    print_status(f"Could not remove temporary .ts file: {str(e)}", "warning")
                return True, final_path
            else:
                print_status(f"Conversion failed for episode {episode_num}, keeping .ts file: {output_path}", "error")
                return False, output_path
        else:
            print_status(f"Keeping .ts file for episode {episode_num}: {output_path}", "info")
            return True, output_path
    else:
        print_status(f"Episode {episode_num} successfully saved to: {save_path}", "success")
        try:
            from utils.media_refresh import schedule_media_refresh

            schedule_media_refresh()
        except Exception:
            pass
        return True, save_path

def complete_anime_url(base_url):
    """
    Complete a base anime URL with season and language selection.
    Takes: https://anime-sama.tv/catalogue/sword-art-online/
    Returns: https://anime-sama.tv/catalogue/sword-art-online/saison1/vostfr/
    """
    from utils.var import Colors
    
    # Remove trailing slash
    base_url = base_url.rstrip('/')
    
    print(f"\n{Colors.BOLD}{Colors.HEADER}📺 SEASON & LANGUAGE SELECTION{Colors.ENDC}")
    print_separator()
    
    # Ask for season
    while True:
        season = input(f"{Colors.BOLD}Enter season number (default: 1): {Colors.ENDC}").strip()
        if not season:
            season = "1"
        
        if season.isdigit():
            break
        else:
            print_status("Please enter a valid number", "error")
    
    # Ask for language
    print(f"\n{Colors.BOLD}Available languages:{Colors.ENDC}")
    print("  1. VOSTFR (Japanese with French subtitles)")
    print("  2. VF (French dub)")
    print("  3. VO (Original Japanese, no subtitles)")
    
    while True:
        lang_choice = input(f"\n{Colors.BOLD}Select language (1-3, default: 1): {Colors.ENDC}").strip()
        if not lang_choice:
            lang_choice = "1"
        
        if lang_choice == "1":
            language = "vostfr"
            break
        elif lang_choice == "2":
            language = "vf"
            break
        elif lang_choice == "3":
            language = "vo"
            break
        else:
            print_status("Please enter 1, 2, or 3", "error")
    
    # Build complete URL
    complete_url = f"{base_url}/saison{season}/{language}/"
    print_status(f"Complete URL: {complete_url}", "info")
    
    return complete_url

def parse_episode_range(episode_str, max_episodes):
    """Parse episode string like '1-5', '3,5,7', or 'all' and return list of indices."""
    if episode_str.lower() == 'all':
        return list(range(max_episodes))
    
    indices = []
    parts = episode_str.split(',')
    
    for part in parts:
        part = part.strip()
        if '-' in part:
            # Range like "1-5"
            start, end = map(str.strip, part.split('-'))
            start_idx = int(start) - 1  # Convert to 0-indexed
            end_idx = int(end)  # Inclusive end
            indices.extend(range(start_idx, end_idx))
        else:
            # Single episode like "3"
            indices.append(int(part) - 1)  # Convert to 0-indexed
    
    # Remove duplicates and sort
    indices = sorted(set(indices))
    
    # Validate indices
    valid_indices = [i for i in indices if 0 <= i < max_episodes]
    
    if len(valid_indices) != len(indices):
        print_status(f"Warning: Some episode numbers are out of range (max: {max_episodes})", "warning")
    
    return valid_indices


def _clamp_jobs(n: int) -> int:
    try:
        n = int(n)
    except Exception:
        return 1
    return max(1, min(10, n))


def _ensure_list(value):
    if value is None:
        return []
    if isinstance(value, list):
        return [v for v in value if v]
    return [value]


def _slug_from_catalogue_url(url: str) -> str:
    m = re.search(r"/catalogue/([^/]+)/", url)
    return m.group(1) if m else "anime"


def run_batch_download(args, urls: list[str], searches: list[str]) -> int:
    """Mode batch/queue: plusieurs animes + jobs parallèles (max 10)."""
    from utils.fetch import select_best_player

    jobs = _clamp_jobs(getattr(args, 'jobs', 1))

    # Résoudre les searches en URLs catalogue
    targets: list[str] = []
    if searches:
        from utils.search import quick_search
        for q in searches:
            print_status(f"Searching for: {q}", "loading")
            found = quick_search(q, provider=args.search_provider)
            if not found:
                print_status(f"No good match found for '{q}'", "error")
                continue

            found = found.rstrip('/')
            base_url = f"{found}/saison{args.season}/{args.lang}/"
            print_status(f"Found match: {base_url}", "success")

            if (not args.yes) and sys.stdin.isatty():
                ans = input(f"Télécharger '{q}' → {base_url} ? (Y/n): ").strip().lower()
                if ans == 'n':
                    print_status("Skipped.", "info")
                    continue

            targets.append(base_url)

    targets.extend(urls)

    if not targets:
        print_status("Aucune URL à traiter.", "error")
        return 1

    # Répertoire racine
    save_root = os.path.abspath(os.path.expanduser(args.directory or './videos'))
    os.makedirs(save_root, exist_ok=True)
    print_status(f"Save directory: {save_root}", "info")

    # Canonicalize domain (anime-sama.tv → .si)
    try:
        from utils.search import canonicalize_site_url
        targets = [canonicalize_site_url(u) for u in targets]
    except Exception:
        pass

    # Construire et enqueuer les jobs
    mgr = DownloadManager(max_parallel=jobs, executor_name="cli-dl")

    # On limite la concurrence interne quand jobs est élevé.
    mp4_workers = 2 if jobs >= 5 else 4
    ts_workers = 2 if jobs >= 5 else 10
    use_ts_threading = bool(getattr(args, 'ts_threaded', False)) and jobs <= 3
    use_mp4_threading = bool(getattr(args, 'mp4_threaded', False))
    automatic_mp4 = bool(getattr(args, 'auto_mp4', False))

    enqueued = 0
    try:
        for base_url in targets:
            is_valid, error_msg = validate_anime_sama_url(base_url)
            if not is_valid:
                print_status(error_msg, "error")
                continue

            # Derive season/lang from the base_url (for output naming/layout).
            season = 1
            lang = "vostfr"
            m = re.search(r"/saison(?P<s>\d+)/(?:vostfr|vf|vo)/", base_url or "")
            if m:
                try:
                    season = int(m.group("s"))
                except Exception:
                    season = 1
            m2 = re.search(r"/saison\d+/(?P<lang>vostfr|vf|vo)/", base_url or "", re.IGNORECASE)
            if m2:
                lang = str(m2.group("lang")).lower()

            slug = _slug_from_catalogue_url(base_url)
            print_status(f"Fetching episodes for: {slug}", "loading")
            episodes = fetch_episodes(base_url)
            if not episodes:
                print_status(f"Failed to fetch episodes for {base_url}", "error")
                continue

            player_choice, _ = select_best_player(episodes)
            if not player_choice:
                print_status(f"No valid player found for {slug}", "error")
                continue

            max_eps = len(episodes[player_choice])
            wanted = args.episodes or 'all'
            episode_indices = parse_episode_range(wanted, max_eps)
            if not episode_indices:
                print_status(f"No valid episodes selected for {slug}", "error")
                continue

            page_urls = [episodes[player_choice][i] for i in episode_indices]
            episode_numbers = [i + 1 for i in episode_indices]

            for ep_num, page_url in zip(episode_numbers, page_urls):
                if not page_url:
                    continue

                out_dir, save_path = build_episode_output_path(save_root, slug, season, lang, ep_num, ext="mp4")
                os.makedirs(out_dir, exist_ok=True)
                label = f"{slug} S{season}E{ep_num}"

                def _make_runner(page_url=page_url, save_path=save_path, label=label):
                    def _runner(cancel_event: threading.Event):
                        if cancel_event.is_set():
                            return None

                        try:
                            src = fetch_video_source(page_url)
                        except Exception as e:
                            print_status(f"{label}: source error: {e}", "error")
                            return None

                        if not src:
                            print_status(f"{label}: no source", "error")
                            return None

                        def _cb(message: str, status_type: str = "info") -> None:
                            # Sortie CLI simple (peut s'entremêler en parallèle)
                            print_status(f"{label}: {message}", status_type if status_type else "info")

                        ok, out = download_video(
                            src,
                            save_path,
                            use_ts_threading=use_ts_threading,
                            url=page_url,
                            automatic_mp4=automatic_mp4,
                            log_callback=_cb,
                            use_tqdm=False,
                            use_mp4_threading=use_mp4_threading,
                            mp4_workers=mp4_workers,
                            ts_workers=ts_workers,
                            cancel_event=cancel_event,
                        )
                        return out if ok else None

                    return _runner

                mgr.enqueue(DownloadJob(label=label, run=_make_runner()))
                enqueued += 1

        if not enqueued:
            print_status("Aucun job n'a été ajouté à la file.", "error")
            return 1

        print_status(f"Queue démarrée: {enqueued} job(s), {jobs} en parallèle.", "success")
        mgr.wait()
        try:
            from utils.media_refresh import flush_media_refresh, shutdown_media_refresh

            flush_media_refresh(timeout_s=15)
            shutdown_media_refresh(timeout_s=2)
        except Exception:
            pass
        return 0

    except KeyboardInterrupt:
        print_status("\nAnnulation demandée (Ctrl+C)…", "warning")
        mgr.cancel_all()
        mgr.shutdown()
        return 1


def main():
    args = parse_arguments()

    # If user passed --config, expose it to config loader via env.
    if getattr(args, "config", None):
        os.environ["ASD_CONFIG"] = str(args.config)

    # Optional UI plugin (keeps CLI as default behavior)
    ui_name = (getattr(args, 'ui', None) or '').strip()
    if ui_name or getattr(args, 'tui', False):
        if not ui_name:
            ui_name = 'textual'
        try:
            from utils.ui.registry import run_ui
            return int(run_ui(ui_name))
        except Exception as e:
            print_status(f"UI failed to start: {e}", "error")
            print_status("Install dependencies: pip install -r requirements.txt", "info")
            return 1
    
    urls = _ensure_list(getattr(args, 'url', None))
    searches = _ensure_list(getattr(args, 'search', None))
    jobs = _clamp_jobs(getattr(args, 'jobs', None) or get_max_concurrent_downloads(default=1))

    # Batch mode: multiple animes or explicit parallel jobs
    if (len(urls) + len(searches) > 1) or (jobs > 1):
        return run_batch_download(args, urls, searches)

    # Check if running in CLI mode or interactive mode
    cli_mode = bool(urls) or bool(searches)

    # Resolve default output root from config if user didn't pass --directory.
    if not getattr(args, "directory", None):
        args.directory = get_default_download_path()

    try:
        # Print header (skip in quick mode for cleaner output)
        if not cli_mode and not args.quick:
            print_header()
        
        # Get URL (from args, search, or interactive input)
        if cli_mode and searches:
            # CLI search mode
            from utils.search import quick_search
            print_status(f"Searching for: {searches[0]}", "loading")
            search_result = quick_search(searches[0], provider=args.search_provider)
            
            if not search_result:
                print_status(f"No good match found for '{searches[0]}'", "error")
                print_status("Try being more specific or use interactive search", "info")
                return 1
            
            print_status(f"Found match: {search_result}", "success")
            
            # Complete URL with season and language from args
            search_result = search_result.rstrip('/')
            base_url = f"{search_result}/saison{args.season}/{args.lang}/"
            print_status(f"Complete URL: {base_url}", "info")

            if (not getattr(args, 'yes', False)) and sys.stdin.isatty():
                ans = input(f"Confirmer le téléchargement: {base_url} ? (Y/n): ").strip().lower()
                if ans == 'n':
                    print_status("Annulé.", "warning")
                    return 1
            
        elif cli_mode:
            base_url = urls[0]
            print_status(f"Using URL from arguments: {base_url}", "info")
        else:
            # Interactive mode: always start with search
            from utils.search import interactive_search
            
            search_url = interactive_search()
            if search_url:
                # Complete the URL with season and language
                base_url = complete_anime_url(search_url)
            else:
                # User cancelled search, fall back to manual URL
                print(f"\n{Colors.BOLD}{Colors.HEADER}🔗 MANUAL URL INPUT{Colors.ENDC}")
                print_separator()
                base_url = None
            
            # If no URL from search, ask for manual input
            if not base_url:
                while True:
                    base_url = input(f"{Colors.BOLD}Enter the complete anime-sama URL: {Colors.ENDC}").strip()
                    
                    if not base_url:
                        print_status("URL cannot be empty", "error")
                        continue
                        
                    is_valid, error_msg = validate_anime_sama_url(base_url)
                    if not is_valid:
                        print_status(error_msg, "error")
                        print_status("Example: https://anime-sama.si/catalogue/roshidere/saison1/vostfr/", "info")
                        continue
                    
                    break
        
        # Canonicalize domain (anime-sama.tv currently redirects to .si without preserving path)
        try:
            from utils.search import canonicalize_site_url
            if base_url:
                base_url = canonicalize_site_url(base_url)
        except Exception:
            pass

        # Validate URL
        is_valid, error_msg = validate_anime_sama_url(base_url)
        if not is_valid:
            print_status(error_msg, "error")
            return 1
        
        anime_name = extract_anime_name(base_url)
        print_status(f"Detected anime: {anime_name}", "info")
        
        episodes = fetch_episodes(base_url)
        if not episodes:
            print_status("Failed to fetch episodes. Please check the URL and try again.", "error")
            return 1
        
        if not cli_mode:
            print_episodes(episodes)
        
        # Get player choice (from args or interactive)
        # Auto-select best player or use argument
        if cli_mode and args.player is not None:
            player_name = f"Player {args.player}"
            if player_name not in episodes:
                print_status(f"Player {args.player} not found. Available players: {', '.join(episodes.keys())}", "error")
                return 1
            player_choice = player_name
            print_status(f"Using {player_choice} (from arguments)", "info")
        else:
            # Auto-select best player
            from utils.fetch import select_best_player
            player_choice, _ = select_best_player(episodes)
            if not player_choice:
                print_status("No valid player found", "error")
                return 1
            
            available_eps = sum(1 for url in episodes[player_choice] if url and url.strip())
            print_status(f"Auto-selected {player_choice} ({available_eps} episodes available)", "success")
        
        # Get episode selection (from args or interactive)
        if cli_mode and args.episodes:
            max_eps = len(episodes[player_choice])
            episode_indices = parse_episode_range(args.episodes, max_eps)
            if not episode_indices:
                print_status("No valid episodes selected", "error")
                return 1
            print_status(f"Selected episodes: {', '.join(str(i+1) for i in episode_indices)}", "info")
        else:
            episode_indices = get_episode_choice(episodes, player_choice)
            if episode_indices is None:
                return 1
        
        # Get save directory (from args or interactive with favorites)
        if cli_mode:
            save_dir = os.path.expanduser(args.directory)
            save_dir = os.path.abspath(save_dir)
            os.makedirs(save_dir, exist_ok=True)
            print_status(f"Save directory: {save_dir}", "info")
        else:
            from utils.path_input import get_save_directory_interactive
            save_dir = get_save_directory_interactive()
        
        if isinstance(episode_indices, int):
            episode_indices = [episode_indices]
        
        urls = [episodes[player_choice][index] for index in episode_indices]
        episode_numbers = [index + 1 for index in episode_indices]
        
        print(f"\n{Colors.BOLD}{Colors.HEADER}🎬 PROCESSING EPISODES{Colors.ENDC}")
        print_separator()
        print_status(f"Player: {player_choice}", "info")
        print_status(f"Episodes selected: {', '.join(map(str, episode_numbers))}", "info")
        
        video_sources = fetch_video_source(urls)
        if not video_sources:
            print_status("Could not extract video sources for selected episodes", "error")
            return 1
        
        if isinstance(video_sources, str):
            video_sources = [video_sources]
        
        # Threading options (from args or interactive)
        if cli_mode:
            use_threading = args.threaded
            use_ts_threading = args.ts_threaded
            use_mp4_threading = args.mp4_threaded
            automatic_mp4 = args.auto_mp4
            pre_selected_tool = None
            
            if args.moviepy:
                pre_selected_tool = 'moviepy'
            elif args.ffmpeg or automatic_mp4:
                pre_selected_tool = 'ffmpeg'
                if not check_ffmpeg_installed():
                    print_status("ffmpeg is not installed. Fallback to moviepy", "warning")
                    pre_selected_tool = 'moviepy'
            
            if use_threading:
                print_status("Using threaded episode downloads", "info")
            if use_ts_threading:
                print_status("Using threaded .ts segment downloads", "info")
            if use_mp4_threading:
                print_status("Using multi-part MP4 downloads (Range)", "info")
            if automatic_mp4:
                print_status(f"Auto-converting to MP4 using {pre_selected_tool}", "info")
        else:
            # Interactive mode with smart defaults
            use_threading = False
            use_ts_threading = False
            use_mp4_threading = False
            automatic_mp4 = False
            pre_selected_tool = None
            
            # Check if quick mode or has M3U8 sources
            has_m3u8 = any('m3u8' in src for src in video_sources if src)
            quick_mode = args.quick if hasattr(args, 'quick') else False
            
            if quick_mode:
                # Quick mode: use best defaults, no questions!
                print_status("Quick mode: Using optimal defaults", "info")
                use_threading = len(episode_indices) > 1  # Auto-thread if multiple episodes
                use_ts_threading = has_m3u8  # Auto-thread .ts if M3U8
                use_mp4_threading = True  # Try to speed up direct MP4 when server supports Range
                automatic_mp4 = has_m3u8  # Auto-convert M3U8
                pre_selected_tool = 'ffmpeg' if check_ffmpeg_installed() else 'moviepy'
                
                if use_threading:
                    print_status("✓ Multi-episode threading enabled", "success")
                if use_ts_threading:
                    print_status("✓ Fast .ts segment downloads enabled", "success")
                if automatic_mp4:
                    print_status(f"✓ Auto MP4 conversion with {pre_selected_tool}", "success")
            else:
                # Standard interactive mode with fewer questions
                if len(episode_indices) > 1:
                    thread_choice = input(f"{Colors.BOLD}Use fast multi-episode download? (Y/n, default: Y): {Colors.ENDC}").strip().lower()
                    use_threading = thread_choice != 'n'  # Default to YES

                if has_m3u8:
                    # Ask once for both threading and conversion
                    print(f"\n{Colors.BOLD}{Colors.OKCYAN}M3U8 sources detected - Recommended settings:{Colors.ENDC}")
                    print(f"  • Fast .ts downloads (10x faster)")
                    print(f"  • Auto MP4 conversion")
                    
                    optimize_choice = input(f"{Colors.BOLD}Use recommended settings? (Y/n, default: Y): {Colors.ENDC}").strip().lower()
                    
                    if optimize_choice != 'n':
                        # Use recommended settings
                        use_ts_threading = True
                        automatic_mp4 = True
                        pre_selected_tool = 'ffmpeg' if check_ffmpeg_installed() else 'moviepy'
                        print_status(f"✓ Using optimized settings with {pre_selected_tool}", "success")
                    else:
                        # Manual choices
                        ts_thread_choice = input(f"{Colors.BOLD}Fast .ts downloads? (Y/n, default: Y): {Colors.ENDC}").strip().lower()
                        use_ts_threading = ts_thread_choice != 'n'
                        
                        auto_mp4_choice = input(f"{Colors.BOLD}Auto-convert to MP4? (Y/n, default: Y): {Colors.ENDC}").strip().lower()
                        automatic_mp4 = auto_mp4_choice != 'n'
                        
                        if automatic_mp4:
                            if check_ffmpeg_installed():
                                pre_selected_tool = 'ffmpeg'
                                print_status("Using ffmpeg for conversion", "info")
                            else:
                                pre_selected_tool = 'moviepy'
                                print_status("Using moviepy for conversion", "info")

        failed_downloads = 0
        try:
            if use_threading and len(episode_indices) > 1:
                print_status("Starting threaded downloads...", "info")
                with ThreadPoolExecutor() as executor:
                    future_to_episode = {
                        executor.submit(
                            download_episode,
                            ep_num,
                            url,
                            video_src,
                            anime_name,
                            save_dir,
                            use_ts_threading,
                            use_mp4_threading,
                            automatic_mp4,
                            pre_selected_tool,
                        ): ep_num
                        for ep_num, url, video_src in zip(episode_numbers, urls, video_sources)
                    }
                    for future in as_completed(future_to_episode):
                        ep_num = future_to_episode[future]
                        try:
                            success, _ = future.result()
                            if not success:
                                failed_downloads += 1
                        except Exception as e:
                            print_status(f"Episode {ep_num} generated an exception: {str(e)}", "error")
                            failed_downloads += 1
            else:
                for episode_num, url, video_source in zip(episode_numbers, urls, video_sources):
                    success, _ = download_episode(
                        episode_num,
                        url,
                        video_source,
                        anime_name,
                        save_dir,
                        use_ts_threading,
                        use_mp4_threading,
                        automatic_mp4,
                        pre_selected_tool,
                    )
                    if not success:
                        failed_downloads += 1

            print_separator()
            if failed_downloads == 0:
                print_status("All downloads completed successfully! Enjoy watching! 🎉", "success")
                try:
                    from utils.media_refresh import flush_media_refresh, shutdown_media_refresh

                    flush_media_refresh(timeout_s=15)
                    shutdown_media_refresh(timeout_s=2)
                except Exception:
                    pass
                if not cli_mode:
                    input(f"{Colors.BOLD}Press Enter to exit...{Colors.ENDC}")
                return 0
            else:
                print_status(f"Completed with {failed_downloads} failed downloads", "warning")
                try:
                    from utils.media_refresh import flush_media_refresh, shutdown_media_refresh

                    flush_media_refresh(timeout_s=15)
                    shutdown_media_refresh(timeout_s=2)
                except Exception:
                    pass
                if not cli_mode:
                    input(f"{Colors.BOLD}Press Enter to exit...{Colors.ENDC}")
                return 1

        except KeyboardInterrupt:
            print_status("\n\nProgram interrupted by user", "error")
            return 1
        except Exception as e:
            print_status(f"Unexpected error: {str(e)}", "error")
            return 1
    except Exception as e:
        print_status(f"Fatal error: {str(e)}", "error")
        return 1
if __name__ == "__main__":

    sys.exit(main())
