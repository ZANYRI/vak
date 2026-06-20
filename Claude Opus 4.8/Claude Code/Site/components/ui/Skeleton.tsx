import type { HTMLAttributes } from 'react';
import { cn } from '@/lib/cn';

export function Skeleton({ className, ...props }: HTMLAttributes<HTMLDivElement>) {
  return (
    <div
      aria-hidden="true"
      className={cn('skeleton rounded-md', className)}
      {...props}
    />
  );
}

/** Card-shaped placeholder used while a recipe grid is loading. */
export function RecipeCardSkeleton() {
  return (
    <div className="overflow-hidden rounded-card border border-beige bg-card shadow-soft">
      <Skeleton className="aspect-[4/3] w-full rounded-none" />
      <div className="space-y-3 p-5">
        <Skeleton className="h-4 w-24" />
        <Skeleton className="h-6 w-3/4" />
        <Skeleton className="h-4 w-full" />
        <Skeleton className="h-4 w-2/3" />
      </div>
    </div>
  );
}
