import { RecipeCard } from "./RecipeCard";
import { Skeleton } from "@/components/ui/Skeleton";
import type { Recipe } from "@/lib/types";

export interface RecipeGridProps {
  recipes: Recipe[];
  loading?: boolean;
  skeletonCount?: number;
}

export function RecipeGrid({
  recipes,
  loading = false,
  skeletonCount = 6,
}: RecipeGridProps) {
  if (loading) {
    return (
      <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
        {Array.from({ length: skeletonCount }).map((_, i) => (
          <div
            key={i}
            className="overflow-hidden rounded-card border border-border bg-card"
          >
            <Skeleton className="aspect-[4/3] w-full rounded-none" />
            <div className="space-y-3 p-4">
              <Skeleton className="h-5 w-3/4" />
              <Skeleton className="h-4 w-full" />
              <Skeleton className="h-4 w-2/3" />
            </div>
          </div>
        ))}
      </div>
    );
  }

  return (
    <ul className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3" role="list">
      {recipes.map((recipe, i) => (
        <li key={recipe.slug} className="relative">
          <RecipeCard recipe={recipe} index={i} priority={i < 3} />
        </li>
      ))}
    </ul>
  );
}
