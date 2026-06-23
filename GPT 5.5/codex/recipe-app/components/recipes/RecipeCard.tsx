"use client";

import Link from "next/link";
import type { Recipe } from "@/data/recipes";
import type { Dictionary } from "@/i18n/dictionary";
import type { Locale } from "@/i18n/config";
import { ImageWithFallback } from "@/components/ui/ImageWithFallback";
import { Badge } from "@/components/ui/Badge";
import { FavoriteButton } from "./FavoriteButton";

export function RecipeCard({ recipe, locale, dictionary }: { recipe: Recipe; locale: Locale; dictionary: Dictionary }) {
  const t = dictionary.recipe;
  const time = `${recipe.totalTimeMinutes} ${t.min}`;
  return <article className="recipe-card">
    <Link className="recipe-image" href={`/${locale}/recipes/${recipe.slug}`} aria-label={recipe.title[locale]}>
      <ImageWithFallback src={recipe.image} alt={recipe.imageAlt[locale]} />
      <span className="rating">★ {recipe.rating.toFixed(1)}</span>
    </Link>
    <div className="recipe-card-content">
      <div className="card-topline"><Badge>{dictionary.cuisines[recipe.cuisine]}</Badge><FavoriteButton slug={recipe.slug} compact addLabel={t.addFavorite} removeLabel={t.removeFavorite} /></div>
      <h3><Link href={`/${locale}/recipes/${recipe.slug}`}>{recipe.title[locale]}</Link></h3>
      <p>{recipe.description[locale]}</p>
      <div className="card-meta"><span>◷ {time}</span><span>◌ {dictionary.difficulty[recipe.difficulty]}</span></div>
    </div>
  </article>;
}
