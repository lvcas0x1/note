/** A neutral message box used for empty and loading states. */
export function Placeholder({ title, children }: { title: string; children?: React.ReactNode }) {
  return (
    <div className="rounded-lg border border-line bg-surface p-6 text-muted">
      <p className="font-medium text-ink">{title}</p>
      {children ? <div className="mt-1 text-sm">{children}</div> : null}
    </div>
  );
}
