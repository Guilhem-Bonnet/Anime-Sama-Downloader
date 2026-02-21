import { create } from 'zustand';
import {
  apiListSubscriptions,
  apiCreateSubscription,
  apiDeleteSubscription,
  apiSyncSubscription,
  apiSyncAll,
  type Subscription,
  type SyncResult,
} from '../api';

// Re-export the canonical type from api.ts
export type { Subscription } from '../api';

export interface SubscriptionState {
  subscriptions: Subscription[];
  isLoading: boolean;
  error?: string;
  filter: 'all' | 'active' | 'new-episodes';
  sortBy: 'label' | 'next-check' | 'last-episode';
}

export interface SubscriptionActions {
  loadSubscriptions: () => Promise<void>;
  addSubscription: (params: { baseUrl: string; label?: string; player?: string }) => Promise<void>;
  removeSubscription: (id: string) => Promise<void>;
  syncSubscription: (id: string) => Promise<SyncResult | undefined>;
  syncAll: () => Promise<void>;
  setFilter: (filter: SubscriptionState['filter']) => void;
  setSortBy: (sortBy: SubscriptionState['sortBy']) => void;
  setError: (error?: string) => void;
}

export const useSubscriptionsStore = create<SubscriptionState & SubscriptionActions>((set, get) => ({
  // State
  subscriptions: [],
  isLoading: false,
  filter: 'all',
  sortBy: 'label',

  // Actions
  loadSubscriptions: async () => {
    set({ isLoading: true, error: undefined });
    try {
      const subscriptions = await apiListSubscriptions();
      set({ subscriptions });
    } catch (error) {
      set({ error: error instanceof Error ? error.message : 'Échec du chargement' });
    } finally {
      set({ isLoading: false });
    }
  },

  addSubscription: async (params) => {
    set({ isLoading: true, error: undefined });
    try {
      const sub = await apiCreateSubscription(params);
      set((state) => ({ subscriptions: [...state.subscriptions, sub] }));
    } catch (error) {
      set({ error: error instanceof Error ? error.message : 'Échec de la création' });
    } finally {
      set({ isLoading: false });
    }
  },

  removeSubscription: async (id) => {
    set({ isLoading: true, error: undefined });
    try {
      await apiDeleteSubscription(id);
      set((state) => ({
        subscriptions: state.subscriptions.filter((s) => s.id !== id),
      }));
    } catch (error) {
      set({ error: error instanceof Error ? error.message : 'Échec de la suppression' });
    } finally {
      set({ isLoading: false });
    }
  },

  syncSubscription: async (id) => {
    set({ isLoading: true, error: undefined });
    try {
      const result = await apiSyncSubscription(id);
      // Update the subscription in-place with the returned data
      set((state) => ({
        subscriptions: state.subscriptions.map((s) =>
          s.id === id ? result.subscription : s
        ),
      }));
      return result;
    } catch (error) {
      set({ error: error instanceof Error ? error.message : 'Échec de la synchronisation' });
      return undefined;
    } finally {
      set({ isLoading: false });
    }
  },

  syncAll: async () => {
    set({ isLoading: true, error: undefined });
    try {
      const { results, errors } = await apiSyncAll({ enqueue: true });
      // Refresh the full list after sync
      const subscriptions = await apiListSubscriptions();
      set({ subscriptions });
      if (errors.length > 0) {
        set({ error: `${errors.length} erreur(s) lors de la synchronisation` });
      }
    } catch (error) {
      set({ error: error instanceof Error ? error.message : 'Échec de la synchronisation' });
    } finally {
      set({ isLoading: false });
    }
  },

  setFilter: (filter) => set({ filter }),
  setSortBy: (sortBy) => set({ sortBy }),
  setError: (error) => set({ error }),
}));
