import { Component, type ReactNode } from 'react';
import { log } from '../../lib/log.ts';
import { embedLabels, type EmbedKind } from './types.ts';

interface Props {
  kind: EmbedKind;
  source: string;
  children: ReactNode;
}

/** Shows an inline error box when one embed fails, so the rest of the article still renders. */
export class EmbedBoundary extends Component<Props, { error: Error | null }> {
  override state = { error: null };

  static getDerivedStateFromError(error: Error) {
    return { error };
  }

  override componentDidCatch(error: Error) {
    log.warn(
      { kind: this.props.kind, source: this.props.source, err: error.message },
      'embed failed',
    );
  }

  override render() {
    if (this.state.error) return <EmbedError kind={this.props.kind} source={this.props.source} />;
    return this.props.children;
  }
}

export function EmbedError({ kind, source }: { kind: EmbedKind; source: string }) {
  return (
    <div
      role="alert"
      className="rounded-lg border border-danger/40 border-l-4 border-l-danger bg-danger-soft p-3 text-sm text-ink"
    >
      <p className="font-medium">This {embedLabels[kind]} could not be shown.</p>
      <p className="mt-1 break-all font-mono text-xs">{source}</p>
    </div>
  );
}
