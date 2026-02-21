# Direction artistique — Anime‑Sama Downloader (Web)

## Intention
Une UI “nocturne”, moderne et lisible, orientée productivité : on comprend vite quoi faire (ajouter un abonnement, synchroniser, choisir des épisodes, suivre les jobs).

## Nom de thème
**Aurora Noir** : fond bleu nuit + halo violet discret + surfaces « verre sombre ».

## Palette (tokens)
- **Fond**: `--bg0` (bleu nuit), `--bg1` (bleu profond)
- **Surface**: `--panel` / `--panel2` (verre sombre)
- **Texte**: `--text` (principal), `--muted` (secondaire)
- **Accent**: `--accent` (violet), `--accent2` (cyan/teal léger)
- **États**: `--ok`, `--warn`, `--err`
- **Bordures**: `--border` (trait fin), `--borderStrong` (trait hover)
- **Ombres**: `--shadow1`, `--shadow2`
- **Focus**: `--ring`

## Typographie
- Police système (rapide, nette)
- Hiérarchie :
  - Titre page : 22–24px / gras
  - Titres sections : 14–16px / semi‑gras
  - Texte : 13–14px
  - Aide/secondaire : 12–13px (`--muted`)

## Grille & spacing
- Largeur max : 1100px
- Espacements : 8 / 12 / 16 / 24 px (rythme constant)
- Rayons : 14px (cartes), 12px (inputs/boutons)

## Composants
- **Topbar** : sticky, fond flouté, onglets en “pills” (active = accent)
- **Cards** : surfaces légèrement dégradées + bordure fine + ombre douce
- **Buttons** : hover/active/focus visibles, style primaire = accent
- **Inputs** : focus ring violet, fond semi‑transparent
- **Tables** : en-tête discret, zebra léger, hover de ligne
- **Modal** : overlay sombre + carte centrée, fermeture clic dehors + `Escape`

## Accessibilité
- `:focus-visible` sur boutons et champs (ring)
- Contraste : texte principal clair, `--muted` suffisamment lisible
- `prefers-reduced-motion` : animations réduites
