import { useTranslations } from 'next-intl';
import { Badge } from '@/components/ui/Badge';
import { ImageWithFallback } from '@/components/ui/ImageWithFallback';
import { FavoriteButton } from '@/components/ui/FavoriteButton';
import { ShareButton } from '@/components/ui/ShareButton';
import { PrintButton } from '@/components/ui/PrintButton';
import type { Locale, Recipe } from '@/lib/types';

type Props = {
  recipe: Recipe;
  locale: Locale;
};

export function RecipeHero({ recipe, locale }: Props) {
  const t = useTranslations('Recipe');
  const tc = useTranslations('Common');
  const tCuisine = useTranslations('Cuisines.names');
  const tCategory = useTranslations('Categories.names');
  const tDiet = useTranslations('Diet');

  return (
    <header className="grid gap-8 lg:grid-cols-2 lg:items-center">
      <div className="order-2 lg:order-1">
        <div className="flex flex-wrap items-center gap-2">
          <Badge tone="tomato">{tCuisine(recipe.cuisine)}</Badge>
          <Badge tone="olive">{tCategory(recipe.category)}</Badge>
          <span className="inline-flex items-center gap-1 text-sm text-saffron">
            <svg viewBox="0 0 24 24" className="h-4 w-4 fill-saffron" aria-hidden="true">
              <path d="M12 2l2.9 6.1 6.6.9-4.8 4.6 1.2 6.6L12 18.6 6.1 20.8l1.2-6.6L2.5 9l6.6-.9z" />
            </svg>
            <span className="font-semibold text-charcoal">{recipe.rating.toFixed(1)}</span>
          </span>
        </div>

        <h1 className="mt-4 font-display text-4xl font-semibold leading-tight text-charcoal text-balance sm:text-5xl">
          {recipe.title[locale]}
        </h1>
        <p className="mt-4 max-w-prose text-lg text-muted">{recipe.description[locale]}</p>

        {recipe.diet.length > 0 && (
          <ul className="mt-5 flex flex-wrap gap-2" aria-label={t('diet')}>
            {recipe.diet.map((diet) => (
              <li key={diet}>
                <Badge tone="saffron">{tDiet(diet)}</Badge>
              </li>
            ))}
          </ul>
        )}

        <div className="mt-7 flex flex-wrap gap-3" data-no-print>
          <FavoriteButton slug={recipe.slug} variant="full" />
          <ShareButton title={recipe.title[locale]} text={recipe.description[locale]} />
          <PrintButton />
        </div>
        <p className="mt-4 text-xs text-muted">
          {t('servings')}: {tc('servings', { count: recipe.servings })}
        </p>
      </div>

      <div className="order-1 lg:order-2">
        <ImageWithFallback
          src={recipe.image}
          alt={recipe.imageAlt[locale]}
          width={900}
          height={675}
          priority
          sizes="(max-width: 1024px) 100vw, 50vw"
          wrapperClassName="aspect-[4/3] w-full rounded-card shadow-lift"
        />
      </div>
    </header>
  );
}
