import { useTranslations } from 'next-intl';
import { Link } from '@/i18n/navigation';

const EXPLORE = [
  { href: '/recipes', key: 'recipes' as const },
  { href: '/categories', key: 'categories' as const },
  { href: '/cuisines', key: 'cuisines' as const },
  { href: '/favorites', key: 'favorites' as const },
];

export function Footer() {
  const t = useTranslations('Footer');
  const tn = useTranslations('Nav');
  const tc = useTranslations('Common');
  const year = 2026;

  return (
    <footer
      data-site-footer
      className="mt-20 border-t border-beige bg-ivory"
    >
      <div className="container-page grid gap-10 py-14 sm:grid-cols-2 lg:grid-cols-4">
        <div className="space-y-3">
          <Link href="/" className="flex items-center gap-2 text-lg font-semibold">
            <span aria-hidden="true" className="text-xl">🍅</span>
            <span className="font-display">{tc('siteName')}</span>
          </Link>
          <p className="max-w-xs text-sm text-muted">{t('tagline')}</p>
        </div>

        <nav aria-label={t('explore')} className="space-y-3">
          <h2 className="text-sm font-semibold text-charcoal">{t('explore')}</h2>
          <ul className="space-y-2">
            {EXPLORE.map((item) => (
              <li key={item.href}>
                <Link
                  href={item.href}
                  className="text-sm text-muted transition-colors hover:text-tomato"
                >
                  {tn(item.key)}
                </Link>
              </li>
            ))}
          </ul>
        </nav>

        <nav aria-label={t('about')} className="space-y-3">
          <h2 className="text-sm font-semibold text-charcoal">{t('about')}</h2>
          <ul className="space-y-2">
            <li>
              <Link
                href="/about"
                className="text-sm text-muted transition-colors hover:text-tomato"
              >
                {tn('about')}
              </Link>
            </li>
          </ul>
        </nav>

        <div className="space-y-3">
          <h2 className="text-sm font-semibold text-charcoal">{tc('siteName')}</h2>
          <p className="text-sm text-muted">{t('madeWith')}</p>
          <p className="text-xs text-muted/80">{t('attribution')}</p>
        </div>
      </div>

      <div className="border-t border-beige">
        <div className="container-page py-5 text-center text-xs text-muted">
          {t('rights', { year })}
        </div>
      </div>
    </footer>
  );
}
