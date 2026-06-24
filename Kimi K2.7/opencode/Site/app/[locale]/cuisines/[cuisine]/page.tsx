import { notFound } from "next/navigation";
import { setRequestLocale } from "next-intl/server";
import { getTranslations } from "next-intl/server";
import { routing } from "@/i18n/routing";
import { getAllRecipes, CUISINES } from "@/lib/recipes";
import { getSiteUrl } from "@/lib/seo";
import type { Metadata } from "next";
import { AnimatedSection } from "@/components/animations/AnimatedSection";
import { RecipeGrid } from "@/components/recipes/RecipeGrid";

export async function generateMetadata({
  params,
}: {
  params: Promise<{ locale: string; cuisine: string }>;
}): Promise<Metadata> {
  const { locale, cuisine } = await params;
  const t = await getTranslations({ locale, namespace: "cuisines" });
  const label = t(`labels.${cuisine}`);
  const url = `${getSiteUrl()}/${locale}/cuisines/${cuisine}`;
  return {
    title: label,
    description: t("intro", { cuisine: label }),
    alternates: { canonical: url },
  };
}

export function generateStaticParams() {
  const params: Array<{ locale: string; cuisine: string }> = [];
  for (const locale of routing.locales) {
    for (const cuisine of CUISINES) {
      params.push({ locale, cuisine });
    }
  }
  return params;
}

export default async function CuisinePage({
  params,
}: {
  params: Promise<{ locale: string; cuisine: string }>;
}) {
  const { locale, cuisine } = await params;
  if (!CUISINES.includes(cuisine as (typeof CUISINES)[number])) {
    notFound();
  }
  setRequestLocale(locale);
  const t = await getTranslations({ locale, namespace: "cuisines" });
  const label = t(`labels.${cuisine}`);
  const recipes = getAllRecipes().filter((r) => r.cuisine === cuisine);

  return (
    <div className="mx-auto max-w-7xl px-4 py-10 md:px-6">
      <AnimatedSection>
        <h1 className="mb-2 text-3xl font-bold md:text-4xl">{label}</h1>
        <p className="text-foreground/70 mb-8">
          {t("intro", { cuisine: label })}
        </p>
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
