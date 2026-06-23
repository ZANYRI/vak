import { setRequestLocale, getTranslations, getLocale } from "next-intl/server";
import { RecipeGrid } from "@/components/recipes/RecipeGrid";
import { EmptyState } from "@/components/ui/EmptyState";
import { Link } from "@/i18n/navigation";
import { Button } from "@/components/ui/Button";
import { AnimatedSection } from "@/components/animations/AnimatedSection";
import { getRecipesByCuisine } from "@/lib/recipes";
import { getCuisineName } from "@/lib/taxonomy";
import { buildMetadata } from "@/lib/seo";
import type { Locale } from "@/lib/types";

export async function generateStaticParams() {
  const { CUISINES } = await import("@/lib/taxonomy");
  return CUISINES.map((c) => ({ cuisine: c.slug }));
}

export async function generateMetadata({
  params
}: {
  params: Promise<{ locale: string; cuisine: string }>;
}) {
  const { locale, cuisine } = await params;
  const name = getCuisineName(cuisine)[locale as Locale];
  return buildMetadata({
    locale: locale as Locale,
    title: name,
    description: `Authentic ${name} recipes.`,
    path: `cuisines/${cuisine}`
  });
}

export default async function CuisineDetailPage({
  params
}: {
  params: Promise<{ locale: string; cuisine: string }>;
}) {
  const { locale, cuisine } = await params;
  setRequestLocale(locale);
  const recipes = getRecipesByCuisine(cuisine);
  const t = await getTranslations({ locale, namespace: "cuisines" });
  const currentLocale = (await getLocale()) as Locale;
  const name = getCuisineName(cuisine)[currentLocale];

  return (
    <div className="container-page py-10">
      <AnimatedSection className="mb-8">
        <p className="text-sm font-medium uppercase tracking-wide text-primary">
          {t("title")}
        </p>
        <h1 className="mt-1 font-serif text-3xl font-semibold text-foreground sm:text-4xl">
          {name}
        </h1>
        <p className="mt-2 max-w-2xl text-muted">
          {t("intro", { name: name.toLowerCase() })}
        </p>
        <p className="mt-1 text-sm text-muted">
          {t("recipesCount", { count: recipes.length })}
        </p>
      </AnimatedSection>

      {recipes.length === 0 ? (
        <EmptyState
          title={t("empty")}
          action={
            <Link href="/cuisines">
              <Button variant="outline">{t("title")}</Button>
            </Link>
          }
          icon="🌍"
        />
      ) : (
        <RecipeGrid recipes={recipes} />
      )}
    </div>
  );
}
