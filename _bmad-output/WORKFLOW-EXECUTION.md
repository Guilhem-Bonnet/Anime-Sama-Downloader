# 🚀 Workflow d'Exécution — Ordre des Agents BMAD

**Version** : 1.0  
**Date** : 31 janvier 2026

---

## 📋 TL;DR — Par où commencer ?

```bash
# 1️⃣ Démarrer Sprint 1
@pm Sprint Planning : Sprint 1

# 2️⃣ Lancer les tasks critiques (backend foundation)
@architect Review TASK-001
@bmad-dev TASK-001: Backend refactoring
@bmad-dev TASK-002: Auth JWT

# 3️⃣ Paralléliser frontend pendant que backend avance
@ux-designer Mockup Dashboard Page
@quick-flow TASK-007: Design system tokens
```

---

## 🎯 Phase 0 : Préparation (AVANT Sprint 1)

### Étape 0.1 : Validation Architecture

**Agent** : @architect (Winston)  
**Brief** : [for-architect.md](./implementation-artifacts/agent-briefs/for-architect.md)

**Prompt** :
```
🏗️ Architecture Review : Projet Global

📚 Contexte :
- PRD : _bmad-output/planning-artifacts/prd.md
- Architecture : _bmad-output/planning-artifacts/03-TECHNICAL-ARCHITECTURE.md

🎯 Objectif :
Valider architecture globale avant implémentation :
1. Backend Clean Architecture (domain/ports/adapters)
2. Frontend Atomic Design
3. SQLite schema
4. API endpoints structure
5. Dependencies externes (anime-sama, ffmpeg)

✅ Checklist :
- [ ] Clean Architecture layers respectées
- [ ] Boundaries claires (domain ne dépend pas d'adapters)
- [ ] API versioning strategy validée
- [ ] Database schema normalisé
- [ ] Risques identifiés (anime-sama scraping, ffmpeg encoding)

📦 Délivrables :
- Architecture review report (validations + alertes)
- ADRs (Architecture Decision Records) si changements nécessaires
```

**Durée estimée** : 2-3h  
**Blocker ?** : Non, mais recommandé

---

## 🏃 Sprint 1 : Foundation Backend (2 semaines, 89 points)

### Étape 1.1 : Sprint Planning

**Agent** : @pm (John)  
**Brief** : [for-pm.md](./implementation-artifacts/agent-briefs/for-pm.md)

**Prompt** :
```
📋 Sprint Planning : Sprint 1

📚 Contexte :
- Backlog : _bmad-output/planning-artifacts/04-BACKLOG-SPRINTS.md (Sprint 1)
- PM Guide : _bmad-output/planning-artifacts/05-PM-EXECUTION-GUIDE.md

🎯 Objectif :
Préparer Sprint 1 (89 story points, 2 semaines) :
1. Review 12 tasks (TASK-001 à TASK-012)
2. Identifier dependencies (TASK-002 auth bloque TASK-003)
3. Assigner agents (Amelia pour backend, Barry pour quick wins)
4. Créer kickoff doc

✅ Checklist :
- [ ] AC vérifiés pour chaque task
- [ ] TASK-002 (auth) prioritaire (bloque autres)
- [ ] Winston (architect) reviewe TASK-001
- [ ] Sprint goal : "Foundation solide (auth + backend clean)"

📦 Délivrables :
- Sprint-1-Kickoff.md
- Assignments détaillés
- Risques identifiés + mitigation
```

**Durée** : 1h  
**Blocker ?** : Non

---

### Étape 1.2 : TASK-001 — Backend Refactoring (Critical)

**Agent** : @bmad-dev (Amelia)  
**Brief** : [for-bmad-dev.md](./implementation-artifacts/agent-briefs/for-bmad-dev.md)

**Prompt** :
```
🔧 TASK-001 : Backend Refactoring (Clean Architecture)

📋 Contexte :
- Epic : Epic 1 (Foundation)
- Backlog : _bmad-output/planning-artifacts/04-BACKLOG-SPRINTS.md (TASK-001)
- Architecture : _bmad-output/planning-artifacts/03-TECHNICAL-ARCHITECTURE.md

🎯 Objectif :
Restructurer backend Go en Clean Architecture :
1. **Domain layer** : Entities (Job, Subscription, Settings)
2. **Ports layer** : Interfaces (JobService, SubscriptionService)
3. **Adapters layer** : SQLite repo, HTTP API, anime-sama scraper
4. **Migrations** : Code existant → nouvelle structure

✅ Acceptance Criteria (100% requis) :
- [ ] Domain entities pures (pas de dépendances externes)
- [ ] Ports définis (interfaces Go)
- [ ] Adapters implémentent ports
- [ ] Tests unitaires domain (70% coverage minimum)
- [ ] Tests intégration adapters (SQLite, HTTP)
- [ ] Code existant migré sans régression
- [ ] Winston (architect) a validé structure

📦 Délivrables :
- internal/domain/*.go
- internal/ports/*.go
- internal/adapters/*/*.go
- Tests (*_test.go)
- Migration doc (si breaking changes)
```

**Durée estimée** : 2-3 jours  
**Blocker ?** : OUI — Bloque TASK-004, TASK-005, TASK-006  
**Review** : Winston (architect) AVANT merge

---

### Étape 1.3 : TASK-002 — Auth JWT (Critical)

**Agent** : @bmad-dev (Amelia)  
**Brief** : [for-bmad-dev.md](./implementation-artifacts/agent-briefs/for-bmad-dev.md)

**Prompt** :
```
🔧 TASK-002 : JWT Authentication

📋 Contexte :
- Epic : Epic 8 (Multi-User Support)
- Backlog : _bmad-output/planning-artifacts/04-BACKLOG-SPRINTS.md (TASK-002)

🎯 Objectif :
Implémenter auth JWT :
1. `/auth/register` (email, password) → JWT token
2. `/auth/login` (email, password) → JWT token
3. `/auth/refresh` (refresh token) → nouveau access token
4. Middleware auth (vérifie JWT sur routes protégées)

✅ Acceptance Criteria :
- [ ] Passwords hashés (bcrypt)
- [ ] JWT avec expiry 1h (access), 7j (refresh)
- [ ] Middleware auth valide token + extrait userId
- [ ] Tests unitaires (login success, wrong password, expired token)
- [ ] Tests intégration (E2E register → login → protected route)
- [ ] Winston (architect) a validé sécurité

📦 Délivrables :
- internal/app/auth.go
- internal/adapters/httpapi/auth_handlers.go
- Tests (*_test.go)
- Documentation API (endpoints)
```

**Durée** : 1-2 jours  
**Blocker ?** : OUI — Bloque TASK-003 (RBAC), TASK-010 (frontend protected routes)  
**Review** : Winston (architect) — SÉCURITÉ CRITIQUE

---

### Étape 1.4 : Paralléliser Frontend (pendant backend)

#### TASK-007 : Design System Tokens

**Agent** : @quick-flow (Barry)  
**Brief** : [for-quick-flow.md](./implementation-artifacts/agent-briefs/for-quick-flow.md)

**Prompt** :
```
⚡ TASK-007 : Design System Tokens

📋 Contexte :
- Epic : Epic 1 (Foundation)
- Design : _bmad-output/planning-artifacts/02-DESIGN-SYSTEM-SAKURA-NIGHT.md

🎯 Objectif :
Créer tokens CSS (palette Sakura Night) :
1. Colors (primary, secondary, accent, text, bg)
2. Typography (font scales, line heights)
3. Spacing (4px scale)
4. Shadows, radius, transitions

✅ Acceptance Criteria :
- [ ] Fichier webapp/src/tokens.css
- [ ] Variables CSS --sakura-*
- [ ] Dark mode par défaut
- [ ] Import dans App.tsx

📦 Délivrables :
- webapp/src/tokens.css
- Documentation usage (README)
```

**Durée** : 2-4h  
**Blocker ?** : Non  
**Review** : Sally (ux-designer) — validation visuelle

---

#### TASK-008 : UI Components Library

**Agent** : @quick-flow (Barry)  
**Brief** : [for-quick-flow.md](./implementation-artifacts/agent-briefs/for-quick-flow.md)

**Prompt** :
```
⚡ TASK-008 : UI Components Library

📋 Contexte :
- Epic : Epic 1 (Foundation)
- Design : _bmad-output/planning-artifacts/02-DESIGN-SYSTEM-SAKURA-NIGHT.md

🎯 Objectif :
Créer composants UI de base :
1. Button (variants: primary, secondary, danger)
2. Input (text, password)
3. Card (container)
4. Badge (status)

✅ Acceptance Criteria :
- [ ] Composants dans webapp/src/components/ui/
- [ ] Props TypeScript interfaces
- [ ] Utilise tokens CSS
- [ ] Tests React Testing Library
- [ ] Storybook examples (optionnel)

📦 Délivrables :
- webapp/src/components/ui/*.tsx
- Tests (*.test.tsx)
```

**Durée** : 4-6h  
**Blocker ?** : Non  
**Dépend de** : TASK-007 (tokens)

---

### Étape 1.5 : TASK-003 — RBAC (High Priority)

**Agent** : @bmad-dev (Amelia)

**Prompt** :
```
🔧 TASK-003 : RBAC (Role-Based Access Control)

📋 Contexte :
- Epic : Epic 8
- Backlog : TASK-003

🎯 Objectif :
Implémenter roles admin/user :
1. Users table : role (admin/user)
2. Middleware RBAC (vérifie role)
3. Routes admin-only (/admin/*)

✅ Acceptance Criteria :
- [ ] Middleware authorize(role)
- [ ] Admin peut CRUD users
- [ ] User ne peut modifier que ses propres settings
- [ ] Tests (admin access OK, user access forbidden)

📦 Délivrables :
- internal/app/rbac.go
- Middleware + tests
```

**Durée** : 1 jour  
**Dépend de** : TASK-002 (auth)

---

## 📊 Suivi Sprint 1

### Daily Standup (Async)

**Agent** : @pm (John)

**Prompt** :
```
📋 Daily Standup : [DATE]

🎯 Objectif :
Status report rapide :
1. Yesterday : Quelles tasks complétées ?
2. Today : Quelles tasks en cours ?
3. Blockers : Problèmes identifiés ?

📦 Format court (< 200 mots)
```

**Fréquence** : Chaque jour de sprint  
**Durée** : 5 min

---

### Sprint Review (Fin Sprint 1)

**Agent** : @pm (John)

**Prompt** :
```
📋 Sprint Review : Sprint 1

📚 Contexte :
- Sprint goals : Sprint-1-Kickoff.md
- Backlog : _bmad-output/planning-artifacts/04-BACKLOG-SPRINTS.md

🎯 Objectif :
Analyser Sprint 1 :
1. Completed : Quelles tasks done (story points)
2. Incomplete : Quelles tasks reportées (pourquoi ?)
3. Velocity : Points complétés / points prévus
4. Lessons learned : Quoi améliorer ?

📦 Délivrables :
- Sprint-1-Review.md
- Recommandations Sprint 2
```

**Durée** : 1-2h  
**Blocker ?** : Non

---

## 🏃 Sprint 2 : Frontend Restructuration (2 semaines, 103 points)

### Étape 2.1 : Sprint Planning

**Agent** : @pm (John)

(Même pattern que Sprint 1)

---

### Étape 2.2 : TASK-013 — React Router Setup

**Agent** : @quick-flow (Barry)

**Prompt** :
```
⚡ TASK-013 : React Router Setup

📋 Contexte :
- Epic : Epic 2 (Search & Discovery)
- Backlog : TASK-013

🎯 Objectif :
Configurer routing :
1. Routes : /, /search, /jobs, /settings
2. Protected routes (auth required)
3. 404 page

✅ Acceptance Criteria :
- [ ] React Router v6
- [ ] Routes définies
- [ ] ProtectedRoute HOC (redirects si pas auth)
- [ ] Navigation menu

📦 Délivrables :
- webapp/src/routes.tsx
- ProtectedRoute component
```

**Durée** : 3-4h  
**Dépend de** : TASK-002 (auth backend)

---

### Étape 2.3 : TASK-014 — Zustand State Management

**Agent** : @bmad-dev (Amelia)

**Prompt** :
```
🔧 TASK-014 : Zustand State Management

📋 Contexte :
- Epic : Epic 1 (Foundation)
- Backlog : TASK-014

🎯 Objectif :
Remplacer useState hooks par Zustand :
1. authStore (user, token, login, logout)
2. jobsStore (jobs, addJob, updateJob)
3. settingsStore (settings, updateSettings)

✅ Acceptance Criteria :
- [ ] Stores dans webapp/src/stores/
- [ ] TypeScript interfaces
- [ ] Tests (store mutations)
- [ ] Persistence (localStorage si nécessaire)

📦 Délivrables :
- webapp/src/stores/*.ts
- Tests
```

**Durée** : 1 jour  
**Blocker ?** : Non

---

## 🎨 Workflows Spéciaux

### Workflow Design → Dev

**Quand** : Nouvelle feature UI complexe

**Étapes** :
1. @ux-designer : Créer mockups haute-fidélité
2. @architect : Review composants (architecture)
3. @quick-flow : Implémenter UI rapide
4. @quinn : Tests E2E user journey

**Exemple** : Dashboard Page (Sprint 2)

```
# 1️⃣ Design
@ux-designer Mockup Dashboard Page

# 2️⃣ Review
@architect Review Dashboard Components

# 3️⃣ Implémentation
@quick-flow TASK-016: Dashboard UI

# 4️⃣ Tests
@quinn E2E Tests : Dashboard User Journey
```

---

### Workflow Bug Fix

**Quand** : Bug identifié en prod/dev

**Étapes** :
1. @pm : Triage (priorité, sprint)
2. @bmad-dev : Fix + tests régression
3. @quinn : Validation E2E

**Exemple** :

```
# 1️⃣ Triage
@pm Bug Triage : [DESCRIPTION]

# 2️⃣ Fix
@bmad-dev Fix Bug : [DESCRIPTION]

# 3️⃣ Validation
@quinn Regression Tests : [BUG]
```

---

## 📞 Escalation & Déblocage

### Blocker Technique

**Si** : Agent bloqué > 4h

**Action** :
```
@pm Blocker Escalation : [TASK]

Contexte : [Détails blocker]
Action : Escalate à @architect (review technique)
Mitigation : Alternative approach ou reassign
```

---

### Scope Creep

**Si** : Task dépasse estimation > 50%

**Action** :
```
@pm Scope Review : [TASK]

Contexte : Task [X] estimée [Y] points, déjà [Z] points
Action : Split task ou ajuster acceptance criteria
```

---

## 🎯 Checklist Exécution Complète

### Sprint 1
- [ ] @pm : Sprint 1 Planning
- [ ] @architect : Architecture review
- [ ] @bmad-dev : TASK-001 (backend refactor)
- [ ] @bmad-dev : TASK-002 (auth JWT)
- [ ] @bmad-dev : TASK-003 (RBAC)
- [ ] @quick-flow : TASK-007 (design tokens)
- [ ] @quick-flow : TASK-008 (UI components)
- [ ] @ux-designer : Dashboard mockups
- [ ] @pm : Daily standups (async)
- [ ] @quinn : E2E tests auth flow
- [ ] @pm : Sprint 1 Review

### Sprint 2
- [ ] @pm : Sprint 2 Planning
- [ ] @quick-flow : TASK-013 (routing)
- [ ] @bmad-dev : TASK-014 (Zustand)
- [ ] @quick-flow : TASK-016 (Dashboard UI)
- [ ] @bmad-dev : TASK-017 (Search API)
- [ ] @quinn : E2E tests search flow
- [ ] @pm : Sprint 2 Review

### Sprint 3
- [ ] @pm : Sprint 3 Planning
- [ ] (... Suite des tasks selon backlog)

---

## 💡 Best Practices

### ✅ DO

- **Paralléliser** : Backend (Amelia) + Frontend (Barry) en même temps
- **Review critique** : Winston review TASK-001, TASK-002 (architecture + sécurité)
- **Tests first** : Chaque task doit avoir tests AVANT merge
- **Daily standup** : John track tous les jours (async Slack/doc)

### ❌ DON'T

- **Pas de skip tests** : Jamais merge sans tests (70% coverage minimum)
- **Pas de blocker silencieux** : Escalate immédiatement si bloqué > 4h
- **Pas de scope creep** : Respecter AC strictement, split si nécessaire

---

## 🚀 TL;DR — Démarrage Rapide

```bash
# Aujourd'hui (31 janvier 2026)
1️⃣ @pm Sprint Planning : Sprint 1
2️⃣ @architect Architecture Review : Projet Global
3️⃣ @bmad-dev TASK-001 : Backend Refactoring
4️⃣ @quick-flow TASK-007 : Design System Tokens

# Demain
5️⃣ @bmad-dev TASK-002 : Auth JWT
6️⃣ @quick-flow TASK-008 : UI Components Library

# Dans 3 jours
7️⃣ @bmad-dev TASK-003 : RBAC
8️⃣ @ux-designer Mockup Dashboard Page
```

**Durée Sprint 1** : 2 semaines (jusqu'au 14 février 2026)  
**Velocity cible** : 44.5 pts/semaine

---

**Bon sprint ! 🚀**
