import type {
  CategorySlug,
  CuisineSlug,
  Diet,
  Difficulty,
  MealType,
  SortKey,
} from './types';

export const CATEGORIES: CategorySlug[] = [
  'breakfast',
  'lunch',
  'dinner',
  'desserts',
  'soups',
  'salads',
  'vegetarian',
  'quick-meals',
  'baking',
  'drinks',
];

export const CUISINES: CuisineSlug[] = [
  'italian',
  'french',
  'georgian',
  'japanese',
  'mexican',
  'indian',
  'mediterranean',
  'ukrainian',
  'latvian',
  'american',
];

export const DIFFICULTIES: Difficulty[] = ['easy', 'medium', 'hard'];

export const DIETS: Diet[] = [
  'vegetarian',
  'vegan',
  'gluten-free',
  'dairy-free',
  'pescatarian',
  'high-protein',
  'low-carb',
];

export const MEAL_TYPES: MealType[] = [
  'breakfast',
  'lunch',
  'dinner',
  'dessert',
  'snack',
  'appetizer',
  'drink',
  'side',
];

export const SORT_KEYS: SortKey[] = ['newest', 'fastest', 'easiest', 'rating'];

export const TIME_OPTIONS: number[] = [15, 30, 45, 60, 90];

/** Number of recipes shown per page in listing views. */
export const PAGE_SIZE = 9;

/** Cuisines highlighted on the home page. */
export const POPULAR_CUISINES: CuisineSlug[] = [
  'italian',
  'japanese',
  'mexican',
  'georgian',
  'indian',
  'french',
];

/** Reliable fallback image (Wikimedia Commons, public domain). */
export const FALLBACK_IMAGE =
  'https://upload.wikimedia.org/wikipedia/commons/a/a9/No_image_available.svg';

export const SITE_URL =
  process.env.NEXT_PUBLIC_SITE_URL?.replace(/\/$/, '') ?? 'http://localhost:3000';
