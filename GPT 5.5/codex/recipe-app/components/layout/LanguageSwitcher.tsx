"use client";

import Link from "next/link";
import { usePathname, useSearchParams } from "next/navigation";
import type { Locale } from "@/i18n/config";
import type { Dictionary } from "@/i18n/dictionary";

export function LanguageSwitcher({ locale, dictionary, onNavigate }: { locale: Locale; dictionary: Dictionary; onNavigate?: () => void }) {
  const pathname = usePathname(); const search = useSearchParams();
  const suffix = search.size ? `?${search.toString()}` : "";
  return <div className="language-switcher" aria-label={dictionary.common.language}>{(["ru", "en"] as const).map((target) => <Link key={target} href={`${pathname.replace(/^\/(en|ru)(?=\/|$)/, `/${target}`)}${suffix}`} className={target === locale ? "active" : ""} aria-current={target === locale ? "true" : undefined} onClick={onNavigate}>{target.toUpperCase()}</Link>)}</div>;
}
