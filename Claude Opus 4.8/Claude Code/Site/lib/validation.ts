import { z } from 'zod';
import {
  CATEGORIES,
  CUISINES,
  DIETS,
  DIFFICULTIES,
  MEAL_TYPES,
  SORT_KEYS,
} from './constants';
import type { RecipeFilters } from './types';

const enumFrom = <T extends string>(values: readonly T[]) =>
  z.enum(values as [T, ...T[]]);

/** Zod schema validating a single recipe record (used to guard the dataset). */
export const recipeSchema = z.object({
  id: z.string().min(1),
  slug: z
    .string()
    .regex(/^[a-z0-9]+(?:-[a-z0-9]+)*$/, 'slug must be kebab-case'),
  title: z.object({ en: z.string().min(1), ru: z.string().min(1) }),
  description: z.object({ en: z.string().min(1), ru: z.string().min(1) }),
  image: z.string().url(),
  imageAlt: z.object({ en: z.string().min(1), ru: z.string().min(1) }),
  cuisine: enumFrom(CUISINES),
  category: enumFrom(CATEGORIES),
  mealType: z.array(enumFrom(MEAL_TYPES)).min(1),
  diet: z.array(enumFrom(DIETS)),
  difficulty: enumFrom(DIFFICULTIES),
  prepTimeMinutes: z.number().int().nonnegative(),
  cookTimeMinutes: z.number().int().nonnegative(),
  totalTimeMinutes: z.number().int().positive(),
  servings: z.number().int().positive(),
  rating: z.number().min(0).max(5),
  publishedAt: z.string().regex(/^\d{4}-\d{2}-\d{2}$/),
  ingredients: z.object({
    en: z.array(z.string().min(1)).min(1),
    ru: z.array(z.string().min(1)).min(1),
  }),
  steps: z.object({
    en: z.array(z.string().min(1)).min(1),
    ru: z.array(z.string().min(1)).min(1),
  }),
  nutrition: z.object({
    calories: z.number().nonnegative(),
    protein: z.number().nonnegative(),
    fat: z.number().nonnegative(),
    carbs: z.number().nonnegative(),
  }),
  tips: z.object({
    en: z.array(z.string()),
    ru: z.array(z.string()),
  }),
  tags: z.array(z.string()),
});

export type RecipeInput = z.infer<typeof recipeSchema>;

/** Schema for the raw URL search params on the recipes listing page. */
export const filterParamsSchema = z.object({
  search: z.string().trim().min(1).optional(),
  cuisine: enumFrom(CUISINES).optional(),
  category: enumFrom(CATEGORIES).optional(),
  difficulty: enumFrom(DIFFICULTIES).optional(),
  diet: enumFrom(DIETS).optional(),
  mealType: enumFrom(MEAL_TYPES).optional(),
  maxTime: z.coerce.number().int().positive().optional(),
  sort: enumFrom(SORT_KEYS).optional(),
});

type RawSearchParams = Record<string, string | string[] | undefined>;

const firstValue = (value: string | string[] | undefined): string | undefined =>
  Array.isArray(value) ? value[0] : value;

/**
 * Safely coerce arbitrary URL search params into typed RecipeFilters.
 * Invalid values are dropped rather than throwing, so a bad URL never 500s.
 */
export function parseFilters(params: RawSearchParams): RecipeFilters {
  const normalized: RawSearchParams = {};
  for (const key of Object.keys(params)) {
    normalized[key] = firstValue(params[key]);
  }
  const result = filterParamsSchema.safeParse(normalized);
  return result.success ? result.data : {};
}
