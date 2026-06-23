"use client";

import { useLocale, useTranslations } from "next-intl";
import { usePathname, useRouter } from "@/i18n/navigation";
import { locales, type Locale } from "@/i18n/routing";
import { useTransition } from "react";
import { cn } from "@/lib/cn";

export function LanguageSwitcher({ className }: { className?: string }) {
  const t = useTranslations("common");
  const locale = useLocale() as Locale;
  const pathname = usePathname();
  const router = useRouter();
  const [isPending, startTransition] = useTransition();

  const switchTo = (next: Locale) => {
    if (next === locale) return;
    startTransition(() => {
      // Passing the current pathname preserves the route across locales.
      router.replace(pathname, { locale: next });
    });
  };

  return (
    <div
      role="group"
      aria-label={t("language")}
      className={cn(
        "inline-flex items-center rounded-full border border-border bg-surface p-0.5 text-sm",
        className
      )}
    >
      {locales.map((l) => {
        const isActive = l === locale;
        return (
          <button
            key={l}
            type="button"
            onClick={() => switchTo(l)}
            aria-current={isActive ? "true" : undefined}
            aria-pressed={isActive}
            disabled={isPending}
            className={cn(
              "min-w-9 rounded-full px-3 py-1 font-medium uppercase tracking-wide transition-colors",
              isActive
                ? "bg-primary text-primary-foreground shadow-soft"
                : "text-muted hover:text-foreground",
              isPending && "opacity-60"
            )}
          >
            {l}
          </button>
        );
      })}
    </div>
  );
}
