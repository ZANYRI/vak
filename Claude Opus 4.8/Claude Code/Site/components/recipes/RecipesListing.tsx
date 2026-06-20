import { getTranslations } from 'next-intl/server';
import { Link } from '@/i18n/navigation';
import { buttonClasses } from '@/components/ui/Button';
import { EmptyState } from '@/components/ui/EmptyState';
import { RecipeFilters } from './RecipeFilters';
import { RecipeGrid } from './RecipeGrid';
import { RecipePagination } from './RecipePagination';
import { filterRecipes, paginateRecipes } from '@/lib/recipes';
import type { Locale, RecipeFilters as Filters } from '@/lib/types';

type Props = {
  locale: Locale;
  filters: Filters;
  page: number;
};

/** Turn typed filters back into a string query for pagination links. */
function filtersToQuery(filters: Filters): Record<string, string> {
  const query: Record<string, string> = {};
  if (filters.search) query.search = filters.search;
  if (filters.cuisine) query.cuisine = filters.cuisine;
  if (filters.category) query.category = filters.category;
  if (filters.difficulty) query.difficulty = filters.difficulty;
  if (filters.diet) query.diet = filters.diet;
  if (filters.mealType) query.mealType = filters.mealType;
  if (filters.maxTime) query.maxTime = String(filters.maxTime);
  if (filters.sort) query.sort = filters.sort;
  return query;
}

export async function RecipesListing({ locale, filters, page }: Props) {
  const t = await getTranslations({ locale, namespace: 'Recipes' });

  const filtered = filterRecipes(locale, filters);
  const { items, page: safePage, totalPages, totalItems } = paginateRecipes(
    filtered,
    page,
  );

  return (
    <div className="space-y-8">
      <RecipeFilters current={filters} />

      <p className="text-sm font-medium text-muted" aria-live="polite">
        {t('resultsCount', { count: totalItems })}
      </p>

      {items.length > 0 ? (
        <>
          <RecipeGrid recipes={items} locale={locale} />
          <RecipePagination
            page={safePage}
            totalPages={totalPages}
            query={filtersToQuery(filters)}
          />
        </>
      ) : (
        <EmptyState
          title={t('noResultsTitle')}
          body={t('noResultsBody')}
          icon={
            <svg viewBox="0 0 24 24" className="h-7 w-7" fill="none" stroke="currentColor" strokeWidth="1.8" aria-hidden="true">
              <circle cx="11" cy="11" r="7" />
              <path d="M21 21l-4.3-4.3" strokeLinecap="round" />
            </svg>
          }
          action={
            <Link href="/recipes" className={buttonClasses('primary')}>
              {t('clearFilters')}
            </Link>
          }
        />
      )}
    </div>
  );
}
