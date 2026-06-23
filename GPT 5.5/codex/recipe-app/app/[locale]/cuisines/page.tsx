import Link from "next/link";
import { notFound } from "next/navigation";
import { cuisines, recipes } from "@/data/recipes";
import { isLocale } from "@/i18n/config";
import { getDictionary } from "@/i18n/dictionary";
import { localizedMetadata } from "@/lib/seo";

export async function generateMetadata({ params }: { params: Promise<{ locale: string }> }) { const { locale } = await params; if (!isLocale(locale)) return {}; const d = getDictionary(locale); return localizedMetadata(locale, "/cuisines", d.cuisines.title, d.cuisines.description); }
export default async function CuisinesPage({ params }: { params: Promise<{ locale: string }> }) { const { locale } = await params; if (!isLocale(locale)) notFound(); const d = getDictionary(locale); return <><section className="page-intro shell"><p className="eyebrow">{d.nav.cuisines}</p><h1>{d.cuisines.title}</h1><p>{d.cuisines.description}</p></section><section className="shell cuisine-list">{cuisines.map((cuisine, index) => <Link key={cuisine} href={`/${locale}/cuisines/${cuisine}`}><span>{String(index + 1).padStart(2, "0")}</span><div><h2>{d.cuisines[cuisine]}</h2><p>{d.cuisineIntro[cuisine]}</p></div><em>{recipes.filter((recipe) => recipe.cuisine === cuisine).length} →</em></Link>)}</section></>; }
