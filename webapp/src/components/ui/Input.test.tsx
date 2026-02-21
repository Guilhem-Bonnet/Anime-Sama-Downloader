import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { Input } from './Input'

describe('Input Component', () => {
  it('should render input element', () => {
    render(<Input />)
    const input = screen.getByRole('textbox')
    expect(input).toBeInTheDocument()
  })

  it('should handle text input', async () => {
    const user = userEvent.setup()
    render(<Input />)
    const input = screen.getByRole('textbox') as HTMLInputElement
    
    await user.type(input, 'test value')
    expect(input.value).toBe('test value')
  })

  it('should render with placeholder', () => {
    render(<Input placeholder="Enter text" />)
    const input = screen.getByPlaceholderText('Enter text')
    expect(input).toBeInTheDocument()
  })

  it('should be disabled when disabled prop is true', () => {
    render(<Input disabled />)
    const input = screen.getByRole('textbox')
    expect(input).toBeDisabled()
  })

  it('should handle different input types', () => {
    const { container } = render(<Input type="password" />)
    const input = container.querySelector('input[type="password"]')
    expect(input).toHaveAttribute('type', 'password')
  })

  it('should handle onChange callback', async () => {
    const user = userEvent.setup()
    const handleChange = () => {}
    render(
      <Input 
        onChange={handleChange}
        placeholder="Search"
      />
    )
    const input = screen.getByPlaceholderText('Search')
    
    await user.type(input, 'test')
    expect((input as HTMLInputElement).value).toBe('test')
  })

  it('should have readonly when specified', () => {
    render(<Input readOnly value="readonly text" />)
    const input = screen.getByDisplayValue('readonly text')
    expect(input).toHaveAttribute('readonly')
  })

  it('should apply custom className', () => {
    const { container } = render(<Input className="custom-input" />)
    const input = container.querySelector('input')
    expect(input).toHaveClass('custom-input')
  })

  it('should work with email type', () => {
    const { container } = render(<Input type="email" />)
    const input = container.querySelector('input[type="email"]')
    expect(input).toHaveAttribute('type', 'email')
  })

  it('should work with number type', () => {
    const { container } = render(<Input type="number" min="0" max="100" />)
    const input = container.querySelector('input[type="number"]')
    expect(input).toHaveAttribute('type', 'number')
    expect(input).toHaveAttribute('min', '0')
    expect(input).toHaveAttribute('max', '100')
  })
})
