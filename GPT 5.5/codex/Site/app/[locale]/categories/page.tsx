import Link from "next/link";
import { notFound } from "next/navigation";
import { categories, recipes } from "@/data/recipes";
import { isLocale } from "@/i18n/config";
import { getDictionary } from "@/i18n/dictionary";
import { localizedMetadata } from "@/lib/seo";

export async function generateMetadata({ params }: { params: Promise<{ locale: string }> }) { const { locale } = await params; if (!isLocale(locale)) return {}; const d = getDictionary(locale); return localizedMetadata(locale, "/categories", d.categories.title, d.categories.description); }
export default async function CategoriesPage({ params }: { params: Promise<{ locale: string }> }) { const { locale } = await params; if (!isLocale(locale)) notFound(); const d = getDictionary(locale); return <><section className="page-intro shell"><p className="eyebrow">{d.nav.recipes}</p><h1>{d.categories.title}</h1><p>{d.categories.description}</p></section><section className="shell category-grid">{categories.map((category, index) => { const count = recipes.filter((recipe) => recipe.category === category).length; return <Link key={category} className="category-tile" href={`/${locale}/categories/${category}`}><span>{String(index + 1).padStart(2, "0")}</span><h2>{d.categories[category]}</h2><p>{count} {d.recipes.results}</p><b>→</b></Link>; })}</section></>; }
