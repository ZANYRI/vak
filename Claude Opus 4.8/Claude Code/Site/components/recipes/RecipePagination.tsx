import { useTranslations } from 'next-intl';
import { Link } from '@/i18n/navigation';
import { getPageRange } from '@/lib/pagination';
import { cn } from '@/lib/cn';

type Props = {
  page: number;
  totalPages: number;
  /** Filter query params to preserve across page links. */
  query?: Record<string, string>;
};

function hrefForPage(n: number, query: Record<string, string>) {
  const pathname = n <= 1 ? '/recipes' : `/recipes/page/${n}`;
  return { pathname, query };
}

export function RecipePagination({ page, totalPages, query = {} }: Props) {
  const t = useTranslations('Pagination');
  if (totalPages <= 1) return null;

  const tokens = getPageRange(page, totalPages);
  const hasPrev = page > 1;
  const hasNext = page < totalPages;

  const arrowBase =
    'inline-flex h-10 items-center justify-center gap-1 rounded-full border px-4 text-sm font-medium transition-colors';

  return (
    <nav aria-label={t('label')} className="mt-10 flex flex-col items-center gap-4">
      <ul className="flex flex-wrap items-center justify-center gap-2">
        <li>
          {hasPrev ? (
            <Link
              href={hrefForPage(page - 1, query)}
              rel="prev"
              className={cn(arrowBase, 'border-beige-dark bg-card text-charcoal hover:border-tomato hover:text-tomato')}
            >
              <span aria-hidden="true">←</span>
              {t('previous')}
            </Link>
          ) : (
            <span
              aria-disabled="true"
              className={cn(arrowBase, 'cursor-not-allowed border-beige bg-beige/40 text-muted/60')}
            >
              <span aria-hidden="true">←</span>
              {t('previous')}
            </span>
          )}
        </li>

        {tokens.map((token, index) =>
          token === 'ellipsis' ? (
            <li key={`ellipsis-${index}`} aria-hidden="true" className="px-2 text-muted">
              …
            </li>
          ) : (
            <li key={token}>
              {token === page ? (
                <span
                  aria-current="page"
                  aria-label={t('currentPage', { number: token })}
                  className="inline-flex h-10 w-10 items-center justify-center rounded-full bg-tomato text-sm font-semibold text-white shadow-soft"
                >
                  {token}
                </span>
              ) : (
                <Link
                  href={hrefForPage(token, query)}
                  aria-label={t('goToPage', { number: token })}
                  className="inline-flex h-10 w-10 items-center justify-center rounded-full border border-beige-dark bg-card text-sm font-medium text-charcoal transition-colors hover:border-tomato hover:text-tomato"
                >
                  {token}
                </Link>
              )}
            </li>
          ),
        )}

        <li>
          {hasNext ? (
            <Link
              href={hrefForPage(page + 1, query)}
              rel="next"
              className={cn(arrowBase, 'border-beige-dark bg-card text-charcoal hover:border-tomato hover:text-tomato')}
            >
              {t('next')}
              <span aria-hidden="true">→</span>
            </Link>
          ) : (
            <span
              aria-disabled="true"
              className={cn(arrowBase, 'cursor-not-allowed border-beige bg-beige/40 text-muted/60')}
            >
              {t('next')}
              <span aria-hidden="true">→</span>
            </span>
          )}
        </li>
      </ul>
      <p className="text-sm text-muted">{t('summary', { current: page, total: totalPages })}</p>
    </nav>
  );
}
