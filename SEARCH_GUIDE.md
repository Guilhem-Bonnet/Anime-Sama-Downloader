# 🔍 Guide de recherche (AniList + résolution d’URL)

Le projet propose une recherche par nom qui **résout automatiquement** l’URL anime-sama (domaine par défaut : **anime-sama.si**).

## Comment ça marche (résumé)

1. Requête AniList (titres, synonymes, variantes)
2. Normalisation / scoring (fuzzy)
3. Résolution du bon “slug” Anime‑Sama pour construire une URL du type :

```
https://anime-sama.si/catalogue/<slug>/saisonN/<lang>/
```

Si le domaine Anime‑Sama change : tu peux le surcharger via config ou variable d’environnement (voir plus bas).

## Utilisation

### Recherche simple

```bash
python main.py -s "one piece" -e 1-12
```

### Saison + langue

```bash
python main.py -s "sword art online" --season 2 --lang vf -e 1-10
```

Langues supportées :

- `vostfr`
- `vf`
- `vo`

### Plusieurs animes (batch)

```bash
python main.py --jobs 5 \
    -s "kaiju" \
    -s "naruto" \
    -e 1-6 --yes
```

## Fournisseurs de recherche

- `anilist` (défaut) : meilleurs résultats, pas de clé API.
- `local` : fallback simplifié.

Exemple :

```bash
python main.py -s "attack on titan" --search-provider anilist -e 1-5
```

## Configuration du domaine (si ça bouge)

### Via variable d’environnement

```bash
ASD_SITE_BASE_URL=https://anime-sama.si python main.py -s "kaiju" -e 1-3
```

### Via config.ini

Dans `config.ini` :

```ini
[SITE]
base_url = https://anime-sama.si
```

## Dépannage rapide

- **Aucun résultat / mauvais anime** : essaie un titre anglais/japonais, ou précise `--season`.
- **Le site a changé de domaine** : règle `ASD_SITE_BASE_URL` ou `[SITE] base_url`.
- **Erreur réseau/timeout** : ajuste `[NETWORK] timeout` dans `config.ini`.

