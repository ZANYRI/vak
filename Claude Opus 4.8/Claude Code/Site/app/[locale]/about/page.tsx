import type { Metadata } from 'next';
import { getTranslations, setRequestLocale } from 'next-intl/server';
import { PageHeader } from '@/components/layout/PageHeader';
import { AnimatedSection } from '@/components/animations/AnimatedSection';
import { buildMetadata } from '@/lib/seo';
import type { Locale } from '@/lib/types';

type Props = {
  params: Promise<{ locale: string }>;
};

export async function generateMetadata({ params }: Props): Promise<Metadata> {
  const { locale } = await params;
  const t = await getTranslations({ locale, namespace: 'About' });
  return buildMetadata({
    locale: locale as Locale,
    path: '/about',
    title: t('title'),
    description: t('metaDescription'),
  });
}

export default async function AboutPage({ params }: Props) {
  const { locale } = await params;
  setRequestLocale(locale);
  const t = await getTranslations({ locale, namespace: 'About' });

  const sections = [
    { icon: '📖', title: t('introTitle'), body: t('introBody') },
    { icon: '🌿', title: t('philosophyTitle'), body: t('philosophyBody') },
    { icon: '🗂️', title: t('dataTitle'), body: t('dataBody') },
    { icon: '📷', title: t('imageTitle'), body: t('imageBody') },
    { icon: '♿', title: t('a11yTitle'), body: t('a11yBody') },
  ];

  return (
    <>
      <PageHeader title={t('title')} />
      <div className="grid gap-6 sm:grid-cols-2">
        {sections.map((section, i) => (
          <AnimatedSection
            as="article"
            key={section.title}
            delay={Math.min(i, 6) * 0.05}
            className="rounded-card border border-beige bg-card p-7 shadow-soft"
          >
            <span aria-hidden="true" className="text-3xl">{section.icon}</span>
            <h2 className="mt-3 font-display text-xl font-semibold text-charcoal">
              {section.title}
            </h2>
            <p className="mt-2 text-muted leading-relaxed">{section.body}</p>
          </AnimatedSection>
        ))}
      </div>
    </>
  );
}
