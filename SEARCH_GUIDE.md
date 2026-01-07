# 🔍 Guide du Moteur de Recherche

## Vue d'ensemble

Le moteur de recherche intégré permet de trouver rapidement des animes sans avoir à chercher l'URL exacte sur anime-sama.tv. Il utilise un **fuzzy matching** intelligent pour trouver les meilleures correspondances.

## Fonctionnalités

### ✨ Recherche Intelligente
- **Fuzzy matching** : trouve des correspondances même avec des fautes de frappe
- **Traductions automatiques** : reconnaît les titres français et japonais
- **Scoring de pertinence** : affiche un pourcentage de correspondance
- **Cache local** : résultats instantanés (pas de scraping web)
- **AniList (par défaut)** : enrichit la recherche avec titres + synonymes pour résoudre l'URL anime-sama plus facilement

### 🌍 Traductions Supportées

Le moteur reconnaît automatiquement certaines traductions courantes :

| Recherche | Traduit en |
|-----------|------------|
| `kaiju` | Kaiju No. 8 |
| `l'attaque des titans` | Shingeki no Kyojin |
| `attaque des titans` | Shingeki no Kyojin |
| `attack on titan` | Shingeki no Kyojin |
| `demon slayer` | Kimetsu no Yaiba |
| `my hero academia` | Boku no Hero Academia |
| `promised neverland` | Yakusoku no Neverland |
| `sao` | Sword Art Online |

## Utilisation

### Mode Interactif

Lancez le programme normalement :

```bash
python main.py
```

Quand demandé, choisissez "Oui" pour utiliser la recherche :

```
🔗 ANIME-SAMA URL INPUT
────────────────────────────────────────
Use search? (Y/n): 
```

Puis tapez votre recherche :

```
🔍 ANIME SEARCH
Search the catalogue by name (French or Japanese)

Enter anime name to search: kaiju
```

Le programme affiche les résultats avec leur score :

```
Search Results:

1. Kaiju No. 8 (94% match)
2. Chainsaw Man (40% match)
3. Haikyuu (40% match)

0. Cancel

Select anime (1-3, or 0 to cancel): 1
```

### Mode CLI avec Recherche

Utilisez l'argument `-s` ou `--search` :

```bash
# Recherche simple
python main.py -s "kaiju" -e 1-5

# Recherche avec titre français
python main.py -s "l'attaque des titans" -e 1 --quick

# Recherche avec options complètes
python main.py -s "demon slayer" -e 1-10 -t --auto-mp4

# Forcer le mode local (sans AniList)
python main.py -s "kaiju" --search-provider local -e 1-5
```

## Comment l'URL est trouvée ?

Le moteur tente (dans cet ordre) :
1. **AniList** (titres + synonymes) → génère des slugs plausibles
2. Teste l'existence de `https://anime-sama.tv/catalogue/<slug>/`
3. Fallback sur la recherche fuzzy existante

Les résultats sont mis en cache dans `~/.anime-sama-downloader.json`.

### Mode CLI Rapide (Quick Search)

Le mode CLI avec `--search` utilise automatiquement la **meilleure correspondance** si le score dépasse 50% :

```bash
python main.py -s "kaiju" -e 1-3 --quick
```

Ceci va :
1. Chercher "kaiju" → trouve "Kaiju No. 8" (94%)
2. Utiliser automatiquement cet anime
3. Télécharger les épisodes 1-3 avec les paramètres optimaux

## Exemples Pratiques

### Recherche par titre partiel

```bash
# Trouve "One Punch Man"
python main.py -s "one punch" -e 1

# Trouve "Jujutsu Kaisen"
python main.py -s "jujutsu" -e 1-5

# Trouve "Sword Art Online"
python main.py -s "sao" -e 1-10
```

### Recherche avec traduction

```bash
# Titre français → trouve l'anime japonais
python main.py -s "l'attaque des titans" -e 1-3

# Titre anglais → trouve l'anime japonais
python main.py -s "demon slayer" -e 1-5
```

### Recherche floue (fuzzy)

Même avec des fautes de frappe ou titres incomplets :

```bash
# "kaij" trouve "Kaiju No. 8"
python main.py -s "kaij" -e 1

# "tokyo ghol" trouve "Tokyo Ghoul"
python main.py -s "tokyo ghol" -e 1
```

## Scoring de Pertinence

Le moteur calcule un score basé sur plusieurs critères :

1. **Similarité globale** (algorithme SequenceMatcher)
2. **+30% bonus** : si le terme est contenu dans le titre
3. **+20% bonus** : si le titre commence par le terme recherché
4. **+10% par mot** : pour chaque mot qui correspond exactement

### Exemples de scores

| Recherche | Anime trouvé | Score | Raison |
|-----------|--------------|-------|--------|
| `kaiju` | Kaiju No. 8 | 94% | Traduction automatique + correspondance exacte |
| `l'attaque des titans` | Shingeki no Kyojin | 141% | Traduction + titre contient les deux |
| `one piece` | One Piece | 130% | Correspondance exacte + bonus |
| `naruto` | Naruto | 120% | Match exact + début du titre |

## Cache Local

Le moteur utilise un **cache local** contenant 55+ animes populaires :

- ✅ **Instantané** : pas d'attente réseau
- ✅ **Fiable** : pas de problèmes de scraping/403
- ✅ **Hors-ligne** : fonctionne sans connexion (pour la recherche)
- ⚠️ **Limité** : uniquement les animes du cache

### Animes dans le cache

Liste actuelle (55+ animes) :
- Kaiju No. 8
- Sword Art Online
- Shingeki no Kyojin (L'Attaque des Titans)
- One Piece, Naruto, Dragon Ball
- Demon Slayer, My Hero Academia
- Death Note, Tokyo Ghoul
- Fullmetal Alchemist, Hunter x Hunter
- Jujutsu Kaisen, Chainsaw Man
- Spy x Family, Frieren, Solo Leveling
- Blue Lock, Haikyuu, Slam Dunk
- Et bien d'autres...

## Ajouter des Animes au Cache

Pour ajouter un anime, éditez `utils/search.py` :

```python
ANIME_CACHE = [
    {"title": "Nom de l'Anime", "url": "https://anime-sama.tv/catalogue/anime-slug/"},
    # ... autres animes
]
```

## Limitations

1. **URLs de base** : retourne l'URL principale (pas la saison/langue spécifique)
2. **Réseau requis (AniList)** : le mode par défaut utilise AniList (désactivable via `--search-provider local`)

## Tips & Astuces

### Améliorer les résultats

1. **Soyez spécifique** : "kaiju 8" est mieux que "kaiju"
2. **Utilisez les traductions** : le système connaît les titres courants
3. **Essayez plusieurs variantes** : "sao", "sword art", "sword art online"

### Mode rapide recommandé

Pour une expérience optimale :

```bash
python main.py -s "VOTRE_RECHERCHE" -e EPISODES --quick
```

Ceci combine :
- ✅ Recherche automatique
- ✅ Meilleure correspondance
- ✅ Paramètres optimaux
- ✅ Pas de questions supplémentaires

### Vérifier avant de télécharger

Si vous voulez voir les résultats sans télécharger, utilisez Python :

```bash
python -c "from utils.search import search_anime; results = search_anime('kaiju'); [print(f'{i+1}. {r[\"title\"]} - {int(r[\"score\"]*100)}%') for i, r in enumerate(results[:5])]"
```

## Dépannage

### "No good match found"

Le score est trop faible (<50%). Solutions :

1. Utilisez le **mode interactif** pour voir tous les résultats
2. Essayez un nom plus **spécifique**
3. Vérifiez si l'anime est dans le **cache** (voir liste ci-dessus)
4. Ajoutez l'anime au cache si nécessaire

### Résultats inattendus

Le fuzzy matching peut donner des résultats surprenants. Utilisez :

```bash
# Mode interactif pour choisir manuellement
python main.py
# Puis rechercher et sélectionner dans la liste
```

## Exemples Complets

### Télécharger Kaiju No. 8 épisodes 1-10

```bash
# Recherche automatique + téléchargement rapide
python main.py -s "kaiju" -e 1-10 --quick

# Avec threading pour plus de vitesse
python main.py -s "kaiju" -e 1-10 -t --ts-threaded --auto-mp4
```

### Découvrir un nouvel anime

```bash
# Mode interactif pour explorer
python main.py

# Répondre "Y" à "Use search?"
# Taper le nom partiel
# Voir tous les résultats avec scores
# Choisir celui qui vous intéresse
```

### Télécharger L'Attaque des Titans

```bash
# Titre français reconnu automatiquement
python main.py -s "l'attaque des titans" -e 1-5 --quick

# Ou titre anglais
python main.py -s "attack on titan" -e 1-5 --quick

# Ou titre japonais
python main.py -s "shingeki no kyojin" -e 1-5 --quick
```

## Intégration avec autres outils

### Avec l'alias shell

Si vous avez configuré `anime-dl` :

```bash
anime-dl -s "kaiju" -e 1-10 --quick
```

### Script pour téléchargements en série

```bash
#!/bin/bash
# download_series.sh
anime-dl -s "$1" -e "$2" --quick
```

Usage :
```bash
./download_series.sh "kaiju" "1-10"
./download_series.sh "demon slayer" "1-26"
```

## Version et Compatibilité

- **Disponible depuis** : v2.6-optimized
- **Nécessite** : Python 3.6+
- **Cache** : `~/.anime-sama-downloader.json`

---

🎌 **Note** : Le moteur de recherche est optimisé pour une utilisation rapide et simple. Pour des besoins avancés ou des animes non présents dans le cache, utilisez directement l'URL anime-sama.tv.
