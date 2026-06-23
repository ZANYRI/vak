"use client";

import { motion, useReducedMotion } from "framer-motion";
import { type ReactNode } from "react";

export function AnimatedSection({
  children,
  className,
  delay = 0,
}: {
  children: ReactNode;
  className?: string;
  delay?: number;
}) {
  const shouldReduceMotion = useReducedMotion();

  const initial = shouldReduceMotion ? undefined : { opacity: 0, y: 16 };
  const whileInView = shouldReduceMotion ? undefined : { opacity: 1, y: 0 };

  return (
    <motion.section
      className={className}
      initial={initial}
      whileInView={whileInView}
      viewport={{ once: true, margin: "-50px" }}
      transition={{ duration: 0.5, delay, ease: "easeOut" }}
    >
      {children}
    </motion.section>
  );
}
