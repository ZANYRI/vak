"use client";

const STORAGE_KEY = "saveur.favorites";

/**
 * Favorites are stored in localStorage and exposed to React via
 * `useSyncExternalStore`. The snapshot array is cached module-side so the
 * reference stays stable across renders until the data actually changes.
 */

let cachedSnapshot: string[] = [];
let initialized = false;

function readRaw(): string[] {
  if (typeof window === "undefined") return [];
  try {
    const raw = window.localStorage.getItem(STORAGE_KEY);
    if (!raw) return [];
    const parsed = JSON.parse(raw);
    if (!Array.isArray(parsed)) return [];
    return parsed.filter((v): v is string => typeof v === "string");
  } catch {
    return [];
  }
}

function ensureInitialized() {
  if (!initialized && typeof window !== "undefined") {
    cachedSnapshot = readRaw();
    initialized = true;
  }
}

const listeners = new Set<() => void>();

function notify() {
  for (const l of listeners) l();
}

function writeStorage(slugs: string[]): void {
  if (typeof window === "undefined") return;
  try {
    window.localStorage.setItem(STORAGE_KEY, JSON.stringify(slugs));
  } catch {
    /* ignore quota / privacy mode errors */
  }
  cachedSnapshot = slugs;
  notify();
}

export function subscribe(callback: () => void): () => void {
  ensureInitialized();
  listeners.add(callback);
  const storageHandler = () => {
    cachedSnapshot = readRaw();
    notify();
  };
  window.addEventListener("storage", storageHandler);
  return () => {
    listeners.delete(callback);
    window.removeEventListener("storage", storageHandler);
  };
}

export function getSnapshot(): string[] {
  ensureInitialized();
  return cachedSnapshot;
}

export function getServerSnapshot(): string[] {
  return [];
}

export function getFavoriteSlugs(): string[] {
  return getSnapshot();
}

export function isFavorite(slug: string): boolean {
  return getSnapshot().includes(slug);
}

export function toggleFavorite(slug: string): boolean {
  const current = getSnapshot();
  const exists = current.includes(slug);
  const next = exists ? current.filter((s) => s !== slug) : [...current, slug];
  writeStorage(next);
  return !exists;
}

export function addFavorite(slug: string): void {
  const current = getSnapshot();
  if (!current.includes(slug)) {
    writeStorage([...current, slug]);
  }
}

export function removeFavorite(slug: string): void {
  writeStorage(getSnapshot().filter((s) => s !== slug));
}

/** Back-compat alias kept for older call sites. */
export function subscribeToFavorites(callback: () => void): () => void {
  return subscribe(callback);
}
