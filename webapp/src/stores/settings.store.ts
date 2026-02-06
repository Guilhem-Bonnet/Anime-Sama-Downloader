import { create } from 'zustand';
import { persist } from 'zustand/middleware';

export interface Settings {
  downloadFolder: string;
  maxConcurrentDownloads: number;
  preferredQuality: '480p' | '720p' | '1080p';
  enableNotifications: boolean;
  autoDeleteAfterDownload: boolean;
  autoCheckUpdates: boolean;
  theme: 'light' | 'dark';
}

export interface SettingsState extends Settings {
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
  downloadFolder: '/media/anime',
  maxConcurrentDownloads: 2,
  preferredQuality: '1080p',
  enableNotifications: true,
  autoDeleteAfterDownload: false,
  autoCheckUpdates: true,
  theme: 'dark',
};

export const useSettingsStore = create<SettingsState & SettingsActions>()(
  persist(
    (set, get) => ({
      // State
      ...defaultSettings,
      isLoading: false,
      isDirty: false,

      // Actions
      updateSettings: (updates) =>
        set((state) => ({
          ...updates,
          isDirty: true,
        })),

      saveSettings: async () => {
        set({ isLoading: true, error: undefined });
        try {
          const state = get();
          const settings: Settings = {
            downloadFolder: state.downloadFolder,
            maxConcurrentDownloads: state.maxConcurrentDownloads,
            preferredQuality: state.preferredQuality,
            enableNotifications: state.enableNotifications,
            autoDeleteAfterDownload: state.autoDeleteAfterDownload,
            autoCheckUpdates: state.autoCheckUpdates,
            theme: state.theme,
          };

          const response = await fetch('/api/v1/settings', {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(settings),
          });

          if (!response.ok) {
            throw new Error(`Save failed: ${response.statusText}`);
          }

          set({ isDirty: false });
        } catch (error) {
          set({
            error: error instanceof Error ? error.message : 'Save failed',
          });
        } finally {
          set({ isLoading: false });
        }
      },

      loadSettings: async () => {
        set({ isLoading: true, error: undefined });
        try {
          const response = await fetch('/api/v1/settings');
          if (!response.ok) {
            throw new Error(`Load failed: ${response.statusText}`);
          }
          const settings = await response.json();
          set({ ...settings, isDirty: false });
        } catch (error) {
          set({
            error: error instanceof Error ? error.message : 'Load failed',
          });
        } finally {
          set({ isLoading: false });
        }
      },

      resetSettings: () =>
        set({
          ...defaultSettings,
          isDirty: true,
        }),

      setLoading: (isLoading) => set({ isLoading }),

      setError: (error) => set({ error }),
    }),
    {
      name: 'anime-sama-settings',
      partialize: (state) => ({
        downloadFolder: state.downloadFolder,
        maxConcurrentDownloads: state.maxConcurrentDownloads,
        preferredQuality: state.preferredQuality,
        enableNotifications: state.enableNotifications,
        autoDeleteAfterDownload: state.autoDeleteAfterDownload,
        autoCheckUpdates: state.autoCheckUpdates,
        theme: state.theme,
      }),
    }
  )
);
