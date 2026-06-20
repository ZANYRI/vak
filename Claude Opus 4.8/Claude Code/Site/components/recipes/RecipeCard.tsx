import { useTranslations } from 'next-intl';
import { Link } from '@/i18n/navigation';
import { Badge } from '@/components/ui/Badge';
import { FavoriteButton } from '@/components/ui/FavoriteButton';
import { ImageWithFallback } from '@/components/ui/ImageWithFallback';
import type { Locale, Recipe } from '@/lib/types';

type Props = {
  recipe: Recipe;
  locale: Locale;
  priority?: boolean;
};

export function RecipeCard({ recipe, locale, priority = false }: Props) {
  const t = useTranslations('Recipes');
  const tc = useTranslations('Common');
  const td = useTranslations('Difficulty');
  const tcat = useTranslations('Cuisines.names');

  return (
    <article className="group relative h-full">
      <FavoriteButton slug={recipe.slug} className="absolute right-3 top-3 z-10" />
      <Link
        href={`/recipes/${recipe.slug}`}
        className="flex h-full flex-col overflow-hidden rounded-card border border-beige bg-card shadow-soft transition-all duration-300 hover:-translate-y-1 hover:shadow-lift motion-reduce:hover:translate-y-0"
      >
        <div className="relative aspect-[4/3] overflow-hidden">
          <ImageWithFallback
            src={recipe.image}
            alt={recipe.imageAlt[locale]}
            fill
            sizes="(max-width: 640px) 100vw, (max-width: 1024px) 50vw, 33vw"
            priority={priority}
            wrapperClassName="h-full w-full"
            className="transition-transform duration-500 group-hover:scale-105 motion-reduce:group-hover:scale-100"
          />
          <div className="absolute left-3 top-3 flex gap-2">
            <Badge tone="tomato">{tcat(recipe.cuisine)}</Badge>
          </div>
        </div>

        <div className="flex flex-1 flex-col gap-3 p-5">
          <div className="flex items-center justify-between text-xs text-muted">
            <span className="inline-flex items-center gap-1">
              <svg viewBox="0 0 24 24" className="h-3.5 w-3.5" fill="none" stroke="currentColor" strokeWidth="2" aria-hidden="true">
                <circle cx="12" cy="12" r="9" />
                <path d="M12 7v5l3 2" strokeLinecap="round" strokeLinejoin="round" />
              </svg>
              {recipe.totalTimeMinutes} {tc('minutes')}
            </span>
            <span className="inline-flex items-center gap-1 text-saffron">
              <svg viewBox="0 0 24 24" className="h-3.5 w-3.5 fill-saffron" aria-hidden="true">
                <path d="M12 2l2.9 6.1 6.6.9-4.8 4.6 1.2 6.6L12 18.6 6.1 20.8l1.2-6.6L2.5 9l6.6-.9z" />
              </svg>
              <span className="font-semibold text-charcoal">{recipe.rating.toFixed(1)}</span>
            </span>
          </div>

          <h3 className="font-display text-lg font-semibold leading-snug text-charcoal text-balance">
            {recipe.title[locale]}
          </h3>
          <p className="line-clamp-2 text-sm text-muted">{recipe.description[locale]}</p>

          <div className="mt-auto flex items-center justify-between pt-2">
            <Badge tone="olive">{td(recipe.difficulty)}</Badge>
            <span className="text-sm font-medium text-tomato transition-transform group-hover:translate-x-0.5 motion-reduce:group-hover:translate-x-0">
              {t('viewRecipe')} →
            </span>
          </div>
        </div>
      </Link>
    </article>
  );
}
