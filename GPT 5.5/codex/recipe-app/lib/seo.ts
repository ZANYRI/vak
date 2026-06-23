import type { Metadata } from "next";
import type { Locale } from "@/i18n/config";

const fallbackUrl = "http://localhost:3000";
export const siteUrl = process.env.NEXT_PUBLIC_SITE_URL ?? fallbackUrl;

export function localizedMetadata(locale: Locale, path: string, title: string, description: string): Metadata {
  const canonical = `${siteUrl}/${locale}${path}`;
  return {
    metadataBase: new URL(siteUrl), title, description,
    alternates: { canonical, languages: { en: `${siteUrl}/en${path}`, ru: `${siteUrl}/ru${path}` } },
    openGraph: { type: "website", url: canonical, title, description, locale: locale === "ru" ? "ru_RU" : "en_US", siteName: "Misen" },
    twitter: { card: "summary_large_image", title, description }
  };
}
