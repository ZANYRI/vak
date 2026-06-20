/**
 * Build a compact list of page tokens for a pagination control, e.g.
 * [1, 'ellipsis', 4, 5, 6, 'ellipsis', 12].
 */
export type PageToken = number | 'ellipsis';

export function getPageRange(
  current: number,
  total: number,
  siblings = 1,
): PageToken[] {
  const totalNumbers = siblings * 2 + 5; // first, last, current, 2x ellipsis
  if (total <= totalNumbers) {
    return Array.from({ length: total }, (_, i) => i + 1);
  }

  const left = Math.max(current - siblings, 1);
  const right = Math.min(current + siblings, total);
  const showLeftDots = left > 2;
  const showRightDots = right < total - 1;

  const tokens: PageToken[] = [1];
  if (showLeftDots) tokens.push('ellipsis');

  for (let p = showLeftDots ? left : 2; p <= (showRightDots ? right : total - 1); p++) {
    tokens.push(p);
  }

  if (showRightDots) tokens.push('ellipsis');
  tokens.push(total);
  return tokens;
}

/** Parse a page number from a route segment, defaulting to 1. */
export function parsePage(value: string | undefined): number {
  const n = Number(value);
  return Number.isInteger(n) && n > 0 ? n : 1;
}
