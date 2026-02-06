import React from 'react';
import { DownloadMonitor } from '../components/DownloadMonitor';
import { Card, CardBody, CardHeader, CardTitle } from '../components/ui/Card';
import { EmptyDownloadsIllustration } from '../components/illustrations/SakuraIllustrations';

export function DownloadsPage() {
  return (
    <Card className="frame-ornate">
      <CardHeader>
        <CardTitle level="h2">Téléchargements</CardTitle>
        <div className="kanji-brush" style={{ marginTop: '6px' }}>
          保存
        </div>
      </CardHeader>
      <CardBody>
        <div style={{ marginBottom: '24px', maxHeight: '180px', overflow: 'hidden', opacity: 0.8 }}>
          <EmptyDownloadsIllustration />
        </div>
        <DownloadMonitor />
      </CardBody>
    </Card>
  );
}
