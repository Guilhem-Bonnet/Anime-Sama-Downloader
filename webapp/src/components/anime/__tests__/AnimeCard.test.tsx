import React from 'react';
import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import AnimeCard from '../AnimeCard';

describe('AnimeCard', () => {
  const defaultProps = {
    id: 'anime-1',
    title: 'Demon Slayer',
    coverUrl: 'https://example.com/cover.jpg',
    season: 'S1',
    language: 'VOSTFR' as const,
    status: 'ongoing' as const,
    onDetails: vi.fn(),
    onDownload: vi.fn(),
  };

  it('renders with default props', () => {
    const { container } = render(<AnimeCard {...defaultProps} />);
    
    expect(screen.getByText('Demon Slayer')).toBeInTheDocument();
    expect(screen.getByText('S1')).toBeInTheDocument();
    expect(screen.getByText('VOSTFR')).toBeInTheDocument();
    expect(screen.getByText('En cours')).toBeInTheDocument();
    
    const article = container.querySelector('article');
    expect(article).toHaveAttribute('aria-label', expect.stringContaining('Demon Slayer'));
  });

  it('renders with list variant', () => {
    const { container } = render(
      <AnimeCard {...defaultProps} variant="list" />
    );
    
    const card = container.querySelector('.anime-card--list');
    expect(card).toBeInTheDocument();
  });

  it('renders with grid variant (default)', () => {
    const { container } = render(
      <AnimeCard {...defaultProps} variant="grid" />
    );
    
    const card = container.querySelector('.anime-card--grid');
    expect(card).toBeInTheDocument();
  });

  it('handles click on Details button', async () => {
    const user = userEvent.setup();
    render(<AnimeCard {...defaultProps} />);
    
    const detailsBtn = screen.getByRole('button', { name: /détails/i });
    await user.click(detailsBtn);
    
    expect(defaultProps.onDetails).toHaveBeenCalledTimes(1);
  });

  it('handles click on Download button', async () => {
    const user = userEvent.setup();
    render(<AnimeCard {...defaultProps} />);
    
    const downloadBtn = screen.getByRole('button', { name: /télécharger/i });
    await user.click(downloadBtn);
    
    expect(defaultProps.onDownload).toHaveBeenCalledTimes(1);
  });

  it('applies disabled state', () => {
    render(<AnimeCard {...defaultProps} disabled={true} />);
    
    const detailsBtn = screen.getByRole('button', { name: /détails/i });
    const downloadBtn = screen.getByRole('button', { name: /télécharger/i });
    
    expect(detailsBtn).toBeDisabled();
    expect(downloadBtn).toBeDisabled();
  });

  it('applies selected state', () => {
    const { container } = render(
      <AnimeCard {...defaultProps} selected={true} />
    );
    
    const card = container.querySelector('.anime-card--selected');
    expect(card).toBeInTheDocument();
  });

  it('renders different status colors', () => {
    const statuses = ['ongoing', 'completed', 'upcoming'] as const;
    
    statuses.forEach((status) => {
      const { unmount } = render(
        <AnimeCard {...defaultProps} status={status} />
      );
      
      const badge = screen.getByLabelText(new RegExp(`Statut`));
      expect(badge).toBeInTheDocument();
      
      unmount();
    });
  });

  it('snapshot test', () => {
    const { container } = render(<AnimeCard {...defaultProps} />);
    expect(container.firstChild).toMatchSnapshot();
  });
});
