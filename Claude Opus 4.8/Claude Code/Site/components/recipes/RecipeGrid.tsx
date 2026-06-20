import { AnimatedSection } from '@/components/animations/AnimatedSection';
import { RecipeCard } from './RecipeCard';
import type { Locale, Recipe } from '@/lib/types';

type Props = {
  recipes: Recipe[];
  locale: Locale;
  /** Mark the first N images as priority for LCP. */
  priorityCount?: number;
};

export function RecipeGrid({ recipes, locale, priorityCount = 3 }: Props) {
  return (
    <ul className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
      {recipes.map((recipe, index) => (
        <AnimatedSection
          as="li"
          key={recipe.id}
          delay={Math.min(index, 6) * 0.05}
          className="h-full list-none"
        >
          <RecipeCard
            recipe={recipe}
            locale={locale}
            priority={index < priorityCount}
          />
        </AnimatedSection>
      ))}
    </ul>
  );
}
