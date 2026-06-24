import type { Metadata } from "next";
import { getTranslations } from "next-intl/server";

export type SiteLocale = "en" | "ru";

export function getSiteUrl(): string {
  return process.env.NEXT_PUBLIC_SITE_URL ?? "http://localhost:3000";
}

export async function createMetadata(
  locale: SiteLocale,
  options: {
    titleKey: string | [string, string];
    descriptionKey: string | [string, string];
    path?: string;
    image?: string;
  },
): Promise<Metadata> {
  const metaT = await getTranslations({ locale, namespace: "metadata" });

  let title: string;
  let description: string;

  if (Array.isArray(options.titleKey)) {
    const [namespace, key] = options.titleKey;
    const t = await getTranslations({ locale, namespace });
    title = t(key);
  } else {
    title = metaT(options.titleKey);
  }

  if (Array.isArray(options.descriptionKey)) {
    const [namespace, key] = options.descriptionKey;
    const t = await getTranslations({ locale, namespace });
    description = t(key);
  } else {
    description = metaT(options.descriptionKey);
  }

  const siteUrl = getSiteUrl();
  const canonical = options.path
    ? `${siteUrl}/${locale}${options.path}`
    : `${siteUrl}/${locale}`;
  const alternates = {
    canonical,
    languages: {
      en: options.path ? `${siteUrl}/en${options.path}` : `${siteUrl}/en`,
      ru: options.path ? `${siteUrl}/ru${options.path}` : `${siteUrl}/ru`,
    },
  };

  return {
    title,
    description,
    openGraph: {
      title,
      description,
      url: canonical,
      type: "website",
      locale,
      images: options.image ? [{ url: options.image }] : undefined,
    },
    twitter: {
      card: "summary_large_image",
      title,
      description,
      images: options.image ? [options.image] : undefined,
    },
    alternates,
  };
}

export function buildJsonLdRecipe(recipe: {
  title: string;
  description: string;
  image: string;
  cuisine: string;
  category: string;
  ingredients: string[];
  instructions: string[];
  prepTimeMinutes: number;
  cookTimeMinutes: number;
  totalTimeMinutes: number;
  servings: number;
  rating: number;
  nutrition: { calories: number; protein: number; fat: number; carbs: number };
}): Record<string, unknown> {
  const isoDuration = (minutes: number) => `PT${minutes}M`;

  return {
    "@context": "https://schema.org",
    "@type": "Recipe",
    name: recipe.title,
    description: recipe.description,
    image: recipe.image,
    author: { "@type": "Organization", name: "Recipe Haven" },
    datePublished: new Date().toISOString().split("T")[0],
    prepTime: isoDuration(recipe.prepTimeMinutes),
    cookTime: isoDuration(recipe.cookTimeMinutes),
    totalTime: isoDuration(recipe.totalTimeMinutes),
    recipeYield: `${recipe.servings} servings`,
    recipeCuisine: recipe.cuisine,
    recipeCategory: recipe.category,
    recipeIngredient: recipe.ingredients,
    recipeInstructions: recipe.instructions.map((step, index) => ({
      "@type": "HowToStep",
      position: index + 1,
      text: step,
    })),
    nutrition: {
      "@type": "NutritionInformation",
      calories: `${recipe.nutrition.calories} kcal`,
      proteinContent: `${recipe.nutrition.protein} g`,
      fatContent: `${recipe.nutrition.fat} g`,
      carbohydrateContent: `${recipe.nutrition.carbs} g`,
    },
    aggregateRating: {
      "@type": "AggregateRating",
      ratingValue: recipe.rating,
      reviewCount: 1,
    },
  };
}
