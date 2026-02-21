import { create } from 'zustand';
import { apiListJobs, apiCancelJob, type Job, type JobState } from '../api';

// Re-export the canonical types from api.ts
export type { Job, JobState } from '../api';

export interface JobsState {
  jobs: Job[];
  isLoading: boolean;
  error?: string;
}

export interface JobsActions {
  loadJobs: (limit?: number) => Promise<void>;
  cancelJob: (id: string) => Promise<void>;
  updateJobFromSSE: (data: Partial<Job> & { id: string }) => void;
  setError: (error?: string) => void;
}

/** Map backend progress (0.0-1.0) to display percentage (0-100) */
export function progressPercent(job: Job): number {
  return Math.round((job.progress ?? 0) * 100);
}

export const useJobsStore = create<JobsState & JobsActions>((set, get) => ({
  // State
  jobs: [],
  isLoading: false,

  // Actions
  loadJobs: async (limit = 200) => {
    set({ isLoading: true, error: undefined });
    try {
      const jobs = await apiListJobs(limit);
      set({ jobs });
    } catch (error) {
      set({ error: error instanceof Error ? error.message : 'Échec du chargement des jobs' });
    } finally {
      set({ isLoading: false });
    }
  },

  cancelJob: async (id) => {
    try {
      const updated = await apiCancelJob(id);
      set((state) => ({
        jobs: state.jobs.map((j) => (j.id === id ? updated : j)),
      }));
    } catch (error) {
      set({ error: error instanceof Error ? error.message : 'Échec de l\'annulation' });
    }
  },

  updateJobFromSSE: (data) => {
    set((state) => {
      const exists = state.jobs.some((j) => j.id === data.id);
      if (!exists) {
        // New job from SSE — insert at top
        return { jobs: [data as Job, ...state.jobs] };
      }
      return {
        jobs: state.jobs.map((j) =>
          j.id === data.id ? { ...j, ...data } : j
        ),
      };
    });
  },

  setError: (error) => set({ error }),
}));
