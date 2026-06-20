'use client';

import Image, { type ImageProps } from 'next/image';
import { useState } from 'react';
import { useTranslations } from 'next-intl';
import { cn } from '@/lib/cn';
import { FALLBACK_IMAGE } from '@/lib/constants';

type Props = Omit<ImageProps, 'onError' | 'onLoad'> & {
  /** Wrapper className (the image itself fills the wrapper). */
  wrapperClassName?: string;
};

/**
 * next/image with a skeleton loading state and a graceful fallback when the
 * remote URL fails to load.
 */
export function ImageWithFallback({
  src,
  alt,
  wrapperClassName,
  className,
  ...props
}: Props) {
  const t = useTranslations('Image');
  const [failed, setFailed] = useState(false);
  const [loaded, setLoaded] = useState(false);

  const effectiveSrc = failed ? FALLBACK_IMAGE : src;

  return (
    <div className={cn('relative overflow-hidden bg-beige', wrapperClassName)}>
      {!loaded && (
        <div className="skeleton absolute inset-0" aria-hidden="true" />
      )}
      <Image
        src={effectiveSrc}
        alt={alt}
        className={cn(
          'h-full w-full object-cover transition-opacity duration-500',
          loaded ? 'opacity-100' : 'opacity-0',
          className,
        )}
        onLoad={() => setLoaded(true)}
        onError={() => {
          setFailed(true);
          setLoaded(true);
        }}
        {...props}
      />
      {failed && (
        <span className="sr-only" role="status">
          {t('unavailable')}
        </span>
      )}
    </div>
  );
}
