import type { Difficulty, LocalizedText, TaxonomyEntry } from "./types";

export const CUISINES: TaxonomyEntry[] = [
  { slug: "italian", name: { en: "Italian", ru: "Итальянская" } },
  { slug: "french", name: { en: "French", ru: "Французская" } },
  { slug: "georgian", name: { en: "Georgian", ru: "Грузинская" } },
  { slug: "japanese", name: { en: "Japanese", ru: "Японская" } },
  { slug: "mexican", name: { en: "Mexican", ru: "Мексиканская" } },
  { slug: "indian", name: { en: "Indian", ru: "Индийская" } },
  { slug: "mediterranean", name: { en: "Mediterranean", ru: "Средиземноморская" } },
  { slug: "ukrainian", name: { en: "Ukrainian", ru: "Украинская" } },
  { slug: "latvian", name: { en: "Latvian", ru: "Латышская" } },
  { slug: "american", name: { en: "American", ru: "Американская" } }
];

export const CATEGORIES: TaxonomyEntry[] = [
  { slug: "breakfast", name: { en: "Breakfast", ru: "Завтрак" } },
  { slug: "lunch", name: { en: "Lunch", ru: "Обед" } },
  { slug: "dinner", name: { en: "Dinner", ru: "Ужин" } },
  { slug: "desserts", name: { en: "Desserts", ru: "Десерты" } },
  { slug: "soups", name: { en: "Soups", ru: "Супы" } },
  { slug: "salads", name: { en: "Salads", ru: "Салаты" } },
  { slug: "vegetarian", name: { en: "Vegetarian", ru: "Вегетарианское" } },
  { slug: "quick-meals", name: { en: "Quick meals", ru: "Быстрые блюда" } },
  { slug: "baking", name: { en: "Baking", ru: "Выпечка" } },
  { slug: "drinks", name: { en: "Drinks", ru: "Напитки" } }
];

export const DIETS: TaxonomyEntry[] = [
  { slug: "vegetarian", name: { en: "Vegetarian", ru: "Вегетарианская" } },
  { slug: "vegan", name: { en: "Vegan", ru: "Веганская" } },
  { slug: "gluten-free", name: { en: "Gluten-free", ru: "Без глютена" } },
  { slug: "dairy-free", name: { en: "Dairy-free", ru: "Без молока" } },
  { slug: "pescatarian", name: { en: "Pescatarian", ru: "Пескетарианская" } },
  { slug: "keto", name: { en: "Keto", ru: "Кето" } },
  { slug: "low-carb", name: { en: "Low-carb", ru: "Низкоуглеводная" } },
  { slug: "nut-free", name: { en: "Nut-free", ru: "Без орехов" } }
];

export const MEAL_TYPES: TaxonomyEntry[] = [
  { slug: "breakfast", name: { en: "Breakfast", ru: "Завтрак" } },
  { slug: "lunch", name: { en: "Lunch", ru: "Обед" } },
  { slug: "dinner", name: { en: "Dinner", ru: "Ужин" } },
  { slug: "dessert", name: { en: "Dessert", ru: "Десерт" } },
  { slug: "snack", name: { en: "Snack", ru: "Перекус" } },
  { slug: "drink", name: { en: "Drink", ru: "Напиток" } }
];

export const DIFFICULTIES: { slug: Difficulty; name: { en: string; ru: string } }[] = [
  { slug: "easy", name: { en: "Easy", ru: "Лёгкий" } },
  { slug: "medium", name: { en: "Medium", ru: "Средний" } },
  { slug: "hard", name: { en: "Hard", ru: "Сложный" } }
];

function findName(entries: TaxonomyEntry[], slug: string): LocalizedText {
  const found = entries.find((e) => e.slug === slug);
  return found?.name ?? { en: slug, ru: slug };
}

export function getCuisineName(slug: string) {
  return findName(CUISINES, slug);
}
export function getCategoryName(slug: string) {
  return findName(CATEGORIES, slug);
}
export function getDietName(slug: string) {
  return findName(DIETS, slug);
}
export function getMealTypeName(slug: string) {
  return findName(MEAL_TYPES, slug);
}
export function getDifficultyName(slug: Difficulty) {
  return findName(
    DIFFICULTIES.map((d) => ({ slug: d.slug, name: d.name })),
    slug
  );
}

export function isValidCuisine(slug: string) {
  return CUISINES.some((c) => c.slug === slug);
}
export function isValidCategory(slug: string) {
  return CATEGORIES.some((c) => c.slug === slug);
}
