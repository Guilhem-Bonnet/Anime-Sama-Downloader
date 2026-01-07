"""
Configuration management for favorite paths and settings
"""

import json
import os
from pathlib import Path

CONFIG_FILE = os.path.expanduser("~/.anime-sama-downloader.json")

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
    # Try last used
    last_path = get_last_used_path()
    if last_path and os.path.exists(last_path):
        return last_path
    
    # Try first favorite
    favorites = get_favorite_paths()
    if favorites and os.path.exists(favorites[0]):
        return favorites[0]
    
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
