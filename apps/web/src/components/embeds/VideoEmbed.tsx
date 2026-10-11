import { useEffect, useRef, useState } from 'react';
import type { EmbedProps } from './types.ts';

/** A video player with the browser's controls. It never autoplays and loads only metadata first. */
export default function VideoEmbed({ src }: EmbedProps) {
  const [failed, setFailed] = useState(false);
  const video = useRef<HTMLVideoElement>(null);

  // Media errors can fire before React's handler is attached, so listen
  // directly and also check for an error that already happened.
  useEffect(() => {
    const el = video.current;
    if (!el) return undefined;
    const onError = () => setFailed(true);
    el.addEventListener('error', onError);
    if (el.error) queueMicrotask(onError);
    return () => el.removeEventListener('error', onError);
  }, []);

  return (
    <figure>
      <video
        ref={video}
        controls
        preload="metadata"
        playsInline
        src={src}
        className="w-full rounded-lg bg-ink"
      >
        <a href={src} download>
          Download the video
        </a>
      </video>
      {failed ? (
        <figcaption className="mt-2 text-sm text-muted">
          This browser cannot play the video.{' '}
          <a href={src} download className="text-accent underline">
            Download the video
          </a>
        </figcaption>
      ) : null}
    </figure>
  );
}
