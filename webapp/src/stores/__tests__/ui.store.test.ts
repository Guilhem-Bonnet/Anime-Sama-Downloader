import { describe, it, expect, beforeEach } from 'vitest';
import { useUIStore } from '../ui.store';

describe('useUIStore', () => {
  beforeEach(() => {
    useUIStore.setState({
      mode: 'dark',
      activeView: 'search',
      isModalOpen: false,
      modalContent: undefined,
      isLoading: false,
      error: undefined,
    });
  });

  it('a un état initial correct', () => {
    const state = useUIStore.getState();
    expect(state.mode).toBe('dark');
    expect(state.activeView).toBe('search');
    expect(state.isModalOpen).toBe(false);
    expect(state.isLoading).toBe(false);
    expect(state.error).toBeUndefined();
  });

  describe('setMode', () => {
    it('passe en light', () => {
      useUIStore.getState().setMode('light');
      expect(useUIStore.getState().mode).toBe('light');
    });

    it('passe en dark', () => {
      useUIStore.setState({ mode: 'light' });
      useUIStore.getState().setMode('dark');
      expect(useUIStore.getState().mode).toBe('dark');
    });
  });

  describe('setActiveView', () => {
    it.each(['search', 'downloads', 'rules', 'settings'] as const)(
      'met la vue sur %s',
      (view) => {
        useUIStore.getState().setActiveView(view);
        expect(useUIStore.getState().activeView).toBe(view);
      }
    );
  });

  describe('setModalOpen', () => {
    it('ouvre la modale avec du contenu', () => {
      useUIStore.getState().setModalOpen(true, 'confirm-delete');
      const state = useUIStore.getState();
      expect(state.isModalOpen).toBe(true);
      expect(state.modalContent).toBe('confirm-delete');
    });

    it('ferme la modale', () => {
      useUIStore.setState({ isModalOpen: true, modalContent: 'x' });
      useUIStore.getState().setModalOpen(false);
      expect(useUIStore.getState().isModalOpen).toBe(false);
    });
  });

  describe('setLoading', () => {
    it('active le chargement', () => {
      useUIStore.getState().setLoading(true);
      expect(useUIStore.getState().isLoading).toBe(true);
    });
  });

  describe('setError / clearError', () => {
    it('définit une erreur', () => {
      useUIStore.getState().setError('oops');
      expect(useUIStore.getState().error).toBe('oops');
    });

    it('efface l\'erreur', () => {
      useUIStore.setState({ error: 'oops' });
      useUIStore.getState().clearError();
      expect(useUIStore.getState().error).toBeUndefined();
    });
  });
});
