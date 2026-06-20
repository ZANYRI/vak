'use client';

import { useTranslations } from 'next-intl';
import { motion, useReducedMotion } from 'motion/react';
import { useFavorites } from '@/components/providers/FavoritesProvider';
import { cn } from '@/lib/cn';

type Props = {
  slug: string;
  variant?: 'icon' | 'full';
  className?: string;
};

export function FavoriteButton({ slug, variant = 'icon', className }: Props) {
  const t = useTranslations('Recipe');
  const { isFavorite, toggleFavorite, ready } = useFavorites();
  const reduce = useReducedMotion();
  const active = isFavorite(slug);
  const label = active ? t('removeFavorite') : t('addFavorite');

  const heart = (
    <motion.svg
      viewBox="0 0 24 24"
      className={cn('h-5 w-5', active ? 'fill-tomato stroke-tomato' : 'fill-none stroke-current')}
      strokeWidth="2"
      aria-hidden="true"
      animate={reduce ? undefined : { scale: active ? [1, 1.25, 1] : 1 }}
      transition={{ duration: 0.3 }}
    >
      <path
        d="M12 20.5l-1.45-1.32C5.4 14.5 2 11.4 2 7.6 2 4.9 4.1 2.8 6.8 2.8c1.5 0 3 .7 3.9 1.9.9-1.2 2.4-1.9 3.9-1.9 2.7 0 4.8 2.1 4.8 4.8 0 3.8-3.4 6.9-8.55 11.58L12 20.5z"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
    </motion.svg>
  );

  if (variant === 'full') {
    return (
      <button
        type="button"
        onClick={() => toggleFavorite(slug)}
        aria-pressed={active}
        disabled={!ready}
        className={cn(
          'inline-flex h-11 items-center justify-center gap-2 rounded-full border px-6 text-sm font-medium transition-all',
          active
            ? 'border-tomato bg-tomato/10 text-tomato'
            : 'border-beige-dark bg-card text-charcoal hover:border-tomato hover:text-tomato',
          'disabled:opacity-50',
          className,
        )}
      >
        {heart}
        {label}
      </button>
    );
  }

  return (
    <button
      type="button"
      onClick={(e) => {
        e.preventDefault();
        e.stopPropagation();
        toggleFavorite(slug);
      }}
      aria-pressed={active}
      aria-label={label}
      title={label}
      disabled={!ready}
      className={cn(
        'inline-flex h-10 w-10 items-center justify-center rounded-full bg-card/90 text-charcoal shadow-soft backdrop-blur transition-all hover:scale-110 motion-reduce:hover:scale-100 disabled:opacity-50',
        className,
      )}
    >
      {heart}
    </button>
  );
}
