import { notFound } from "next/navigation";
import { setRequestLocale } from "next-intl/server";
import { getTranslations } from "next-intl/server";
import { routing } from "@/i18n/routing";
import { getAllRecipes, getRecipeBySlug } from "@/lib/recipes";
import { getSiteUrl, buildJsonLdRecipe } from "@/lib/seo";
import { AnimatedSection } from "@/components/animations/AnimatedSection";
import { RecipeHero } from "@/components/recipes/RecipeHero";
import { RecipeMeta } from "@/components/recipes/RecipeMeta";
import { IngredientList } from "@/components/recipes/IngredientList";
import { InstructionSteps } from "@/components/recipes/InstructionSteps";
import { NutritionCard } from "@/components/recipes/NutritionCard";
import { RelatedRecipes } from "@/components/recipes/RelatedRecipes";
import { FavoriteButton } from "@/components/recipes/FavoriteButton";
import { ShareButton } from "@/components/recipes/ShareButton";
import { PrintButton } from "@/components/recipes/PrintButton";
import { Badge } from "@/components/ui/Badge";
import type { Metadata } from "next";

export async function generateMetadata({
  params,
}: {
  params: Promise<{ locale: string; slug: string }>;
}): Promise<Metadata> {
  const { locale, slug } = await params;
  const recipe = getRecipeBySlug(slug);
  if (!recipe) return {};

  const title = recipe.title[locale as "en" | "ru"];
  const description = recipe.description[locale as "en" | "ru"];
  const siteUrl = getSiteUrl();
  const path = `/recipes/${recipe.slug}`;
  const canonical = `${siteUrl}/${locale}${path}`;

  return {
    title,
    description,
    openGraph: {
      title,
      description,
      url: canonical,
      type: "article",
      locale,
      images: [{ url: recipe.image }],
    },
    twitter: {
      card: "summary_large_image",
      title,
      description,
      images: [recipe.image],
    },
    alternates: {
      canonical,
      languages: {
        en: `${siteUrl}/en${path}`,
        ru: `${siteUrl}/ru${path}`,
      },
    },
  };
}

export function generateStaticParams() {
  const recipes = getAllRecipes();
  const params: Array<{ locale: string; slug: string }> = [];
  for (const locale of routing.locales) {
    for (const recipe of recipes) {
      params.push({ locale, slug: recipe.slug });
    }
  }
  return params;
}

export default async function RecipeDetailPage({
  params,
}: {
  params: Promise<{ locale: string; slug: string }>;
}) {
  const { locale, slug } = await params;
  setRequestLocale(locale);
  const typedLocale = locale as "en" | "ru";
  const recipe = getRecipeBySlug(slug);
  if (!recipe) {
    notFound();
  }

  const t = await getTranslations({ locale, namespace: "recipe" });
  const siteUrl = getSiteUrl();
  const shareUrl = `${siteUrl}/${locale}/recipes/${recipe.slug}`;
  const jsonLd = buildJsonLdRecipe({
    title: recipe.title[typedLocale],
    description: recipe.description[typedLocale],
    image: recipe.image,
    cuisine: recipe.cuisine,
    category: recipe.category,
    ingredients: recipe.ingredients[typedLocale],
    instructions: recipe.steps[typedLocale],
    prepTimeMinutes: recipe.prepTimeMinutes,
    cookTimeMinutes: recipe.cookTimeMinutes,
    totalTimeMinutes: recipe.totalTimeMinutes,
    servings: recipe.servings,
    rating: recipe.rating,
    nutrition: recipe.nutrition,
  });

  return (
    <>
      <script
        type="application/ld+json"
        dangerouslySetInnerHTML={{ __html: JSON.stringify(jsonLd) }}
      />
      <article className="mx-auto max-w-7xl px-4 py-10 md:px-6 print:px-0">
        <AnimatedSection>
          <RecipeHero recipe={recipe} locale={typedLocale} />
        </AnimatedSection>

        <AnimatedSection className="mt-8">
          <div className="mb-6 flex flex-wrap items-center gap-2">
            {recipe.diet.map((d) => (
              <Badge key={d} variant="secondary">
                {d}
              </Badge>
            ))}
            {recipe.mealType.map((m) => (
              <Badge key={m}>{m}</Badge>
            ))}
          </div>
          <RecipeMeta recipe={recipe} />
        </AnimatedSection>

        <AnimatedSection className="mt-8 flex flex-wrap gap-3">
          <FavoriteButton id={recipe.id} />
          <ShareButton title={recipe.title[typedLocale]} url={shareUrl} />
          <PrintButton label={t("print")} />
        </AnimatedSection>

        <div className="mt-10 grid gap-8 lg:grid-cols-3">
          <div className="lg:col-span-1">
            <IngredientList recipe={recipe} locale={typedLocale} />
          </div>
          <div className="space-y-8 lg:col-span-2">
            <InstructionSteps recipe={recipe} locale={typedLocale} />
            <NutritionCard recipe={recipe} />
            {recipe.tips[typedLocale].length > 0 && (
              <section>
                <h2 className="mb-3 text-xl font-bold">{t("tips")}</h2>
                <ul className="list-disc space-y-1 pl-5">
                  {recipe.tips[typedLocale].map((tip, i) => (
                    <li key={i} className="text-foreground/90">
                      {tip}
                    </li>
                  ))}
                </ul>
              </section>
            )}
          </div>
        </div>

        <AnimatedSection className="mt-16">
          <RelatedRecipes recipe={recipe} locale={typedLocale} />
        </AnimatedSection>
      </article>
    </>
  );
}
