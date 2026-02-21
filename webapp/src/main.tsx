import React from 'react';
import ReactDOM from 'react-dom/client';
import { AppRouter } from './AppRouter';
import { Demo } from './Demo';
// Import Tailwind + Nocturne Core design system
import './index.css';
import './styles/tokens.css';
import './styles/globals.css';
import './styles/components.css';
import './styles/animations.css';

// Uncomment to view demo or set a flag to toggle
const getShowDemo = () => {
  try {
    return typeof window !== 'undefined' && localStorage.getItem('SHOW_DESIGN_DEMO') === 'true';
  } catch {
    return false;
  }
};

const SHOW_DEMO = getShowDemo();

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    {SHOW_DEMO ? <Demo /> : <AppRouter />}
  </React.StrictMode>
);
