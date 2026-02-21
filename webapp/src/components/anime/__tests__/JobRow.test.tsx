import React from 'react';
import { describe, it, expect, vi } from 'vitest';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import JobRow from '../JobRow';

describe('JobRow', () => {
  const defaultProps = {
    id: 'job-1',
    animeTitle: 'Demon Slayer',
    episode: 1,
    progress: 50,
    eta: '5 min',
    speed: '5 MB/s',
    status: 'downloading' as const,
    onCancel: vi.fn(),
  };

  it('renders with default props', () => {
    render(<JobRow {...defaultProps} />);
    
    expect(screen.getByText('Demon Slayer')).toBeInTheDocument();
    expect(screen.getByText('Ép. 1')).toBeInTheDocument();
    expect(screen.getByText('50%')).toBeInTheDocument();
    expect(screen.getByText('Téléchargement')).toBeInTheDocument();
  });

  it('renders progress bar with correct value', () => {
    const { container } = render(<JobRow {...defaultProps} progress={75} />);
    
    const progressBar = container.querySelector('[role="progressbar"]');
    expect(progressBar).toHaveAttribute('aria-valuenow', '75');
  });

  it('renders optional metadata (ETA and Speed)', () => {
    render(<JobRow {...defaultProps} eta="3 min" speed="8 MB/s" />);
    
    expect(screen.getByText('3 min')).toBeInTheDocument();
    expect(screen.getByText('8 MB/s')).toBeInTheDocument();
  });

  it('handles pause action for downloading status', async () => {
    const user = userEvent.setup();
    const onPause = vi.fn();
    
    render(
      <JobRow {...defaultProps} status="downloading" onPause={onPause} />
    );
    
    const pauseBtn = screen.getByRole('button', { name: /pause/i });
    await user.click(pauseBtn);
    
    expect(onPause).toHaveBeenCalledOnce();
  });

  it('handles resume action for paused status', async () => {
    const user = userEvent.setup();
    const onResume = vi.fn();
    
    render(
      <JobRow {...defaultProps} status="paused" onResume={onResume} />
    );
    
    const resumeBtn = screen.getByRole('button', { name: /reprendre/i });
    await user.click(resumeBtn);
    
    expect(onResume).toHaveBeenCalledOnce();
  });

  it('handles retry action for failed status', async () => {
    const user = userEvent.setup();
    const onRetry = vi.fn();
    
    render(
      <JobRow {...defaultProps} status="failed" onRetry={onRetry} />
    );
    
    const retryBtn = screen.getByRole('button', { name: /réessayer/i });
    await user.click(retryBtn);
    
    expect(onRetry).toHaveBeenCalledOnce();
  });

  it('handles cancel action', async () => {
    const user = userEvent.setup();
    const onCancel = vi.fn();
    
    render(
      <JobRow {...defaultProps} onCancel={onCancel} />
    );
    
    const cancelBtn = screen.getByRole('button', { name: /annuler/i });
    await user.click(cancelBtn);
    
    expect(onCancel).toHaveBeenCalledOnce();
  });

  it('hides actions for queued status', () => {
    render(
      <JobRow {...defaultProps} status="queued" />
    );
    
    const buttons = screen.queryAllByRole('button');
    // No action buttons for queued status
    expect(buttons.length).toBe(0);
  });

  it('hides actions for completed status', () => {
    render(
      <JobRow {...defaultProps} status="completed" />
    );
    
    const buttons = screen.queryAllByRole('button');
    expect(buttons.length).toBe(0);
  });

  it('renders different status colors', () => {
    const statuses = ['queued', 'downloading', 'paused', 'completed', 'failed'] as const;
    
    statuses.forEach((status) => {
      const { unmount } = render(
        <JobRow {...defaultProps} status={status} />
      );
      
      const matches = screen.getAllByText(/en attente|téléchargement|mis en pause|terminé|erreur/i);
      expect(matches.length).toBeGreaterThan(0);
      
      unmount();
    });
  });

  it('snapshot test', () => {
    const { container } = render(<JobRow {...defaultProps} />);
    expect(container.firstChild).toMatchSnapshot();
  });
});
