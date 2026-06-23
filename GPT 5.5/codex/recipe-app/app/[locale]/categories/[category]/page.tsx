import { notFound } from "next/navigation";
import { RecipeGrid } from "@/components/recipes/RecipeGrid";
import { categories, recipes } from "@/data/recipes";
import { isLocale, locales } from "@/i18n/config";
import { getDictionary } from "@/i18n/dictionary";
import { localizedMetadata } from "@/lib/seo";

export function generateStaticParams() { return locales.flatMap((locale) => categories.map((category) => ({ locale, category }))); }
export async function generateMetadata({ params }: { params: Promise<{ locale: string; category: string }> }) { const { locale, category } = await params; if (!isLocale(locale) || !categories.includes(category as never)) return {}; const d = getDictionary(locale); return localizedMetadata(locale, `/categories/${category}`, d.categories[category as (typeof categories)[number]], d.categories.description); }
export default async function CategoryPage({ params }: { params: Promise<{ locale: string; category: string }> }) { const { locale, category } = await params; if (!isLocale(locale) || !categories.includes(category as never)) notFound(); const d = getDictionary(locale); const key = category as (typeof categories)[number]; return <><section className="page-intro shell"><p className="eyebrow">{d.nav.categories}</p><h1>{d.categories[key]}</h1><p>{d.categories.description}</p></section><section className="shell"><RecipeGrid recipes={recipes.filter((recipe) => recipe.category === key)} locale={locale} dictionary={d} /></section></>; }
