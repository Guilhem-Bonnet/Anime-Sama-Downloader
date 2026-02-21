import React from 'react';
import { ModeToggle } from '../custom/ModeToggle';

export function Header() {
  return (
    <header
      style={{
        position: 'sticky',
        top: 0,
        zIndex: 40,
        borderBottom: '1px solid var(--night-border-default)',
        background: 'linear-gradient(180deg, rgba(10,14,26,0.96), rgba(26,31,46,0.92))',
        backdropFilter: 'blur(12px)',
      }}
    >
      <div
        className="flex"
        style={{
          maxWidth: '1200px',
          margin: '0 auto',
          padding: '20px 24px',
          alignItems: 'center',
          justifyContent: 'space-between',
        }}
      >
        <div className="title-easter">
          <h1
            style={{
              fontSize: 'var(--text-display)',
              fontWeight: 700,
              fontFamily: 'var(--font-display)',
              letterSpacing: '0.04em',
            }}
          >
            🎌 Anime-Sama Downloader
          </h1>
          <div className="easter-ink">墨の道</div>
          <div style={{ display: 'flex', alignItems: 'center', gap: '10px', marginTop: '8px' }}>
            <span className="night-stamp">Nocturne Core</span>
            <span className="night-divider" style={{ width: '80px' }}></span>
            <span style={{ color: 'var(--night-text-secondary)', fontSize: '12px' }}>MVP v1.0</span>
          </div>
        </div>
        <ModeToggle />
      </div>
    </header>
  );
}
