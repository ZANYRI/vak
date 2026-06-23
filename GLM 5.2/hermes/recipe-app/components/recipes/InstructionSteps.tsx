"use client";

import { useLocale, useTranslations } from "next-intl";
import { motion } from "motion/react";
import type { Recipe } from "@/lib/types";

export function InstructionSteps({ recipe }: { recipe: Recipe }) {
  const locale = useLocale() as "en" | "ru";
  const t = useTranslations("recipe");
  const steps = recipe.steps[locale];

  return (
    <section
      aria-labelledby="instructions-heading"
      className="rounded-card border border-border bg-card p-5 shadow-soft"
    >
      <h2 id="instructions-heading" className="font-serif text-xl font-semibold text-foreground">
        {t("instructions")}
      </h2>
      <ol className="mt-4 space-y-5">
        {steps.map((step, i) => (
          <motion.li
            key={i}
            initial={{ opacity: 0, x: -12 }}
            whileInView={{ opacity: 1, x: 0 }}
            viewport={{ once: true, margin: "-30px" }}
            transition={{ duration: 0.4, delay: Math.min(i * 0.05, 0.3) }}
            className="flex gap-4"
          >
            <span
              aria-hidden="true"
              className="flex h-8 w-8 flex-none items-center justify-center rounded-full bg-primary font-sans text-sm font-semibold text-primary-foreground"
            >
              {i + 1}
            </span>
            <div className="pt-1">
              <p className="sr-only">
                {t("step")} {i + 1}
              </p>
              <p className="text-sm leading-relaxed text-foreground">{step}</p>
            </div>
          </motion.li>
        ))}
      </ol>
    </section>
  );
}
