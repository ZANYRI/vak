import type { routing } from "@/i18n/routing";

export type Locale = (typeof routing.locales)[number];

export type LocalizedText = {
  en: string;
  ru: string;
};

export type Difficulty = "easy" | "medium" | "hard";

export type Nutrition = {
  calories: number;
  protein: number;
  fat: number;
  carbs: number;
};

export type Recipe = {
  id: string;
  slug: string;
  title: LocalizedText;
  description: LocalizedText;
  image: string;
  imageAlt: LocalizedText;
  cuisine: string;
  category: string;
  mealType: string[];
  diet: string[];
  difficulty: Difficulty;
  prepTimeMinutes: number;
  cookTimeMinutes: number;
  totalTimeMinutes: number;
  servings: number;
  rating: number;
  ingredients: {
    en: string[];
    ru: string[];
  };
  steps: {
    en: string[];
    ru: string[];
  };
  nutrition: Nutrition;
  tips: {
    en: string[];
    ru: string[];
  };
  tags: string[];
  /** ISO date the recipe was published; used for "newest" sorting. */
  datePublished: string;
};

export type SortOption = "newest" | "fastest" | "easiest" | "rating";

export type FilterParams = {
  search?: string;
  cuisine?: string;
  category?: string;
  difficulty?: Difficulty | "";
  diet?: string;
  mealType?: string;
  maxTime?: string;
  sort?: SortOption | "";
};

export type TaxonomyEntry = {
  slug: string;
  name: LocalizedText;
};

export type PaginatedResult<T> = {
  items: T[];
  total: number;
  page: number;
  pageSize: number;
  totalPages: number;
  from: number;
  to: number;
};

export const PAGE_SIZE = 9;
