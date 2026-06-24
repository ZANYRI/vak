"use client";

import Image from "next/image";
import { useState } from "react";
import { cn } from "@/lib/utils";

type ImageWithFallbackProps = {
  src: string;
  alt: string;
  fill?: boolean;
  width?: number;
  height?: number;
  className?: string;
  containerClassName?: string;
  priority?: boolean;
};

export function ImageWithFallback({
  src,
  alt,
  fill,
  width,
  height,
  className,
  containerClassName,
  priority,
}: ImageWithFallbackProps) {
  const [loaded, setLoaded] = useState(false);
  const [error, setError] = useState(false);

  return (
    <div
      className={cn(
        "bg-muted relative overflow-hidden",
        fill && "h-full w-full",
        containerClassName,
      )}
    >
      {!loaded && !error && (
        <div
          className="bg-muted absolute inset-0 animate-pulse"
          aria-hidden="true"
        />
      )}
      {error ? (
        <div className="bg-muted text-foreground/60 flex h-full w-full items-center justify-center p-4 text-center text-sm">
          <span>{alt}</span>
        </div>
      ) : (
        <Image
          src={src}
          alt={alt}
          fill={fill}
          width={!fill ? width : undefined}
          height={!fill ? height : undefined}
          priority={priority}
          onLoad={() => setLoaded(true)}
          onError={() => setError(true)}
          className={cn(
            "object-cover transition-transform duration-500",
            className,
          )}
          sizes={
            fill
              ? "(max-width: 768px) 100vw, (max-width: 1200px) 50vw, 33vw"
              : undefined
          }
          unoptimized
        />
      )}
    </div>
  );
}
