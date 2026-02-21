import React from 'react';
import { DownloadMonitor } from '../components/DownloadMonitor';
import { Card, CardBody, CardHeader, CardTitle } from '../components/ui/Card';

export const DownloadsPage = React.memo(function DownloadsPage() {
  return (
    <Card className="frame-ornate">
      <CardHeader>
        <CardTitle level="h2">Téléchargements</CardTitle>
        <div className="kanji-brush" style={{ marginTop: '6px' }}>
          保存
        </div>
      </CardHeader>
      <CardBody>
        <DownloadMonitor />
      </CardBody>
    </Card>
  );
});
