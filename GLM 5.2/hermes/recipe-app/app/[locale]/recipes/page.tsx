import { setRequestLocale, getTranslations } from "next-intl/server";
import { RecipesListing } from "@/components/recipes/RecipesListing";
import { parseFilterParams, parsePageNumber } from "@/lib/validation";
import { buildMetadata } from "@/lib/seo";
import type { Locale } from "@/lib/types";

export async function generateMetadata({
  params
}: {
  params: Promise<{ locale: string }>;
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

export default async function RecipesPage({
  params,
  searchParams
}: {
  params: Promise<{ locale: string }>;
  searchParams: Promise<Record<string, string | string[] | undefined>>;
}) {
  const { locale } = await params;
  setRequestLocale(locale);
  const raw = await searchParams;
  const filters = parseFilterParams(raw);
  const page = parsePageNumber(typeof raw.page === "string" ? raw.page : undefined);

  return <RecipesListing filters={filters} page={page} />;
}
