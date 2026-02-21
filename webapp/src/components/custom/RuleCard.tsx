import React from 'react';
import { useJobsStore } from '../../stores/jobs.store';
import { useUIStore } from '../../stores/ui.store';
import { useSSE } from '../../hooks/useSSE';

export interface RuleCardProps {
  id: string;
  name: string;
  enabled: boolean;
  animePattern: string;
  onToggle: (enabled: boolean) => void;
  onEdit: () => void;
  onDelete: () => void;
}

export const RuleCard: React.FC<RuleCardProps> = ({
  id,
  name,
  enabled,
  animePattern,
  onToggle,
  onEdit,
  onDelete,
}) => {
  return (
    <div className="p-4 bg-gray-100 dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700">
      <div className="flex justify-between items-start mb-3">
        <div className="flex-1">
          <h3 className="font-semibold text-gray-900 dark:text-white">{name}</h3>
          <p className="text-sm text-gray-600 dark:text-gray-400 mt-1">Pattern: {animePattern}</p>
        </div>
        <label className="flex items-center gap-2 cursor-pointer">
          <input
            type="checkbox"
            checked={enabled}
            onChange={(e) => onToggle(e.target.checked)}
            className="w-4 h-4 rounded"
          />
          <span className="text-sm text-gray-600 dark:text-gray-400">Active</span>
        </label>
      </div>
      <div className="flex gap-2">
        <button
          onClick={onEdit}
          className="flex-1 px-3 py-2 bg-cyan-500 hover:bg-cyan-600 text-white rounded text-sm transition-colors"
        >
          Edit
        </button>
        <button
          onClick={onDelete}
          className="flex-1 px-3 py-2 bg-red-500 hover:bg-red-600 text-white rounded text-sm transition-colors"
        >
          Delete
        </button>
      </div>
    </div>
  );
};
