# 🚀 Quick Start (CLI / TUI / Web / Docker)

Ce guide donne les commandes “prêtes à copier-coller”. Pour le détail : [README.md](README.md).

## 1) Installer

```bash
python3 -m pip install -r requirements.txt
```

Optionnel (recommandé) : installer `ffmpeg` pour convertir les flux `.ts` en `.mp4` plus rapidement.

## 2) Télécharger (CLI)

### Mode interactif (simple)

```bash
python main.py
```

### Mode rapide (recommandé au quotidien)

```bash
python main.py --quick
```

### Recherche par nom (AniList)

```bash
python main.py -s "one piece" --season 1 --lang vostfr -e 1-12
```

### URL directe

```bash
python main.py -u "https://anime-sama.si/catalogue/roshidere/saison1/vostfr/" -e 1-12 -t
```

### Batch / parallèle (jusqu’à 10)

```bash
python main.py --jobs 5 \
  -s "kaiju" \
  -s "naruto" \
  -e 1-6 --yes
```

## 3) Interface terminal (TUI)

```bash
python main.py --tui
```

## 4) Interface Web (dev local)

Backend :
```bash
./scripts/dev-backend.sh
```

Variables optionnelles (par défaut : `127.0.0.1:8000`) : `ASD_WEB_HOST`, `ASD_WEB_PORT`.

Frontend :
```bash
./scripts/dev-frontend.sh
```

Variables optionnelles (par défaut : `127.0.0.1:5173`) : `ASD_WEBAPP_HOST`, `ASD_WEBAPP_PORT`.

Ouvre ensuite :
- http://127.0.0.1:5173 (SPA)
- http://127.0.0.1:8000 (API + fallback minimal)

## 5) Docker

### Dev

```bash
docker compose up --build
```

Accès : http://localhost:5173

### Prod

```bash
docker compose -f docker-compose.prod.yml up --build
```

Accès : http://localhost:8000

### Dossier de sortie (Docker)

- Dans le conteneur : `/data/videos`
- Sur l’hôte : configuré par `ASD_HOST_DOWNLOAD_ROOT`

```bash
cp .env.example .env
# éditer .env puis relancer docker compose
```

Dans l’interface Web en Docker : la destination est un **sous-dossier** sous `/data/videos` (pas un chemin absolu hôte).

---

## ⚡ Astuces

1. **Appuie sur Entrée** : les défauts sont optimaux
2. **Besoin d'aide ?** : Tapez `h` quand demandé
3. **Pressé ?** : utilise `--quick`
4. **Automatiser ?** : Mode CLI complet
5. **Erreur ?** : retries automatiques (si temporaire)

---

**Version** : 2.6-optimized  
**Date** : Janvier 2026  
**🎯 Recommandation** : Mode `--quick` pour 90% des cas d'usage
