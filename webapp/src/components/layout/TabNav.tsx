import React from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { Button } from '../ui/Button';
import { Search, Download, Settings, Zap } from 'lucide-react';

type ViewType = 'search' | 'downloads' | 'rules' | 'settings';

const navItems: { path: string; view: ViewType; icon: React.ReactNode; label: string }[] = [
  { path: '/', view: 'search', icon: <Search className="w-5 h-5" />, label: 'Rechercher' },
  { path: '/downloads', view: 'downloads', icon: <Download className="w-5 h-5" />, label: 'Téléchargements' },
  { path: '/rules', view: 'rules', icon: <Settings className="w-5 h-5" />, label: 'Règles' },
  { path: '/settings', view: 'settings', icon: <Zap className="w-5 h-5" />, label: 'Paramètres' },
];

export function TabNav() {
  const navigate = useNavigate();
  const location = useLocation();

  return (
    <nav
      style={{
        position: 'sticky',
        top: '86px',
        zIndex: 30,
        borderBottom: '1px solid var(--sakura-border-default)',
        background: 'rgba(10,14,26,0.85)',
        backdropFilter: 'blur(12px)',
      }}
    >
      <div
        className="flex"
        style={{
          maxWidth: '1200px',
          margin: '0 auto',
          padding: '12px 24px',
          gap: '8px',
          flexWrap: 'wrap',
        }}
      >
        {navItems.map(({ path, view, icon, label }) => {
          const isActive = location.pathname === path;
          return (
            <Button
              key={view}
              variant={isActive ? 'primary' : 'ghost'}
              onClick={() => navigate(path)}
              style={{
                padding: '10px 16px',
                borderBottom: isActive ? `2px solid var(--sakura-accent-brown-500)` : '2px solid transparent',
                color: isActive ? 'var(--sakura-accent-brown-500)' : 'var(--sakura-text-secondary)',
                transition: 'all 200ms ease',
                fontWeight: isActive ? 600 : 500,
              }}
            >
              <span style={{ marginRight: '8px' }}>{icon}</span>
              {label}
            </Button>
          );
        })}
      </div>
    </nav>
  );
}
