import { setRequestLocale, getLocale } from "next-intl/server";
import { notFound } from "next/navigation";
import { getRecipeBySlug, getAllRecipes } from "@/lib/recipes";
import { RecipeHero } from "@/components/recipes/RecipeHero";
import { RecipeMeta } from "@/components/recipes/RecipeMeta";
import { IngredientList } from "@/components/recipes/IngredientList";
import { InstructionSteps } from "@/components/recipes/InstructionSteps";
import { NutritionCard } from "@/components/recipes/NutritionCard";
import { TipsSection } from "@/components/recipes/TipsSection";
import { RelatedRecipes } from "@/components/recipes/RelatedRecipes";
import { ShareButton } from "@/components/ui/ShareButton";
import { PrintButton } from "@/components/ui/PrintButton";
import { AnimatedSection } from "@/components/animations/AnimatedSection";
import { buildRecipeMetadata, buildRecipeJsonLd } from "@/lib/seo";
import type { Locale } from "@/lib/types";

export function generateStaticParams() {
  return getAllRecipes().map((r) => ({ slug: r.slug }));
}

export async function generateMetadata({
  params
}: {
  params: Promise<{ locale: string; slug: string }>;
}) {
  const { locale, slug } = await params;
  const recipe = getRecipeBySlug(slug);
  if (!recipe) return {};
  return buildRecipeMetadata(recipe, locale as Locale);
}

export default async function RecipeDetailPage({
  params
}: {
  params: Promise<{ locale: string; slug: string }>;
}) {
  const { locale, slug } = await params;
  setRequestLocale(locale);
  const recipe = getRecipeBySlug(slug);
  if (!recipe) notFound();

  const currentLocale = (await getLocale()) as Locale;

  const jsonLd = buildRecipeJsonLd(recipe, currentLocale);

  return (
    <article className="container-page py-8">
      <script
        type="application/ld+json"
        dangerouslySetInnerHTML={{ __html: JSON.stringify(jsonLd) }}
      />

      <RecipeHero recipe={recipe} />

      <div className="mt-4 flex flex-wrap items-center gap-3">
        <ShareButton />
        <PrintButton />
      </div>

      <AnimatedSection className="mt-6">
        <RecipeMeta recipe={recipe} />
      </AnimatedSection>

      <div className="mt-6 grid gap-6 lg:grid-cols-[1fr_360px]">
        <div className="space-y-6">
          <InstructionSteps recipe={recipe} />
          <TipsSection recipe={recipe} />
        </div>
        <aside className="space-y-6">
          <IngredientList recipe={recipe} />
          <NutritionCard recipe={recipe} />
        </aside>
      </div>

      <RelatedRecipes recipe={recipe} />
    </article>
  );
}
