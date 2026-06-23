"use client";

import { useTranslations, useLocale } from "next-intl";
import { useFavorites } from "@/components/providers/FavoritesProvider";
import { getAllRecipes } from "@/lib/recipes";
import { RecipeGrid } from "@/components/recipes/RecipeGrid";
import { EmptyState } from "@/components/recipes/EmptyState";
import { AnimatedSection } from "@/components/animations/AnimatedSection";

export default function FavoritesPage() {
  const t = useTranslations("favorites");
  const locale = useLocale();
  const { favorites } = useFavorites();

  const recipes = getAllRecipes().filter((r) => favorites.includes(r.id));

  return (
    <div className="mx-auto max-w-7xl px-4 py-10 md:px-6">
      <AnimatedSection>
        <h1 className="mb-2 text-3xl font-bold md:text-4xl">{t("title")}</h1>
        <p className="text-foreground/70 mb-8">{t("description")}</p>
      </AnimatedSection>
      <AnimatedSection>
        {recipes.length === 0 ? (
          <EmptyState
            title={t("emptyTitle")}
            description={t("emptyDescription")}
          />
        ) : (
          <RecipeGrid recipes={recipes} locale={locale as "en" | "ru"} />
        )}
      </AnimatedSection>
    </div>
  );
}
