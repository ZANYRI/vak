import { notFound } from "next/navigation";
import { Suspense } from "react";
import { Footer } from "@/components/layout/Footer";
import { Header } from "@/components/layout/Header";
import { isLocale, locales } from "@/i18n/config";
import { getDictionary } from "@/i18n/dictionary";

export function generateStaticParams() { return locales.map((locale) => ({ locale })); }

export default async function LocaleLayout({ children, params }: Readonly<{ children: React.ReactNode; params: Promise<{ locale: string }> }>) {
  const { locale: rawLocale } = await params;
  if (!isLocale(rawLocale)) notFound();
  const dictionary = getDictionary(rawLocale);
  return <><a className="skip-link" href="#main-content">{dictionary.common.skip}</a><Suspense fallback={<header className="site-header"><div className="shell header-inner"><span className="brand">mis<span>e</span>n</span></div></header>}><Header locale={rawLocale} dictionary={dictionary} /></Suspense><main id="main-content">{children}</main><Footer locale={rawLocale} dictionary={dictionary} /></>;
}
