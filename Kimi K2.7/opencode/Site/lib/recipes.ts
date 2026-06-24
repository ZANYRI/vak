import { recipes } from "@/data/recipes";
import type { Recipe } from "@/lib/validation";

export const PAGE_SIZE = 9;

export const DIFFICULTY_ORDER = { easy: 1, medium: 2, hard: 3 };

export function getAllRecipes(): Recipe[] {
  return recipes;
}

export function getRecipeBySlug(slug: string): Recipe | undefined {
  return recipes.find((recipe) => recipe.slug === slug);
}

export function getUniqueValues<T extends keyof Recipe>(
  key: T,
): Array<Recipe[T]> {
  const values = recipes.map((recipe) => recipe[key]);
  return Array.from(new Set(values.flat() as Array<Recipe[T]>));
}

export type FilterParams = {
  search?: string;
  cuisine?: string;
  category?: string;
  difficulty?: string;
  diet?: string;
  mealType?: string;
  time?: string;
};

export function filterAndSortRecipes(
  locale: "en" | "ru",
  params: FilterParams & { sort?: string },
): Recipe[] {
  const { search, cuisine, category, difficulty, diet, mealType, time, sort } =
    params;
  let result = [...recipes];

  if (search) {
    const q = search.toLowerCase();
    result = result.filter((recipe) => {
      const title = recipe.title[locale].toLowerCase();
      const description = recipe.description[locale].toLowerCase();
      const ingredients = recipe.ingredients[locale].join(" ").toLowerCase();
      const tags = recipe.tags.join(" ").toLowerCase();
      const cuisineMatches = recipe.cuisine.toLowerCase().includes(q);
      const categoryMatches = recipe.category.toLowerCase().includes(q);
      return (
        title.includes(q) ||
        description.includes(q) ||
        ingredients.includes(q) ||
        tags.includes(q) ||
        cuisineMatches ||
        categoryMatches
      );
    });
  }

  if (cuisine) result = result.filter((r) => r.cuisine === cuisine);
  if (category) result = result.filter((r) => r.category === category);
  if (difficulty) result = result.filter((r) => r.difficulty === difficulty);
  if (diet) result = result.filter((r) => r.diet.includes(diet));
  if (mealType) result = result.filter((r) => r.mealType.includes(mealType));
  if (time) {
    const max = Number(time);
    if (!Number.isNaN(max)) {
      result = result.filter((r) => r.totalTimeMinutes <= max);
    }
  }

  switch (sort) {
    case "fastest":
      result.sort((a, b) => a.totalTimeMinutes - b.totalTimeMinutes);
      break;
    case "easiest":
      result.sort(
        (a, b) =>
          DIFFICULTY_ORDER[a.difficulty] - DIFFICULTY_ORDER[b.difficulty],
      );
      break;
    case "highestRated":
      result.sort((a, b) => b.rating - a.rating);
      break;
    case "newest":
    default:
      result.sort((a, b) => Number(b.id) - Number(a.id));
      break;
  }

  return result;
}

export function paginateRecipes(
  items: Recipe[],
  pageNumber: number,
  pageSize = PAGE_SIZE,
): {
  items: Recipe[];
  currentPage: number;
  totalPages: number;
  totalCount: number;
} {
  const totalCount = items.length;
  const totalPages = Math.max(1, Math.ceil(totalCount / pageSize));
  const currentPage = Math.min(Math.max(pageNumber, 1), totalPages);
  const start = (currentPage - 1) * pageSize;
  const end = start + pageSize;
  return {
    items: items.slice(start, end),
    currentPage,
    totalPages,
    totalCount,
  };
}

export function getRelatedRecipes(recipe: Recipe, count = 3): Recipe[] {
  return recipes
    .filter((r) => r.id !== recipe.id)
    .sort((a, b) => {
      const aScore =
        (a.cuisine === recipe.cuisine ? 2 : 0) +
        (a.category === recipe.category ? 1 : 0);
      const bScore =
        (b.cuisine === recipe.cuisine ? 2 : 0) +
        (b.category === recipe.category ? 1 : 0);
      return bScore - aScore || b.rating - a.rating;
    })
    .slice(0, count);
}

export const CATEGORIES = [
  "breakfast",
  "lunch",
  "dinner",
  "desserts",
  "soups",
  "salads",
  "vegetarian",
  "quick-meals",
  "baking",
  "drinks",
] as const;

export const CUISINES = [
  "italian",
  "french",
  "georgian",
  "japanese",
  "mexican",
  "indian",
  "mediterranean",
  "ukrainian",
  "latvian",
  "american",
] as const;
