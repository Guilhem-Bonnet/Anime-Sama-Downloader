import React, { Suspense, useEffect } from 'react';
import { Outlet } from 'react-router-dom';
import { Header } from './Header';
import { TabNav } from './TabNav';
import { Footer } from './Footer';
import { ErrorBoundary } from './ErrorBoundary';
import { useJobsStore } from '../../stores/jobs.store';

export function Layout() {
  const loadJobs = useJobsStore((s) => s.loadJobs);

  // Charger les jobs dès le démarrage de l'app (Layout ne démonte jamais).
  // Polling toutes les 15s pour rester à jour même sans SSE.
  useEffect(() => {
    loadJobs();
    const interval = setInterval(loadJobs, 15_000);
    return () => clearInterval(interval);
  }, []);

  return (
    <ErrorBoundary>
      <div
        className="night-ink-bg"
        style={{
          minHeight: '100vh',
          background: 'var(--night-bg-base)',
          color: 'var(--night-text-primary)',
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
                <p style={{ color: 'var(--night-text-secondary)' }}>Loading...</p>
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
