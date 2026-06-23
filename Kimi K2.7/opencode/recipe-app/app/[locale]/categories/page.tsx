import { setRequestLocale } from "next-intl/server";
import { getTranslations } from "next-intl/server";
import { routing } from "@/i18n/routing";
import { getAllRecipes, CATEGORIES } from "@/lib/recipes";
import { createMetadata } from "@/lib/seo";
import { Link } from "@/i18n/navigation";
import { AnimatedSection } from "@/components/animations/AnimatedSection";

export async function generateMetadata({
  params,
}: {
  params: Promise<{ locale: string }>;
}) {
  const { locale } = await params;
  return createMetadata(locale as "en" | "ru", {
    titleKey: ["categories", "title"],
    descriptionKey: ["categories", "description"],
    path: "/categories",
  });
}

export function generateStaticParams() {
  return routing.locales.map((locale) => ({ locale }));
}

export default async function CategoriesPage({
  params,
}: {
  params: Promise<{ locale: string }>;
}) {
  const { locale } = await params;
  setRequestLocale(locale);
  const t = await getTranslations({ locale, namespace: "categories" });
  const recipes = getAllRecipes();

  const counts = CATEGORIES.reduce<Record<string, number>>((acc, category) => {
    acc[category] = recipes.filter((r) => r.category === category).length;
    return acc;
  }, {});

  return (
    <div className="mx-auto max-w-7xl px-4 py-10 md:px-6">
      <AnimatedSection>
        <h1 className="mb-2 text-3xl font-bold md:text-4xl">{t("title")}</h1>
        <p className="text-foreground/70 mb-8">{t("description")}</p>
      </AnimatedSection>
      <AnimatedSection>
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
          {CATEGORIES.map((category) => (
            <Link
              key={category}
              href={`/categories/${category}`}
              className="border-border bg-card hover:border-primary flex items-center justify-between rounded-lg border p-6 shadow-sm transition-all hover:shadow-md"
            >
              <span className="text-lg font-medium">
                {t(`labels.${category}`)}
              </span>
              <span className="bg-muted rounded-full px-3 py-1 text-sm">
                {counts[category]}
              </span>
            </Link>
          ))}
        </div>
      </AnimatedSection>
    </div>
  );
}
