import { defineRouting } from "next-intl/routing";

export const routing = defineRouting({
  locales: ["en", "ru"] as const,
  defaultLocale: "en",
  localePrefix: "always",
  localeDetection: true
});

export type Locale = (typeof routing.locales)[number];
export type AppLocale = Locale;

export const locales = routing.locales;
export const defaultLocale = routing.defaultLocale;
