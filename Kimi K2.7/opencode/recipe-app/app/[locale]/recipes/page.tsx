import { setRequestLocale } from "next-intl/server";
import { getTranslations } from "next-intl/server";
import { routing } from "@/i18n/routing";
import {
  filterAndSortRecipes,
  paginateRecipes,
  PAGE_SIZE,
} from "@/lib/recipes";
import { createMetadata } from "@/lib/seo";
import { RecipeFilters } from "@/components/filters/RecipeFilters";
import { RecipePagination } from "@/components/filters/RecipePagination";
import { RecipeGrid } from "@/components/recipes/RecipeGrid";
import { EmptyState } from "@/components/recipes/EmptyState";
import { AnimatedSection } from "@/components/animations/AnimatedSection";

export async function generateMetadata({
  params,
}: {
  params: Promise<{ locale: string }>;
}) {
  const { locale } = await params;
  return createMetadata(locale as "en" | "ru", {
    titleKey: ["recipes", "title"],
    descriptionKey: ["recipes", "description"],
    path: "/recipes",
  });
}

export function generateStaticParams() {
  return routing.locales.map((locale) => ({ locale }));
}

export default async function RecipesPage({
  params,
  searchParams,
}: {
  params: Promise<{ locale: string }>;
  searchParams: Promise<Record<string, string | string[] | undefined>>;
}) {
  const { locale } = await params;
  const query = await searchParams;
  setRequestLocale(locale);

  const t = await getTranslations({ locale, namespace: "recipes" });

  const filtered = filterAndSortRecipes(locale as "en" | "ru", {
    search: typeof query.search === "string" ? query.search : undefined,
    cuisine: typeof query.cuisine === "string" ? query.cuisine : undefined,
    category: typeof query.category === "string" ? query.category : undefined,
    difficulty:
      typeof query.difficulty === "string" ? query.difficulty : undefined,
    diet: typeof query.diet === "string" ? query.diet : undefined,
    mealType: typeof query.mealType === "string" ? query.mealType : undefined,
    time: typeof query.time === "string" ? query.time : undefined,
    sort: typeof query.sort === "string" ? query.sort : undefined,
  });

  const { items, totalPages, totalCount } = paginateRecipes(
    filtered,
    1,
    PAGE_SIZE,
  );

  return (
    <div className="mx-auto max-w-7xl px-4 py-10 md:px-6">
      <AnimatedSection>
        <h1 className="mb-2 text-3xl font-bold md:text-4xl">{t("title")}</h1>
        <p className="text-foreground/70 mb-6">{t("description")}</p>
      </AnimatedSection>

      <AnimatedSection className="mb-8">
        <RecipeFilters />
      </AnimatedSection>

      <AnimatedSection>
        <p className="text-foreground/70 mb-4 text-sm">
          {t("results", { count: totalCount })}
        </p>
        {items.length > 0 ? (
          <RecipeGrid recipes={items} locale={locale as "en" | "ru"} />
        ) : (
          <EmptyState title={t("noResults")} />
        )}
      </AnimatedSection>

      {totalPages > 1 && (
        <div className="mt-10">
          <RecipePagination currentPage={1} totalPages={totalPages} />
        </div>
      )}
    </div>
  );
}
