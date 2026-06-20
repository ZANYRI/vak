import type { Metadata } from 'next';
import { getTranslations, setRequestLocale } from 'next-intl/server';
import { PageHeader } from '@/components/layout/PageHeader';
import { CuisineCard } from '@/components/recipes/CuisineCard';
import { AnimatedSection } from '@/components/animations/AnimatedSection';
import { CUISINES } from '@/lib/constants';
import { countByCuisine, getRecipesByCuisine } from '@/lib/recipes';
import { buildMetadata } from '@/lib/seo';
import type { Locale } from '@/lib/types';

type Props = {
  params: Promise<{ locale: string }>;
};

export async function generateMetadata({ params }: Props): Promise<Metadata> {
  const { locale } = await params;
  const t = await getTranslations({ locale, namespace: 'Cuisines' });
  return buildMetadata({
    locale: locale as Locale,
    path: '/cuisines',
    title: t('title'),
    description: t('metaDescription'),
  });
}

export default async function CuisinesPage({ params }: Props) {
  const { locale } = await params;
  setRequestLocale(locale);
  const loc = locale as Locale;
  const t = await getTranslations({ locale, namespace: 'Cuisines' });

  return (
    <>
      <PageHeader title={t('title')} subtitle={t('subtitle')} />
      <div className="grid grid-cols-2 gap-4 sm:grid-cols-3 lg:grid-cols-5">
        {CUISINES.map((cuisine, i) => (
          <AnimatedSection key={cuisine} delay={Math.min(i, 9) * 0.04}>
            <CuisineCard
              cuisine={cuisine}
              locale={loc}
              count={countByCuisine(cuisine)}
              image={getRecipesByCuisine(cuisine)[0]?.image ?? ''}
            />
          </AnimatedSection>
        ))}
      </div>
    </>
  );
}
