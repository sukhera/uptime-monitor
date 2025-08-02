/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  darkMode: 'class',
  theme: {
    extend: {
      colors: {
        // Status colors from current design
        operational: {
          50: '#f0fdf4',
          500: '#22c55e',
          600: '#16a34a',
          900: '#14532d'
        },
        degraded: {
          50: '#fffbeb',
          500: '#f59e0b',
          600: '#d97706',
          900: '#78350f'
        },
        down: {
          50: '#fef2f2',
          500: '#ef4444',
          600: '#dc2626',
          900: '#7f1d1d'
        }
      },
      animation: {
        'pulse-slow': 'pulse 3s infinite',
        'bounce-subtle': 'bounce 2s infinite'
      }
    },
  },
  plugins: [],
}