import type { SelectHTMLAttributes } from 'react';
import { cn } from '@/lib/cn';

export type SelectOption = {
  value: string;
  label: string;
};

export type SelectProps = SelectHTMLAttributes<HTMLSelectElement> & {
  options: SelectOption[];
};

export function Select({ className, options, ...props }: SelectProps) {
  return (
    <div className="relative">
      <select
        className={cn(
          'h-11 w-full appearance-none rounded-full border border-beige-dark bg-card px-5 pr-10 text-sm text-charcoal',
          'transition-colors focus:border-tomato focus:outline-none',
          className,
        )}
        {...props}
      >
        {options.map((option) => (
          <option key={option.value} value={option.value}>
            {option.label}
          </option>
        ))}
      </select>
      <svg
        aria-hidden="true"
        viewBox="0 0 20 20"
        className="pointer-events-none absolute right-4 top-1/2 h-4 w-4 -translate-y-1/2 text-muted"
        fill="none"
        stroke="currentColor"
        strokeWidth="2"
      >
        <path d="M6 8l4 4 4-4" strokeLinecap="round" strokeLinejoin="round" />
      </svg>
    </div>
  );
}
