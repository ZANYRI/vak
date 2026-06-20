'use client';

import { useTranslations } from 'next-intl';
import { Link } from '@/i18n/navigation';
import { buttonClasses } from '@/components/ui/Button';

export default function LocaleNotFound() {
  const t = useTranslations('NotFound');

  return (
    <div className="flex min-h-[50vh] flex-col items-center justify-center text-center">
      <p className="font-display text-7xl font-semibold text-tomato">404</p>
      <h1 className="mt-4 font-display text-3xl font-semibold text-charcoal">
        {t('title')}
      </h1>
      <p className="mt-3 max-w-md text-muted">{t('body')}</p>
      <Link href="/" className={`mt-7 ${buttonClasses('primary')}`}>
        {t('home')}
      </Link>
    </div>
  );
}
