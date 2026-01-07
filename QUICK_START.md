# 🚀 Guide Rapide - 3 Façons d'Utiliser

## Mode 1 : Interactif Standard (Pour Débutants)
**Le plus simple, avec aide intégrée**

```bash
python main.py
```

**Ce qui se passe :**
1. Besoin d'aide ? Appuyez sur `h` (ou Enter pour passer)
2. Entrez l'URL de l'anime
3. Choisissez le player
4. Sélectionnez les épisodes
5. Questions avec défauts intelligents (juste appuyez sur Enter !)

**Temps : ~1 minute d'interaction**

---

## Mode 2 : Quick Mode (Recommandé) ⚡
**Optimal avec minimum de questions**

```bash
python main.py --quick
```

**Ce qui se passe :**
1. Entrez l'URL
2. Choisissez player et épisodes
3. **C'est tout !** Le reste est automatique avec les meilleurs réglages

**Temps : ~30 secondes d'interaction**

**Utilise automatiquement :**
- ✓ Threading si plusieurs épisodes
- ✓ Fast .ts downloads pour M3U8
- ✓ Conversion MP4 automatique
- ✓ ffmpeg si installé, sinon moviepy

---

## Mode 3 : CLI Full Auto (Pour Experts) 🔥
**Zéro interaction, scriptable**

```bash
python main.py -u "URL_ANIME" -e 1-10 -t --auto-mp4 --ts-threaded
```

**Ce qui se passe :**
1. Rien ! Tout est spécifié en arguments
2. Le téléchargement démarre immédiatement

**Temps : 0 seconde d'interaction**

---

## Mode 4 : UI Terminal Moderne (TUI) 🖥️
**Interface moderne en terminal (optionnelle)**

```bash
python main.py --tui
```

**Note :** la TUI est optionnelle. La CLI reste le mode par défaut.

Si tu n'as pas encore installé les dépendances :

```bash
python3 -m pip install -r requirements.txt
```

---

## 📊 Comparaison Rapide

| Mode | Questions | Temps Setup | Idéal Pour |
|------|-----------|-------------|------------|
| **Interactif** | 5 | ~1 min | Débutants, découverte |
| **Quick** | 3 | ~30 sec | Utilisation quotidienne ⭐ |
| **CLI** | 0 | 0 sec | Scripts, automation |

---

## 💡 Exemples Concrets

### Débutant - Première fois
```bash
python main.py
# Suivez les instructions, appuyez sur Enter pour les défauts
```

### Utilisateur Régulier - Téléchargement Rapide
```bash
python main.py --quick
# URL + Player + Episodes = c'est parti !
```

### Power User - Script Quotidien
```bash
# Télécharger les nouveaux épisodes automatiquement
python main.py \
  -u "https://anime-sama.tv/catalogue/one-piece/saison1/vostfr/" \
  -e 1010-1015 \
  -t --auto-mp4 --ts-threaded \
  -d ~/Videos/OnePiece
```

### Batch Download - Plusieurs Animes
```bash
# Script shell pour télécharger plusieurs animes
#!/bin/bash
python main.py -u "URL_ANIME1" -e 1-12 -t --auto-mp4 &
python main.py -u "URL_ANIME2" -e 1-12 -t --auto-mp4 &
python main.py -u "URL_ANIME3" -e 1-12 -t --auto-mp4 &
wait
echo "Tous les téléchargements terminés!"
```

---

## 🎯 Recommandations

### Pour 90% des Utilisateurs
```bash
python main.py --quick
```
**Pourquoi ?** : Parfait équilibre entre simplicité et contrôle.

### Pour Automatisation / Cron Jobs
```bash
python main.py -u "URL" -e all -t --auto-mp4 -d ~/Videos
```
**Pourquoi ?** : Scriptable, aucune interaction requise.

### Pour Première Utilisation
```bash
python main.py
# Tapez 'h' pour voir le tutoriel
```
**Pourquoi ?** : Guidé pas à pas avec aide intégrée.

---

## 🔥 Alias Pratique

Ajoutez à votre `~/.zshrc` ou `~/.bashrc` :

```bash
# Mode quick par défaut
alias anime-dl="cd ~/Anime-Sama-Downloader && python main.py --quick"

# Mode CLI complet
alias anime-get="cd ~/Anime-Sama-Downloader && python main.py"
```

Utilisation :
```bash
anime-dl                    # Lance en mode quick
anime-get -u "URL" -e 1-10  # CLI full
```

---

## 📝 Tableau des Arguments

| Argument | Court | Description | Exemple |
|----------|-------|-------------|---------|
| `--url` | `-u` | URL anime-sama | `-u "https://..."` |
| `--search` | `-s` | Recherche par nom | `-s "kaiju"` |
| `--episodes` | `-e` | Episodes (range/liste/all) | `-e 1-5` `-e 3,5,7` `-e all` |
| `--player` | `-p` | Numéro player | `-p 2` |
| `--directory` | `-d` | Dossier de sauvegarde | `-d ~/Downloads` |
| `--threaded` | `-t` | Threading épisodes | `-t` |
| `--ts-threaded` | - | Threading .ts segments | `--ts-threaded` |
| `--auto-mp4` | - | Conversion auto MP4 | `--auto-mp4` |
| `--quick` | - | Mode quick (défauts smart) | `--quick` |
| `--search-provider` | - | Provider recherche (anilist/local) | `--search-provider anilist` |
| `--tui` | - | UI terminal moderne (Textual) | `--tui` |
| `--ffmpeg` | - | Force ffmpeg | `--ffmpeg` |
| `--moviepy` | - | Force moviepy | `--moviepy` |
| `--no-tutorial` | - | Skip tutoriel | `--no-tutorial` |
| `--help` | `-h` | Aide complète | `--help` |
| `--version` | - | Version du programme | `--version` |

---

## ⚡ Quick Tips

1. **Juste appuyez sur Enter** : Les défauts sont optimaux
2. **Besoin d'aide ?** : Tapez `h` quand demandé
3. **Pressé ?** : Utilisez `--quick`
4. **Automatiser ?** : Mode CLI complet
5. **Erreur ?** : Le programme retry automatiquement

---

**Version** : 2.6-optimized  
**Date** : Janvier 2026  
**🎯 Recommandation** : Mode `--quick` pour 90% des cas d'usage
