import { useTranslations } from "next-intl";
import type { Recipe } from "@/lib/validation";
import { formatDuration } from "@/lib/utils";

type RecipeMetaProps = {
  recipe: Recipe;
};

export function RecipeMeta({ recipe }: RecipeMetaProps) {
  const t = useTranslations("recipe");

  const items = [
    { label: t("prepTime"), value: formatDuration(recipe.prepTimeMinutes) },
    { label: t("cookTime"), value: formatDuration(recipe.cookTimeMinutes) },
    { label: t("totalTime"), value: formatDuration(recipe.totalTimeMinutes) },
    { label: t("servings"), value: recipe.servings.toString() },
    { label: t("difficulty"), value: recipe.difficulty },
    { label: t("cuisine"), value: recipe.cuisine },
    { label: t("category"), value: recipe.category },
  ];

  return (
    <div className="grid grid-cols-2 gap-4 sm:grid-cols-3 lg:grid-cols-4">
      {items.map((item) => (
        <div
          key={item.label}
          className="border-border bg-card rounded-md border p-3 text-center"
        >
          <p className="text-foreground/60 text-xs tracking-wide uppercase">
            {item.label}
          </p>
          <p className="mt-1 font-semibold capitalize">{item.value}</p>
        </div>
      ))}
    </div>
  );
}
