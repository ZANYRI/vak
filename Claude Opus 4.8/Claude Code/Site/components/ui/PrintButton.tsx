'use client';

import { useTranslations } from 'next-intl';
import { cn } from '@/lib/cn';

export function PrintButton({ className }: { className?: string }) {
  const t = useTranslations('Recipe');

  return (
    <button
      type="button"
      onClick={() => window.print()}
      className={cn(
        'inline-flex h-11 items-center justify-center gap-2 rounded-full border border-beige-dark bg-card px-6 text-sm font-medium text-charcoal transition-all hover:border-saffron hover:text-[#8a5d00]',
        className,
      )}
    >
      <svg viewBox="0 0 24 24" className="h-4 w-4" fill="none" stroke="currentColor" strokeWidth="2" aria-hidden="true">
        <path d="M6 9V3h12v6" strokeLinecap="round" strokeLinejoin="round" />
        <rect x="6" y="13" width="12" height="8" rx="1" />
        <path d="M6 17H4a2 2 0 01-2-2v-3a3 3 0 013-3h14a3 3 0 013 3v3a2 2 0 01-2 2h-2" />
      </svg>
      {t('print')}
    </button>
  );
}
