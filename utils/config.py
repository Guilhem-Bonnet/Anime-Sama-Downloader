"""Configuration management for favorite paths and settings.

This project historically used a JSON file in the user's home directory.
We keep that for backward compatibility (favorites/last used path) and add
optional INI-based configuration for global settings (download root, site URL
override, etc.).
"""

import configparser
import json
import os
import shutil
from pathlib import Path

CONFIG_FILE = os.path.expanduser("~/.anime-sama-downloader.json")


def _first_existing(paths: list[str]) -> str | None:
    for p in paths:
        if p and os.path.exists(p):
            return p
    return None


def find_config_ini(explicit_path: str | None = None) -> str | None:
    """Find an INI config file.

    Precedence:
    - explicit_path (if provided)
    - env ASD_CONFIG
    - ./config.ini (repo/local)
    - ~/.config/anime-sama-downloader/config.ini
    - ~/.anime-sama-downloader.ini
    """
    env_path = os.environ.get("ASD_CONFIG")
    home = os.path.expanduser("~")
    return _first_existing(
        [
            explicit_path,
            env_path,
            os.path.abspath("config.ini"),
            os.path.join(home, ".config", "anime-sama-downloader", "config.ini"),
            os.path.join(home, ".anime-sama-downloader.ini"),
        ]
    )


def load_ini_config(explicit_path: str | None = None) -> configparser.ConfigParser:
    cfg = configparser.ConfigParser()
    cfg.read_dict({"DEFAULT": {}})
    path = find_config_ini(explicit_path)
    if not path:
        return cfg
    try:
        cfg.read(path, encoding="utf-8")
    except Exception:
        # Never crash on config read errors.
        return configparser.ConfigParser()
    return cfg


def _ini_get(
    cfg: configparser.ConfigParser,
    section: str,
    key: str,
    default: str | None = None,
) -> str | None:
    try:
        if cfg.has_option(section, key):
            v = cfg.get(section, key)
            return v
        if cfg.has_option("DEFAULT", key):
            return cfg.get("DEFAULT", key)
    except Exception:
        return default
    return default


def _env_get(key: str) -> str | None:
    v = os.environ.get(key)
    if v is None:
        return None
    v = str(v).strip()
    return v or None


def get_site_base_url_override(config_path: str | None = None) -> str | None:
    """Override Anime-Sama base URL if DNS/domain changes.

    Supported:
    - env ASD_SITE_BASE_URL (e.g. https://anime-sama.si)
    - env ASD_SITE_DOMAIN (e.g. anime-sama.si)
    - INI [SITE] base_url or domain
    """
    base = _env_get("ASD_SITE_BASE_URL")
    if base:
        return base.rstrip("/")

    domain = _env_get("ASD_SITE_DOMAIN")
    if domain:
        domain = domain.replace("https://", "").replace("http://", "").strip("/")
        return f"https://{domain}" if domain else None

    ini = load_ini_config(config_path)
    base = _ini_get(ini, "SITE", "base_url")
    if base:
        return str(base).strip().rstrip("/")

    domain = _ini_get(ini, "SITE", "domain")
    if domain:
        domain = str(domain).replace("https://", "").replace("http://", "").strip().strip("/")
        return f"https://{domain}" if domain else None

    return None


def get_max_concurrent_downloads(config_path: str | None = None, default: int = 3) -> int:
    """Max téléchargements en parallèle (1-10)."""
    env_v = _env_get("ASD_MAX_CONCURRENT_DOWNLOADS")
    if env_v is not None:
        try:
            v = int(env_v)
            return max(1, min(v, 10))
        except ValueError:
            pass

    ini = load_ini_config(config_path)
    ini_v = _ini_get(ini, "DEFAULT", "max_concurrent_downloads")
    if ini_v is not None:
        try:
            v = int(str(ini_v).strip())
            return max(1, min(v, 10))
        except ValueError:
            pass

    try:
        return max(1, min(int(default), 10))
    except Exception:
        return 3


def get_web_bind(config_path: str | None = None) -> tuple[str, int]:
    """Retourne (host, port) pour l'UI web."""
    host = _env_get("ASD_WEB_HOST")
    port_s = _env_get("ASD_WEB_PORT")

    ini = load_ini_config(config_path)
    if not host:
        host = _ini_get(ini, "WEB", "host")
    if not port_s:
        port_s = _ini_get(ini, "WEB", "port")

    if not host:
        host = "127.0.0.1"

    port = 8000
    if port_s:
        try:
            port = int(str(port_s).strip())
        except ValueError:
            port = 8000

    return (str(host).strip() or "127.0.0.1", port)


def _parse_bool(v: str | None, default: bool = False) -> bool:
    if v is None:
        return bool(default)
    s = str(v).strip().lower()
    if s in {"1", "true", "yes", "y", "on"}:
        return True
    if s in {"0", "false", "no", "n", "off"}:
        return False
    return bool(default)


def get_output_naming_mode(config_path: str | None = None, default: str = "legacy") -> str:
    """Naming mode for output paths.

    Values (case-insensitive):
    - legacy (default): <root>/<slug>/Saison <n>/<lang>/<slug>-S<n>E<ep>.ext
    - media / jellyfin / plex: <root>/<Series>/Season 01/<Series> - S01E01.ext

    Sources:
    - env ASD_OUTPUT_NAMING_MODE or ASD_NAMING_MODE
    - INI [OUTPUT] naming_mode
    """
    env_v = _env_get("ASD_OUTPUT_NAMING_MODE") or _env_get("ASD_NAMING_MODE")
    if env_v:
        return str(env_v).strip().lower()

    ini = load_ini_config(config_path)
    ini_v = _ini_get(ini, "OUTPUT", "naming_mode")
    if ini_v:
        return str(ini_v).strip().lower()

    return str(default).strip().lower()


def get_media_separate_lang(config_path: str | None = None, default: bool = False) -> bool:
    """If True in media-server naming, create one series folder per language.

    Example: "My Show [VOSTFR]" and "My Show [VF]".

    Sources:
    - env ASD_MEDIA_SEPARATE_LANG
    - INI [OUTPUT] media_separate_lang
    """
    env_v = _env_get("ASD_MEDIA_SEPARATE_LANG")
    if env_v is not None:
        return _parse_bool(env_v, default=default)

    ini = load_ini_config(config_path)
    ini_v = _ini_get(ini, "OUTPUT", "media_separate_lang")
    if ini_v is not None:
        return _parse_bool(str(ini_v), default=default)

    return bool(default)


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


def _iter_linux_mountpoints() -> list[str]:
    mounts: list[str] = []
    try:
        with open("/proc/mounts", "r", encoding="utf-8") as f:
            for line in f:
                parts = line.split()
                if len(parts) < 3:
                    continue
                mountpoint = parts[1]
                fstype = parts[2]
                if fstype in {
                    "proc",
                    "sysfs",
                    "devtmpfs",
                    "tmpfs",
                    "cgroup",
                    "cgroup2",
                    "overlay",
                    "squashfs",
                    "autofs",
                    "fusectl",
                    "mqueue",
                    "debugfs",
                    "tracefs",
                    "securityfs",
                    "pstore",
                    "ramfs",
                }:
                    continue
                if mountpoint.startswith("/proc") or mountpoint.startswith("/sys") or mountpoint.startswith("/dev"):
                    continue
                mounts.append(mountpoint)
    except Exception:
        return []

    # Deduplicate while keeping order
    seen: set[str] = set()
    out: list[str] = []
    for m in mounts:
        if m in seen:
            continue
        seen.add(m)
        out.append(m)
    return out


def _smart_disk_root_linux() -> str | None:
    """Try to pick a writable mountpoint with lots of free space.

    Preference order:
    - /run/media, /media, /mnt (external/secondary disks)
    - fallback: any writable mountpoint
    """

    mounts = _iter_linux_mountpoints()
    if not mounts:
        return None

    def score_mount(mountpoint: str) -> tuple[int, int]:
        # prefer removable-ish locations
        preferred = int(
            mountpoint.startswith("/run/media")
            or mountpoint.startswith("/media")
            or mountpoint.startswith("/mnt")
        )
        try:
            usage = shutil.disk_usage(mountpoint)
            free = int(usage.free)
        except Exception:
            free = -1
        return (preferred, free)

    candidates = [m for m in mounts if os.path.isdir(m) and os.access(m, os.W_OK)]
    if not candidates:
        return None

    best = max(candidates, key=score_mount)
    if best and os.path.isdir(best):
        return best
    return None


def get_preferred_download_root() -> str | None:
    config = load_config()
    root = _config_get(config, "download", "root")
    if isinstance(root, str) and root.strip():
        return os.path.abspath(os.path.expanduser(root.strip()))
    return None


def set_preferred_download_root(path: str) -> None:
    config = load_config()
    _config_set(config, os.path.abspath(os.path.expanduser(path)), "download", "root")
    save_config(config)

def load_config():
    """Load configuration from file."""
    if os.path.exists(CONFIG_FILE):
        try:
            with open(CONFIG_FILE, 'r', encoding='utf-8') as f:
                return json.load(f)
        except Exception as e:
            print(f"Warning: Could not load config: {e}")
            return {}
    return {}

def save_config(config):
    """Save configuration to file."""
    try:
        with open(CONFIG_FILE, 'w', encoding='utf-8') as f:
            json.dump(config, f, indent=2, ensure_ascii=False)
        return True
    except Exception as e:
        print(f"Warning: Could not save config: {e}")
        return False

def get_favorite_paths():
    """Get list of favorite download paths."""
    config = load_config()
    return config.get('favorite_paths', [])

def add_favorite_path(path):
    """Add a path to favorites."""
    path = os.path.abspath(os.path.expanduser(path))
    config = load_config()
    
    if 'favorite_paths' not in config:
        config['favorite_paths'] = []
    
    # Don't add duplicates
    if path not in config['favorite_paths']:
        config['favorite_paths'].append(path)
        save_config(config)
        return True
    return False

def remove_favorite_path(path):
    """Remove a path from favorites."""
    path = os.path.abspath(os.path.expanduser(path))
    config = load_config()
    
    if 'favorite_paths' in config and path in config['favorite_paths']:
        config['favorite_paths'].remove(path)
        save_config(config)
        return True
    return False

def get_last_used_path():
    """Get the last used download path."""
    config = load_config()
    return config.get('last_used_path', None)

def set_last_used_path(path):
    """Save the last used download path."""
    path = os.path.abspath(os.path.expanduser(path))
    config = load_config()
    config['last_used_path'] = path
    save_config(config)

def get_default_download_path():
    """Get the default download path with smart defaults."""
    # INI/env override (new)
    ini_root = _env_get("ASD_DOWNLOAD_ROOT")
    if not ini_root:
        ini_cfg = load_ini_config()
        ini_root = _ini_get(ini_cfg, "DEFAULT", "save_directory")
    if isinstance(ini_root, str) and ini_root.strip():
        expanded = os.path.abspath(os.path.expanduser(ini_root.strip()))
        return expanded

    # Explicit user preference
    preferred = get_preferred_download_root()
    if preferred and os.path.exists(preferred):
        return preferred

    # Try last used
    last_path = get_last_used_path()
    if last_path and os.path.exists(last_path):
        return last_path
    
    # Try first favorite
    favorites = get_favorite_paths()
    if favorites and os.path.exists(favorites[0]):
        return favorites[0]
    
    # Linux: try a mounted disk with lots of free space
    if os.name == "posix" and os.path.exists("/proc/mounts"):
        disk_root = _smart_disk_root_linux()
        if disk_root:
            candidate = os.path.join(disk_root, "Anime-Sama-Downloader")
            return candidate

    # Default fallbacks
    possible_paths = [
        os.path.expanduser("~/Téléchargements"),
        os.path.expanduser("~/Downloads"),
        os.path.expanduser("~/videos"),
        "./videos"
    ]
    
    for path in possible_paths:
        if os.path.exists(path):
            return path
    
    # Last resort: create videos in current directory
    return "./videos"
