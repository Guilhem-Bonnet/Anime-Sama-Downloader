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
    const badge = container.querySelector('div')
    
    expect(badge).toHaveClass('badge')
  })

  it('should apply custom variant classes', () => {
    const { container } = render(
      <Badge variant="destructive">Error</Badge>
    )
    const badge = container.querySelector('div')
    
    expect(badge).toHaveClass('destructive')
  })

  it('should apply outline variant', () => {
    const { container } = render(
      <Badge variant="outline">Outline Badge</Badge>
    )
    const badge = container.querySelector('div')
    
    expect(badge).toHaveClass('outline')
  })

  it('should render with custom className', () => {
    const { container } = render(
      <Badge className="custom-class">Badge</Badge>
    )
    const badge = container.querySelector('div')
    
    expect(badge).toHaveClass('custom-class')
  })

  it('should render secondary variant', () => {
    const { container } = render(
      <Badge variant="secondary">Secondary</Badge>
    )
    const badge = container.querySelector('div')
    
    expect(badge).toHaveClass('secondary')
  })
})
