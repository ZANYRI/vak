// Flat ESLint config for Next.js 16.
// `eslint-config-next` ships a native flat-config array (includes core-web-vitals,
// react, react-hooks, @next/eslint-plugin-next, import, jsx-a11y, and typescript).
import nextConfig from "eslint-config-next";

const eslintConfig = [
  ...nextConfig,
  {
    ignores: ["**/*.config.mjs", ".next/**", "node_modules/**", "out/**", "next-env.d.ts"]
  },
  {
    rules: {
      "@typescript-eslint/no-explicit-any": "warn",
      "@typescript-eslint/no-unused-vars": [
        "warn",
        { argsIgnorePattern: "^_", varsIgnorePattern: "^_" }
      ]
    }
  }
];

export default eslintConfig;
