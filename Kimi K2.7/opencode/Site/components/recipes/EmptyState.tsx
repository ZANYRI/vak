"use client";

import { motion, useReducedMotion } from "framer-motion";

type EmptyStateProps = {
  title: string;
  description?: string;
};

export function EmptyState({ title, description }: EmptyStateProps) {
  const shouldReduceMotion = useReducedMotion();

  const initial = shouldReduceMotion ? undefined : { opacity: 0, scale: 0.98 };

  return (
    <motion.div
      initial={initial}
      animate={{ opacity: 1, scale: 1 }}
      transition={{ duration: 0.4 }}
      className="border-border bg-card flex flex-col items-center justify-center rounded-xl border border-dashed p-8 text-center md:p-12"
    >
      <div className="bg-muted mb-4 flex h-16 w-16 items-center justify-center rounded-full">
        <svg
          xmlns="http://www.w3.org/2000/svg"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          strokeWidth={1.5}
          className="text-foreground/50 h-8 w-8"
        >
          <path d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
        </svg>
      </div>
      <h2 className="text-xl font-semibold">{title}</h2>
      {description && (
        <p className="text-foreground/70 mt-2 max-w-md">{description}</p>
      )}
    </motion.div>
  );
}
