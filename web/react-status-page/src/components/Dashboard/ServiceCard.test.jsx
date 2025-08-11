import { describe, it, expect, vi } from 'vitest'
import { render, screen, fireEvent } from '@testing-library/react'
import ServiceCard from './ServiceCard'

// Mock StatusIndicator component
vi.mock('./StatusIndicator', () => ({
  default: ({ status }) => <div data-testid="status-indicator">{status}</div>
}))

describe('ServiceCard Desktop Features', () => {
  const mockService = {
    name: 'Test Service',
    status: 'operational',
    latency_ms: 125,
    updated_at: '2024-01-01T10:00:00Z',
    category: 'Core',
    icon: '🌐',
    uptime: '99.9',
    error: null
  }

  describe('Desktop-First Styling', () => {
    it('should apply glassmorphism styling classes', () => {
      const { container } = render(<ServiceCard service={mockService} />)
      const card = container.firstChild
      
      expect(card).toHaveClass('glass-card')
      expect(card).toHaveClass('rounded-2xl')
      expect(card).toHaveClass('desktop-hover-lift')
      expect(card).toHaveClass('animate-desktop-fade-in')
    })

    it('should have desktop-specific minimum heights', () => {
      const { container } = render(<ServiceCard service={mockService} />)
      const card = container.firstChild
      
      expect(card).toHaveClass('min-h-[280px]')
      expect(card).toHaveClass('lg:min-h-[320px]')
    })

    it('should display larger icons for desktop', () => {
      render(<ServiceCard service={mockService} />)
      const iconContainer = screen.getByText('🌐')
      
      expect(iconContainer).toHaveClass('w-12', 'h-12', 'lg:w-16', 'lg:h-16')
      expect(iconContainer).toHaveClass('glass-card')
    })

    it('should use desktop typography scales', () => {
      render(<ServiceCard service={mockService} />)
      const serviceName = screen.getByText('Test Service')
      
      expect(serviceName).toHaveClass('text-xl', 'lg:text-2xl', 'font-semibold', 'text-white')
    })
  })

  describe('Status-Specific Glassmorphism', () => {
    it('should apply operational glassmorphism colors', () => {
      const { container } = render(<ServiceCard service={{...mockService, status: 'operational'}} />)
      const statusIndicator = container.querySelector('.bg-operational-glass')
      
      expect(statusIndicator).toHaveClass('bg-operational-glass')
      expect(statusIndicator).toHaveClass('border-operational-border')
      expect(statusIndicator).toHaveClass('text-green-300')
    })

    it('should apply degraded glassmorphism colors', () => {
      const { container } = render(<ServiceCard service={{...mockService, status: 'degraded'}} />)
      const statusIndicator = container.querySelector('.bg-degraded-glass')
      
      expect(statusIndicator).toHaveClass('bg-degraded-glass')
      expect(statusIndicator).toHaveClass('border-degraded-border')
      expect(statusIndicator).toHaveClass('text-yellow-300')
    })

    it('should apply down glassmorphism colors', () => {
      const { container } = render(<ServiceCard service={{...mockService, status: 'down'}} />)
      const statusIndicator = container.querySelector('.bg-down-glass')
      
      expect(statusIndicator).toHaveClass('bg-down-glass')
      expect(statusIndicator).toHaveClass('border-down-border')
      expect(statusIndicator).toHaveClass('text-red-300')
    })
  })

  describe('Desktop Metrics Display', () => {
    it('should display uptime and response time in grid layout', () => {
      render(<ServiceCard service={mockService} />)
      
      expect(screen.getByText('99.9%')).toBeInTheDocument()
      expect(screen.getByText('125ms')).toBeInTheDocument()
      expect(screen.getByText('Uptime')).toBeInTheDocument()
      expect(screen.getByText('Response')).toBeInTheDocument()
    })

    it('should apply glass-secondary styling to metrics containers', () => {
      const { container } = render(<ServiceCard service={mockService} />)
      const metricsGrid = container.querySelector('.grid.grid-cols-2.gap-4')
      const metricContainers = metricsGrid.querySelectorAll('.glass-secondary')
      
      expect(metricContainers).toHaveLength(2)
    })
  })

  describe('Desktop Hover Interactions', () => {
    it('should show hover overlay on mouse enter', () => {
      const { container } = render(<ServiceCard service={mockService} />)
      const card = container.firstChild
      
      // Initially no overlay
      expect(screen.queryByText('View Details')).not.toBeInTheDocument()
      
      // Hover should show overlay
      fireEvent.mouseEnter(card)
      expect(screen.getByText('View Details')).toBeInTheDocument()
      expect(screen.getByText('Click to see full service metrics')).toBeInTheDocument()
    })

    it('should hide hover overlay on mouse leave', () => {
      const { container } = render(<ServiceCard service={mockService} />)
      const card = container.firstChild
      
      fireEvent.mouseEnter(card)
      expect(screen.getByText('View Details')).toBeInTheDocument()
      
      fireEvent.mouseLeave(card)
      expect(screen.queryByText('View Details')).not.toBeInTheDocument()
    })
  })

  describe('Error Handling', () => {
    it('should display error details with glassmorphism styling', () => {
      const serviceWithError = {
        ...mockService,
        error: 'Connection timeout'
      }
      
      render(<ServiceCard service={serviceWithError} />)
      
      expect(screen.getByText('Error Details:')).toBeInTheDocument()
      expect(screen.getByText('Connection timeout')).toBeInTheDocument()
      
      const errorContainer = screen.getByText('Error Details:').parentElement
      expect(errorContainer).toHaveClass('bg-down-glass', 'border-down-border')
    })

    it('should handle invalid service data with glassmorphism', () => {
      const { container } = render(<ServiceCard service={null} />)
      const errorCard = container.firstChild
      
      expect(errorCard).toHaveClass('glass-card')
      expect(screen.getByText('Invalid service data')).toBeInTheDocument()
      expect(screen.getByText('⚠️')).toBeInTheDocument()
    })
  })

  describe('Desktop Category Display', () => {
    it('should display service category', () => {
      render(<ServiceCard service={mockService} />)
      expect(screen.getByText('Core')).toBeInTheDocument()
    })

    it('should show default category when none provided', () => {
      const serviceWithoutCategory = { ...mockService, category: undefined }
      render(<ServiceCard service={serviceWithoutCategory} />)
      expect(screen.getByText('Service')).toBeInTheDocument()
    })
  })

  describe('Desktop Data Formatting', () => {
    it('should format timestamps correctly', () => {
      render(<ServiceCard service={mockService} />)
      // Check that timestamp is formatted (actual time depends on locale)
      expect(screen.getByText(/\d{1,2}:\d{2}:\d{2}/)).toBeInTheDocument()
    })

    it('should handle missing data gracefully', () => {
      const incompleteService = {
        name: 'Incomplete Service',
        status: 'operational'
      }
      
      render(<ServiceCard service={incompleteService} />)
      
      expect(screen.getAllByText('N/A')).toHaveLength(2) // Response time and timestamp
      expect(screen.getByText('99.9%')).toBeInTheDocument() // Default uptime
    })
  })

  describe('Accessibility', () => {
    it('should have proper cursor pointer for interactive elements', () => {
      const { container } = render(<ServiceCard service={mockService} />)
      const card = container.firstChild
      
      expect(card).toHaveClass('cursor-pointer')
    })

    it('should maintain semantic structure', () => {
      render(<ServiceCard service={mockService} />)
      
      expect(screen.getByRole('heading')).toHaveTextContent('Test Service')
    })
  })
})