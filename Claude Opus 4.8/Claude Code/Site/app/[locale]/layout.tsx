import type { Metadata } from 'next';
import { NextIntlClientProvider, hasLocale } from 'next-intl';
import { getTranslations, setRequestLocale } from 'next-intl/server';
import { notFound } from 'next/navigation';
import { Inter, Fraunces } from 'next/font/google';
import { routing } from '@/i18n/routing';
import { FavoritesProvider } from '@/components/providers/FavoritesProvider';
import { Header } from '@/components/layout/Header';
import { Footer } from '@/components/layout/Footer';
import { PageTransition } from '@/components/animations/PageTransition';
import { buildMetadata } from '@/lib/seo';
import type { Locale } from '@/lib/types';
import '../globals.css';

const inter = Inter({
  subsets: ['latin', 'cyrillic'],
  variable: '--font-inter',
  display: 'swap',
});

const fraunces = Fraunces({
  subsets: ['latin'],
  variable: '--font-fraunces',
  display: 'swap',
  axes: ['opsz'],
});

export function generateStaticParams() {
  return routing.locales.map((locale) => ({ locale }));
}

type LayoutProps = {
  children: React.ReactNode;
  params: Promise<{ locale: string }>;
};

export async function generateMetadata({
  params,
}: {
  params: Promise<{ locale: string }>;
}): Promise<Metadata> {
  const { locale } = await params;
  const t = await getTranslations({ locale, namespace: 'Common' });
  return {
    metadataBase: new URL(
      process.env.NEXT_PUBLIC_SITE_URL ?? 'http://localhost:3000',
    ),
    ...buildMetadata({
      locale: locale as Locale,
      title: `${t('siteName')} — ${t('tagline')}`,
      description: t('tagline'),
    }),
    title: {
      default: `${t('siteName')} — ${t('tagline')}`,
      template: `%s · ${t('siteName')}`,
    },
    icons: {
      icon: [{ url: '/favicon.svg', type: 'image/svg+xml' }],
    },
  };
}

export default async function LocaleLayout({ children, params }: LayoutProps) {
  const { locale } = await params;
  if (!hasLocale(routing.locales, locale)) {
    notFound();
  }
  setRequestLocale(locale);

  const t = await getTranslations({ locale, namespace: 'Common' });

  return (
    <html lang={locale} className={`${inter.variable} ${fraunces.variable}`}>
      <body className="min-h-screen bg-cream antialiased">
        <NextIntlClientProvider>
          <FavoritesProvider>
            <a
              href="#main-content"
              className="sr-only focus:not-sr-only focus:absolute focus:left-4 focus:top-4 focus:z-[100] focus:rounded-full focus:bg-tomato focus:px-5 focus:py-2 focus:text-sm focus:font-medium focus:text-white"
            >
              {t('skipToContent')}
            </a>
            <Header />
            <main id="main-content" className="container-page py-10 sm:py-14">
              <PageTransition>{children}</PageTransition>
            </main>
            <Footer />
          </FavoritesProvider>
        </NextIntlClientProvider>
      </body>
    </html>
  );
}
