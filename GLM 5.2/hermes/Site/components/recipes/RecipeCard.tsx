"use client";

import { useLocale, useTranslations } from "next-intl";
import { Link } from "@/i18n/navigation";
import { motion } from "motion/react";
import { ImageWithFallback } from "@/components/ui/ImageWithFallback";
import { Badge } from "@/components/ui/Badge";
import type { Recipe } from "@/lib/types";

export interface RecipeCardProps {
  recipe: Recipe;
  priority?: boolean;
  index?: number;
}

export function RecipeCard({ recipe, priority = false, index = 0 }: RecipeCardProps) {
  const locale = useLocale() as "en" | "ru";
  const t = useTranslations();

  return (
    <motion.article
      initial={{ opacity: 0, y: 24 }}
      whileInView={{ opacity: 1, y: 0 }}
      viewport={{ once: true, margin: "-40px" }}
      transition={{ duration: 0.45, delay: Math.min(index * 0.06, 0.3), ease: [0.22, 1, 0.36, 1] }}
      className="group flex flex-col overflow-hidden rounded-card border border-border bg-card shadow-soft transition-shadow hover:shadow-lift"
    >
      <Link
        href={`/recipes/${recipe.slug}`}
        className="relative block aspect-[4/3] overflow-hidden"
        aria-label={recipe.title[locale]}
      >
        <ImageWithFallback
          src={recipe.image}
          alt={recipe.imageAlt[locale]}
          fill
          priority={priority}
          sizes="(max-width: 768px) 100vw, (max-width: 1200px) 50vw, 33vw"
          className="object-cover transition-transform duration-500 group-hover:scale-105"
        />
        <div className="pointer-events-none absolute inset-x-0 bottom-0 h-16 bg-gradient-to-t from-black/30 to-transparent" />
        <Badge
          variant="muted"
          className="absolute left-3 top-3 bg-white/90 backdrop-blur"
        >
          {t(`difficulty.${recipe.difficulty}`)}
        </Badge>
      </Link>

      <div className="flex flex-1 flex-col p-4">
        <h3 className="font-serif text-lg font-semibold leading-snug text-foreground">
          <Link
            href={`/recipes/${recipe.slug}`}
            className="after:absolute after:inset-0 after:content-['']"
          >
            {recipe.title[locale]}
          </Link>
        </h3>
        <p className="mt-1.5 line-clamp-2 text-sm text-muted">
          {recipe.description[locale]}
        </p>

        <dl className="mt-4 flex flex-wrap items-center gap-x-4 gap-y-1 text-xs text-muted">
          <div className="flex items-center gap-1">
            <dt aria-hidden="true">⏱</dt>
            <dd>
              {recipe.totalTimeMinutes} {t("common.minutes")}
            </dd>
          </div>
          <div className="flex items-center gap-1">
            <dt aria-hidden="true">★</dt>
            <dd>{recipe.rating.toFixed(1)}</dd>
          </div>
          <div className="flex items-center gap-1">
            <dt aria-hidden="true">🍽</dt>
            <dd>
              {recipe.servings} {t("common.servings")}
            </dd>
          </div>
        </dl>
      </div>
    </motion.article>
  );
}
