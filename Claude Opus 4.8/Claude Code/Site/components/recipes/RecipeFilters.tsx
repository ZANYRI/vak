'use client';

import { useTranslations } from 'next-intl';
import { useEffect, useRef, useState } from 'react';
import { AnimatePresence, motion } from 'motion/react';
import { useRouter } from '@/i18n/navigation';
import { Select } from '@/components/ui/Select';
import {
  CATEGORIES,
  CUISINES,
  DIETS,
  DIFFICULTIES,
  MEAL_TYPES,
  SORT_KEYS,
  TIME_OPTIONS,
} from '@/lib/constants';
import type { RecipeFilters as Filters } from '@/lib/types';
import { cn } from '@/lib/cn';

type Props = {
  current: Filters;
};

export function RecipeFilters({ current }: Props) {
  const t = useTranslations('Filters');
  const tr = useTranslations('Recipes');
  const tCuisine = useTranslations('Cuisines.names');
  const tCategory = useTranslations('Categories.names');
  const tDifficulty = useTranslations('Difficulty');
  const tDiet = useTranslations('Diet');
  const tMeal = useTranslations('MealType');
  const tSort = useTranslations('Recipes.sort');

  const router = useRouter();
  const [open, setOpen] = useState(false);
  const [search, setSearch] = useState(current.search ?? '');
  const firstRender = useRef(true);

  // Build the next query object from the current filters plus an override.
  const pushFilters = (overrides: Partial<Filters>) => {
    const next: Filters = { ...current, ...overrides };
    const query: Record<string, string> = {};
    if (next.search) query.search = next.search;
    if (next.cuisine) query.cuisine = next.cuisine;
    if (next.category) query.category = next.category;
    if (next.difficulty) query.difficulty = next.difficulty;
    if (next.diet) query.diet = next.diet;
    if (next.mealType) query.mealType = next.mealType;
    if (next.maxTime) query.maxTime = String(next.maxTime);
    if (next.sort) query.sort = next.sort;
    // Always navigate back to page 1 when filters change.
    router.push({ pathname: '/recipes', query });
  };

  // Debounce the free-text search so we don't navigate on every keystroke.
  useEffect(() => {
    if (firstRender.current) {
      firstRender.current = false;
      return;
    }
    const id = window.setTimeout(() => {
      if ((current.search ?? '') !== search.trim()) {
        pushFilters({ search: search.trim() || undefined });
      }
    }, 400);
    return () => window.clearTimeout(id);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [search]);

  const allOption = { value: '', label: t('all') };

  const hasFilters = Boolean(
    current.search ||
      current.cuisine ||
      current.category ||
      current.difficulty ||
      current.diet ||
      current.mealType ||
      current.maxTime,
  );

  const fields = (
    <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
      <Field label={t('cuisine')}>
        <Select
          aria-label={t('cuisine')}
          value={current.cuisine ?? ''}
          onChange={(e) => pushFilters({ cuisine: (e.target.value || undefined) as Filters['cuisine'] })}
          options={[allOption, ...CUISINES.map((c) => ({ value: c, label: tCuisine(c) }))]}
        />
      </Field>
      <Field label={t('category')}>
        <Select
          aria-label={t('category')}
          value={current.category ?? ''}
          onChange={(e) => pushFilters({ category: (e.target.value || undefined) as Filters['category'] })}
          options={[allOption, ...CATEGORIES.map((c) => ({ value: c, label: tCategory(c) }))]}
        />
      </Field>
      <Field label={t('mealType')}>
        <Select
          aria-label={t('mealType')}
          value={current.mealType ?? ''}
          onChange={(e) => pushFilters({ mealType: (e.target.value || undefined) as Filters['mealType'] })}
          options={[allOption, ...MEAL_TYPES.map((m) => ({ value: m, label: tMeal(m) }))]}
        />
      </Field>
      <Field label={t('difficulty')}>
        <Select
          aria-label={t('difficulty')}
          value={current.difficulty ?? ''}
          onChange={(e) => pushFilters({ difficulty: (e.target.value || undefined) as Filters['difficulty'] })}
          options={[allOption, ...DIFFICULTIES.map((d) => ({ value: d, label: tDifficulty(d) }))]}
        />
      </Field>
      <Field label={t('diet')}>
        <Select
          aria-label={t('diet')}
          value={current.diet ?? ''}
          onChange={(e) => pushFilters({ diet: (e.target.value || undefined) as Filters['diet'] })}
          options={[allOption, ...DIETS.map((d) => ({ value: d, label: tDiet(d) }))]}
        />
      </Field>
      <Field label={t('maxTime')}>
        <Select
          aria-label={t('maxTime')}
          value={current.maxTime ? String(current.maxTime) : ''}
          onChange={(e) => pushFilters({ maxTime: e.target.value ? Number(e.target.value) : undefined })}
          options={[
            { value: '', label: t('anyTime') },
            ...TIME_OPTIONS.map((m) => ({ value: String(m), label: t('upToMinutes', { count: m }) })),
          ]}
        />
      </Field>
    </div>
  );

  return (
    <section
      aria-label={t('title')}
      className="rounded-card border border-beige bg-card p-5 shadow-soft"
    >
      <div className="flex flex-col gap-4 lg:flex-row lg:items-center">
        {/* Search */}
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
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            placeholder={tr('searchPlaceholder')}
            aria-label={tr('searchPlaceholder')}
            className="h-11 w-full rounded-full border border-beige-dark bg-cream pl-12 pr-4 text-sm text-charcoal placeholder:text-muted/70 focus:border-tomato focus:outline-none"
          />
        </div>

        {/* Sort */}
        <div className="flex items-center gap-2">
          <label htmlFor="sort" className="sr-only">
            {tr('sortLabel')}
          </label>
          <Select
            id="sort"
            aria-label={tr('sortLabel')}
            className="min-w-44"
            value={current.sort ?? 'newest'}
            onChange={(e) => pushFilters({ sort: e.target.value as Filters['sort'] })}
            options={SORT_KEYS.map((s) => ({ value: s, label: tSort(s) }))}
          />
        </div>

        {/* Toggle (mobile) */}
        <button
          type="button"
          onClick={() => setOpen((v) => !v)}
          aria-expanded={open}
          aria-controls="filter-fields"
          className="inline-flex h-11 items-center justify-center gap-2 rounded-full border border-beige-dark bg-card px-5 text-sm font-medium text-charcoal lg:hidden"
        >
          <svg viewBox="0 0 24 24" className="h-4 w-4" fill="none" stroke="currentColor" strokeWidth="2" aria-hidden="true">
            <path d="M4 6h16M7 12h10M10 18h4" strokeLinecap="round" />
          </svg>
          {open ? t('close') : t('open')}
        </button>
      </div>

      {/* Fields: always visible on desktop, collapsible on mobile */}
      <div id="filter-fields" className="hidden lg:mt-5 lg:block">
        {fields}
      </div>
      <AnimatePresence initial={false}>
        {open && (
          <motion.div
            className="overflow-hidden lg:hidden"
            initial={{ height: 0, opacity: 0 }}
            animate={{ height: 'auto', opacity: 1 }}
            exit={{ height: 0, opacity: 0 }}
            transition={{ duration: 0.25, ease: 'easeOut' }}
          >
            <div className="mt-5">{fields}</div>
          </motion.div>
        )}
      </AnimatePresence>

      {hasFilters && (
        <div className="mt-5 flex justify-end">
          <button
            type="button"
            onClick={() => {
              setSearch('');
              router.push({ pathname: '/recipes', query: {} });
            }}
            className={cn(
              'inline-flex items-center gap-1.5 rounded-full px-4 py-2 text-sm font-medium text-tomato transition-colors hover:bg-tomato/10',
            )}
          >
            <svg viewBox="0 0 24 24" className="h-4 w-4" fill="none" stroke="currentColor" strokeWidth="2" aria-hidden="true">
              <path d="M6 6l12 12M18 6L6 18" strokeLinecap="round" />
            </svg>
            {tr('clearFilters')}
          </button>
        </div>
      )}
    </section>
  );
}

function Field({ label, children }: { label: string; children: React.ReactNode }) {
  return (
    <label className="block space-y-1.5">
      <span className="text-xs font-medium text-muted">{label}</span>
      {children}
    </label>
  );
}
