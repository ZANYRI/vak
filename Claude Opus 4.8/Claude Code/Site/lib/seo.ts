import type { Metadata } from 'next';
import { locales } from '@/i18n/routing';
import { SITE_URL } from './constants';
import type { Locale, Recipe } from './types';

const SITE_NAME = 'Savora';

/** Build the absolute URL for a locale-prefixed path (path starts with "/"). */
export function localizedUrl(locale: Locale, path = ''): string {
  const clean = path === '/' ? '' : path;
  return `${SITE_URL}/${locale}${clean}`;
}

/** Map every locale to its absolute URL — used for hreflang alternates. */
function languageAlternates(path = ''): Record<string, string> {
  return Object.fromEntries(
    locales.map((locale) => [locale, localizedUrl(locale, path)]),
  );
}

type BuildMetadataOptions = {
  locale: Locale;
  /** Path without the locale prefix, e.g. "/recipes" or "/recipes/classic-carbonara". */
  path?: string;
  title: string;
  description: string;
  images?: string[];
  type?: 'website' | 'article';
};

export function buildMetadata({
  locale,
  path = '',
  title,
  description,
  images,
  type = 'website',
}: BuildMetadataOptions): Metadata {
  const url = localizedUrl(locale, path);
  const ogImages = images?.map((image) => ({ url: image }));

  return {
    title,
    description,
    alternates: {
      canonical: url,
      languages: { ...languageAlternates(path), 'x-default': localizedUrl('en', path) },
    },
    openGraph: {
      type,
      siteName: SITE_NAME,
      locale,
      url,
      title,
      description,
      images: ogImages,
    },
    twitter: {
      card: 'summary_large_image',
      title,
      description,
      images,
    },
  };
}

/** Convert a number of minutes to an ISO 8601 duration, e.g. 25 -> "PT25M". */
export function toIsoDuration(minutes: number): string {
  if (minutes <= 0) return 'PT0M';
  const hours = Math.floor(minutes / 60);
  const mins = minutes % 60;
  return `PT${hours ? `${hours}H` : ''}${mins ? `${mins}M` : ''}`;
}

/** Build schema.org Recipe JSON-LD for a recipe detail page. */
export function buildRecipeJsonLd(recipe: Recipe, locale: Locale) {
  return {
    '@context': 'https://schema.org',
    '@type': 'Recipe',
    name: recipe.title[locale],
    description: recipe.description[locale],
    image: [recipe.image],
    author: { '@type': 'Organization', name: SITE_NAME },
    datePublished: recipe.publishedAt,
    prepTime: toIsoDuration(recipe.prepTimeMinutes),
    cookTime: toIsoDuration(recipe.cookTimeMinutes),
    totalTime: toIsoDuration(recipe.totalTimeMinutes),
    recipeYield: `${recipe.servings}`,
    recipeCuisine: recipe.cuisine,
    recipeCategory: recipe.category,
    keywords: recipe.tags.join(', '),
    recipeIngredient: recipe.ingredients[locale],
    recipeInstructions: recipe.steps[locale].map((step, index) => ({
      '@type': 'HowToStep',
      position: index + 1,
      text: step,
    })),
    nutrition: {
      '@type': 'NutritionInformation',
      calories: `${recipe.nutrition.calories} kcal`,
      proteinContent: `${recipe.nutrition.protein} g`,
      fatContent: `${recipe.nutrition.fat} g`,
      carbohydrateContent: `${recipe.nutrition.carbs} g`,
    },
    aggregateRating: {
      '@type': 'AggregateRating',
      ratingValue: recipe.rating,
      bestRating: 5,
      ratingCount: Math.round(recipe.rating * 24),
    },
    url: localizedUrl(locale, `/recipes/${recipe.slug}`),
  };
}
