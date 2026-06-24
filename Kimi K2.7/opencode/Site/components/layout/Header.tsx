"use client";

import { useLocale, useTranslations } from "next-intl";
import { getPathname, Link, usePathname } from "@/i18n/navigation";
import { cn } from "@/lib/utils";
import { LanguageSwitcher } from "./LanguageSwitcher";

const navItems = [
  { href: "/", label: "home" },
  { href: "/recipes", label: "recipes" },
  { href: "/categories", label: "categories" },
  { href: "/cuisines", label: "cuisines" },
  { href: "/favorites", label: "favorites" },
  { href: "/about", label: "about" },
] as const;

export function Header() {
  const t = useTranslations("nav");
  const locale = useLocale();
  const pathname = usePathname();

  return (
    <header className="border-border bg-background/95 sticky top-0 z-50 border-b backdrop-blur">
      <div className="mx-auto flex max-w-7xl items-center justify-between px-4 py-3 md:px-6">
        <Link
          href="/"
          className="text-primary text-xl font-bold tracking-tight"
        >
          Recipe Haven
        </Link>
        <nav aria-label={t("home")} className="hidden gap-1 md:flex">
          {navItems.map(({ href, label }) => {
            const localizedHref = getPathname({ href, locale });
            const isActive =
              pathname === localizedHref ||
              (href !== "/" && pathname.startsWith(`${localizedHref}/`));
            return (
              <Link
                key={href}
                href={href}
                className={cn(
                  "focus-visible:ring-accent rounded-md px-3 py-2 text-sm font-medium transition-colors focus-visible:ring-2 focus-visible:outline-none",
                  isActive
                    ? "bg-muted text-foreground"
                    : "text-foreground/70 hover:bg-muted hover:text-foreground",
                )}
              >
                {t(label)}
              </Link>
            );
          })}
        </nav>
        <LanguageSwitcher />
      </div>
    </header>
  );
}
