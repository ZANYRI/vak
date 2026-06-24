import { notFound } from "next/navigation";
import { RecipeDetail } from "@/components/recipes/RecipeDetail";
import { recipeBySlug, recipes } from "@/data/recipes";
import { isLocale, locales } from "@/i18n/config";
import { getDictionary } from "@/i18n/dictionary";
import { localizedMetadata } from "@/lib/seo";

export function generateStaticParams() { return locales.flatMap((locale) => recipes.map((recipe) => ({ locale, slug: recipe.slug }))); }
export async function generateMetadata({ params }: { params: Promise<{ locale: string; slug: string }> }) { const { locale, slug } = await params; const recipe = recipeBySlug(slug); if (!isLocale(locale) || !recipe) return {}; return { ...localizedMetadata(locale, `/recipes/${slug}`, recipe.title[locale], recipe.description[locale]), openGraph: { images: [{ url: recipe.image, alt: recipe.imageAlt[locale] }] } }; }
export default async function RecipePage({ params }: { params: Promise<{ locale: string; slug: string }> }) { const { locale, slug } = await params; const recipe = recipeBySlug(slug); if (!isLocale(locale) || !recipe) notFound(); return <RecipeDetail recipe={recipe} locale={locale} dictionary={getDictionary(locale)} />; }
