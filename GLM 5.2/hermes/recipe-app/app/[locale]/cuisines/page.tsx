import { setRequestLocale, getTranslations } from "next-intl/server";
import { useLocale, useTranslations } from "next-intl";
import { Link } from "@/i18n/navigation";
import { CUISINES } from "@/lib/taxonomy";
import { getRecipesCountByCuisine } from "@/lib/recipes";
import { AnimatedSection } from "@/components/animations/AnimatedSection";
import { buildMetadata } from "@/lib/seo";
import type { Locale } from "@/lib/types";

export async function generateMetadata({
  params
}: {
  params: Promise<{ locale: string }>;
}) {
  const { locale } = await params;
  const t = await getTranslations({ locale, namespace: "seo" });
  return buildMetadata({
    locale: locale as Locale,
    title: "Cuisines",
    description: t("cuisinesDescription"),
    path: "cuisines"
  });
}

export default async function CuisinesPage({
  params
}: {
  params: Promise<{ locale: string }>;
}) {
  const { locale } = await params;
  setRequestLocale(locale);
  return <CuisinesContent />;
}

function CuisinesContent() {
  const t = useTranslations();
  const locale = useLocale() as Locale;

  return (
    <div className="container-page py-10">
      <AnimatedSection className="mb-8">
        <h1 className="font-serif text-3xl font-semibold text-foreground sm:text-4xl">
          {t("cuisines.title")}
        </h1>
        <p className="mt-2 text-muted">{t("cuisines.subtitle")}</p>
      </AnimatedSection>

      <ul className="grid gap-5 sm:grid-cols-2 lg:grid-cols-3" role="list">
        {CUISINES.map((c, i) => {
          const count = getRecipesCountByCuisine(c.slug);
          return (
            <AnimatedSection key={c.slug} delay={i * 0.04}>
              <li>
                <Link
                  href={`/cuisines/${c.slug}`}
                  className="group flex items-center justify-between rounded-card border border-border bg-card p-5 shadow-soft transition-all hover:-translate-y-0.5 hover:shadow-lift"
                >
                  <div>
                    <h2 className="font-serif text-xl font-semibold text-foreground">
                      {c.name[locale]}
                    </h2>
                    <p className="mt-1 text-sm text-muted">
                      {t("cuisines.recipesCount", { count })}
                    </p>
                  </div>
                  <span
                    aria-hidden="true"
                    className="text-muted transition-transform group-hover:translate-x-1"
                  >
                    →
                  </span>
                </Link>
              </li>
            </AnimatedSection>
          );
        })}
      </ul>
    </div>
  );
}
