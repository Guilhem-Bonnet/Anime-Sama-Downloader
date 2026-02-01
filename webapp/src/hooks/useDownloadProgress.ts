import { useEffect, useState } from 'react';
import { useJobsStore } from '../stores/jobs.store';
import { useSSE } from './useSSE';

export interface DownloadProgressData {
  jobId: string;
  progress: number;
  estimatedTimeRemaining?: number;
}

export function useDownloadProgress(jobId?: string) {
  const [progress, setProgress] = useState(0);
  const { updateJobProgress } = useJobsStore();

  const { close } = useSSE(
    jobId ? `/api/jobs/${jobId}/progress` : '',
    jobId ? (data: DownloadProgressData) => {
      setProgress(data.progress);
      updateJobProgress(jobId, data.progress);
    } : undefined
  );

  return { progress, close };
}
