"use client";

import { FormEvent, useState } from "react";
import { useRouter } from "next/navigation";
import { useTranslations } from "next-intl";
import { Button } from "@/components/ui/Button";
import { Input } from "@/components/ui/Input";

export function HomeSearch() {
  const [value, setValue] = useState("");
  const router = useRouter();
  const t = useTranslations("home");

  const handleSubmit = (e: FormEvent) => {
    e.preventDefault();
    const query = value.trim();
    router.push(
      `/recipes${query ? `?search=${encodeURIComponent(query)}` : ""}`,
    );
  };

  return (
    <form onSubmit={handleSubmit} className="flex w-full max-w-xl gap-2">
      <Input
        type="search"
        placeholder={t("searchPlaceholder")}
        value={value}
        onChange={(e) => setValue(e.target.value)}
        className="flex-1"
      />
      <Button type="submit" size="lg">
        {t("browseAll")}
      </Button>
    </form>
  );
}
