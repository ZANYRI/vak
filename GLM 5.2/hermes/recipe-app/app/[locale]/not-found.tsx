import { useTranslations } from "next-intl";
import { Link } from "@/i18n/navigation";
import { Button } from "@/components/ui/Button";
import { AnimatedSection } from "@/components/animations/AnimatedSection";

export default function NotFound() {
  const t = useTranslations("notFound");

  return (
    <div className="container-page flex min-h-[50vh] items-center justify-center py-16">
      <AnimatedSection className="text-center">
        <p className="font-serif text-7xl font-semibold text-primary">404</p>
        <h1 className="mt-4 font-serif text-2xl font-semibold text-foreground">
          {t("title")}
        </h1>
        <p className="mt-2 text-muted">{t("description")}</p>
        <div className="mt-6">
          <Link href="/">
            <Button>{t("goHome")}</Button>
          </Link>
        </div>
      </AnimatedSection>
    </div>
  );
}
