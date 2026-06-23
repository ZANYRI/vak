"use client";

import { useTranslations } from "next-intl";
import type { Recipe } from "@/lib/validation";

type NutritionCardProps = {
  recipe: Recipe;
};

export function NutritionCard({ recipe }: NutritionCardProps) {
  const t = useTranslations("recipe");

  const items = [
    { label: t("calories"), value: recipe.nutrition.calories, unit: "kcal" },
    { label: t("protein"), value: recipe.nutrition.protein, unit: "g" },
    { label: t("fat"), value: recipe.nutrition.fat, unit: "g" },
    { label: t("carbs"), value: recipe.nutrition.carbs, unit: "g" },
  ];

  return (
    <section className="border-border bg-card rounded-xl border p-6">
      <h2 className="mb-4 text-xl font-bold">{t("nutrition")}</h2>
      <div className="grid grid-cols-2 gap-4 sm:grid-cols-4">
        {items.map((item) => (
          <div key={item.label} className="text-center">
            <p className="text-primary text-2xl font-bold">
              {item.value}
              <span className="text-foreground/60 text-sm font-medium">
                {" "}
                {item.unit}
              </span>
            </p>
            <p className="text-foreground/70 text-sm">{item.label}</p>
          </div>
        ))}
      </div>
    </section>
  );
}
