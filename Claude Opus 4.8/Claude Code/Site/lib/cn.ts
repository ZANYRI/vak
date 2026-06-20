type ClassValue = string | number | false | null | undefined;

/** Minimal classnames joiner — no external dependency needed. */
export function cn(...values: ClassValue[]): string {
  return values.filter(Boolean).join(' ');
}
