import type { ReactNode } from "react";
import "./globals.css";

export const metadata = {
  title: {
    default: "Saveur — Bilingual Recipes",
    absolute: "Saveur — Bilingual Recipes"
  },
  description: "Discover bilingual recipes from around the world.",
  robots: { index: true, follow: true },
  icons: { icon: "/favicon.svg" }
};

export default function RootLayout({ children }: { children: ReactNode }) {
  return (
    <html lang="en" suppressHydrationWarning>
      <body>{children}</body>
    </html>
  );
}
