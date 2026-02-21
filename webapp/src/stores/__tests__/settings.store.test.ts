import { describe, it, expect, vi, beforeEach } from 'vitest';
import { useSettingsStore } from '../settings.store';
import * as api from '../../api';

vi.mock('../../api', () => ({
  apiGetSettings: vi.fn(),
  apiPutSettings: vi.fn(),
}));

const defaultSettings = {
  destination: 'videos',
  outputNamingMode: 'legacy' as const,
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

describe('useSettingsStore', () => {
  beforeEach(() => {
    useSettingsStore.setState({
      settings: { ...defaultSettings },
      isLoading: false,
      isDirty: false,
      error: undefined,
    });
    vi.clearAllMocks();
  });

  it('a les valeurs par défaut correctes', () => {
    const state = useSettingsStore.getState();
    expect(state.settings.destination).toBe('videos');
    expect(state.settings.maxWorkers).toBe(2);
    expect(state.isDirty).toBe(false);
  });

  describe('updateSettings', () => {
    it('met à jour partiellement les settings et marque dirty', () => {
      useSettingsStore.getState().updateSettings({ maxWorkers: 8 });

      const state = useSettingsStore.getState();
      expect(state.settings.maxWorkers).toBe(8);
      expect(state.settings.destination).toBe('videos');
      expect(state.isDirty).toBe(true);
    });
  });

  describe('loadSettings', () => {
    it('charge les settings depuis l\'API', async () => {
      const remote = { ...defaultSettings, maxWorkers: 4 };
      vi.mocked(api.apiGetSettings).mockResolvedValue(remote);

      await useSettingsStore.getState().loadSettings();

      expect(useSettingsStore.getState().settings.maxWorkers).toBe(4);
      expect(useSettingsStore.getState().isDirty).toBe(false);
      expect(useSettingsStore.getState().isLoading).toBe(false);
    });

    it('gère les erreurs de chargement', async () => {
      vi.mocked(api.apiGetSettings).mockRejectedValue(new Error('offline'));

      await useSettingsStore.getState().loadSettings();

      expect(useSettingsStore.getState().error).toBe('offline');
      expect(useSettingsStore.getState().isLoading).toBe(false);
    });
  });

  describe('saveSettings', () => {
    it('sauvegarde et réinitialise dirty', async () => {
      const saved = { ...defaultSettings, maxWorkers: 8 };
      vi.mocked(api.apiPutSettings).mockResolvedValue(saved);
      useSettingsStore.setState({
        settings: saved,
        isDirty: true,
      });

      await useSettingsStore.getState().saveSettings();

      expect(api.apiPutSettings).toHaveBeenCalledWith(saved);
      expect(useSettingsStore.getState().isDirty).toBe(false);
      expect(useSettingsStore.getState().isLoading).toBe(false);
    });

    it('gère les erreurs de sauvegarde', async () => {
      vi.mocked(api.apiPutSettings).mockRejectedValue(new Error('403'));

      await useSettingsStore.getState().saveSettings();

      expect(useSettingsStore.getState().error).toBe('403');
    });
  });

  describe('resetSettings', () => {
    it('rétablit les valeurs par défaut et marque dirty', () => {
      useSettingsStore.setState({
        settings: { ...defaultSettings, maxWorkers: 99 },
        isDirty: false,
      });

      useSettingsStore.getState().resetSettings();

      expect(useSettingsStore.getState().settings.maxWorkers).toBe(2);
      expect(useSettingsStore.getState().isDirty).toBe(true);
    });
  });
});
