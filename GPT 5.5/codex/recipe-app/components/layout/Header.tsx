"use client";

import Link from "next/link";
import { useState } from "react";
import type { Locale } from "@/i18n/config";
import type { Dictionary } from "@/i18n/dictionary";
import { LanguageSwitcher } from "./LanguageSwitcher";

export function Header({ locale, dictionary }: { locale: Locale; dictionary: Dictionary }) {
  const [open, setOpen] = useState(false);
  const navItems = [["", dictionary.nav.home], ["/recipes", dictionary.nav.recipes], ["/categories", dictionary.nav.categories], ["/cuisines", dictionary.nav.cuisines], ["/favorites", dictionary.nav.favorites], ["/about", dictionary.nav.about]];
  return <header className="site-header"><div className="shell header-inner">
    <Link className="brand" href={`/${locale}`} onClick={() => setOpen(false)} aria-label={dictionary.siteName}>mis<span>e</span>n</Link>
    <button className="menu-button" type="button" onClick={() => setOpen((value) => !value)} aria-expanded={open} aria-controls="primary-navigation"><span className="sr-only">{open ? dictionary.common.close : dictionary.common.menu}</span><span aria-hidden="true">{open ? "×" : "☰"}</span></button>
    <nav id="primary-navigation" className={open ? "nav-links is-open" : "nav-links"} aria-label="Primary navigation">{navItems.map(([path, label]) => <Link key={path} href={`/${locale}${path}`} onClick={() => setOpen(false)}>{label}</Link>)}<LanguageSwitcher locale={locale} dictionary={dictionary} onNavigate={() => setOpen(false)} /></nav>
  </div></header>;
}
