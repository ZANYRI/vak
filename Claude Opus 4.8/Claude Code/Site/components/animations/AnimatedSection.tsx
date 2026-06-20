'use client';

import { motion, useReducedMotion } from 'motion/react';
import type { ReactNode } from 'react';

type Props = {
  children: ReactNode;
  className?: string;
  /** Stagger delay in seconds. */
  delay?: number;
  as?: 'div' | 'section' | 'article' | 'li';
};

/**
 * Fades and lifts its children into view once on scroll. Fully disabled when
 * the user prefers reduced motion.
 */
export function AnimatedSection({
  children,
  className,
  delay = 0,
  as = 'div',
}: Props) {
  const reduce = useReducedMotion();
  const MotionTag = motion[as];

  if (reduce) {
    const Tag = as;
    return <Tag className={className}>{children}</Tag>;
  }

  return (
    <MotionTag
      className={className}
      initial={{ opacity: 0, y: 16 }}
      whileInView={{ opacity: 1, y: 0 }}
      viewport={{ once: true, margin: '-60px' }}
      transition={{ duration: 0.5, delay, ease: [0.22, 1, 0.36, 1] }}
    >
      {children}
    </MotionTag>
  );
}
