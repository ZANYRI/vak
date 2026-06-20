'use client';

import { useTranslations } from 'next-intl';
import { Suspense, useState } from 'react';
import { AnimatePresence, motion } from 'motion/react';
import { Link, usePathname } from '@/i18n/navigation';
import { cn } from '@/lib/cn';
import { LanguageSwitcher } from './LanguageSwitcher';

const NAV = [
  { href: '/', key: 'home' as const },
  { href: '/recipes', key: 'recipes' as const },
  { href: '/categories', key: 'categories' as const },
  { href: '/cuisines', key: 'cuisines' as const },
  { href: '/favorites', key: 'favorites' as const },
  { href: '/about', key: 'about' as const },
];

export function Header() {
  const t = useTranslations('Nav');
  const tc = useTranslations('Common');
  const pathname = usePathname();
  const [open, setOpen] = useState(false);

  const isActive = (href: string) =>
    href === '/' ? pathname === '/' : pathname.startsWith(href);

  return (
    <header
      data-site-header
      className="sticky top-0 z-50 border-b border-beige/80 bg-cream/85 backdrop-blur-md"
    >
      <div className="container-page flex h-16 items-center justify-between gap-4">
        <Link
          href="/"
          className="flex items-center gap-2 text-xl font-semibold tracking-tight text-charcoal"
        >
          <span aria-hidden="true" className="text-2xl">🍅</span>
          <span className="font-display">{tc('siteName')}</span>
        </Link>

        <nav aria-label={t('primary')} className="hidden md:block">
          <ul className="flex items-center gap-1">
            {NAV.map((item) => (
              <li key={item.href}>
                <Link
                  href={item.href}
                  aria-current={isActive(item.href) ? 'page' : undefined}
                  className={cn(
                    'relative rounded-full px-3.5 py-2 text-sm font-medium transition-colors',
                    isActive(item.href)
                      ? 'text-tomato'
                      : 'text-charcoal/80 hover:text-tomato',
                  )}
                >
                  {t(item.key)}
                  {isActive(item.href) && (
                    <motion.span
                      layoutId="nav-underline"
                      className="absolute inset-x-3 -bottom-0.5 h-0.5 rounded-full bg-tomato"
                    />
                  )}
                </Link>
              </li>
            ))}
          </ul>
        </nav>

        <div className="flex items-center gap-2">
          <Suspense
            fallback={
              <div className="h-9 w-[88px] rounded-full border border-beige-dark bg-card" />
            }
          >
            <LanguageSwitcher />
          </Suspense>
          <button
            type="button"
            className="inline-flex h-10 w-10 items-center justify-center rounded-full border border-beige-dark bg-card text-charcoal md:hidden"
            aria-expanded={open}
            aria-controls="mobile-nav"
            aria-label={open ? t('closeMenu') : t('openMenu')}
            onClick={() => setOpen((v) => !v)}
          >
            <svg viewBox="0 0 24 24" className="h-5 w-5" fill="none" stroke="currentColor" strokeWidth="2">
              {open ? (
                <path d="M6 6l12 12M18 6L6 18" strokeLinecap="round" />
              ) : (
                <path d="M4 7h16M4 12h16M4 17h16" strokeLinecap="round" />
              )}
            </svg>
          </button>
        </div>
      </div>

      <AnimatePresence>
        {open && (
          <motion.nav
            id="mobile-nav"
            aria-label={t('primary')}
            className="md:hidden"
            initial={{ height: 0, opacity: 0 }}
            animate={{ height: 'auto', opacity: 1 }}
            exit={{ height: 0, opacity: 0 }}
            transition={{ duration: 0.25, ease: 'easeOut' }}
          >
            <ul className="container-page flex flex-col gap-1 border-t border-beige/80 py-3">
              {NAV.map((item) => (
                <li key={item.href}>
                  <Link
                    href={item.href}
                    onClick={() => setOpen(false)}
                    aria-current={isActive(item.href) ? 'page' : undefined}
                    className={cn(
                      'block rounded-xl px-4 py-3 text-sm font-medium transition-colors',
                      isActive(item.href)
                        ? 'bg-tomato/10 text-tomato'
                        : 'text-charcoal hover:bg-beige/60',
                    )}
                  >
                    {t(item.key)}
                  </Link>
                </li>
              ))}
            </ul>
          </motion.nav>
        )}
      </AnimatePresence>
    </header>
  );
}
