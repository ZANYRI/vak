import type { Metadata } from 'next';
import { getTranslations, setRequestLocale } from 'next-intl/server';
import { Link } from '@/i18n/navigation';
import { buttonClasses } from '@/components/ui/Button';
import { AnimatedSection } from '@/components/animations/AnimatedSection';
import { SearchBar } from '@/components/recipes/SearchBar';
import { RecipeGrid } from '@/components/recipes/RecipeGrid';
import { CuisineCard } from '@/components/recipes/CuisineCard';
import { CategoryCard } from '@/components/recipes/CategoryCard';
import {
  getAllRecipes,
  getFeaturedRecipes,
  getRecipesByCuisine,
  countByCategory,
  countByCuisine,
} from '@/lib/recipes';
import { CATEGORIES, POPULAR_CUISINES, CUISINES } from '@/lib/constants';
import { buildMetadata } from '@/lib/seo';
import type { Locale } from '@/lib/types';

type Props = {
  params: Promise<{ locale: string }>;
};

export async function generateMetadata({ params }: Props): Promise<Metadata> {
  const { locale } = await params;
  const t = await getTranslations({ locale, namespace: 'Home' });
  const tc = await getTranslations({ locale, namespace: 'Common' });
  return buildMetadata({
    locale: locale as Locale,
    title: `${tc('siteName')} — ${tc('tagline')}`,
    description: `${t('heroTitle')} ${t('heroSubtitle')}`,
  });
}

export default async function HomePage({ params }: Props) {
  const { locale } = await params;
  setRequestLocale(locale);
  const loc = locale as Locale;

  const t = await getTranslations({ locale, namespace: 'Home' });
  const featured = getFeaturedRecipes(6);

  const stats = [
    { value: getAllRecipes().length, label: t('statsRecipes') },
    { value: CUISINES.length, label: t('statsCuisines') },
    { value: CATEGORIES.length, label: t('statsCategories') },
  ];

  return (
    <div className="space-y-20">
      {/* Hero */}
      <section className="grid items-center gap-10 lg:grid-cols-2">
        <div>
          <AnimatedSection>
            <p className="mb-3 inline-flex items-center gap-2 rounded-full bg-tomato/10 px-4 py-1.5 text-sm font-medium text-tomato-dark">
              🍅 {t('featuredTitle')}
            </p>
          </AnimatedSection>
          <AnimatedSection delay={0.05}>
            <h1 className="font-display text-4xl font-semibold leading-tight text-charcoal text-balance sm:text-5xl lg:text-6xl">
              {t('heroTitle')}
            </h1>
          </AnimatedSection>
          <AnimatedSection delay={0.12}>
            <p className="mt-5 max-w-xl text-lg text-muted">{t('heroSubtitle')}</p>
          </AnimatedSection>
          <AnimatedSection delay={0.2} className="mt-8">
            <SearchBar className="max-w-xl" />
          </AnimatedSection>
          <AnimatedSection delay={0.28} className="mt-6 flex flex-wrap items-center gap-3">
            <Link href="/recipes" className={buttonClasses('primary', 'lg')}>
              {t('browseAll')} →
            </Link>
            <Link href="/cuisines" className={buttonClasses('outline', 'lg')}>
              {t('popularCuisinesTitle')}
            </Link>
          </AnimatedSection>

          <dl className="mt-10 flex gap-8">
            {stats.map((stat) => (
              <div key={stat.label}>
                <dt className="font-display text-3xl font-semibold text-tomato">{stat.value}</dt>
                <dd className="text-sm text-muted">{stat.label}</dd>
              </div>
            ))}
          </dl>
        </div>

        <AnimatedSection delay={0.1} className="relative">
          <div className="grid grid-cols-2 gap-4">
            {featured.slice(0, 2).map((recipe, i) => (
              <div
                key={recipe.id}
                className={i === 0 ? 'mt-8' : ''}
              >
                <Link
                  href={`/recipes/${recipe.slug}`}
                  className="group block overflow-hidden rounded-card shadow-lift"
                >
                  <div className="relative aspect-[3/4]">
                    {/* eslint-disable-next-line @next/next/no-img-element */}
                    <img
                      src={recipe.image}
                      alt={recipe.imageAlt[loc]}
                      className="h-full w-full object-cover transition-transform duration-700 group-hover:scale-105 motion-reduce:group-hover:scale-100"
                    />
                    <div className="absolute inset-0 bg-gradient-to-t from-charcoal/70 to-transparent" />
                    <div className="absolute inset-x-0 bottom-0 p-4">
                      <h2 className="font-display text-lg font-semibold text-white">
                        {recipe.title[loc]}
                      </h2>
                    </div>
                  </div>
                </Link>
              </div>
            ))}
          </div>
        </AnimatedSection>
      </section>

      {/* Featured recipes */}
      <section aria-labelledby="featured-heading">
        <AnimatedSection className="mb-8 flex items-end justify-between gap-4">
          <div>
            <h2 id="featured-heading" className="font-display text-3xl font-semibold text-charcoal">
              {t('featuredTitle')}
            </h2>
            <p className="mt-2 text-muted">{t('featuredSubtitle')}</p>
          </div>
          <Link
            href="/recipes"
            className="hidden shrink-0 text-sm font-medium text-tomato hover:underline sm:inline"
          >
            {t('browseAll')} →
          </Link>
        </AnimatedSection>
        <RecipeGrid recipes={featured} locale={loc} priorityCount={3} />
      </section>

      {/* Popular cuisines */}
      <section aria-labelledby="cuisines-heading">
        <AnimatedSection className="mb-8">
          <h2 id="cuisines-heading" className="font-display text-3xl font-semibold text-charcoal">
            {t('popularCuisinesTitle')}
          </h2>
          <p className="mt-2 text-muted">{t('popularCuisinesSubtitle')}</p>
        </AnimatedSection>
        <div className="grid grid-cols-2 gap-4 sm:grid-cols-3 lg:grid-cols-6">
          {POPULAR_CUISINES.map((cuisine, i) => (
            <AnimatedSection key={cuisine} delay={Math.min(i, 6) * 0.05}>
              <CuisineCard
                cuisine={cuisine}
                locale={loc}
                count={countByCuisine(cuisine)}
                image={getRecipesByCuisine(cuisine)[0]?.image ?? ''}
              />
            </AnimatedSection>
          ))}
        </div>
      </section>

      {/* Categories */}
      <section aria-labelledby="categories-heading">
        <AnimatedSection className="mb-8">
          <h2 id="categories-heading" className="font-display text-3xl font-semibold text-charcoal">
            {t('categoriesTitle')}
          </h2>
          <p className="mt-2 text-muted">{t('categoriesSubtitle')}</p>
        </AnimatedSection>
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
          {CATEGORIES.map((category, i) => (
            <AnimatedSection key={category} delay={Math.min(i, 9) * 0.04}>
              <CategoryCard category={category} count={countByCategory(category)} />
            </AnimatedSection>
          ))}
        </div>
      </section>

      {/* CTA */}
      <AnimatedSection>
        <section className="overflow-hidden rounded-card bg-olive px-8 py-14 text-center text-white">
          <h2 className="font-display text-3xl font-semibold text-balance sm:text-4xl">
            {t('heroTitle')}
          </h2>
          <p className="mx-auto mt-3 max-w-xl text-white/85">{t('heroSubtitle')}</p>
          <Link
            href="/recipes"
            className="mt-7 inline-flex h-13 items-center justify-center rounded-full bg-white px-8 text-base font-medium text-olive-dark transition-transform hover:-translate-y-0.5 motion-reduce:hover:translate-y-0"
          >
            {t('browseAll')} →
          </Link>
        </section>
      </AnimatedSection>
    </div>
  );
}
