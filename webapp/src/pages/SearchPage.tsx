import React from 'react';
import { SearchBar } from '../components/SearchBar';
import { SearchResultsGrid } from '../components/SearchResults';
import { Card, CardBody, CardHeader, CardTitle } from '../components/ui/Card';
import { HeroLandscapeIllustration } from '../components/illustrations/SakuraIllustrations';

export function SearchPage() {
  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: '32px' }}>
      <Card className="frame-ornate">
        <CardHeader>
          <CardTitle level="h2">Rechercher & Télécharger</CardTitle>
          <div className="kanji-brush" style={{ marginTop: '6px' }}>
            検索
          </div>
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
  );
}
