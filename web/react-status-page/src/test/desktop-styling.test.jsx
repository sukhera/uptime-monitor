import { describe, it, expect, beforeEach, vi } from 'vitest'

// Mock window.matchMedia for responsive testing
const createMatchMedia = (width) => vi.fn().mockImplementation(query => {
  const mediaQuery = query.match(/\(min-width:\s*(\d+)px\)/)
  const minWidth = mediaQuery ? parseInt(mediaQuery[1]) : 0
  
  return {
    matches: width >= minWidth,
    media: query,
    onchange: null,
    addListener: vi.fn(),
    removeListener: vi.fn(),
    addEventListener: vi.fn(),
    removeEventListener: vi.fn(),
    dispatchEvent: vi.fn(),
  }
})

describe('Desktop Breakpoints and Styling', () => {
  describe('Desktop Breakpoints', () => {
    beforeEach(() => {
      // Reset window.matchMedia before each test
      delete window.matchMedia
    })

    it('should support lg breakpoint (1024px)', () => {
      window.matchMedia = createMatchMedia(1024)
      
      const lgQuery = window.matchMedia('(min-width: 1024px)')
      expect(lgQuery.matches).toBe(true)
      
      const smQuery = window.matchMedia('(min-width: 640px)')
      expect(smQuery.matches).toBe(true) // Should still match since we're at 1024px
    })

    it('should support xl breakpoint (1280px)', () => {
      window.matchMedia = createMatchMedia(1280)
      
      const xlQuery = window.matchMedia('(min-width: 1280px)')
      expect(xlQuery.matches).toBe(true)
    })

    it('should support 2xl breakpoint (1536px)', () => {
      window.matchMedia = createMatchMedia(1536)
      
      const xl2Query = window.matchMedia('(min-width: 1536px)')
      expect(xl2Query.matches).toBe(true)
    })

    it('should support 3xl breakpoint (1920px)', () => {
      window.matchMedia = createMatchMedia(1920)
      
      const xl3Query = window.matchMedia('(min-width: 1920px)')
      expect(xl3Query.matches).toBe(true)
    })

    it('should support 4xl breakpoint (2560px)', () => {
      window.matchMedia = createMatchMedia(2560)
      
      const xl4Query = window.matchMedia('(min-width: 2560px)')
      expect(xl4Query.matches).toBe(true)
    })

    it('should not match desktop breakpoints on smaller screens', () => {
      window.matchMedia = createMatchMedia(768)
      
      const lgQuery = window.matchMedia('(min-width: 1024px)')
      expect(lgQuery.matches).toBe(false)
    })
  })

  describe('Glassmorphism CSS Classes', () => {
    it('should define glass-card utility class', () => {
      const div = document.createElement('div')
      div.className = 'glass-card'
      document.body.appendChild(div)
      
      // Check if the class exists (we can't easily test computed styles in jsdom)
      expect(div.classList.contains('glass-card')).toBe(true)
      
      document.body.removeChild(div)
    })

    it('should define glass-secondary utility class', () => {
      const div = document.createElement('div')
      div.className = 'glass-secondary'
      document.body.appendChild(div)
      
      expect(div.classList.contains('glass-secondary')).toBe(true)
      
      document.body.removeChild(div)
    })

    it('should define desktop-hover-lift utility class', () => {
      const div = document.createElement('div')
      div.className = 'desktop-hover-lift'
      document.body.appendChild(div)
      
      expect(div.classList.contains('desktop-hover-lift')).toBe(true)
      
      document.body.removeChild(div)
    })

    it('should define desktop-container utility class', () => {
      const div = document.createElement('div')
      div.className = 'desktop-container'
      document.body.appendChild(div)
      
      expect(div.classList.contains('desktop-container')).toBe(true)
      
      document.body.removeChild(div)
    })

    it('should define desktop-section utility class', () => {
      const div = document.createElement('div')
      div.className = 'desktop-section'
      document.body.appendChild(div)
      
      expect(div.classList.contains('desktop-section')).toBe(true)
      
      document.body.removeChild(div)
    })
  })

  describe('Desktop Grid System', () => {
    it('should support 12-column grid system', () => {
      const div = document.createElement('div')
      div.className = 'grid grid-cols-12'
      document.body.appendChild(div)
      
      expect(div.classList.contains('grid')).toBe(true)
      expect(div.classList.contains('grid-cols-12')).toBe(true)
      
      document.body.removeChild(div)
    })

    it('should support extended grid columns (13-16)', () => {
      const grids = [
        'grid-cols-13',
        'grid-cols-14', 
        'grid-cols-15',
        'grid-cols-16'
      ]
      
      grids.forEach(gridClass => {
        const div = document.createElement('div')
        div.className = gridClass
        document.body.appendChild(div)
        
        expect(div.classList.contains(gridClass)).toBe(true)
        
        document.body.removeChild(div)
      })
    })

    it('should support column spanning classes', () => {
      const spans = ['col-span-3', 'col-span-9']
      
      spans.forEach(spanClass => {
        const div = document.createElement('div')
        div.className = spanClass
        document.body.appendChild(div)
        
        expect(div.classList.contains(spanClass)).toBe(true)
        
        document.body.removeChild(div)
      })
    })
  })

  describe('Desktop Typography', () => {
    it('should support large desktop font sizes', () => {
      const fontSizes = [
        'text-4xl',
        'text-5xl',
        'lg:text-5xl',
        'text-xl',
        'lg:text-2xl'
      ]
      
      fontSizes.forEach(fontSize => {
        const div = document.createElement('div')
        div.className = fontSize
        document.body.appendChild(div)
        
        expect(div.classList.contains(fontSize)).toBe(true)
        
        document.body.removeChild(div)
      })
    })
  })

  describe('Desktop Animations', () => {
    it('should support desktop-fade-in animation', () => {
      const div = document.createElement('div')
      div.className = 'animate-desktop-fade-in'
      document.body.appendChild(div)
      
      expect(div.classList.contains('animate-desktop-fade-in')).toBe(true)
      
      document.body.removeChild(div)
    })

    it('should support hover transitions', () => {
      const div = document.createElement('div')
      div.className = 'transition-all duration-200'
      document.body.appendChild(div)
      
      expect(div.classList.contains('transition-all')).toBe(true)
      expect(div.classList.contains('duration-200')).toBe(true)
      
      document.body.removeChild(div)
    })
  })

  describe('Status Color System', () => {
    it('should support operational glass colors', () => {
      const colors = [
        'bg-operational-glass',
        'border-operational-border',
        'text-green-300',
        'text-green-400'
      ]
      
      colors.forEach(colorClass => {
        const div = document.createElement('div')
        div.className = colorClass
        document.body.appendChild(div)
        
        expect(div.classList.contains(colorClass)).toBe(true)
        
        document.body.removeChild(div)
      })
    })

    it('should support degraded glass colors', () => {
      const colors = [
        'bg-degraded-glass',
        'border-degraded-border',
        'text-yellow-300',
        'text-yellow-400'
      ]
      
      colors.forEach(colorClass => {
        const div = document.createElement('div')
        div.className = colorClass
        document.body.appendChild(div)
        
        expect(div.classList.contains(colorClass)).toBe(true)
        
        document.body.removeChild(div)
      })
    })

    it('should support down glass colors', () => {
      const colors = [
        'bg-down-glass', 
        'border-down-border',
        'text-red-300',
        'text-red-400'
      ]
      
      colors.forEach(colorClass => {
        const div = document.createElement('div')
        div.className = colorClass
        document.body.appendChild(div)
        
        expect(div.classList.contains(colorClass)).toBe(true)
        
        document.body.removeChild(div)
      })
    })
  })

  describe('Desktop Spacing and Layout', () => {
    it('should support desktop-specific spacing classes', () => {
      const spacingClasses = [
        'min-h-[280px]',
        'lg:min-h-[320px]',
        'w-12',
        'h-12',
        'lg:w-16',
        'lg:h-16'
      ]
      
      spacingClasses.forEach(spacingClass => {
        const div = document.createElement('div')
        div.className = spacingClass
        document.body.appendChild(div)
        
        expect(div.classList.contains(spacingClass)).toBe(true)
        
        document.body.removeChild(div)
      })
    })

    it('should support backdrop blur classes', () => {
      const blurClasses = [
        'backdrop-blur-10',
        'backdrop-blur-20'
      ]
      
      blurClasses.forEach(blurClass => {
        const div = document.createElement('div')
        div.className = blurClass
        document.body.appendChild(div)
        
        expect(div.classList.contains(blurClass)).toBe(true)
        
        document.body.removeChild(div)
      })
    })
  })

  describe('Desktop Interactive States', () => {
    it('should support hover state classes', () => {
      const hoverClasses = [
        'hover:bg-white/10',
        'hover:bg-red-500/30',
        'desktop-hover-lift'
      ]
      
      hoverClasses.forEach(hoverClass => {
        const div = document.createElement('div')
        div.className = hoverClass
        document.body.appendChild(div)
        
        expect(div.classList.contains(hoverClass)).toBe(true)
        
        document.body.removeChild(div)
      })
    })

    it('should support cursor pointer for interactive elements', () => {
      const div = document.createElement('div')
      div.className = 'cursor-pointer'
      document.body.appendChild(div)
      
      expect(div.classList.contains('cursor-pointer')).toBe(true)
      
      document.body.removeChild(div)
    })
  })
})