import { notFound } from "next/navigation";
import { RecipeGrid } from "@/components/recipes/RecipeGrid";
import { cuisines, recipes } from "@/data/recipes";
import { isLocale, locales } from "@/i18n/config";
import { getDictionary } from "@/i18n/dictionary";
import { localizedMetadata } from "@/lib/seo";

export function generateStaticParams() { return locales.flatMap((locale) => cuisines.map((cuisine) => ({ locale, cuisine }))); }
export async function generateMetadata({ params }: { params: Promise<{ locale: string; cuisine: string }> }) { const { locale, cuisine } = await params; if (!isLocale(locale) || !cuisines.includes(cuisine as never)) return {}; const d = getDictionary(locale); const key = cuisine as (typeof cuisines)[number]; return localizedMetadata(locale, `/cuisines/${cuisine}`, d.cuisines[key], d.cuisineIntro[key]); }
export default async function CuisinePage({ params }: { params: Promise<{ locale: string; cuisine: string }> }) { const { locale, cuisine } = await params; if (!isLocale(locale) || !cuisines.includes(cuisine as never)) notFound(); const d = getDictionary(locale); const key = cuisine as (typeof cuisines)[number]; return <><section className="page-intro shell"><p className="eyebrow">{d.nav.cuisines}</p><h1>{d.cuisines[key]}</h1><p>{d.cuisineIntro[key]}</p></section><section className="shell"><RecipeGrid recipes={recipes.filter((recipe) => recipe.cuisine === key)} locale={locale} dictionary={d} /></section></>; }
