import { recipes as allRecipes } from "@/data/recipes";
import { PAGE_SIZE, type FilterParams, type Locale, type PaginatedResult, type Recipe, type SortOption } from "./types";

export function getAllRecipes(): Recipe[] {
  return allRecipes;
}

export function getRecipeBySlug(slug: string): Recipe | undefined {
  return allRecipes.find((r) => r.slug === slug);
}

export function getRecipesByCuisine(cuisineSlug: string): Recipe[] {
  return allRecipes.filter((r) => r.cuisine === cuisineSlug);
}

export function getRecipesByCategory(categorySlug: string): Recipe[] {
  return allRecipes.filter((r) => r.category === categorySlug);
}

export function getRecipesBySlug(slugs: string[]): Recipe[] {
  const set = new Set(slugs);
  return allRecipes.filter((r) => set.has(r.slug));
}

export function getFeaturedRecipes(limit = 6): Recipe[] {
  return [...allRecipes]
    .sort((a, b) => b.rating - a.rating)
    .slice(0, limit);
}

export function getRelatedRecipes(recipe: Recipe, limit = 3): Recipe[] {
  return allRecipes
    .filter((r) => r.slug !== recipe.slug)
    .map((r) => {
      let score = 0;
      if (r.cuisine === recipe.cuisine) score += 3;
      if (r.category === recipe.category) score += 2;
      if (r.mealType.some((m) => recipe.mealType.includes(m))) score += 1;
      if (r.diet.some((d) => recipe.diet.includes(d))) score += 1;
      return { r, score };
    })
    .sort((a, b) => b.score - a.score || b.r.rating - a.r.rating)
    .slice(0, limit)
    .map((x) => x.r);
}

function normalize(value: string): string {
  return value.toLowerCase().trim();
}

function recipeMatchesSearch(recipe: Recipe, query: string, _locale: Locale): boolean {
  const q = normalize(query);
  if (!q) return true;
  const haystack = [
    recipe.title.en,
    recipe.title.ru,
    recipe.description.en,
    recipe.description.ru,
    recipe.tags.join(" "),
    recipe.cuisine,
    recipe.category,
    recipe.ingredients.en.join(" "),
    recipe.ingredients.ru.join(" ")
  ];
  // Always search both languages so users can find a recipe regardless of UI locale.
  return haystack.some((field) => normalize(field).includes(q));
}

function applyFilters(recipes: Recipe[], params: FilterParams, locale: Locale): Recipe[] {
  let result = recipes;

  if (params.search) {
    result = result.filter((r) => recipeMatchesSearch(r, params.search!, locale));
  }
  if (params.cuisine) {
    result = result.filter((r) => r.cuisine === params.cuisine);
  }
  if (params.category) {
    result = result.filter((r) => r.category === params.category);
  }
  if (params.difficulty) {
    result = result.filter((r) => r.difficulty === params.difficulty);
  }
  if (params.diet) {
    result = result.filter((r) => r.diet.includes(params.diet!));
  }
  if (params.mealType) {
    result = result.filter((r) => r.mealType.includes(params.mealType!));
  }
  if (params.maxTime) {
    const max = Number(params.maxTime);
    if (!Number.isNaN(max) && max > 0) {
      result = result.filter((r) => r.totalTimeMinutes <= max);
    }
  }
  return result;
}

function sortRecipes(recipes: Recipe[], sort: SortOption): Recipe[] {
  const copy = [...recipes];
  switch (sort) {
    case "newest":
      return copy.sort((a, b) => b.datePublished.localeCompare(a.datePublished));
    case "fastest":
      return copy.sort((a, b) => a.totalTimeMinutes - b.totalTimeMinutes);
    case "easiest":
      return copy.sort(
        (a, b) => difficultyRank(a.difficulty) - difficultyRank(b.difficulty) || a.totalTimeMinutes - b.totalTimeMinutes
      );
    case "rating":
      return copy.sort((a, b) => b.rating - a.rating);
    default:
      return copy;
  }
}

function difficultyRank(d: Recipe["difficulty"]): number {
  return d === "easy" ? 0 : d === "medium" ? 1 : 2;
}

export function queryRecipes(
  params: FilterParams,
  page: number,
  locale: Locale,
  pageSize: number = PAGE_SIZE
): PaginatedResult<Recipe> {
  const filtered = applyFilters(allRecipes, params, locale);
  const sorted = sortRecipes(filtered, params.sort || "newest");
  const total = sorted.length;
  const totalPages = Math.max(1, Math.ceil(total / pageSize));
  const safePage = Math.min(Math.max(1, page), totalPages);
  const start = (safePage - 1) * pageSize;
  const items = sorted.slice(start, start + pageSize);
  const from = total === 0 ? 0 : start + 1;
  const to = Math.min(start + pageSize, total);
  return { items, total, page: safePage, pageSize, totalPages, from, to };
}

export function getRecipesCountByCuisine(slug: string): number {
  return allRecipes.filter((r) => r.cuisine === slug).length;
}

export function getRecipesCountByCategory(slug: string): number {
  return allRecipes.filter((r) => r.category === slug).length;
}
