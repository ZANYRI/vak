import type { ReactNode } from 'react';

type Props = {
  icon?: ReactNode;
  title: string;
  body?: string;
  action?: ReactNode;
};

export function EmptyState({ icon, title, body, action }: Props) {
  return (
    <div className="flex flex-col items-center justify-center rounded-card border border-dashed border-beige-dark bg-ivory px-6 py-16 text-center">
      {icon && (
        <div className="mb-4 flex h-16 w-16 items-center justify-center rounded-full bg-beige text-tomato">
          {icon}
        </div>
      )}
      <h2 className="text-xl font-semibold text-charcoal">{title}</h2>
      {body && <p className="mt-2 max-w-md text-sm text-muted">{body}</p>}
      {action && <div className="mt-6">{action}</div>}
    </div>
  );
}
