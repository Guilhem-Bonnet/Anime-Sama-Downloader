# 🎯 GUIDE D'EXÉCUTION — Pour le Product Manager

**Projet** : Anime-Sama Downloader v1.0 Refonte  
**Date** : 31 janvier 2026  
**Durée** : 10-13 semaines (2,5-3 mois)

---

## 📋 TU ES ICI : RÔLE DU PM

En tant que **Product Manager**, tu es responsable de :
1. ✅ **Piloter l'exécution** du projet selon le plan défini
2. ✅ **Coordonner les équipes** (dev, design, QA)
3. ✅ **Prioriser le backlog** et ajuster si besoin
4. ✅ **Suivre l'avancement** (sprints, burndown)
5. ✅ **Communiquer avec le Product Owner** (Guilhem)
6. ✅ **Garantir la qualité** des livrables

---

## 📚 ÉTAPE 1 : LIRE TOUTE LA DOCUMENTATION

### Documents à lire dans l'ordre :

1. **[00-PROJECT-BRIEF.md](00-PROJECT-BRIEF.md)** ← **COMMENCE ICI**
   - Vue d'ensemble projet
   - Objectifs, personas, roadmap
   - Durée : 10 min

2. **[01-PERSONAS-AND-JOURNEYS.md](01-PERSONAS-AND-JOURNEYS.md)**
   - Comprendre les utilisateurs cibles (Alex & Maya)
   - User journeys avant/après
   - Matrice de décision features
   - Durée : 15 min

3. **[02-DESIGN-SYSTEM-SAKURA-NIGHT.md](02-DESIGN-SYSTEM-SAKURA-NIGHT.md)**
   - Charte graphique complète
   - Composants UI, couleurs, typographie
   - Mockups wireframes
   - Durée : 20 min

4. **[03-TECHNICAL-ARCHITECTURE.md](03-TECHNICAL-ARCHITECTURE.md)**
   - Architecture backend (Go) + frontend (React)
   - API endpoints, database schema
   - Sécurité, performance
   - Durée : 30 min

5. **[04-BACKLOG-SPRINTS.md](04-BACKLOG-SPRINTS.md)** ← **LE PLUS IMPORTANT**
   - Backlog détaillé par sprint
   - User stories avec acceptance criteria
   - Story points estimés
   - Durée : 45 min

**⏱️ Temps total de lecture** : ~2h

---

## 🎬 ÉTAPE 2 : SETUP PROJET

### 2.1 GitHub Issues & Milestones

**Action** : Créer les milestones + issues sur GitHub

1. Crée 4 **milestones** GitHub :
   - `Sprint 1 - Fondations` (Deadline : +2 semaines)
   - `Sprint 2 - UI/UX Refonte` (Deadline : +5 semaines)
   - `Sprint 3 - Features Avancées` (Deadline : +8 semaines)
   - `Sprint 4 - Polish & Release` (Deadline : +10 semaines)

2. Pour **chaque task** dans [04-BACKLOG-SPRINTS.md](04-BACKLOG-SPRINTS.md) :
   - Crée une **GitHub Issue**
   - Titre : `[TASK-XXX] Nom de la tâche`
   - Description : Copie Story Points + Acceptance Criteria
   - Labels : `backend`, `frontend`, `infrastructure`, `design`
   - Assignee : À déterminer (voir section 3)
   - Milestone : Sprint correspondant

**Exemple GitHub Issue** :

```markdown
**Story Points** : 8

**Description**
Découper App.tsx (1922 lignes) en pages + composants réutilisables.

**Acceptance Criteria**
- [ ] Pages créées (Dashboard, Search, Subscriptions, Calendar, Jobs, Settings)
- [ ] Composants UI extraits (Header, TabNav, Card, Button, Input)
- [ ] App.tsx < 200 lignes (juste routing)
- [ ] Aucune régression fonctionnelle

**Sprint** : Sprint 1
```

**🛠️ Outil recommandé** : Script automation pour créer issues en masse (GitHub CLI ou API)

### 2.2 Project Board

**Action** : Créer un **GitHub Project Board** (Kanban)

Colonnes :
- 📥 **Backlog** (toutes les issues non démarrées)
- 🏃 **Sprint actuel** (issues du sprint en cours)
- 🚧 **In Progress** (en développement)
- 👀 **In Review** (en code review)
- ✅ **Done** (terminées)

### 2.3 Communication Channels

**Action** : Setup canaux de communication

- **Slack/Discord** : Channel `#asd-refonte`
- **Daily standups** : 15 min tous les matins (async si remote)
- **Sprint Planning** : Début de chaque sprint (2h)
- **Sprint Review** : Fin de chaque sprint (1h)
- **Sprint Retro** : Après review (30 min)

---

## 👥 ÉTAPE 3 : CONSTITUER L'ÉQUIPE

### Rôles nécessaires

| Rôle | Responsabilités | Nombre |
|------|----------------|--------|
| **Tech Lead** | Architecture, code reviews, décisions techniques | 1 |
| **Frontend Dev** | React, TypeScript, UI components | 1-2 |
| **Backend Dev** | Go, APIs, database, workers | 1-2 |
| **UX/UI Designer** | Mockups, design system, animations | 1 |
| **QA Engineer** | Tests manuels + auto, E2E | 1 |
| **Product Manager** | Toi ! Pilotage, priorités, communication | 1 |

**Équipe idéale** : 5-7 personnes

### Assigner les tâches

**Répartition recommandée** :

**Sprint 1** :
- **Backend Dev** : TASK-001 à TASK-007 (audit, tests, refactor)
- **Frontend Dev** : TASK-101 à TASK-107 (découpage, design system, Zustand)
- **DevOps** : TASK-201 à TASK-204 (Docker, CI/CD)

**Sprint 2** :
- **UX/UI Designer** : TASK-301 à TASK-302 (thème, composants)
- **Frontend Dev #1** : TASK-303, TASK-304, TASK-305 (pages Dashboard, Search, Subs)
- **Frontend Dev #2** : TASK-306, TASK-307, TASK-308 (pages Calendar, Jobs, Settings)
- **Frontend Dev (tous)** : TASK-309, TASK-310 (animations, responsive)

**Sprint 3** :
- **Backend Dev #1** : TASK-401, TASK-402, TASK-403 (auth, users, AniList)
- **Backend Dev #2** : TASK-404, TASK-405, TASK-406 (scheduler, Jellyfin, naming)
- **Backend Dev** : TASK-407 (notifications)
- **Frontend Dev** : TASK-408 à TASK-412 (auth UI, profil, AniList import UI)

**Sprint 4** :
- **Toute l'équipe** : TASK-501 à TASK-508 (perf, a11y, doc, E2E, release)

---

## 📊 ÉTAPE 4 : SPRINT PLANNING #1

**Quand** : Avant de commencer Sprint 1  
**Durée** : 2 heures  
**Participants** : Toute l'équipe

### Agenda

1. **Présentation projet** (15 min)
   - Toi (PM) présente le brief + personas + objectifs
   - Montre mockups design system

2. **Review backlog Sprint 1** (30 min)
   - Passe en revue chaque task (TASK-001 à TASK-204)
   - Équipe pose questions
   - Tech Lead valide faisabilité

3. **Estimation effort** (30 min)
   - Équipe valide/ajuste story points
   - Planning Poker si besoin

4. **Assignation tâches** (30 min)
   - Chaque dev prend 2-3 tasks
   - Équilibre charge de travail

5. **Définition Sprint Goal** (15 min)
   - **Sprint 1 Goal** : "Codebase propre, tests 70%+, design system prêt"

### Output

- ✅ Toutes les issues Sprint 1 assignées
- ✅ Sprint Goal défini et partagé
- ✅ Équipe comprend les priorités

---

## 📈 ÉTAPE 5 : SUIVRE L'AVANCEMENT

### Daily Standups (15 min/jour)

Format :
1. Chaque personne répond :
   - ✅ Qu'ai-je fait hier ?
   - 🎯 Que vais-je faire aujourd'hui ?
   - 🚧 Y a-t-il des blockers ?

2. PM note les blockers → résout dans la journée

### Outils de suivi

**GitHub Project Board** :
- Update colonnes en temps réel
- Burndown chart (story points restants)

**Métriques à tracker** :
- **Vélocité** : Story points complétés par sprint
- **Burndown** : Graphe SP restants vs temps
- **Cycle time** : Temps moyen In Progress → Done
- **Blockers** : Nombre + durée moyenne résolution

**Dashboard recommandé** :
```
Sprint 1 - Jour 5/10
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━ 50%
✅ Done: 45 SP
🚧 In Progress: 20 SP
📥 To Do: 24 SP

🚀 Vélocité projetée: 40 SP/sprint
🎯 On track pour Sprint Goal
```

---

## 🎭 ÉTAPE 6 : COORDONNER LES ÉQUIPES

### Communication PM ↔ Équipes

**Avec les devs** :
- Daily standups (résoudre blockers)
- Code reviews (valider qualité)
- Questions techniques → escalade Tech Lead

**Avec le designer** :
- Review mockups (Sprint 2)
- Feedback itératif
- Validation finale avant dev

**Avec QA** :
- Définir stratégie tests (Sprint 1)
- Review test plans
- Validation sprints (regression testing)

**Avec Guilhem (Product Owner)** :
- Démo fin de sprint (review)
- Validation priorités si changements
- Feedback utilisateur

### Résoudre les blockers

**Types de blockers courants** :
1. **Technique** : Bug bloquant, API down
   → Escalade Tech Lead, trouve workaround

2. **Dépendances** : Task A bloque Task B
   → Reprio, parallélise si possible

3. **Scope creep** : Nouvelle feature demandée
   → Discute avec PO, ajoute au backlog futur

4. **Ressources** : Quelqu'un malade/indispo
   → Réassigne tâches, ajuste sprint

---

## 🏁 ÉTAPE 7 : SPRINT REVIEWS & RETROS

### Sprint Review (fin de sprint)

**Durée** : 1 heure  
**Participants** : Équipe + Guilhem (PO)

**Agenda** :
1. **Démo** (30 min)
   - Chaque dev montre ce qu'il a fait
   - Live demo si possible

2. **Metrics** (10 min)
   - Vélocité sprint
   - Story points done vs planned

3. **Feedback PO** (20 min)
   - Guilhem valide livrables
   - Ajustements si besoin

### Sprint Retro (après review)

**Durée** : 30 min  
**Participants** : Équipe seulement (sans PO)

**Format** :
1. **What went well?** ✅
2. **What could be improved?** ⚠️
3. **Action items** 🎯

**Output** : 2-3 actions concrètes pour le prochain sprint

---

## 📋 ÉTAPE 8 : CHECKLIST PAR SPRINT

### Sprint 1 Checklist

- [ ] Sprint Planning #1 fait
- [ ] Toutes les issues créées + assignées
- [ ] Daily standups tous les jours
- [ ] Code reviews pour chaque PR
- [ ] Tests coverage ≥ 70% atteint
- [ ] Design system base prêt
- [ ] Sprint Review + Retro
- [ ] Démo à Guilhem

### Sprint 2 Checklist

- [ ] Sprint Planning #2 fait
- [ ] Mockups validés par Guilhem
- [ ] Thème "Sakura Night" implémenté
- [ ] Toutes les pages refaites
- [ ] UX : Télécharger 1 anime en ≤ 3 clics validé
- [ ] Responsive testing (mobile + desktop)
- [ ] Sprint Review + Retro
- [ ] Démo à Guilhem

### Sprint 3 Checklist

- [ ] Sprint Planning #3 fait
- [ ] Multi-users + auth fonctionnel
- [ ] AniList auto-import testé
- [ ] Jellyfin integration testée (avec instance réelle)
- [ ] Scheduler hebdo configuré
- [ ] Tests intégration passent
- [ ] Sprint Review + Retro
- [ ] Démo à Guilhem

### Sprint 4 Checklist

- [ ] Sprint Planning #4 fait
- [ ] Lighthouse score ≥ 90
- [ ] Accessibilité WCAG AA validée
- [ ] Documentation complète (user + dev)
- [ ] E2E tests passent
- [ ] Security audit fait
- [ ] Release v1.0 créée (Git tag + Docker images)
- [ ] Sprint Review + Retro FINAL
- [ ] 🎉 **RELEASE PARTY!**

---

## 🎁 ÉTAPE 9 : LIVRABLES FINAUX

À la fin du projet, tu dois avoir :

### Code
- ✅ Backend Go refactorisé + tests 70%+
- ✅ Frontend React modulaire + design system
- ✅ Multi-users + auth JWT
- ✅ Intégration Jellyfin fonctionnelle
- ✅ CI/CD (GitHub Actions)

### Documentation
- ✅ README complet
- ✅ QUICK_START.md
- ✅ ARCHITECTURE.md
- ✅ API.md (OpenAPI)
- ✅ CONTRIBUTING.md
- ✅ CHANGELOG.md

### Release
- ✅ Git tag v1.0.0
- ✅ Docker images published
- ✅ GitHub Release avec binaries
- ✅ Migration guide (v0.x → v1.0)

### Metrics
- ✅ Lighthouse ≥ 90
- ✅ Test coverage ≥ 70%
- ✅ WCAG AA compliance
- ✅ Vélocité : 35-40 SP/sprint

---

## 🚨 RISQUES & MITIGATION

| Risque | Probabilité | Impact | Mitigation |
|--------|-------------|--------|------------|
| Anime-sama change structure HTML | 🟠 Moyen | 🔴 Haut | Abstraire scraper, ajouter tests, monitor |
| Scope creep (nouvelles features) | 🟠 Moyen | 🟠 Moyen | Backlog futur, ne pas toucher MVP |
| Dev malade/indispo | 🟡 Faible | 🟠 Moyen | Cross-training, doc claire |
| Jellyfin API change | 🟡 Faible | 🟠 Moyen | Tester avec plusieurs versions |
| Performance issues (scale) | 🟡 Faible | 🟡 Faible | Load testing Sprint 4 |

---

## 📞 CONTACTS & ESCALATION

**Product Owner (Guilhem)** : Décisions produit, validation features  
**Tech Lead** : Décisions techniques, architecture  
**PM (toi)** : Coordination, priorités, blockers

**Escalation path** :
1. Blocker dev → Tech Lead (< 1h)
2. Blocker Tech Lead → PM (< 4h)
3. Décision produit → Guilhem (< 24h)

---

## 🎯 TL;DR — LES 10 COMMANDEMENTS DU PM

1. ✅ **Lire toute la doc** avant de commencer
2. ✅ **Créer issues GitHub** depuis le backlog
3. ✅ **Constituer l'équipe** (5-7 personnes)
4. ✅ **Sprint Planning** au début de chaque sprint
5. ✅ **Daily standups** (15 min/jour)
6. ✅ **Résoudre blockers** dans la journée
7. ✅ **Code reviews** obligatoires (qualité++)
8. ✅ **Sprint Review + Retro** à chaque fin de sprint
9. ✅ **Communiquer avec Guilhem** (validation)
10. ✅ **Célébrer les wins** 🎉

---

## 📚 ANNEXE : DOCUMENTS À DISTRIBUER

### Qui reçoit quoi ?

**Toute l'équipe** :
- [00-PROJECT-BRIEF.md](00-PROJECT-BRIEF.md)
- [01-PERSONAS-AND-JOURNEYS.md](01-PERSONAS-AND-JOURNEYS.md)

**Tech Lead + Devs** :
- [03-TECHNICAL-ARCHITECTURE.md](03-TECHNICAL-ARCHITECTURE.md)
- [04-BACKLOG-SPRINTS.md](04-BACKLOG-SPRINTS.md)

**UX/UI Designer** :
- [01-PERSONAS-AND-JOURNEYS.md](01-PERSONAS-AND-JOURNEYS.md)
- [02-DESIGN-SYSTEM-SAKURA-NIGHT.md](02-DESIGN-SYSTEM-SAKURA-NIGHT.md)

**QA Engineer** :
- [04-BACKLOG-SPRINTS.md](04-BACKLOG-SPRINTS.md) (acceptance criteria)
- [03-TECHNICAL-ARCHITECTURE.md](03-TECHNICAL-ARCHITECTURE.md) (API endpoints pour tests)

**Product Owner (Guilhem)** :
- [00-PROJECT-BRIEF.md](00-PROJECT-BRIEF.md)
- Ce document ([05-PM-EXECUTION-GUIDE.md](05-PM-EXECUTION-GUIDE.md))

---

## 🎬 PRÊT À DÉMARRER ?

### Prochaines actions immédiates :

1. ✅ Lire tous les docs (2h)
2. ✅ Créer GitHub issues (1h)
3. ✅ Setup Project Board (30 min)
4. ✅ Constituer l'équipe (1-2 jours)
5. ✅ Planifier Sprint Planning #1 (date + heure)
6. 🚀 **GO GO GO!**

---

**Bonne chance et que la refonte soit avec toi !** 🌸✨

*— L'équipe BMAD (Mary, Winston, Bob, Sally, Carson, Amelia)*
