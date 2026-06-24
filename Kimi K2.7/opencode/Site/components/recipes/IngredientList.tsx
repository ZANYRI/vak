"use client";

import { useTranslations } from "next-intl";
import type { Recipe } from "@/lib/validation";

type IngredientListProps = {
  recipe: Recipe;
  locale: "en" | "ru";
};

export function IngredientList({ recipe, locale }: IngredientListProps) {
  const t = useTranslations("recipe");

  return (
    <section className="border-border bg-card rounded-xl border p-6">
      <h2 className="mb-4 text-xl font-bold">{t("ingredients")}</h2>
      <ul className="space-y-2">
        {recipe.ingredients[locale].map((ingredient, index) => (
          <li key={index} className="flex items-start gap-3">
            <input
              type="checkbox"
              id={`ingredient-${index}`}
              className="accent-primary mt-1 h-4 w-4"
            />
            <label
              htmlFor={`ingredient-${index}`}
              className="text-foreground/90"
            >
              {ingredient}
            </label>
          </li>
        ))}
      </ul>
    </section>
  );
}
