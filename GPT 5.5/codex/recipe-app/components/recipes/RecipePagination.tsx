import Link from "next/link";
import type { Dictionary } from "@/i18n/dictionary";
import type { Locale } from "@/i18n/config";

type Props = { locale: Locale; page: number; totalPages: number; query: Record<string, string | undefined>; dictionary: Dictionary };

function href(locale: Locale, page: number, query: Props["query"]) {
  const params = new URLSearchParams();
  Object.entries(query).forEach(([key, value]) => { if (value) params.set(key, value); });
  const suffix = params.size ? `?${params}` : "";
  return `/${locale}/recipes${page === 1 ? "" : `/page/${page}`}${suffix}`;
}

export function RecipePagination({ locale, page, totalPages, query, dictionary }: Props) {
  if (totalPages < 2) return null;
  const numbers = Array.from({ length: totalPages }, (_, index) => index + 1);
  const link = (target: number, label: string, disabled = false) => disabled ? <span className="pagination-button is-disabled" aria-disabled="true">{label}</span> : <Link className="pagination-button" href={href(locale, target, query)}>{label}</Link>;
  return <nav className="pagination" aria-label="Pagination"><div>{link(page - 1, `← ${dictionary.recipes.previous}`, page === 1)}</div><div className="page-numbers">{numbers.map((number) => number === page ? <span key={number} className="pagination-button is-active" aria-current="page"><span className="sr-only">{dictionary.recipes.page} </span>{number}</span> : <Link key={number} className="pagination-button" href={href(locale, number, query)}>{number}</Link>)}</div><div>{link(page + 1, `${dictionary.recipes.next} →`, page === totalPages)}</div></nav>;
}
