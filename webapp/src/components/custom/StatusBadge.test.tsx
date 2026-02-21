import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/react'
import { StatusBadge } from './StatusBadge'

describe('StatusBadge Component', () => {
  it('should render status badge', () => {
    const { container } = render(<StatusBadge status="active" />)
    expect(container.querySelector('span')).toBeInTheDocument()
  })

  it('should display active status', () => {
    const { container } = render(<StatusBadge status="active" />)
    expect(container.textContent?.toLowerCase()).toContain('active')
  })

  it('should display inactive status', () => {
    const { container } = render(<StatusBadge status="inactive" />)
    expect(container.textContent?.toLowerCase()).toContain('inactive')
  })

  it('should display pending status', () => {
    const { container } = render(<StatusBadge status="pending" />)
    expect(container.textContent?.toLowerCase()).toContain('pending')
  })

  it('should display completed status', () => {
    const { container } = render(<StatusBadge status="completed" />)
    expect(container.textContent?.toLowerCase()).toContain('completed')
  })

  it('should render with custom size', () => {
    const { container } = render(
      <StatusBadge status="active" size="lg" />
    )
    expect(container.querySelector('span')).toBeInTheDocument()
  })

  it('should apply custom className', () => {
    const { container } = render(
      <StatusBadge status="active" className="custom" />
    )
    const badge = container.querySelector('span')
    expect(badge).toHaveClass('custom')
  })

  it('should handle error status', () => {
    const { container } = render(<StatusBadge status="error" />)
    expect(container).toBeInTheDocument()
  })
})
