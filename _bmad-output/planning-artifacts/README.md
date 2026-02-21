# 📦 INDEX — Documentation Projet ASD v1.0 Refonte

**Projet** : Anime-Sama Downloader v1.0  
**Date** : 31 janvier 2026  
**Status** : Planning complet ✅

---

## 🎯 POINT D'ENTRÉE PRINCIPAL

### Pour le Product Manager (Guilhem ou délégué)

👉 **COMMENCE ICI** : [05-PM-EXECUTION-GUIDE.md](05-PM-EXECUTION-GUIDE.md)

Ce guide te dit **exactement** :
- Quoi lire et dans quel ordre
- Comment setup le projet (GitHub, équipe, sprints)
- Comment piloter l'exécution au quotidien
- Checklist complète par sprint

**Temps de lecture** : 20 min  
**Temps de setup** : 2-3h

---

## 📚 DOCUMENTATION COMPLÈTE

### Documents stratégiques

| Document | Description | Lecteurs | Durée |
|----------|-------------|----------|-------|
| [00-PROJECT-BRIEF.md](00-PROJECT-BRIEF.md) | Vue d'ensemble projet (objectifs, vision, roadmap) | Tous | 10 min |
| [01-PERSONAS-AND-JOURNEYS.md](01-PERSONAS-AND-JOURNEYS.md) | Personas utilisateurs + user journeys | Tous | 15 min |
| [05-PM-EXECUTION-GUIDE.md](05-PM-EXECUTION-GUIDE.md) | **Guide d'exécution PM** (le plus important) | PM | 20 min |

### Documents techniques

| Document | Description | Lecteurs | Durée |
|----------|-------------|----------|-------|
| [03-TECHNICAL-ARCHITECTURE.md](03-TECHNICAL-ARCHITECTURE.md) | Architecture backend + frontend, API, database | Tech Lead, Devs | 30 min |
| [04-BACKLOG-SPRINTS.md](04-BACKLOG-SPRINTS.md) | Backlog détaillé avec user stories | Tous, surtout Devs | 45 min |

### Documents design

| Document | Description | Lecteurs | Durée |
|----------|-------------|----------|-------|
| [02-DESIGN-SYSTEM-SAKURA-NIGHT.md](02-DESIGN-SYSTEM-SAKURA-NIGHT.md) | Charte graphique complète + composants UI | Designer, Frontend | 20 min |

---

## 🎭 QUI LIT QUOI ?

### Product Owner (Guilhem)
1. ✅ [00-PROJECT-BRIEF.md](00-PROJECT-BRIEF.md) — Vue d'ensemble
2. ✅ [05-PM-EXECUTION-GUIDE.md](05-PM-EXECUTION-GUIDE.md) — Comment piloter
3. 📌 Validation à chaque Sprint Review

### Product Manager (Délégué)
1. ✅ **TOUS LES DOCUMENTS** dans l'ordre (2h total)
2. 🎯 Focus sur [05-PM-EXECUTION-GUIDE.md](05-PM-EXECUTION-GUIDE.md)
3. 📊 Setup GitHub + équipe selon instructions

### Tech Lead / Architect
1. ✅ [00-PROJECT-BRIEF.md](00-PROJECT-BRIEF.md)
2. ✅ [03-TECHNICAL-ARCHITECTURE.md](03-TECHNICAL-ARCHITECTURE.md) — Architecture complète
3. ✅ [04-BACKLOG-SPRINTS.md](04-BACKLOG-SPRINTS.md) — Backlog technique
4. 🛠️ Valider faisabilité, ajuster estimations si besoin

### Backend Developers
1. ✅ [00-PROJECT-BRIEF.md](00-PROJECT-BRIEF.md)
2. ✅ [03-TECHNICAL-ARCHITECTURE.md](03-TECHNICAL-ARCHITECTURE.md) — Section Backend
3. ✅ [04-BACKLOG-SPRINTS.md](04-BACKLOG-SPRINTS.md) — Tasks TASK-001 à TASK-007, TASK-401 à TASK-407
4. 💻 Implémenter selon acceptance criteria

### Frontend Developers
1. ✅ [00-PROJECT-BRIEF.md](00-PROJECT-BRIEF.md)
2. ✅ [01-PERSONAS-AND-JOURNEYS.md](01-PERSONAS-AND-JOURNEYS.md) — Comprendre utilisateurs
3. ✅ [02-DESIGN-SYSTEM-SAKURA-NIGHT.md](02-DESIGN-SYSTEM-SAKURA-NIGHT.md) — Design system
4. ✅ [03-TECHNICAL-ARCHITECTURE.md](03-TECHNICAL-ARCHITECTURE.md) — Section Frontend
5. ✅ [04-BACKLOG-SPRINTS.md](04-BACKLOG-SPRINTS.md) — Tasks TASK-101 à TASK-107, TASK-301 à TASK-412
6. 🎨 Implémenter selon design system

### UX/UI Designer
1. ✅ [01-PERSONAS-AND-JOURNEYS.md](01-PERSONAS-AND-JOURNEYS.md) — Personas + journeys
2. ✅ [02-DESIGN-SYSTEM-SAKURA-NIGHT.md](02-DESIGN-SYSTEM-SAKURA-NIGHT.md) — Charte graphique
3. 🎨 Créer mockups haute-fidélité (Sprint 2)
4. 🤝 Collaboration avec Frontend Devs

### QA Engineer
1. ✅ [00-PROJECT-BRIEF.md](00-PROJECT-BRIEF.md)
2. ✅ [01-PERSONAS-AND-JOURNEYS.md](01-PERSONAS-AND-JOURNEYS.md) — User journeys = test scenarios
3. ✅ [04-BACKLOG-SPRINTS.md](04-BACKLOG-SPRINTS.md) — Acceptance criteria = tests
4. 🧪 Écrire test plans + tests E2E

---

## 📅 PLANNING RECAP

| Sprint | Durée | Objectif | Story Points |
|--------|-------|----------|--------------|
| **Sprint 1** | 2 sem | Fondations & Cleanup | 89 SP |
| **Sprint 2** | 3 sem | Redesign UI/UX | 103 SP |
| **Sprint 3** | 3 sem | Features Avancées | 110 SP |
| **Sprint 4** | 2 sem | Polish & Release | 60 SP |
| **TOTAL** | 10 sem | | **362 SP** |

---

## 🎯 MÉTRIQUES DE SUCCÈS

### Technique
- ✅ Test coverage ≥ 70%
- ✅ Lighthouse score ≥ 90
- ✅ WCAG AA compliance
- ✅ Zero critical security issues

### Produit
- ✅ UX : Télécharger 1 anime en ≤ 3 clics (persona Alex)
- ✅ Automation : AniList sync → subscriptions auto (persona Maya)
- ✅ Jellyfin integration fonctionnelle
- ✅ Multi-users avec auth JWT

### Business
- ✅ Release v1.0 dans 10-13 semaines
- ✅ Documentation complète (user + dev)
- ✅ Feedback positif des 2 personas

---

## 🚀 QUICK START

### Si tu es le PM :

```bash
# 1. Lis le guide d'exécution (20 min)
open 05-PM-EXECUTION-GUIDE.md

# 2. Lis tous les autres docs (2h)
# (dans l'ordre indiqué dans le guide)

# 3. Setup GitHub (1h)
# - Créer milestones
# - Créer issues depuis backlog
# - Setup Project Board

# 4. Constituer l'équipe (1-2 jours)

# 5. Sprint Planning #1
# Date : ________________
# Heure : _______________
# Participants : Toute l'équipe

# 6. GO GO GO! 🚀
```

### Si tu es un Dev :

```bash
# 1. Lis les docs pour ton rôle (voir section "Qui lit quoi?")

# 2. Attends Sprint Planning #1
# Le PM te contactera avec date + heure

# 3. Prends tes tasks assignées

# 4. Code + tests + PR

# 5. Profit! 💰
```

---

## 📞 SUPPORT

### Questions ?

- **Produit / Priorités** → Product Owner (Guilhem)
- **Technique / Architecture** → Tech Lead
- **Coordination / Blockers** → Product Manager
- **Design / UX** → UX Designer

### Documents manquants ?

Tous les documents sont dans :
```
_bmad-output/planning-artifacts/
├── 00-PROJECT-BRIEF.md
├── 01-PERSONAS-AND-JOURNEYS.md
├── 02-DESIGN-SYSTEM-SAKURA-NIGHT.md
├── 03-TECHNICAL-ARCHITECTURE.md
├── 04-BACKLOG-SPRINTS.md
├── 05-PM-EXECUTION-GUIDE.md
└── README.md (ce fichier)
```

---

## 🎉 ON EST PRÊTS !

Toute la documentation nécessaire est créée. L'équipe peut démarrer dès que :

1. ✅ PM a lu tous les docs
2. ✅ GitHub issues créées
3. ✅ Équipe constituée
4. ✅ Sprint Planning #1 planifié

**Prochaine étape** : Le PM suit [05-PM-EXECUTION-GUIDE.md](05-PM-EXECUTION-GUIDE.md) pas à pas.

---

**Bonne chance pour cette refonte épique !** 🌸✨

*— L'équipe BMAD Party Mode*  
*Mary, Winston, Bob, Sally, Carson, Amelia, Quinn, John, et BMad Master*
