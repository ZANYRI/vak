import type { MetadataRoute } from 'next';
import { locales } from '@/i18n/routing';
import { getAllSlugs, getTotalPages } from '@/lib/recipes';
import { CATEGORIES, CUISINES } from '@/lib/constants';
import { localizedUrl } from '@/lib/seo';

const STATIC_PATHS = ['', '/recipes', '/categories', '/cuisines', '/favorites', '/about'];

export default function sitemap(): MetadataRoute.Sitemap {
  const entries: MetadataRoute.Sitemap = [];

  const addForAllLocales = (path: string, priority: number) => {
    for (const locale of locales) {
      entries.push({
        url: localizedUrl(locale, path),
        changeFrequency: 'weekly',
        priority,
        alternates: {
          languages: Object.fromEntries(
            locales.map((l) => [l, localizedUrl(l, path)]),
          ),
        },
      });
    }
  };

  for (const path of STATIC_PATHS) {
    addForAllLocales(path, path === '' ? 1 : 0.8);
  }

  const totalPages = getTotalPages();
  for (let p = 2; p <= totalPages; p++) {
    addForAllLocales(`/recipes/page/${p}`, 0.5);
  }

  for (const slug of getAllSlugs()) {
    addForAllLocales(`/recipes/${slug}`, 0.7);
  }

  for (const category of CATEGORIES) {
    addForAllLocales(`/categories/${category}`, 0.6);
  }

  for (const cuisine of CUISINES) {
    addForAllLocales(`/cuisines/${cuisine}`, 0.6);
  }

  return entries;
}
