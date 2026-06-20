import { useTranslations } from 'next-intl';
import { AnimatedSection } from '@/components/animations/AnimatedSection';
import type { Locale, Recipe } from '@/lib/types';

type Props = {
  recipe: Recipe;
  locale: Locale;
};

export function InstructionSteps({ recipe, locale }: Props) {
  const t = useTranslations('Recipe');
  const steps = recipe.steps[locale];

  return (
    <section aria-labelledby="instructions-heading">
      <h2 id="instructions-heading" className="font-display text-2xl font-semibold text-charcoal">
        {t('instructions')}
      </h2>
      <ol className="mt-6 space-y-5">
        {steps.map((step, index) => (
          <AnimatedSection as="li" key={index} delay={Math.min(index, 6) * 0.04} className="flex gap-4 list-none">
            <span
              aria-hidden="true"
              className="flex h-10 w-10 shrink-0 items-center justify-center rounded-full bg-tomato font-display text-lg font-semibold text-white shadow-soft"
            >
              {index + 1}
            </span>
            <div className="pt-1">
              <p className="sr-only">{t('step', { number: index + 1 })}</p>
              <p className="text-charcoal leading-relaxed">{step}</p>
            </div>
          </AnimatedSection>
        ))}
      </ol>
    </section>
  );
}
