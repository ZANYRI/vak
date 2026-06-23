"use client";

import Image, { type ImageProps } from "next/image";
import { useState } from "react";
import { cn } from "@/lib/cn";

type ImageWithFallbackProps = Omit<ImageProps, "onError"> & {
  fallbackSrc?: string;
  /** Optional skeleton shown while the primary image loads. */
  showSkeleton?: boolean;
};

const DEFAULT_FALLBACK =
  "data:image/svg+xml;utf8," +
  encodeURIComponent(
    `<svg xmlns="http://www.w3.org/2000/svg" width="1200" height="800" viewBox="0 0 1200 800">
      <rect width="1200" height="800" fill="#f6efe1"/>
      <g fill="none" stroke="#c9bda6" stroke-width="3">
        <circle cx="600" cy="360" r="90"/>
        <path d="M600 270 v-40 M600 450 v40 M510 360 h-40 M690 360 h40 M536 296 l-28-28 M664 424 l28 28 M536 424 l-28 28 M664 296 l28-28"/>
      </g>
      <text x="600" y="540" font-family="Georgia, serif" font-size="44" fill="#7c7268" text-anchor="middle">Saveur</text>
      <text x="600" y="590" font-family="sans-serif" font-size="24" fill="#a89c87" text-anchor="middle">image unavailable</text>
    </svg>`
  );

export function ImageWithFallback({
  fallbackSrc = DEFAULT_FALLBACK,
  showSkeleton = true,
  className,
  alt,
  ...props
}: ImageWithFallbackProps) {
  const [errored, setErrored] = useState(false);
  const [loaded, setLoaded] = useState(false);

  const src = errored ? fallbackSrc : props.src;

  return (
    <span
      className={cn(
        "relative block overflow-hidden",
        showSkeleton && !loaded && !errored ? "skeleton" : "",
        className
      )}
    >
      <Image
        {...props}
        src={src}
        alt={alt}
        onLoad={() => setLoaded(true)}
        onError={() => setErrored(true)}
        className={cn(
          "transition-opacity duration-500",
          loaded || errored ? "opacity-100" : "opacity-0",
          className
        )}
      />
    </span>
  );
}
