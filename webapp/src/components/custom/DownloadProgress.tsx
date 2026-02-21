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
    <div className={className} style={{ width: '100%' }}>
      <div
        style={{
          height: '10px',
          width: '100%',
          background: 'var(--night-bg-tertiary, #2a2a2a)',
          borderRadius: '8px',
          overflow: 'hidden',
        }}
      >
        <div
          role="progressbar"
          aria-valuenow={clampedProgress}
          aria-valuemin={0}
          aria-valuemax={100}
          style={{
            height: '100%',
            width: `${clampedProgress}%`,
            background: 'linear-gradient(90deg, var(--night-accent-brown-400, #b8860b), var(--night-accent-brown-600, #8b6914))',
            borderRadius: '8px',
            transition: 'width 0.4s ease-out',
            ...(animated && clampedProgress > 0 && clampedProgress < 100
              ? { animation: 'pulse 2s ease-in-out infinite' }
              : {}),
          }}
        />
      </div>
      {showLabel && (
        <p
          style={{
            fontSize: '12px',
            color: 'var(--night-text-secondary, #999)',
            marginTop: '4px',
            margin: '4px 0 0 0',
          }}
        >
          {clampedProgress}%
        </p>
      )}
    </div>
  );
};
