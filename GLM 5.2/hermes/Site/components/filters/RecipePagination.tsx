"use client";

import { useTranslations } from "next-intl";
import { Link } from "@/i18n/navigation";
import { getPageList } from "@/lib/pagination";
import { cn } from "@/lib/cn";

export interface RecipePaginationProps {
  basePath: string;
  queryString: string;
  currentPage: number;
  totalPages: number;
}

export function RecipePagination({
  basePath,
  queryString,
  currentPage,
  totalPages
}: RecipePaginationProps) {
  const t = useTranslations("pagination");
  if (totalPages <= 1) return null;

  const pages = getPageList(currentPage, totalPages);

  const buildHref = (page: number) => {
    const query = queryString.startsWith("?") ? queryString.slice(1) : queryString;
    const params = new URLSearchParams(query);
    if (page <= 1) params.delete("page");
    else params.set("page", String(page));
    const qs = params.toString();
    return qs ? `${basePath}?${qs}` : basePath;
  };

  const prevDisabled = currentPage <= 1;
  const nextDisabled = currentPage >= totalPages;

  return (
    <nav
      role="navigation"
      aria-label={t("page", { number: currentPage })}
      className="mt-10 flex items-center justify-center gap-1.5"
    >
      <PageLink
        href={buildHref(currentPage - 1)}
        disabled={prevDisabled}
        ariaLabel={t("previous")}
      >
        <svg viewBox="0 0 24 24" width="18" height="18" fill="none" stroke="currentColor" strokeWidth={2} aria-hidden="true">
          <path strokeLinecap="round" strokeLinejoin="round" d="M15 18l-6-6 6-6" />
        </svg>
      </PageLink>

      {pages.map((p, i) =>
        p === null ? (
          <span key={`gap-${i}`} className="px-1.5 text-muted" aria-hidden="true">
            …
          </span>
        ) : (
          <PageLink
            key={p}
            href={buildHref(p)}
            active={p === currentPage}
            ariaLabel={p === currentPage ? t("currentPage", { number: p }) : t("page", { number: p })}
            ariaCurrent={p === currentPage ? "page" : undefined}
          >
            {p}
          </PageLink>
        )
      )}

      <PageLink
        href={buildHref(currentPage + 1)}
        disabled={nextDisabled}
        ariaLabel={t("next")}
      >
        <svg viewBox="0 0 24 24" width="18" height="18" fill="none" stroke="currentColor" strokeWidth={2} aria-hidden="true">
          <path strokeLinecap="round" strokeLinejoin="round" d="M9 6l6 6-6 6" />
        </svg>
      </PageLink>
    </nav>
  );
}

function PageLink({
  href,
  children,
  active,
  disabled,
  ariaLabel,
  ariaCurrent
}: {
  href: string;
  children: React.ReactNode;
  active?: boolean;
  disabled?: boolean;
  ariaLabel: string;
  ariaCurrent?: "page";
}) {
  const base =
    "inline-flex h-10 min-w-10 items-center justify-center rounded-full px-3 text-sm font-medium transition-colors";
  if (disabled) {
    return (
      <span
        aria-disabled="true"
        aria-label={ariaLabel}
        className={cn(base, "cursor-not-allowed text-muted/40")}
      >
        {children}
      </span>
    );
  }
  return (
    <Link
      href={href}
      aria-label={ariaLabel}
      aria-current={ariaCurrent}
      className={cn(
        base,
        active
          ? "bg-primary text-primary-foreground shadow-soft"
          : "border border-border bg-surface text-foreground hover:bg-surface-alt"
      )}
    >
      {children}
    </Link>
  );
}
