"use client";

import { useTranslations } from "next-intl";
import { RecipeCard } from "./RecipeCard";
import { getRelatedRecipes } from "@/lib/recipes";
import type { Recipe } from "@/lib/types";

export function RelatedRecipes({ recipe }: { recipe: Recipe }) {
  const t = useTranslations("recipe");
  const related = getRelatedRecipes(recipe, 3);

  if (related.length === 0) return null;

  return (
    <section aria-labelledby="related-heading" className="mt-12">
      <h2 id="related-heading" className="font-serif text-2xl font-semibold text-foreground">
        {t("related")}
      </h2>
      <ul className="mt-6 grid gap-6 sm:grid-cols-2 lg:grid-cols-3" role="list">
        {related.map((r, i) => (
          <li key={r.slug} className="relative">
            <RecipeCard recipe={r} index={i} />
          </li>
        ))}
      </ul>
    </section>
  );
}
