/**
 * Catppuccin Latte (https://catppuccin.com/palette/, palette v1.8.0), the
 * only light Catppuccin flavor. styles/index.css maps these to the theme
 * tokens; charts and diagrams use them directly.
 */
export const latte = {
  rosewater: '#dc8a78',
  flamingo: '#dd7878',
  pink: '#ea76cb',
  mauve: '#8839ef',
  red: '#d20f39',
  maroon: '#e64553',
  peach: '#fe640b',
  yellow: '#df8e1d',
  green: '#40a02b',
  teal: '#179299',
  sky: '#04a5e5',
  sapphire: '#209fb5',
  blue: '#1e66f5',
  lavender: '#7287fd',
  text: '#4c4f69',
  subtext1: '#5c5f77',
  subtext0: '#6c6f85',
  overlay2: '#7c7f93',
  overlay1: '#8c8fa1',
  overlay0: '#9ca0b0',
  surface2: '#acb0be',
  surface1: '#bcc0cc',
  surface0: '#ccd0da',
  base: '#eff1f5',
  mantle: '#e6e9ef',
  crust: '#dce0e8',
} as const;

/** Theme tokens (also defined in styles/index.css; a test keeps them equal). */
export const tokens = {
  /** base mixed 30% with white, like the terminal's 90%-opaque Latte background. */
  page: '#f4f5f8',
  surface: latte.base,
  crust: latte.crust,
  line: latte.surface0,
  ink: latte.text,
  muted: latte.subtext1,
  accent: latte.mauve,
  /** mauve mixed 10% into base; text on it must be ink, not accent (contrast). */
  'accent-soft': '#e5dff4',
  /** yellow mixed 30% into base, for search highlights. */
  mark: '#ead3b4',
  danger: latte.red,
  /** red mixed 8% into base. */
  'danger-soft': '#eddfe6',
} as const;

/** Series colors for charts, in order. */
export const seriesColors = [
  latte.blue,
  latte.mauve,
  latte.teal,
  latte.peach,
  latte.green,
  latte.red,
  latte.sky,
  latte.pink,
  latte.yellow,
  latte.lavender,
];
