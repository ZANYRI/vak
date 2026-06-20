import type { Metadata } from 'next';
import { getTranslations, setRequestLocale } from 'next-intl/server';
import { PageHeader } from '@/components/layout/PageHeader';
import { CategoryCard } from '@/components/recipes/CategoryCard';
import { AnimatedSection } from '@/components/animations/AnimatedSection';
import { CATEGORIES } from '@/lib/constants';
import { countByCategory } from '@/lib/recipes';
import { buildMetadata } from '@/lib/seo';
import type { Locale } from '@/lib/types';

type Props = {
  params: Promise<{ locale: string }>;
};

export async function generateMetadata({ params }: Props): Promise<Metadata> {
  const { locale } = await params;
  const t = await getTranslations({ locale, namespace: 'Categories' });
  return buildMetadata({
    locale: locale as Locale,
    path: '/categories',
    title: t('title'),
    description: t('metaDescription'),
  });
}

export default async function CategoriesPage({ params }: Props) {
  const { locale } = await params;
  setRequestLocale(locale);
  const t = await getTranslations({ locale, namespace: 'Categories' });

  return (
    <>
      <PageHeader title={t('title')} subtitle={t('subtitle')} />
      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
        {CATEGORIES.map((category, i) => (
          <AnimatedSection key={category} delay={Math.min(i, 9) * 0.04}>
            <CategoryCard category={category} count={countByCategory(category)} />
          </AnimatedSection>
        ))}
      </div>
    </>
  );
}
