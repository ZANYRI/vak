import { useLocale, useTranslations } from "next-intl";
import type { Recipe } from "@/lib/types";

export function IngredientList({ recipe }: { recipe: Recipe }) {
  const locale = useLocale() as "en" | "ru";
  const t = useTranslations("recipe");
  const items = recipe.ingredients[locale];

  return (
    <section aria-labelledby="ingredients-heading" className="rounded-card border border-border bg-card p-5 shadow-soft">
      <h2 id="ingredients-heading" className="font-serif text-xl font-semibold text-foreground">
        {t("ingredients")}
      </h2>
      <ul className="mt-4 space-y-2.5">
        {items.map((item, i) => (
          <li key={i} className="flex items-start gap-3 text-sm text-foreground">
            <span
              aria-hidden="true"
              className="mt-1 flex h-4 w-4 flex-none items-center justify-center rounded-full border border-primary text-[10px] text-primary"
            >
              ●
            </span>
            <span>{item}</span>
          </li>
        ))}
      </ul>
    </section>
  );
}
