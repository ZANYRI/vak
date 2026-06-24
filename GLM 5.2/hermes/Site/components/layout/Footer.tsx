import { useTranslations } from "next-intl";
import { Link } from "@/i18n/navigation";

export function Footer() {
  const t = useTranslations();
  const year = new Date().getFullYear();

  return (
    <footer className="mt-20 border-t border-border bg-surface-alt">
      <div className="container-page grid gap-8 py-12 sm:grid-cols-2 lg:grid-cols-3">
        <div>
          <p className="font-serif text-lg font-semibold text-foreground">
            {t("common.siteName")}
          </p>
          <p className="mt-2 max-w-xs text-sm text-muted">
            {t("common.tagline")}
          </p>
        </div>
        <nav aria-label={t("common.menu")} className="flex flex-col gap-2">
          <Link href="/recipes" className="text-sm text-muted hover:text-foreground">
            {t("nav.recipes")}
          </Link>
          <Link href="/categories" className="text-sm text-muted hover:text-foreground">
            {t("nav.categories")}
          </Link>
          <Link href="/cuisines" className="text-sm text-muted hover:text-foreground">
            {t("nav.cuisines")}
          </Link>
          <Link href="/about" className="text-sm text-muted hover:text-foreground">
            {t("nav.about")}
          </Link>
        </nav>
        <p className="text-xs text-muted sm:col-span-2 lg:col-span-1 lg:text-right">
          © {year} {t("common.siteName")}. Next.js · React · Tailwind CSS.
        </p>
      </div>
    </footer>
  );
}
