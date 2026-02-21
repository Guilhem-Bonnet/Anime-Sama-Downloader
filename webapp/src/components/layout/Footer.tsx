import React from 'react';

export function Footer() {
  return (
    <footer
      style={{
        borderTop: '1px solid var(--night-border-default)',
        marginTop: '48px',
        background: 'var(--night-bg-surface)',
      }}
    >
      <div
        style={{
          maxWidth: '1200px',
          margin: '0 auto',
          padding: '24px',
          textAlign: 'center',
          fontSize: '12px',
          color: 'var(--night-text-secondary)',
        }}
      >
        <p>Anime-Sama Downloader MVP • Built with ❤️ using Go + React</p>
      </div>
    </footer>
  );
}
