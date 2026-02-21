import React from 'react';

export interface DownloadProgressProps {
  progress: number;
  showLabel?: boolean;
  animated?: boolean;
  className?: string;
}

export const DownloadProgress: React.FC<DownloadProgressProps> = ({
  progress,
  showLabel = true,
  animated = true,
  className = '',
}) => {
  const clampedProgress = Math.min(100, Math.max(0, progress));

  return (
    <div className={`w-full ${className}`}>
      <div className="h-2 w-full bg-gray-200 rounded-full overflow-hidden">
        <div
          className={`h-full bg-gradient-to-r from-cyan-500 to-magenta-500 transition-all duration-300 ease-out ${
            animated ? 'animate-pulse' : ''
          }`}
          style={{ width: `${clampedProgress}%` }}
          role="progressbar"
          aria-valuenow={clampedProgress}
          aria-valuemin={0}
          aria-valuemax={100}
        />
      </div>
      {showLabel && <p className="text-xs text-gray-500 mt-1">{clampedProgress}%</p>}
    </div>
  );
};
