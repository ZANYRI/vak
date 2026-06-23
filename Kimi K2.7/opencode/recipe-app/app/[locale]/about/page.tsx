import { setRequestLocale } from "next-intl/server";
import { getTranslations } from "next-intl/server";
import { routing } from "@/i18n/routing";
import { createMetadata } from "@/lib/seo";
import { AnimatedSection } from "@/components/animations/AnimatedSection";

export async function generateMetadata({
  params,
}: {
  params: Promise<{ locale: string }>;
}) {
  const { locale } = await params;
  return createMetadata(locale as "en" | "ru", {
    titleKey: ["about", "title"],
    descriptionKey: ["about", "intro"],
    path: "/about",
  });
}

export function generateStaticParams() {
  return routing.locales.map((locale) => ({ locale }));
}

export default async function AboutPage({
  params,
}: {
  params: Promise<{ locale: string }>;
}) {
  const { locale } = await params;
  setRequestLocale(locale);
  const t = await getTranslations({ locale, namespace: "about" });

  return (
    <div className="mx-auto max-w-3xl px-4 py-10 md:px-6">
      <AnimatedSection>
        <h1 className="mb-6 text-3xl font-bold md:text-4xl">{t("title")}</h1>
        <p className="text-foreground/80 mb-8 text-lg">{t("intro")}</p>
      </AnimatedSection>

      <AnimatedSection className="space-y-8">
        <section>
          <h2 className="mb-2 text-xl font-semibold">{t("philosophyTitle")}</h2>
          <p className="text-foreground/80">{t("philosophy")}</p>
        </section>
        <section>
          <h2 className="mb-2 text-xl font-semibold">{t("dataTitle")}</h2>
          <p className="text-foreground/80">{t("data")}</p>
        </section>
        <section>
          <h2 className="mb-2 text-xl font-semibold">
            {t("accessibilityTitle")}
          </h2>
          <p className="text-foreground/80">{t("accessibility")}</p>
        </section>
      </AnimatedSection>
    </div>
  );
}
