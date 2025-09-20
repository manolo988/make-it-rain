/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        rain: {
          pink: '#ff0080',
          'pink-light': '#ff66b3',
          'pink-dark': '#cc0066',
          black: '#0a0a0a',
          'gray-dark': '#1a1a1a',
          'gray-medium': '#2a2a2a',
          'gray-light': '#3a3a3a',
          'gray-border': '#404040',
        }
      },
      animation: {
        'fade-in': 'fadeIn 0.5s cubic-bezier(0.4, 0, 0.2, 1)',
        'slide-up': 'slideUp 0.3s cubic-bezier(0.4, 0, 0.2, 1)',
        'glow': 'glow 2s ease-in-out infinite',
      },
      keyframes: {
        fadeIn: {
          '0%': { opacity: '0' },
          '100%': { opacity: '1' },
        },
        slideUp: {
          '0%': { transform: 'translateY(10px)', opacity: '0' },
          '100%': { transform: 'translateY(0)', opacity: '1' },
        },
        glow: {
          '0%, 100%': { boxShadow: '0 0 20px rgba(255, 0, 128, 0.3)' },
          '50%': { boxShadow: '0 0 30px rgba(255, 0, 128, 0.5)' },
        },
      },
      backgroundImage: {
        'gradient-radial': 'radial-gradient(var(--tw-gradient-stops))',
        'gradient-rain': 'linear-gradient(135deg, #ff0080 0%, #ff66b3 100%)',
        'gradient-dark': 'linear-gradient(180deg, #0a0a0a 0%, #1a1a1a 100%)',
      },
      fontFamily: {
        'sans': ['-apple-system', 'BlinkMacSystemFont', 'Inter', 'system-ui', 'sans-serif'],
      },
    },
  },
  plugins: [],
}