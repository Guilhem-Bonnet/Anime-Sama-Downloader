import React, { createContext, useContext, useState, useCallback } from 'react';
import { Badge } from './Badge';

export interface Toast {
  id: string;
  variant: 'success' | 'warning' | 'error' | 'info';
  title: string;
  message?: string;
  duration?: number;
}

interface ToastContextType {
  toasts: Toast[];
  addToast: (toast: Omit<Toast, 'id'>) => void;
  removeToast: (id: string) => void;
}

const ToastContext = createContext<ToastContextType | undefined>(undefined);

export function ToastProvider({ children }: { children: React.ReactNode }) {
  const [toasts, setToasts] = useState<Toast[]>([]);

  const addToast = useCallback((toast: Omit<Toast, 'id'>) => {
    const id = Math.random().toString(36).substring(2, 9);
    const newToast: Toast = { ...toast, id };
    setToasts((prev) => [...prev, newToast]);

    const duration = toast.duration || 5000;
    setTimeout(() => {
      setToasts((prev) => prev.filter((t) => t.id !== id));
    }, duration);
  }, []);

  const removeToast = useCallback((id: string) => {
    setToasts((prev) => prev.filter((t) => t.id !== id));
  }, []);

  return (
    <ToastContext.Provider value={{ toasts, addToast, removeToast }}>
      {children}
      <ToastContainer toasts={toasts} onRemove={removeToast} />
    </ToastContext.Provider>
  );
}

export function useToast() {
  const context = useContext(ToastContext);
  if (!context) {
    throw new Error('useToast must be used within ToastProvider');
  }
  return context;
}

function ToastContainer({ toasts, onRemove }: { toasts: Toast[]; onRemove: (id: string) => void }) {
  if (toasts.length === 0) return null;

  return (
    <div
      style={{
        position: 'fixed',
        bottom: 'var(--space-4)',
        right: 'var(--space-4)',
        zIndex: 'var(--z-toast)',
        display: 'flex',
        flexDirection: 'column',
        gap: 'var(--space-2)',
        maxWidth: '400px',
      }}
    >
      {toasts.map((toast) => (
        <ToastItem key={toast.id} toast={toast} onRemove={onRemove} />
      ))}
    </div>
  );
}

function ToastItem({ toast, onRemove }: { toast: Toast; onRemove: (id: string) => void }) {
  return (
    <div
      className="frame-ornate"
      style={{
        background: 'var(--night-bg-surface)',
        border: '1px solid var(--night-border-default)',
        borderRadius: 'var(--radius-lg)',
        padding: 'var(--space-4)',
        boxShadow: 'var(--shadow-lg)',
        animation: 'toast-slide-in 300ms ease-out',
      }}
    >
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', gap: 'var(--space-3)' }}>
        <div style={{ flex: 1 }}>
          <div style={{ display: 'flex', alignItems: 'center', gap: 'var(--space-2)', marginBottom: 'var(--space-1)' }}>
            <Badge variant={toast.variant}>{toast.variant.toUpperCase()}</Badge>
            <h4 style={{ fontSize: 'var(--text-body)', fontWeight: 600, color: 'var(--night-text-primary)' }}>
              {toast.title}
            </h4>
          </div>
          {toast.message && (
            <p style={{ fontSize: 'var(--text-body-sm)', color: 'var(--night-text-secondary)', marginTop: 'var(--space-1)' }}>
              {toast.message}
            </p>
          )}
        </div>
        <button
          onClick={() => onRemove(toast.id)}
          style={{
            background: 'transparent',
            border: 'none',
            color: 'var(--night-text-secondary)',
            cursor: 'pointer',
            fontSize: '20px',
            padding: 0,
            lineHeight: 1,
          }}
          aria-label="Close toast"
        >
          ×
        </button>
      </div>
    </div>
  );
}
