# Anime-Sama Video Downloader

[![Python 3.6+](https://img.shields.io/badge/Python-3.6+-3776AB?logo=python&logoColor=white)](https://www.python.org/)
[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)
[![GitHub issues](https://img.shields.io/github/issues/Guilhem-Bonnet/Anime-Sama-Downloader)](https://github.com/Guilhem-Bonnet/Anime-Sama-Downloader/issues)
[![GitHub stars](https://img.shields.io/github/stars/Guilhem-Bonnet/Anime-Sama-Downloader)](https://github.com/Guilhem-Bonnet/Anime-Sama-Downloader/stargazers)

A command-line tool to download anime episodes from anime-sama.si. Supports multiple video sources, parallel downloads, AniList-backed search, and includes both CLI and web interfaces.

**Fork of:** [sertrafurr/Anime-Sama-Downloader](https://github.com/sertrafurr/Anime-Sama-Downloader)

> **v2.6:** AniList search, modern TUI, parallel downloads (up to 10 concurrent), queue management. See [CHANGELOG.md](CHANGELOG.md).

## Features

### Core Functionality
- **Multiple interfaces**: Interactive CLI, command-line mode, terminal UI (`--tui`), and web UI
- **Smart search**: AniList-backed search with title synonyms and automatic URL resolution
- **Parallel downloads**: Download up to 10 anime/episodes simultaneously with `--jobs`
- **Multi-source support**: SendVid, Sibnet, Vidmoly, OneUpload, and more
- **Episode selection**: Ranges (1-5), specific episodes (3,5,7), or all
- **Threading**: Multi-threaded segment downloads for M3U8/HLS streams
- **Format conversion**: FFmpeg or MoviePy for .ts to .mp4 conversion
- **Progress tracking**: Real-time download progress with speed indicators
- **Queue management**: Cancel operations with Ctrl+C, automatic cleanup of partial files

### Technical Improvements Over Original
- **HTTP connection pooling**: 30% faster with persistent connections and smart caching
- **AniList integration**: Better search results without requiring API keys
- **Concurrent queue system**: Download manager with configurable parallelism (1-10 jobs)
- **Modern web UI**: React + FastAPI stack with TypeScript
- **Configuration system**: INI files + environment variables with layered priority
- **Robust error handling**: Automatic retries, graceful failures, partial download cleanup
- **Flexible output paths**: Smart directory structure with season/episode organization

## Quick Start

### Requirements

- Python 3.6+
- pip (Python package manager)

### Installation

```bash
git clone https://github.com/Guilhem-Bonnet/Anime-Sama-Downloader.git
cd Anime-Sama-Downloader
pip install -r requirements.txt
```

**Required packages:** `requests`, `beautifulsoup4`, `tqdm`, `moviepy`, `textual`, `fastapi`, `uvicorn`

### Basic Usage

```bash
# Interactive mode (guided)
python main.py

# Search and download
python main.py -s "demon slayer" -e 1-5 --quick

# Direct URL download
python main.py -u "https://anime-sama.si/catalogue/sword-art-online/saison1/vostfr/" -e 1-10 -t

# Parallel downloads
python main.py --jobs 5 -s "kaiju" -s "one piece" -e 1-12 --yes

# Modern terminal UI
python main.py --tui

# Web interface
python main.py --ui web
```

## Usage Modes

### 1. Interactive Mode

```bash
python main.py
```

Step-by-step guided interface. Best for beginners.

### 2. Search Mode

```bash
# Basic search
python main.py -s "kaiju" -e 1-5

# With season and language
python main.py -s "sword art online" --season 2 --lang vf -e 1-10

# Multiple anime
python main.py --jobs 5 -s "one piece" -s "naruto" -e 1-12 --yes
```

Works with Japanese, French, and English titles thanks to AniList.

### 3. URL Mode

```bash
# Single anime
python main.py -u "URL" -e 1-10 -t

# Multiple URLs
python main.py --jobs 3 -u "URL1" -u "URL2" -e 1-5 --mp4-threaded
```

URLs must follow format: `https://anime-sama.si/catalogue/{anime}/{season}/{lang}/`

### 4. Terminal UI (TUI)

```bash
python main.py --tui
```

Modern terminal interface powered by Textual.

### 5. Web Interface

```bash
python main.py --ui web
```

Opens web UI at `http://localhost:8000`. See [Web UI](#web-ui) section for details.

## CLI Arguments

```
-s, --search          Search anime by name (repeatable)
-u, --url             Direct anime-sama URL (repeatable)
-e, --episodes        Episodes to download: "1-5", "3,5,7", "all"
-p, --player          Player number (1, 2, 3...), auto-selected if omitted
-d, --directory       Output directory (default: config.ini or ./videos/)
-t, --threaded        Use multi-threaded downloads
--ts-threaded         Multi-threaded .ts segment downloads (faster for M3U8)
--mp4-threaded        Same as --ts-threaded --auto-mp4
--auto-mp4            Auto-convert to MP4 after download
--ffmpeg              Use FFmpeg instead of MoviePy for conversion
--jobs                Concurrent downloads (1-10, default: from config)
--quick               Skip interactive prompts, use first result
--yes                 Auto-answer yes to all prompts
--season              Season number for search (default: 1)
--lang                Language: vostfr, vf, vo (default: vostfr)
--search-provider     Search engine: anilist, local (default: anilist)
--tui                 Launch terminal UI
--ui web              Launch web interface
--config              Path to config.ini
```

## Configuration

Configuration priority (highest to lowest):
1. CLI arguments
2. Environment variables
3. `config.ini` file

### config.ini

Create from template:
```bash
cp config.ini.example config.ini
```

Example configuration:
```ini
[download]
save_directory = ./videos
max_concurrent_downloads = 5

[site]
base_url = https://anime-sama.si

[web]
host = 0.0.0.0
port = 8000
```

See [config.ini.example](config.ini.example) for all options.

### Environment Variables

```bash
ASD_CONFIG=./custom-config.ini          # Config file path
ASD_DOWNLOAD_ROOT=~/Downloads/anime     # Output directory
ASD_MAX_CONCURRENT_DOWNLOADS=10         # Parallel downloads
ASD_SITE_BASE_URL=https://anime-sama.si # Site URL
ASD_WEB_HOST=0.0.0.0                    # Web UI host
ASD_WEB_PORT=8000                       # Web UI port
```

Example:
```bash
ASD_MAX_CONCURRENT_DOWNLOADS=10 python main.py --ui web
```

## Web UI

### Development

Start backend:
```bash
./scripts/dev-backend.sh
```

Start frontend (in another terminal):
```bash
./scripts/dev-frontend.sh
```

Access at `http://localhost:5173`

### Docker (Development)

```bash
docker compose up --build
```

Access at `http://localhost:5173`

### Docker (Production)

```bash
docker compose -f docker-compose.prod.yml up --build
```

Access at `http://localhost:8000`

Frontend is built as static files and served by FastAPI backend.

## Video Source Support

| Platform | Status | Performance | Notes |
|----------|--------|-------------|-------|
| SendVid | ✅ Working | Good | Primary source |
| Sibnet | ✅ Working | Good | Reliable backup |
| Vidmoly | ✅ Working | Fast with threading | M3U8/HLS, use `--ts-threaded` |
| OneUpload | ✅ Working | Fast with threading | M3U8/HLS, use `--ts-threaded` |
| MoveArnPre | ✅ Working | Fast with threading | M3U8/HLS, use `--ts-threaded` |
| SmoothPre | ✅ Working | Fast with threading | M3U8/HLS, use `--ts-threaded` |
| Mivalyo | ✅ Working | Medium with threading | M3U8/HLS, use `--ts-threaded` |
| MYVI | ❌ Deprecated | None | Redirects to ads |
| VK.com | ❌ Unsupported | None | No working URLs found |

**Tip:** For M3U8/HLS sources (Vidmoly, OneUpload, etc.), use `--ts-threaded` for significantly faster downloads.

## Examples

### Search Examples
```bash
# Japanese titles
python main.py -s "shingeki no kyojin" -e 1-5

# French titles
python main.py -s "l'attaque des titans" -e 1-10

# English titles
python main.py -s "attack on titan" -e 1-5

# With season/language
python main.py -s "naruto" --season 2 --lang vf -e 1-20
```

### Batch Downloads
```bash
# Multiple anime, 5 parallel downloads
python main.py --jobs 5 \
  -s "one piece" \
  -s "naruto" \
  -u "https://anime-sama.si/catalogue/bleach/saison1/vostfr/" \
  -e 1-10 --mp4-threaded --yes
```

### Advanced Usage
```bash
# Download all episodes with FFmpeg conversion
python main.py -u "URL" -e all --auto-mp4 --ffmpeg -d ~/anime

# Specific episodes with custom player
python main.py -u "URL" -e 3,5,7,10 -p 2 --ts-threaded

# Quick search mode (no prompts)
python main.py -s "kaiju" -e 1-5 --quick --mp4-threaded
```

## Technical Architecture

### Core Components

- **main.py**: Entry point, argument parsing, mode selection
- **utils/download_manager.py**: Queue management, parallel downloads, cancellation handling
- **utils/downloaders/**: Video extraction and download logic for each source type
- **utils/search.py**: AniList API integration and local fuzzy search
- **utils/http_pool.py**: Connection pooling for 30% performance improvement
- **utils/ts_to_mp4.py**: Video conversion (FFmpeg or MoviePy)
- **utils/tui.py**: Terminal UI using Textual framework
- **utils/ui/web/app.py**: FastAPI backend
- **webapp/**: React + TypeScript frontend

### Key Improvements

1. **Connection Pooling**: Persistent HTTP connections reduce overhead
2. **Smart Caching**: Response caching for repeated requests
3. **Concurrent Queue**: Download manager with configurable parallelism
4. **AniList Integration**: Public API for better anime metadata (no key required)
5. **Flexible Configuration**: Layered config system (CLI > env > INI)
6. **Error Recovery**: Automatic retries, partial download cleanup
7. **Modern Stack**: React, TypeScript, FastAPI, Textual

## Search Guide

The search feature uses AniList to find anime and resolve anime-sama.si URLs automatically.

**How it works:**
1. Query AniList for anime titles and synonyms
2. Match against anime-sama.si catalog
3. Build URL with correct slug, season, and language

**Supported:**
- Multiple languages (Japanese, English, French)
- Synonym matching
- Season and language specification
- Fuzzy matching fallback

See [SEARCH_GUIDE.md](SEARCH_GUIDE.md) for details.

## Contributing

Contributions welcome. Open an issue before starting major changes.

**Issues:** [GitHub Issues](https://github.com/Guilhem-Bonnet/Anime-Sama-Downloader/issues)

## License

This project is licensed under the GNU General Public License v3.0. See [LICENSE](LICENSE) for details.

## Disclaimer

This tool is for educational purposes. Respect copyright laws in your jurisdiction and anime-sama.si's terms of service.

## Credits

**Original fork:** [sertrafurr/Anime-Sama-Downloader](https://github.com/sertrafurr/Anime-Sama-Downloader)

Core video extraction logic: Human-developed  
Enhancements and UI improvements: Mix of manual development and AI assistance
