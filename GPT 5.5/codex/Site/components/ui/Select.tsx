import type { ReactNode, SelectHTMLAttributes } from "react";

export function Select({ children, className = "", ...props }: SelectHTMLAttributes<HTMLSelectElement> & { children: ReactNode }) {
  return <select className={className} {...props}>{children}</select>;
}
