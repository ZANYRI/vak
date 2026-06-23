"use client";

import { motion, useReducedMotion } from "framer-motion";
import { useTranslations } from "next-intl";
import type { Recipe } from "@/lib/validation";

type InstructionStepsProps = {
  recipe: Recipe;
  locale: "en" | "ru";
};

export function InstructionSteps({ recipe, locale }: InstructionStepsProps) {
  const t = useTranslations("recipe");
  const shouldReduceMotion = useReducedMotion();

  return (
    <section>
      <h2 className="mb-4 text-xl font-bold">{t("instructions")}</h2>
      <ol className="space-y-4">
        {recipe.steps[locale].map((step, index) => {
          const initial = shouldReduceMotion
            ? undefined
            : { opacity: 0, x: -12 };
          const whileInView = shouldReduceMotion
            ? undefined
            : { opacity: 1, x: 0 };
          return (
            <motion.li
              key={index}
              initial={initial}
              whileInView={whileInView}
              viewport={{ once: true }}
              transition={{ duration: 0.3, delay: index * 0.05 }}
              className="border-border bg-card flex gap-4 rounded-md border p-4"
            >
              <span className="bg-primary text-primary-foreground flex h-8 w-8 shrink-0 items-center justify-center rounded-full text-sm font-bold">
                {index + 1}
              </span>
              <p className="text-foreground/90">{step}</p>
            </motion.li>
          );
        })}
      </ol>
    </section>
  );
}
