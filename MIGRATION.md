# 🧭 Migration (versions précédentes → état actuel)

Ce document résume les points à connaître si tu viens d’une ancienne version (ou d’un ancien fork).

## 1) Domaine Anime‑Sama

Le domaine par défaut est **anime-sama.si**.

Si tes scripts utilisaient `anime-sama.tv`, remplace simplement l’URL. Sinon, tu peux surcharger le domaine :

- variable d’environnement :

```bash
ASD_SITE_BASE_URL=https://anime-sama.si python main.py -s "kaiju" -e 1-3
```

- `config.ini` :

```ini
[SITE]
base_url = https://anime-sama.si
```

## 2) Interfaces disponibles

- **CLI interactif** : `python main.py`
- **CLI scriptable** : `python main.py -s ...` ou `python main.py -u ...`
- **TUI (Textual)** : `python main.py --tui`
- **Interface Web** : `python main.py --ui web` (ou via Docker)

## 3) Téléchargements en parallèle

La file de téléchargement permet de lancer plusieurs épisodes en parallèle :

```bash
python main.py --jobs 5 -s "one piece" -e 1-12 --yes
```

## 4) Docker : changement important sur le dossier de sortie

En Docker :

- sortie dans le conteneur : `/data/videos`
- sortie sur l’hôte : dossier monté via `ASD_HOST_DOWNLOAD_ROOT`

```bash
cp .env.example .env
# ASD_HOST_DOWNLOAD_ROOT=/chemin/absolu/sur/hote
docker compose up --build
```

L’interface Web en Docker n’accepte pas de chemin absolu “hôte” : on choisit uniquement un **sous-dossier** sous `/data/videos`.

## 5) Compatibilité

- Tes usages “classiques” (`python main.py` et téléchargement par URL) restent valables.
- Si tu vois une doc qui parle de `anime-sama.tv`, considère-la comme obsolète : utilise `anime-sama.si`.
