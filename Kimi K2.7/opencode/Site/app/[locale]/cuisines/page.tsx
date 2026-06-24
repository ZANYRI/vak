import { setRequestLocale } from "next-intl/server";
import { getTranslations } from "next-intl/server";
import { routing } from "@/i18n/routing";
import { getAllRecipes, CUISINES } from "@/lib/recipes";
import { createMetadata } from "@/lib/seo";
import { Link } from "@/i18n/navigation";
import { AnimatedSection } from "@/components/animations/AnimatedSection";
import { ImageWithFallback } from "@/components/recipes/ImageWithFallback";

export async function generateMetadata({
  params,
}: {
  params: Promise<{ locale: string }>;
}) {
  const { locale } = await params;
  return createMetadata(locale as "en" | "ru", {
    titleKey: ["cuisines", "title"],
    descriptionKey: ["cuisines", "description"],
    path: "/cuisines",
  });
}

export function generateStaticParams() {
  return routing.locales.map((locale) => ({ locale }));
}

export default async function CuisinesPage({
  params,
}: {
  params: Promise<{ locale: string }>;
}) {
  const { locale } = await params;
  setRequestLocale(locale);
  const t = await getTranslations({ locale, namespace: "cuisines" });
  const recipes = getAllRecipes();

  const cuisineWithImage = CUISINES.map((cuisine) => {
    const first = recipes.find((r) => r.cuisine === cuisine);
    return { cuisine, image: first?.image ?? "" };
  });

  return (
    <div className="mx-auto max-w-7xl px-4 py-10 md:px-6">
      <AnimatedSection>
        <h1 className="mb-2 text-3xl font-bold md:text-4xl">{t("title")}</h1>
        <p className="text-foreground/70 mb-8">{t("description")}</p>
      </AnimatedSection>
      <AnimatedSection>
        <div className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
          {cuisineWithImage.map(({ cuisine, image }) => (
            <Link
              key={cuisine}
              href={`/cuisines/${cuisine}`}
              className="group border-border bg-card hover:border-primary overflow-hidden rounded-lg border shadow-sm transition-all hover:shadow-md"
            >
              <div className="relative aspect-[16/9]">
                <ImageWithFallback
                  src={image}
                  alt={t(`labels.${cuisine}`)}
                  fill
                  className="group-hover:scale-105"
                />
              </div>
              <div className="p-4">
                <h2 className="text-lg font-semibold capitalize">
                  {t(`labels.${cuisine}`)}
                </h2>
              </div>
            </Link>
          ))}
        </div>
      </AnimatedSection>
    </div>
  );
}
