import { useEffect, useState } from 'react';
import { useUIStore } from '../stores/ui.store';

type Theme = 'light' | 'dark';

export function useMode(): [Theme, (mode: Theme) => void] {
  const { mode, setMode } = useUIStore();
  const [isDarkMode, setIsDarkMode] = useState(mode === 'dark');

  useEffect(() => {
    // Apply theme to document
    document.documentElement.setAttribute('data-theme', mode);
    
    // Sync system preference
    if (mode === 'dark') {
      document.documentElement.classList.add('dark');
    } else {
      document.documentElement.classList.remove('dark');
    }
  }, [mode]);

  const toggleMode = (newMode: Theme) => {
    setMode(newMode);
    setIsDarkMode(newMode === 'dark');
  };

  return [mode as Theme, toggleMode];
}
