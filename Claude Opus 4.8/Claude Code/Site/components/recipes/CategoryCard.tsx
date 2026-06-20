import { useTranslations } from 'next-intl';
import { Link } from '@/i18n/navigation';
import type { CategorySlug } from '@/lib/types';

const EMOJI: Record<CategorySlug, string> = {
  breakfast: '🍳',
  lunch: '🥪',
  dinner: '🍽️',
  desserts: '🍰',
  soups: '🥣',
  salads: '🥗',
  vegetarian: '🥦',
  'quick-meals': '⚡',
  baking: '🥐',
  drinks: '🥤',
};

type Props = {
  category: CategorySlug;
  count: number;
};

export function CategoryCard({ category, count }: Props) {
  const tName = useTranslations('Categories.names');
  const t = useTranslations('Categories');

  return (
    <Link
      href={`/categories/${category}`}
      className="group flex items-center gap-4 rounded-card border border-beige bg-card p-5 shadow-soft transition-all duration-300 hover:-translate-y-1 hover:border-tomato/40 hover:shadow-lift motion-reduce:hover:translate-y-0"
    >
      <span
        aria-hidden="true"
        className="flex h-12 w-12 shrink-0 items-center justify-center rounded-2xl bg-beige text-2xl transition-transform group-hover:scale-110 motion-reduce:group-hover:scale-100"
      >
        {EMOJI[category]}
      </span>
      <span className="min-w-0">
        <span className="block font-semibold text-charcoal">{tName(category)}</span>
        <span className="block text-xs text-muted">{t('recipeCount', { count })}</span>
      </span>
    </Link>
  );
}
