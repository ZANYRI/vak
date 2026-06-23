"use client";

import {
  createContext,
  useCallback,
  useContext,
  useMemo,
  useState,
  type ReactNode,
} from "react";

type FavoritesContextValue = {
  favorites: string[];
  isFavorite: (id: string) => boolean;
  toggle: (id: string) => void;
  add: (id: string) => void;
  remove: (id: string) => void;
};

const FavoritesContext = createContext<FavoritesContextValue | null>(null);

const STORAGE_KEY = "recipe-app-favorites";

function readStoredFavorites(): string[] {
  if (typeof window === "undefined") return [];
  try {
    const stored = window.localStorage.getItem(STORAGE_KEY);
    return stored ? (JSON.parse(stored) as string[]) : [];
  } catch {
    return [];
  }
}

function writeStoredFavorites(favorites: string[]) {
  if (typeof window === "undefined") return;
  try {
    window.localStorage.setItem(STORAGE_KEY, JSON.stringify(favorites));
  } catch {
    // ignore storage errors
  }
}

export function FavoritesProvider({ children }: { children: ReactNode }) {
  const [favorites, setFavorites] = useState<string[]>(() =>
    readStoredFavorites(),
  );

  const isFavorite = useCallback(
    (id: string) => favorites.includes(id),
    [favorites],
  );

  const toggle = useCallback((id: string) => {
    setFavorites((prev) => {
      const next = prev.includes(id)
        ? prev.filter((item) => item !== id)
        : [...prev, id];
      writeStoredFavorites(next);
      return next;
    });
  }, []);

  const add = useCallback((id: string) => {
    setFavorites((prev) => {
      if (prev.includes(id)) return prev;
      const next = [...prev, id];
      writeStoredFavorites(next);
      return next;
    });
  }, []);

  const remove = useCallback((id: string) => {
    setFavorites((prev) => {
      const next = prev.filter((item) => item !== id);
      writeStoredFavorites(next);
      return next;
    });
  }, []);

  const value = useMemo(
    () => ({ favorites, isFavorite, toggle, add, remove }),
    [favorites, isFavorite, toggle, add, remove],
  );

  return (
    <FavoritesContext.Provider value={value}>
      {children}
    </FavoritesContext.Provider>
  );
}

export function useFavorites(): FavoritesContextValue {
  const context = useContext(FavoritesContext);
  if (!context) {
    throw new Error("useFavorites must be used within FavoritesProvider");
  }
  return context;
}
