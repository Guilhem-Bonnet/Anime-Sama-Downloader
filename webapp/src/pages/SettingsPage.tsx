import React, { useEffect } from 'react';
import { Button } from '../components/ui/Button';
import { Card, CardBody, CardHeader, CardTitle } from '../components/ui/Card';
import { Input } from '../components/ui/Input';
import { useSettingsStore } from '../stores/settings.store';

export function SettingsPage() {
  const {
    settings,
    isDirty,
    isLoading,
    error,
    updateSettings,
    saveSettings,
    loadSettings,
    resetSettings,
  } = useSettingsStore();

  useEffect(() => {
    loadSettings();
  }, []);

  const handleSave = async () => {
    await saveSettings();
  };

  const handleReset = () => {
    if (confirm('Réinitialiser tous les paramètres à leurs valeurs par défaut ?')) {
      resetSettings();
    }
  };

  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: '24px' }}>
      <Card className="frame-ornate">
        <CardHeader>
          <CardTitle level="h2">Paramètres</CardTitle>
          <div className="kanji-brush" style={{ marginTop: '6px' }}>
            設定
          </div>
        </CardHeader>
        <CardBody>
          <p style={{ color: 'var(--night-text-secondary)' }}>
            Configurez le service de téléchargement et les intégrations media server.
          </p>
          {error && (
            <p style={{ color: '#dc2626', fontSize: '13px', marginTop: '8px' }}>
              ⚠ {error}
            </p>
          )}
        </CardBody>
      </Card>

      {/* Téléchargement */}
      <div className="grid" style={{ gridTemplateColumns: 'repeat(auto-fit, minmax(240px, 1fr))', gap: '16px' }}>
        <Card className="frame-ornate">
          <CardBody>
            <label style={{ display: 'flex', flexDirection: 'column', gap: '8px' }}>
              <span style={{ fontSize: '13px', fontWeight: 600, color: 'var(--night-text-secondary)' }}>
                Dossier de destination
              </span>
              <Input
                placeholder="videos"
                value={settings.destination}
                onChange={(e) => updateSettings({ destination: e.target.value })}
              />
            </label>
          </CardBody>
        </Card>

        <Card className="frame-ornate">
          <CardBody>
            <label style={{ display: 'flex', flexDirection: 'column', gap: '8px' }}>
              <span style={{ fontSize: '13px', fontWeight: 600, color: 'var(--night-text-secondary)' }}>
                Mode de nommage
              </span>
              <select
                value={settings.outputNamingMode}
                onChange={(e) => updateSettings({ outputNamingMode: e.target.value as 'legacy' | 'media-server' })}
                style={{
                  padding: '8px 12px',
                  borderRadius: '12px',
                  border: '1px solid var(--night-border-default)',
                  background: 'var(--night-bg-surface)',
                  color: 'var(--night-text-primary)',
                  fontSize: '13px',
                  cursor: 'pointer',
                }}
              >
                <option value="legacy">Legacy</option>
                <option value="media-server">Media Server (Jellyfin/Plex)</option>
              </select>
            </label>
          </CardBody>
        </Card>

        <Card className="frame-ornate">
          <CardBody>
            <label style={{ display: 'flex', flexDirection: 'column', gap: '8px' }}>
              <span style={{ fontSize: '13px', fontWeight: 600, color: 'var(--night-text-secondary)' }}>
                Workers max
              </span>
              <Input
                type="number"
                min="1"
                max="8"
                value={String(settings.maxWorkers)}
                onChange={(e) => updateSettings({ maxWorkers: parseInt(e.target.value) || 1 })}
              />
            </label>
          </CardBody>
        </Card>

        <Card className="frame-ornate">
          <CardBody>
            <label style={{ display: 'flex', flexDirection: 'column', gap: '8px' }}>
              <span style={{ fontSize: '13px', fontWeight: 600, color: 'var(--night-text-secondary)' }}>
                Téléchargements simultanés
              </span>
              <Input
                type="number"
                min="1"
                max="10"
                value={String(settings.maxConcurrentDownloads)}
                onChange={(e) => updateSettings({ maxConcurrentDownloads: parseInt(e.target.value) || 1 })}
              />
            </label>
          </CardBody>
        </Card>
      </div>

      {/* Options */}
      <Card className="frame-ornate">
        <CardBody>
          <label className="flex" style={{ alignItems: 'center', gap: '12px', cursor: 'pointer' }}>
            <input
              type="checkbox"
              checked={settings.separateLang}
              onChange={(e) => updateSettings({ separateLang: e.target.checked })}
              style={{ width: '18px', height: '18px', cursor: 'pointer' }}
            />
            <span style={{ color: 'var(--night-text-primary)' }}>
              Séparer par langue
            </span>
          </label>
        </CardBody>
      </Card>

      {/* Jellyfin */}
      <Card className="frame-ornate">
        <CardHeader>
          <CardTitle level="h3">Jellyfin</CardTitle>
        </CardHeader>
        <CardBody>
          <div style={{ display: 'flex', flexDirection: 'column', gap: '12px' }}>
            <label style={{ display: 'flex', flexDirection: 'column', gap: '4px' }}>
              <span style={{ fontSize: '13px', fontWeight: 600, color: 'var(--night-text-secondary)' }}>URL</span>
              <Input
                placeholder="http://localhost:8096"
                value={settings.jellyfinUrl || ''}
                onChange={(e) => updateSettings({ jellyfinUrl: e.target.value })}
              />
            </label>
            <label style={{ display: 'flex', flexDirection: 'column', gap: '4px' }}>
              <span style={{ fontSize: '13px', fontWeight: 600, color: 'var(--night-text-secondary)' }}>Clé API</span>
              <Input
                type="password"
                placeholder="Clé API Jellyfin"
                value={settings.jellyfinApiKey || ''}
                onChange={(e) => updateSettings({ jellyfinApiKey: e.target.value })}
              />
            </label>
          </div>
        </CardBody>
      </Card>

      {/* Plex */}
      <Card className="frame-ornate">
        <CardHeader>
          <CardTitle level="h3">Plex</CardTitle>
        </CardHeader>
        <CardBody>
          <div style={{ display: 'flex', flexDirection: 'column', gap: '12px' }}>
            <label style={{ display: 'flex', flexDirection: 'column', gap: '4px' }}>
              <span style={{ fontSize: '13px', fontWeight: 600, color: 'var(--night-text-secondary)' }}>URL</span>
              <Input
                placeholder="http://localhost:32400"
                value={settings.plexUrl || ''}
                onChange={(e) => updateSettings({ plexUrl: e.target.value })}
              />
            </label>
            <label style={{ display: 'flex', flexDirection: 'column', gap: '4px' }}>
              <span style={{ fontSize: '13px', fontWeight: 600, color: 'var(--night-text-secondary)' }}>Token</span>
              <Input
                type="password"
                placeholder="Token Plex"
                value={settings.plexToken || ''}
                onChange={(e) => updateSettings({ plexToken: e.target.value })}
              />
            </label>
            <label style={{ display: 'flex', flexDirection: 'column', gap: '4px' }}>
              <span style={{ fontSize: '13px', fontWeight: 600, color: 'var(--night-text-secondary)' }}>Section ID</span>
              <Input
                placeholder="ID de la section bibliothèque"
                value={settings.plexSectionId || ''}
                onChange={(e) => updateSettings({ plexSectionId: e.target.value })}
              />
            </label>
          </div>
        </CardBody>
      </Card>

      {/* AniList */}
      <Card className="frame-ornate">
        <CardHeader>
          <CardTitle level="h3">AniList</CardTitle>
        </CardHeader>
        <CardBody>
          <label style={{ display: 'flex', flexDirection: 'column', gap: '4px' }}>
            <span style={{ fontSize: '13px', fontWeight: 600, color: 'var(--night-text-secondary)' }}>Token d'accès</span>
            <Input
              type="password"
              placeholder="Token AniList"
              value={settings.anilistToken || ''}
              onChange={(e) => updateSettings({ anilistToken: e.target.value })}
            />
          </label>
        </CardBody>
      </Card>

      <div style={{ display: 'flex', gap: '12px', justifyContent: 'flex-end' }}>
        <Button variant="secondary" onClick={handleReset} disabled={isLoading}>
          Réinitialiser les paramètres
        </Button>
        <Button variant="primary" onClick={handleSave} disabled={!isDirty || isLoading}>
          💾 {isLoading ? 'Enregistrement...' : 'Enregistrer les paramètres'}
        </Button>
      </div>
    </div>
  );
}
