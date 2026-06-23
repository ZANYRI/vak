import type { Recipe } from "@/data/recipes";
import type { Locale } from "@/i18n/config";
import type { Dictionary } from "@/i18n/dictionary";
import { RecipeCard } from "./RecipeCard";

export function RecipeGrid({ recipes, locale, dictionary }: { recipes: Recipe[]; locale: Locale; dictionary: Dictionary }) {
  return <div className="recipe-grid">{recipes.map((recipe, index) => <div className="card-reveal" style={{ animationDelay: `${Math.min(index * 55, 360)}ms` }} key={recipe.id}><RecipeCard recipe={recipe} locale={locale} dictionary={dictionary} /></div>)}</div>;
}
