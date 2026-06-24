"use client";

import { motion, useReducedMotion } from "framer-motion";
import { useTranslations } from "next-intl";
import { Link } from "@/i18n/navigation";
import type { Recipe } from "@/lib/validation";
import { formatDuration } from "@/lib/utils";
import { Badge } from "@/components/ui/Badge";
import { ImageWithFallback } from "./ImageWithFallback";
import { FavoriteButton } from "./FavoriteButton";

type RecipeCardProps = {
  recipe: Recipe;
  locale: "en" | "ru";
  index?: number;
};

export function RecipeCard({ recipe, locale, index = 0 }: RecipeCardProps) {
  const t = useTranslations("recipe");
  const shouldReduceMotion = useReducedMotion();

  const initial = shouldReduceMotion ? undefined : { opacity: 0, y: 20 };
  const whileInView = shouldReduceMotion ? undefined : { opacity: 1, y: 0 };

  return (
    <motion.article
      initial={initial}
      whileInView={whileInView}
      viewport={{ once: true }}
      transition={{ duration: 0.4, delay: index * 0.05 }}
      className="group border-border bg-card relative flex flex-col overflow-hidden rounded-lg border shadow-sm transition-shadow hover:shadow-md"
    >
      <Link
        href={`/recipes/${recipe.slug}`}
        className="relative aspect-[4/3] overflow-hidden"
      >
        <ImageWithFallback
          src={recipe.image}
          alt={recipe.imageAlt[locale]}
          fill
          className="group-hover:scale-105"
        />
        <div className="absolute top-2 right-2">
          <FavoriteButton
            id={recipe.id}
            className="bg-white/90 backdrop-blur"
          />
        </div>
      </Link>
      <div className="flex flex-1 flex-col p-4">
        <div className="mb-2 flex flex-wrap gap-2">
          <Badge variant="secondary">{recipe.cuisine}</Badge>
          <Badge>{recipe.difficulty}</Badge>
        </div>
        <h3 className="mb-1 text-lg leading-tight font-semibold">
          <Link
            href={`/recipes/${recipe.slug}`}
            className="focus-visible:ring-accent hover:underline focus-visible:ring-2 focus-visible:outline-none"
          >
            {recipe.title[locale]}
          </Link>
        </h3>
        <p className="text-foreground/70 mb-3 line-clamp-2 flex-1 text-sm">
          {recipe.description[locale]}
        </p>
        <div className="text-foreground/70 flex items-center justify-between text-sm">
          <span>
            {t("totalTime")}: {formatDuration(recipe.totalTimeMinutes)}
          </span>
          <span>★ {recipe.rating.toFixed(1)}</span>
        </div>
      </div>
    </motion.article>
  );
}
