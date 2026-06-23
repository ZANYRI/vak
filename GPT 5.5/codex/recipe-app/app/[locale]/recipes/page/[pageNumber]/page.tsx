import { notFound } from "next/navigation";
import { RecipeListing } from "@/components/recipes/RecipeListing";
import { isLocale } from "@/i18n/config";

export default async function NumberedRecipesPage({ params, searchParams }: { params: Promise<{ locale: string; pageNumber: string }>; searchParams: Promise<Record<string, string | string[] | undefined>> }) { const { locale, pageNumber } = await params; if (!isLocale(locale)) notFound(); return <RecipeListing locale={locale} page={Number(pageNumber)} searchParams={await searchParams} />; }
