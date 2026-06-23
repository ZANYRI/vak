"use client";

import { useSearchParams } from "next/navigation";
import { useTranslations } from "next-intl";
import { Link } from "@/i18n/navigation";
import { cn } from "@/lib/utils";

type RecipePaginationProps = {
  currentPage: number;
  totalPages: number;
};

export function RecipePagination({
  currentPage,
  totalPages,
}: RecipePaginationProps) {
  const t = useTranslations("pagination");
  const searchParams = useSearchParams();

  const makeHref = (page: number) => {
    const params = new URLSearchParams(searchParams.toString());
    const base = page === 1 ? "/recipes" : `/recipes/page/${page}`;
    return `${base}?${params.toString()}`;
  };

  const pages = Array.from({ length: totalPages }, (_, i) => i + 1);

  return (
    <nav
      aria-label={t("page", { page: currentPage })}
      className="flex items-center justify-center gap-2"
    >
      <Link
        href={makeHref(currentPage - 1)}
        aria-disabled={currentPage <= 1}
        tabIndex={currentPage <= 1 ? -1 : undefined}
        className={cn(
          "border-border bg-card hover:bg-muted focus-visible:ring-accent rounded-md border px-4 py-2 text-sm font-medium transition-colors focus-visible:ring-2 focus-visible:outline-none",
          currentPage <= 1 && "pointer-events-none opacity-50",
        )}
      >
        {t("previous")}
      </Link>

      {pages.map((page) => (
        <Link
          key={page}
          href={makeHref(page)}
          aria-current={page === currentPage ? "page" : undefined}
          aria-label={t("goToPage", { page })}
          className={cn(
            "focus-visible:ring-accent hidden rounded-md px-4 py-2 text-sm font-medium transition-colors focus-visible:ring-2 focus-visible:outline-none sm:inline-flex",
            page === currentPage
              ? "bg-primary text-primary-foreground"
              : "border-border bg-card hover:bg-muted border",
          )}
        >
          {page}
        </Link>
      ))}

      <Link
        href={makeHref(currentPage + 1)}
        aria-disabled={currentPage >= totalPages}
        tabIndex={currentPage >= totalPages ? -1 : undefined}
        className={cn(
          "border-border bg-card hover:bg-muted focus-visible:ring-accent rounded-md border px-4 py-2 text-sm font-medium transition-colors focus-visible:ring-2 focus-visible:outline-none",
          currentPage >= totalPages && "pointer-events-none opacity-50",
        )}
      >
        {t("next")}
      </Link>
    </nav>
  );
}
