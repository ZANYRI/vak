import { notFound } from "next/navigation";
import { Suspense } from "react";
import type { Locale } from "@/i18n/config";
import { getDictionary } from "@/i18n/dictionary";
import { filterRecipes, perPage, type RecipeQuery } from "@/lib/recipes";
import { firstValue, recipeQuerySchema } from "@/lib/validation";
import { EmptyState } from "@/components/ui/EmptyState";
import { RecipeGrid } from "./RecipeGrid";
import { RecipePagination } from "./RecipePagination";
import { RecipeFilters } from "@/components/filters/RecipeFilters";

type SearchParams = Record<string, string | string[] | undefined>;

export function RecipeListing({ locale, page, searchParams }: { locale: Locale; page: number; searchParams: SearchParams }) {
  if (!Number.isInteger(page) || page < 1) notFound();
  const dictionary = getDictionary(locale);
  const raw = Object.fromEntries(Object.entries(searchParams).map(([key, value]) => [key, firstValue(value)]));
  const parsed = recipeQuerySchema.safeParse(raw);
  const query: RecipeQuery = parsed.success ? Object.fromEntries(Object.entries(parsed.data).map(([key, value]) => [key, value === undefined ? undefined : String(value)])) : {};
  const filtered = filterRecipes(query, locale);
  const totalPages = Math.max(1, Math.ceil(filtered.length / perPage));
  if (page > totalPages) notFound();
  const result = filtered.slice((page - 1) * perPage, page * perPage);
  return <>
    <section className="page-intro shell"><p className="eyebrow">{dictionary.nav.recipes}</p><h1>{dictionary.recipes.title}</h1><p>{dictionary.recipes.description}</p></section>
    <section className="shell listing-content"><Suspense fallback={<div className="search-field">{dictionary.common.loading}</div>}><RecipeFilters dictionary={dictionary} /></Suspense><p className="result-count">{filtered.length} {dictionary.recipes.results}</p>{result.length ? <RecipeGrid recipes={result} locale={locale} dictionary={dictionary} /> : <EmptyState title={dictionary.recipes.noResults} description={dictionary.recipes.description} href={`/${locale}/recipes`} action={dictionary.recipes.clear} />}<RecipePagination locale={locale} page={page} totalPages={totalPages} query={query} dictionary={dictionary} /></section>
  </>;
}
