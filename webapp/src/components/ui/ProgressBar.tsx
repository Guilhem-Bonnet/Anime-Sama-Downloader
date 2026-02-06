import React from 'react';

export interface ProgressBarProps {
  value: number;
  max?: number;
  variant?: 'primary' | 'success' | 'warning' | 'error';
  size?: 'sm' | 'md' | 'lg';
  showLabel?: boolean;
  label?: string;
  className?: string;
}

export function ProgressBar({
  value,
  max = 100,
  variant = 'primary',
  size = 'md',
  showLabel = false,
  label,
  className = '',
}: ProgressBarProps) {
  const percentage = Math.min(Math.max((value / max) * 100, 0), 100);

  const variantColors = {
    primary: 'var(--sakura-accent-magenta-500)',
    success: 'var(--sakura-success-text)',
    warning: 'var(--sakura-warning-text)',
    error: 'var(--sakura-error-text)',
  };

  const sizeMap = {
    sm: '6px',
    md: '8px',
    lg: '12px',
  };

  return (
    <div className={className} style={{ width: '100%' }}>
      {(showLabel || label) && (
        <div
          style={{
            display: 'flex',
            justifyContent: 'space-between',
            marginBottom: 'var(--space-2)',
            fontSize: 'var(--text-body-sm)',
            color: 'var(--sakura-text-secondary)',
          }}
        >
          {label && <span>{label}</span>}
          {showLabel && <span>{Math.round(percentage)}%</span>}
        </div>
      )}
      <div
        className="progress-bar-track"
        style={{
          width: '100%',
          height: sizeMap[size],
          background: 'var(--sakura-bg-elevated)',
          borderRadius: 'var(--radius-full)',
          border: '1px solid var(--sakura-border-subtle)',
          overflow: 'hidden',
          position: 'relative',
        }}
      >
        <div
          className="progress-bar-fill"
          style={{
            height: '100%',
            width: `${percentage}%`,
            background: variantColors[variant],
            borderRadius: 'var(--radius-full)',
            transition: 'width var(--transition-normal)',
          }}
        />
      </div>
    </div>
  );
}
