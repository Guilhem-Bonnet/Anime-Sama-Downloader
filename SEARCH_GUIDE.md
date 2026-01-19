# 🔍 Guide de recherche (AniList + résolution d’URL, serveur Go)

Le serveur expose des endpoints pour :

- récupérer un planning AniList (`GET /api/v1/anilist/airing`)
- résoudre un anime vers des candidates Anime‑Sama (`POST /api/v1/animesama/resolve`)
- créer un abonnement à partir d’une base URL (saison/langue) (`POST /api/v1/subscriptions`)

L’UI intègre ce flux dans l’onglet **Calendrier → AniList** (bouton “Résoudre”, puis “Créer”).

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
        "baseUrl": "https://anime-sama.si/catalogue/.../saison1/vostfr/",
        "label": "Frieren (S1 vostfr)",
        "player": "auto"
    }'
```

## Notes

- Le domaine utilisé est `anime-sama.si` (certaines URLs historiques sont normalisées vers ce domaine).
- Pour le détail des paramètres, consulte la spec : http://127.0.0.1:8080/api/v1/openapi.json

