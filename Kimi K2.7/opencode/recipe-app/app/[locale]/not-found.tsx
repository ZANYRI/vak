"use client";

import { useTranslations } from "next-intl";
import { Link } from "@/i18n/navigation";
import { Button } from "@/components/ui/Button";

export default function NotFoundPage() {
  const t = useTranslations("errors");

  return (
    <div className="flex flex-1 flex-col items-center justify-center px-4 py-20 text-center">
      <h1 className="mb-4 text-4xl font-bold">404</h1>
      <p className="text-foreground/70 mb-8 text-lg">{t("notFound")}</p>
      <Link href="/">
        <Button>{t("notFound")}</Button>
      </Link>
    </div>
  );
}
