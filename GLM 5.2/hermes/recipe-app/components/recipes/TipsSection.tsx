import { useLocale, useTranslations } from "next-intl";
import type { Recipe } from "@/lib/types";

export function TipsSection({ recipe }: { recipe: Recipe }) {
  const locale = useLocale() as "en" | "ru";
  const t = useTranslations("recipe");
  const tips = recipe.tips[locale];
  if (tips.length === 0) return null;

  return (
    <section
      aria-labelledby="tips-heading"
      className="rounded-card border border-accent/30 bg-accent/5 p-5"
    >
      <h2 id="tips-heading" className="font-serif text-xl font-semibold text-foreground">
        {t("tips")}
      </h2>
      <ul className="mt-4 space-y-3">
        {tips.map((tip, i) => (
          <li key={i} className="flex items-start gap-3 text-sm text-foreground">
            <span aria-hidden="true" className="text-accent">
              💡
            </span>
            <span>{tip}</span>
          </li>
        ))}
      </ul>
    </section>
  );
}
