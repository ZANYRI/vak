import { setRequestLocale, getTranslations } from "next-intl/server";
import { FavoritesView } from "@/components/favorites/FavoritesView";
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
    title: "Favorites",
    description: t("favoritesDescription"),
    path: "favorites"
  });
}

export default async function FavoritesPage({
  params
}: {
  params: Promise<{ locale: string }>;
}) {
  const { locale } = await params;
  setRequestLocale(locale);
  return <FavoritesView />;
}
