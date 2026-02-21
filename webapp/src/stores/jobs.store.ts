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
  addJobs: (jobs: Job[]) => void;
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
      const fetched = await apiListJobs(limit);
      // API is source of truth — replace store jobs with fresh data
      set({ jobs: fetched });
    } catch (error) {
      // Don't clear existing jobs on failure — stale data is better than nothing
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

  addJobs: (newJobs) => {
    set((state) => {
      const existingIds = new Set(state.jobs.map((j) => j.id));
      const fresh = newJobs.filter((j) => !existingIds.has(j.id));
      return { jobs: [...fresh, ...state.jobs] };
    });
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
