import React from 'react';
import { Button } from '../ui/Button';
import { Card, CardBody, CardFooter, CardHeader, CardTitle } from '../ui/Card';

export class ErrorBoundary extends React.Component<
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
            background: 'var(--night-bg-base)',
          }}
        >
          <Card style={{ maxWidth: '520px', width: '100%' }}>
            <CardHeader>
              <CardTitle level="h2">Something went wrong</CardTitle>
            </CardHeader>
            <CardBody>
              <p style={{ color: 'var(--night-error-text)' }}>{this.state.errorMessage}</p>
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
