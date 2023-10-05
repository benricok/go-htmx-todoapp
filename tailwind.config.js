/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./web/public/*.{html,js}"],
  theme: {
    extend: {},
  },
  plugins: [
    require('@tailwindcss/forms'),
  ],
}

