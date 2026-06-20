import { recipes } from '@/data/recipes';
import type {
  CategorySlug,
  CuisineSlug,
  Locale,
  Recipe,
  RecipeFilters,
  SortKey,
} from './types';
import { PAGE_SIZE } from './constants';

const DIFFICULTY_ORDER: Record<Recipe['difficulty'], number> = {
  easy: 0,
  medium: 1,
  hard: 2,
};

export function getAllRecipes(): Recipe[] {
  return recipes;
}

export function getRecipeBySlug(slug: string): Recipe | undefined {
  return recipes.find((r) => r.slug === slug);
}

export function getAllSlugs(): string[] {
  return recipes.map((r) => r.slug);
}

export function getRecipesByCuisine(cuisine: CuisineSlug): Recipe[] {
  return recipes.filter((r) => r.cuisine === cuisine);
}

export function getRecipesByCategory(category: CategorySlug): Recipe[] {
  return recipes.filter((r) => r.category === category);
}

export function countByCuisine(cuisine: CuisineSlug): number {
  return getRecipesByCuisine(cuisine).length;
}

export function countByCategory(category: CategorySlug): number {
  return getRecipesByCategory(category).length;
}

/** Highest-rated recipes for the home page "featured" section. */
export function getFeaturedRecipes(limit = 6): Recipe[] {
  return [...recipes].sort((a, b) => b.rating - a.rating).slice(0, limit);
}

/** Related recipes share a cuisine or category but exclude the current recipe. */
export function getRelatedRecipes(recipe: Recipe, limit = 3): Recipe[] {
  return recipes
    .filter((r) => r.slug !== recipe.slug)
    .map((r) => {
      let score = 0;
      if (r.cuisine === recipe.cuisine) score += 2;
      if (r.category === recipe.category) score += 1;
      if (r.diet.some((d) => recipe.diet.includes(d))) score += 1;
      return { r, score };
    })
    .filter((x) => x.score > 0)
    .sort((a, b) => b.score - a.score || b.r.rating - a.r.rating)
    .slice(0, limit)
    .map((x) => x.r);
}

function sortRecipes(list: Recipe[], sort: SortKey = 'newest'): Recipe[] {
  const sorted = [...list];
  switch (sort) {
    case 'fastest':
      return sorted.sort((a, b) => a.totalTimeMinutes - b.totalTimeMinutes);
    case 'easiest':
      return sorted.sort(
        (a, b) => DIFFICULTY_ORDER[a.difficulty] - DIFFICULTY_ORDER[b.difficulty],
      );
    case 'rating':
      return sorted.sort((a, b) => b.rating - a.rating);
    case 'newest':
    default:
      return sorted.sort(
        (a, b) => Date.parse(b.publishedAt) - Date.parse(a.publishedAt),
      );
  }
}

function matchesSearch(recipe: Recipe, locale: Locale, query: string): boolean {
  const q = query.trim().toLowerCase();
  if (!q) return true;
  const haystack = [
    recipe.title[locale],
    recipe.description[locale],
    ...recipe.ingredients[locale],
    ...recipe.tags,
    recipe.cuisine,
    recipe.category,
    // Search both languages so a query works regardless of UI locale.
    recipe.title.en,
    recipe.title.ru,
  ]
    .join(' ')
    .toLowerCase();
  return haystack.includes(q);
}

/**
 * Filter + sort the full recipe set for a given locale. Pure and synchronous so
 * it can run on the server for SEO-friendly, no-JS listing pages.
 */
export function filterRecipes(locale: Locale, filters: RecipeFilters): Recipe[] {
  const filtered = recipes.filter((recipe) => {
    if (filters.search && !matchesSearch(recipe, locale, filters.search)) {
      return false;
    }
    if (filters.cuisine && recipe.cuisine !== filters.cuisine) return false;
    if (filters.category && recipe.category !== filters.category) return false;
    if (filters.difficulty && recipe.difficulty !== filters.difficulty) {
      return false;
    }
    if (filters.diet && !recipe.diet.includes(filters.diet)) return false;
    if (filters.mealType && !recipe.mealType.includes(filters.mealType)) {
      return false;
    }
    if (filters.maxTime && recipe.totalTimeMinutes > filters.maxTime) {
      return false;
    }
    return true;
  });

  return sortRecipes(filtered, filters.sort);
}

export type PaginatedRecipes = {
  items: Recipe[];
  page: number;
  totalPages: number;
  totalItems: number;
  pageSize: number;
};

export function paginateRecipes(
  list: Recipe[],
  page: number,
  pageSize: number = PAGE_SIZE,
): PaginatedRecipes {
  const totalItems = list.length;
  const totalPages = Math.max(1, Math.ceil(totalItems / pageSize));
  const safePage = Math.min(Math.max(1, page), totalPages);
  const start = (safePage - 1) * pageSize;
  return {
    items: list.slice(start, start + pageSize),
    page: safePage,
    totalPages,
    totalItems,
    pageSize,
  };
}

export function getTotalPages(pageSize: number = PAGE_SIZE): number {
  return Math.max(1, Math.ceil(recipes.length / pageSize));
}
