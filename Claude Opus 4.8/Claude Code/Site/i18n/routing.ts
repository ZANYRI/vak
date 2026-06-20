import { defineRouting } from 'next-intl/routing';

export const locales = ['en', 'ru'] as const;
export type Locale = (typeof locales)[number];

export const defaultLocale: Locale = 'en';

export const routing = defineRouting({
  locales,
  defaultLocale,
  // Always prefix the locale in the URL (e.g. /en, /ru) so routes are explicit.
  localePrefix: 'always',
});
