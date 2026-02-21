# 🔍 Guide de recherche (AniList + résolution d’URL, serveur Go)

Le serveur expose des endpoints pour :

- récupérer un planning AniList (`GET /api/v1/anilist/airing`)
- résoudre un anime vers des candidates Anime‑Sama (`POST /api/v1/animesama/resolve`)
- créer un abonnement à partir d’une base URL (saison/langue) (`POST /api/v1/subscriptions`)
- lister les épisodes disponibles d’un abonnement (`GET /api/v1/subscriptions/{id}/episodes`)
- lancer le téléchargement d’une sélection d’épisodes (`POST /api/v1/subscriptions/{id}/enqueue`)

L’UI intègre ce flux :

- dans **Abonnements** via la section **Recherche (nom → URL)**
- dans **Calendrier → AniList** (bouton “Résoudre”, puis “Créer”)

## Résolution via API (exemple curl)

```bash
curl -sS -X POST http://127.0.0.1:8080/api/v1/animesama/resolve \
    -H 'content-type: application/json' \
    -d '{
        "titles": ["Sousou no Frieren", "Frieren: Beyond Journey\u0027s End"],
        "season": 1,
        "lang": "vostfr",
        "maxCandidates": 5
    }'
```

La réponse contient une liste de `candidates` avec : `catalogueUrl`, `baseUrl`, `matchedTitle`, `score`.

## Créer un abonnement depuis une candidate

```bash
curl -sS -X POST http://127.0.0.1:8080/api/v1/subscriptions \
    -H 'content-type: application/json' \
    -d '{
                "baseUrl": "https://anime-sama.si/catalogue/.../saison1/vostfr/"
    }'
```

`label` est optionnel (auto-généré depuis l’URL) et `player` est géré automatiquement.

## Sélection d’épisodes

Dans l’UI : bouton **Épisodes** sur une subscription → coche ce que tu veux → **Télécharger la sélection**.

Via API (exemple) :

```bash
curl -sS http://127.0.0.1:8080/api/v1/subscriptions/<id>/episodes

curl -sS -X POST http://127.0.0.1:8080/api/v1/subscriptions/<id>/enqueue \
    -H 'content-type: application/json' \
    -d '{"episodes":[1,2,3]}'
```

## Notes

- Le domaine utilisé est `anime-sama.si` (certaines URLs historiques sont normalisées vers ce domaine).
- Pour le détail des paramètres, consulte la spec : http://127.0.0.1:8080/api/v1/openapi.json

