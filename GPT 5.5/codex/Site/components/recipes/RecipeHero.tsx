import type { ReactNode } from "react";

export function RecipeHero({ children }: { children: ReactNode }) {
  return <section className="recipe-hero">{children}</section>;
}
