import { createNavigation } from 'next-intl/navigation';
import { routing } from './routing';

// Locale-aware navigation APIs. Use these instead of next/link & next/navigation
// so that the active locale is preserved automatically.
export const { Link, redirect, usePathname, useRouter, getPathname } =
  createNavigation(routing);
