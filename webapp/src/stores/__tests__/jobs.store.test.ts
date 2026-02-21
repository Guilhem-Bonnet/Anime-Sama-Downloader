import { describe, it, expect, vi, beforeEach } from 'vitest';
import { useJobsStore, progressPercent, type Job } from '../jobs.store';
import * as api from '../../api';

vi.mock('../../api', () => ({
  apiListJobs: vi.fn(),
  apiCancelJob: vi.fn(),
}));

const mockJob = (overrides: Partial<Job> = {}): Job => ({
  id: 'j-1',
  type: 'download',
  state: 'running',
  progress: 0.5,
  createdAt: '2024-01-01T00:00:00Z',
  updatedAt: '2024-01-01T00:00:00Z',
  ...overrides,
});

describe('progressPercent', () => {
  it('convertit 0.0 → 0', () => {
    expect(progressPercent(mockJob({ progress: 0 }))).toBe(0);
  });

  it('convertit 0.5 → 50', () => {
    expect(progressPercent(mockJob({ progress: 0.5 }))).toBe(50);
  });

  it('convertit 1.0 → 100', () => {
    expect(progressPercent(mockJob({ progress: 1 }))).toBe(100);
  });

  it('arrondit 0.333 → 33', () => {
    expect(progressPercent(mockJob({ progress: 0.333 }))).toBe(33);
  });

  it('gère undefined → 0', () => {
    expect(progressPercent(mockJob({ progress: undefined as any }))).toBe(0);
  });
});

describe('useJobsStore', () => {
  beforeEach(() => {
    useJobsStore.setState({ jobs: [], isLoading: false, error: undefined });
    vi.clearAllMocks();
  });

  describe('loadJobs', () => {
    it('charge les jobs avec succès', async () => {
      const jobs = [mockJob({ id: 'j-1' }), mockJob({ id: 'j-2' })];
      vi.mocked(api.apiListJobs).mockResolvedValue(jobs);

      await useJobsStore.getState().loadJobs();

      expect(api.apiListJobs).toHaveBeenCalledWith(200);
      expect(useJobsStore.getState().jobs).toEqual(jobs);
      expect(useJobsStore.getState().isLoading).toBe(false);
      expect(useJobsStore.getState().error).toBeUndefined();
    });

    it('passe la limite personnalisée', async () => {
      vi.mocked(api.apiListJobs).mockResolvedValue([]);
      await useJobsStore.getState().loadJobs(10);
      expect(api.apiListJobs).toHaveBeenCalledWith(10);
    });

    it('gère les erreurs API', async () => {
      vi.mocked(api.apiListJobs).mockRejectedValue(new Error('réseau'));

      await useJobsStore.getState().loadJobs();

      expect(useJobsStore.getState().error).toBe('réseau');
      expect(useJobsStore.getState().isLoading).toBe(false);
    });

    it('gère les erreurs non-Error', async () => {
      vi.mocked(api.apiListJobs).mockRejectedValue('boom');

      await useJobsStore.getState().loadJobs();

      expect(useJobsStore.getState().error).toBe('Échec du chargement des jobs');
    });
  });

  describe('cancelJob', () => {
    it('met à jour le job annulé dans la liste', async () => {
      const job1 = mockJob({ id: 'j-1', state: 'running' });
      const cancelled = mockJob({ id: 'j-1', state: 'canceled' });
      useJobsStore.setState({ jobs: [job1] });
      vi.mocked(api.apiCancelJob).mockResolvedValue(cancelled);

      await useJobsStore.getState().cancelJob('j-1');

      expect(useJobsStore.getState().jobs[0].state).toBe('canceled');
    });

    it('gère les erreurs d\'annulation', async () => {
      useJobsStore.setState({ jobs: [mockJob()] });
      vi.mocked(api.apiCancelJob).mockRejectedValue(new Error('interdit'));

      await useJobsStore.getState().cancelJob('j-1');

      expect(useJobsStore.getState().error).toBe('interdit');
    });
  });

  describe('updateJobFromSSE', () => {
    it('met à jour un job existant', () => {
      useJobsStore.setState({ jobs: [mockJob({ id: 'j-1', progress: 0.1 })] });

      useJobsStore.getState().updateJobFromSSE({ id: 'j-1', progress: 0.9 });

      expect(useJobsStore.getState().jobs[0].progress).toBe(0.9);
    });

    it('insère un nouveau job en tête', () => {
      useJobsStore.setState({ jobs: [mockJob({ id: 'j-1' })] });

      useJobsStore.getState().updateJobFromSSE({
        id: 'j-new',
        type: 'download',
        state: 'queued',
        progress: 0,
        createdAt: '',
        updatedAt: '',
      });

      const { jobs } = useJobsStore.getState();
      expect(jobs).toHaveLength(2);
      expect(jobs[0].id).toBe('j-new');
    });

    it('ne touche pas les autres jobs', () => {
      useJobsStore.setState({
        jobs: [
          mockJob({ id: 'j-1', progress: 0.1 }),
          mockJob({ id: 'j-2', progress: 0.2 }),
        ],
      });

      useJobsStore.getState().updateJobFromSSE({ id: 'j-1', progress: 0.9 });

      expect(useJobsStore.getState().jobs[1].progress).toBe(0.2);
    });
  });

  describe('setError', () => {
    it('définit et efface l\'erreur', () => {
      useJobsStore.getState().setError('oops');
      expect(useJobsStore.getState().error).toBe('oops');

      useJobsStore.getState().setError(undefined);
      expect(useJobsStore.getState().error).toBeUndefined();
    });
  });
});
