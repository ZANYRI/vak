"use client";

import Image from "next/image";
import { useState } from "react";

type Props = { src: string; alt: string; className?: string; sizes?: string };

export function ImageWithFallback({ src, alt, className = "", sizes = "(max-width: 700px) 100vw, 33vw" }: Props) {
  const [failed, setFailed] = useState(false);
  if (failed) {
    return <div className={`image-fallback ${className}`} role="img" aria-label={alt}><span>✦</span></div>;
  }
  return <Image src={src} alt={alt} fill sizes={sizes} className={className} onError={() => setFailed(true)} />;
}
