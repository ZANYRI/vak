import { z } from "zod";

const optionalText = z.string().trim().max(80).optional();

export const recipeQuerySchema = z.object({
  search: optionalText,
  cuisine: optionalText,
  category: optionalText,
  difficulty: optionalText,
  diet: optionalText,
  meal: optionalText,
  time: z.coerce.number().int().min(0).max(360).optional(),
  sort: z.enum(["newest", "fastest", "easiest", "rated"]).optional()
});

export function firstValue(value: string | string[] | undefined) {
  return Array.isArray(value) ? value[0] : value;
}
