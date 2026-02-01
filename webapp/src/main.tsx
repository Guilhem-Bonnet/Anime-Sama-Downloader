import React from 'react';
import ReactDOM from 'react-dom/client';
import { AppRouter } from './AppRouter';
import { Demo } from './Demo';
// Import Tailwind + Sakura Night design system
import './index.css';
import './styles/tokens.css';
import './styles/globals.css';
import './styles/components.css';
import './styles/animations.css';

// Uncomment to view demo or set a flag to toggle
const SHOW_DEMO = localStorage.getItem('SHOW_DESIGN_DEMO') === 'true';

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    {SHOW_DEMO ? <Demo /> : <AppRouter />}
  </React.StrictMode>
);
