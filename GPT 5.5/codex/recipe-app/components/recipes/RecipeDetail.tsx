import Link from "next/link";
import type { Recipe } from "@/data/recipes";
import type { Locale } from "@/i18n/config";
import type { Dictionary } from "@/i18n/dictionary";
import { ImageWithFallback } from "@/components/ui/ImageWithFallback";
import { Badge } from "@/components/ui/Badge";
import { FavoriteButton } from "./FavoriteButton";
import { ShareButton } from "./ShareButton";
import { PrintButton } from "./PrintButton";
import { relatedRecipes } from "@/lib/recipes";
import { RecipeHero } from "./RecipeHero";
import { RecipeMeta } from "./RecipeMeta";
import { IngredientList } from "./IngredientList";
import { InstructionSteps } from "./InstructionSteps";
import { NutritionCard } from "./NutritionCard";
import { RelatedRecipes } from "./RelatedRecipes";

export function RecipeDetail({ recipe, locale, dictionary }: { recipe: Recipe; locale: Locale; dictionary: Dictionary }) {
  const t = dictionary.recipe; const cuisine = dictionary.cuisines[recipe.cuisine]; const category = dictionary.categories[recipe.category];
  const jsonLd = { "@context": "https://schema.org", "@type": "Recipe", name: recipe.title[locale], description: recipe.description[locale], image: [recipe.image], author: { "@type": "Organization", name: "Misen" }, datePublished: recipe.publishedAt, prepTime: `PT${recipe.prepTimeMinutes}M`, cookTime: `PT${recipe.cookTimeMinutes}M`, totalTime: `PT${recipe.totalTimeMinutes}M`, recipeYield: `${recipe.servings}`, recipeCuisine: cuisine, recipeCategory: category, recipeIngredient: recipe.ingredients[locale], recipeInstructions: recipe.steps[locale].map((text, index) => ({ "@type": "HowToStep", position: index + 1, text })), nutrition: { "@type": "NutritionInformation", calories: `${recipe.nutrition.calories} calories`, proteinContent: `${recipe.nutrition.protein} g`, fatContent: `${recipe.nutrition.fat} g`, carbohydrateContent: `${recipe.nutrition.carbs} g` } };
  return <article>
    <script type="application/ld+json" dangerouslySetInnerHTML={{ __html: JSON.stringify(jsonLd).replace(/</g, "\\u003c") }} />
    <RecipeHero><div className="shell recipe-hero-grid"><div><Link className="back-link" href={`/${locale}/recipes`}>← {dictionary.common.back}</Link><div className="badge-row"><Badge>{cuisine}</Badge><Badge>{dictionary.difficulty[recipe.difficulty]}</Badge></div><h1>{recipe.title[locale]}</h1><p>{recipe.description[locale]}</p><div className="recipe-actions"><FavoriteButton slug={recipe.slug} addLabel={t.addFavorite} removeLabel={t.removeFavorite} /><ShareButton label={t.share} title={recipe.title[locale]} /><PrintButton label={t.print} /></div></div><div className="detail-image"><ImageWithFallback src={recipe.image} alt={recipe.imageAlt[locale]} sizes="(max-width: 850px) 100vw, 50vw" /></div></div></RecipeHero>
    <section className="shell recipe-layout"><aside className="recipe-facts"><dl><RecipeMeta label={t.prep} value={`${recipe.prepTimeMinutes} ${t.min}`} /><RecipeMeta label={t.cook} value={`${recipe.cookTimeMinutes} ${t.min}`} /><RecipeMeta label={t.total} value={`${recipe.totalTimeMinutes} ${t.min}`} /><RecipeMeta label={t.servings} value={recipe.servings} /><RecipeMeta label={t.cuisine} value={cuisine} /><RecipeMeta label={t.meal} value={recipe.mealType.map((item) => dictionary.categories[item as "breakfast" | "lunch" | "dinner"] ?? item).join(", ")} /></dl></aside><div className="recipe-body"><section><h2>{t.ingredients}</h2><IngredientList ingredients={recipe.ingredients[locale]} /></section><section><h2>{t.instructions}</h2><InstructionSteps steps={recipe.steps[locale]} /></section><NutritionCard title={t.nutrition} labels={[t.calories, t.protein, t.fat, t.carbs]} values={[recipe.nutrition.calories, `${recipe.nutrition.protein}g`, `${recipe.nutrition.fat}g`, `${recipe.nutrition.carbs}g`]} /><section className="tips"><h2>{t.tips}</h2><ul>{recipe.tips[locale].map((tip) => <li key={tip}>{tip}</li>)}</ul></section></div></section>
    <RelatedRecipes title={t.related} recipes={relatedRecipes(recipe)} locale={locale} dictionary={dictionary} />
  </article>;
}
