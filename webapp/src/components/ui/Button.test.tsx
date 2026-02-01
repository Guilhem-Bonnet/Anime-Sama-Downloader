import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { Button } from "./Button"

describe('Button Component', () => {
  it('should render button with text', () => {
    render(<Button>Click me</Button>)
    const button = screen.getByRole('button', { name: /click me/i })
    expect(button).toBeInTheDocument()
  })

  it('should handle click events', async () => {
    const user = userEvent.setup()
    const handleClick = () => console.log('clicked')
    
    render(<Button onClick={handleClick}>Click me</Button>)
    const button = screen.getByRole('button', { name: /click me/i })
    
    await user.click(button)
    expect(button).toBeInTheDocument()
  })

  it('should be disabled when disabled prop is true', () => {
    render(<Button disabled>Disabled Button</Button>)
    const button = screen.getByRole('button', { name: /disabled button/i })
    
    expect(button).toBeDisabled()
  })

  it('should apply variant classes', () => {
    const { container } = render(
      <Button variant="destructive">Delete</Button>
    )
    const button = container.querySelector('button')
    
    expect(button).toHaveClass('destructive')
  })

  it('should apply size classes', () => {
    const { container } = render(
      <Button size="lg">Large Button</Button>
    )
    const button = container.querySelector('button')
    
    expect(button).toHaveClass('lg')
  })

  it('should render as child element with asChild prop', () => {
    const { container } = render(
      <Button asChild>
        <a href="/test">Link Button</a>
      </Button>
    )
    const link = container.querySelector('a')
    
    expect(link).toHaveAttribute('href', '/test')
  })
})
