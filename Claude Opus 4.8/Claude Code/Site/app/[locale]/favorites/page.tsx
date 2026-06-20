import type { Metadata } from 'next';
import { getTranslations, setRequestLocale } from 'next-intl/server';
import { PageHeader } from '@/components/layout/PageHeader';
import { FavoritesList } from '@/components/recipes/FavoritesList';
import { getAllRecipes } from '@/lib/recipes';
import { buildMetadata } from '@/lib/seo';
import type { Locale } from '@/lib/types';

type Props = {
  params: Promise<{ locale: string }>;
};

export async function generateMetadata({ params }: Props): Promise<Metadata> {
  const { locale } = await params;
  const t = await getTranslations({ locale, namespace: 'Favorites' });
  return buildMetadata({
    locale: locale as Locale,
    path: '/favorites',
    title: t('title'),
    description: t('metaDescription'),
  });
}

export default async function FavoritesPage({ params }: Props) {
  const { locale } = await params;
  setRequestLocale(locale);
  const t = await getTranslations({ locale, namespace: 'Favorites' });

  return (
    <>
      <PageHeader title={t('title')} subtitle={t('subtitle')} />
      <FavoritesList recipes={getAllRecipes()} locale={locale as Locale} />
    </>
  );
}
