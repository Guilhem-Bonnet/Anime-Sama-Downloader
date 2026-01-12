"""Configuration management for favorite paths and settings."""

import json
import os
import shutil
from pathlib import Path

CONFIG_FILE = os.path.expanduser("~/.anime-sama-downloader.json")


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
