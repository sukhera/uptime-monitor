import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen } from '@testing-library/react'
import StatusDashboard from './StatusDashboard'

// Mock child components
vi.mock('./ServiceCard', () => ({
  default: ({ service }) => (
    <div data-testid="service-card">
      {service.name} - {service.status}
    </div>
  )
}))

vi.mock('../Desktop/DesktopSidebar', () => ({
  default: ({ services }) => (
    <div data-testid="desktop-sidebar">
      Sidebar with {services.length} services
    </div>
  )
}))

vi.mock('../common/LoadingSpinner', () => ({
  default: () => <div data-testid="loading-spinner">Loading...</div>
}))

// Mock hooks
vi.mock('../../hooks/useApi', () => ({
  useApi: vi.fn()
}))

vi.mock('../../hooks/usePolling', () => ({
  usePolling: vi.fn()
}))

import { useApi } from '../../hooks/useApi'
import { usePolling } from '../../hooks/usePolling'

describe('StatusDashboard Desktop Layout', () => {
  const mockServices = [
    {
      name: 'API Service',
      status: 'operational',
      latency_ms: 120,
      category: 'Backend'
    },
    {
      name: 'Web Service', 
      status: 'degraded',
      latency_ms: 250,
      category: 'Frontend'
    },
    {
      name: 'Database',
      status: 'down',
      latency_ms: null,
      category: 'Infrastructure'
    }
  ]

  beforeEach(() => {
    usePolling.mockImplementation(() => {})
  })

  describe('Desktop-First Layout Structure', () => {
    it('should apply desktop gradient background', () => {
      useApi.mockReturnValue({
        data: mockServices,
        loading: false,
        error: null,
        refetch: vi.fn()
      })

      const { container } = render(<StatusDashboard />)
      const mainContainer = container.firstChild
      
      expect(mainContainer).toHaveClass(
        'min-h-screen',
        'bg-gradient-to-br',
        'from-slate-900',
        'via-purple-900',
        'to-slate-900'
      )
    })

    it('should implement 12-column grid system', () => {
      useApi.mockReturnValue({
        data: mockServices,
        loading: false,
        error: null,
        refetch: vi.fn()
      })

      render(<StatusDashboard />)
      
      const gridContainer = screen.getByText('Sidebar with 3 services').parentElement.parentElement
      expect(gridContainer).toHaveClass('grid', 'grid-cols-12', 'gap-8')
    })

    it('should have 3-column sidebar and 9-column main content', () => {
      useApi.mockReturnValue({
        data: mockServices,
        loading: false,
        error: null,
        refetch: vi.fn()
      })

      render(<StatusDashboard />)
      
      const sidebar = screen.getByTestId('desktop-sidebar').parentElement
      const mainContent = screen.getByText('System Status').parentElement.parentElement
      
      expect(sidebar).toHaveClass('col-span-3')
      expect(mainContent).toHaveClass('col-span-9')
    })

    it('should use desktop container and section classes', () => {
      useApi.mockReturnValue({
        data: mockServices,
        loading: false,
        error: null,
        refetch: vi.fn()
      })

      const { container } = render(<StatusDashboard />)
      const desktopContainer = container.querySelector('.desktop-container')
      
      expect(desktopContainer).toHaveClass('desktop-container', 'desktop-section')
    })
  })

  describe('Desktop Typography and Styling', () => {
    it('should use large desktop typography for main heading', () => {
      useApi.mockReturnValue({
        data: mockServices,
        loading: false,
        error: null,
        refetch: vi.fn()
      })

      render(<StatusDashboard />)
      const heading = screen.getByText('System Status')
      
      expect(heading).toHaveClass('text-4xl', 'lg:text-5xl', 'font-bold', 'text-white', 'mb-4')
    })

    it('should display real-time status indicator', () => {
      useApi.mockReturnValue({
        data: mockServices,
        loading: false,
        error: null,
        refetch: vi.fn()
      })

      render(<StatusDashboard />)
      
      expect(screen.getByText('Real-time monitoring of all services')).toBeInTheDocument()
      expect(screen.getByText(/Last updated:/)).toBeInTheDocument()
      
      // Check for pulsing indicator
      const { container } = render(<StatusDashboard />)
      const pulseElement = container.querySelector('.animate-pulse')
      expect(pulseElement).toHaveClass('bg-green-400', 'rounded-full')
    })

    it('should apply desktop fade-in animation', () => {
      useApi.mockReturnValue({
        data: mockServices,
        loading: false,
        error: null,
        refetch: vi.fn()
      })

      render(<StatusDashboard />)
      const headerSection = screen.getByText('System Status').parentElement
      
      expect(headerSection).toHaveClass('animate-desktop-fade-in')
    })
  })

  describe('Desktop Services Grid', () => {
    it('should use desktop grid layout for services', () => {
      useApi.mockReturnValue({
        data: mockServices,
        loading: false,
        error: null,
        refetch: vi.fn()
      })

      render(<StatusDashboard />)
      const servicesGrid = screen.getAllByTestId('service-card')[0].parentElement
      
      expect(servicesGrid).toHaveClass('grid', 'gap-6', 'lg:grid-cols-2', 'xl:grid-cols-3')
    })

    it('should render all service cards', () => {
      useApi.mockReturnValue({
        data: mockServices,
        loading: false,
        error: null,
        refetch: vi.fn()
      })

      render(<StatusDashboard />)
      
      expect(screen.getByText('API Service - operational')).toBeInTheDocument()
      expect(screen.getByText('Web Service - degraded')).toBeInTheDocument()
      expect(screen.getByText('Database - down')).toBeInTheDocument()
    })

    it('should pass services data to sidebar', () => {
      useApi.mockReturnValue({
        data: mockServices,
        loading: false,
        error: null,
        refetch: vi.fn()
      })

      render(<StatusDashboard />)
      
      expect(screen.getByText('Sidebar with 3 services')).toBeInTheDocument()
    })
  })

  describe('Desktop Error States', () => {
    it('should display glassmorphism error state', () => {
      useApi.mockReturnValue({
        data: null,
        loading: false,
        error: { message: 'Connection failed' },
        refetch: vi.fn()
      })

      render(<StatusDashboard />)
      
      expect(screen.getByText('⚠️')).toBeInTheDocument()
      expect(screen.getByText('Error Loading Status')).toBeInTheDocument()
      expect(screen.getByText('Connection failed')).toBeInTheDocument()
      
      // Check for glassmorphism styling
      const errorContainer = screen.getByText('Error Loading Status').parentElement
      expect(errorContainer).toHaveClass('glass-card', 'rounded-2xl')
    })

    it('should have desktop-styled retry button', () => {
      const mockRefetch = vi.fn()
      useApi.mockReturnValue({
        data: null,
        loading: false,
        error: { message: 'Connection failed' },
        refetch: mockRefetch
      })

      render(<StatusDashboard />)
      const retryButton = screen.getByText('🔄 Retry Connection')
      
      expect(retryButton).toHaveClass('desktop-hover-lift')
      expect(retryButton).toHaveClass('bg-red-500/20', 'hover:bg-red-500/30')
    })
  })

  describe('Desktop Empty State', () => {
    it('should display glassmorphism empty state', () => {
      useApi.mockReturnValue({
        data: [],
        loading: false,
        error: null,
        refetch: vi.fn()
      })

      render(<StatusDashboard />)
      
      expect(screen.getByText('📡')).toBeInTheDocument()
      expect(screen.getByText('No Services Available')).toBeInTheDocument()
      expect(screen.getByText(/Configure your first service/)).toBeInTheDocument()
      
      // Check for glassmorphism styling
      const emptyContainer = screen.getByText('No Services Available').parentElement
      expect(emptyContainer).toHaveClass('glass-card', 'rounded-2xl')
    })

    it('should have desktop-styled configure button', () => {
      useApi.mockReturnValue({
        data: [],
        loading: false,
        error: null,
        refetch: vi.fn()
      })

      render(<StatusDashboard />)
      const configButton = screen.getByText('⚙️ Configure Services')
      
      expect(configButton).toHaveClass('desktop-hover-lift')
      expect(configButton).toHaveClass('glass-secondary', 'hover:bg-white/10')
    })
  })

  describe('Desktop Polling Integration', () => {
    it('should setup polling with 60 second interval', () => {
      const mockRefetch = vi.fn()
      useApi.mockReturnValue({
        data: mockServices,
        loading: false,
        error: null,
        refetch: mockRefetch
      })

      render(<StatusDashboard />)
      
      expect(usePolling).toHaveBeenCalledWith(expect.any(Function), 60000)
    })

    it('should update timestamp on polling', async () => {
      const mockRefetch = vi.fn()
      let pollingCallback
      
      usePolling.mockImplementation((callback) => {
        pollingCallback = callback
      })
      
      useApi.mockReturnValue({
        data: mockServices,
        loading: false,
        error: null,
        refetch: mockRefetch
      })

      const { act } = await import('@testing-library/react')
      render(<StatusDashboard />)
      
      // Simulate polling callback with act
      await act(async () => {
        pollingCallback()
      })
      
      expect(mockRefetch).toHaveBeenCalled()
    })
  })

  describe('Loading State', () => {
    it('should show loading spinner', () => {
      useApi.mockReturnValue({
        data: null,
        loading: true,
        error: null,
        refetch: vi.fn()
      })

      render(<StatusDashboard />)
      
      expect(screen.getByTestId('loading-spinner')).toBeInTheDocument()
    })
  })
})