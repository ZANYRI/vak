import Link from "next/link";
import { notFound } from "next/navigation";
import { AnimatedSection } from "@/components/ui/AnimatedSection";
import { ImageWithFallback } from "@/components/ui/ImageWithFallback";
import { RecipeGrid } from "@/components/recipes/RecipeGrid";
import { cuisines, recipes } from "@/data/recipes";
import { isLocale, type Locale } from "@/i18n/config";
import { getDictionary } from "@/i18n/dictionary";
import { localizedMetadata } from "@/lib/seo";

export async function generateMetadata({ params }: { params: Promise<{ locale: string }> }) {
  const { locale } = await params; if (!isLocale(locale)) return {};
  const dictionary = getDictionary(locale); return localizedMetadata(locale, "", dictionary.hero.title, dictionary.hero.description);
}

export default async function HomePage({ params }: { params: Promise<{ locale: string }> }) {
  const { locale: rawLocale } = await params; if (!isLocale(rawLocale)) notFound(); const locale: Locale = rawLocale; const dictionary = getDictionary(locale);
  const featured = recipes.slice(0, 6); const heroRecipe = recipes[0];
  return <>
    <section className="hero"><div className="shell hero-grid"><AnimatedSection className="hero-copy"><p className="eyebrow">{dictionary.hero.eyebrow}</p><h1>{dictionary.hero.title}</h1><p>{dictionary.hero.description}</p><Link className="button button-primary" href={`/${locale}/recipes`}>{dictionary.hero.cta} <span aria-hidden="true">→</span></Link><div className="hero-note"><span>✦</span>{dictionary.hero.featured}</div></AnimatedSection><div className="hero-visual"><div className="hero-image main"><ImageWithFallback src={heroRecipe.image} alt={heroRecipe.imageAlt[locale]} sizes="(max-width: 800px) 100vw, 50vw" /></div><div className="hero-image floating"><ImageWithFallback src={recipes[6].image} alt={recipes[6].imageAlt[locale]} sizes="220px" /></div><div className="hero-orbit">fresh<br />everyday<br />food</div></div></div></section>
    <AnimatedSection className="shell section-head" delay={90}><div><p className="eyebrow">{dictionary.hero.featured}</p><h2>{dictionary.home.featured}</h2><p>{dictionary.home.featuredDescription}</p></div><Link className="text-link" href={`/${locale}/recipes`}>{dictionary.home.viewAll} →</Link></AnimatedSection>
    <section className="shell"><RecipeGrid recipes={featured} locale={locale} dictionary={dictionary} /></section>
    <AnimatedSection className="shell cuisine-section" delay={140}><div className="section-head"><div><p className="eyebrow">{dictionary.nav.cuisines}</p><h2>{dictionary.home.cuisines}</h2></div><Link className="text-link" href={`/${locale}/cuisines`}>{dictionary.home.viewAll} →</Link></div><div className="cuisine-pills">{cuisines.map((cuisine, index) => <Link key={cuisine} href={`/${locale}/cuisines/${cuisine}`}><span>{String(index + 1).padStart(2, "0")}</span>{dictionary.cuisines[cuisine]}</Link>)}</div></AnimatedSection>
  </>;
}
