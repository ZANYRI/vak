// Next.js 16 renamed the `middleware` file convention to `proxy`. next-intl's
// request handler is still imported from `next-intl/middleware`; it works
// unchanged as the proxy default export (locale detection & redirects).
import createMiddleware from 'next-intl/middleware';
import { routing } from './i18n/routing';

export default createMiddleware(routing);

export const config = {
  // Match all pathnames except for
  // - API routes (/api, /trpc)
  // - Next.js internals (/_next, /_vercel)
  // - files with an extension (e.g. favicon.ico, images)
  matcher: '/((?!api|trpc|_next|_vercel|.*\\..*).*)',
};
