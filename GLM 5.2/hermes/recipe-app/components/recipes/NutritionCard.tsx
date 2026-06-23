import { useTranslations } from "next-intl";
import type { Nutrition, Recipe } from "@/lib/types";

export function NutritionCard({ recipe }: { recipe: Recipe }) {
  const t = useTranslations("recipe");
  const n: Nutrition = recipe.nutrition;

  const rows: Array<{ label: string; value: string; accent?: boolean }> = [
    { label: t("calories"), value: `${n.calories}`, accent: true },
    { label: t("protein"), value: `${n.protein} g` },
    { label: t("fat"), value: `${n.fat} g` },
    { label: t("carbs"), value: `${n.carbs} g` }
  ];

  return (
    <section
      aria-labelledby="nutrition-heading"
      className="rounded-card border border-border bg-card p-5 shadow-soft"
    >
      <h2 id="nutrition-heading" className="font-serif text-xl font-semibold text-foreground">
        {t("nutrition")}
      </h2>
      <dl className="mt-4 grid grid-cols-2 gap-3">
        {rows.map((r) => (
          <div
            key={r.label}
            className={`rounded-xl p-3 text-center ${r.accent ? "bg-primary/10" : "bg-surface-alt"}`}
          >
            <dt className="text-xs uppercase tracking-wide text-muted">{r.label}</dt>
            <dd className={`mt-0.5 text-lg font-semibold ${r.accent ? "text-primary" : "text-foreground"}`}>
              {r.value}
            </dd>
          </div>
        ))}
      </dl>
      <p className="mt-3 text-xs text-muted">
        {t("calories")}: per serving · approximate values
      </p>
    </section>
  );
}
