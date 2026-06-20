import type { ButtonHTMLAttributes } from 'react';
import { cn } from '@/lib/cn';

type Variant = 'primary' | 'secondary' | 'outline' | 'ghost';
type Size = 'sm' | 'md' | 'lg';

const base =
  'inline-flex items-center justify-center gap-2 rounded-full font-medium transition-all duration-200 disabled:pointer-events-none disabled:opacity-50 active:scale-[0.98] motion-reduce:active:scale-100';

const variants: Record<Variant, string> = {
  primary:
    'bg-tomato text-white shadow-soft hover:bg-tomato-dark hover:-translate-y-0.5 motion-reduce:hover:translate-y-0',
  secondary:
    'bg-olive text-white shadow-soft hover:bg-olive-dark hover:-translate-y-0.5 motion-reduce:hover:translate-y-0',
  outline:
    'border border-beige-dark bg-card text-charcoal hover:border-tomato hover:text-tomato',
  ghost: 'text-charcoal hover:bg-beige/60',
};

const sizes: Record<Size, string> = {
  sm: 'h-9 px-4 text-sm',
  md: 'h-11 px-6 text-sm',
  lg: 'h-13 px-8 text-base',
};

export type ButtonProps = ButtonHTMLAttributes<HTMLButtonElement> & {
  variant?: Variant;
  size?: Size;
};

export function Button({
  className,
  variant = 'primary',
  size = 'md',
  type = 'button',
  ...props
}: ButtonProps) {
  return (
    <button
      type={type}
      className={cn(base, variants[variant], sizes[size], className)}
      {...props}
    />
  );
}

/** Shared classes so links (e.g. next-intl <Link>) can be styled as buttons. */
export function buttonClasses(
  variant: Variant = 'primary',
  size: Size = 'md',
  className?: string,
) {
  return cn(base, variants[variant], sizes[size], className);
}
