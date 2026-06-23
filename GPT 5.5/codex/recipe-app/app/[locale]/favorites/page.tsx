import { notFound } from "next/navigation";
import { FavoritesClient } from "@/components/recipes/FavoritesClient";
import { recipes } from "@/data/recipes";
import { isLocale } from "@/i18n/config";
import { getDictionary } from "@/i18n/dictionary";
import { localizedMetadata } from "@/lib/seo";

export async function generateMetadata({ params }: { params: Promise<{ locale: string }> }) { const { locale } = await params; if (!isLocale(locale)) return {}; const d = getDictionary(locale); return localizedMetadata(locale, "/favorites", d.favorites.title, d.favorites.description); }
export default async function FavoritesPage({ params }: { params: Promise<{ locale: string }> }) { const { locale } = await params; if (!isLocale(locale)) notFound(); const d = getDictionary(locale); return <><section className="page-intro shell"><p className="eyebrow">{d.nav.favorites}</p><h1>{d.favorites.title}</h1><p>{d.favorites.description}</p></section><section className="shell"><FavoritesClient recipes={recipes} locale={locale} dictionary={d} /></section></>; }
