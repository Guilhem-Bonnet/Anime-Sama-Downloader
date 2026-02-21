import React from 'react';

export interface StatusBadgeProps {
  status: 'pending' | 'queued' | 'running' | 'muxing' | 'completed' | 'failed' | 'canceled' | string;
  size?: 'sm' | 'md' | 'lg';
  className?: string;
}

const statusColors: Record<string, string> = {
  pending: 'bg-gray-600 text-white',
  queued: 'bg-gray-600 text-white',
  running: 'bg-blue-600 text-white',
  muxing: 'bg-indigo-600 text-white',
  completed: 'bg-green-600 text-white',
  failed: 'bg-red-600 text-white',
  canceled: 'bg-yellow-600 text-white',
};

export const StatusBadge: React.FC<StatusBadgeProps> = ({ status, size = 'md', className = '' }) => {
  const colorClass = statusColors[status] || 'bg-gray-600 text-white';

  const sizeClasses = {
    sm: 'px-2 py-1 text-xs',
    md: 'px-3 py-1.5 text-sm',
    lg: 'px-4 py-2 text-base',
  };

  return (
    <span
      className={`inline-flex items-center rounded-full font-medium ${colorClass} ${sizeClasses[size]} ${className}`}
    >
      {status}
    </span>
  );
};
