import { describe, it, expect, vi, beforeEach } from 'vitest';
import { useSubscriptionsStore } from '../subscriptions.store';
import * as api from '../../api';

vi.mock('../../api', () => ({
  apiListSubscriptions: vi.fn(),
  apiCreateSubscription: vi.fn(),
  apiDeleteSubscription: vi.fn(),
  apiSyncSubscription: vi.fn(),
  apiSyncAll: vi.fn(),
}));

const store = useSubscriptionsStore;

const mockSub = (overrides: Partial<api.Subscription> = {}): api.Subscription => ({
  id: 's-1',
  baseUrl: 'https://anime-sama.fr/catalogue/naruto/',
  label: 'Naruto',
  player: 'sibnet',
  lastScheduledEpisode: 10,
  lastDownloadedEpisode: 8,
  lastAvailableEpisode: 220,
  nextCheckAt: '2024-06-01T00:00:00Z',
  lastCheckedAt: '2024-05-01T00:00:00Z',
  createdAt: '2024-01-01T00:00:00Z',
  updatedAt: '2024-05-01T00:00:00Z',
  ...overrides,
});

describe('useSubscriptionsStore', () => {
  beforeEach(() => {
    store.setState({
      subscriptions: [],
      isLoading: false,
      error: undefined,
      filter: 'all',
      sortBy: 'label',
    });
    vi.clearAllMocks();
  });

  describe('loadSubscriptions', () => {
    it('charge les abonnements avec succès', async () => {
      const subs = [mockSub({ id: 's-1' }), mockSub({ id: 's-2', label: 'Bleach' })];
      vi.mocked(api.apiListSubscriptions).mockResolvedValue(subs);

      await store.getState().loadSubscriptions();

      expect(store.getState().subscriptions).toEqual(subs);
      expect(store.getState().isLoading).toBe(false);
    });

    it('gère les erreurs', async () => {
      vi.mocked(api.apiListSubscriptions).mockRejectedValue(new Error('500'));

      await store.getState().loadSubscriptions();

      expect(store.getState().error).toBe('500');
    });
  });

  describe('addSubscription', () => {
    it('ajoute un abonnement à la liste', async () => {
      const newSub = mockSub({ id: 's-new' });
      vi.mocked(api.apiCreateSubscription).mockResolvedValue(newSub);

      await store.getState().addSubscription({
        baseUrl: 'https://anime-sama.fr/catalogue/naruto/',
      });

      expect(store.getState().subscriptions).toHaveLength(1);
      expect(store.getState().subscriptions[0].id).toBe('s-new');
    });

    it('gère les erreurs de création', async () => {
      vi.mocked(api.apiCreateSubscription).mockRejectedValue(new Error('duplicate'));

      await store.getState().addSubscription({ baseUrl: 'x' });

      expect(store.getState().error).toBe('duplicate');
    });
  });

  describe('removeSubscription', () => {
    it('supprime un abonnement de la liste', async () => {
      store.setState({ subscriptions: [mockSub({ id: 's-1' }), mockSub({ id: 's-2' })] });
      vi.mocked(api.apiDeleteSubscription).mockResolvedValue(undefined);

      await store.getState().removeSubscription('s-1');

      expect(store.getState().subscriptions).toHaveLength(1);
      expect(store.getState().subscriptions[0].id).toBe('s-2');
    });
  });

  describe('syncSubscription', () => {
    it('met à jour l\'abonnement après sync', async () => {
      const original = mockSub({ id: 's-1', lastAvailableEpisode: 10 });
      const synced = mockSub({ id: 's-1', lastAvailableEpisode: 15 });
      store.setState({ subscriptions: [original] });

      const result: api.SyncResult = {
        subscription: synced,
        selectedPlayer: 'sibnet',
        maxAvailableEpisode: 15,
        enqueuedEpisodes: [11, 12, 13, 14, 15],
        enqueuedJobIDs: ['j-1'],
        message: 'OK',
      };
      vi.mocked(api.apiSyncSubscription).mockResolvedValue(result);

      const res = await store.getState().syncSubscription('s-1');

      expect(res?.subscription.lastAvailableEpisode).toBe(15);
      expect(store.getState().subscriptions[0].lastAvailableEpisode).toBe(15);
    });

    it('retourne undefined en cas d\'erreur', async () => {
      store.setState({ subscriptions: [mockSub()] });
      vi.mocked(api.apiSyncSubscription).mockRejectedValue(new Error('timeout'));

      const res = await store.getState().syncSubscription('s-1');

      expect(res).toBeUndefined();
      expect(store.getState().error).toBe('timeout');
    });
  });

  describe('filtre et tri', () => {
    it('setFilter change le filtre', () => {
      store.getState().setFilter('active');
      expect(store.getState().filter).toBe('active');
    });

    it('setSortBy change le tri', () => {
      store.getState().setSortBy('next-check');
      expect(store.getState().sortBy).toBe('next-check');
    });
  });
});
