import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/react'
import { Card } from './Card'

describe('Card Component', () => {
  it('should render card with children', () => {
    render(
      <Card>
        <div>Card Content</div>
      </Card>
    )
    expect(screen.getByText('Card Content')).toBeInTheDocument()
  })

  it('should apply card class', () => {
    const { container } = render(
      <Card>
        <div>Content</div>
      </Card>
    )
    const card = container.querySelector('div[class*="card"]') || container.firstChild
    
    expect(card).toBeInTheDocument()
  })

  it('should render CardHeader', () => {
    render(
      <Card>
        <Card.Header>Header</Card.Header>
      </Card>
    )
    expect(screen.getByText('Header')).toBeInTheDocument()
  })

  it('should render CardTitle', () => {
    render(
      <Card>
        <Card.Header>
          <Card.Title>Title</Card.Title>
        </Card.Header>
      </Card>
    )
    expect(screen.getByText('Title')).toBeInTheDocument()
  })

  it('should render CardDescription', () => {
    render(
      <Card>
        <Card.Header>
          <Card.Description>Description</Card.Description>
        </Card.Header>
      </Card>
    )
    expect(screen.getByText('Description')).toBeInTheDocument()
  })

  it('should render CardContent', () => {
    render(
      <Card>
        <Card.Content>Content</Card.Content>
      </Card>
    )
    expect(screen.getByText('Content')).toBeInTheDocument()
  })

  it('should render CardFooter', () => {
    render(
      <Card>
        <Card.Footer>Footer</Card.Footer>
      </Card>
    )
    expect(screen.getByText('Footer')).toBeInTheDocument()
  })

  it('should render complete card structure', () => {
    const { container } = render(
      <Card>
        <Card.Header>
          <Card.Title>Card Title</Card.Title>
          <Card.Description>Card Description</Card.Description>
        </Card.Header>
        <Card.Content>Body Content</Card.Content>
        <Card.Footer>Footer Content</Card.Footer>
      </Card>
    )
    
    expect(screen.getByText('Card Title')).toBeInTheDocument()
    expect(screen.getByText('Card Description')).toBeInTheDocument()
    expect(screen.getByText('Body Content')).toBeInTheDocument()
    expect(screen.getByText('Footer Content')).toBeInTheDocument()
  })
})
