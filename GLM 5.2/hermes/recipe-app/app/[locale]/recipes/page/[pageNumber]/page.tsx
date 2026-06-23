import { setRequestLocale, getTranslations } from "next-intl/server";
import { RecipesListing } from "@/components/recipes/RecipesListing";
import { parseFilterParams } from "@/lib/validation";
import { buildMetadata } from "@/lib/seo";
import type { Locale } from "@/lib/types";

export async function generateMetadata({
  params
}: {
  params: Promise<{ locale: string; pageNumber: string }>;
}) {
  const { locale } = await params;
  const t = await getTranslations({ locale, namespace: "seo" });
  return buildMetadata({
    locale: locale as Locale,
    title: "Recipes",
    description: t("recipesDescription"),
    path: "recipes"
  });
}

export default async function RecipesPagedPage({
  params,
  searchParams
}: {
  params: Promise<{ locale: string; pageNumber: string }>;
  searchParams: Promise<Record<string, string | string[] | undefined>>;
}) {
  const { locale, pageNumber } = await params;
  setRequestLocale(locale);
  const raw = await searchParams;
  const filters = parseFilterParams(raw);
  const page = Math.max(1, parseInt(pageNumber, 10) || 1);

  return <RecipesListing filters={filters} page={page} />;
}
