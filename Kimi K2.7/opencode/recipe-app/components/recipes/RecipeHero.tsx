"use client";

import { motion, useReducedMotion } from "framer-motion";
import type { Recipe } from "@/lib/validation";
import { ImageWithFallback } from "./ImageWithFallback";

type RecipeHeroProps = {
  recipe: Recipe;
  locale: "en" | "ru";
};

export function RecipeHero({ recipe, locale }: RecipeHeroProps) {
  const shouldReduceMotion = useReducedMotion();
  const initial = shouldReduceMotion ? undefined : { opacity: 0, y: 24 };

  return (
    <motion.div
      initial={initial}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.6, ease: "easeOut" }}
      className="grid gap-8 lg:grid-cols-2"
    >
      <div className="relative aspect-[4/3] overflow-hidden rounded-xl lg:aspect-square">
        <ImageWithFallback
          src={recipe.image}
          alt={recipe.imageAlt[locale]}
          fill
          priority
        />
      </div>
      <div className="flex flex-col justify-center">
        <h1 className="mb-4 text-3xl leading-tight font-bold md:text-4xl lg:text-5xl">
          {recipe.title[locale]}
        </h1>
        <p className="text-foreground/80 mb-6 text-lg">
          {recipe.description[locale]}
        </p>
      </div>
    </motion.div>
  );
}
