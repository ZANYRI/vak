import en from "@/i18n/messages/en.json";
import ru from "@/i18n/messages/ru.json";
import type { Locale } from "@/i18n/config";

export type Dictionary = typeof en;

const dictionaries: Record<Locale, Dictionary> = { en, ru };

export function getDictionary(locale: Locale): Dictionary {
  return dictionaries[locale];
}
