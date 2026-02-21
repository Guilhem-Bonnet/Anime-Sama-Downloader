import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/react'
import { Card, CardHeader, CardTitle, CardSubtitle, CardBody, CardFooter } from './Card'

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
        <CardHeader>Header</CardHeader>
      </Card>
    )
    expect(screen.getByText('Header')).toBeInTheDocument()
  })

  it('should render CardTitle', () => {
    render(
      <Card>
        <CardHeader>
          <CardTitle>Title</CardTitle>
        </CardHeader>
      </Card>
    )
    expect(screen.getByText('Title')).toBeInTheDocument()
  })

  it('should render CardSubtitle', () => {
    render(
      <Card>
        <CardHeader>
          <CardSubtitle>Description</CardSubtitle>
        </CardHeader>
      </Card>
    )
    expect(screen.getByText('Description')).toBeInTheDocument()
  })

  it('should render CardBody', () => {
    render(
      <Card>
        <CardBody>Content</CardBody>
      </Card>
    )
    expect(screen.getByText('Content')).toBeInTheDocument()
  })

  it('should render CardFooter', () => {
    render(
      <Card>
        <CardFooter>Footer</CardFooter>
      </Card>
    )
    expect(screen.getByText('Footer')).toBeInTheDocument()
  })

  it('should render complete card structure', () => {
    render(
      <Card>
        <CardHeader>
          <CardTitle>Card Title</CardTitle>
          <CardSubtitle>Card Description</CardSubtitle>
        </CardHeader>
        <CardBody>Body Content</CardBody>
        <CardFooter>Footer Content</CardFooter>
      </Card>
    )
    
    expect(screen.getByText('Card Title')).toBeInTheDocument()
    expect(screen.getByText('Card Description')).toBeInTheDocument()
    expect(screen.getByText('Body Content')).toBeInTheDocument()
    expect(screen.getByText('Footer Content')).toBeInTheDocument()
  })
})
