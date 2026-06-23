import { notFound } from "next/navigation";
import { RecipeListing } from "@/components/recipes/RecipeListing";
import { isLocale } from "@/i18n/config";
import { getDictionary } from "@/i18n/dictionary";
import { localizedMetadata } from "@/lib/seo";

export async function generateMetadata({ params }: { params: Promise<{ locale: string }> }) { const { locale } = await params; if (!isLocale(locale)) return {}; const d = getDictionary(locale); return localizedMetadata(locale, "/recipes", d.recipes.title, d.recipes.description); }
export default async function RecipesPage({ params, searchParams }: { params: Promise<{ locale: string }>; searchParams: Promise<Record<string, string | string[] | undefined>> }) { const { locale } = await params; if (!isLocale(locale)) notFound(); return <RecipeListing locale={locale} page={1} searchParams={await searchParams} />; }
