import type { InputHTMLAttributes } from 'react';
import { cn } from '@/lib/cn';

export type InputProps = InputHTMLAttributes<HTMLInputElement>;

export function Input({ className, type = 'text', ...props }: InputProps) {
  return (
    <input
      type={type}
      className={cn(
        'h-11 w-full rounded-full border border-beige-dark bg-card px-5 text-sm text-charcoal',
        'placeholder:text-muted/70 transition-colors',
        'focus:border-tomato focus:outline-none',
        className,
      )}
      {...props}
    />
  );
}
