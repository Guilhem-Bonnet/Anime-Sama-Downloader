import React, { useEffect } from 'react';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { Layout } from './components/layout';
import { SearchPage, DownloadsPage, RulesPage, SettingsPage, AnimeDetailPage } from './pages';
import { useSearchStore } from './stores/search.store';

function InitialDataLoader({ children }: { children: React.ReactNode }) {
  const { performSearch, results } = useSearchStore();

  useEffect(() => {
    if (results.length === 0) {
      performSearch('anime'); // Load default results
    }
  }, []);

  return <>{children}</>;
}

export function AppRouter() {
  return (
    <BrowserRouter>
      <InitialDataLoader>
        <Routes>
          <Route path="/" element={<Layout />}>
            <Route index element={<SearchPage />} />
            <Route path="downloads" element={<DownloadsPage />} />
            <Route path="rules" element={<RulesPage />} />
            <Route path="settings" element={<SettingsPage />} />
          </Route>
          <Route path="/anime/:id" element={<AnimeDetailPage />} />
          <Route path="*" element={<Navigate to="/" replace />} />
        </Routes>
      </InitialDataLoader>
    </BrowserRouter>
  );
}
