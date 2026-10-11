// npm run export: builds the web app and the content, then writes the static
// site to ./static and its videos to ./static-videos (both replaced on every run).
//
// Settings (environment variables):
//   NOTES_BASE_PATH       URL path the site is served under (default "/")
//   NOTES_VIDEO_BASE_URL  absolute URL where ./static-videos will be uploaded
//                         (required when an article embeds a video)
import { execFileSync } from 'node:child_process';
import { rmSync } from 'node:fs';

const out = 'static';
const videos = 'static-videos';
const basePath = process.env.NOTES_BASE_PATH ?? '/';
const videoBaseUrl = process.env.NOTES_VIDEO_BASE_URL ?? '';

// Only these two generated folders are removed.
rmSync(out, { recursive: true, force: true });
rmSync(videos, { recursive: true, force: true });

execFileSync('npm', ['run', 'build', '-w', '@notes/web'], { stdio: 'inherit' });
execFileSync('npm', ['run', 'build:content'], { stdio: 'inherit' });
const args = [
  'run',
  './cmd/notes',
  'export',
  '--data',
  '../data',
  '--out',
  `../${out}`,
  '--base-path',
  basePath,
];
if (videoBaseUrl) args.push('--video-base-url', videoBaseUrl);
execFileSync('go', args, { cwd: 'server', stdio: 'inherit' });
