import { create } from 'zustand';

export interface UIState {
  mode: 'light' | 'dark';
  activeView: 'search' | 'downloads' | 'rules' | 'settings';
  isModalOpen: boolean;
  modalContent?: string;
  isLoading: boolean;
  error?: string;
}

export interface UIActions {
  setMode: (mode: 'light' | 'dark') => void;
  setActiveView: (view: UIState['activeView']) => void;
  setModalOpen: (open: boolean, content?: string) => void;
  setLoading: (loading: boolean) => void;
  setError: (error?: string) => void;
  clearError: () => void;
}

export const useUIStore = create<UIState & UIActions>((set) => ({
  // State
  mode: 'dark',
  activeView: 'search',
  isModalOpen: false,
  isLoading: false,

  // Actions
  setMode: (mode) => set({ mode }),
  setActiveView: (activeView) => set({ activeView }),
  setModalOpen: (isModalOpen, modalContent) => set({ isModalOpen, modalContent }),
  setLoading: (isLoading) => set({ isLoading }),
  setError: (error) => set({ error }),
  clearError: () => set({ error: undefined }),
}));
