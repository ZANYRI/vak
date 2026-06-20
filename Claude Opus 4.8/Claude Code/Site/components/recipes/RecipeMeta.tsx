import { useTranslations } from 'next-intl';
import type { Recipe } from '@/lib/types';

type Props = {
  recipe: Recipe;
};

function ClockIcon() {
  return (
    <svg viewBox="0 0 24 24" className="h-5 w-5" fill="none" stroke="currentColor" strokeWidth="1.8" aria-hidden="true">
      <circle cx="12" cy="12" r="9" />
      <path d="M12 7v5l3 2" strokeLinecap="round" strokeLinejoin="round" />
    </svg>
  );
}

function GaugeIcon() {
  return (
    <svg viewBox="0 0 24 24" className="h-5 w-5" fill="none" stroke="currentColor" strokeWidth="1.8" aria-hidden="true">
      <path d="M5 18a8 8 0 1114 0" strokeLinecap="round" />
      <path d="M12 14l4-4" strokeLinecap="round" />
    </svg>
  );
}

function ServingsIcon() {
  return (
    <svg viewBox="0 0 24 24" className="h-5 w-5" fill="none" stroke="currentColor" strokeWidth="1.8" aria-hidden="true">
      <path d="M4 11h16M5 11a7 7 0 0114 0M12 4v1M3 15h18l-1 2a3 3 0 01-3 2H7a3 3 0 01-3-2l-1-2z" strokeLinecap="round" strokeLinejoin="round" />
    </svg>
  );
}

export function RecipeMeta({ recipe }: Props) {
  const t = useTranslations('Recipe');
  const tc = useTranslations('Common');
  const td = useTranslations('Difficulty');

  const items = [
    { icon: <ClockIcon />, label: t('prepTime'), value: `${recipe.prepTimeMinutes} ${tc('minutes')}` },
    { icon: <ClockIcon />, label: t('cookTime'), value: `${recipe.cookTimeMinutes} ${tc('minutes')}` },
    { icon: <ClockIcon />, label: t('totalTime'), value: `${recipe.totalTimeMinutes} ${tc('minutes')}` },
    { icon: <GaugeIcon />, label: t('difficulty'), value: td(recipe.difficulty) },
    { icon: <ServingsIcon />, label: t('servings'), value: String(recipe.servings) },
  ];

  return (
    <dl className="grid grid-cols-2 gap-3 sm:grid-cols-3 lg:grid-cols-5">
      {items.map((item) => (
        <div
          key={item.label}
          className="flex flex-col items-center gap-1 rounded-2xl border border-beige bg-card px-3 py-4 text-center shadow-soft"
        >
          <span className="text-tomato">{item.icon}</span>
          <dt className="text-xs text-muted">{item.label}</dt>
          <dd className="text-sm font-semibold text-charcoal">{item.value}</dd>
        </div>
      ))}
    </dl>
  );
}
