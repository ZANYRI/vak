"use client";

import { useFavorites } from "@/components/providers/FavoritesProvider";
import { cn } from "@/lib/utils";

type FavoriteButtonProps = {
  id: string;
  className?: string;
};

export function FavoriteButton({ id, className }: FavoriteButtonProps) {
  const { isFavorite, toggle } = useFavorites();
  const active = isFavorite(id);

  return (
    <button
      type="button"
      onClick={() => toggle(id)}
      aria-pressed={active}
      aria-label={active ? "Remove from favorites" : "Add to favorites"}
      className={cn(
        "bg-card hover:bg-muted focus-visible:ring-accent flex h-10 w-10 items-center justify-center rounded-full shadow-sm transition-colors focus-visible:ring-2 focus-visible:outline-none",
        active && "text-red-500",
        className,
      )}
    >
      <svg
        xmlns="http://www.w3.org/2000/svg"
        viewBox="0 0 24 24"
        fill={active ? "currentColor" : "none"}
        stroke="currentColor"
        strokeWidth={2}
        className="h-5 w-5"
      >
        <path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78Z" />
      </svg>
    </button>
  );
}
