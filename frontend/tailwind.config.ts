import type { Config } from "tailwindcss";

const config: Config = {
  content: [
    "./app/**/*.{ts,tsx}",
    "./components/**/*.{ts,tsx}",
    "./lib/**/*.{ts,tsx}"
  ],
  theme: {
    extend: {
      colors: {
        base: "#081217",
        panel: "#10252b",
        panelAlt: "#16333a",
        line: "#2d5b63",
        glow: "#f6c453",
        signal: "#6af2d4",
        danger: "#ff7a66"
      },
      boxShadow: {
        panel: "0 18px 60px rgba(0, 0, 0, 0.22)"
      },
      backgroundImage: {
        grid: "linear-gradient(rgba(106,242,212,0.08) 1px, transparent 1px), linear-gradient(90deg, rgba(106,242,212,0.08) 1px, transparent 1px)"
      }
    }
  },
  plugins: []
};

export default config;
