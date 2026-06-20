import { useTranslations } from 'next-intl';
import { RecipeCard } from './RecipeCard';
import type { Locale, Recipe } from '@/lib/types';

type Props = {
  recipes: Recipe[];
  locale: Locale;
};

export function RelatedRecipes({ recipes, locale }: Props) {
  const t = useTranslations('Recipe');
  if (recipes.length === 0) return null;

  return (
    <section aria-labelledby="related-heading" className="mt-16">
      <h2 id="related-heading" className="font-display text-2xl font-semibold text-charcoal">
        {t('related')}
      </h2>
      <ul className="mt-6 grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
        {recipes.map((recipe) => (
          <li key={recipe.id} className="list-none">
            <RecipeCard recipe={recipe} locale={locale} />
          </li>
        ))}
      </ul>
    </section>
  );
}
