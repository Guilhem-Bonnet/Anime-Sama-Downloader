# Changelog - Version Optimisée

## v2.6 (Janvier 2026)

### 🔍 Recherche améliorée (AniList)
- Recherche enrichie par titres + synonymes (AniList, sans clé)
- Résolution automatique de l'URL anime-sama (`/catalogue/<slug>/`) avec cache local
- Option `--search-provider anilist|local`

### 🖥️ TUI (Terminal UI) moderne
- Interface terminal moderne (Textual), optionnelle via `--tui`
- La CLI reste le comportement par défaut

---

## v2.5 (Optimized)

## 🚀 Nouvelles Fonctionnalités

### Interface en Ligne de Commande (CLI)
Vous pouvez maintenant utiliser le script avec des arguments en ligne de commande pour automatiser les téléchargements !

**Exemples d'utilisation :**

```bash
# Mode interactif (comportement par défaut)
python main.py

# Télécharger des épisodes spécifiques avec URL
python main.py -u "https://anime-sama.tv/catalogue/sword-art-online/saison1/vostfr/" -e 1-5

# Télécharger avec threading activé
python main.py -u "URL" -e 3,5,7 -t -d ~/Downloads

# Télécharger tous les épisodes avec conversion automatique en MP4
python main.py -u "URL" -e all --auto-mp4 --ts-threaded

# Utiliser un player spécifique
python main.py -u "URL" -e 1-10 -p 2 --threaded
```

**Arguments disponibles :**

- `-u, --url` : URL anime-sama (ex: https://anime-sama.tv/catalogue/...)
- `-e, --episodes` : Épisodes à télécharger (ex: "1-5", "3,5,7", "all")
- `-p, --player` : Numéro du player à utiliser (si omis: auto-sélection)
- `-d, --directory` : Répertoire de sauvegarde (défaut: ./videos)
- `-t, --threaded` : Activer les téléchargements threadés (plus rapide)
- `--ts-threaded` : Activer le téléchargement threadé des segments .ts (beaucoup plus rapide pour M3U8)
- `--auto-mp4` : Convertir automatiquement les fichiers .ts en .mp4
- `--ffmpeg` : Utiliser ffmpeg pour la conversion (plus rapide, défaut)
- `--moviepy` : Utiliser moviepy pour la conversion (plus lent mais plus léger)
- `--no-tutorial` : Ignorer l'invite du tutoriel
- `--search-provider` : Provider pour `--search` (anilist/local)
- `--tui` : Lancer la TUI (Textual)
- `--version` : Afficher la version

## 🐛 Corrections de Bugs

### 1. **Expansion du tilde (~) corrigée**
- ✅ Le bug où `~/Téléchargement` créait un dossier littéral au lieu d'utiliser le chemin absolu a été corrigé
- Les chemins sont maintenant correctement expansés avec `os.path.expanduser()`

### 2. **Imports dupliqués nettoyés**
- ❌ Avant : `import os`, `import sys`, `import subprocess` apparaissaient plusieurs fois
- ✅ Maintenant : Imports propres et organisés en haut du fichier

## ⚡ Optimisations de Performance

### 1. **Pool de Connexions HTTP**
- Implémentation d'un pool de connexions réutilisables
- Réduction de la latence réseau grâce à la réutilisation des connexions
- Configuration automatique des retries avec backoff exponentiel

### 2. **Cache des Réponses HTTP**
- Cache intelligent pour les requêtes fréquentes (episodes.js, playlists M3U8)
- Réduction du nombre de requêtes réseau
- Cache LRU avec limite de taille (100 entrées max)

### 3. **Retry Strategy**
- Retry automatique avec stratégie de backoff exponentiel
- Gestion des erreurs temporaires (429, 500, 502, 503, 504)
- Maximum 3 tentatives par requête

## 📦 Nouveaux Fichiers

- `utils/http_pool.py` : Gestionnaire de pool HTTP et cache
- `CHANGELOG.md` : Ce fichier de changelog

## 🔧 Améliorations Techniques

### Robustesse
- Meilleure gestion des timeouts
- Validation améliorée des chemins de fichiers
- Gestion des erreurs plus explicite

### Code Quality
- Imports organisés et dédupliqués
- Meilleure séparation des responsabilités
- Documentation améliorée

## 📊 Gains de Performance Estimés

- **Pool HTTP** : ~20-30% plus rapide sur les requêtes répétées
- **Cache** : Jusqu'à 90% de réduction des requêtes pour les métadonnées
- **Threading .ts** : 5-10x plus rapide (déjà présent, maintenant plus accessible via CLI)

## 🔄 Rétrocompatibilité

✅ **100% compatible** avec la version précédente !
- Le mode interactif fonctionne exactement comme avant
- Toutes les fonctionnalités existantes sont préservées
- La nouvelle interface CLI est optionnelle

## 🎯 Utilisation Recommandée

### Pour un usage ponctuel :
```bash
python main.py  # Mode interactif classique
```

### Pour automatiser / scripter :
```bash
# Télécharger les 5 premiers épisodes rapidement
python main.py -u "URL_ANIME" -e 1-5 --threaded --ts-threaded --auto-mp4

# Télécharger toute une saison dans un dossier spécifique
python main.py -u "URL_ANIME" -e all -d ~/Animes/SwordArtOnline --threaded
```

## 🚧 Prochaines Améliorations Prévues

- [ ] Reprise des téléchargements interrompus
- [ ] Fichier de configuration (config.json/yaml)
- [ ] Progress bar global pour les téléchargements multiples
- [ ] Export/import de listes d'épisodes
- [ ] Notification de fin de téléchargement

## 📝 Notes de Version

**Version** : 2.6-optimized  
**Date** : Janvier 2026  
**Branche** : `main`  
**Status** : Stable

---

## 🙏 Comment Tester

1. Tester le mode CLI :
   ```bash
   python main.py -u "https://anime-sama.tv/catalogue/roshidere/saison1/vostfr/" -e 1 --no-tutorial
   ```

2. Tester le mode interactif (doit fonctionner comme avant) :
   ```bash
   python main.py
   ```

3. Tester la TUI :
   ```bash
   python main.py --tui
   ```

## 📮 Feedback

Des bugs ? Des suggestions ? Ouvrez une issue sur GitHub !
