"use client";

import { useState } from "react";

export function ShareButton({ label, title }: { label: string; title: string }) {
  const [done, setDone] = useState(false);
  async function share() {
    const payload = { title, url: window.location.href };
    if (navigator.share) await navigator.share(payload);
    else { await navigator.clipboard.writeText(payload.url); setDone(true); window.setTimeout(() => setDone(false), 1600); }
  }
  return <button type="button" className="button button-secondary" onClick={() => void share()}>{done ? "✓" : "↗"} {done ? "Copied" : label}</button>;
}
