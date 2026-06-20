'use client';

import { useLocale, useTranslations } from 'next-intl';
import { useSearchParams } from 'next/navigation';
import { useTransition } from 'react';
import { usePathname, useRouter } from '@/i18n/navigation';
import { locales, type Locale } from '@/i18n/routing';
import { cn } from '@/lib/cn';

const LABELS: Record<Locale, string> = { en: 'EN', ru: 'RU' };
const FULL: Record<Locale, 'en' | 'ru'> = { en: 'en', ru: 'ru' };

export function LanguageSwitcher() {
  const activeLocale = useLocale() as Locale;
  const t = useTranslations('LanguageSwitcher');
  const pathname = usePathname();
  const searchParams = useSearchParams();
  const router = useRouter();
  const [isPending, startTransition] = useTransition();

  const switchTo = (next: Locale) => {
    if (next === activeLocale) return;
    const query = Object.fromEntries(searchParams.entries());
    startTransition(() => {
      // next-intl preserves the current pathname and swaps the locale prefix.
      router.replace({ pathname, query }, { locale: next });
    });
  };

  return (
    <div
      role="group"
      aria-label={t('label')}
      className="inline-flex items-center rounded-full border border-beige-dark bg-card p-1"
    >
      {locales.map((locale) => {
        const isActive = locale === activeLocale;
        return (
          <button
            key={locale}
            type="button"
            onClick={() => switchTo(locale)}
            disabled={isPending}
            aria-pressed={isActive}
            aria-label={t('switchTo', { language: t(FULL[locale]) })}
            className={cn(
              'rounded-full px-3 py-1.5 text-xs font-semibold transition-all duration-200',
              'focus-visible:outline-2 disabled:opacity-60',
              isActive
                ? 'bg-tomato text-white shadow-soft'
                : 'text-muted hover:text-charcoal hover:bg-beige/60',
            )}
          >
            {LABELS[locale]}
          </button>
        );
      })}
    </div>
  );
}
