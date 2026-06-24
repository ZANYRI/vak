import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: { default: "Misen — Recipes worth making", template: "%s | Misen" },
  description: "A bilingual collection of considered recipes for the curious cook."
};

export default function RootLayout({ children }: Readonly<{ children: React.ReactNode }>) {
  return <html lang="en"><body>{children}</body></html>;
}
