import type { Metadata } from 'next';
import { notFound } from 'next/navigation';
import { getTranslations, setRequestLocale } from 'next-intl/server';
import { Link } from '@/i18n/navigation';
import { RecipeHero } from '@/components/recipes/RecipeHero';
import { RecipeMeta } from '@/components/recipes/RecipeMeta';
import { IngredientList } from '@/components/recipes/IngredientList';
import { InstructionSteps } from '@/components/recipes/InstructionSteps';
import { NutritionCard } from '@/components/recipes/NutritionCard';
import { RelatedRecipes } from '@/components/recipes/RelatedRecipes';
import { getAllSlugs, getRecipeBySlug, getRelatedRecipes } from '@/lib/recipes';
import { routing } from '@/i18n/routing';
import { buildMetadata, buildRecipeJsonLd } from '@/lib/seo';
import type { Locale } from '@/lib/types';

type Props = {
  params: Promise<{ locale: string; slug: string }>;
};

export function generateStaticParams() {
  const params: { locale: string; slug: string }[] = [];
  for (const locale of routing.locales) {
    for (const slug of getAllSlugs()) {
      params.push({ locale, slug });
    }
  }
  return params;
}

export async function generateMetadata({ params }: Props): Promise<Metadata> {
  const { locale, slug } = await params;
  const recipe = getRecipeBySlug(slug);
  if (!recipe) return {};
  return buildMetadata({
    locale: locale as Locale,
    path: `/recipes/${recipe.slug}`,
    title: recipe.title[locale as Locale],
    description: recipe.description[locale as Locale],
    images: [recipe.image],
    type: 'article',
  });
}

export default async function RecipeDetailPage({ params }: Props) {
  const { locale, slug } = await params;
  setRequestLocale(locale);

  const recipe = getRecipeBySlug(slug);
  if (!recipe) {
    notFound();
  }

  const loc = locale as Locale;
  const t = await getTranslations({ locale, namespace: 'Recipe' });
  const related = getRelatedRecipes(recipe, 3);
  const jsonLd = buildRecipeJsonLd(recipe, loc);

  return (
    <article className="space-y-12">
      {/* Recipe JSON-LD structured data */}
      <script
        type="application/ld+json"
        dangerouslySetInnerHTML={{ __html: JSON.stringify(jsonLd) }}
      />

      <nav aria-label="Breadcrumb" data-no-print>
        <Link
          href="/recipes"
          className="inline-flex items-center gap-1 text-sm font-medium text-muted transition-colors hover:text-tomato"
        >
          <span aria-hidden="true">←</span>
          {t('backToRecipes')}
        </Link>
      </nav>

      <RecipeHero recipe={recipe} locale={loc} />

      <RecipeMeta recipe={recipe} />

      <div className="grid gap-8 lg:grid-cols-[1fr_1.6fr]">
        <div className="space-y-8">
          <IngredientList recipe={recipe} locale={loc} />
          <NutritionCard recipe={recipe} />
        </div>

        <div className="space-y-10">
          <InstructionSteps recipe={recipe} locale={loc} />

          {recipe.tips[loc].length > 0 && (
            <section
              aria-labelledby="tips-heading"
              className="rounded-card border border-saffron/30 bg-saffron/10 p-6"
            >
              <h2 id="tips-heading" className="font-display text-xl font-semibold text-charcoal">
                {t('tips')}
              </h2>
              <ul className="mt-4 space-y-3">
                {recipe.tips[loc].map((tip, index) => (
                  <li key={index} className="flex items-start gap-3 text-charcoal">
                    <span aria-hidden="true" className="text-lg">💡</span>
                    <span className="text-sm leading-relaxed">{tip}</span>
                  </li>
                ))}
              </ul>
            </section>
          )}
        </div>
      </div>

      <RelatedRecipes recipes={related} locale={loc} />
    </article>
  );
}
