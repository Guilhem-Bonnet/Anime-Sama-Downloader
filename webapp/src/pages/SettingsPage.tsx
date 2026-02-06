import React from 'react';
import { Button } from '../components/ui/Button';
import { Card, CardBody, CardHeader, CardTitle } from '../components/ui/Card';
import { Input } from '../components/ui/Input';

export function SettingsPage() {
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
              <Input placeholder="/media/anime" defaultValue="/media/anime" />
            </label>
          </CardBody>
        </Card>

        <Card className="frame-ornate">
          <CardBody>
            <label style={{ display: 'flex', flexDirection: 'column', gap: '8px' }}>
              <span style={{ fontSize: '13px', fontWeight: 600, color: 'var(--sakura-text-secondary)' }}>
                Téléchargements simultanés
              </span>
              <Input type="number" min="1" max="5" defaultValue="2" />
            </label>
          </CardBody>
        </Card>

        <Card className="frame-ornate">
          <CardBody>
            <label style={{ display: 'flex', flexDirection: 'column', gap: '8px' }}>
              <span style={{ fontSize: '13px', fontWeight: 600, color: 'var(--sakura-text-secondary)' }}>
                Qualité préférée
              </span>
              <select
                style={{
                  padding: '8px 12px',
                  borderRadius: '12px',
                  border: '1px solid var(--sakura-border-default)',
                  background: 'var(--sakura-bg-surface)',
                  color: 'var(--sakura-text-primary)',
                  fontSize: '13px',
                  cursor: 'pointer',
                }}
              >
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

      <Card
        className="frame-ornate"
        style={{ background: 'linear-gradient(135deg, rgba(143,106,61,0.05), rgba(125,114,103,0.05))' }}
      >
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
  );
}
