"use client";

import { useTranslations } from "next-intl";
import type { Recipe } from "@/lib/validation";
import { getRelatedRecipes } from "@/lib/recipes";
import { RecipeGrid } from "./RecipeGrid";

type RelatedRecipesProps = {
  recipe: Recipe;
  locale: "en" | "ru";
};

export function RelatedRecipes({ recipe, locale }: RelatedRecipesProps) {
  const t = useTranslations("recipe");
  const related = getRelatedRecipes(recipe, 3);

  if (related.length === 0) return null;

  return (
    <section>
      <h2 className="mb-4 text-2xl font-bold">{t("related")}</h2>
      <RecipeGrid recipes={related} locale={locale} />
    </section>
  );
}
