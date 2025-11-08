/** @type {import('tailwindcss').Config} */
export default {
  content: ['./src/**/*.{html,js,svelte,ts}'],
  theme: {
    extend: {
      fontFamily: {
        'mono': ['Courier New', 'Monaco', 'Consolas', 'monospace'],
      },
      colors: {
        'portal-blue': '#00BFFF',
        'portal-light-blue': '#E0F7FF',
      },
    },
  },
  plugins: [],
}