import { Suspense } from "react";
import { useLocale, useTranslations } from "next-intl";
import { RecipeGrid } from "./RecipeGrid";
import { RecipeFilters } from "@/components/filters/RecipeFilters";
import { RecipePagination } from "@/components/filters/RecipePagination";
import { SearchBar } from "@/components/filters/SearchBar";
import { EmptyState } from "@/components/ui/EmptyState";
import { Button } from "@/components/ui/Button";
import { AnimatedSection } from "@/components/animations/AnimatedSection";
import { queryRecipes } from "@/lib/recipes";
import type { FilterParams, Locale } from "@/lib/types";

export const RECIPES_BASE_PATH = "/recipes";

export function RecipesListing({
  filters,
  page
}: {
  filters: FilterParams;
  page: number;
}) {
  const t = useTranslations();
  const locale = useLocale() as Locale;
  const result = queryRecipes(filters, page, locale);
  const queryString = buildQueryString(filters);

  return (
    <div className="container-page py-10">
      <AnimatedSection className="mb-6">
        <h1 className="font-serif text-3xl font-semibold text-foreground sm:text-4xl">
          {t("recipes.title")}
        </h1>
        <p className="mt-2 text-muted">{t("recipes.subtitle")}</p>
      </AnimatedSection>

      <div className="mb-6">
        <SearchBar
          basePath={RECIPES_BASE_PATH}
          initialQuery={filters.search ?? ""}
          placeholder={t("filters.searchPlaceholder")}
        />
      </div>

      <div className="grid gap-6 lg:grid-cols-[280px_1fr]">
        <Suspense
          fallback={
            <div className="h-64 rounded-card border border-border bg-card" />
          }
        >
          <RecipeFilters basePath={RECIPES_BASE_PATH} />
        </Suspense>

        <div>
          <p className="mb-4 text-sm text-muted" aria-live="polite">
            {t("recipes.showing", {
              from: result.from,
              to: result.to,
              total: result.total
            })}
          </p>

          {result.items.length === 0 ? (
            <EmptyState
              title={t("recipes.empty")}
              description={t("recipes.emptyHint")}
              action={
                <a href={RECIPES_BASE_PATH}>
                  <Button variant="outline">{t("recipes.clearFilters")}</Button>
                </a>
              }
              icon="🍽"
            />
          ) : (
            <RecipeGrid recipes={result.items} />
          )}

          <Suspense fallback={null}>
            <RecipePagination
              basePath={RECIPES_BASE_PATH}
              queryString={queryString}
              currentPage={result.page}
              totalPages={result.totalPages}
            />
          </Suspense>
        </div>
      </div>
    </div>
  );
}

function buildQueryString(filters: FilterParams): string {
  const params = new URLSearchParams();
  for (const [k, v] of Object.entries(filters)) {
    if (v && k !== "page") params.set(k, String(v));
  }
  const qs = params.toString();
  return qs ? `?${qs}` : "";
}
