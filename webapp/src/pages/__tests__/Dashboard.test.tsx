import { render, screen, waitFor } from '@testing-library/react'
import { BrowserRouter } from 'react-router-dom'
import { describe, it, expect, beforeEach, vi, afterEach } from 'vitest'
import Dashboard from '../Dashboard'
import * as api from '../../api'

// Mock the canonical API module
vi.mock('../../api', () => ({
  apiListJobs: vi.fn(),
  apiListSubscriptions: vi.fn(),
  apiCancelJob: vi.fn(),
  apiCreateSubscription: vi.fn(),
  apiDeleteSubscription: vi.fn(),
  apiSyncSubscription: vi.fn(),
  apiSyncAll: vi.fn(),
}))

// Mock AniList GraphQL
const originalFetch = globalThis.fetch
beforeEach(() => {
  globalThis.fetch = vi.fn() as any
})
afterEach(() => {
  globalThis.fetch = originalFetch
})

const mockJobs: api.Job[] = [
  {
    id: '1',
    type: 'download',
    state: 'completed',
    progress: 1.0,
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString(),
    params: { animeTitle: 'Jujutsu Kaisen', episodeNumber: 1 },
  },
  {
    id: '2',
    type: 'download',
    state: 'running',
    progress: 0.45,
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString(),
    params: { animeTitle: 'Attack on Titan', episodeNumber: 5 },
  },
]

const mockSubscriptions: api.Subscription[] = [
  {
    id: '1',
    baseUrl: 'https://anime-sama.fr/catalogue/one-piece/vostfr/episodes/',
    label: 'One Piece',
    player: 'sendvid',
    lastScheduledEpisode: 999,
    lastDownloadedEpisode: 999,
    lastAvailableEpisode: 1000,
    nextCheckAt: new Date().toISOString(),
    lastCheckedAt: new Date().toISOString(),
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString(),
  },
  {
    id: '2',
    baseUrl: 'https://anime-sama.fr/catalogue/bleach/vostfr/episodes/',
    label: 'Bleach',
    player: 'sendvid',
    lastScheduledEpisode: 50,
    lastDownloadedEpisode: 50,
    lastAvailableEpisode: 50,
    nextCheckAt: new Date().toISOString(),
    lastCheckedAt: new Date().toISOString(),
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString(),
  },
]

const mockTrendingAnime = {
  id: 123,
  title: 'Trending Anime',
  coverImage: 'https://example.com/trending.jpg',
  season: 'FALL',
  year: 2024,
}

describe('Dashboard Page', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    ;(api.apiListJobs as any).mockResolvedValue([])
    ;(api.apiListSubscriptions as any).mockResolvedValue([])
    ;(globalThis.fetch as any).mockResolvedValue({
      json: async () => ({
        data: { Page: { media: [mockTrendingAnime] } },
      }),
    })
  })

  describe('Rendering', () => {
    it('should render dashboard container', async () => {
      render(
        <BrowserRouter>
          <Dashboard />
        </BrowserRouter>
      )

      const dashboard = screen.getByRole('main')
      expect(dashboard).toBeInTheDocument()
      expect(dashboard).toHaveClass('dashboard')
    })

    it('should render all main sections', async () => {
      ;(api.apiListJobs as any).mockResolvedValue(mockJobs)
      ;(api.apiListSubscriptions as any).mockResolvedValue(mockSubscriptions)

      render(
        <BrowserRouter>
          <Dashboard />
        </BrowserRouter>
      )

      await waitFor(() => {
        const elements = screen.getAllByText(/jujutsu kaisen/i)
        expect(elements.length).toBeGreaterThan(0)
      })
    })
  })

  describe('Data Loading', () => {
    it('should fetch jobs on mount', async () => {
      ;(api.apiListJobs as any).mockResolvedValue(mockJobs)
      ;(api.apiListSubscriptions as any).mockResolvedValue(mockSubscriptions)

      render(
        <BrowserRouter>
          <Dashboard />
        </BrowserRouter>
      )

      await waitFor(() => {
        expect(api.apiListJobs).toHaveBeenCalled()
      })
    })

    it('should fetch subscriptions on mount', async () => {
      ;(api.apiListJobs as any).mockResolvedValue(mockJobs)
      ;(api.apiListSubscriptions as any).mockResolvedValue(mockSubscriptions)

      render(
        <BrowserRouter>
          <Dashboard />
        </BrowserRouter>
      )

      await waitFor(() => {
        expect(api.apiListSubscriptions).toHaveBeenCalled()
      })
    })

    it('should handle API errors gracefully', async () => {
      ;(api.apiListJobs as any).mockRejectedValue(new Error('API Error'))

      render(
        <BrowserRouter>
          <Dashboard />
        </BrowserRouter>
      )

      // Dashboard should not crash — stores handle errors internally
      await waitFor(() => {
        const main = screen.getByRole('main')
        expect(main).toBeInTheDocument()
      })
    })
  })

  describe('Empty States', () => {
    it('should show empty state when no jobs', async () => {
      ;(api.apiListJobs as any).mockResolvedValue([])
      ;(api.apiListSubscriptions as any).mockResolvedValue(mockSubscriptions)

      render(
        <BrowserRouter>
          <Dashboard />
        </BrowserRouter>
      )

      await waitFor(() => {
        expect(
          screen.getByText(/aucun téléchargement récent|no recent downloads/i)
        ).toBeInTheDocument()
      })
    })

    it('should show empty state when no subscriptions', async () => {
      ;(api.apiListJobs as any).mockResolvedValue(mockJobs)
      ;(api.apiListSubscriptions as any).mockResolvedValue([])

      render(
        <BrowserRouter>
          <Dashboard />
        </BrowserRouter>
      )

      await waitFor(() => {
        expect(
          screen.getByText(/aucun abonnement actif|no active subscription/i)
        ).toBeInTheDocument()
      })
    })
  })

  describe('Accessibility', () => {
    it('should have proper aria-label', async () => {
      ;(api.apiListJobs as any).mockResolvedValue(mockJobs)
      ;(api.apiListSubscriptions as any).mockResolvedValue(mockSubscriptions)

      render(
        <BrowserRouter>
          <Dashboard />
        </BrowserRouter>
      )

      await waitFor(() => {
        const main = screen.getByRole('main')
        expect(main).toHaveAttribute('aria-label', 'Tableau de bord principal')
      })
    })

    it('should have proper heading hierarchy', async () => {
      ;(api.apiListJobs as any).mockResolvedValue(mockJobs)
      ;(api.apiListSubscriptions as any).mockResolvedValue(mockSubscriptions)

      render(
        <BrowserRouter>
          <Dashboard />
        </BrowserRouter>
      )

      await waitFor(() => {
        const headings = screen.getAllByRole('heading')
        expect(headings.length).toBeGreaterThan(0)
      })
    })
  })
})
