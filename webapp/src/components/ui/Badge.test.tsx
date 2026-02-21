import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/react'
import { Badge } from './Badge'

describe('Badge Component', () => {
  it('should render badge with text', () => {
    render(<Badge>New</Badge>)
    expect(screen.getByText('New')).toBeInTheDocument()
  })

  it('should apply default variant', () => {
    const { container } = render(<Badge>Tag</Badge>)
    const badge = container.querySelector('span')
    
    expect(badge).toHaveClass('badge', 'badge-primary')
  })

  it('should apply error variant classes', () => {
    const { container } = render(
      <Badge variant="error">Error</Badge>
    )
    const badge = container.querySelector('span')
    
    expect(badge).toHaveClass('badge-error')
  })

  it('should apply warning variant', () => {
    const { container } = render(
      <Badge variant="warning">Warning Badge</Badge>
    )
    const badge = container.querySelector('span')
    
    expect(badge).toHaveClass('badge-warning')
  })

  it('should render with custom className', () => {
    const { container } = render(
      <Badge className="custom-class">Badge</Badge>
    )
    const badge = container.querySelector('span')
    
    expect(badge).toHaveClass('custom-class')
  })

  it('should render secondary variant', () => {
    const { container } = render(
      <Badge variant="secondary">Secondary</Badge>
    )
    const badge = container.querySelector('span')
    
    expect(badge).toHaveClass('badge-secondary')
  })
})
