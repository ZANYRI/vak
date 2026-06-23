"use client";

import { useTranslations } from "next-intl";
import { useSyncExternalStore } from "react";
import { motion } from "motion/react";
import {
  getSnapshot,
  subscribe,
  toggleFavorite
} from "@/lib/favorites";
import { cn } from "@/lib/cn";

export interface FavoriteButtonProps {
  slug: string;
  className?: string;
  /** Size in pixels of the icon button. */
  size?: number;
}

export function FavoriteButton({
  slug,
  className,
  size = 44
}: FavoriteButtonProps) {
  const t = useTranslations("recipe");
  // Subscribe to the favorites store so the heart stays in sync across
  // components and tabs. SSR returns an empty snapshot (no favorites),
  // so the button renders in its default state on the server.
  const slugs = useSyncExternalStore(subscribe, getSnapshot, getSnapshot);
  const active = slugs.includes(slug);

  const handleClick = () => {
    toggleFavorite(slug);
  };

  const label = active ? t("removeFromFavorites") : t("addToFavorites");

  return (
    <motion.button
      type="button"
      onClick={handleClick}
      aria-pressed={active}
      aria-label={label}
      title={label}
      whileTap={{ scale: 0.85 }}
      animate={active ? { scale: [1, 1.25, 1] } : { scale: 1 }}
      transition={{ duration: 0.3 }}
      className={cn(
        "inline-flex items-center justify-center rounded-full border transition-colors",
        active
          ? "border-primary bg-primary/10 text-primary"
          : "border-border bg-surface text-muted hover:text-primary",
        className
      )}
      style={{ width: size, height: size }}
    >
      <svg
        viewBox="0 0 24 24"
        width={size * 0.5}
        height={size * 0.5}
        fill={active ? "currentColor" : "none"}
        stroke="currentColor"
        strokeWidth={2}
        aria-hidden="true"
      >
        <path
          strokeLinecap="round"
          strokeLinejoin="round"
          d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z"
        />
      </svg>
    </motion.button>
  );
}
