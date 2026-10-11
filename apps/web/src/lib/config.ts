/** Runtime configuration from /config.js (written by the Go server or by `notes export`). */
export interface NotesConfig {
  mode: 'server' | 'static';
  basePath: string;
  videoBaseUrl?: string;
}

declare global {
  interface Window {
    __NOTES_CONFIG__?: NotesConfig;
  }
}

export function readConfig(): NotesConfig {
  const c = typeof window === 'undefined' ? undefined : window.__NOTES_CONFIG__;
  const basePath = c?.basePath && c.basePath.startsWith('/') ? c.basePath : '/';
  return {
    mode: c?.mode === 'static' ? 'static' : 'server',
    basePath: basePath.endsWith('/') ? basePath : `${basePath}/`,
  };
}
