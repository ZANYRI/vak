import type { Metadata } from "next";
import type { Locale, Recipe } from "./types";
import { getDifficultyName, getCuisineName, getCategoryName } from "./taxonomy";

const siteUrl = process.env.NEXT_PUBLIC_SITE_URL?.replace(/\/$/, "") || "http://localhost:3000";

export function getSiteUrl(): string {
  return siteUrl;
}

export function localizedPath(locale: Locale, path: string = ""): string {
  const clean = path.startsWith("/") ? path.slice(1) : path;
  return `/${locale}${clean ? `/${clean}` : ""}`;
}

export function absoluteUrl(locale: Locale, path: string = ""): string {
  return `${siteUrl}${localizedPath(locale, path)}`;
}

type SeoInput = {
  locale: Locale;
  title: string;
  description: string;
  path?: string;
  image?: string;
};

export function buildMetadata({ locale, title, description, path = "", image }: SeoInput): Metadata {
  const url = absoluteUrl(locale, path);
  const ogImage = image || `${siteUrl}/og-default.svg`;
  const otherLocale: Locale = locale === "en" ? "ru" : "en";
  const fullTitle = title.toLowerCase().includes("saveur")
    ? title
    : `${title} — Saveur`;

  return {
    title: fullTitle,
    description,
    alternates: {
      canonical: url,
      languages: {
        [locale]: url,
        [otherLocale]: absoluteUrl(otherLocale, path),
        "x-default": absoluteUrl("en", path)
      }
    },
    openGraph: {
      title,
      description,
      url,
      siteName: "Saveur",
      locale: locale === "ru" ? "ru_RU" : "en_US",
      type: "website",
      images: [{ url: ogImage, width: 1200, height: 630, alt: title }]
    },
    twitter: {
      card: "summary_large_image",
      title,
      description,
      images: [ogImage]
    }
  };
}

export function buildRecipeMetadata(recipe: Recipe, locale: Locale): Metadata {
  const title = recipe.title[locale];
  const description = recipe.description[locale];
  const path = `recipes/${recipe.slug}`;
  return buildMetadata({
    locale,
    title: `${title} — Saveur`,
    description,
    path,
    image: recipe.image
  });
}

/** Builds a schema.org `Recipe` JSON-LD object for structured data. */
export function buildRecipeJsonLd(recipe: Recipe, locale: Locale) {
  const lang = locale === "ru" ? "ru-RU" : "en-US";
  const cuisine = getCuisineName(recipe.cuisine)[locale];
  const category = getCategoryName(recipe.category)[locale];

  return {
    "@context": "https://schema.org",
    "@type": "Recipe",
    name: recipe.title[locale],
    description: recipe.description[locale],
    image: [recipe.image],
    author: {
      "@type": "Organization",
      name: "Saveur"
    },
    datePublished: recipe.datePublished,
    inLanguage: lang,
    prepTime: `PT${recipe.prepTimeMinutes}M`,
    cookTime: `PT${recipe.cookTimeMinutes}M`,
    totalTime: `PT${recipe.totalTimeMinutes}M`,
    recipeYield: `${recipe.servings} ${locale === "ru" ? "порций" : "servings"}`,
    recipeCuisine: cuisine,
    recipeCategory: category,
    keywords: recipe.tags.join(", "),
    aggregateRating: {
      "@type": "AggregateRating",
      ratingValue: recipe.rating,
      reviewCount: Math.round(recipe.rating * 37)
    },
    nutrition: {
      "@type": "NutritionInformation",
      calories: `${recipe.nutrition.calories} kcal`,
      proteinContent: `${recipe.nutrition.protein} g`,
      fatContent: `${recipe.nutrition.fat} g`,
      carbohydrateContent: `${recipe.nutrition.carbs} g`
    },
    recipeIngredient: recipe.ingredients[locale],
    recipeInstructions: recipe.steps[locale].map((step, index) => ({
      "@type": "HowToStep",
      position: index + 1,
      text: step
    })),
    difficulty: getDifficultyName(recipe.difficulty)[locale]
  };
}
