import type { Metadata } from 'next';
import { notFound } from 'next/navigation';
import { getTranslations, setRequestLocale } from 'next-intl/server';
import { PageHeader } from '@/components/layout/PageHeader';
import { RecipesListing } from '@/components/recipes/RecipesListing';
import { parseFilters } from '@/lib/validation';
import { parsePage } from '@/lib/pagination';
import { buildMetadata } from '@/lib/seo';
import { getAllRecipes, getTotalPages } from '@/lib/recipes';
import { routing } from '@/i18n/routing';
import type { Locale } from '@/lib/types';

type SearchParams = Record<string, string | string[] | undefined>;

type Props = {
  params: Promise<{ locale: string; pageNumber: string }>;
  searchParams: Promise<SearchParams>;
};

/** Pre-render unfiltered pages 2..N for each locale. */
export function generateStaticParams() {
  const totalPages = getTotalPages();
  const params: { locale: string; pageNumber: string }[] = [];
  for (const locale of routing.locales) {
    for (let p = 2; p <= totalPages; p++) {
      params.push({ locale, pageNumber: String(p) });
    }
  }
  return params;
}

export async function generateMetadata({
  params,
}: {
  params: Promise<{ locale: string; pageNumber: string }>;
}): Promise<Metadata> {
  const { locale, pageNumber } = await params;
  const t = await getTranslations({ locale, namespace: 'Recipes' });
  const page = parsePage(pageNumber);
  return buildMetadata({
    locale: locale as Locale,
    path: `/recipes/page/${page}`,
    title: `${t('metaTitle')} — ${t('title')} ${page}`,
    description: t('metaDescription'),
  });
}

export default async function RecipesPaginatedPage({ params, searchParams }: Props) {
  const { locale, pageNumber } = await params;
  setRequestLocale(locale);

  const page = parsePage(pageNumber);
  // Page 1 lives at /recipes — reject it (and bad values) here.
  if (page < 2 || String(page) !== pageNumber) {
    notFound();
  }

  const sp = await searchParams;
  const filters = parseFilters(sp);
  const t = await getTranslations({ locale, namespace: 'Recipes' });

  return (
    <>
      <PageHeader
        title={t('title')}
        subtitle={t('description', { count: getAllRecipes().length })}
      />
      <RecipesListing locale={locale as Locale} filters={filters} page={page} />
    </>
  );
}
