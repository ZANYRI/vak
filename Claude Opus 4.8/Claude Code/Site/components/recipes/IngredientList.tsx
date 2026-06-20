import { useTranslations } from 'next-intl';
import type { Locale, Recipe } from '@/lib/types';

type Props = {
  recipe: Recipe;
  locale: Locale;
};

export function IngredientList({ recipe, locale }: Props) {
  const t = useTranslations('Recipe');
  const items = recipe.ingredients[locale];

  return (
    <section aria-labelledby="ingredients-heading" className="rounded-card border border-beige bg-card p-6 shadow-soft">
      <h2 id="ingredients-heading" className="font-display text-2xl font-semibold text-charcoal">
        {t('ingredients')}
      </h2>
      <ul className="mt-4 space-y-3">
        {items.map((ingredient, index) => (
          <li key={index} className="flex items-start gap-3 text-charcoal">
            <span aria-hidden="true" className="mt-1.5 h-2 w-2 shrink-0 rounded-full bg-tomato" />
            <span className="text-sm leading-relaxed">{ingredient}</span>
          </li>
        ))}
      </ul>
    </section>
  );
}
