import React from 'react';
import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import EpisodeRow from '../EpisodeRow';

describe('EpisodeRow', () => {
  const defaultProps = {
    number: 1,
    title: 'Première attaque',
    duration: '24min',
    status: 'available' as const,
    selected: false,
    onChange: vi.fn(),
  };

  it('renders with default props', () => {
    render(<EpisodeRow {...defaultProps} />);
    
    expect(screen.getByText('Ép. 1')).toBeInTheDocument();
    expect(screen.getByText('Première attaque')).toBeInTheDocument();
    expect(screen.getByText('24min')).toBeInTheDocument();
    expect(screen.getByText('Disponible')).toBeInTheDocument();
  });

  it('renders without optional title and duration', () => {
    render(
      <EpisodeRow
        number={2}
        status="available"
        selected={false}
        onChange={vi.fn()}
      />
    );
    
    expect(screen.getByText('Ép. 2')).toBeInTheDocument();
    expect(screen.queryByText('Première attaque')).not.toBeInTheDocument();
  });

  it('handles checkbox change', async () => {
    const user = userEvent.setup();
    const onChange = vi.fn();
    
    render(
      <EpisodeRow {...defaultProps} onChange={onChange} />
    );
    
    const checkbox = screen.getByRole('checkbox');
    await user.click(checkbox);
    
    expect(onChange).toHaveBeenCalledWith(true);
  });

  it('reflects selected state', () => {
    const { container } = render(
      <EpisodeRow {...defaultProps} selected={true} />
    );
    
    const checkbox = screen.getByRole('checkbox') as HTMLInputElement;
    expect(checkbox.checked).toBe(true);
    
    const row = container.querySelector('.episode-row--selected');
    expect(row).toBeInTheDocument();
  });

  it('handles disabled state', () => {
    render(
      <EpisodeRow {...defaultProps} disabled={true} />
    );
    
    const checkbox = screen.getByRole('checkbox');
    expect(checkbox).toBeDisabled();
  });

  it('renders different status badges', () => {
    const statuses = ['available', 'downloading', 'downloaded'] as const;
    
    statuses.forEach((status) => {
      const { unmount } = render(
        <EpisodeRow {...defaultProps} status={status} />
      );
      
      // Verify badge is rendered
      const badge = screen.getByText(/disponible|téléchargement|téléchargé/i);
      expect(badge).toBeInTheDocument();
      
      unmount();
    });
  });

  it('shows loading spinner for downloading status', () => {
    const { container } = render(
      <EpisodeRow {...defaultProps} status="downloading" />
    );
    
    const spinner = container.querySelector('.episode-row__spinner');
    expect(spinner).toBeInTheDocument();
  });

  it('snapshot test', () => {
    const { container } = render(<EpisodeRow {...defaultProps} />);
    expect(container.firstChild).toMatchSnapshot();
  });
});
