import { create } from 'zustand';
import { apiGetSettings, apiPutSettings, type Settings } from '../api';

// Re-export the canonical Settings type from api.ts
export type { Settings } from '../api';

export interface SettingsState {
  settings: Settings;
  isLoading: boolean;
  error?: string;
  isDirty: boolean;
}

export interface SettingsActions {
  updateSettings: (updates: Partial<Settings>) => void;
  saveSettings: () => Promise<void>;
  loadSettings: () => Promise<void>;
  resetSettings: () => void;
  setLoading: (loading: boolean) => void;
  setError: (error?: string) => void;
}

const defaultSettings: Settings = {
  destination: 'videos',
  outputNamingMode: 'legacy',
  separateLang: false,
  maxWorkers: 2,
  maxConcurrentDownloads: 4,
  jellyfinUrl: '',
  jellyfinApiKey: '',
  plexUrl: '',
  plexToken: '',
  plexSectionId: '',
  anilistToken: '',
};

export const useSettingsStore = create<SettingsState & SettingsActions>()(
  (set, get) => ({
    // State
    settings: { ...defaultSettings },
    isLoading: false,
    isDirty: false,

    // Actions
    updateSettings: (updates) =>
      set((state) => ({
        settings: { ...state.settings, ...updates },
        isDirty: true,
      })),

    saveSettings: async () => {
      set({ isLoading: true, error: undefined });
      try {
        const saved = await apiPutSettings(get().settings);
        set({ settings: saved, isDirty: false });
      } catch (error) {
        set({
          error: error instanceof Error ? error.message : 'Échec de la sauvegarde',
        });
      } finally {
        set({ isLoading: false });
      }
    },

    loadSettings: async () => {
      set({ isLoading: true, error: undefined });
      try {
        const settings = await apiGetSettings();
        set({ settings, isDirty: false });
      } catch (error) {
        set({
          error: error instanceof Error ? error.message : 'Échec du chargement',
        });
      } finally {
        set({ isLoading: false });
      }
    },

    resetSettings: () =>
      set({
        settings: { ...defaultSettings },
        isDirty: true,
      }),

    setLoading: (isLoading) => set({ isLoading }),

    setError: (error) => set({ error }),
  })
);
