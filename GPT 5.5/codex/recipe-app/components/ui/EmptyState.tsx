import Link from "next/link";

export function EmptyState({ title, description, href, action }: { title: string; description: string; href?: string; action?: string }) {
  return <div className="empty-state"><span aria-hidden="true">✦</span><h2>{title}</h2><p>{description}</p>{href && action ? <Link className="button button-primary" href={href}>{action}</Link> : null}</div>;
}
