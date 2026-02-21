import React from 'react';

export interface SkeletonProps {
  width?: string | number;
  height?: string | number;
  variant?: 'text' | 'circular' | 'rectangular';
  className?: string;
  style?: React.CSSProperties;
}

export function Skeleton({
  width = '100%',
  height = '20px',
  variant = 'rectangular',
  className = '',
  style,
}: SkeletonProps) {
  const variantStyles = {
    text: { borderRadius: 'var(--radius-sm)' },
    circular: { borderRadius: '50%' },
    rectangular: { borderRadius: 'var(--radius-md)' },
  };

  return (
    <div
      className={`skeleton ${className}`}
      style={{
        width,
        height,
        background: 'var(--night-bg-elevated)',
        ...variantStyles[variant],
        animation: 'skeleton-pulse 1.5s ease-in-out infinite',
        ...style,
      }}
    />
  );
}

export function SkeletonCard() {
  return (
    <div
      style={{
        padding: 'var(--space-4)',
        background: 'var(--night-bg-surface)',
        border: '1px solid var(--night-border-default)',
        borderRadius: 'var(--radius-lg)',
      }}
    >
      <Skeleton height="160px" style={{ marginBottom: 'var(--space-3)' }} />
      <Skeleton width="60%" height="24px" style={{ marginBottom: 'var(--space-2)' }} />
      <Skeleton width="40%" height="16px" style={{ marginBottom: 'var(--space-3)' }} />
      <div style={{ display: 'flex', gap: 'var(--space-2)' }}>
        <Skeleton width="80px" height="32px" />
        <Skeleton width="80px" height="32px" />
      </div>
    </div>
  );
}

export function SkeletonList({ count = 3 }: { count?: number }) {
  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: 'var(--space-3)' }}>
      {Array.from({ length: count }, (_, i) => (
        <div
          key={i}
          style={{
            padding: 'var(--space-4)',
            background: 'var(--night-bg-surface)',
            border: '1px solid var(--night-border-default)',
            borderRadius: 'var(--radius-md)',
            display: 'flex',
            gap: 'var(--space-3)',
          }}
        >
          <Skeleton width="80px" height="80px" />
          <div style={{ flex: 1 }}>
            <Skeleton width="70%" height="20px" style={{ marginBottom: 'var(--space-2)' }} />
            <Skeleton width="50%" height="16px" style={{ marginBottom: 'var(--space-2)' }} />
            <Skeleton width="30%" height="16px" />
          </div>
        </div>
      ))}
    </div>
  );
}
