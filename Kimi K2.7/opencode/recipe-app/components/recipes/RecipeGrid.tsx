import type { Recipe } from "@/lib/validation";
import { RecipeCard } from "./RecipeCard";

type RecipeGridProps = {
  recipes: Recipe[];
  locale: "en" | "ru";
};

export function RecipeGrid({ recipes, locale }: RecipeGridProps) {
  return (
    <div className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
      {recipes.map((recipe, index) => (
        <RecipeCard
          key={recipe.id}
          recipe={recipe}
          locale={locale}
          index={index}
        />
      ))}
    </div>
  );
}
