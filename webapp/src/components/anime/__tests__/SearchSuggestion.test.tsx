import React from 'react';
import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import SearchSuggestion from '../SearchSuggestion';

describe('SearchSuggestion', () => {
  const defaultProps = {
    title: 'Demon Slayer',
    season: 'S1',
    language: 'VOSTFR' as const,
    query: 'demon',
    onSelect: vi.fn(),
  };

  it('renders with default props', () => {
    render(<SearchSuggestion {...defaultProps} />);
    
    expect(screen.getByText(/demon slayer/i)).toBeInTheDocument();
    expect(screen.getByText('S1')).toBeInTheDocument();
    expect(screen.getByText('VOSTFR')).toBeInTheDocument();
  });

  it('highlights query in title', () => {
    const { container } = render(
      <SearchSuggestion {...defaultProps} query="demon" />
    );
    
    const mark = container.querySelector('mark');
    expect(mark).toBeInTheDocument();
    expect(mark?.textContent).toMatch(/demon/i);
  });

  it('renders without season (optional)', () => {
    render(
      <SearchSuggestion
        title="Jujutsu Kaisen"
        language="VOSTFR"
        query="jujutsu"
        onSelect={vi.fn()}
      />
    );
    
    expect(screen.getByText(/jujutsu kaisen/i)).toBeInTheDocument();
    expect(screen.queryByText('S1')).not.toBeInTheDocument();
  });

  it('handles click to select', async () => {
    const user = userEvent.setup();
    const onSelect = vi.fn();
    
    render(
      <SearchSuggestion {...defaultProps} onSelect={onSelect} />
    );
    
    const suggestion = screen.getByRole('option');
    await user.click(suggestion);
    
    expect(onSelect).toHaveBeenCalledOnce();
  });

  it('handles Enter key to select', async () => {
    const user = userEvent.setup();
    const onSelect = vi.fn();
    
    render(
      <SearchSuggestion {...defaultProps} onSelect={onSelect} />
    );
    
    const suggestion = screen.getByRole('option');
    suggestion.focus();
    await user.keyboard('{Enter}');
    
    expect(onSelect).toHaveBeenCalledOnce();
  });

  it('handles Space key to select', async () => {
    const user = userEvent.setup();
    const onSelect = vi.fn();
    
    render(
      <SearchSuggestion {...defaultProps} onSelect={onSelect} />
    );
    
    const suggestion = screen.getByRole('option');
    suggestion.focus();
    await user.keyboard(' ');
    
    expect(onSelect).toHaveBeenCalledOnce();
  });

  it('is keyboard focusable', () => {
    render(<SearchSuggestion {...defaultProps} />);
    
    const suggestion = screen.getByRole('option');
    expect(suggestion).toHaveAttribute('tabIndex', '0');
  });

  it('snapshot test', () => {
    const { container } = render(<SearchSuggestion {...defaultProps} />);
    expect(container.firstChild).toMatchSnapshot();
  });
});
