import type { HTMLAttributes } from 'react';
import { cn } from '@/lib/cn';

type Tone = 'neutral' | 'tomato' | 'olive' | 'saffron';

const tones: Record<Tone, string> = {
  neutral: 'bg-beige text-charcoal',
  tomato: 'bg-tomato/12 text-tomato-dark',
  olive: 'bg-olive/12 text-olive-dark',
  saffron: 'bg-saffron/18 text-[#8a5d00]',
};

export type BadgeProps = HTMLAttributes<HTMLSpanElement> & {
  tone?: Tone;
};

export function Badge({ className, tone = 'neutral', ...props }: BadgeProps) {
  return (
    <span
      className={cn(
        'inline-flex items-center rounded-full px-3 py-1 text-xs font-medium',
        tones[tone],
        className,
      )}
      {...props}
    />
  );
}
