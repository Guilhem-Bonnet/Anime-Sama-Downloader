# 📋 PROJECT BRIEF — Anime-Sama Downloader v1.0 Refonte

**Date** : 31 janvier 2026  
**Project Owner** : Guilhem Bonnet  
**Durée estimée** : 10-13 semaines (2,5-3 mois)  
**Budget** : Aucune deadline stricte — qualité prioritaire

---

## 🎯 OBJECTIF PRINCIPAL

Refondre complètement **Anime-Sama Downloader** pour en faire une application moderne, intuitive et visuellement inspirée par l'esthétique anime/manga, tout en nettoyant les résidus techniques et en améliorant l'expérience utilisateur de A à Z.

---

## 📊 CONTEXTE

### État actuel
- **Backend** : Go 1.22+ avec architecture propre (Clean Architecture)
- **Frontend** : React SPA avec UI "désastreuse" (monolithe 1922 lignes, UX confuse)
- **Problèmes** : Layouts cassés, navigation complexe, pas d'identité visuelle, bugs dans la sélection d'épisodes

### Vision cible
Une application **élégante, rapide et intuitive** qui :
- Permet de télécharger des animes depuis anime-sama.si en **3 clics max** (casual users)
- Offre une gestion avancée pour les power users (abonnements, scheduler, intégration Jellyfin)
- Dégage une **identité visuelle forte** inspirée par le style anime moderne

---

## 👥 PERSONAS CIBLES

### Alex le Casual Fan 🍿
- 22-28 ans, étudiant/jeune actif
- Veut télécharger les dernières sorties VOSTFR simplement
- Compétences : Basiques (peut lancer Docker)
- Besoin : Simplicité, rapidité, notifications

### Maya la Jellyfin Power User 📚
- 28-40 ans, self-hoster avec NAS
- Veut automatiser sa médiathèque anime
- Compétences : Techniques (Linux, APIs, Docker)
- Besoin : Abonnements auto, sync AniList, intégration Jellyfin

---

## 🎨 DIRECTION ARTISTIQUE

**Thème** : **"Sakura Night"**  
Fusion entre minimalisme moderne et esthétique anime artisanale

**Palette** :
- Base : Noir profond + Gris bleuté sombre
- Accents : Magenta électrique (#D946EF) + Cyan néon (#06B6D4)
- Highlights : Or doux + Rose Sakura

**Éléments visuels** :
- Brush strokes subtils en arrière-plan
- Animation pétales de sakura (légère)
- Gradients magenta→cyan sur interactions
- Typographie : Inter + calligraphie japonaise pour titres

---

## 🏗️ ARCHITECTURE TECHNIQUE

### Backend (Go)
- Clean Architecture maintenue
- Nouveaux services : UserService, JellyfinWebhookService
- Rate limiting anime-sama
- Test coverage 70%+

### Frontend (React + TypeScript)
- Refactoring complet : Atomic Design
- State management centralisé (Zustand)
- React Router pour navigation
- Design system avec tokens CSS

### Intégration Jellyfin
- Service indépendant avec webhooks
- Événement `download.completed` → scan Jellyfin auto
- Naming strategy compatible Jellyfin

---

## 📅 ROADMAP

**Phase 1 : Fondations** (2-3 semaines)
- Cleanup backend + frontend
- Design system
- Tests 70%+

**Phase 2 : UI/UX Refonte** (3-4 semaines)
- Charte graphique "Sakura Night"
- Composants modulaires
- Nouvelles pages

**Phase 3 : Features Avancées** (3-4 semaines)
- Multi-users + auth
- Jellyfin integration
- AniList auto-sync

**Phase 4 : Polish & Release** (1-2 semaines)
- Performance, accessibility
- Documentation
- Release v1.0

---

## 📦 LIVRABLES ATTENDUS

### Documentation
- ✅ Personas & User Journeys
- ✅ Design System complet
- ✅ Architecture technique
- ✅ Backlog détaillé (GitHub Issues)
- ✅ Roadmap par sprint

### Code
- Backend Go refactorisé + tests
- Frontend React modulaire + design system
- Intégration Jellyfin fonctionnelle
- CI/CD (linting + tests auto)

### UI/UX
- Mockups haute-fidélité (4 pages principales)
- Composants UI réutilisables
- Animations & micro-interactions
- Responsive (desktop + mobile)

---

## 🎯 MÉTRIQUES DE SUCCÈS

- **UX** : Télécharger 1 anime en ≤ 3 clics (Alex)
- **Performance** : Lighthouse score ≥ 90
- **Tests** : Coverage ≥ 70%
- **Accessibilité** : WCAG AA compliance
- **Satisfaction** : Feedbacks positifs des 2 personas

---

## 🚨 RISQUES & CONTRAINTES

### Risques
- ⚠️ Anime-sama peut changer structure HTML → scraping break
- ⚠️ AniList API rate limiting
- ⚠️ Complexité Jellyfin integration (multiples versions)

### Contraintes
- 10 téléchargements simultanés max (suffisant)
- 1-3 utilisateurs concurrents
- Dépendance ffmpeg (doit être documentée/installée)

---

## 📞 CONTACTS & RESSOURCES

**Repository** : https://github.com/Guilhem-Bonnet/Anime-Sama-Downloader  
**Branch** : `go-rewrite`  
**Documentation** : Voir `_bmad-output/planning-artifacts/`

---

## 🎬 NEXT STEPS

1. **PM** : Lire ce brief + tous les docs dans `planning-artifacts/`
2. **PM** : Créer milestones GitHub + issues depuis backlog
3. **Architect** : Revoir architecture technique avec équipe
4. **Designer** : Finaliser mockups haute-fidélité
5. **Dev Lead** : Estimer efforts par tâche
6. **Team** : Sprint Planning #1

---

**Ce document est le point d'entrée.** Consultez les autres artéfacts pour les détails techniques, design, et exécution.
