import Link from "next/link";
import type { Locale } from "@/i18n/config";
import type { Dictionary } from "@/i18n/dictionary";

export function Footer({ locale, dictionary }: { locale: Locale; dictionary: Dictionary }) {
  const year = new Date().getFullYear();
  return <footer className="site-footer"><div className="shell footer-grid"><div><Link className="brand" href={`/${locale}`}>mis<span>e</span>n</Link><p>{dictionary.footer.tagline}</p></div><nav aria-label="Footer navigation"><Link href={`/${locale}/recipes`}>{dictionary.nav.recipes}</Link><Link href={`/${locale}/categories`}>{dictionary.nav.categories}</Link><Link href={`/${locale}/about`}>{dictionary.nav.about}</Link></nav><p>{dictionary.footer.copyright.replace("{year}", String(year))}</p></div></footer>;
}
