import { z } from "zod";

const difficultySchema = z.enum(["easy", "medium", "hard"]);

const sortSchema = z.enum(["newest", "fastest", "easiest", "rating"]);

/**
 * Schema for the query-string filter state shared by the recipes listing,
 * search bar, and filter panel. Accepts empty strings (treated as "unset")
 * so a bare `<select>` option with value "" validates cleanly.
 */
export const filterParamsSchema = z.object({
  search: z.string().trim().optional().catch(""),
  cuisine: z.string().trim().optional().catch(""),
  category: z.string().trim().optional().catch(""),
  difficulty: z
    .union([difficultySchema, z.literal("")])
    .optional()
    .catch(""),
  diet: z.string().trim().optional().catch(""),
  mealType: z.string().trim().optional().catch(""),
  maxTime: z
    .string()
    .trim()
    .optional()
    .catch("")
    .refine(
      (v) => v === "" || v === undefined || (!Number.isNaN(Number(v)) && Number(v) > 0),
      "maxTime must be a positive number"
    ),
  sort: z.union([sortSchema, z.literal("")]).optional().catch("")
});

export type ParsedFilterParams = z.infer<typeof filterParamsSchema>;

export function parseFilterParams(raw: Record<string, string | string[] | undefined>) {
  const flat: Record<string, string | undefined> = {};
  for (const key of Object.keys(raw)) {
    const value = raw[key];
    if (Array.isArray(value)) {
      flat[key] = value[0];
    } else {
      flat[key] = value;
    }
  }
  return filterParamsSchema.parse(flat);
}

const pageNumberSchema = z.coerce.number().int().min(1).max(10000);

export function parsePageNumber(value: string | undefined): number {
  const result = pageNumberSchema.safeParse(value);
  return result.success ? result.data : 1;
}
