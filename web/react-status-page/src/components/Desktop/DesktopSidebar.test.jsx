import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/react'
import DesktopSidebar from './DesktopSidebar'

describe('DesktopSidebar Component', () => {
  const mockServices = [
    {
      name: 'API Gateway',
      status: 'operational',
      category: 'Backend'
    },
    {
      name: 'Web Service',
      status: 'operational', 
      category: 'Frontend'
    },
    {
      name: 'Database Primary',
      status: 'degraded',
      category: 'Infrastructure'
    },
    {
      name: 'Database Replica',
      status: 'down',
      category: 'Infrastructure'
    },
    {
      name: 'Uncategorized Service',
      status: 'operational'
    }
  ]

  describe('Desktop Styling and Layout', () => {
    it('should apply sticky positioning and spacing', () => {
      const { container } = render(<DesktopSidebar services={mockServices} />)
      const sidebar = container.firstChild
      
      expect(sidebar).toHaveClass('sticky', 'top-8', 'space-y-6')
    })

    it('should apply glassmorphism styling to all sections', () => {
      const { container } = render(<DesktopSidebar services={mockServices} />)
      const glassCards = container.querySelectorAll('.glass-card')
      
      expect(glassCards.length).toBeGreaterThanOrEqual(3) // Overview + Categories + Quick Actions
      glassCards.forEach(card => {
        expect(card).toHaveClass('rounded-2xl', 'p-6', 'animate-desktop-fade-in')
      })
    })
  })

  describe('System Overview Section', () => {
    it('should display correct overall status for mixed services', () => {
      render(<DesktopSidebar services={mockServices} />)
      
      expect(screen.getByText('Issues Detected')).toBeInTheDocument()
      expect(screen.getByText('3/5 services operational')).toBeInTheDocument()
    })

    it('should display all systems operational when no issues', () => {
      const allOperationalServices = mockServices.map(s => ({...s, status: 'operational'}))
      render(<DesktopSidebar services={allOperationalServices} />)
      
      expect(screen.getByText('All Systems Operational')).toBeInTheDocument()
      expect(screen.getByText('5/5 services operational')).toBeInTheDocument()
    })

    it('should display degraded performance when only degraded services', () => {
      const degradedServices = [
        { name: 'Service 1', status: 'operational' },
        { name: 'Service 2', status: 'degraded' }
      ]
      render(<DesktopSidebar services={degradedServices} />)
      
      expect(screen.getByText('Degraded Performance')).toBeInTheDocument()
    })

    it('should apply correct color classes for overall status', () => {
      render(<DesktopSidebar services={mockServices} />)
      const statusText = screen.getByText('Issues Detected')
      
      expect(statusText).toHaveClass('text-red-400')
    })

    it('should display service statistics with glassmorphism styling', () => {
      const { container } = render(<DesktopSidebar services={mockServices} />)
      
      expect(screen.getByText('Total Services')).toBeInTheDocument()
      expect(screen.getByText('5')).toBeInTheDocument()
      expect(screen.getByText('Operational')).toBeInTheDocument()
      expect(screen.getByText('3')).toBeInTheDocument()
      
      // Check for glass-secondary styling on stat items
      const statItems = container.querySelectorAll('.glass-secondary')
      expect(statItems.length).toBeGreaterThan(0)
    })

    it('should show degraded count when present', () => {
      const { container } = render(<DesktopSidebar services={mockServices} />)
      
      expect(screen.getByText('Degraded')).toBeInTheDocument()
      const degradedCount = container.querySelector('.text-yellow-400')
      expect(degradedCount).toHaveTextContent('1')
    })

    it('should show down count when present', () => {
      const { container } = render(<DesktopSidebar services={mockServices} />)
      
      expect(screen.getByText('Down')).toBeInTheDocument()
      const downCount = container.querySelector('.text-red-400.font-medium')
      expect(downCount).toHaveTextContent('1')
    })
  })

  describe('Categories Section', () => {
    it('should display categories when services have categories', () => {
      render(<DesktopSidebar services={mockServices} />)
      
      expect(screen.getByText('Categories')).toBeInTheDocument()
      expect(screen.getByText('Backend')).toBeInTheDocument()
      expect(screen.getByText('Frontend')).toBeInTheDocument()
      expect(screen.getByText('Infrastructure')).toBeInTheDocument()
    })

    it('should display service count per category', () => {
      render(<DesktopSidebar services={mockServices} />)
      
      // Infrastructure has 2 services
      expect(screen.getByText('2 services')).toBeInTheDocument()
      
      // Backend and Frontend have 1 service each
      expect(screen.getAllByText('1 services')).toHaveLength(2)
    })

    it('should show correct status indicators for categories', () => {
      const { container } = render(<DesktopSidebar services={mockServices} />)
      
      // Backend category (all operational) should be green
      const backendIndicator = container.querySelector('.text-lg.text-green-400')
      expect(backendIndicator).toBeInTheDocument()
      
      // Infrastructure category (has down service) should be red
      const infraIndicator = container.querySelector('.text-lg.text-red-400')
      expect(infraIndicator).toBeInTheDocument()
    })

    it('should apply hover effects to category items', () => {
      const { container } = render(<DesktopSidebar services={mockServices} />)
      const categoryItems = container.querySelectorAll('.hover\\:bg-white\\/5')
      
      expect(categoryItems.length).toBeGreaterThan(0)
      categoryItems.forEach(item => {
        expect(item).toHaveClass('transition-colors', 'cursor-pointer')
      })
    })

    it('should not display categories section when no categories exist', () => {
      const servicesWithoutCategories = [
        { name: 'Service 1', status: 'operational' },
        { name: 'Service 2', status: 'operational' }
      ]
      render(<DesktopSidebar services={servicesWithoutCategories} />)
      
      expect(screen.queryByText('Categories')).not.toBeInTheDocument()
    })
  })

  describe('Quick Actions Section', () => {
    it('should display all quick action buttons', () => {
      render(<DesktopSidebar services={mockServices} />)
      
      expect(screen.getByText('Quick Actions')).toBeInTheDocument()
      expect(screen.getByText('📊 View All Metrics')).toBeInTheDocument()
      expect(screen.getByText('🔄 Refresh All Services')).toBeInTheDocument()
      expect(screen.getByText('⚙️ Service Configuration')).toBeInTheDocument()
    })

    it('should apply desktop hover styling to action buttons', () => {
      const { container } = render(<DesktopSidebar services={mockServices} />)
      const actionButtons = container.querySelectorAll('.hover\\:bg-white\\/10')
      
      expect(actionButtons.length).toBe(3)
      actionButtons.forEach(button => {
        expect(button).toHaveClass('glass-secondary', 'transition-colors')
      })
    })
  })

  describe('Empty State Handling', () => {
    it('should handle empty services array', () => {
      const { container } = render(<DesktopSidebar services={[]} />)
      
      expect(screen.getByText('All Systems Operational')).toBeInTheDocument()
      expect(screen.getByText('0/0 services operational')).toBeInTheDocument()
      
      // Check total services count specifically
      const totalServicesCount = container.querySelector('.text-white.font-medium')
      expect(totalServicesCount).toHaveTextContent('0')
    })

    it('should handle undefined services', () => {
      render(<DesktopSidebar services={undefined} />)
      
      expect(screen.getByText('All Systems Operational')).toBeInTheDocument()
      expect(screen.getByText('0/0 services operational')).toBeInTheDocument()
    })

    it('should handle null services', () => {
      render(<DesktopSidebar services={null} />)
      
      expect(screen.getByText('All Systems Operational')).toBeInTheDocument()
    })
  })

  describe('Status Calculations', () => {
    it('should correctly calculate service counts', () => {
      const { container } = render(<DesktopSidebar services={mockServices} />)
      
      // 5 total, 3 operational, 1 degraded, 1 down
      expect(screen.getByText('5')).toBeInTheDocument() // Total
      expect(screen.getByText('3')).toBeInTheDocument() // Operational  
      
      // Check specific color classes for counts
      const degradedCount = container.querySelector('.text-yellow-400.font-medium')
      const downCount = container.querySelector('.text-red-400.font-medium')
      expect(degradedCount).toHaveTextContent('1')
      expect(downCount).toHaveTextContent('1')
    })

    it('should prioritize down status over degraded in overall status', () => {
      const mixedServices = [
        { name: 'Service 1', status: 'operational' },
        { name: 'Service 2', status: 'degraded' },
        { name: 'Service 3', status: 'down' }
      ]
      render(<DesktopSidebar services={mixedServices} />)
      
      // Should show "Issues Detected" (down status) rather than "Degraded Performance"
      expect(screen.getByText('Issues Detected')).toBeInTheDocument()
      const statusText = screen.getByText('Issues Detected')
      expect(statusText).toHaveClass('text-red-400')
    })
  })

  describe('Accessibility', () => {
    it('should have proper heading hierarchy', () => {
      render(<DesktopSidebar services={mockServices} />)
      
      const headings = screen.getAllByRole('heading')
      expect(headings.length).toBeGreaterThanOrEqual(3) // Overview, Categories, Quick Actions
    })

    it('should have semantic button elements for actions', () => {
      render(<DesktopSidebar services={mockServices} />)
      
      const buttons = screen.getAllByRole('button')
      expect(buttons.length).toBe(3) // Three quick action buttons
    })
  })
})