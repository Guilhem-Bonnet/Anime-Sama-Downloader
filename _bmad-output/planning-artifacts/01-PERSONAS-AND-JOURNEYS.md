# 👥 PERSONAS & USER JOURNEYS

**Projet** : Anime-Sama Downloader v1.0  
**Date** : 31 janvier 2026  
**Auteurs** : Mary (Analyst) + Sally (UX Designer)

---

## 🍿 PERSONA 1 : Alex le Casual Fan

### Profil

**Démographie**
- Âge : 22-28 ans
- Occupation : Étudiant en informatique / Jeune développeur
- Localisation : France métropolitaine
- Équipement : PC Windows/Mac, connaissances Docker basiques

**Comportements**
- Regarde 3-5 animes simultanément (séries en cours)
- Préfère VOSTFR (sous-titres français)
- Regarde en soirée (19h-23h)
- Partage rarement (usage perso uniquement)

**Compétences techniques**
- Niveau : ⭐⭐☆☆☆ (Basique)
- Sait installer Docker Desktop
- Copier-coller des commandes terminal
- Éditer un fichier `.env`
- **Ne sait PAS** : Éditer config avancée, debug logs

**Objectifs**
1. Télécharger les derniers épisodes de ses animes préférés **rapidement**
2. Être notifié quand un nouvel épisode est disponible
3. Organiser ses téléchargements simplement

**Frustrations actuelles**
- ❌ "Je ne comprends pas ce qu'est un 'baseURL' ou un 'player'"
- ❌ "Pourquoi y a-t-il plusieurs groupes de sous-titres ? Lequel est le meilleur ?"
- ❌ "L'interface est moche et confuse, je clique partout sans comprendre"
- ❌ "Impossible de savoir si le téléchargement a réussi ou échoué"

**Citations**
> "Je veux juste télécharger *Demon Slayer* épisode 12 en 3 clics, pas passer 10 minutes à comprendre l'interface."

> "Pourquoi ça ressemble à un logiciel d'entreprise des années 2000 alors qu'on parle d'anime ?"

---

### User Journey : Télécharger un nouvel épisode

#### Scénario
Alex vient d'apprendre qu'un nouvel épisode de **Jujutsu Kaisen S2** est sorti. Il veut le télécharger en VOSTFR.

#### Étapes (État actuel — problématique)

| Étape | Action | Ressenti | Problème |
|-------|--------|----------|----------|
| 1 | Ouvre l'app | 😐 Neutre | Interface austère, pas engageante |
| 2 | Cherche l'anime | 😕 Confus | "Je tape le titre... où ?" |
| 3 | Voit plusieurs résultats | 😰 Stressé | "Season 1, 2, VF, VOSTFR... lequel choisir ?" |
| 4 | Sélectionne épisode | 😡 Frustré | "Range input + champ texte, c'est quoi la différence ?" |
| 5 | Clique "Télécharger" | 🤔 Incertain | "Ça a marché ? Où est-ce que je vois la progression ?" |
| 6 | Attend | 😴 Ennuyé | Aucun feedback visuel, juste un spinner |
| 7 | Vérifie le dossier | 😤 Agacé | "Fichier nommé bizarrement, impossible à retrouver" |

**🎯 Objectif de la refonte** : Ramener ce parcours à **3 clics + 1 notification**

#### Étapes (Nouveau design — cible)

| Étape | Action | Ressenti | Solution |
|-------|--------|----------|----------|
| 1 | Ouvre l'app | 😍 Émerveillé | **Dashboard visuel** avec cards anime stylisées |
| 2 | Voit "Nouveautés" | 🎉 Excité | **Card "Jujutsu Kaisen S2 Ep12 disponible"** en highlight |
| 3 | Clique sur la card | ✅ Confiant | Modal avec preview + metadata (durée, qualité, groupe) |
| 4 | Clique "Télécharger" | 😌 Serein | **Confirmation visuelle + ajout à la queue** |
| 5 | Voit progression | 📊 Informé | **Progress bar animée + ETA** dans section Jobs |
| 6 | Reçoit notification | 🔔 Satisfait | **Toast "Épisode prêt !"** + lien direct vers fichier |

**Temps total** : < 30 secondes  
**Clics** : 3  
**Satisfaction** : 9/10

---

## 📚 PERSONA 2 : Maya la Jellyfin Power User

### Profil

**Démographie**
- Âge : 28-40 ans
- Occupation : DevOps / SysAdmin / Développeur senior
- Localisation : France / Belgique / Suisse
- Équipement : NAS Synology/TrueNAS, serveur Linux dédié, Jellyfin/Plex

**Comportements**
- Suit 20-30 animes simultanément (dont 10+ en cours)
- Archive tout (2To+ de médiathèque anime)
- Automatise tout (cron jobs, webhooks, scripts)
- Partage avec famille (3-5 utilisateurs Jellyfin)

**Compétences techniques**
- Niveau : ⭐⭐⭐⭐⭐ (Expert)
- Maîtrise Docker Compose, Linux, APIs REST
- Écrit des scripts Bash/Python
- Connaît metadata (AniDB, AniList, TVDB)

**Objectifs**
1. Automatiser complètement les téléchargements (sync AniList → subscriptions)
2. Intégrer proprement à Jellyfin (naming, metadata, scan auto)
3. Gérer plusieurs saisons/langues par anime
4. Monitoring et alertes (download failures, espace disque)

**Frustrations actuelles**
- ❌ "Je dois créer manuellement chaque abonnement, pas de sync AniList watchlist"
- ❌ "Les fichiers sont mal nommés, Jellyfin ne les reconnaît pas"
- ❌ "Impossible de gérer 20 animes proprement, l'interface saturée"
- ❌ "Pas de webhooks, je dois polling pour intégrer à mon stack"
- ❌ "Un seul utilisateur = ma femme ne peut pas avoir sa propre liste"

**Citations**
> "Je veux que ça 'just works' : je mets un anime dans ma watchlist AniList, il se télécharge automatiquement et apparaît dans Jellyfin."

> "L'UI peut être simple pour les débutants, mais donnez-moi un mode 'expert' avec API complète et webhooks."

---

### User Journey : Automatiser sa médiathèque

#### Scénario
Maya découvre 5 nouveaux animes dans sa watchlist AniList. Elle veut qu'ils se téléchargent automatiquement chaque semaine et apparaissent dans Jellyfin.

#### Étapes (État actuel — problématique)

| Étape | Action | Ressenti | Problème |
|-------|--------|----------|----------|
| 1 | Ajoute animes sur AniList | ✅ OK | (AniList fonctionne bien) |
| 2 | Ouvre ASD | 😐 Neutre | "Bon, faut tout configurer manuellement..." |
| 3 | Cherche chaque anime | 😒 Résigné | 5× recherche + sélection manuelle |
| 4 | Crée 5 abonnements | 😤 Frustré | "Pourquoi pas d'import auto ?" |
| 5 | Configure scheduler | 🤔 Confus | "C'est où le scheduler ? Ah, pas implémenté." |
| 6 | Télécharge épisodes | ⏳ Patient | Attend que les jobs finissent |
| 7 | Vérifie Jellyfin | 😡 Énervé | "Rien n'apparaît, il faut renommer manuellement" |
| 8 | Renomme fichiers | 😫 Épuisé | Script Python custom pour fixer les noms |
| 9 | Scan Jellyfin | 🥱 Blasé | "Enfin ça marche, mais quel bordel" |

**🎯 Objectif de la refonte** : Automatiser 100% de ce workflow

#### Étapes (Nouveau design — cible)

| Étape | Action | Ressenti | Solution |
|-------|--------|----------|----------|
| 1 | Ajoute animes sur AniList | ✅ OK | (AniList inchangé) |
| 2 | Ouvre ASD | 😊 Satisfait | **Badge "5 nouveaux animes détectés dans watchlist"** |
| 3 | Clique "Importer" | 🎯 Focus | **Modal avec preview + matching auto anime-sama** |
| 4 | Valide import | ✅ Confiant | **5 abonnements créés en 1 clic** |
| 5 | Configure une fois | ⚙️ Efficient | **Settings : scheduler hebdo + webhook Jellyfin** |
| 6 | Laisse tourner | 😌 Détendu | **Jobs auto chaque semaine (cron)** |
| 7 | Reçoit notification | 📧 Informé | **Email/Discord "25 épisodes téléchargés"** |
| 8 | Ouvre Jellyfin | 🎉 Ravi | **Tout est là, bien nommé, metadata correcte** |

**Temps setup initial** : 5 minutes  
**Temps maintenance** : 0 minutes (automatique)  
**Satisfaction** : 10/10

---

## 🎯 SYNTHÈSE : PRIORITÉS PAR PERSONA

### Features critiques pour Alex 🍿
1. ✅ Interface visuelle moderne (charte "Sakura Night")
2. ✅ Dashboard avec "Nouveautés" automatique
3. ✅ Téléchargement en 3 clics max
4. ✅ Progress bars + notifications
5. ✅ Recherche simple avec suggestions

### Features critiques pour Maya 📚
1. ✅ Sync AniList watchlist → subscriptions auto
2. ✅ Scheduler hebdomadaire + cron configurable
3. ✅ Webhooks Jellyfin + naming strategy
4. ✅ Multi-users avec droits différenciés
5. ✅ API complète + documentation OpenAPI
6. ✅ Logs structurés + monitoring

---

## 📊 MATRICE DE DÉCISION FEATURES

| Feature | Alex | Maya | Priorité | Sprint |
|---------|------|------|----------|--------|
| UI moderne | ⭐⭐⭐⭐⭐ | ⭐⭐☆☆☆ | 🔴 P0 | 2 |
| Dashboard nouveautés | ⭐⭐⭐⭐⭐ | ⭐⭐⭐☆☆ | 🔴 P0 | 2 |
| Recherche UX | ⭐⭐⭐⭐⭐ | ⭐⭐☆☆☆ | 🔴 P0 | 2 |
| Progress bars | ⭐⭐⭐⭐⭐ | ⭐⭐⭐☆☆ | 🔴 P0 | 2 |
| AniList sync | ⭐☆☆☆☆ | ⭐⭐⭐⭐⭐ | 🟠 P1 | 3 |
| Scheduler | ⭐☆☆☆☆ | ⭐⭐⭐⭐⭐ | 🟠 P1 | 3 |
| Jellyfin webhooks | ⭐☆☆☆☆ | ⭐⭐⭐⭐⭐ | 🟠 P1 | 3 |
| Multi-users | ⭐☆☆☆☆ | ⭐⭐⭐⭐☆ | 🟡 P2 | 3 |
| Mobile app | ⭐⭐⭐☆☆ | ⭐☆☆☆☆ | 🟢 P3 | 4+ |

---

**Conclusion** : Ces 2 personas couvrent 95% des use cases. L'objectif est de satisfaire **Alex en priorité** (UX simple) tout en donnant les outils puissants à **Maya** (automatisation).
