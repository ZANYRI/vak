"use client";

import { useLocale, useTranslations } from "next-intl";
import { ImageWithFallback } from "@/components/ui/ImageWithFallback";
import { FavoriteButton } from "@/components/ui/FavoriteButton";
import { Badge } from "@/components/ui/Badge";
import { motion } from "motion/react";
import type { Recipe } from "@/lib/types";

export interface RecipeHeroProps {
  recipe: Recipe;
}

export function RecipeHero({ recipe }: RecipeHeroProps) {
  const locale = useLocale() as "en" | "ru";
  const t = useTranslations();

  return (
    <header className="relative overflow-hidden rounded-card border border-border bg-card shadow-soft">
      <div className="relative aspect-[16/9] w-full sm:aspect-[21/9]">
        <ImageWithFallback
          src={recipe.image}
          alt={recipe.imageAlt[locale]}
          fill
          priority
          sizes="(max-width: 768px) 100vw, 1200px"
          className="object-cover"
        />
        <div className="absolute inset-0 bg-gradient-to-t from-black/70 via-black/20 to-transparent" />
        <div className="absolute inset-x-0 bottom-0 p-5 sm:p-8">
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.5, ease: [0.22, 1, 0.36, 1] }}
            className="flex flex-wrap items-center gap-2"
          >
            <Badge variant="primary" className="bg-white/90 text-primary">
              {t(`difficulty.${recipe.difficulty}`)}
            </Badge>
            <Badge variant="muted" className="bg-white/80 text-foreground">
              {recipe.totalTimeMinutes} min
            </Badge>
          </motion.div>
          <motion.h1
            initial={{ opacity: 0, y: 24 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.55, delay: 0.05, ease: [0.22, 1, 0.36, 1] }}
            className="mt-3 max-w-3xl font-serif text-3xl font-semibold leading-tight text-white sm:text-4xl"
          >
            {recipe.title[locale]}
          </motion.h1>
          <p className="mt-2 max-w-2xl text-sm text-white/85 sm:text-base">
            {recipe.description[locale]}
          </p>
        </div>
        <div className="absolute right-4 top-4">
          <FavoriteButton slug={recipe.slug} className="backdrop-blur" />
        </div>
      </div>
    </header>
  );
}
