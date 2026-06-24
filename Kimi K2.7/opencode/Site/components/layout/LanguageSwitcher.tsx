"use client";

import { useLocale, useTranslations } from "next-intl";
import { Link, usePathname } from "@/i18n/navigation";
import { cn } from "@/lib/utils";

const locales = ["en", "ru"] as const;

export function LanguageSwitcher() {
  const t = useTranslations("language");
  const pathname = usePathname();
  const currentLocale = useLocale();

  return (
    <nav
      aria-label={t("label")}
      className="bg-muted flex items-center gap-1 rounded-md p-1"
    >
      {locales.map((locale) => {
        const isActive = locale === currentLocale;
        return (
          <Link
            key={locale}
            href={pathname}
            locale={locale}
            aria-current={isActive ? "page" : undefined}
            className={cn(
              "focus-visible:ring-accent rounded px-2 py-1 text-sm font-medium uppercase transition-colors focus-visible:ring-2 focus-visible:outline-none",
              isActive
                ? "bg-card text-foreground shadow-sm"
                : "text-foreground/70 hover:bg-card/50 hover:text-foreground",
            )}
          >
            {locale}
          </Link>
        );
      })}
    </nav>
  );
}
