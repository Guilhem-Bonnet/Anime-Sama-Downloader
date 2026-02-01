import { create } from 'zustand';

export interface Job {
  id: string;
  status: 'pending' | 'running' | 'completed' | 'failed';
  progress: number;
  downloadId: string;
  animeId: string;
  episodeNumber: number;
  errorMessage?: string;
  startedAt?: string;
  completedAt?: string;
}

export interface JobsState {
  jobs: Job[];
  activeJobIds: Set<string>;
  jobProgress: Map<string, number>;
}

export interface JobsActions {
  addJob: (job: Job) => void;
  updateJobProgress: (jobId: string, progress: number) => void;
  updateJobStatus: (jobId: string, status: Job['status'], error?: string) => void;
  removeJob: (jobId: string) => void;
  getJob: (jobId: string) => Job | undefined;
  subscribeToProgress: (callback: (jobId: string, progress: number) => void) => () => void;
}

export const useJobsStore = create<JobsState & JobsActions>((set, get) => {
  const progressSubscribers: Array<(jobId: string, progress: number) => void> = [];

  return {
    // State
    jobs: [],
    activeJobIds: new Set(),
    jobProgress: new Map(),

    // Actions
    addJob: (job) => {
      set((state) => ({
        jobs: [job, ...state.jobs],
        activeJobIds: new Set([...state.activeJobIds, job.id]),
      }));
    },

    updateJobProgress: (jobId, progress) => {
      set((state) => ({
        jobProgress: new Map(state.jobProgress).set(jobId, progress),
      }));
      progressSubscribers.forEach((cb) => cb(jobId, progress));
    },

    updateJobStatus: (jobId, status, error) => {
      set((state) => ({
        jobs: state.jobs.map((job) =>
          job.id === jobId ? { ...job, status, errorMessage: error } : job
        ),
        activeJobIds: status === 'completed' || status === 'failed' 
          ? new Set([...state.activeJobIds].filter((id) => id !== jobId))
          : state.activeJobIds,
      }));
    },

    removeJob: (jobId) => {
      set((state) => ({
        jobs: state.jobs.filter((j) => j.id !== jobId),
        activeJobIds: new Set([...state.activeJobIds].filter((id) => id !== jobId)),
      }));
    },

    getJob: (jobId) => {
      return get().jobs.find((j) => j.id === jobId);
    },

    subscribeToProgress: (callback) => {
      progressSubscribers.push(callback);
      return () => {
        progressSubscribers.splice(progressSubscribers.indexOf(callback), 1);
      };
    },
  };
});
