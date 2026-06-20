import type { Metadata } from 'next';
import { notFound } from 'next/navigation';
import { getTranslations, setRequestLocale } from 'next-intl/server';
import { PageHeader } from '@/components/layout/PageHeader';
import { RecipeGrid } from '@/components/recipes/RecipeGrid';
import { CUISINES } from '@/lib/constants';
import { getRecipesByCuisine } from '@/lib/recipes';
import { routing } from '@/i18n/routing';
import { buildMetadata } from '@/lib/seo';
import type { CuisineSlug, Locale } from '@/lib/types';

type Props = {
  params: Promise<{ locale: string; cuisine: string }>;
};

export function generateStaticParams() {
  const params: { locale: string; cuisine: string }[] = [];
  for (const locale of routing.locales) {
    for (const cuisine of CUISINES) {
      params.push({ locale, cuisine });
    }
  }
  return params;
}

function isCuisine(value: string): value is CuisineSlug {
  return (CUISINES as string[]).includes(value);
}

export async function generateMetadata({ params }: Props): Promise<Metadata> {
  const { locale, cuisine } = await params;
  if (!isCuisine(cuisine)) return {};
  const tName = await getTranslations({ locale, namespace: 'Cuisines.names' });
  const tIntro = await getTranslations({ locale, namespace: 'Cuisines.intros' });
  return buildMetadata({
    locale: locale as Locale,
    path: `/cuisines/${cuisine}`,
    title: tName(cuisine),
    description: tIntro(cuisine),
  });
}

export default async function CuisinePage({ params }: Props) {
  const { locale, cuisine } = await params;
  setRequestLocale(locale);
  if (!isCuisine(cuisine)) {
    notFound();
  }

  const t = await getTranslations({ locale, namespace: 'Cuisines' });
  const tName = await getTranslations({ locale, namespace: 'Cuisines.names' });
  const tIntro = await getTranslations({ locale, namespace: 'Cuisines.intros' });
  const recipes = getRecipesByCuisine(cuisine);

  return (
    <>
      <PageHeader
        eyebrow={t('title')}
        title={tName(cuisine)}
        subtitle={tIntro(cuisine)}
      />
      <RecipeGrid recipes={recipes} locale={locale as Locale} />
    </>
  );
}
