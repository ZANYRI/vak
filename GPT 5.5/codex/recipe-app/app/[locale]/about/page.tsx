import { notFound } from "next/navigation";
import { isLocale } from "@/i18n/config";
import { getDictionary } from "@/i18n/dictionary";
import { localizedMetadata } from "@/lib/seo";

export async function generateMetadata({ params }: { params: Promise<{ locale: string }> }) { const { locale } = await params; if (!isLocale(locale)) return {}; const d = getDictionary(locale); return localizedMetadata(locale, "/about", d.about.title, d.about.intro); }
export default async function AboutPage({ params }: { params: Promise<{ locale: string }> }) { const { locale } = await params; if (!isLocale(locale)) notFound(); const d = getDictionary(locale); const sections = [[d.about.philosophyTitle, d.about.philosophy], [d.about.dataTitle, d.about.data], [d.about.imagesTitle, d.about.images], [d.about.accessibilityTitle, d.about.accessibility]]; return <section className="shell about"><div className="about-intro"><p className="eyebrow">Misen</p><h1>{d.about.title}</h1><p>{d.about.intro}</p></div><div className="about-sections">{sections.map(([title, body], index) => <article key={title}><span>{String(index + 1).padStart(2, "0")}</span><h2>{title}</h2><p>{body}</p></article>)}</div></section>; }
