import React, { useEffect, useRef } from 'react';

export interface LogViewerProps {
  logs: string[];
  maxLines?: number;
  autoScroll?: boolean;
  className?: string;
}

export const LogViewer: React.FC<LogViewerProps> = ({
  logs,
  maxLines = 100,
  autoScroll = true,
  className = '',
}) => {
  const scrollRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (autoScroll && scrollRef.current) {
      scrollRef.current.scrollTop = scrollRef.current.scrollHeight;
    }
  }, [logs, autoScroll]);

  const displayedLogs = logs.slice(-maxLines);

  return (
    <div
      ref={scrollRef}
      className={`p-4 bg-black dark:bg-gray-900 rounded-lg border border-gray-700 dark:border-gray-800 overflow-y-auto font-mono text-sm h-64 ${className}`}
    >
      {displayedLogs.length === 0 ? (
        <p className="text-gray-500">No logs yet...</p>
      ) : (
        <div className="space-y-1">
          {displayedLogs.map((log, index) => (
            <div
              key={index}
              className={`text-xs ${
                log.includes('ERROR')
                  ? 'text-red-400'
                  : log.includes('WARN')
                    ? 'text-yellow-400'
                    : log.includes('SUCCESS')
                      ? 'text-green-400'
                      : 'text-gray-300'
              }`}
            >
              <span className="text-gray-600">[{new Date().toLocaleTimeString()}]</span> {log}
            </div>
          ))}
        </div>
      )}
    </div>
  );
};
