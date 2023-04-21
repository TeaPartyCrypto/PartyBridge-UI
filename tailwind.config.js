/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./src/**/*.{js,jsx,ts,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        'accent-primary': '#CC6228',
        'dark-primary': '#181C27',
        'dark-secondary': '#202020',
        'dark-third': '#181C27',
        'dark-fourth': '#646464',
        'light-primary':  '#FEFEFE',
        'light-secondary': '#FFFFFF'
      }
    },
  },
  plugins: [],
}

