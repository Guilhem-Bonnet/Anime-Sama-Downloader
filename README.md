<div align="center">
   
# 🎌 Anime-Sama Video Downloader

<img src="https://img.shields.io/badge/Python-3.6+-blue.svg?style=for-the-badge&logo=python" alt="Python Version">
<img src="https://img.shields.io/badge/Platform-Windows%20%7C%20macOS%20%7C%20Linux_(mostly_windows)-lightgrey.svg?style=for-the-badge" alt="Platform">
<img src="https://img.shields.io/badge/License-GPL_V3-green.svg?style=for-the-badge" alt="License">
<img src="https://img.shields.io/badge/Version-2.6_Optimized-orange.svg?style=for-the-badge" alt="Version">

**🚀 A powerful, beautiful and simple CLI tool to download anime episodes from anime-sama.si**

*Enhanced with colorful interface, smart source detection, CLI mode, and robust error handling*

*Questions? Unworking urls? Open an issue, will be added fastly (hopefully)*

### 🌟 Star this repo if it helped you!

![Website Support](https://img.shields.io/badge/Website%20Support-100%25-brightgreen)

> **✨ NEW in v2.6**: AniList-backed search + URL resolution, optional modern TUI (`--tui`)  
> See [CHANGELOG.md](CHANGELOG.md) for full details | [UX improvements](UX_IMPROVEMENTS.md)

> **⚡ NEW**: Queue + téléchargements parallèles (jusqu’à `--jobs 10`) + annulation (Ctrl+C)

### Scans support ? 5 stars and it will be added !
## ✨ Features

<table>
<tr>
<td width="50%">

### 🎯 **Smart & Intuitive**
- 🌈 **Beautiful CLI Interface** with colors and emojis
- 💻 **Triple Mode**: Interactive / Quick / CLI (NEW!)
- 🔍 **Smart Search Engine** - Find animes by name! (NEW!)
- 🧠 **AniList-backed Search** - Better titles & synonyms (NEW!)
- 🖥️ **Modern Terminal UI (TUI)** - Optional `--tui` mode (NEW!)
- 🎯 **Smart Defaults** - Just press Enter! (NEW!)
- ✅ **Auto URL Validation** with helpful error messages
- 📝 **Built-in Tutorial** for first-time users
- ⚡ **Multi-threaded Downloads** for blazing fast performance
- 🧵 **Download Queue (1–10 jobs)** - Plusieurs animes/épisodes en parallèle (NEW!)
- 🛑 **Cancel/Stop** - Ctrl+C en CLI, “Annuler tout” en TUI (NEW!)
</td>
<td width="50%">

### ⚡ **Powerful & Reliable**  
- 🎪 **Multiple Player Support** (Player 1, 2, 3...)
- 🔄 **Smart Source Detection** (SendVid, Sibnet and others)
- 📊 **Real-time Progress** with download speeds
- 🛡️ **Robust Error Handling** with retry logic
- 📺 **Multiple Episode Selection** with threads supports
- 😊 **FFmpeg support** choose between 2 converters
- 🚀 **HTTP Connection Pooling** - 30% faster (NEW!)
- 💾 **Smart Response Caching** (NEW!)

</tr>
</table>

<details>
  <summary><strong>🚩 Stealing/Malware List (Shows every fake/copy without credit made out of my code.) (Click to reveal)</strong></summary>

  | Username | Link | Description |
  |----------|------|-------------|
  | `OMTSE` | [Repo](https://github.com/OMTSE/Anime-Sama-Downloader) | Used code without credit |

DO NOT HARASS ANY INDIVIDUAL IN THIS LIST (They probably may be ban at some point aswell.)
</details>



---

## 🚀 Quick Start

### 📋 Prerequisites

<details>
<summary>🐍 <strong>Python Requirements</strong></summary>

Make sure you have **Python 3.6+** installed:

```bash
# Check Python version
python --version

# Install required packages
pip install -r requirements.txt
```

**Required Libraries:**
- `requests` - HTTP requests handling
- `beautifulsoup4` - HTML parsing
- `tqdm` - Progress bar display

</details>

### ⚡ Installation & Usage

```bash
# 1. Clone the repository
git clone https://github.com/Guilhem-Bonnet/Anime-Sama-Downloader.git

# 2. Navigate into the project directory
cd Anime-Sama-Downloader

# 3. Install dependencies
python3 -m pip install -r requirements.txt

# 3. Run the magic! ✨
python3 main.py

# OR use search to find animes! (NEW!)
python3 main.py -s "kaiju" -e 1-5 --quick

# OR batch mode: multiple animes in one command (NEW!)
python3 main.py --jobs 10 \
  -s "kaiju" -s "one piece" \
  -u "https://anime-sama.si/catalogue/roshidere/saison1/vostfr/" \
  -e 1-5 --mp4-threaded --yes

# OR use CLI mode with URL
python3 main.py -u "ANIME_URL" -e 1-5 -t --auto-mp4

# OR launch the modern TUI (optional)
python3 main.py --tui

# OR launch the web UI (backend)
python3 main.py --ui web
```

---

## ⚙️ Configuration (INI + variables d'env)

Un fichier prêt à l'emploi est fourni : [config.ini](config.ini). Tu peux aussi partir de [config.ini.example](config.ini.example).

### Ordre de priorité

1) Arguments CLI (ex: `--directory`, `--jobs`)
2) Variables d'environnement (ex: `ASD_DOWNLOAD_ROOT`)
3) `config.ini` (ou un autre INI via `--config` / `ASD_CONFIG`)

### Variables d'environnement supportées

- `ASD_CONFIG` : chemin vers un `config.ini` à utiliser
- `ASD_DOWNLOAD_ROOT` : dossier racine de sortie (équivalent à `save_directory`)
- `ASD_MAX_CONCURRENT_DOWNLOADS` : concurrence globale (1–10)
- `ASD_SITE_BASE_URL` / `ASD_SITE_DOMAIN` : override du domaine Anime‑Sama si ça change
- `ASD_WEB_HOST` / `ASD_WEB_PORT` : bind du serveur web

### Exemples

```bash
# Utiliser un autre fichier de config
python3 main.py --config ./config.ini --ui web

# Override rapide sans toucher au fichier
ASD_SITE_BASE_URL=https://anime-sama.si ASD_MAX_CONCURRENT_DOWNLOADS=10 python3 main.py --ui web
```

---

## 🌐 Web UI (Dev / Docker)

### Dev local (2 terminaux)

Backend:
```bash
./scripts/dev-backend.sh
```

Frontend:
```bash
./scripts/dev-frontend.sh
```

Ouvre `http://127.0.0.1:5173`.

### Dev Docker

```bash
docker compose up --build
```

Ouvre `http://127.0.0.1:5173`.

### Prod Docker (SPA servie par le backend)

```bash
docker compose -f docker-compose.prod.yml up --build
```

Ouvre `http://127.0.0.1:8000/`.

### 💡 Three Ways to Use

<table>
<tr>
<td width="33%">

#### 🎨 Interactive Mode
```bash
python main.py
```
- Guided step-by-step
- Built-in search
- Perfect for beginners

</td>
<td width="33%">

#### 🔍 Search Mode (NEW!)
```bash
python main.py -s "kaiju" -e 1-10
```
- Find by name
- Fuzzy matching
- Auto-translations

</td>
<td width="33%">

#### ⚡ CLI Mode
```bash
python main.py -u "URL" -e 1-10
```
- Fully automated
- Script-friendly
- One-command downloads

</td>
</tr>
</table>

#### 📝 CLI Examples

```bash
# 🔍 Search by anime name (NEW!)
python main.py -s "kaiju" -e 1-5 --quick
python main.py -s "l'attaque des titans" -e 1-10 --quick
python main.py -s "demon slayer" -e 1-26 -t --auto-mp4

# 🎯 Search with custom season/language (NEW!)
python main.py -s "sword art online" --season 2 --lang vf -e 1-5
python main.py -s "naruto" --season 1 --lang vostfr -e 1-10 --quick

# Force legacy local search (no AniList)
python main.py -s "kaiju" --search-provider local -e 1-3

# Download episodes 1 to 5 with threading
python main.py -u "https://anime-sama.si/catalogue/sword-art-online/saison1/vostfr/" -e 1-5 -t

# Download specific episodes
python main.py -u "URL" -e 3,5,7,10 --ts-threaded

# Download all episodes with auto-conversion
python main.py -u "URL" -e all --auto-mp4 --ffmpeg -d ~/Downloads

# Get help
python main.py --help

# Launch TUI
python main.py --tui

# Batch download (queue) - 10 parallèles
python main.py --jobs 10 -s "one piece" -s "naruto" -e 1-12 --mp4-threaded --yes
```

> **🔍 See [SEARCH_GUIDE.md](SEARCH_GUIDE.md) for complete search documentation**

---

## 📖 Complete Usage Guide

<div align="center">
<h3>🎯 Three Simple Steps</h3>
</div>

<table>
<tr>
<td width="33%" align="center">

### 1️⃣ Find Anime
<img src="https://img.shields.io/badge/Step-1-blue?style=for-the-badge">

Visit **[anime-sama.tv](https://anime-sama.tv/catalogue/)**

🔍 Search your anime  
📺 Select season & language  
📋 Copy the complete URL

</td>
<td width="33%" align="center">

### 2️⃣ Run Script  
<img src="https://img.shields.io/badge/Step-2-green?style=for-the-badge">

Launch the downloader

🔎 Search or paste URL (NEW!)
🎮 Choose player & episode  
📁 Set download folder

</td>
<td width="33%" align="center">

### 3️⃣ Enjoy!
<img src="https://img.shields.io/badge/Step-3-purple?style=for-the-badge">

Watch the magic happen

⬇️ Auto-download starts  
📊 Real-time progress  
🎉 Episode ready to watch!

</td>
</tr>
</table>

### 🔎 Search Examples (NEW!)

```bash
# 🎯 Search by anime name
python main.py -s "kaiju" -e 1-5
python main.py -s "demon slayer" -e 1-10 --quick

# 🌍 French titles work too!
python main.py -s "l'attaque des titans" -e 1-5
python main.py -s "attaque des titans" -e 1-10

# 🇯🇵 Japanese titles
python main.py -s "shingeki no kyojin" -e 1-5
python main.py -s "kimetsu no yaiba" -e 1-26

# ⚡ Quick search (best for automation)
python main.py -s "one piece" -e 1-10 --quick
```

### 🔗 Example URLs

```bash
# ✅ Perfect URL format
https://anime-sama.tv/catalogue/roshidere/saison1/vostfr/
https://anime-sama.tv/catalogue/demon-slayer/saison1/vf/
https://anime-sama.tv/catalogue/attack-on-titan/saison3/vostfr/
https://anime-sama.tv/catalogue/one-piece/saison1/vostfr/

# ❌ Won't work
https://anime-sama.tv/catalogue/roshidere/  # Missing season/language
https://anime-sama.tv/  # Just homepage
```

---

## 🛠️ Video Source Support

<div align="center">

| Platform | Status | Performance | Notes |
|:--------:|:------:|:-----------:|:------|
| 📹 **SendVid** | ![Working](https://img.shields.io/badge/Status-✅_Working-brightgreen) | 🔄 Good | Primary recommended source |
| 🎬 **Sibnet** | ![Working](https://img.shields.io/badge/Status-✅_Working-brightgreen) | 🔄 Good | Reliable backup source |
| 🎬 **Vidmoly** | ![Working](https://img.shields.io/badge/Status-✅_Working-brightgreen) | 🔄 SLOW if not threaded. FASTEST if | Download .ts file then make them into an mp4 back. |
| 🎬 **ONEUPLOAD** | ![Working](https://img.shields.io/badge/Status-✅_Working-brightgreen) | 🔄 SLOW if not threaded. Very fast if | Download .ts file then make them into an mp4 back. |
| 🎬 **MOVEARNPRE** | ![Working](https://img.shields.io/badge/Status-✅_Working-brightgreen) | 🔄 SLOW if not threaded. Very fast if  | Download .ts file then make them into an mp4 back. |
| 🎬 **SMOOTHPRE** | ![Working](https://img.shields.io/badge/Status-✅_Working-brightgreen) | 🔄 SLOW if not threaded. Very fast if | Download .ts file then make them into an mp4 back. |
| 🎬 **MIVALYO** | ![Working](https://img.shields.io/badge/Status-✅_Working-brightgreen) | 🔄 SLOW if not threaded. Faster if | Download .ts file then make them into an mp4 back. |
| 🚫 **MYVI** | ![Deprecated](https://img.shields.io/badge/Status-❌_Deprecated-red) | ❌ None | Scam website, only redirect to advertisement. |
| 🤔 **VK.com** | ![Deprecated](https://img.shields.io/badge/Status-❌_Unsupported-red) | ❌ None | Could try, but did not find any working URL. |

</div>

---

## 📸 Screenshots

<details>
<summary>🖼️ <strong>View CLI Interface Screenshots</strong></summary>

### 🎨 Main Interface
```
╔══════════════════════════════════════════════════════════════╗
║                 ANIME-SAMA VIDEO DOWNLOADER                  ║
║                       Enhanced CLI v2.0                      ║
╚══════════════════════════════════════════════════════════════╝

📺 Download anime episodes from anime-sama.tv easily!
```

### 🎮 Player Selection
```
🎮 SELECT PLAYER
─────────────────────────────────────────────────────────────────
  1. Player 1 (12/15 working episodes)
  2. Player 2 (8/15 working episodes)  
  3. Player 3 (15/15 working episodes)

Enter player number (1-3) or type player name:
```

### 📊 Download Progress
```
⬇️ DOWNLOADING
─────────────────────────────────────────────────────────────────
📥 roshidere_episode_1.mp4: 100%|████████| 145M/145M [02:15<00:00, 1.07MB/s]
✅ Download completed successfully!
```

</details>

---

## ⚙️ Configuration

<details>
<summary>🔧 <strong>Customization Options</strong></summary>

### 📁 Default Settings
- **Download Directory**: `./videos/`
- **Video Format**: `.mp4`
- **Naming Convention**: `{anime_name}_episode_{number}.mp4`

### 🎨 Color Themes
The script uses a beautiful color scheme:
- 🔵 **Info**: Cyan messages
- ✅ **Success**: Green confirmations  
- ⚠️ **Warning**: Yellow alerts
- ❌ **Error**: Red error messages
- 💜 **Headers**: Purple titles

</details>

---

## 🤝 Contributing

<div align="center">

We welcome contributions! Here's how you can help:

[![Issues](https://img.shields.io/badge/Issues-Welcome-blue?style=for-the-badge)](https://github.com/sertrafurr/Anime-Sama-Downloader/issues)
[![Pull Requests](https://img.shields.io/badge/PRs-Welcome-green?style=for-the-badge)](https://github.com/sertrafurr/Anime-Sama-Downloader/pulls)
[![Discussions](https://img.shields.io/badge/Discussions-Join-purple?style=for-the-badge)](https://github.com/sertrafurr/Anime-Sama-Downloader/discussions)

</div>

### 🐛 Found a Bug?
 Check existing [issues](https://github.com/sertrafurr/issues)
 Create a new issue with:
    📝 Clear description
    🔄 Steps to reproduce
    💻 System information

### 💡 Feature Request?
 Open a [discussion](https://github.com/sertrafurr/discussions)
 Explain your idea
 Community feedback welcome!

---

## 📄 License

<div align="center">

This project is licensed under the **GPL v3 License**

[![License: GPL](https://img.shields.io/badge/License-GPL_V3-yellow.svg?style=for-the-badge)](https://opensource.org/licenses/MIT)

*Feel free to use, modify, and distribute!*

</div>

---

## ⚠️ Disclaimer

<div align="center">
<strong>📢 Important Notice</strong>
</div>

- 🎯 This tool is for **educational purposes** only
- 📺 Respect **copyright laws** in your jurisdiction  
- 🔒 Use responsibly and in compliance with anime-sama.tv's terms

---

<div align="center">

## 🙏 Acknowledgments

<img src="https://img.shields.io/badge/Made_with-❤️-red?style=for-the-badge">

**🧠 Core algorithms and video extraction logic: Human-developed**  
**🎨 Code restructuring and user interface enhancements: AI-assisted**

---

### 🌟 Star this repo if it helped you!

[![Stars](https://img.shields.io/github/stars/sertrafurr/anime-sama-downloader?style=for-the-badge&logo=github)](https://github.com/sertrafurr/anime-sama-downloader/stargazers)

</div>

You wish for something/a service to get removed/added, open an issue.
