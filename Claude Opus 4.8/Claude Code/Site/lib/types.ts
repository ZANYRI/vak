import type { Locale } from '@/i18n/routing';

export type { Locale };

export type LocalizedText = {
  en: string;
  ru: string;
};

export type LocalizedList = {
  en: string[];
  ru: string[];
};

export type Difficulty = 'easy' | 'medium' | 'hard';

export type CategorySlug =
  | 'breakfast'
  | 'lunch'
  | 'dinner'
  | 'desserts'
  | 'soups'
  | 'salads'
  | 'vegetarian'
  | 'quick-meals'
  | 'baking'
  | 'drinks';

export type CuisineSlug =
  | 'italian'
  | 'french'
  | 'georgian'
  | 'japanese'
  | 'mexican'
  | 'indian'
  | 'mediterranean'
  | 'ukrainian'
  | 'latvian'
  | 'american';

export type MealType =
  | 'breakfast'
  | 'lunch'
  | 'dinner'
  | 'dessert'
  | 'snack'
  | 'appetizer'
  | 'drink'
  | 'side';

export type Diet =
  | 'vegetarian'
  | 'vegan'
  | 'gluten-free'
  | 'dairy-free'
  | 'pescatarian'
  | 'high-protein'
  | 'low-carb';

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
  cuisine: CuisineSlug;
  category: CategorySlug;
  mealType: MealType[];
  diet: Diet[];
  difficulty: Difficulty;
  prepTimeMinutes: number;
  cookTimeMinutes: number;
  totalTimeMinutes: number;
  servings: number;
  rating: number;
  /** ISO date string used for "newest" sorting and SEO datePublished. */
  publishedAt: string;
  ingredients: LocalizedList;
  steps: LocalizedList;
  nutrition: Nutrition;
  tips: LocalizedList;
  tags: string[];
};

export type SortKey = 'newest' | 'fastest' | 'easiest' | 'rating';

export type RecipeFilters = {
  search?: string;
  cuisine?: CuisineSlug;
  category?: CategorySlug;
  difficulty?: Difficulty;
  diet?: Diet;
  mealType?: MealType;
  maxTime?: number;
  sort?: SortKey;
};
