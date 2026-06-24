"use client";

import { useMemo, useSyncExternalStore } from "react";
import type { Recipe } from "@/data/recipes";
import type { Locale } from "@/i18n/config";
import type { Dictionary } from "@/i18n/dictionary";
import { getFavoritesSnapshot, subscribeToFavorites } from "./FavoriteButton";
import { RecipeGrid } from "./RecipeGrid";
import { EmptyState } from "@/components/ui/EmptyState";

export function FavoritesClient({ recipes, locale, dictionary }: { recipes: Recipe[]; locale: Locale; dictionary: Dictionary }) {
  const snapshot = useSyncExternalStore(subscribeToFavorites, getFavoritesSnapshot, () => "[]");
  const slugs = useMemo(() => { try { const parsed: unknown = JSON.parse(snapshot); return Array.isArray(parsed) ? parsed.filter((item): item is string => typeof item === "string") : []; } catch { return []; } }, [snapshot]);
  const saved = recipes.filter((recipe) => slugs.includes(recipe.slug));
  return saved.length ? <RecipeGrid recipes={saved} locale={locale} dictionary={dictionary} /> : <EmptyState title={dictionary.favorites.emptyTitle} description={dictionary.favorites.emptyDescription} href={`/${locale}/recipes`} action={dictionary.favorites.browse} />;
}
