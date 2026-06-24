"use client";

import { useTranslations } from "next-intl";
import { useRouter, useSearchParams } from "next/navigation";
import { useTransition } from "react";
import { Select } from "@/components/ui/Select";
import {
  CATEGORIES,
  CUISINES,
  DIETS,
  DIFFICULTIES,
  MEAL_TYPES
} from "@/lib/taxonomy";
import { useLocale } from "next-intl";
import { Button } from "@/components/ui/Button";

const TIME_OPTIONS = [
  { value: "", labelKey: "anyTime" },
  { value: "15", label: "≤ 15 min" },
  { value: "30", label: "≤ 30 min" },
  { value: "45", label: "≤ 45 min" },
  { value: "60", label: "≤ 60 min" },
  { value: "90", label: "≤ 90 min" }
];

const SORT_OPTIONS = [
  { value: "newest", key: "newest" },
  { value: "fastest", key: "fastest" },
  { value: "easiest", key: "easiest" },
  { value: "rating", key: "rating" }
];

export interface RecipeFiltersProps {
  basePath: string;
}

export function RecipeFilters({ basePath }: RecipeFiltersProps) {
  const t = useTranslations();
  const locale = useLocale() as "en" | "ru";
  const router = useRouter();
  const searchParams = useSearchParams();
  const [isPending, startTransition] = useTransition();

  const updateParam = (key: string, value: string) => {
    const params = new URLSearchParams(searchParams.toString());
    if (value && value !== "page") {
      params.set(key, value);
    } else {
      params.delete(key);
    }
    params.delete("page"); // reset to first page on filter change
    const qs = params.toString();
    router.replace(qs ? `${basePath}?${qs}` : basePath, { scroll: false });
  };

  const clearAll = () => {
    startTransition(() => router.replace(basePath, { scroll: false }));
  };

  const hasFilters =
    searchParams.get("cuisine") ||
    searchParams.get("category") ||
    searchParams.get("difficulty") ||
    searchParams.get("diet") ||
    searchParams.get("mealType") ||
    searchParams.get("maxTime");

  const toOptions = (
    entries: Array<{ slug: string; name: { en: string; ru: string } }>,
    anyKey: string
  ) => [
    { value: "", label: t(`filters.${anyKey}`) },
    ...entries.map((e) => ({ value: e.slug, label: e.name[locale] }))
  ];

  return (
    <div className="space-y-4 rounded-card border border-border bg-card p-4 shadow-soft">
      <div className="flex items-center justify-between">
        <h2 className="font-serif text-lg font-semibold text-foreground">
          {t("filters.title")}
        </h2>
        {hasFilters ? (
          <Button variant="ghost" size="sm" onClick={clearAll} disabled={isPending}>
            {t("filters.clear")}
          </Button>
        ) : null}
      </div>

      <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
        <Select
          label={t("filters.cuisine")}
          value={searchParams.get("cuisine") ?? ""}
          options={toOptions(CUISINES, "anyCuisine")}
          onChange={(e) => updateParam("cuisine", e.target.value)}
        />
        <Select
          label={t("filters.category")}
          value={searchParams.get("category") ?? ""}
          options={toOptions(CATEGORIES, "anyCategory")}
          onChange={(e) => updateParam("category", e.target.value)}
        />
        <Select
          label={t("filters.difficulty")}
          value={searchParams.get("difficulty") ?? ""}
          options={[
            { value: "", label: t("filters.anyDifficulty") },
            ...DIFFICULTIES.map((d) => ({ value: d.slug, label: d.name[locale] }))
          ]}
          onChange={(e) => updateParam("difficulty", e.target.value)}
        />
        <Select
          label={t("filters.diet")}
          value={searchParams.get("diet") ?? ""}
          options={toOptions(DIETS, "anyDiet")}
          onChange={(e) => updateParam("diet", e.target.value)}
        />
        <Select
          label={t("filters.mealType")}
          value={searchParams.get("mealType") ?? ""}
          options={toOptions(MEAL_TYPES, "anyMealType")}
          onChange={(e) => updateParam("mealType", e.target.value)}
        />
        <Select
          label={t("filters.maxTime")}
          value={searchParams.get("maxTime") ?? ""}
          options={TIME_OPTIONS.map((o) => ({
            value: o.value,
            label: o.label ?? t(`filters.${o.labelKey}`)
          }))}
          onChange={(e) => updateParam("maxTime", e.target.value)}
        />
      </div>

      <div className="flex flex-wrap items-center justify-between gap-3 border-t border-border pt-3">
        <Select
          label={t("filters.sortBy")}
          value={searchParams.get("sort") ?? "newest"}
          options={SORT_OPTIONS.map((o) => ({
            value: o.value,
            label: t(`sort.${o.key}`)
          }))}
          onChange={(e) => updateParam("sort", e.target.value)}
          className="max-w-56"
        />
      </div>
    </div>
  );
}
