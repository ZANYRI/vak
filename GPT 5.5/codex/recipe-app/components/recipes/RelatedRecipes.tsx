import type { Recipe } from "@/data/recipes";
import type { Locale } from "@/i18n/config";
import type { Dictionary } from "@/i18n/dictionary";
import { RecipeGrid } from "./RecipeGrid";

export function RelatedRecipes({ title, recipes, locale, dictionary }: { title: string; recipes: Recipe[]; locale: Locale; dictionary: Dictionary }) {
  return <section className="shell related"><h2>{title}</h2><RecipeGrid recipes={recipes} locale={locale} dictionary={dictionary} /></section>;
}
