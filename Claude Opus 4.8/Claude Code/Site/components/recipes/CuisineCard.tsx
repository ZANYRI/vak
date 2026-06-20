import { useTranslations } from 'next-intl';
import { Link } from '@/i18n/navigation';
import { ImageWithFallback } from '@/components/ui/ImageWithFallback';
import type { CuisineSlug, Locale } from '@/lib/types';

type Props = {
  cuisine: CuisineSlug;
  locale: Locale;
  image: string;
  count: number;
};

export function CuisineCard({ cuisine, image, count }: Props) {
  const tName = useTranslations('Cuisines.names');
  const t = useTranslations('Cuisines');

  return (
    <Link
      href={`/cuisines/${cuisine}`}
      className="group relative block aspect-[4/5] overflow-hidden rounded-card shadow-soft transition-all duration-300 hover:-translate-y-1 hover:shadow-lift motion-reduce:hover:translate-y-0"
    >
      <ImageWithFallback
        src={image}
        alt={tName(cuisine)}
        fill
        sizes="(max-width: 640px) 50vw, (max-width: 1024px) 33vw, 16vw"
        wrapperClassName="absolute inset-0 h-full w-full"
        className="transition-transform duration-500 group-hover:scale-110 motion-reduce:group-hover:scale-100"
      />
      <div className="absolute inset-0 bg-gradient-to-t from-charcoal/80 via-charcoal/20 to-transparent" />
      <div className="absolute inset-x-0 bottom-0 p-4">
        <h3 className="font-display text-lg font-semibold text-white">{tName(cuisine)}</h3>
        <p className="text-xs text-white/80">{t('recipeCount', { count })}</p>
      </div>
    </Link>
  );
}
