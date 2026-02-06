import React from 'react';

export function Footer() {
  return (
    <footer
      style={{
        borderTop: '1px solid var(--sakura-border-default)',
        marginTop: '48px',
        background: 'var(--sakura-bg-surface)',
      }}
    >
      <div
        style={{
          maxWidth: '1200px',
          margin: '0 auto',
          padding: '24px',
          textAlign: 'center',
          fontSize: '12px',
          color: 'var(--sakura-text-secondary)',
        }}
      >
        <p>Anime-Sama Downloader MVP • Built with ❤️ using Go + React</p>
      </div>
    </footer>
  );
}
