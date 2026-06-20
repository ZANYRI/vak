'use client';

import { motion, useReducedMotion } from 'motion/react';
import { usePathname } from '@/i18n/navigation';
import type { ReactNode } from 'react';

/**
 * Subtle fade/slide applied to page content on each navigation. The pathname is
 * used as a key so the animation re-runs when the route changes.
 */
export function PageTransition({ children }: { children: ReactNode }) {
  const pathname = usePathname();
  const reduce = useReducedMotion();

  if (reduce) return <>{children}</>;

  return (
    <motion.div
      key={pathname}
      initial={{ opacity: 0, y: 8 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.35, ease: [0.22, 1, 0.36, 1] }}
    >
      {children}
    </motion.div>
  );
}
