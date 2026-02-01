import React, { Suspense, useEffect } from 'react';
import { useUIStore } from './stores/ui.store';
import { useSearchStore } from './stores/search.store';
import { ModeToggle } from './components/custom/ModeToggle';
import { SearchBar } from './components/SearchBar';
import { SearchResultsGrid } from './components/SearchResults';
import { DownloadMonitor } from './components/DownloadMonitor';
import { Button } from './components/ui/Button';
import { Card, CardBody, CardFooter, CardHeader, CardTitle } from './components/ui/Card';
import { Input } from './components/ui/Input';
import { Badge } from './components/ui/Badge';
import { Search, Download, Settings, Zap } from 'lucide-react';

const ErrorBoundary: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [hasError, setHasError] = React.useState(false);
  const [errorMessage, setErrorMessage] = React.useState('');

  React.useEffect(() => {
    const errorHandler = (event: ErrorEvent) => {
      setHasError(true);
      setErrorMessage(event.message);
    };
    window.addEventListener('error', errorHandler);
    return () => window.removeEventListener('error', errorHandler);
  }, []);

  if (hasError) {
    return (
      <div
        className="flex"
        style={{
          minHeight: '100vh',
          alignItems: 'center',
          justifyContent: 'center',
          padding: '24px',
          background: 'var(--sakura-bg-base)',
        }}
      >
        <Card style={{ maxWidth: '520px', width: '100%' }}>
          <CardHeader>
            <CardTitle level="h2">Something went wrong</CardTitle>
          </CardHeader>
          <CardBody>
            <p style={{ color: 'var(--sakura-error-text)' }}>{errorMessage}</p>
          </CardBody>
          <CardFooter>
            <Button variant="danger" onClick={() => window.location.reload()}>
              Reload Page
            </Button>
          </CardFooter>
        </Card>
      </div>
    );
  }

  return <>{children}</>;
};

export default function App() {
  const { activeView, setActiveView } = useUIStore();
  const { performSearch, results } = useSearchStore();

  // Load initial data on mount
  useEffect(() => {
    if (results.length === 0) {
      performSearch('anime'); // Load default results
    }
  }, []);

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
        {/* Header */}
        <header
          style={{
            position: 'sticky',
            top: 0,
            zIndex: 40,
            borderBottom: '1px solid var(--sakura-border-default)',
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
            <div>
              <h1 style={{ fontSize: 'var(--text-display)', fontWeight: 700 }}>
                🎌 Anime-Sama Downloader
              </h1>
              <div style={{ display: 'flex', alignItems: 'center', gap: '10px', marginTop: '8px' }}>
                <span className="sakura-stamp">Sakura Night</span>
                <span className="sakura-divider" style={{ width: '80px' }}></span>
                <span style={{ color: 'var(--sakura-text-secondary)', fontSize: '12px' }}>MVP v1.0</span>
              </div>
            </div>
            <ModeToggle />
          </div>
        </header>

        {/* Navigation Tabs */}
        <nav
          style={{
            position: 'sticky',
            top: '86px',
            zIndex: 30,
            borderBottom: '1px solid var(--sakura-border-default)',
            background: 'rgba(10,14,26,0.85)',
            backdropFilter: 'blur(12px)',
          }}
        >
          <div
            className="flex"
            style={{
              maxWidth: '1200px',
              margin: '0 auto',
              padding: '12px 24px',
              gap: '8px',
              flexWrap: 'wrap',
            }}
          >
            {(['search', 'downloads', 'rules', 'settings'] as const).map((view) => (
              <Button
                key={view}
                variant={activeView === view ? 'primary' : 'ghost'}
                onClick={() => setActiveView(view)}
                style={{
                  padding: '10px 16px',
                  borderBottom: activeView === view ? '2px solid var(--sakura-accent-cyan-500)' : '2px solid transparent',
                }}
              >
                <span style={{ marginRight: '8px' }}>
                  {view === 'search' && <Search className="w-5 h-5" />}
                  {view === 'downloads' && <Download className="w-5 h-5" />}
                  {view === 'rules' && <Settings className="w-5 h-5" />}
                  {view === 'settings' && <Zap className="w-5 h-5" />}
                </span>
                {view.charAt(0).toUpperCase() + view.slice(1)}
              </Button>
            ))}
          </div>
        </nav>

        {/* Main Content */}
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
            {activeView === 'search' && (
              <div style={{ display: 'flex', flexDirection: 'column', gap: '32px' }}>
                <Card>
                  <CardHeader>
                    <CardTitle level="h2">Rechercher & Télécharger</CardTitle>
                  </CardHeader>
                  <CardBody>
                    <p style={{ color: 'var(--sakura-text-secondary)', marginBottom: '16px' }}>
                      Trouvez vos animes préférés et téléchargez-les facilement
                    </p>
                    <SearchBar />
                  </CardBody>
                </Card>
                <SearchResultsGrid />
              </div>
            )}

            {activeView === 'downloads' && (
              <Card>
                <CardHeader>
                  <CardTitle level="h2">Download Monitor</CardTitle>
                </CardHeader>
                <CardBody>
                  <DownloadMonitor />
                </CardBody>
              </Card>
            )}

            {activeView === 'rules' && (
              <Card>
                <CardHeader>
                  <CardTitle level="h2">Automation Rules</CardTitle>
                </CardHeader>
                <CardBody>
                  <p style={{ color: 'var(--sakura-text-secondary)', marginBottom: '16px' }}>
                    📋 Create rules to automatically download new episodes matching your criteria.
                  </p>
                  <Card>
                    <CardBody>
                      <div className="flex" style={{ gap: '12px', alignItems: 'center' }}>
                        <Badge variant="info">Example</Badge>
                        <div>
                          <p style={{ fontWeight: 600, marginBottom: '6px' }}>
                            Auto-download new episodes
                          </p>
                          <p style={{ color: 'var(--sakura-text-secondary)', fontSize: '12px' }}>
                            • Pattern: "Attack on Titan*" • Limit: 2 concurrent
                          </p>
                        </div>
                      </div>
                    </CardBody>
                  </Card>
                </CardBody>
                <CardFooter>
                  <Button variant="primary">+ Add New Rule</Button>
                </CardFooter>
              </Card>
            )}

            {activeView === 'settings' && (
              <Card>
                <CardHeader>
                  <CardTitle level="h2">⚙️ Settings</CardTitle>
                </CardHeader>
                <CardBody>
                  <div className="grid" style={{ gridTemplateColumns: 'repeat(auto-fit, minmax(240px, 1fr))', gap: '16px' }}>
                    <Card>
                      <CardBody>
                        <Input
                          label="Download Path"
                          placeholder="/media/anime"
                          defaultValue="/media/anime"
                        />
                      </CardBody>
                    </Card>

                    <Card>
                      <CardBody>
                        <Input
                          label="Concurrent Downloads"
                          type="number"
                          min="1"
                          max="5"
                          defaultValue="2"
                        />
                      </CardBody>
                    </Card>

                    <Card style={{ gridColumn: '1 / -1' }}>
                      <CardBody>
                        <label className="flex" style={{ alignItems: 'center', gap: '12px' }}>
                          <input type="checkbox" defaultChecked />
                          <span>Enable browser notifications</span>
                        </label>
                      </CardBody>
                    </Card>

                    <Card style={{ gridColumn: '1 / -1' }}>
                      <CardBody>
                        <label className="flex" style={{ alignItems: 'center', gap: '12px' }}>
                          <input type="checkbox" />
                          <span>Auto-delete after download</span>
                        </label>
                      </CardBody>
                    </Card>
                  </div>
                </CardBody>
                <CardFooter>
                  <Button variant="secondary">💾 Save Settings</Button>
                </CardFooter>
              </Card>
            )}
          </Suspense>
        </main>

        {/* Footer */}
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
      </div>
    </ErrorBoundary>
  );
}
