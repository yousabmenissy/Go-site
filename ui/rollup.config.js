import typescript from "rollup-plugin-typescript2";

export default {
  input: "./static/ts/base.ts", // Entry point
  output: {
    file: "./static/js/base.js", // Output file
    format: "es", // Output format (ES Module)
  },
  plugins: [
    typescript(), // TypeScript plugin
  ],
};
