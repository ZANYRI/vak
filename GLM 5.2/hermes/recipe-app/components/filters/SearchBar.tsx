"use client";

import { useTranslations } from "next-intl";
import { useRouter } from "next/navigation";
import { useState, useTransition } from "react";
import { cn } from "@/lib/cn";

export interface SearchBarProps {
  basePath: string;
  initialQuery?: string;
  placeholder?: string;
  className?: string;
  /** When true, the form navigates to `basePath?search=…`. */
  navigateOnSubmit?: boolean;
}

export function SearchBar({
  basePath,
  initialQuery = "",
  placeholder,
  className,
  navigateOnSubmit = true
}: SearchBarProps) {
  const t = useTranslations("filters");
  const router = useRouter();
  const [value, setValue] = useState(initialQuery);
  const [isPending, startTransition] = useTransition();

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!navigateOnSubmit) return;
    const params = new URLSearchParams();
    if (value.trim()) params.set("search", value.trim());
    const qs = params.toString();
    startTransition(() => {
      router.replace(qs ? `${basePath}?${qs}` : basePath);
    });
  };

  return (
    <form
      role="search"
      onSubmit={handleSubmit}
      className={cn("relative w-full", className)}
    >
      <label htmlFor="recipe-search" className="sr-only">
        {t("search")}
      </label>
      <span className="pointer-events-none absolute left-4 top-1/2 -translate-y-1/2 text-muted">
        <svg viewBox="0 0 24 24" width="20" height="20" fill="none" stroke="currentColor" strokeWidth={2} aria-hidden="true">
          <circle cx="11" cy="11" r="7" />
          <path strokeLinecap="round" d="M21 21l-4.3-4.3" />
        </svg>
      </span>
      <input
        id="recipe-search"
        type="search"
        value={value}
        onChange={(e) => setValue(e.target.value)}
        placeholder={placeholder ?? t("searchPlaceholder")}
        className="h-12 w-full rounded-full border border-border bg-surface pl-12 pr-4 text-sm text-foreground placeholder:text-muted transition-colors focus:border-primary focus:outline-none focus:ring-2 focus:ring-primary/30"
      />
      {isPending ? (
        <span className="absolute right-4 top-1/2 -translate-y-1/2 text-xs text-muted">
          …
        </span>
      ) : null}
    </form>
  );
}
