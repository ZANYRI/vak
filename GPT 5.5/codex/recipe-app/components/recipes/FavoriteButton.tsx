"use client";

import { useSyncExternalStore } from "react";

const storageKey = "misen-favorite-recipe-slugs";

function parseFavorites(value: string | null): string[] {
  try {
    const parsed: unknown = JSON.parse(value ?? "[]");
    return Array.isArray(parsed) ? parsed.filter((item): item is string => typeof item === "string") : [];
  } catch { return []; }
}

export function getFavoritesSnapshot() { return localStorage.getItem(storageKey) ?? "[]"; }
export function readFavorites(): string[] { return parseFavorites(getFavoritesSnapshot()); }
export function subscribeToFavorites(callback: () => void) {
  window.addEventListener("misen:favorites-changed", callback);
  window.addEventListener("storage", callback);
  return () => { window.removeEventListener("misen:favorites-changed", callback); window.removeEventListener("storage", callback); };
}

export function FavoriteButton({ slug, addLabel, removeLabel, compact = false }: { slug: string; addLabel: string; removeLabel: string; compact?: boolean }) {
  const snapshot = useSyncExternalStore(subscribeToFavorites, getFavoritesSnapshot, () => "[]");
  const saved = parseFavorites(snapshot).includes(slug);
  const toggle = () => {
    const current = readFavorites();
    const next = current.includes(slug) ? current.filter((item) => item !== slug) : [...current, slug];
    localStorage.setItem(storageKey, JSON.stringify(next));
    window.dispatchEvent(new Event("misen:favorites-changed"));
  };
  return <button type="button" className={`favorite-button ${saved ? "is-saved" : ""} ${compact ? "compact" : ""}`} onClick={toggle} aria-pressed={saved} aria-label={saved ? removeLabel : addLabel}><span aria-hidden="true">{saved ? "♥" : "♡"}</span>{compact ? null : <span>{saved ? removeLabel : addLabel}</span>}</button>;
}

export { storageKey };
