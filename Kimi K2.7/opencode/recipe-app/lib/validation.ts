import { z } from "zod";

const localizedTextSchema = z.object({ en: z.string(), ru: z.string() });

const localizedStringArraySchema = z.object({
  en: z.array(z.string()),
  ru: z.array(z.string()),
});

export const recipeSchema = z.object({
  id: z.string(),
  slug: z.string(),
  title: localizedTextSchema,
  description: localizedTextSchema,
  image: z.string().url(),
  imageAlt: localizedTextSchema,
  cuisine: z.string(),
  category: z.string(),
  mealType: z.array(z.string()),
  diet: z.array(z.string()),
  difficulty: z.enum(["easy", "medium", "hard"]),
  prepTimeMinutes: z.number().nonnegative(),
  cookTimeMinutes: z.number().nonnegative(),
  totalTimeMinutes: z.number().nonnegative(),
  servings: z.number().positive(),
  rating: z.number().min(0).max(5),
  ingredients: localizedStringArraySchema,
  steps: localizedStringArraySchema,
  nutrition: z.object({
    calories: z.number(),
    protein: z.number(),
    fat: z.number(),
    carbs: z.number(),
  }),
  tips: localizedStringArraySchema,
  tags: z.array(z.string()),
});

export type Recipe = z.infer<typeof recipeSchema>;
export type LocalizedText = z.infer<typeof localizedTextSchema>;
