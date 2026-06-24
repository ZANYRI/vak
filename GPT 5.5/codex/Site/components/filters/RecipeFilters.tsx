"use client";

import { usePathname, useRouter, useSearchParams } from "next/navigation";
import { categories, cuisines } from "@/data/recipes";
import type { Dictionary } from "@/i18n/dictionary";
import { SearchBar } from "@/components/recipes/SearchBar";

const selectKeys = ["cuisine", "category", "difficulty", "diet", "meal", "time", "sort"] as const;

export function RecipeFilters({ dictionary }: { dictionary: Dictionary }) {
  const router = useRouter(); const pathname = usePathname(); const searchParams = useSearchParams();
  function onSubmit(formData: FormData) {
    const params = new URLSearchParams();
    for (const key of ["search", ...selectKeys]) { const value = formData.get(key); if (typeof value === "string" && value) params.set(key, value); }
    const cleanPath = pathname.replace(/\/page\/\d+$/, "");
    router.replace(`${cleanPath}${params.size ? `?${params}` : ""}`);
  }
  const current = (key: string) => searchParams.get(key) ?? "";
  const label = (key: string) => dictionary.recipes[key as keyof typeof dictionary.recipes] as string;
  return <form className="filters" action={onSubmit}>
    <SearchBar label={dictionary.recipes.search} defaultValue={current("search")} placeholder={dictionary.recipes.search} />
    <details className="filter-details"><summary>{dictionary.recipes.filters}<span>+</span></summary><div className="filter-panel">
      <label>{dictionary.recipe.cuisine}<select name="cuisine" defaultValue={current("cuisine")}><option value="">{dictionary.recipes.any}</option>{cuisines.map((item) => <option key={item} value={item}>{dictionary.cuisines[item]}</option>)}</select></label>
      <label>{dictionary.categories.title}<select name="category" defaultValue={current("category")}><option value="">{dictionary.recipes.any}</option>{categories.map((item) => <option key={item} value={item}>{dictionary.categories[item]}</option>)}</select></label>
      <label>{dictionary.recipe.difficulty}<select name="difficulty" defaultValue={current("difficulty")}><option value="">{dictionary.recipes.any}</option>{(["easy", "medium", "hard"] as const).map((item) => <option key={item} value={item}>{dictionary.difficulty[item]}</option>)}</select></label>
      <label>{dictionary.recipe.diet}<select name="diet" defaultValue={current("diet")}><option value="">{dictionary.recipes.any}</option>{(["vegetarian", "vegan", "pescatarian", "omnivore", "gluten-free"] as const).map((item) => <option key={item} value={item}>{dictionary.diet[item]}</option>)}</select></label>
      <label>{dictionary.recipe.meal}<select name="meal" defaultValue={current("meal")}><option value="">{dictionary.recipes.any}</option>{["breakfast", "lunch", "dinner"].map((item) => <option key={item} value={item}>{dictionary.categories[item as "breakfast" | "lunch" | "dinner"]}</option>)}</select></label>
      <label>{dictionary.recipe.total}<select name="time" defaultValue={current("time")}><option value="">{dictionary.recipes.any}</option>{[20, 30, 45, 60].map((item) => <option key={item} value={item}>{dictionary.recipes.maxTime.replace("{minutes}", String(item))}</option>)}</select></label>
      <label>{dictionary.recipes.sort}<select name="sort" defaultValue={current("sort") || "newest"}>{["newest", "fastest", "easiest", "rated"].map((item) => <option key={item} value={item}>{label(item)}</option>)}</select></label>
      <button className="button button-primary" type="submit">{dictionary.recipes.apply}</button>
    </div></details>
  </form>;
}
