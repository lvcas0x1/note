const SVG_NS = 'http://www.w3.org/2000/svg';

type IconName = 'copy' | 'check' | 'error';

// 16×16 outline icons drawn with the current text color.
const icons: Record<IconName, string[]> = {
  // Two overlapping rounded sheets (clipboard copy).
  copy: [
    'M5.5 5.5h7a1 1 0 0 1 1 1v7a1 1 0 0 1-1 1h-7a1 1 0 0 1-1-1v-7a1 1 0 0 1 1-1Z',
    'M11.5 5.5v-2a1 1 0 0 0-1-1h-7a1 1 0 0 0-1 1v7a1 1 0 0 0 1 1h2',
  ],
  check: ['M3 8.5l3.25 3.25L13 5'],
  error: ['M4 4l8 8', 'M12 4l-8 8'],
};

const labels: Record<IconName, string> = {
  copy: 'Copy code',
  check: 'Copied',
  error: 'Copy failed',
};

function icon(name: IconName): SVGSVGElement {
  const svg = document.createElementNS(SVG_NS, 'svg');
  svg.setAttribute('viewBox', '0 0 16 16');
  svg.setAttribute('width', '16');
  svg.setAttribute('height', '16');
  svg.setAttribute('fill', 'none');
  svg.setAttribute('stroke', 'currentColor');
  svg.setAttribute('stroke-width', '1.5');
  svg.setAttribute('stroke-linecap', 'round');
  svg.setAttribute('stroke-linejoin', 'round');
  svg.setAttribute('aria-hidden', 'true');
  for (const d of icons[name]) {
    const path = document.createElementNS(SVG_NS, 'path');
    path.setAttribute('d', d);
    svg.append(path);
  }
  return svg;
}

function show(button: HTMLButtonElement, name: IconName) {
  button.replaceChildren(icon(name));
  button.setAttribute('aria-label', labels[name]);
  button.title = labels[name];
  button.dataset.state = name;
}

/** Creates the icon button that copies a code block's text. */
export function createCopyButton(pre: Element): HTMLButtonElement {
  const button = document.createElement('button');
  button.type = 'button';
  button.className =
    'absolute right-2 top-2 rounded-md border border-line bg-page p-1.5 text-muted opacity-70 transition hover:text-ink hover:opacity-100 focus-visible:opacity-100 data-[state=check]:text-accent data-[state=check]:opacity-100';
  show(button, 'copy');
  let reset: ReturnType<typeof setTimeout> | undefined;
  button.addEventListener('click', () => {
    const text = pre.querySelector('code')?.textContent ?? pre.textContent ?? '';
    clearTimeout(reset);
    void navigator.clipboard
      .writeText(text)
      .then(
        () => show(button, 'check'),
        () => show(button, 'error'),
      )
      .finally(() => {
        reset = setTimeout(() => show(button, 'copy'), 1500);
      });
  });
  return button;
}
