import { AnimatedSection } from '@/components/animations/AnimatedSection';

type Props = {
  title: string;
  subtitle?: string;
  eyebrow?: string;
};

export function PageHeader({ title, subtitle, eyebrow }: Props) {
  return (
    <AnimatedSection className="mb-10">
      {eyebrow && (
        <p className="mb-2 text-sm font-semibold uppercase tracking-wide text-tomato">
          {eyebrow}
        </p>
      )}
      <h1 className="font-display text-4xl font-semibold text-charcoal text-balance sm:text-5xl">
        {title}
      </h1>
      {subtitle && (
        <p className="mt-3 max-w-2xl text-lg text-muted">{subtitle}</p>
      )}
    </AnimatedSection>
  );
}
