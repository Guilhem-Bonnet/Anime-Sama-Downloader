# 📋 BACKLOG DÉTAILLÉ — Sprints & User Stories

**Projet** : Anime-Sama Downloader v1.0  
**Date** : 31 janvier 2026  
**Auteur** : Bob (Scrum Master)

---

## 🏃 SPRINT 1 : FONDATIONS & CLEANUP (2 semaines)

**Objectif** : Nettoyer la codebase, augmenter la qualité, préparer le refactoring

### Backend Go

**TASK-001** : Audit dépendances + documentation ffmpeg  
**Story Points** : 3  
**Description** : Auditer toutes les dépendances Go, documenter requirement ffmpeg dans README, ajouter check au démarrage  
**Acceptance Criteria** :
- [ ] Liste complète des dépendances documentée
- [ ] README mentionne ffmpeg comme requirement
- [ ] Server log warning si ffmpeg manquant
- [ ] Script `scripts/check-deps.sh` créé

**TASK-002** : Augmenter test coverage à 70%+  
**Story Points** : 8  
**Description** : Ajouter tests manquants (unit + integration), particulièrement sur SubscriptionService, JobService, scrapers  
**Acceptance Criteria** :
- [ ] Coverage ≥ 70% (mesure avec `go test -cover`)
- [ ] Tests d'intégration pour tous les endpoints API
- [ ] Tests pour AnimeSamaScraper (mock HTTP)
- [ ] Tests pour AniListClient (mock GraphQL)

**TASK-003** : Rate limiting anime-sama  
**Story Points** : 5  
**Description** : Implémenter rate limiter (1 req/sec) pour éviter ban anime-sama  
**Acceptance Criteria** :
- [ ] Rate limiter configurable (env var `RATE_LIMIT_RPS`)
- [ ] Logs quand rate limit atteint
- [ ] Tests pour vérifier respect limite

**TASK-004** : Refactor SubscriptionService  
**Story Points** : 5  
**Description** : SubscriptionService fait trop de choses → extraire `SchedulerService` et `EpisodeResolver`  
**Acceptance Criteria** :
- [ ] SchedulerService gère `nextCheckAt` logic
- [ ] EpisodeResolver gère fetching episodes anime-sama
- [ ] Tests passent après refactor
- [ ] Aucune régression fonctionnelle

**TASK-005** : Compléter OpenAPI schema  
**Story Points** : 3  
**Description** : Documenter tous les endpoints dans openapi.json, ajouter examples  
**Acceptance Criteria** :
- [ ] Tous les endpoints documentés
- [ ] Request/response examples présents
- [ ] Schémas pour tous les DTOs
- [ ] Validation avec Swagger Editor

**TASK-006** : Logging structuré avec contexte  
**Story Points** : 5  
**Description** : Ajouter request_id à tous les logs, logger headers pertinents  
**Acceptance Criteria** :
- [ ] Middleware génère `request_id` (UUID)
- [ ] Tous les logs incluent `request_id`
- [ ] Logs incluent `user_id` si authentifié
- [ ] Format JSON pour production

**TASK-007** : Graceful shutdown  
**Story Points** : 3  
**Description** : Implémenter shutdown propre (drain workers, close DB, etc.)  
**Acceptance Criteria** :
- [ ] Signal SIGTERM/SIGINT capturé
- [ ] Workers drainés (max 30s timeout)
- [ ] Database connections closed proprement
- [ ] HTTP server shutdown graceful

### Frontend React

**TASK-101** : Découper App.tsx en composants modulaires  
**Story Points** : 13  
**Description** : App.tsx = 1922 lignes → découper en pages + composants réutilisables  
**Acceptance Criteria** :
- [ ] Pages créées (Dashboard, Search, Subscriptions, Calendar, Jobs, Settings)
- [ ] Composants UI extraits (Header, TabNav, Card, Button, Input)
- [ ] App.tsx < 200 lignes (juste routing)
- [ ] Aucune régression fonctionnelle

**TASK-102** : Créer design system (tokens + composants base)  
**Story Points** : 8  
**Description** : Implémenter tokens CSS (couleurs, spacing, etc.) + composants UI primitifs  
**Acceptance Criteria** :
- [ ] Fichier `tokens.css` avec variables CSS
- [ ] Composants : Button, Input, Card, Modal, Badge, ProgressBar
- [ ] Storybook ou page de démo pour chaque composant
- [ ] Documentation usage

**TASK-103** : Migrer state vers Zustand  
**Story Points** : 8  
**Description** : Remplacer les 20+ `useState` par stores Zustand centralisés  
**Acceptance Criteria** :
- [ ] Stores : subscriptionsStore, jobsStore, settingsStore
- [ ] Tous les composants utilisent stores au lieu de useState
- [ ] Performance améliorée (moins de re-renders)
- [ ] Tests pour stores

**TASK-104** : Refonte EpisodeSelector component  
**Story Points** : 5  
**Description** : UX intuitive pour sélectionner épisodes (remplacer range + input texte)  
**Acceptance Criteria** :
- [ ] Checkbox list avec preview
- [ ] Raccourcis : "Tous", "Derniers 5", "Custom range"
- [ ] Validation inline (rouge si invalide)
- [ ] Accessible (keyboard navigation)

**TASK-105** : Ajouter React Router  
**Story Points** : 5  
**Description** : Remplacer tabs inline par vraies routes  
**Acceptance Criteria** :
- [ ] Routes : /, /search, /subscriptions, /calendar, /jobs, /settings
- [ ] Navigation avec React Router Link
- [ ] URL reflète la page active
- [ ] Back button fonctionne

**TASK-106** : Error boundaries + loading states  
**Story Points** : 5  
**Description** : Gérer erreurs proprement + loading states cohérents  
**Acceptance Criteria** :
- [ ] Error boundary au niveau app
- [ ] Skeleton loaders pour listes
- [ ] Toast notifications pour erreurs
- [ ] Retry button si erreur API

**TASK-107** : Optimiser re-renders  
**Story Points** : 5  
**Description** : React.memo, useMemo sur listes longues, éviter re-renders inutiles  
**Acceptance Criteria** :
- [ ] React.memo sur composants lourds
- [ ] useMemo sur listes triées/filtrées
- [ ] React DevTools Profiler confirme amélioration
- [ ] Lighthouse performance score ≥ 85

### Infrastructure

**TASK-201** : Optimiser Dockerfiles  
**Story Points** : 3  
**Description** : Multi-stage build, reduce image size, cache layers  
**Acceptance Criteria** :
- [ ] Image finale < 50MB
- [ ] Build time < 2 min
- [ ] Cache Docker layers efficace
- [ ] README instructions claires

**TASK-202** : Health checks endpoints  
**Story Points** : 2  
**Description** : `/health` et `/ready` endpoints pour monitoring  
**Acceptance Criteria** :
- [ ] GET /health → 200 si server UP
- [ ] GET /ready → 200 si DB accessible + workers ready
- [ ] Docker Compose health check configuré

**TASK-203** : Documenter variables env  
**Story Points** : 2  
**Description** : `.env.example` complet avec descriptions  
**Acceptance Criteria** :
- [ ] Toutes les vars env documentées
- [ ] Valeurs par défaut indiquées
- [ ] Exemples pour chaque type de déploiement

**TASK-204** : CI/CD linting + tests  
**Story Points** : 5  
**Description** : GitHub Actions pour lint + tests auto sur PR  
**Acceptance Criteria** :
- [ ] Workflow `.github/workflows/ci.yml`
- [ ] Lint Go (golangci-lint)
- [ ] Lint TypeScript (ESLint)
- [ ] Tests Go + React
- [ ] Fail si coverage < 70%

---

## 🎨 SPRINT 2 : REDESIGN UI/UX (3 semaines)

**Objectif** : Implémenter charte "Sakura Night", refaire toutes les pages

### Design System

**TASK-301** : Implémenter thème "Sakura Night"  
**Story Points** : 13  
**Description** : Appliquer palette, typographie, composants selon design system  
**Acceptance Criteria** :
- [ ] Tous les tokens CSS implémentés
- [ ] Gradients + brush strokes background
- [ ] Animations sakura petals (Canvas)
- [ ] Dark mode par défaut
- [ ] Lighthouse score ≥ 90

**TASK-302** : Créer composants anime-specific  
**Story Points** : 8  
**Description** : AnimeCard, EpisodeList, SubscriptionCard avec nouveau design  
**Acceptance Criteria** :
- [ ] AnimeCard avec cover image + metadata
- [ ] EpisodeList avec checkboxes + preview
- [ ] SubscriptionCard avec status badges
- [ ] Hover effects + micro-interactions

### Pages

**TASK-303** : Dashboard refonte  
**Story Points** : 13  
**Description** : Page d'accueil avec sections "Nouveautés", "Abonnements", "Téléchargements"  
**Acceptance Criteria** :
- [ ] Section "Nouveautés de la semaine" (cards visuelles)
- [ ] Section "Mes abonnements" (liste avec status)
- [ ] Section "Téléchargements en cours" (progress bars)
- [ ] UX : Télécharger 1 anime en ≤ 3 clics
- [ ] Responsive (mobile + desktop)

**TASK-304** : Page Recherche refonte  
**Story Points** : 8  
**Description** : Recherche avec autocomplete + filtres  
**Acceptance Criteria** :
- [ ] Input avec autocomplete (debounced)
- [ ] Filtres : Saison, Langue, Status
- [ ] Cards résultats avec cover + metadata
- [ ] Boutons "Télécharger" et "S'abonner"

**TASK-305** : Page Abonnements refonte  
**Story Points** : 8  
**Description** : Liste abonnements avec tri + filtres  
**Acceptance Criteria** :
- [ ] Tri : Nom, Prochaine vérif, Derniers épisodes
- [ ] Filtres : Status (à jour, nouveaux épisodes)
- [ ] Actions : Sync, Enqueue, Delete
- [ ] Bulk actions (select multiple)

**TASK-306** : Page Calendar refonte  
**Story Points** : 8  
**Description** : Calendrier visuel des sorties hebdomadaires  
**Acceptance Criteria** :
- [ ] Timeline par jour (colonnes)
- [ ] Cards anime avec heure de sortie
- [ ] Filtres : Abonnements seulement, Tous
- [ ] Sync AniList airing schedule

**TASK-307** : Page Jobs refonte  
**Story Points** : 8  
**Description** : Queue jobs avec progress bars + logs  
**Acceptance Criteria** :
- [ ] Onglets : En cours, Queue, Terminés, Échoués
- [ ] Progress bars animées + ETA
- [ ] Logs streamés (SSE)
- [ ] Actions : Cancel, Retry

**TASK-308** : Page Settings refonte  
**Story Points** : 5  
**Description** : Formulaire settings avec validation inline  
**Acceptance Criteria** :
- [ ] Sections : Général, AniList, Jellyfin, Avancé
- [ ] Validation inline (rouge si invalide)
- [ ] Tooltips pour chaque setting
- [ ] Test connection buttons

### Animations

**TASK-309** : Animations & micro-interactions  
**Story Points** : 8  
**Description** : Framer Motion pour transitions + hover effects  
**Acceptance Criteria** :
- [ ] Page transitions (fade in)
- [ ] Card hover (lift + glow)
- [ ] Button press (scale)
- [ ] Loading states (skeleton + shimmer)
- [ ] Respect `prefers-reduced-motion`

**TASK-310** : Responsive design  
**Story Points** : 8  
**Description** : Mobile-first, adaptatif tablette + desktop  
**Acceptance Criteria** :
- [ ] Breakpoints : 640px, 768px, 1024px
- [ ] Grids adaptatives
- [ ] Navigation mobile (hamburger menu)
- [ ] Touch-friendly (buttons ≥ 44px)

---

## 🚀 SPRINT 3 : FEATURES AVANCÉES (3 semaines)

**Objectif** : Multi-users, Jellyfin, AniList auto-sync

### Backend

**TASK-401** : Multi-users + authentication  
**Story Points** : 13  
**Description** : JWT auth, rôles (admin/user), isolation subscriptions  
**Acceptance Criteria** :
- [ ] Table `users` créée
- [ ] Endpoints : /auth/login, /auth/logout, /auth/me
- [ ] JWT tokens (access + refresh)
- [ ] Middleware auth sur toutes les routes
- [ ] RBAC : admin full access, user own data only

**TASK-402** : UserService + tests  
**Story Points** : 8  
**Description** : Service CRUD users, password hashing (bcrypt)  
**Acceptance Criteria** :
- [ ] CRUD complet
- [ ] Password hashing bcrypt
- [ ] Tests unitaires + intégration
- [ ] Validation (email, password strength)

**TASK-403** : AniList auto-import amélioré  
**Story Points** : 13  
**Description** : Sync watchlist → subscriptions auto, scheduler hebdo  
**Acceptance Criteria** :
- [ ] Endpoint POST /anilist/import-auto
- [ ] Match anime → anime-sama (fuzzy search)
- [ ] User valide matches → create subs
- [ ] Scheduler hebdo (cron) pour re-sync
- [ ] Tests avec mocks AniList API

**TASK-404** : Scheduler calendrier  
**Story Points** : 8  
**Description** : Cron-like scheduler pour checks périodiques  
**Acceptance Criteria** :
- [ ] Config : interval (daily, weekly, custom cron)
- [ ] Scheduler tourne en background
- [ ] Logs chaque run
- [ ] Graceful shutdown du scheduler

**TASK-405** : JellyfinWebhookService  
**Story Points** : 13  
**Description** : Webhook POST après download, trigger scan Jellyfin  
**Acceptance Criteria** :
- [ ] Config : jellyfin_url, api_key, library_id
- [ ] Event `download.completed` → POST webhook
- [ ] Jellyfin library scan triggered
- [ ] Retry logic si Jellyfin down
- [ ] Tests avec mock Jellyfin API

**TASK-406** : Naming strategy Jellyfin  
**Story Points** : 5  
**Description** : Renommer fichiers selon format Jellyfin-compatible  
**Acceptance Criteria** :
- [ ] Format : `[Title] - S[Season]E[Episode].ext`
- [ ] Metadata AniList/AniDB embedded (optionnel)
- [ ] Config : custom naming template
- [ ] Tests pour vérifier noms corrects

**TASK-407** : Notification system  
**Story Points** : 5  
**Description** : Push notifications via SSE (download completed, new episode)  
**Acceptance Criteria** :
- [ ] SSE events : `download.completed`, `episode.available`
- [ ] Frontend toast notifications
- [ ] Optionnel : Email/Discord webhooks
- [ ] Tests SSE

### Frontend

**TASK-408** : Page Login/Register  
**Story Points** : 8  
**Description** : Formulaires login + register + forgot password  
**Acceptance Criteria** :
- [ ] Formulaire login (email + password)
- [ ] Formulaire register (username, email, password)
- [ ] Validation inline
- [ ] Redirection après login
- [ ] Error handling (mauvais credentials)

**TASK-409** : Auth context + protected routes  
**Story Points** : 5  
**Description** : Context React pour auth, redirect si non authentifié  
**Acceptance Criteria** :
- [ ] AuthContext avec login/logout/refresh
- [ ] Protected routes (redirect vers /login)
- [ ] Persist auth dans localStorage
- [ ] Auto-refresh token

**TASK-410** : Page Profil utilisateur  
**Story Points** : 5  
**Description** : Éditer profil, changer password, déconnecter  
**Acceptance Criteria** :
- [ ] Form edit username/email
- [ ] Form change password
- [ ] Button logout
- [ ] Button connect AniList (OAuth)

**TASK-411** : AniList import UI  
**Story Points** : 8  
**Description** : Modal pour importer watchlist AniList  
**Acceptance Criteria** :
- [ ] Button "Importer depuis AniList"
- [ ] Modal avec preview matches
- [ ] User peut valider/rejeter chaque match
- [ ] Confirmation + création subs
- [ ] Loading states + errors

**TASK-412** : Jellyfin settings UI  
**Story Points** : 3  
**Description** : Section settings pour config Jellyfin  
**Acceptance Criteria** :
- [ ] Inputs : URL, API Key, Library ID
- [ ] Button "Test connection"
- [ ] Validation + feedback (success/error)

---

## 🧪 SPRINT 4 : POLISH & RELEASE (2 semaines)

**Objectif** : Perf, accessibilité, doc, release v1.0

**TASK-501** : Performance audit  
**Story Points** : 8  
**Description** : Lighthouse score ≥ 90, optimisations  
**Acceptance Criteria** :
- [ ] Lighthouse Performance ≥ 90
- [ ] Bundle size < 300KB (gzipped)
- [ ] Code splitting (React.lazy)
- [ ] Image optimization (WebP + lazy load)

**TASK-502** : Accessibility audit  
**Story Points** : 8  
**Description** : WCAG AA compliance, screen readers  
**Acceptance Criteria** :
- [ ] Contrastes respectés (AAA pour texte)
- [ ] Focus visible sur tous les interactifs
- [ ] Aria labels sur icônes
- [ ] Keyboard navigation complète
- [ ] Tests avec screen reader

**TASK-503** : i18n support  
**Story Points** : 8  
**Description** : Internationalisation (FR + EN)  
**Acceptance Criteria** :
- [ ] react-i18next configuré
- [ ] Toutes les strings traduites
- [ ] Détection langue navigateur
- [ ] Sélecteur langue dans settings

**TASK-504** : Documentation utilisateur  
**Story Points** : 5  
**Description** : Quick Start, FAQ, Troubleshooting  
**Acceptance Criteria** :
- [ ] README complet (installation, usage)
- [ ] QUICK_START.md (< 5 min setup)
- [ ] FAQ.md (questions fréquentes)
- [ ] TROUBLESHOOTING.md (debug common issues)

**TASK-505** : Documentation développeur  
**Story Points** : 5  
**Description** : Architecture, Contributing, API docs  
**Acceptance Criteria** :
- [ ] ARCHITECTURE.md (diagrammes)
- [ ] CONTRIBUTING.md (guidelines)
- [ ] API.md (référence complète)
- [ ] CHANGELOG.md (depuis v0.x)

**TASK-506** : E2E tests  
**Story Points** : 13  
**Description** : Tests Playwright pour user flows critiques  
**Acceptance Criteria** :
- [ ] Test : Télécharger 1 anime (search → download)
- [ ] Test : Créer abonnement
- [ ] Test : Importer depuis AniList
- [ ] Test : Gérer jobs (cancel, retry)
- [ ] CI : Run E2E tests sur PR

**TASK-507** : Security audit  
**Story Points** : 8  
**Description** : Audit sécurité (input validation, XSS, SQL injection)  
**Acceptance Criteria** :
- [ ] Audit input validation (tous les endpoints)
- [ ] Protection XSS (sanitize HTML)
- [ ] Protection SQL injection (parameterized queries)
- [ ] HTTPS enforcement (production)
- [ ] Security headers (CSP, HSTS)

**TASK-508** : Release v1.0  
**Story Points** : 5  
**Description** : Préparer release, migration guide, tag Git  
**Acceptance Criteria** :
- [ ] CHANGELOG.md v1.0
- [ ] MIGRATION.md (depuis v0.x)
- [ ] Git tag v1.0.0
- [ ] Docker images published (Docker Hub)
- [ ] GitHub Release avec binaries

---

## 📊 RÉCAPITULATIF

| Sprint | Objectif | Story Points | Durée |
|--------|----------|--------------|-------|
| Sprint 1 | Fondations & Cleanup | 89 | 2 sem |
| Sprint 2 | Redesign UI/UX | 103 | 3 sem |
| Sprint 3 | Features Avancées | 110 | 3 sem |
| Sprint 4 | Polish & Release | 60 | 2 sem |
| **TOTAL** | | **362 SP** | **10 sem** |

**Vélocité estimée** : 35-40 SP/semaine (équipe 2-3 devs)

---

## 🎯 DÉFINITION OF DONE

Une story est **Done** quand :
- [ ] Code écrit + tests passent
- [ ] Code review approuvé
- [ ] Tests unitaires + intégration (si applicable)
- [ ] Documentation mise à jour
- [ ] Déployé sur env staging + validé
- [ ] Aucune régression détectée

---

**🎬 Prochaine étape** : Créer GitHub Issues depuis ce backlog, assigner milestones par sprint.
