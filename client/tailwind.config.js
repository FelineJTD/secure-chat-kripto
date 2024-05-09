/** @type {import('tailwindcss').Config} */
export default {
  content: ["./src/**/*.{html,js,svelte,ts}"],
  theme: {
    colors: {
      primary: "#D7E2E3",
      secondary: "#E3DBD7",
      neutral: {
        100: "#FDFCFB",
        200: "#F5EFE4",
        300: "#CFC9C2",
        400: "#A8A3A1",
        500: "#827D7F",
        600: "#5B565D",
        700: "#34303C",
        800: "#0E0A1A",
        900: "#020204"
      },
      red: "#8A2939",
      green: "#298A7A"
    },
    extend: {}
  },
  plugins: []
};
