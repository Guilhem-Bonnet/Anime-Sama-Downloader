import { create } from 'zustand';
import { persist } from 'zustand/middleware';

export interface Subscription {
  id: string;
  anime_id: string;
  anime_title: string;
  check_interval_minutes: number;
  last_scheduled_episode: number;
  next_check_at: string;
  created_at: string;
  updated_at: string;
}

export interface SubscriptionState {
  subscriptions: Subscription[];
  isLoading: boolean;
  error?: string;
  filter: 'all' | 'active' | 'new-episodes';
  sortBy: 'name' | 'next-check' | 'last-episode';
}

export interface SubscriptionActions {
  setSubscriptions: (subs: Subscription[]) => void;
  addSubscription: (sub: Subscription) => void;
  updateSubscription: (id: string, updates: Partial<Subscription>) => void;
  removeSubscription: (id: string) => void;
  setLoading: (loading: boolean) => void;
  setError: (error?: string) => void;
  setFilter: (filter: SubscriptionState['filter']) => void;
  setSortBy: (sortBy: SubscriptionState['sortBy']) => void;
  syncSubscription: (id: string) => Promise<void>;
  bulkSync: (ids: string[]) => Promise<void>;
}

export const useSubscriptionsStore = create<SubscriptionState & SubscriptionActions>((set, get) => ({
  // State
  subscriptions: [],
  isLoading: false,
  filter: 'all',
  sortBy: 'name',

  // Actions
  setSubscriptions: (subscriptions) => set({ subscriptions }),
  
  addSubscription: (subscription) =>
    set((state) => ({
      subscriptions: [...state.subscriptions, subscription],
    })),

  updateSubscription: (id, updates) =>
    set((state) => ({
      subscriptions: state.subscriptions.map((sub) =>
        sub.id === id ? { ...sub, ...updates } : sub
      ),
    })),

  removeSubscription: (id) =>
    set((state) => ({
      subscriptions: state.subscriptions.filter((sub) => sub.id !== id),
    })),

  setLoading: (isLoading) => set({ isLoading }),

  setError: (error) => set({ error }),

  setFilter: (filter) => set({ filter }),

  setSortBy: (sortBy) => set({ sortBy }),

  syncSubscription: async (id) => {
    const state = get();
    set({ isLoading: true, error: undefined });
    try {
      const response = await fetch(`/api/v1/subscriptions/${id}/sync`, {
        method: 'POST',
      });
      if (!response.ok) {
        throw new Error(`Sync failed: ${response.statusText}`);
      }
      const updated = await response.json();
      state.updateSubscription(id, updated);
    } catch (error) {
      set({ error: error instanceof Error ? error.message : 'Sync failed' });
    } finally {
      set({ isLoading: false });
    }
  },

  bulkSync: async (ids) => {
    set({ isLoading: true, error: undefined });
    try {
      await Promise.all(ids.map((id) => get().syncSubscription(id)));
    } finally {
      set({ isLoading: false });
    }
  },
}));
