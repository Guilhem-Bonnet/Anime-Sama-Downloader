import { useEffect, useState } from 'react';
import { useJobsStore } from '../stores/jobs.store';
import { useSSE } from './useSSE';
import type { Job, JobState } from '../api';

export interface DownloadProgressData {
  jobId: string;
  progress: number;
  estimatedTimeRemaining?: number;
}

export function useDownloadProgress(jobId?: string) {
  const [progress, setProgress] = useState(0);
  const { updateJobFromSSE } = useJobsStore();

  // Single SSE connection — handles both progress and completion events
  const { close } = useSSE(
    jobId ? `/api/v1/events` : '',
    jobId ? (data: { id: string; progress?: number; state?: JobState }) => {
      if (data.id !== jobId) return;
      const pct = Math.round((data.progress || 0) * 100);
      setProgress(pct);
      updateJobFromSSE(data as Partial<Job> & { id: string });
    } : undefined,
  );

  return { progress, close };
}
