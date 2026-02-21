# 🌸 Sakura Night Design System — Quick Start

## Voir la Démo Complète

Pour afficher la page de démo interactive (Storybook-like), suivez ces étapes :

### Option 1 : Via localStorage (Recommandé)

1. Ouvrir l'app dans le navigateur (http://localhost:5173 si dev)
2. Ouvrir la **Console du navigateur** (F12 ou Cmd+Option+J)
3. Copier-coller :
```javascript
localStorage.setItem('SHOW_DESIGN_DEMO', 'true');
location.reload();
```
4. La démo s'affichera à la place de l'app principale

### Option 2 : Retour à l'App

Pour revenir à l'application normale :
```javascript
localStorage.removeItem('SHOW_DESIGN_DEMO');
location.reload();
```

---

## Composants Disponibles

Tous les composants sont dans `webapp/src/components/ui/` :

### Button
```tsx
import { Button } from '@/components/ui';

<Button variant="primary">Primary</Button>
<Button variant="secondary" size="sm">Small</Button>
<Button variant="danger" isLoading>Loading...</Button>
```

**Variants** : `primary`, `secondary`, `ghost`, `danger`  
**Sizes** : `sm`, `md`, `lg`  
**Props** : `isLoading`, `disabled`, tout ce que supporte `<button>`

### Card
```tsx
import { Card, CardHeader, CardTitle, CardBody, CardFooter } from '@/components/ui';

<Card>
  <CardHeader>
    <CardTitle>Mon Card</CardTitle>
  </CardHeader>
  <CardBody>Contenu</CardBody>
  <CardFooter>
    <Button>Action</Button>
  </CardFooter>
</Card>
```

### Input, TextArea, Select
```tsx
import { Input, TextArea, Select } from '@/components/ui';

<Input 
  label="Email" 
  type="email"
  placeholder="mail@example.com"
  hint="Un email valide svp"
  error={error && "Email invalide"}
/>

<TextArea label="Message" placeholder="Tapez ici..." />

<Select 
  label="Choix"
  options={[
    { value: 'a', label: 'Option A' },
    { value: 'b', label: 'Option B' }
  ]}
/>
```

### Badge
```tsx
import { Badge } from '@/components/ui';

<Badge variant="primary">Primary</Badge>
<Badge variant="success">Done</Badge>
<Badge variant="error">Error</Badge>
```

**Variants** : `primary`, `secondary`, `success`, `warning`, `error`, `info`

---

## Variables CSS

Tous les tokens sont dans `webapp/src/styles/tokens.css`.

### Utiliser dans vos styles :

```css
.my-element {
  /* Couleurs */
  color: var(--sakura-text-primary);
  background: var(--sakura-bg-surface);
  border-color: var(--sakura-border-default);
  
  /* Spacing */
  padding: var(--space-4);
  margin: var(--space-2);
  
  /* Typography */
  font-size: var(--text-body);
  font-weight: 600;
  
  /* Effects */
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-md);
  transition: all var(--transition-fast);
}

/* Accents */
.primary {
  background: var(--sakura-accent-magenta-500);
  color: white;
}

.secondary {
  background: var(--sakura-accent-cyan-500);
  color: white;
}
```

### Palette Complète

```
🌙 Backgrounds:
  --sakura-bg-base      #0A0E1A
  --sakura-bg-surface   #1A1F2E
  --sakura-bg-elevated  #252A3B

📝 Text:
  --sakura-text-primary      #F5F7FF
  --sakura-text-secondary    #A8B3D1
  --sakura-text-muted        #6B7694

⚡ Accents:
  --sakura-accent-magenta-500  #D946EF
  --sakura-accent-cyan-500     #06B6D4
  --sakura-pink-500            #FB6F8A
  --sakura-gold-500            #F59E0B

✅ Semantic:
  --sakura-success-text   #4ADE80
  --sakura-warning-text   #FBBF24
  --sakura-error-text     #F87171
  --sakura-info-text      #22D3EE
```

---

## Classnames Utilitaires

Des classes utilitaires sont disponibles dans `webapp/src/styles/components.css` :

```html
<!-- Flexbox -->
<div class="flex gap-4">
  <div class="flex-col">Colonne</div>
</div>

<!-- Spacing -->
<div class="p-4 m-4">Padding + Margin</div>
<div class="mt-4 mb-4">Margin top + bottom</div>

<!-- Grid -->
<div class="grid">
  <div>Item 1</div>
  <div>Item 2</div>
</div>
```

---

## Animations

Animations CSS disponibles dans `webapp/src/styles/animations.css` :

```css
.fade-in { animation: fadeIn 0.4s ease-out; }
.slide-in-left { animation: slideInLeft 0.3s ease-out; }
.pulse { animation: pulse 2s infinite; }
.spin { animation: spin 1s linear infinite; }
.loading { background: linear-gradient(...); animation: shimmer 1.5s infinite; }
.glow { animation: glow 2s ease-in-out infinite; }
```

---

## Fichiers Importants

```
webapp/src/
├── styles/
│   ├── tokens.css        ← Variables CSS (couleurs, spacing, etc.)
│   ├── globals.css       ← Reset CSS + styles de base
│   ├── components.css    ← Styles composants (buttons, cards, etc.)
│   └── animations.css    ← Animations et keyframes
├── components/
│   └── ui/
│       ├── Button.tsx    ← Composant Button
│       ├── Card.tsx      ← Composant Card
│       ├── Input.tsx     ← Composant Input, TextArea, Select
│       ├── Badge.tsx     ← Composant Badge
│       └── index.ts      ← Export centralisé
├── Demo.tsx              ← Page démo complète
├── App.tsx               ← App principale
├── main.tsx              ← Entry point
└── api.ts                ← API client
```

---

## Build & Déploiement

Le build TypeScript est **zéro-erreur** ✅

```bash
# Build
cd webapp && npm run build

# Résultat
dist/assets/index.css   11.62 kB (gzipped: 3.14 kB)
dist/assets/index.js    206.14 kB (gzipped: 59.95 kB)
```

---

## Prochaines Étapes

- [ ] Intégrer Storybook pour documentation interactive
- [ ] Ajouter Component tests (React Testing Library)
- [ ] Implémenter dark/light theme toggle
- [ ] WCAG AA accessibilité audit
- [ ] Ajouter plus de composants (Toggle, Checkbox, Modal, etc.)

---

**Quick Flow Sprint COMPLETED** ✅  
Sakura Night Design System — Ready for Production

Contact : Barry (Quick Flow Dev)
