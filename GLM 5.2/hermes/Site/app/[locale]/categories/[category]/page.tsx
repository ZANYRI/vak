import { setRequestLocale, getTranslations, getLocale } from "next-intl/server";
import { RecipeGrid } from "@/components/recipes/RecipeGrid";
import { EmptyState } from "@/components/ui/EmptyState";
import { Link } from "@/i18n/navigation";
import { Button } from "@/components/ui/Button";
import { AnimatedSection } from "@/components/animations/AnimatedSection";
import { getRecipesByCategory } from "@/lib/recipes";
import { getCategoryName } from "@/lib/taxonomy";
import { buildMetadata } from "@/lib/seo";
import type { Locale } from "@/lib/types";

export async function generateStaticParams() {
  const { CATEGORIES } = await import("@/lib/taxonomy");
  return CATEGORIES.map((c) => ({ category: c.slug }));
}

export async function generateMetadata({
  params
}: {
  params: Promise<{ locale: string; category: string }>;
}) {
  const { locale, category } = await params;
  const name = getCategoryName(category)[locale as Locale];
  return buildMetadata({
    locale: locale as Locale,
    title: name,
    description: `Recipes in the ${name} category.`,
    path: `categories/${category}`
  });
}

export default async function CategoryDetailPage({
  params
}: {
  params: Promise<{ locale: string; category: string }>;
}) {
  const { locale, category } = await params;
  setRequestLocale(locale);
  const recipes = getRecipesByCategory(category);
  const t = await getTranslations({ locale, namespace: "categories" });
  const currentLocale = (await getLocale()) as Locale;
  const name = getCategoryName(category)[currentLocale];

  return (
    <div className="container-page py-10">
      <AnimatedSection className="mb-8">
        <p className="text-sm font-medium uppercase tracking-wide text-primary">
          {t("title")}
        </p>
        <h1 className="mt-1 font-serif text-3xl font-semibold text-foreground sm:text-4xl">
          {name}
        </h1>
        <p className="mt-2 text-muted">
          {t("recipesCount", { count: recipes.length })}
        </p>
      </AnimatedSection>

      {recipes.length === 0 ? (
        <EmptyState
          title={t("empty")}
          action={
            <Link href="/categories">
              <Button variant="outline">{t("title")}</Button>
            </Link>
          }
          icon="🏷"
        />
      ) : (
        <RecipeGrid recipes={recipes} />
      )}
    </div>
  );
}
