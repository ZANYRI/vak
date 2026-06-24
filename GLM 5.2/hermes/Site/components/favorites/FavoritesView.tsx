"use client";

import { useTranslations } from "next-intl";
import { useSyncExternalStore } from "react";
import { RecipeGrid } from "@/components/recipes/RecipeGrid";
import { EmptyState } from "@/components/ui/EmptyState";
import { Link } from "@/i18n/navigation";
import { Button } from "@/components/ui/Button";
import { AnimatedSection } from "@/components/animations/AnimatedSection";
import {
  getSnapshot,
  getServerSnapshot,
  subscribe
} from "@/lib/favorites";
import { getRecipesBySlug } from "@/lib/recipes";

export function FavoritesView() {
  const t = useTranslations();
  const slugs = useSyncExternalStore(subscribe, getSnapshot, getServerSnapshot);
  const recipes = getRecipesBySlug(slugs);

  return (
    <div className="container-page py-10">
      <AnimatedSection className="mb-8">
        <h1 className="font-serif text-3xl font-semibold text-foreground sm:text-4xl">
          {t("favorites.title")}
        </h1>
        <p className="mt-2 text-muted">{t("favorites.subtitle")}</p>
      </AnimatedSection>

      {recipes.length === 0 ? (
        <EmptyState
          title={t("favorites.empty")}
          description={t("favorites.emptyDescription")}
          icon="♡"
          action={
            <Link href="/recipes">
              <Button>{t("favorites.browseRecipes")}</Button>
            </Link>
          }
        />
      ) : (
        <RecipeGrid recipes={recipes} />
      )}
    </div>
  );
}
