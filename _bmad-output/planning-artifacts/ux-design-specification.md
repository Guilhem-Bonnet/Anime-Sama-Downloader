---
stepsCompleted: [1, 2]
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

3. **Micro-interactions animées**
   - Progress bars fluides + pétales sakura subtils
   - Rendre le téléchargement *satisfying*
   - Créer une expérience mémorable

---

