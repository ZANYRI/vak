import type { Metadata } from 'next';
import { getTranslations, setRequestLocale } from 'next-intl/server';
import { PageHeader } from '@/components/layout/PageHeader';
import { RecipesListing } from '@/components/recipes/RecipesListing';
import { parseFilters } from '@/lib/validation';
import { buildMetadata } from '@/lib/seo';
import { getAllRecipes } from '@/lib/recipes';
import type { Locale } from '@/lib/types';

type SearchParams = Record<string, string | string[] | undefined>;

type Props = {
  params: Promise<{ locale: string }>;
  searchParams: Promise<SearchParams>;
};

export async function generateMetadata({
  params,
}: {
  params: Promise<{ locale: string }>;
}): Promise<Metadata> {
  const { locale } = await params;
  const t = await getTranslations({ locale, namespace: 'Recipes' });
  return buildMetadata({
    locale: locale as Locale,
    path: '/recipes',
    title: t('metaTitle'),
    description: t('metaDescription'),
  });
}

export default async function RecipesPage({ params, searchParams }: Props) {
  const { locale } = await params;
  setRequestLocale(locale);
  const sp = await searchParams;
  const filters = parseFilters(sp);

  const t = await getTranslations({ locale, namespace: 'Recipes' });

  return (
    <>
      <PageHeader
        title={t('title')}
        subtitle={t('description', { count: getAllRecipes().length })}
      />
      <RecipesListing locale={locale as Locale} filters={filters} page={1} />
    </>
  );
}
