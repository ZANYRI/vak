'use client';

import {
  createContext,
  useCallback,
  useContext,
  useMemo,
  useSyncExternalStore,
} from 'react';

const STORAGE_KEY = 'savora:favorites';

/**
 * A tiny external store backed by localStorage. Using useSyncExternalStore (the
 * idiomatic React 19 approach) avoids reading storage during render and keeps
 * every tab in sync without setState-in-effect patterns.
 */

// Distinct constants so we can tell "server render" from "hydrated, empty".
const SERVER_SNAPSHOT: string[] = [];
let snapshot: string[] = [];
let lastRaw: string | null = null;
const listeners = new Set<() => void>();
let initialized = false;

function parse(raw: string | null): string[] {
  if (!raw) return [];
  try {
    const parsed: unknown = JSON.parse(raw);
    return Array.isArray(parsed)
      ? parsed.filter((x): x is string => typeof x === 'string')
      : [];
  } catch {
    return [];
  }
}

function getSnapshot(): string[] {
  try {
    const raw = window.localStorage.getItem(STORAGE_KEY);
    // Only recompute (new array reference) when the raw value actually changes,
    // so useSyncExternalStore sees a stable snapshot otherwise.
    if (raw !== lastRaw || !initialized) {
      lastRaw = raw;
      snapshot = parse(raw);
      initialized = true;
    }
    return snapshot;
  } catch {
    return snapshot;
  }
}

function getServerSnapshot(): string[] {
  return SERVER_SNAPSHOT;
}

function emit() {
  for (const listener of listeners) listener();
}

function subscribe(listener: () => void): () => void {
  listeners.add(listener);
  const onStorage = (event: StorageEvent) => {
    if (event.key === STORAGE_KEY) emit();
  };
  window.addEventListener('storage', onStorage);
  return () => {
    listeners.delete(listener);
    window.removeEventListener('storage', onStorage);
  };
}

function writeFavorites(next: string[]) {
  try {
    window.localStorage.setItem(STORAGE_KEY, JSON.stringify(next));
  } catch {
    /* storage may be unavailable (private mode) — ignore */
  }
  lastRaw = JSON.stringify(next);
  snapshot = next;
  initialized = true;
  emit();
}

type FavoritesContextValue = {
  favorites: string[];
  /** True once we've read localStorage on the client — avoids a hydration flash. */
  ready: boolean;
  isFavorite: (slug: string) => boolean;
  toggleFavorite: (slug: string) => void;
};

const FavoritesContext = createContext<FavoritesContextValue | null>(null);

export function FavoritesProvider({ children }: { children: React.ReactNode }) {
  const favorites = useSyncExternalStore(
    subscribe,
    getSnapshot,
    getServerSnapshot,
  );
  const ready = favorites !== SERVER_SNAPSHOT;

  const toggleFavorite = useCallback((slug: string) => {
    const current = getSnapshot();
    const next = current.includes(slug)
      ? current.filter((s) => s !== slug)
      : [...current, slug];
    writeFavorites(next);
  }, []);

  const value = useMemo<FavoritesContextValue>(
    () => ({
      favorites,
      ready,
      isFavorite: (slug: string) => favorites.includes(slug),
      toggleFavorite,
    }),
    [favorites, ready, toggleFavorite],
  );

  return <FavoritesContext value={value}>{children}</FavoritesContext>;
}

export function useFavorites(): FavoritesContextValue {
  const ctx = useContext(FavoritesContext);
  if (!ctx) {
    throw new Error('useFavorites must be used within a FavoritesProvider');
  }
  return ctx;
}
