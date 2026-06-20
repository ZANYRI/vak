'use client';

import { useTranslations } from 'next-intl';
import { useState } from 'react';
import { cn } from '@/lib/cn';

type Props = {
  title: string;
  text?: string;
  className?: string;
};

export function ShareButton({ title, text, className }: Props) {
  const t = useTranslations('Recipe');
  const [copied, setCopied] = useState(false);

  const onShare = async () => {
    const url = typeof window !== 'undefined' ? window.location.href : '';
    try {
      if (navigator.share) {
        await navigator.share({ title, text, url });
        return;
      }
      await navigator.clipboard.writeText(url);
      setCopied(true);
      window.setTimeout(() => setCopied(false), 2000);
    } catch {
      /* user cancelled or APIs unavailable — no-op */
    }
  };

  return (
    <button
      type="button"
      onClick={onShare}
      className={cn(
        'inline-flex h-11 items-center justify-center gap-2 rounded-full border border-beige-dark bg-card px-6 text-sm font-medium text-charcoal transition-all hover:border-olive hover:text-olive',
        className,
      )}
    >
      <svg viewBox="0 0 24 24" className="h-4 w-4" fill="none" stroke="currentColor" strokeWidth="2" aria-hidden="true">
        <circle cx="18" cy="5" r="3" />
        <circle cx="6" cy="12" r="3" />
        <circle cx="18" cy="19" r="3" />
        <path d="M8.6 13.5l6.8 4M15.4 6.5l-6.8 4" strokeLinecap="round" />
      </svg>
      {copied ? t('shareCopied') : t('share')}
    </button>
  );
}
