---
stepsCompleted: [1, 2, 3]
inputDocuments:
  - _bmad-output/planning-artifacts/prd.md
  - _bmad-output/planning-artifacts/00-PROJECT-BRIEF.md
  - _bmad-output/planning-artifacts/01-PERSONAS-AND-JOURNEYS.md
  - _bmad-output/planning-artifacts/02-DESIGN-SYSTEM-SAKURA-NIGHT.md
  - _bmad-output/planning-artifacts/03-TECHNICAL-ARCHITECTURE.md
  - _bmad-output/planning-artifacts/04-BACKLOG-SPRINTS.md
workflowType: 'ux-design'
documentCounts:
  briefCount: 1
  personasCount: 1
  designSystemCount: 1
  architectureCount: 1
  backlogCount: 1
  prdCount: 1
lastUpdated: '2026-01-31'
---

# UX Design Specification — Anime-Sama Downloader v1.0

**Author:** Guilhem
**Date:** 31 janvier 2026
**Design Era:** Sakura Night™

---

## Executive Summary

### Project Vision

**Anime-Sama Downloader v1.0** est une refonte complète d'une application de téléchargement d'animes. L'objectif est de transformer une interface "désastreuse" en une application **moderne, intuitive et visuellement inspirée par l'esthétique anime/manga**.

**Problème à résoudre** :
- Interface confuse et austère (années 2000)
- UX complex pour télécharger un épisode
- Aucune identité visuelle anime
- Pas d'automatisation (subscriptions, sync AniList)

**Vision cible** :
- Télécharger un anime en **3 clics** pour les casual users
- Automatisation complète pour power users
- Identité visuelle forte : **"Sakura Night"** (minimalisme moderne + esthétique anime)

### Target Users

**1. Alex le Casual Fan 🍿** (22-28 ans, étudiant/jeune dev)
- Veut télécharger les derniers épisodes **rapidement**
- Compétences basiques (Docker, pas plus)
- Frustration clé : "Pourquoi c'est si confus et moche ?"
- **Besoin UX** : 3 clics, feedback clair, interface engageante

**2. Maya la Jellyfin Power User 📚** (28-40 ans, DevOps/SysAdmin)
- Veut **automatiser** (AniList sync → Jellyfin)
- Compétences avancées (Linux, APIs, webhooks)
- Frustration clé : "Pourquoi pas d'automatisation complète ?"
- **Besoin UX** : Mode expert, APIs, webhooks, monitoring

### Key Design Challenges

1. **Simplifier sans sacrifier la puissance**
   - Alex a besoin de "3 clics"
   - Maya a besoin de "mode expert complet"
   - Solution : Dashboard pour Alex + Mode avancé pour Maya

2. **Clarifier les choix d'anime complexes**
   - Plusieurs saisons, VF/VOSTFR, groupes de sous-titres
   - Frustration actuelle : "Lequel choisir ?"
   - Solution : Preview metadata, filtres explicites, labels clairs

3. **Feedback temps-réel sur les downloads**
   - Actuellement : confusion sur statut job
   - Solution : Progress bars animées + SSE streaming + notifications

4. **Intégration Jellyfin sans douleur**
   - Maya veut metadata correcte + naming strategy
   - Solution : Webhooks + auto-naming conforme Jellyfin

### Design Opportunities

1. **Dashboard "Nouveautés" personnalisé**
   - Card visuelle "Jujutsu Kaisen S2 Ep12 disponible"
   - Engagement vs. interface austère actuelle

2. **One-click subscriptions**
   - AniList sync → subscriptions auto
   - Scheduler configurable
   - Game changer pour Maya
   - **Note** : Phase 2 (Growth), pas MVP

3. **Micro-interactions animées**
   - Progress bars fluides + pétales sakura subtils
   - Rendre le téléchargement *satisfying*
   - Créer une expérience mémorable

---

## Core User Experience

### Defining Experience

**Phase 1 MVP Strategy** : Deux expériences adaptées dans une seule interface, sans calendrier.

**Action Core pour Alex 🍿** : Rechercher un anime → Télécharger → Suivre progression
- Doit être **< 30 secondes** et **≤ 3 clics**
- Interface épurée, minimaliste, sans bruit
- Feedback immédiat sur chaque action

**Action Core pour Maya 📚** : Configurer webhooks Jellyfin → Automatiser downloads → Monitorer queue
- Doit être **complètement robuste** avec logs accessibles
- Dashboard expert avec stats en temps réel
- API REST pour intégrations custom

### Platform Strategy

- **Platform primaire** : Web SPA (React + TypeScript)
- **Responsiveness** : Desktop prioritaire, mobile acceptable (responsive)
- **Real-time requirements** : SSE pour streaming job progress, webhooks pour Jellyfin
- **Offline** : Non nécessaire (dépend Internet pour scraping anime-sama.si)
- **CLI** : Phase 2+ (commandes allégées via API)

### Effortless Interactions

**Pour Alex (Mode Simple)** :
- ✨ Tapez titre → suggestions autocomplete instantanées (< 100ms)
- ✨ Cliquez anime → métadonnées explicites (saison, épisode, langue)
- ✨ Cliquez "Télécharger" → confirmation + ajout à queue visible
- ✨ Progress bar fluide avec ETA + notification quand prêt

**Pour Maya (Mode Expert)** :
- ✨ Toggle "Power User Mode" → Dashboard enrichi (onglets Stats, Logs, Webhooks)
- ✨ Configure Jellyfin webhook une fois → fonctionne automatiquement après
- ✨ Logs structurés filtrables par anime/status/date
- ✨ Webhooks testables avec "Send Test Event"

### Critical Success Moments

1. **Premier téléchargement (Alex)** 
   - Moment : "Je clique télécharger et je vois du feedback IMMÉDIATEMENT"
   - Risque d'échec : silence radio → utilise autre app
   - Solution : toast "Ajouté à la queue" + progress bar immédiatement visible

2. **Clarification choix d'anime complexes**
   - Moment : "Je sélectionne la bonne saison/langue sans ambiguïté"
   - Risque d'échec : mauvais choix → download inutile → frustration
   - Solution : labels explicites, preview metadata, confirmation avant télécharger

3. **Error Transparency**
   - Moment : "Quand quelque chose échoue, je comprends pourquoi"
   - Risque d'échec : erreur silencieuse → utilisateur croit que ça marche
   - Solution : toast d'erreur clair + logs détaillés pour Maya, simplifié pour Alex

4. **Jellyfin Integration Setup (Maya)**
   - Moment : "Je configure webhooks, c'est prêt en 5 minutes"
   - Risque d'échec : wizard confus → abandonne = script custom
   - Solution : wizard pas-à-pas avec test button

5. **Mode Toggle Seamless**
   - Moment : "Je switch en Mode Expert, tout reste accessible"
   - Risque d'échec : toggle casse les states → page blank
   - Solution : persistence localStorage, state consistency tests

### Experience Principles

**Basé sur 3 méthodes avancées (Focus Group, Debate, War Room)** :

1. **Clarté avant pouvoir** ✅
   - Choix d'anime complexes = labels explicites + preview
   - Jamais d'ambiguïté (VOSTFR vs VF vs saison)
   - Mode Simple reste simple, Mode Expert reste accessible

2. **Feedback immédiat + Error Transparency** 🔴 **CRITIQUE**
   - Chaque action utilisateur = réponse visuelle dans 200ms
   - **Les erreurs ne sont JAMAIS silencieuses**
   - Alex : messages simples ("Download failed, retrying in 5min")
   - Maya : logs full stack pour debugger

3. **Power User Mode Toggle (Architecture Pattern)**
   - Menu Settings → toggle "Power User Mode" (default OFF)
   - Quand ON : onglets Webhooks/Logs/Stats apparaissent
   - Alex ne clique jamais dessus = voit jamais la complexité
   - Maya active une fois = reste en mode expert
   - Persistence localStorage + state consistency

4. **Deux expériences, une codebase (Zustand + Composition)**
   - Pas de duplication massive (code partagé 40%, variants 35%, distinct 25%)
   - Architecture : Zustand stores + composants simple/expert séparés
   - Maintenance claire, testing parallèle possible
   - Timeline : ~27 jours (5-6 sprints)

5. **Esthétique & fonction (Sakura Night)**
   - Palette magenta/cyan + dark mode = engagement visuel
   - Priorité pour Alex (casual engagement)
   - Maya accepte du moment que c'est usable

6. **Jellyfin Integration sans douleur**
   - Wizard onboarding pour webhooks (pas de doc manual)
   - Naming strategy documentée + auto-validée
   - Test button pour vérifier connection
   - Support multiples versions Jellyfin (10.8, 10.9, 10.10)

### Scoping Note

**Phase 1 MVP** (sans) :
- ❌ Calendrier airing schedule
- ❌ Subscriptions auto + scheduler
- ❌ AniList sync bidirectionnel
- ❌ Mobile app native
- ❌ Multi-users RBAC avancé

**Phase 1 MVP** (avec) :
- ✅ Mode Simple (Alex) : recherche + download + queue + notifications
- ✅ Mode Expert (Maya) : webhooks Jellyfin + logs + stats + settings avancés
- ✅ Error handling transparent
- ✅ Sakura Night design system
- ✅ SSE streaming pour job progress
- ✅ Jellyfin wizard setup

---

