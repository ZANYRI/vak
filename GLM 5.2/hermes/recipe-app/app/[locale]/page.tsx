import { setRequestLocale, getTranslations } from "next-intl/server";
import { useLocale, useTranslations } from "next-intl";
import { Link } from "@/i18n/navigation";
import { getFeaturedRecipes } from "@/lib/recipes";
import { CUISINES } from "@/lib/taxonomy";
import { RecipeCard } from "@/components/recipes/RecipeCard";
import { SearchBar } from "@/components/filters/SearchBar";
import { AnimatedSection } from "@/components/animations/AnimatedSection";
import { Button } from "@/components/ui/Button";
import { buildMetadata } from "@/lib/seo";

export async function generateMetadata({
  params
}: {
  params: Promise<{ locale: string }>;
}) {
  const { locale } = await params;
  const t = await getTranslations({ locale, namespace: "seo" });
  return buildMetadata({
    locale: locale as "en" | "ru",
    title: "Saveur — Bilingual Recipes",
    description: t("homeDescription")
  });
}

export default async function HomePage({
  params
}: {
  params: Promise<{ locale: string }>;
}) {
  const { locale } = await params;
  setRequestLocale(locale);
  const featured = getFeaturedRecipes(6);
  const popularCuisines = CUISINES.slice(0, 6);

  return <HomeContent featured={featured} popularCuisines={popularCuisines} />;
}

function HomeContent({
  featured,
  popularCuisines
}: {
  featured: ReturnType<typeof getFeaturedRecipes>;
  popularCuisines: typeof CUISINES;
}) {
  const t = useTranslations();
  const locale = useLocale() as "en" | "ru";

  return (
    <>
      {/* Hero */}
      <section className="relative overflow-hidden border-b border-border">
        <div
          aria-hidden="true"
          className="absolute inset-0 bg-gradient-to-br from-accent/15 via-background to-secondary/10"
        />
        <div className="container-page relative grid gap-10 py-16 lg:grid-cols-2 lg:py-24">
          <AnimatedSection className="flex flex-col justify-center">
            <span className="inline-flex w-fit items-center rounded-full border border-border bg-surface/80 px-3 py-1 text-xs font-medium uppercase tracking-wide text-secondary backdrop-blur">
              {t("common.tagline")}
            </span>
            <h1 className="mt-4 font-serif text-4xl font-semibold leading-[1.1] text-foreground sm:text-5xl lg:text-6xl">
              {t("home.heroTitle")}
            </h1>
            <p className="mt-4 max-w-xl text-base text-muted sm:text-lg">
              {t("home.heroSubtitle")}
            </p>
            <div className="mt-6 flex flex-wrap gap-3">
              <Link href="/recipes">
                <Button size="lg">{t("home.heroCta")}</Button>
              </Link>
              <Link href="/categories">
                <Button size="lg" variant="outline">
                  {t("nav.categories")}
                </Button>
              </Link>
            </div>
            <div className="mt-8 max-w-md">
              <SearchBar
                basePath="/recipes"
                placeholder={t("home.searchPlaceholder")}
              />
            </div>
          </AnimatedSection>

          <AnimatedSection delay={0.15} className="relative">
            <div className="relative mx-auto grid w-full max-w-md grid-cols-2 gap-4">
              {featured.slice(0, 4).map((r, i) => (
                <Link
                  key={r.slug}
                  href={`/recipes/${r.slug}`}
                  className={`group relative overflow-hidden rounded-card border border-border shadow-soft ${i % 2 === 1 ? "mt-8" : ""}`}
                >
                  {/* eslint-disable-next-line @next/next/no-img-element */}
                  <img
                    src={r.image}
                    alt={r.imageAlt[locale]}
                    loading="lazy"
                    className="aspect-[4/5] w-full object-cover transition-transform duration-500 group-hover:scale-105"
                  />
                  <div className="absolute inset-0 bg-gradient-to-t from-black/60 to-transparent" />
                  <p className="absolute inset-x-0 bottom-0 p-3 text-sm font-medium text-white">
                    {r.title[locale]}
                  </p>
                </Link>
              ))}
            </div>
          </AnimatedSection>
        </div>
      </section>

      {/* Featured recipes */}
      <section className="container-page py-16">
        <AnimatedSection className="mb-8 flex items-end justify-between">
          <div>
            <h2 className="font-serif text-3xl font-semibold text-foreground">
              {t("home.featured")}
            </h2>
            <p className="mt-1 text-muted">{t("home.featuredSubtitle")}</p>
          </div>
          <Link
            href="/recipes"
            className="hidden text-sm font-medium text-primary hover:underline sm:inline"
          >
            {t("common.viewAll")} →
          </Link>
        </AnimatedSection>
        <ul className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3" role="list">
          {featured.map((r, i) => (
            <li key={r.slug} className="relative">
              <RecipeCard recipe={r} index={i} priority={i < 3} />
            </li>
          ))}
        </ul>
      </section>

      {/* Popular cuisines */}
      <section className="border-y border-border bg-surface-alt/50">
        <div className="container-page py-16">
          <AnimatedSection className="mb-8 text-center">
            <h2 className="font-serif text-3xl font-semibold text-foreground">
              {t("home.popularCuisines")}
            </h2>
            <p className="mt-1 text-muted">{t("home.popularCuisinesSubtitle")}</p>
          </AnimatedSection>
          <ul
            className="grid grid-cols-2 gap-4 sm:grid-cols-3 lg:grid-cols-6"
            role="list"
          >
            {popularCuisines.map((c, i) => (
              <AnimatedSection key={c.slug} delay={i * 0.05}>
                <Link
                  href={`/cuisines/${c.slug}`}
                  className="flex h-24 flex-col items-center justify-center rounded-card border border-border bg-card text-center shadow-soft transition-all hover:-translate-y-0.5 hover:shadow-lift"
                >
                  <span className="font-serif text-lg font-semibold text-foreground">
                    {c.name[locale]}
                  </span>
                </Link>
              </AnimatedSection>
            ))}
          </ul>
        </div>
      </section>
    </>
  );
}
