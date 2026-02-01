# 🎨 DESIGN SYSTEM — "Sakura Night"

**Projet** : Anime-Sama Downloader v1.0  
**Date** : 31 janvier 2026  
**Auteurs** : Sally (UX Designer) + Carson (Brainstorming Coach)

---

## 🌙 CONCEPT

**"Sakura Night"** fusionne le **minimalisme moderne** avec l'**esthétique anime artisanale**. Imaginez un cinéma nocturne japonais — élégant, feutré, avec des touches vibrantes qui évoquent les affiches d'anime et les traits de pinceau traditionnels.

**Inspirations** :
- Anime modernes : *Demon Slayer*, *Jujutsu Kaisen*, *Spy x Family*
- Design japonais : Ukiyo-e, seigaiha (vagues), sakura
- UI modernes : Vercel, Linear, Arc Browser (dark mode)

---

## 🎨 PALETTE DE COULEURS

### Couleurs de base

```css
/* Backgrounds */
--sakura-bg-base: #0A0E1A;        /* Noir profond (ciel nocturne) */
--sakura-bg-surface: #1A1F2E;     /* Gris bleuté sombre (surface) */
--sakura-bg-elevated: #252A3B;    /* Surface surélevée */

/* Texte */
--sakura-text-primary: #F5F7FF;   /* Blanc légèrement teinté bleu */
--sakura-text-secondary: #A8B3D1; /* Gris bleuté clair */
--sakura-text-muted: #6B7694;     /* Gris bleuté foncé */

/* Bordures & Dividers */
--sakura-border-subtle: rgba(255, 255, 255, 0.08);
--sakura-border-default: rgba(255, 255, 255, 0.12);
--sakura-border-strong: rgba(255, 255, 255, 0.20);
```

### Couleurs d'accent

```css
/* Accent primaire — Magenta électrique */
--sakura-accent-magenta-50: #FCE7FF;
--sakura-accent-magenta-100: #F5C7FF;
--sakura-accent-magenta-400: #E879F9;
--sakura-accent-magenta-500: #D946EF;  /* PRIMARY */
--sakura-accent-magenta-600: #C026D3;
--sakura-accent-magenta-900: #701A75;

/* Accent secondaire — Cyan néon */
--sakura-accent-cyan-50: #ECFEFF;
--sakura-accent-cyan-400: #22D3EE;
--sakura-accent-cyan-500: #06B6D4;    /* SECONDARY */
--sakura-accent-cyan-600: #0891B2;

/* Sakura Rose (highlights) */
--sakura-pink-400: #FD8BA0;
--sakura-pink-500: #FB6F8A;

/* Or doux (status, badges) */
--sakura-gold-400: #FBBF24;
--sakura-gold-500: #F59E0B;
```

### Couleurs sémantiques

```css
/* Success */
--sakura-success-bg: rgba(34, 197, 94, 0.10);
--sakura-success-border: rgba(34, 197, 94, 0.30);
--sakura-success-text: #4ADE80;

/* Warning */
--sakura-warning-bg: rgba(251, 191, 36, 0.10);
--sakura-warning-border: rgba(251, 191, 36, 0.30);
--sakura-warning-text: #FBBF24;

/* Error */
--sakura-error-bg: rgba(239, 68, 68, 0.10);
--sakura-error-border: rgba(239, 68, 68, 0.30);
--sakura-error-text: #F87171;

/* Info */
--sakura-info-bg: rgba(6, 182, 212, 0.10);
--sakura-info-border: rgba(6, 182, 212, 0.30);
--sakura-info-text: #22D3EE;
```

---

## ✏️ TYPOGRAPHIE

### Polices

```css
/* Système (corps de texte, UI) */
--font-sans: 'Inter', ui-sans-serif, system-ui, -apple-system, 
             'Segoe UI', Roboto, Arial, sans-serif;

/* Calligraphie japonaise (titres, logo) */
--font-display: 'Noto Serif JP', 'Shippori Mincho', serif;

/* Monospace (code, logs) */
--font-mono: ui-monospace, 'SF Mono', 'Cascadia Code', 
             'Fira Code', Menlo, Monaco, Consolas, monospace;
```

### Échelle typographique

```css
/* Display (page titles) */
--text-display-lg: 32px / 1.2 / 700;     /* Logo, hero */
--text-display: 28px / 1.25 / 700;       /* Page title */

/* Headings */
--text-h1: 24px / 1.3 / 600;             /* Section title */
--text-h2: 20px / 1.35 / 600;            /* Card title */
--text-h3: 16px / 1.4 / 600;             /* Subsection */

/* Body */
--text-body-lg: 16px / 1.5 / 400;        /* Large body */
--text-body: 14px / 1.5 / 400;           /* Default body */
--text-body-sm: 13px / 1.5 / 400;        /* Small body */

/* UI Elements */
--text-label: 13px / 1.4 / 500;          /* Form labels */
--text-caption: 12px / 1.4 / 400;        /* Captions, hints */
--text-overline: 11px / 1.3 / 600;       /* UPPERCASE labels */
```

---

## 📐 SPACING & LAYOUT

### Échelle de spacing (multiple de 4)

```css
--space-1: 4px;
--space-2: 8px;
--space-3: 12px;
--space-4: 16px;
--space-5: 20px;
--space-6: 24px;
--space-8: 32px;
--space-10: 40px;
--space-12: 48px;
--space-16: 64px;
--space-20: 80px;
```

### Border radius

```css
--radius-sm: 8px;      /* Small elements */
--radius-md: 12px;     /* Inputs, buttons */
--radius-lg: 16px;     /* Cards */
--radius-xl: 20px;     /* Large cards */
--radius-full: 9999px; /* Pills, avatars */
```

### Ombres

```css
--shadow-sm: 0 2px 8px rgba(0, 0, 0, 0.12);
--shadow-md: 0 4px 16px rgba(0, 0, 0, 0.20);
--shadow-lg: 0 8px 32px rgba(0, 0, 0, 0.30);
--shadow-xl: 0 16px 64px rgba(0, 0, 0, 0.40);

/* Glow effects */
--glow-magenta: 0 0 24px rgba(217, 70, 239, 0.40);
--glow-cyan: 0 0 24px rgba(6, 182, 212, 0.40);
```

### Z-index

```css
--z-base: 0;
--z-dropdown: 100;
--z-sticky: 200;
--z-overlay: 300;
--z-modal: 400;
--z-toast: 500;
```

---

## 🧱 COMPOSANTS UI

### Buttons

```css
/* Primary button */
.btn-primary {
  background: linear-gradient(135deg, var(--sakura-accent-magenta-500), var(--sakura-accent-magenta-600));
  color: white;
  border: 1px solid var(--sakura-accent-magenta-600);
  border-radius: var(--radius-md);
  padding: 10px 20px;
  font-weight: 600;
  transition: all 0.2s ease;
}
.btn-primary:hover {
  transform: translateY(-2px);
  box-shadow: var(--glow-magenta), var(--shadow-md);
}

/* Secondary button */
.btn-secondary {
  background: var(--sakura-bg-elevated);
  color: var(--sakura-text-primary);
  border: 1px solid var(--sakura-border-default);
  /* ... */
}

/* Ghost button */
.btn-ghost {
  background: transparent;
  color: var(--sakura-text-secondary);
  border: none;
  /* ... */
}
```

### Cards

```css
.card {
  background: linear-gradient(135deg, 
    rgba(255, 255, 255, 0.05), 
    rgba(255, 255, 255, 0.02));
  border: 1px solid var(--sakura-border-default);
  border-radius: var(--radius-lg);
  padding: var(--space-6);
  backdrop-filter: blur(12px);
  transition: all 0.3s ease;
}
.card:hover {
  border-color: var(--sakura-accent-magenta-500);
  box-shadow: var(--shadow-lg);
  transform: translateY(-4px);
}
```

### Inputs

```css
.input {
  background: rgba(0, 0, 0, 0.30);
  border: 1px solid var(--sakura-border-default);
  border-radius: var(--radius-md);
  padding: 12px 16px;
  color: var(--sakura-text-primary);
  font-size: 14px;
  transition: all 0.2s ease;
}
.input:focus {
  outline: none;
  border-color: var(--sakura-accent-magenta-500);
  box-shadow: 0 0 0 4px rgba(217, 70, 239, 0.15);
}
```

### Progress bars

```css
.progress-bar {
  height: 8px;
  background: var(--sakura-bg-surface);
  border-radius: var(--radius-full);
  overflow: hidden;
}
.progress-fill {
  height: 100%;
  background: linear-gradient(90deg, 
    var(--sakura-accent-magenta-500), 
    var(--sakura-accent-cyan-500));
  animation: shimmer 2s infinite;
}

@keyframes shimmer {
  0% { background-position: 0% 50%; }
  100% { background-position: 200% 50%; }
}
```

### Badges

```css
.badge {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 4px 12px;
  border-radius: var(--radius-full);
  font-size: 12px;
  font-weight: 600;
  border: 1px solid;
}

.badge-success {
  background: var(--sakura-success-bg);
  border-color: var(--sakura-success-border);
  color: var(--sakura-success-text);
}
/* ... variants warning, error, info */
```

---

## 🌊 ÉLÉMENTS DÉCORATIFS

### Brush Strokes (SVG)

Traits de pinceau subtils en arrière-plan des sections principales :

```html
<svg class="brush-stroke" viewBox="0 0 1200 400">
  <path d="M0,200 Q300,150 600,200 T1200,200" 
        stroke="url(#gradient-magenta)" 
        stroke-width="80" 
        opacity="0.08" 
        fill="none" />
</svg>
```

### Sakura Petals (Canvas animation)

Animation légère de pétales tombant (performance optimisée avec `requestAnimationFrame`) :

```typescript
// Pseudo-code
class SakuraPetal {
  x, y, rotation, speed, opacity
  update() { /* physics */ }
  draw(ctx) { /* render petal */ }
}
// Max 20 pétales simultanés, respawn au bottom
```

### Wave Pattern (Seigaiha)

Motif ondulé japonais en filigrane (background subtil) :

```css
.wave-pattern {
  background-image: 
    radial-gradient(circle at 50% 100%, 
      transparent 0%, 
      transparent 50%, 
      rgba(217, 70, 239, 0.03) 50%, 
      rgba(217, 70, 239, 0.03) 100%);
  background-size: 40px 40px;
  opacity: 0.3;
}
```

---

## 🎬 ANIMATIONS & TRANSITIONS

### Principes
- **Durée** : 200-300ms (rapide, réactif)
- **Easing** : `cubic-bezier(0.4, 0, 0.2, 1)` (Material ease-out)
- **Réduit si `prefers-reduced-motion`**

### Micro-interactions

```css
/* Button press */
.btn:active {
  transform: scale(0.98);
}

/* Card hover lift */
.card:hover {
  transform: translateY(-4px);
  transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

/* Fade in */
@keyframes fadeIn {
  from { opacity: 0; transform: translateY(8px); }
  to { opacity: 1; transform: translateY(0); }
}
.fade-in {
  animation: fadeIn 0.4s ease-out;
}

/* Shimmer loading */
@keyframes shimmer {
  0% { background-position: -200% 0; }
  100% { background-position: 200% 0; }
}
.loading {
  background: linear-gradient(90deg, 
    transparent, 
    rgba(255,255,255,0.1), 
    transparent);
  background-size: 200% 100%;
  animation: shimmer 1.5s infinite;
}
```

---

## 📱 RESPONSIVE BREAKPOINTS

```css
/* Mobile first */
--breakpoint-sm: 640px;   /* Tablette portrait */
--breakpoint-md: 768px;   /* Tablette landscape */
--breakpoint-lg: 1024px;  /* Desktop */
--breakpoint-xl: 1280px;  /* Large desktop */

/* Grid responsive */
.grid {
  display: grid;
  gap: var(--space-4);
  grid-template-columns: 1fr; /* Mobile */
}
@media (min-width: 768px) {
  .grid { grid-template-columns: repeat(2, 1fr); }
}
@media (min-width: 1024px) {
  .grid { grid-template-columns: repeat(3, 1fr); }
}
```

---

## ♿ ACCESSIBILITÉ

### Contrastes (WCAG AA)
- Texte principal : Ratio ≥ 7:1 (AAA)
- Texte secondaire : Ratio ≥ 4.5:1 (AA)
- Éléments interactifs : Focus visible (outline ou box-shadow)

### Focus states

```css
/* Keyboard navigation */
*:focus-visible {
  outline: 2px solid var(--sakura-accent-magenta-500);
  outline-offset: 2px;
}

/* Skip to main content */
.skip-link {
  position: absolute;
  top: -40px;
  left: 0;
  background: var(--sakura-accent-magenta-500);
  color: white;
  padding: 8px 16px;
  z-index: 1000;
}
.skip-link:focus {
  top: 0;
}
```

### Screen readers

```html
<!-- Aria labels -->
<button aria-label="Télécharger épisode 12">
  <DownloadIcon />
</button>

<!-- Live regions -->
<div role="status" aria-live="polite" aria-atomic="true">
  Téléchargement en cours : 45%
</div>
```

---

## 🎨 MOCKUPS WIREFRAMES

### Page 1 : Dashboard

```
┌─────────────────────────────────────────────────┐
│  🌸 Anime-Sama Downloader    [Search] [Profile] │
├─────────────────────────────────────────────────┤
│                                                 │
│  🆕 NOUVEAUTÉS DE LA SEMAINE                    │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐          │
│  │  Cover  │ │  Cover  │ │  Cover  │          │
│  │  Demon  │ │ Jujutsu │ │   Spy   │          │
│  │ Slayer  │ │ Kaisen  │ │  Family │          │
│  │  S4E12  │ │  S2E18  │ │  S2E03  │          │
│  └─────────┘ └─────────┘ └─────────┘          │
│                                                 │
│  📺 MES ABONNEMENTS                             │
│  ┌────────────────────────────────────┐        │
│  │ One Piece        ⚡ Ep 1102 dispo  │        │
│  │ Attack on Titan  ✅ À jour         │        │
│  │ Tokyo Revengers  ⏳ Ep 24 en cours │        │
│  └────────────────────────────────────┘        │
│                                                 │
│  📊 TÉLÉCHARGEMENTS EN COURS                    │
│  [Progress bar: Demon Slayer S4E12 - 67%]      │
│  [Progress bar: Spy Family S2E03 - 12%]        │
└─────────────────────────────────────────────────┘
```

### Page 2 : Recherche

```
┌─────────────────────────────────────────────────┐
│  🔍 [Rechercher un anime...]          [Filtres] │
├─────────────────────────────────────────────────┤
│                                                 │
│  RÉSULTATS POUR "demon slayer"                  │
│                                                 │
│  ┌─────────────────────────────────────────┐   │
│  │ 🎴 Demon Slayer: Kimetsu no Yaiba      │   │
│  │    Season 4 • VOSTFR • 12 épisodes     │   │
│  │    [Télécharger S4E12] [S'abonner]     │   │
│  └─────────────────────────────────────────┘   │
│                                                 │
│  ┌─────────────────────────────────────────┐   │
│  │ 🎴 Demon Slayer: Mugen Train           │   │
│  │    Film • VOSTFR • 1080p               │   │
│  │    [Télécharger] [Détails]             │   │
│  └─────────────────────────────────────────┘   │
└─────────────────────────────────────────────────┘
```

---

## 📦 IMPLÉMENTATION

### Structure fichiers

```
webapp/src/styles/
├── tokens.css          # Variables CSS (couleurs, spacing, etc.)
├── globals.css         # Reset + base styles
├── components.css      # Composants UI (buttons, cards, etc.)
├── animations.css      # Keyframes & transitions
└── utilities.css       # Classes utilitaires (margins, etc.)
```

### React components

```typescript
// Button.tsx
interface ButtonProps {
  variant: 'primary' | 'secondary' | 'ghost';
  size?: 'sm' | 'md' | 'lg';
  children: React.ReactNode;
  onClick?: () => void;
}

export const Button: React.FC<ButtonProps> = ({ 
  variant, size = 'md', children, onClick 
}) => {
  return (
    <button 
      className={cn('btn', `btn-${variant}`, `btn-${size}`)}
      onClick={onClick}
    >
      {children}
    </button>
  );
};
```

---

**🎯 Prochaine étape** : Implémenter ce design system dans le code (Sprint 2) et créer un Storybook pour documentation.
