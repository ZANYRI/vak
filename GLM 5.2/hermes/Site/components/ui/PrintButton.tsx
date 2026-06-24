"use client";

import { useTranslations } from "next-intl";
import { Button } from "./Button";

export function PrintButton() {
  const t = useTranslations("common");
  return (
    <Button type="button" variant="outline" onClick={() => window.print()}>
      <svg
        viewBox="0 0 24 24"
        width="18"
        height="18"
        fill="none"
        stroke="currentColor"
        strokeWidth={2}
        aria-hidden="true"
      >
        <path
          strokeLinecap="round"
          strokeLinejoin="round"
          d="M6 9V3h12v6M6 18H4a2 2 0 0 1-2-2v-4a2 2 0 0 1 2-2h16a2 2 0 0 1 2 2v4a2 2 0 0 1-2 2h-2M6 14h12v7H6z"
        />
      </svg>
      {t("print")}
    </Button>
  );
}
