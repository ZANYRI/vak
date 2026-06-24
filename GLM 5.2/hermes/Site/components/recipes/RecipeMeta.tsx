import { useLocale, useTranslations } from "next-intl";
import { Badge } from "@/components/ui/Badge";
import {
  getCategoryName,
  getCuisineName,
  getDietName,
  getDifficultyName,
  getMealTypeName
} from "@/lib/taxonomy";
import type { Recipe } from "@/lib/types";

export interface RecipeMetaProps {
  recipe: Recipe;
}

export function RecipeMeta({ recipe }: RecipeMetaProps) {
  const locale = useLocale() as "en" | "ru";
  const t = useTranslations("recipe");

  const stats: Array<{ label: string; value: string; icon: string }> = [
    { label: t("prepTime"), value: `${recipe.prepTimeMinutes} min`, icon: "🔪" },
    { label: t("cookTime"), value: `${recipe.cookTimeMinutes} min`, icon: "🔥" },
    { label: t("totalTime"), value: `${recipe.totalTimeMinutes} min`, icon: "⏱" },
    { label: t("difficulty"), value: getDifficultyName(recipe.difficulty)[locale], icon: "📊" },
    { label: t("servings"), value: String(recipe.servings), icon: "🍽" },
    { label: t("cuisine"), value: getCuisineName(recipe.cuisine)[locale], icon: "🌍" }
  ];

  return (
    <section
      aria-label={t("servings")}
      className="grid grid-cols-2 gap-3 rounded-card border border-border bg-card p-4 shadow-soft sm:grid-cols-3 lg:grid-cols-6"
    >
      {stats.map((s) => (
        <div key={s.label} className="flex flex-col items-center text-center">
          <span aria-hidden="true" className="text-xl">
            {s.icon}
          </span>
          <span className="mt-1 text-xs uppercase tracking-wide text-muted">
            {s.label}
          </span>
          <span className="font-medium text-foreground">{s.value}</span>
        </div>
      ))}
      <div className="col-span-2 mt-2 flex flex-wrap items-center gap-2 sm:col-span-3 lg:col-span-6">
        <span className="text-xs uppercase tracking-wide text-muted">
          {t("category")}:
        </span>
        <Badge variant="secondary">{getCategoryName(recipe.category)[locale]}</Badge>
        <span className="text-xs uppercase tracking-wide text-muted">
          {t("mealType")}:
        </span>
        {recipe.mealType.map((m) => (
          <Badge key={m} variant="muted">
            {getMealTypeName(m)[locale]}
          </Badge>
        ))}
        {recipe.diet.length > 0 ? (
          <>
            <span className="text-xs uppercase tracking-wide text-muted">
              {t("diet")}:
            </span>
            {recipe.diet.map((d) => (
              <Badge key={d} variant="accent">
                {getDietName(d)[locale]}
              </Badge>
            ))}
          </>
        ) : null}
      </div>
    </section>
  );
}
