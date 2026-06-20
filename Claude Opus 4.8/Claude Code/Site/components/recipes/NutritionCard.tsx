import { useTranslations } from 'next-intl';
import type { Recipe } from '@/lib/types';

export function NutritionCard({ recipe }: { recipe: Recipe }) {
  const t = useTranslations('Recipe');
  const { calories, protein, fat, carbs } = recipe.nutrition;

  const rows = [
    { label: t('calories'), value: t('kcal', { count: calories }) },
    { label: t('protein'), value: t('grams', { count: protein }) },
    { label: t('fat'), value: t('grams', { count: fat }) },
    { label: t('carbs'), value: t('grams', { count: carbs }) },
  ];

  return (
    <section aria-labelledby="nutrition-heading" className="rounded-card border border-beige bg-card p-6 shadow-soft">
      <h2 id="nutrition-heading" className="font-display text-xl font-semibold text-charcoal">
        {t('nutrition')}
      </h2>
      <dl className="mt-4 divide-y divide-beige">
        {rows.map((row) => (
          <div key={row.label} className="flex items-center justify-between py-2.5">
            <dt className="text-sm text-muted">{row.label}</dt>
            <dd className="text-sm font-semibold text-charcoal">{row.value}</dd>
          </div>
        ))}
      </dl>
    </section>
  );
}
