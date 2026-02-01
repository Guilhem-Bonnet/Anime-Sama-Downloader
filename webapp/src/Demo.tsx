import React, { useState } from 'react';
import { Button } from './components/ui/Button';
import { Card, CardHeader, CardTitle, CardBody, CardFooter } from './components/ui/Card';
import { Input, TextArea, Select } from './components/ui/Input';
import { Badge } from './components/ui/Badge';

export const Demo: React.FC = () => {
  const [inputValue, setInputValue] = useState('');
  const [selectValue, setSelectValue] = useState('option1');

  return (
    <div className="p-6" style={{ maxWidth: '1200px', margin: '0 auto' }}>
      <style>{`
        .demo-section {
          margin-bottom: 48px;
        }
        .demo-section h2 {
          font-size: var(--text-h1);
          font-weight: var(--text-h1-weight);
          margin-bottom: 24px;
          padding-bottom: 12px;
          border-bottom: 1px solid var(--sakura-border-default);
        }
        .demo-grid {
          display: grid;
          grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
          gap: 16px;
          margin-bottom: 24px;
        }
        .demo-item {
          display: flex;
          gap: 8px;
          flex-wrap: wrap;
        }
        .fade-in {
          animation: fadeIn 0.4s ease-out;
        }
      `}</style>

      {/* HEADER */}
      <div className="fade-in" style={{ marginBottom: '48px' }}>
        <h1 style={{ fontSize: 'var(--text-display)', marginBottom: '8px' }}>
          🌸 Sakura Night Design System
        </h1>
        <p style={{ color: 'var(--sakura-text-secondary)', fontSize: 'var(--text-body)' }}>
          Prototype complet du design system Anime-Sama Downloader
        </p>
      </div>

      {/* BUTTONS */}
      <div className="demo-section">
        <h2>Buttons</h2>

        <h3 style={{ fontSize: 'var(--text-h3)', marginTop: '24px', marginBottom: '12px' }}>
          Primary
        </h3>
        <div className="demo-item">
          <Button variant="primary">Primary Button</Button>
          <Button variant="primary" size="sm">
            Small
          </Button>
          <Button variant="primary" size="lg">
            Large
          </Button>
          <Button variant="primary" disabled>
            Disabled
          </Button>
          <Button variant="primary" isLoading>
            Loading
          </Button>
        </div>

        <h3 style={{ fontSize: 'var(--text-h3)', marginTop: '24px', marginBottom: '12px' }}>
          Secondary
        </h3>
        <div className="demo-item">
          <Button variant="secondary">Secondary Button</Button>
          <Button variant="secondary" size="sm">
            Small
          </Button>
          <Button variant="secondary" size="lg">
            Large
          </Button>
        </div>

        <h3 style={{ fontSize: 'var(--text-h3)', marginTop: '24px', marginBottom: '12px' }}>
          Ghost
        </h3>
        <div className="demo-item">
          <Button variant="ghost">Ghost Button</Button>
          <Button variant="ghost" size="sm">
            Small
          </Button>
          <Button variant="ghost" size="lg">
            Large
          </Button>
        </div>

        <h3 style={{ fontSize: 'var(--text-h3)', marginTop: '24px', marginBottom: '12px' }}>
          Danger
        </h3>
        <div className="demo-item">
          <Button variant="danger">Delete</Button>
          <Button variant="danger" size="sm">
            Small Delete
          </Button>
        </div>
      </div>

      {/* CARDS */}
      <div className="demo-section">
        <h2>Cards</h2>
        <div className="grid" style={{ gridTemplateColumns: 'repeat(auto-fit, minmax(300px, 1fr))' }}>
          <Card>
            <CardHeader>
              <CardTitle>Card Title</CardTitle>
            </CardHeader>
            <CardBody>
              <p>
                Ceci est un contenu de card standard. Les cards utilisent un
                gradient subtil et un backdrop blur pour l'effet verre.
              </p>
            </CardBody>
            <CardFooter>
              <Button variant="ghost" size="sm">
                Cancel
              </Button>
              <Button variant="primary" size="sm">
                Confirm
              </Button>
            </CardFooter>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle>Card with Badge</CardTitle>
            </CardHeader>
            <CardBody>
              <div style={{ display: 'flex', gap: '8px', marginBottom: '12px' }}>
                <Badge variant="primary">Primary</Badge>
                <Badge variant="success">Success</Badge>
                <Badge variant="warning">Warning</Badge>
              </div>
              <p>Les badges s'affichent bien à côté du texte.</p>
            </CardBody>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle>Empty Card</CardTitle>
            </CardHeader>
            <CardBody>
              <p style={{ color: 'var(--sakura-text-muted)' }}>
                Une card vide avec du contenu minimal.
              </p>
            </CardBody>
          </Card>
        </div>
      </div>

      {/* INPUTS */}
      <div className="demo-section">
        <h2>Inputs</h2>
        <div style={{ maxWidth: '400px', display: 'flex', flexDirection: 'column', gap: '20px' }}>
          <Input
            label="Text Input"
            placeholder="Entrez du texte..."
            value={inputValue}
            onChange={(e) => setInputValue(e.target.value)}
            hint="Ceci est un hint utile"
          />

          <Input
            label="Email Input"
            type="email"
            placeholder="mail@example.com"
            error="Email invalide"
          />

          <Input
            label="Password Input"
            type="password"
            placeholder="••••••••"
          />

          <TextArea
            label="Textarea"
            placeholder="Entrez un message plus long..."
            hint="Vous pouvez taper autant de texte que vous voulez"
          />

          <Select
            label="Select Input"
            value={selectValue}
            onChange={(e) => setSelectValue(e.target.value)}
            options={[
              { value: 'option1', label: 'Option 1' },
              { value: 'option2', label: 'Option 2' },
              { value: 'option3', label: 'Option 3' },
            ]}
            hint="Sélectionnez une option"
          />
        </div>
      </div>

      {/* BADGES */}
      <div className="demo-section">
        <h2>Badges</h2>
        <div className="demo-item">
          <Badge variant="primary">Primary</Badge>
          <Badge variant="secondary">Secondary</Badge>
          <Badge variant="success">Success</Badge>
          <Badge variant="warning">Warning</Badge>
          <Badge variant="error">Error</Badge>
          <Badge variant="info">Info</Badge>
        </div>
      </div>

      {/* COLORS */}
      <div className="demo-section">
        <h2>Color Palette</h2>
        <div className="grid">
          {/* Backgrounds */}
          <div
            style={{
              background: 'var(--sakura-bg-base)',
              padding: '16px',
              borderRadius: '8px',
              border: '1px solid var(--sakura-border-default)',
              textAlign: 'center',
            }}
          >
            <p style={{ fontSize: 'var(--text-caption)' }}>--sakura-bg-base</p>
            <p style={{ fontSize: '12px', color: 'var(--sakura-text-muted)' }}>#0A0E1A</p>
          </div>

          <div
            style={{
              background: 'var(--sakura-bg-surface)',
              padding: '16px',
              borderRadius: '8px',
              border: '1px solid var(--sakura-border-default)',
              textAlign: 'center',
            }}
          >
            <p style={{ fontSize: 'var(--text-caption)' }}>--sakura-bg-surface</p>
            <p style={{ fontSize: '12px', color: 'var(--sakura-text-muted)' }}>#1A1F2E</p>
          </div>

          <div
            style={{
              background: 'var(--sakura-bg-elevated)',
              padding: '16px',
              borderRadius: '8px',
              border: '1px solid var(--sakura-border-default)',
              textAlign: 'center',
            }}
          >
            <p style={{ fontSize: 'var(--text-caption)' }}>--sakura-bg-elevated</p>
            <p style={{ fontSize: '12px', color: 'var(--sakura-text-muted)' }}>#252A3B</p>
          </div>

          {/* Accents */}
          <div
            style={{
              background: 'var(--sakura-accent-magenta-500)',
              padding: '16px',
              borderRadius: '8px',
              textAlign: 'center',
              color: 'white',
            }}
          >
            <p style={{ fontSize: 'var(--text-caption)' }}>Magenta Primary</p>
            <p style={{ fontSize: '12px' }}>#D946EF</p>
          </div>

          <div
            style={{
              background: 'var(--sakura-accent-cyan-500)',
              padding: '16px',
              borderRadius: '8px',
              textAlign: 'center',
              color: 'white',
            }}
          >
            <p style={{ fontSize: 'var(--text-caption)' }}>Cyan Secondary</p>
            <p style={{ fontSize: '12px' }}>#06B6D4</p>
          </div>

          <div
            style={{
              background: 'var(--sakura-pink-500)',
              padding: '16px',
              borderRadius: '8px',
              textAlign: 'center',
              color: 'white',
            }}
          >
            <p style={{ fontSize: 'var(--text-caption)' }}>Sakura Pink</p>
            <p style={{ fontSize: '12px' }}>#FB6F8A</p>
          </div>
        </div>
      </div>

      {/* PROGRESS BAR */}
      <div className="demo-section">
        <h2>Progress Bar</h2>
        <div style={{ display: 'flex', flexDirection: 'column', gap: '16px', maxWidth: '400px' }}>
          <div>
            <p style={{ marginBottom: '8px', fontSize: 'var(--text-label)' }}>25%</p>
            <div className="progress-bar">
              <div
                className="progress-fill"
                style={{ width: '25%', animation: 'none' }}
              />
            </div>
          </div>

          <div>
            <p style={{ marginBottom: '8px', fontSize: 'var(--text-label)' }}>50%</p>
            <div className="progress-bar">
              <div
                className="progress-fill"
                style={{ width: '50%', animation: 'none' }}
              />
            </div>
          </div>

          <div>
            <p style={{ marginBottom: '8px', fontSize: 'var(--text-label)' }}>75%</p>
            <div className="progress-bar">
              <div
                className="progress-fill"
                style={{ width: '75%', animation: 'none' }}
              />
            </div>
          </div>

          <div>
            <p style={{ marginBottom: '8px', fontSize: 'var(--text-label)' }}>Animated</p>
            <div className="progress-bar">
              <div className="progress-fill" style={{ width: '100%' }} />
            </div>
          </div>
        </div>
      </div>

      {/* TYPOGRAPHY */}
      <div className="demo-section">
        <h2>Typography</h2>
        <h1 style={{ fontSize: 'var(--text-display)', marginBottom: '16px' }}>
          Display (32px, Bold)
        </h1>
        <h2 style={{ fontSize: 'var(--text-h1)', marginBottom: '16px' }}>
          Heading 1 (24px, Semi-Bold)
        </h2>
        <h3 style={{ fontSize: 'var(--text-h2)', marginBottom: '16px' }}>
          Heading 2 (20px, Semi-Bold)
        </h3>
        <h4 style={{ fontSize: 'var(--text-h3)', marginBottom: '16px' }}>
          Heading 3 (16px, Semi-Bold)
        </h4>
        <p style={{ fontSize: 'var(--text-body-lg)', marginBottom: '12px' }}>
          Body Large (16px, Regular)
        </p>
        <p style={{ fontSize: 'var(--text-body)', marginBottom: '12px' }}>
          Body Default (14px, Regular)
        </p>
        <p style={{ fontSize: 'var(--text-body-sm)', marginBottom: '12px' }}>
          Body Small (13px, Regular)
        </p>
        <p style={{ fontSize: 'var(--text-caption)' }}>Caption (12px, Regular)</p>
      </div>
    </div>
  );
};
