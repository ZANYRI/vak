import type { Metadata } from 'next';
import { notFound } from 'next/navigation';
import { getTranslations, setRequestLocale } from 'next-intl/server';
import { PageHeader } from '@/components/layout/PageHeader';
import { RecipeGrid } from '@/components/recipes/RecipeGrid';
import { CATEGORIES } from '@/lib/constants';
import { getRecipesByCategory } from '@/lib/recipes';
import { routing } from '@/i18n/routing';
import { buildMetadata } from '@/lib/seo';
import type { CategorySlug, Locale } from '@/lib/types';

type Props = {
  params: Promise<{ locale: string; category: string }>;
};

export function generateStaticParams() {
  const params: { locale: string; category: string }[] = [];
  for (const locale of routing.locales) {
    for (const category of CATEGORIES) {
      params.push({ locale, category });
    }
  }
  return params;
}

function isCategory(value: string): value is CategorySlug {
  return (CATEGORIES as string[]).includes(value);
}

export async function generateMetadata({ params }: Props): Promise<Metadata> {
  const { locale, category } = await params;
  if (!isCategory(category)) return {};
  const tName = await getTranslations({ locale, namespace: 'Categories.names' });
  const t = await getTranslations({ locale, namespace: 'Categories' });
  return buildMetadata({
    locale: locale as Locale,
    path: `/categories/${category}`,
    title: tName(category),
    description: t('inCategory', { name: tName(category) }),
  });
}

export default async function CategoryPage({ params }: Props) {
  const { locale, category } = await params;
  setRequestLocale(locale);
  if (!isCategory(category)) {
    notFound();
  }

  const t = await getTranslations({ locale, namespace: 'Categories' });
  const tName = await getTranslations({ locale, namespace: 'Categories.names' });
  const recipes = getRecipesByCategory(category);

  return (
    <>
      <PageHeader
        eyebrow={t('title')}
        title={tName(category)}
        subtitle={t('recipeCount', { count: recipes.length })}
      />
      <RecipeGrid recipes={recipes} locale={locale as Locale} />
    </>
  );
}
