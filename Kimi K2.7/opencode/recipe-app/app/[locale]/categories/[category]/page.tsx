import { notFound } from "next/navigation";
import { setRequestLocale } from "next-intl/server";
import { getTranslations } from "next-intl/server";
import { routing } from "@/i18n/routing";
import { getAllRecipes, CATEGORIES } from "@/lib/recipes";
import { getSiteUrl } from "@/lib/seo";
import type { Metadata } from "next";
import { AnimatedSection } from "@/components/animations/AnimatedSection";
import { RecipeGrid } from "@/components/recipes/RecipeGrid";

export async function generateMetadata({
  params,
}: {
  params: Promise<{ locale: string; category: string }>;
}): Promise<Metadata> {
  const { locale, category } = await params;
  const t = await getTranslations({ locale, namespace: "categories" });
  const label = t(`labels.${category}`);
  const url = `${getSiteUrl()}/${locale}/categories/${category}`;
  return {
    title: label,
    description: t("description"),
    alternates: { canonical: url },
  };
}

export function generateStaticParams() {
  const params: Array<{ locale: string; category: string }> = [];
  for (const locale of routing.locales) {
    for (const category of CATEGORIES) {
      params.push({ locale, category });
    }
  }
  return params;
}

export default async function CategoryPage({
  params,
}: {
  params: Promise<{ locale: string; category: string }>;
}) {
  const { locale, category } = await params;
  if (!CATEGORIES.includes(category as (typeof CATEGORIES)[number])) {
    notFound();
  }
  setRequestLocale(locale);
  const t = await getTranslations({ locale, namespace: "categories" });
  const recipes = getAllRecipes().filter((r) => r.category === category);

  return (
    <div className="mx-auto max-w-7xl px-4 py-10 md:px-6">
      <AnimatedSection>
        <h1 className="mb-2 text-3xl font-bold md:text-4xl">
          {t(`labels.${category}`)}
        </h1>
        <p className="text-foreground/70 mb-8">{t("description")}</p>
      </AnimatedSection>
      <AnimatedSection>
        {recipes.length > 0 ? (
          <RecipeGrid recipes={recipes} locale={locale as "en" | "ru"} />
        ) : (
          <p>No recipes found.</p>
        )}
      </AnimatedSection>
    </div>
  );
}
