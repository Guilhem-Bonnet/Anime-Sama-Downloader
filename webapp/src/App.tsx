import React, { Suspense, useEffect } from 'react';
import { useUIStore } from './stores/ui.store';
import { useSearchStore } from './stores/search.store';
import { ModeToggle } from './components/custom/ModeToggle';
import { SearchBar } from './components/SearchBar';
import { SearchResultsGrid } from './components/SearchResults';
import { DownloadMonitor } from './components/DownloadMonitor';

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
      <div className="min-h-screen bg-gray-100 dark:bg-gray-900 flex items-center justify-center p-4">
        <div className="bg-red-50 dark:bg-red-900 p-8 rounded-lg max-w-md">
          <h1 className="text-2xl font-bold text-red-900 dark:text-red-100">Something went wrong</h1>
          <p className="text-red-700 dark:text-red-200 mt-2 text-sm">{errorMessage}</p>
          <button
            onClick={() => window.location.reload()}
            className="mt-4 px-4 py-2 bg-red-600 text-white rounded hover:bg-red-700 transition-colors"
          >
            Reload Page
          </button>
        </div>
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
      <div className="min-h-screen bg-gradient-to-br from-gray-50 via-pink-50/20 to-cyan-50/20 dark:from-gray-950 dark:via-gray-900 dark:to-gray-950 transition-colors">
        {/* Header avec gradient Sakura Night */}
        <header className="border-b border-gray-200/50 dark:border-gray-700/50 bg-gradient-to-r from-white via-pink-50/30 to-cyan-50/30 dark:from-gray-900 dark:via-gray-900 dark:to-gray-900 backdrop-blur-xl sticky top-0 z-40 shadow-lg">
          <div className="max-w-7xl mx-auto px-6 py-5 flex justify-between items-center">
            <div>
              <h1 className="text-3xl font-bold bg-gradient-to-r from-magenta-600 via-pink-600 to-cyan-600 bg-clip-text text-transparent">
                🎌 Anime-Sama Downloader
              </h1>
              <p className="text-xs text-gray-600 dark:text-gray-400 mt-1 font-medium">Sakura Night Design • MVP v1.0</p>
            </div>
            <ModeToggle />
          </div>
        </header>

        {/* Navigation Tabs avec effet glassmorphism */}
        <nav className="bg-white/80 dark:bg-gray-800/80 backdrop-blur-xl border-b border-gray-200/50 dark:border-gray-700/50 sticky top-[88px] z-30 shadow-sm">
          <div className="max-w-7xl mx-auto px-6 flex gap-2">
            {(['search', 'downloads', 'rules', 'settings'] as const).map((view) => (
              <button
                key={view}
                onClick={() => setActiveView(view)}
                className={`relative px-6 py-4 font-semibold border-b-3 transition-all duration-300 ${
                  activeView === view
                    ? 'border-cyan-500 text-transparent bg-clip-text bg-gradient-to-r from-magenta-600 to-cyan-600'
                    : 'border-transparent text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white hover:bg-gray-100/50 dark:hover:bg-gray-700/50'
                }`}
              >
                {activeView === view && (
                  <span className="absolute inset-x-0 -bottom-[2px] h-[3px] bg-gradient-to-r from-magenta-500 to-cyan-500 rounded-t-full" />
                )}
                <span className="text-xl mr-2">
                  {view === 'search' && '🔍'}
                  {view === 'downloads' && '📥'}
                  {view === 'rules' && '⚙️'}
                  {view === 'settings' && '⚡'}
                </span>
                {view.charAt(0).toUpperCase() + view.slice(1)}
              </button>
            ))}
          </div>
        </nav>

        {/* Main Content */}
        <main className="max-w-7xl mx-auto px-6 py-12">
          <Suspense fallback={
            <div className="flex items-center justify-center py-12">
              <div className="animate-spin text-2xl mr-3">⌛</div>
              <p className="text-gray-500 dark:text-gray-400">Loading...</p>
            </div>
          }>
            {activeView === 'search' && (
              <div className="space-y-10 animate-fadeIn">
                <div className="text-center">
                  <h2 className="text-4xl font-bold mb-3 bg-gradient-to-r from-magenta-600 via-pink-600 to-cyan-600 bg-clip-text text-transparent">
                    Rechercher & Télécharger
                  </h2>
                  <p className="text-gray-600 dark:text-gray-400 text-lg">
                    Trouvez vos animes préférés et téléchargez-les facilement
                  </p>
                </div>
                <SearchBar />
                <SearchResultsGrid />
              </div>
            )}

            {activeView === 'downloads' && (
              <div>
                <h2 className="text-2xl font-bold mb-6 text-gray-900 dark:text-white">
                  Download Monitor
                </h2>
                <DownloadMonitor />
              </div>
            )}

            {activeView === 'rules' && (
              <div className="space-y-6">
                <h2 className="text-2xl font-bold text-gray-900 dark:text-white">
                  Automation Rules
                </h2>
                <div className="p-6 bg-blue-50 dark:bg-blue-900/20 rounded-lg border border-blue-200 dark:border-blue-700">
                  <p className="text-gray-700 dark:text-gray-300 mb-4">
                    📋 Create rules to automatically download new episodes matching your criteria.
                  </p>
                  <div className="space-y-3">
                    <div className="p-4 bg-white dark:bg-gray-800 rounded border border-gray-200 dark:border-gray-700">
                      <p className="text-sm font-medium text-gray-900 dark:text-white mb-2">
                        Example: Auto-download new episodes
                      </p>
                      <p className="text-xs text-gray-600 dark:text-gray-400">
                        • Pattern: "Attack on Titan*" • Limit: 2 concurrent
                      </p>
                    </div>
                  </div>
                </div>
                <button className="px-6 py-3 bg-magenta-600 hover:bg-magenta-700 text-white rounded-lg transition-colors font-medium">
                  + Add New Rule
                </button>
              </div>
            )}

            {activeView === 'settings' && (
              <div className="space-y-6">
                <h2 className="text-2xl font-bold text-gray-900 dark:text-white">
                  ⚙️ Settings
                </h2>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                  <div className="p-4 bg-gray-100 dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700">
                    <label className="block text-sm font-medium text-gray-900 dark:text-white mb-2">
                      Download Path
                    </label>
                    <input
                      type="text"
                      placeholder="/media/anime"
                      className="w-full px-4 py-2 bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded"
                      defaultValue="/media/anime"
                    />
                  </div>

                  <div className="p-4 bg-gray-100 dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700">
                    <label className="block text-sm font-medium text-gray-900 dark:text-white mb-2">
                      Concurrent Downloads
                    </label>
                    <input
                      type="number"
                      min="1"
                      max="5"
                      className="w-full px-4 py-2 bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded"
                      defaultValue="2"
                    />
                  </div>

                  <div className="p-4 bg-gray-100 dark:bg-gray-800 rounded-lg col-span-1 md:col-span-2">
                    <label className="flex items-center gap-3 cursor-pointer">
                      <input
                        type="checkbox"
                        className="w-4 h-4 rounded"
                        defaultChecked
                      />
                      <span className="font-medium text-gray-900 dark:text-white">
                        Enable browser notifications
                      </span>
                    </label>
                  </div>

                  <div className="p-4 bg-gray-100 dark:bg-gray-800 rounded-lg col-span-1 md:col-span-2">
                    <label className="flex items-center gap-3 cursor-pointer">
                      <input
                        type="checkbox"
                        className="w-4 h-4 rounded"
                      />
                      <span className="font-medium text-gray-900 dark:text-white">
                        Auto-delete after download
                      </span>
                    </label>
                  </div>
                </div>

                <button className="px-6 py-3 bg-cyan-600 hover:bg-cyan-700 text-white rounded-lg transition-colors font-medium">
                  💾 Save Settings
                </button>
              </div>
            )}
          </Suspense>
        </main>

        {/* Footer */}
        <footer className="border-t border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-800 mt-12">
          <div className="max-w-7xl mx-auto px-4 py-8 text-center text-sm text-gray-600 dark:text-gray-400">
            <p>Anime-Sama Downloader MVP • Built with ❤️ using Go + React</p>
          </div>
        </footer>
      </div>
    </ErrorBoundary>
  );
}
