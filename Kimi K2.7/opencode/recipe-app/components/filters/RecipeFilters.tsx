"use client";

import { usePathname, useRouter, useSearchParams } from "next/navigation";
import { useTranslations } from "next-intl";
import { Select } from "@/components/ui/Select";
import { Button } from "@/components/ui/Button";

const TIME_OPTIONS = [15, 30, 45, 60, 90, 120];
const SORT_OPTIONS = ["newest", "fastest", "easiest", "highestRated"];

export function RecipeFilters() {
  const router = useRouter();
  const pathname = usePathname();
  const searchParams = useSearchParams();

  const tFilters = useTranslations("filters");
  const tCategories = useTranslations("categories.labels");
  const tCuisines = useTranslations("cuisines.labels");
  const tDifficulties = useTranslations("filters.difficulties");
  const tDiets = useTranslations("filters.diets");
  const tMealTypes = useTranslations("filters.mealTypes");
  const tRecipes = useTranslations("recipes");

  const cuisines = [
    "italian",
    "french",
    "georgian",
    "japanese",
    "mexican",
    "indian",
    "mediterranean",
    "ukrainian",
    "latvian",
    "american",
  ];
  const categories = [
    "breakfast",
    "lunch",
    "dinner",
    "desserts",
    "soups",
    "salads",
    "vegetarian",
    "quick-meals",
    "baking",
    "drinks",
  ];
  const difficulties = ["easy", "medium", "hard"];
  const diets = ["vegetarian", "vegan", "gluten-free", "dairy-free"];
  const mealTypes = [
    "breakfast",
    "brunch",
    "lunch",
    "dinner",
    "snack",
    "dessert",
    "drink",
  ];

  const updateParam = (key: string, value: string) => {
    const params = new URLSearchParams(searchParams.toString());
    if (value && value !== "all") {
      params.set(key, value);
    } else {
      params.delete(key);
    }
    router.push(`${pathname}?${params.toString()}`, { scroll: false });
  };

  const clearFilters = () => {
    const params = new URLSearchParams();
    const search = searchParams.get("search");
    if (search) params.set("search", search);
    router.push(`${pathname}?${params.toString()}`, { scroll: false });
  };

  const hasActive = [
    "cuisine",
    "category",
    "difficulty",
    "diet",
    "mealType",
    "time",
    "sort",
  ].some((key) => searchParams.has(key));

  return (
    <div className="border-border bg-card space-y-4 rounded-xl border p-4 md:p-6">
      <div className="flex items-center justify-between">
        <h2 className="font-bold">{tRecipes("filters")}</h2>
        {hasActive && (
          <Button variant="ghost" size="sm" onClick={clearFilters}>
            {tFilters("clear")}
          </Button>
        )}
      </div>
      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
        <Select
          label={tFilters("cuisine")}
          value={searchParams.get("cuisine") ?? "all"}
          onChange={(e) => updateParam("cuisine", e.target.value)}
        >
          <option value="all">{tFilters("all")}</option>
          {cuisines.map((c) => (
            <option key={c} value={c}>
              {tCuisines(c)}
            </option>
          ))}
        </Select>

        <Select
          label={tFilters("category")}
          value={searchParams.get("category") ?? "all"}
          onChange={(e) => updateParam("category", e.target.value)}
        >
          <option value="all">{tFilters("all")}</option>
          {categories.map((c) => (
            <option key={c} value={c}>
              {tCategories(c)}
            </option>
          ))}
        </Select>

        <Select
          label={tFilters("difficulty")}
          value={searchParams.get("difficulty") ?? "all"}
          onChange={(e) => updateParam("difficulty", e.target.value)}
        >
          <option value="all">{tFilters("all")}</option>
          {difficulties.map((d) => (
            <option key={d} value={d}>
              {tDifficulties(d)}
            </option>
          ))}
        </Select>

        <Select
          label={tFilters("diet")}
          value={searchParams.get("diet") ?? "all"}
          onChange={(e) => updateParam("diet", e.target.value)}
        >
          <option value="all">{tFilters("all")}</option>
          {diets.map((d) => (
            <option key={d} value={d}>
              {tDiets(d)}
            </option>
          ))}
        </Select>

        <Select
          label={tFilters("mealType")}
          value={searchParams.get("mealType") ?? "all"}
          onChange={(e) => updateParam("mealType", e.target.value)}
        >
          <option value="all">{tFilters("all")}</option>
          {mealTypes.map((m) => (
            <option key={m} value={m}>
              {tMealTypes(m)}
            </option>
          ))}
        </Select>

        <Select
          label={tFilters("time")}
          value={searchParams.get("time") ?? "all"}
          onChange={(e) => updateParam("time", e.target.value)}
        >
          <option value="all">{tFilters("all")}</option>
          {TIME_OPTIONS.map((time) => (
            <option key={time} value={time}>
              {time} min
            </option>
          ))}
        </Select>

        <Select
          label={tRecipes("sortLabel")}
          value={searchParams.get("sort") ?? "newest"}
          onChange={(e) => updateParam("sort", e.target.value)}
        >
          {SORT_OPTIONS.map((sort) => (
            <option key={sort} value={sort}>
              {tRecipes(`sort.${sort}`)}
            </option>
          ))}
        </Select>
      </div>
    </div>
  );
}
