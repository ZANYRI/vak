import type { ComponentPropsWithoutRef, ReactNode } from "react";

type Props = ComponentPropsWithoutRef<"section"> & { children: ReactNode; delay?: number };

export function AnimatedSection({ children, className = "", delay = 0, style, ...props }: Props) {
  return (
    <section
      className={`reveal ${className}`}
      style={{ ...style, animationDelay: `${delay}ms` }}
      {...props}
    >
      {children}
    </section>
  );
}
