import React, { Suspense } from 'react';
import { Outlet } from 'react-router-dom';
import { Header } from './Header';
import { TabNav } from './TabNav';
import { Footer } from './Footer';
import { ErrorBoundary } from './ErrorBoundary';

export function Layout() {
  return (
    <ErrorBoundary>
      <div
        className="sakura-ink-bg"
        style={{
          minHeight: '100vh',
          background: 'var(--sakura-bg-base)',
          color: 'var(--sakura-text-primary)',
        }}
      >
        <div className="ink-silhouette" aria-hidden="true" />
        <Header />
        <TabNav />
        <main style={{ maxWidth: '1200px', margin: '0 auto', padding: '32px 24px' }}>
          <Suspense
            fallback={
              <div className="flex" style={{ alignItems: 'center', justifyContent: 'center', padding: '48px 0' }}>
                <div className="spin" style={{ fontSize: '20px', marginRight: '12px' }}>
                  ⌛
                </div>
                <p style={{ color: 'var(--sakura-text-secondary)' }}>Loading...</p>
              </div>
            }
          >
            <Outlet />
          </Suspense>
        </main>
        <Footer />
      </div>
    </ErrorBoundary>
  );
}
