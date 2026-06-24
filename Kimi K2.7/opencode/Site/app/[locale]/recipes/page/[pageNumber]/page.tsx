import { notFound } from "next/navigation";
import { setRequestLocale } from "next-intl/server";
import { getTranslations } from "next-intl/server";
import { routing } from "@/i18n/routing";
import {
  filterAndSortRecipes,
  getAllRecipes,
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
  params: Promise<{ locale: string; pageNumber: string }>;
}) {
  const { locale, pageNumber } = await params;
  return createMetadata(locale as "en" | "ru", {
    titleKey: ["recipes", "title"],
    descriptionKey: ["recipes", "description"],
    path: `/recipes/page/${pageNumber}`,
  });
}

export function generateStaticParams() {
  const totalPages = Math.ceil(getAllRecipes().length / PAGE_SIZE);
  const params: Array<{ locale: string; pageNumber: string }> = [];
  for (const locale of routing.locales) {
    for (let page = 2; page <= totalPages; page++) {
      params.push({ locale, pageNumber: String(page) });
    }
  }
  return params;
}

export default async function RecipesPageNumber({
  params,
  searchParams,
}: {
  params: Promise<{ locale: string; pageNumber: string }>;
  searchParams: Promise<Record<string, string | string[] | undefined>>;
}) {
  const { locale, pageNumber } = await params;
  const query = await searchParams;
  setRequestLocale(locale);

  const page = Number(pageNumber);
  if (Number.isNaN(page) || page < 2) {
    notFound();
  }

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

  const { items, currentPage, totalPages, totalCount } = paginateRecipes(
    filtered,
    page,
    PAGE_SIZE,
  );
  if (currentPage !== page) {
    notFound();
  }

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
          <RecipePagination currentPage={currentPage} totalPages={totalPages} />
        </div>
      )}
    </div>
  );
}
