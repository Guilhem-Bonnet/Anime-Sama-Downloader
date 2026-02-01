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
import { HeroLandscapeIllustration, EmptyDownloadsIllustration } from './components/illustrations/SakuraIllustrations';

class ErrorBoundary extends React.Component<
  { children: React.ReactNode },
  { hasError: boolean; errorMessage: string }
> {
  constructor(props: { children: React.ReactNode }) {
    super(props);
    this.state = { hasError: false, errorMessage: '' };
  }

  static getDerivedStateFromError(error: Error) {
    return { hasError: true, errorMessage: error.message };
  }

  componentDidCatch(error: Error) {
    console.error('App render error:', error);
  }

  render() {
    if (this.state.hasError) {
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
              <p style={{ color: 'var(--sakura-error-text)' }}>{this.state.errorMessage}</p>
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

    return <>{this.props.children}</>;
  }
}

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
        <div className="ink-silhouette" aria-hidden="true" />
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
                  borderBottom: activeView === view ? `2px solid var(--sakura-accent-brown-500)` : '2px solid transparent',
                  color: activeView === view ? 'var(--sakura-accent-brown-500)' : 'var(--sakura-text-secondary)',
                  transition: 'all 200ms ease',
                  fontWeight: activeView === view ? 600 : 500,
                }}
              >
                <span style={{ marginRight: '8px' }}>
                  {view === 'search' && <Search className="w-5 h-5" />}
                  {view === 'downloads' && <Download className="w-5 h-5" />}
                  {view === 'rules' && <Settings className="w-5 h-5" />}
                  {view === 'settings' && <Zap className="w-5 h-5" />}
                </span>
                {view === 'search' && 'Rechercher'}
                {view === 'downloads' && 'Téléchargements'}
                {view === 'rules' && 'Règles'}
                {view === 'settings' && 'Paramètres'}
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
                <Card className="frame-ornate">
                  <CardHeader>
                    <CardTitle level="h2">Rechercher & Télécharger</CardTitle>
                    <div className="kanji-brush" style={{ marginTop: '6px' }}>検索</div>
                  </CardHeader>
                  <CardBody>
                    <p style={{ color: 'var(--sakura-text-secondary)', marginBottom: '16px' }}>
                      Trouvez vos animes préférés et téléchargez-les facilement
                    </p>
                    <div style={{ marginBottom: '18px', maxHeight: '160px', overflow: 'hidden', opacity: 0.9 }}>
                      <HeroLandscapeIllustration />
                    </div>
                    <SearchBar />
                  </CardBody>
                </Card>
                <SearchResultsGrid />
              </div>
            )}

            {activeView === 'downloads' && (
              <Card className="frame-ornate">
                <CardHeader>
                  <CardTitle level="h2">Téléchargements</CardTitle>
                  <div className="kanji-brush" style={{ marginTop: '6px' }}>保存</div>
                </CardHeader>
                <CardBody>
                  <div style={{ marginBottom: '24px', maxHeight: '180px', overflow: 'hidden', opacity: 0.8 }}>
                    <EmptyDownloadsIllustration />
                  </div>
                  <DownloadMonitor />
                </CardBody>
              </Card>
            )}

            {activeView === 'rules' && (
              <div style={{ display: 'flex', flexDirection: 'column', gap: '24px' }}>
                <Card className="frame-ornate">
                  <CardHeader>
                    <CardTitle level="h2">Règles d'Automatisation</CardTitle>
                    <div className="kanji-brush" style={{ marginTop: '6px' }}>規則</div>
                  </CardHeader>
                  <CardBody>
                    <p style={{ color: 'var(--sakura-text-secondary)', marginBottom: '20px' }}>
                      Créez des règles pour télécharger automatiquement les nouveaux épisodes selon vos critères.
                    </p>
                  </CardBody>
                </Card>

                <Card className="frame-ornate">
                  <CardBody>
                    <div className="flex" style={{ gap: '16px', alignItems: 'flex-start' }}>
                      <Badge variant="info">Exemple</Badge>
                      <div style={{ flex: 1 }}>
                        <p style={{ fontWeight: 600, marginBottom: '8px', color: 'var(--sakura-text-primary)' }}>
                          Téléchargement automatique des nouveaux épisodes
                        </p>
                        <p style={{ color: 'var(--sakura-text-secondary)', fontSize: '13px', lineHeight: '1.6' }}>
                          • Motif: "Attack on Titan*" <br/>
                          • Limite: 2 téléchargements simultanés <br/>
                          • Dossier: /media/anime
                        </p>
                      </div>
                    </div>
                  </CardBody>
                </Card>

                <Card className="frame-ornate">
                  <CardBody>
                    <div className="flex" style={{ gap: '16px', alignItems: 'flex-start' }}>
                      <Badge variant="success">Actif</Badge>
                      <div style={{ flex: 1 }}>
                        <p style={{ fontWeight: 600, marginBottom: '8px', color: 'var(--sakura-text-primary)' }}>
                          Jujutsu Kaisen - Nouvelles sorties
                        </p>
                        <p style={{ color: 'var(--sakura-text-secondary)', fontSize: '13px', lineHeight: '1.6' }}>
                          • Dernière vérification: il y a 2 heures <br/>
                          • Épisodes téléchargés: 24 <br/>
                          • Prochaine vérification: dans 1 heure
                        </p>
                      </div>
                    </div>
                  </CardBody>
                </Card>

                <div>
                  <Button variant="primary">+ Ajouter une nouvelle règle</Button>
                </div>
              </div>
            )}

            {activeView === 'settings' && (
              <div style={{ display: 'flex', flexDirection: 'column', gap: '24px' }}>
                <Card className="frame-ornate">
                  <CardHeader>
                    <CardTitle level="h2">Paramètres</CardTitle>
                    <div className="kanji-brush" style={{ marginTop: '6px' }}>設定</div>
                  </CardHeader>
                  <CardBody>
                    <p style={{ color: 'var(--sakura-text-secondary)' }}>
                      Configurez vos préférences de téléchargement et les notifications.
                    </p>
                  </CardBody>
                </Card>

                <div className="grid" style={{ gridTemplateColumns: 'repeat(auto-fit, minmax(240px, 1fr))', gap: '16px' }}>
                  <Card className="frame-ornate">
                    <CardBody>
                      <label style={{ display: 'flex', flexDirection: 'column', gap: '8px' }}>
                        <span style={{ fontSize: '13px', fontWeight: 600, color: 'var(--sakura-text-secondary)' }}>
                          Dossier de téléchargement
                        </span>
                        <Input
                          placeholder="/media/anime"
                          defaultValue="/media/anime"
                        />
                      </label>
                    </CardBody>
                  </Card>

                  <Card className="frame-ornate">
                    <CardBody>
                      <label style={{ display: 'flex', flexDirection: 'column', gap: '8px' }}>
                        <span style={{ fontSize: '13px', fontWeight: 600, color: 'var(--sakura-text-secondary)' }}>
                          Téléchargements simultanés
                        </span>
                        <Input
                          type="number"
                          min="1"
                          max="5"
                          defaultValue="2"
                        />
                      </label>
                    </CardBody>
                  </Card>

                  <Card className="frame-ornate">
                    <CardBody>
                      <label style={{ display: 'flex', flexDirection: 'column', gap: '8px' }}>
                        <span style={{ fontSize: '13px', fontWeight: 600, color: 'var(--sakura-text-secondary)' }}>
                          Qualité préférée
                        </span>
                        <select style={{ 
                          padding: '8px 12px', 
                          borderRadius: '12px',
                          border: '1px solid var(--sakura-border-default)',
                          background: 'var(--sakura-bg-surface)',
                          color: 'var(--sakura-text-primary)',
                          fontSize: '13px',
                          cursor: 'pointer'
                        }}>
                          <option>720p</option>
                          <option selected>1080p</option>
                          <option>480p</option>
                        </select>
                      </label>
                    </CardBody>
                  </Card>
                </div>

                <div style={{ display: 'flex', flexDirection: 'column', gap: '12px' }}>
                  <Card className="frame-ornate">
                    <CardBody>
                      <label className="flex" style={{ alignItems: 'center', gap: '12px', cursor: 'pointer' }}>
                        <input type="checkbox" defaultChecked style={{ width: '18px', height: '18px', cursor: 'pointer' }} />
                        <span style={{ color: 'var(--sakura-text-primary)' }}>
                          Activer les notifications du navigateur
                        </span>
                      </label>
                    </CardBody>
                  </Card>

                  <Card className="frame-ornate">
                    <CardBody>
                      <label className="flex" style={{ alignItems: 'center', gap: '12px', cursor: 'pointer' }}>
                        <input type="checkbox" style={{ width: '18px', height: '18px', cursor: 'pointer' }} />
                        <span style={{ color: 'var(--sakura-text-primary)' }}>
                          Supprimer automatiquement après téléchargement
                        </span>
                      </label>
                    </CardBody>
                  </Card>

                  <Card className="frame-ornate">
                    <CardBody>
                      <label className="flex" style={{ alignItems: 'center', gap: '12px', cursor: 'pointer' }}>
                        <input type="checkbox" defaultChecked style={{ width: '18px', height: '18px', cursor: 'pointer' }} />
                        <span style={{ color: 'var(--sakura-text-primary)' }}>
                          Vérifier les mises à jour automatiquement
                        </span>
                      </label>
                    </CardBody>
                  </Card>
                </div>

                <Card className="frame-ornate" style={{ background: 'linear-gradient(135deg, rgba(143,106,61,0.05), rgba(125,114,103,0.05))' }}>
                  <CardBody>
                    <div style={{ display: 'flex', flexDirection: 'column', gap: '12px' }}>
                      <p style={{ color: 'var(--sakura-text-secondary)', fontSize: '13px' }}>
                        <strong>Version</strong> · v1.0.0
                      </p>
                      <p style={{ color: 'var(--sakura-text-secondary)', fontSize: '13px' }}>
                        <strong>Dernière mise à jour</strong> · 1er février 2026
                      </p>
                      <p style={{ color: 'var(--sakura-text-secondary)', fontSize: '13px' }}>
                        <strong>Statut</strong> · ✨ Stable
                      </p>
                    </div>
                  </CardBody>
                </Card>

                <div style={{ display: 'flex', gap: '12px', justifyContent: 'flex-end' }}>
                  <Button variant="secondary">Réinitialiser les paramètres</Button>
                  <Button variant="primary">💾 Enregistrer les paramètres</Button>
                </div>
              </div>
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
