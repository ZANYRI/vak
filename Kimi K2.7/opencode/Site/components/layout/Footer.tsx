"use client";

import { useTranslations } from "next-intl";
import { Link } from "@/i18n/navigation";

export function Footer() {
  const t = useTranslations("nav");

  return (
    <footer className="border-border bg-muted border-t">
      <div className="mx-auto max-w-7xl px-4 py-8 md:px-6">
        <div className="flex flex-col items-center justify-between gap-4 md:flex-row">
          <Link href="/" className="text-primary text-lg font-bold">
            Recipe Haven
          </Link>
          <nav className="text-foreground/70 flex flex-wrap justify-center gap-4 text-sm">
            <Link href="/" className="hover:text-foreground">
              {t("home")}
            </Link>
            <Link href="/recipes">{t("recipes")}</Link>
            <Link href="/about">{t("about")}</Link>
          </nav>
          <p className="text-foreground/60 text-xs">
            © {new Date().getFullYear()} Recipe Haven. Demo project.
          </p>
        </div>
      </div>
    </footer>
  );
}
