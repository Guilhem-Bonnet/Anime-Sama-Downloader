import React from 'react';
import { Button } from '../components/ui/Button';
import { Card, CardBody, CardHeader, CardTitle } from '../components/ui/Card';
import { Badge } from '../components/ui/Badge';

export const RulesPage = React.memo(function RulesPage() {
  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: '24px' }}>
      <Card className="frame-ornate">
        <CardHeader>
          <CardTitle level="h2">Règles d'Automatisation</CardTitle>
          <div className="kanji-brush" style={{ marginTop: '6px' }}>
            規則
          </div>
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
                • Motif: "Attack on Titan*" <br />
                • Limite: 2 téléchargements simultanés <br />
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
                • Dernière vérification: il y a 2 heures <br />
                • Épisodes téléchargés: 24 <br />
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
  );
});
