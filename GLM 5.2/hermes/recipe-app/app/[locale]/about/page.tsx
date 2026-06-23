import { setRequestLocale, getTranslations } from "next-intl/server";
import { useTranslations } from "next-intl";
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
    title: "About",
    description: t("aboutDescription"),
    path: "about"
  });
}

export default async function AboutPage({
  params
}: {
  params: Promise<{ locale: string }>;
}) {
  const { locale } = await params;
  setRequestLocale(locale);
  return <AboutContent />;
}

function AboutContent() {
  const t = useTranslations("about");

  const sections: Array<{
    title: string;
    body: string;
  }> = [
    { title: t("philosophyTitle"), body: t("philosophy") },
    { title: t("dataSourceTitle"), body: t("dataSource") },
    { title: t("imagesTitle"), body: t("images") },
    { title: t("accessibilityTitle"), body: t("accessibility") }
  ];

  return (
    <div className="container-page py-12">
      <AnimatedSection className="mx-auto max-w-3xl">
        <h1 className="font-serif text-3xl font-semibold text-foreground sm:text-4xl">
          {t("title")}
        </h1>
        <p className="mt-4 text-lg text-muted">{t("lead")}</p>

        <div className="mt-10 space-y-8">
          {sections.map((s) => (
            <section key={s.title}>
              <h2 className="font-serif text-xl font-semibold text-foreground">
                {s.title}
              </h2>
              <p className="mt-2 leading-relaxed text-foreground/80">{s.body}</p>
            </section>
          ))}
        </div>
      </AnimatedSection>
    </div>
  );
}
