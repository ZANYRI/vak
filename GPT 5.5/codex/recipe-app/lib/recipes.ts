import { recipes, type Recipe } from "@/data/recipes";
import type { Locale } from "@/i18n/config";

export type RecipeQuery = {
  search?: string;
  cuisine?: string;
  category?: string;
  difficulty?: string;
  diet?: string;
  meal?: string;
  time?: string;
  sort?: string;
};

export const perPage = 8;

export function filterRecipes(query: RecipeQuery, locale: Locale): Recipe[] {
  const needle = query.search?.trim().toLocaleLowerCase(locale);
  const maxTime = Number(query.time);
  const filtered = recipes.filter((recipe) => {
    const searchable = [
      recipe.title[locale], recipe.description[locale], ...recipe.ingredients[locale], ...recipe.tags,
      recipe.cuisine, recipe.category
    ].join(" ").toLocaleLowerCase(locale);
    return (!needle || searchable.includes(needle)) &&
      (!query.cuisine || recipe.cuisine === query.cuisine) &&
      (!query.category || recipe.category === query.category) &&
      (!query.difficulty || recipe.difficulty === query.difficulty) &&
      (!query.diet || recipe.diet.includes(query.diet as never)) &&
      (!query.meal || recipe.mealType.includes(query.meal)) &&
      (!Number.isFinite(maxTime) || maxTime <= 0 || recipe.totalTimeMinutes <= maxTime);
  });

  return filtered.sort((a, b) => {
    if (query.sort === "fastest") return a.totalTimeMinutes - b.totalTimeMinutes;
    if (query.sort === "easiest") return a.difficulty.localeCompare(b.difficulty);
    if (query.sort === "rated") return b.rating - a.rating;
    return b.publishedAt.localeCompare(a.publishedAt);
  });
}

export function relatedRecipes(recipe: Recipe): Recipe[] {
  return recipes
    .filter((candidate) => candidate.slug !== recipe.slug && (candidate.cuisine === recipe.cuisine || candidate.category === recipe.category))
    .slice(0, 3);
}
