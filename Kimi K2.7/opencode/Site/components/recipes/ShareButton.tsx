"use client";

import { useState } from "react";
import { Button } from "@/components/ui/Button";

type ShareButtonProps = {
  title: string;
  url: string;
};

export function ShareButton({ title, url }: ShareButtonProps) {
  const [copied, setCopied] = useState(false);

  const handleClick = async () => {
    const shareData = { title, url };
    try {
      if (navigator.canShare?.(shareData)) {
        await navigator.share(shareData);
      } else {
        await navigator.clipboard.writeText(url);
        setCopied(true);
        setTimeout(() => setCopied(false), 2000);
      }
    } catch {
      // user cancelled or not supported
    }
  };

  return (
    <Button variant="outline" onClick={handleClick}>
      {copied ? "Copied!" : "Share"}
    </Button>
  );
}
