import { PAGE_SIZE } from "./types";

export function getTotalPages(total: number, pageSize: number = PAGE_SIZE): number {
  if (total <= 0) return 1;
  return Math.ceil(total / pageSize);
}

export function clampPage(page: number, totalPages: number): number {
  if (page < 1) return 1;
  if (page > totalPages) return totalPages;
  return page;
}

export function getPageRange(total: number, page: number, pageSize: number = PAGE_SIZE) {
  const from = total === 0 ? 0 : (page - 1) * pageSize + 1;
  const to = Math.min(page * pageSize, total);
  return { from, to };
}

/**
 * Returns a windowed list of page numbers around `current`, with ellipses
 * represented as `null`. Used by the pagination component.
 */
export function getPageList(current: number, totalPages: number): (number | null)[] {
  if (totalPages <= 7) {
    return Array.from({ length: totalPages }, (_, i) => i + 1);
  }
  const pages: (number | null)[] = [1];
  const left = Math.max(2, current - 1);
  const right = Math.min(totalPages - 1, current + 1);
  if (left > 2) pages.push(null);
  for (let i = left; i <= right; i++) pages.push(i);
  if (right < totalPages - 1) pages.push(null);
  pages.push(totalPages);
  return pages;
}

export function buildPageHref(
  baseHref: string,
  queryString: string,
  page: number
): string {
  const query = queryString.startsWith("?") ? queryString.slice(1) : queryString;
  const params = new URLSearchParams(query);
  if (page <= 1) {
    params.delete("page");
  } else {
    params.set("page", String(page));
  }
  const search = params.toString();
  return search ? `${baseHref}?${search}` : baseHref;
}
