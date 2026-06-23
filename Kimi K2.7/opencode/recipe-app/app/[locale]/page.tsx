import { setRequestLocale } from "next-intl/server";
import { getTranslations } from "next-intl/server";
import { routing } from "@/i18n/routing";
import { getAllRecipes } from "@/lib/recipes";
import { createMetadata } from "@/lib/seo";
import { AnimatedSection } from "@/components/animations/AnimatedSection";
import { HeroSection, HeroItem } from "@/components/animations/HeroSection";
import { RecipeGrid } from "@/components/recipes/RecipeGrid";
import { HomeSearch } from "@/components/filters/HomeSearch";
import { Link } from "@/i18n/navigation";
import { Button } from "@/components/ui/Button";

export async function generateMetadata({
  params,
}: {
  params: Promise<{ locale: string }>;
}) {
  const { locale } = await params;
  return createMetadata(locale as "en" | "ru", {
    titleKey: "defaultTitle",
    descriptionKey: "defaultDescription",
    path: "",
  });
}

export function generateStaticParams() {
  return routing.locales.map((locale) => ({ locale }));
}

export default async function HomePage({
  params,
}: {
  params: Promise<{ locale: string }>;
}) {
  const { locale } = await params;
  setRequestLocale(locale);
  const t = await getTranslations({ locale, namespace: "home" });
  const nav = await getTranslations({ locale, namespace: "nav" });
  const cuisineT = await getTranslations({
    locale,
    namespace: "cuisines.labels",
  });
  const allRecipes = getAllRecipes();
  const featured = allRecipes
    .filter((recipe) => recipe.rating >= 4.7)
    .slice(0, 3);
  const cuisineSlugs = [
    "italian",
    "japanese",
    "mexican",
    "indian",
    "mediterranean",
    "american",
  ];

  return (
    <>
      <HeroSection>
        <HeroItem className="mb-6">
          <h1 className="text-4xl leading-tight font-extrabold md:text-6xl">
            {t("heroTitle")}
          </h1>
        </HeroItem>
        <HeroItem className="mb-8">
          <p className="text-foreground/80 text-lg md:text-xl">
            {t("heroSubtitle")}
          </p>
        </HeroItem>
        <HeroItem className="mb-8 flex justify-center">
          <HomeSearch />
        </HeroItem>
        <HeroItem>
          <Link href="/recipes">
            <Button size="lg">{t("browseAll")}</Button>
          </Link>
        </HeroItem>
      </HeroSection>

      <AnimatedSection className="mx-auto max-w-7xl px-4 py-16 md:px-6">
        <div className="mb-8 flex items-end justify-between">
          <h2 className="text-2xl font-bold md:text-3xl">
            {t("featuredTitle")}
          </h2>
          <Link
            href="/recipes"
            className="text-primary text-sm font-medium hover:underline"
          >
            {nav("recipes")} →
          </Link>
        </div>
        <RecipeGrid recipes={featured} locale={locale as "en" | "ru"} />
      </AnimatedSection>

      <AnimatedSection className="bg-muted px-4 py-16 md:px-6">
        <div className="mx-auto max-w-7xl">
          <h2 className="mb-8 text-2xl font-bold md:text-3xl">
            {t("popularCuisines")}
          </h2>
          <div className="grid grid-cols-2 gap-4 sm:grid-cols-3 lg:grid-cols-6">
            {cuisineSlugs.map((cuisine) => (
              <Link
                key={cuisine}
                href={`/cuisines/${cuisine}`}
                className="group bg-card hover:bg-primary hover:text-primary-foreground flex items-center justify-center rounded-lg p-6 text-center font-medium capitalize shadow-sm transition-all hover:shadow-md"
              >
                {cuisineT(cuisine)}
              </Link>
            ))}
          </div>
        </div>
      </AnimatedSection>
    </>
  );
}
