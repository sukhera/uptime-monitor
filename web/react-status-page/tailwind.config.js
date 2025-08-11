/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  darkMode: 'class',
  theme: {
    screens: {
      // Desktop-only breakpoints (removed sm: 640px, md: 768px)
      'lg': '1024px',   // Minimum desktop
      'xl': '1280px',   // Standard desktop  
      '2xl': '1536px',  // Large desktop
      '3xl': '1920px',  // Ultra-wide
      '4xl': '2560px',  // 4K displays
    },
    extend: {
      colors: {
        // Enhanced status colors with glassmorphism variants
        operational: {
          50: '#f0fdf4',
          500: '#22c55e',
          600: '#16a34a',
          900: '#14532d',
          glass: 'rgba(34, 197, 94, 0.1)',
          border: 'rgba(34, 197, 94, 0.2)',
        },
        degraded: {
          50: '#fffbeb',
          500: '#f59e0b',
          600: '#d97706',
          900: '#78350f',
          glass: 'rgba(245, 158, 11, 0.1)',
          border: 'rgba(245, 158, 11, 0.2)',
        },
        down: {
          50: '#fef2f2',
          500: '#ef4444',
          600: '#dc2626',
          900: '#7f1d1d',
          glass: 'rgba(239, 68, 68, 0.1)',
          border: 'rgba(239, 68, 68, 0.2)',
        },
        // Glassmorphism color system
        glass: {
          primary: 'rgba(255, 255, 255, 0.1)',
          secondary: 'rgba(255, 255, 255, 0.05)',
          border: 'rgba(255, 255, 255, 0.2)',
        },
      },
      // Desktop-optimized grid columns
      gridTemplateColumns: {
        '13': 'repeat(13, minmax(0, 1fr))',
        '14': 'repeat(14, minmax(0, 1fr))',
        '15': 'repeat(15, minmax(0, 1fr))',
        '16': 'repeat(16, minmax(0, 1fr))',
      },
      // Desktop typography scale
      fontSize: {
        '5xl': ['3rem', { lineHeight: '1' }],         // 48px - Section headings
        '4xl': ['2.25rem', { lineHeight: '2.5rem' }], // 36px - Card titles
        '3xl': ['1.875rem', { lineHeight: '2.25rem' }], // 30px - Service names
      },
      // Enhanced backdrop blur for glassmorphism
      backdropBlur: {
        xs: '2px',
        sm: '4px', 
        DEFAULT: '8px',
        md: '12px',
        lg: '16px',
        xl: '24px',
        '2xl': '40px',
      },
      animation: {
        'pulse-slow': 'pulse 3s infinite',
        'bounce-subtle': 'bounce 2s infinite',
        'desktop-fade-in': 'desktop-fade-in 0.5s ease-out',
      },
      keyframes: {
        'desktop-fade-in': {
          from: {
            opacity: '0',
            transform: 'translateY(20px)',
          },
          to: {
            opacity: '1',
            transform: 'translateY(0)',
          },
        },
      },
    },
  },
  plugins: [
    // Glassmorphism utilities plugin
    function({ addUtilities }) {
      const newUtilities = {
        '.glass-card': {
          background: 'rgba(255, 255, 255, 0.1)',
          backdropFilter: 'blur(20px)',
          border: '1px solid rgba(255, 255, 255, 0.2)',
          boxShadow: '0 8px 32px 0 rgba(31, 38, 135, 0.37)',
        },
        '.glass-secondary': {
          background: 'rgba(255, 255, 255, 0.05)',
          backdropFilter: 'blur(10px)',
          border: '1px solid rgba(255, 255, 255, 0.1)',
        },
        '.desktop-hover-lift': {
          transition: 'all 0.3s cubic-bezier(0.4, 0, 0.2, 1)',
        },
        '.desktop-hover-lift:hover': {
          transform: 'translateY(-8px) scale(1.02)',
          boxShadow: '0 20px 25px -5px rgba(0, 0, 0, 0.1)',
        },
        '.desktop-container': {
          maxWidth: '80rem',
          marginLeft: 'auto',
          marginRight: 'auto',
          paddingLeft: '2rem',
          paddingRight: '2rem',
        },
        '.desktop-section': {
          paddingTop: '3rem',
          paddingBottom: '3rem',
        },
      }
      addUtilities(newUtilities)
    },
  ],
}