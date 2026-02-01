import React from 'react';
import { useMode } from '../../hooks/useMode';
import { Moon, Circle } from 'lucide-react';

export const ModeToggle: React.FC = () => {
  const [mode, setMode] = useMode();

  return (
    <button
      onClick={() => setMode(mode === 'dark' ? 'light' : 'dark')}
      className="p-2 rounded-lg bg-gray-200 dark:bg-gray-700 text-gray-800 dark:text-white hover:bg-gray-300 dark:hover:bg-gray-600 transition-colors"
      aria-label={`Switch to ${mode === 'dark' ? 'light' : 'dark'} mode`}
      title={mode === 'dark' ? 'Mode clair' : 'Mode sombre'}
    >
      {mode === 'dark' ? <Circle className="w-4 h-4" /> : <Moon className="w-4 h-4" />}
    </button>
  );
};
