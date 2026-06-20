'use client';

import { useTranslations } from 'next-intl';
import { useFavorites } from '@/components/providers/FavoritesProvider';
import { RecipeGrid } from './RecipeGrid';
import { EmptyState } from '@/components/ui/EmptyState';
import { RecipeCardSkeleton } from '@/components/ui/Skeleton';
import { Link } from '@/i18n/navigation';
import { buttonClasses } from '@/components/ui/Button';
import type { Locale, Recipe } from '@/lib/types';

type Props = {
  recipes: Recipe[];
  locale: Locale;
};

export function FavoritesList({ recipes, locale }: Props) {
  const t = useTranslations('Favorites');
  const { favorites, ready } = useFavorites();

  if (!ready) {
    return (
      <div className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
        {Array.from({ length: 3 }).map((_, i) => (
          <RecipeCardSkeleton key={i} />
        ))}
      </div>
    );
  }

  const saved = recipes.filter((recipe) => favorites.includes(recipe.slug));

  if (saved.length === 0) {
    return (
      <EmptyState
        title={t('emptyTitle')}
        body={t('emptyBody')}
        icon={
          <svg viewBox="0 0 24 24" className="h-7 w-7 fill-none stroke-current" strokeWidth="1.8" aria-hidden="true">
            <path
              d="M12 20.5l-1.45-1.32C5.4 14.5 2 11.4 2 7.6 2 4.9 4.1 2.8 6.8 2.8c1.5 0 3 .7 3.9 1.9.9-1.2 2.4-1.9 3.9-1.9 2.7 0 4.8 2.1 4.8 4.8 0 3.8-3.4 6.9-8.55 11.58L12 20.5z"
              strokeLinecap="round"
              strokeLinejoin="round"
            />
          </svg>
        }
        action={
          <Link href="/recipes" className={buttonClasses('primary')}>
            {t('browse')}
          </Link>
        }
      />
    );
  }

  return (
    <div className="space-y-6">
      <p className="text-sm font-medium text-muted" aria-live="polite">
        {t('count', { count: saved.length })}
      </p>
      <RecipeGrid recipes={saved} locale={locale} />
    </div>
  );
}
