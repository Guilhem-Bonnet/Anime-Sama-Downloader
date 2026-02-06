import React, { useEffect, useRef } from 'react';

export interface ModalProps {
  isOpen: boolean;
  onClose: () => void;
  title?: string;
  children: React.ReactNode;
  size?: 'sm' | 'md' | 'lg' | 'xl';
  showCloseButton?: boolean;
}

export function Modal({
  isOpen,
  onClose,
  title,
  children,
  size = 'md',
  showCloseButton = true,
}: ModalProps) {
  const modalRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (isOpen) {
      document.body.style.overflow = 'hidden';
    } else {
      document.body.style.overflow = '';
    }
    return () => {
      document.body.style.overflow = '';
    };
  }, [isOpen]);

  useEffect(() => {
    const handleEscape = (e: KeyboardEvent) => {
      if (e.key === 'Escape' && isOpen) {
        onClose();
      }
    };
    document.addEventListener('keydown', handleEscape);
    return () => document.removeEventListener('keydown', handleEscape);
  }, [isOpen, onClose]);

  if (!isOpen) return null;

  const sizeMap = {
    sm: '400px',
    md: '560px',
    lg: '720px',
    xl: '960px',
  };

  return (
    <div
      className="modal-overlay"
      style={{
        position: 'fixed',
        top: 0,
        left: 0,
        right: 0,
        bottom: 0,
        background: 'rgba(0, 0, 0, 0.6)',
        backdropFilter: 'blur(4px)',
        zIndex: 'var(--z-modal)',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        padding: 'var(--space-4)',
      }}
      onClick={(e) => {
        if (e.target === e.currentTarget) onClose();
      }}
    >
      <div
        ref={modalRef}
        className="modal-content frame-ornate"
        style={{
          background: 'var(--sakura-bg-surface)',
          borderRadius: 'var(--radius-lg)',
          border: '1px solid var(--sakura-border-default)',
          maxWidth: sizeMap[size],
          width: '100%',
          maxHeight: 'calc(100vh - var(--space-8))',
          overflow: 'auto',
          boxShadow: 'var(--shadow-xl)',
        }}
      >
        {(title || showCloseButton) && (
          <div
            className="modal-header"
            style={{
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'space-between',
              padding: 'var(--space-5) var(--space-6)',
              borderBottom: '1px solid var(--sakura-border-default)',
            }}
          >
            {title && (
              <h2
                style={{
                  fontSize: 'var(--text-h2)',
                  fontWeight: 600,
                  color: 'var(--sakura-text-primary)',
                }}
              >
                {title}
              </h2>
            )}
            {showCloseButton && (
              <button
                onClick={onClose}
                style={{
                  background: 'transparent',
                  border: 'none',
                  color: 'var(--sakura-text-secondary)',
                  cursor: 'pointer',
                  fontSize: '24px',
                  padding: 'var(--space-2)',
                  lineHeight: 1,
                  transition: 'color var(--transition-fast)',
                }}
                onMouseEnter={(e) => {
                  e.currentTarget.style.color = 'var(--sakura-text-primary)';
                }}
                onMouseLeave={(e) => {
                  e.currentTarget.style.color = 'var(--sakura-text-secondary)';
                }}
                aria-label="Close modal"
              >
                ×
              </button>
            )}
          </div>
        )}
        <div
          className="modal-body"
          style={{
            padding: 'var(--space-6)',
          }}
        >
          {children}
        </div>
      </div>
    </div>
  );
}
