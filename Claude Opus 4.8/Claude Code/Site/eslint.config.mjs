import next from 'eslint-config-next';

/**
 * Flat ESLint config. Next 16 ships a native flat-config array from
 * `eslint-config-next` (bundles core-web-vitals + the TypeScript rules), so we
 * spread it directly rather than going through the eslintrc compatibility layer.
 */
const eslintConfig = [
  {
    ignores: ['.next/**', 'node_modules/**', 'next-env.d.ts'],
  },
  ...next,
];

export default eslintConfig;
