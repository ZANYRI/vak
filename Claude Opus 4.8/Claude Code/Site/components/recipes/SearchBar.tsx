'use client';

import { useTranslations } from 'next-intl';
import { useState } from 'react';
import { useRouter } from '@/i18n/navigation';
import { cn } from '@/lib/cn';

type Props = {
  /** Initial value (e.g. when reflecting a URL query). */
  defaultValue?: string;
  placeholder?: string;
  className?: string;
};

/** Hero search box that navigates to the recipes page with a search query. */
export function SearchBar({ defaultValue = '', placeholder, className }: Props) {
  const t = useTranslations('Home');
  const router = useRouter();
  const [value, setValue] = useState(defaultValue);

  const onSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    const q = value.trim();
    router.push({ pathname: '/recipes', query: q ? { search: q } : {} });
  };

  return (
    <form
      role="search"
      onSubmit={onSubmit}
      className={cn('flex w-full items-center gap-2', className)}
    >
      <div className="relative flex-1">
        <svg
          viewBox="0 0 24 24"
          className="pointer-events-none absolute left-4 top-1/2 h-5 w-5 -translate-y-1/2 text-muted"
          fill="none"
          stroke="currentColor"
          strokeWidth="2"
          aria-hidden="true"
        >
          <circle cx="11" cy="11" r="7" />
          <path d="M21 21l-4.3-4.3" strokeLinecap="round" />
        </svg>
        <input
          type="search"
          value={value}
          onChange={(e) => setValue(e.target.value)}
          placeholder={placeholder ?? t('searchPlaceholder')}
          aria-label={t('searchPlaceholder')}
          className="h-13 w-full rounded-full border border-beige-dark bg-card pl-12 pr-4 text-sm text-charcoal shadow-soft placeholder:text-muted/70 focus:border-tomato focus:outline-none"
        />
      </div>
      <button
        type="submit"
        className="h-13 shrink-0 rounded-full bg-tomato px-6 text-sm font-medium text-white shadow-soft transition-all hover:bg-tomato-dark hover:-translate-y-0.5 motion-reduce:hover:translate-y-0"
      >
        {t('searchSubmit')}
      </button>
    </form>
  );
}
