module.exports = {
  content: ["./src/**/*.{js,jsx,ts,tsx}"],
  theme: {
    extend: {
      animation: {
        "spin-slow": "spin 2s linear infinite",
      },
    },
    fontFamily: {
      sans: ["Nunito", "system-ui"],
    },
  },
  variants: {
    extend: {},
  },
  plugins: [],
};
