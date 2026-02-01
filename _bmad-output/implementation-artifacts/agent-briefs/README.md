### 🎨 [for-ux-designer.md](./for-ux-designer.md) — Sally

**Utilise pour** :
- Créer mockups haute-fidélité
- Designer composants UI (tous les états)
- Améliorer UX d'une feature existante
- Audits accessibility

**Spécialité** : Design UX/UI, wireframes, prototypage, design system

### 📚 [for-tech-writer.md](./for-tech-writer.md) — Paige

**Utilise pour** :
- Créer/mettre à jour README
- Documenter API (endpoints, schemas)
- Guides utilisateur (Quick Start, How-To)
- Troubleshooting guides

**Spécialité** : Documentation technique claire, guides utilisateur, API docs

### 🧪 [for-quinn.md](./for-quinn.md) — Quinn

**Utilise pour** :
- Générer tests E2E (Playwright) pour user journeys
- Tests API (Go) : happy path + errors
- Tests frontend (React Testing Library)
- Tests de régression

**Spécialité** : Tests rapides, E2E testing, quality assurance


### 📋 [for-pm.md](./for-pm.md) — John

**Utilise pour** :
- Sprint planning (kickoff docs)
- Daily monitoring (standup reports)
- Sprint reviews (metrics, retrospective)
- Risk management (mitigation plans)
- Stakeholder communication

**Spécialité** : Coordination projet, suivi sprints, gestion risques
# 🤖 BRIEFS AGENTS BMAD — Index

**Projet** : Anime-Sama Downloader v1.0 Refonte  
**Date** : 31 janvier 2026

---

## 📋 COMMENT UTILISER CES BRIEFS

Chaque brief est un **prompt optimisé** pour un agent BMAD spécifique. 

**Workflow** :
1. Choisis l'agent selon la tâche
2. Ouvre le brief correspondant
3. Copie le prompt et donne-le à l'agent
4. L'agent a accès à tous les docs de planning

---

## 🎯 QUEL AGENT POUR QUELLE TÂCHE ?

### Phase Planning (FAIT ✅)
- **Analyst (Mary)** : Personas, user journeys, analyse besoins
- **Architect (Winston)** : Architecture technique
- **PM (John)** : Roadmap, coordination
- **Scrum Master (Bob)** : Backlog, sprints
- **UX Designer (Sally)** : Design system, mockups
- **Brainstorming Coach (Carson)** : Direction artistique

### Phase Implémentation (À FAIRE)

| Agent | Quand l'utiliser | Brief |
|-------|------------------|-------|
| **quick-flow-solo-dev (Barry)** | Feature complète en 1 sprint, prototype rapide | [for-quick-flow.md](for-quick-flow.md) |
| **bmad-dev (Amelia)** | Implémentation story précise, tests exhaustifs | [for-bmad-dev.md](for-bmad-dev.md) |
| **architect (Winston)** | Valider architecture, refactoring complexe | [for-architect.md](for-architect.md) |
| **ux-designer (Sally)** | Mockups haute-fidélité, composants UI | [for-ux-designer.md](for-ux-designer.md) |
| **tech-writer (Paige)** | Documentation (README, API, guides) | [for-tech-writer.md](for-tech-writer.md) |
| **quinn (QA)** | Tests E2E, test plans | [for-quinn.md](for-quinn.md) |

---

## 🚀 EXEMPLES D'UTILISATION

### Exemple 1 : Implémenter une story backend

**Tâche** : TASK-002 (Augmenter test coverage à 70%+)

**Agent** : bmad-dev (Amelia)

**Prompt** :
```
Je te confie la user story TASK-002 du backlog.

Consulte :
- _bmad-output/planning-artifacts/04-BACKLOG-SPRINTS.md (story complète)
- _bmad-output/planning-artifacts/03-TECHNICAL-ARCHITECTURE.md (architecture)

Implémente selon les acceptance criteria.
Tests obligatoires avant de marquer Done.
```

### Exemple 2 : Refaire le design system complet

**Tâche** : Implémenter charte "Sakura Night"

**Agent** : quick-flow-solo-dev (Barry)

**Prompt** :
```
Feature : Implémenter le design system "Sakura Night"

Contexte : _bmad-output/planning-artifacts/02-DESIGN-SYSTEM-SAKURA-NIGHT.md

Objectif : 
- Créer tokens.css avec toutes les variables
- Créer composants UI (Button, Card, Input, etc.)
- Page de démo Storybook-like

Quick Flow : prototype fonctionnel en 1 sprint.
```

### Exemple 3 : Valider refactoring architecture

**Tâche** : Review refactor SubscriptionService

**Agent** : architect (Winston)

**Prompt** :
```
Code review architecture :

Fichier : internal/app/subscriptions.go

Contexte : 
- _bmad-output/planning-artifacts/03-TECHNICAL-ARCHITECTURE.md
- _bmad-output/planning-artifacts/04-BACKLOG-SPRINTS.md (TASK-004)

Valide :
1. Respect Clean Architecture
2. Séparation responsabilités (extraire SchedulerService)
3. Tests coverage
4. Suggestions optimisations
```

### Exemple 4 : Créer documentation API

**Tâche** : Documenter tous les endpoints REST

**Agent** : tech-writer (Paige)

**Prompt** :
```
Documentation API REST complète :

Source :
- _bmad-output/planning-artifacts/03-TECHNICAL-ARCHITECTURE.md (section API Endpoints)
- internal/adapters/httpapi/openapi.go

Produis :
- API.md (référence complète)
- Exemples curl pour chaque endpoint
- Schemas request/response
- Codes erreurs

Format : Markdown, claire, exemples concrets.
```

---

## 📂 STRUCTURE DES BRIEFS

Chaque brief contient :

```markdown
# Brief pour [Agent]

## 🎯 Mission
[Description courte]

## 📚 Documents à consulter
- Lien vers docs planning
- Lien vers code si applicable

## 🔧 Prompt type
[Template de prompt prêt à l'emploi]

## ✅ Acceptance Criteria
- [ ] Critère 1
- [ ] Critère 2

## 📦 Livrables attendus
- Fichier 1
- Fichier 2

## 🧪 Validation
Comment vérifier que c'est bien fait
```

---

## 🎬 WORKFLOW RECOMMANDÉ

### Sprint 1 : Fondations

1. **TASK-101 à TASK-107** (Frontend refactor)
   → Utilise **quick-flow-solo-dev** (Barry) pour prototyper rapidement
   → Puis **bmad-dev** (Amelia) pour polish + tests

2. **TASK-001 à TASK-007** (Backend cleanup)
   → Utilise **bmad-dev** (Amelia) directement (tests critiques)

3. **TASK-201 à TASK-204** (Infra)
   → Utilise **quick-flow-solo-dev** (Barry)

### Sprint 2 : UI/UX

1. **TASK-301 à TASK-302** (Design system)
   → Utilise **ux-designer** (Sally) pour mockups
   → Puis **quick-flow-solo-dev** (Barry) pour implémentation rapide

2. **TASK-303 à TASK-310** (Pages)
   → Utilise **quick-flow-solo-dev** (Barry) par page
   → Validation avec **ux-designer** (Sally)

### Sprint 3 : Features avancées

1. **TASK-401 à TASK-407** (Backend features)
   → Utilise **bmad-dev** (Amelia) (complexité élevée)
   → Review avec **architect** (Winston)

2. **TASK-408 à TASK-412** (Frontend features)
   → Utilise **quick-flow-solo-dev** (Barry)

### Sprint 4 : Polish

1. **TASK-501 à TASK-508** (Doc + tests)
   → **tech-writer** (Paige) pour doc
   → **quinn** (QA) pour tests E2E
   → **bmad-dev** (Amelia) pour perf + security

---

## 🎯 RÉSUMÉ : AGENT PAR PHASE

| Phase | Agent principal | Agent secondaire |
|-------|----------------|------------------|
| **Prototype rapide** | quick-flow-solo-dev | ux-designer |
| **Implémentation robuste** | bmad-dev | architect (review) |
| **Design UI** | ux-designer | quick-flow-solo-dev (impl) |
| **Documentation** | tech-writer | bmad-dev |
| **Tests** | quinn | bmad-dev |
| **Architecture review** | architect | bmad-dev |

---

**Prochaine étape** : Consulte les briefs individuels selon tes besoins.
