import { createBrowserRouter } from 'react-router';
import { appRoutes } from './App.tsx';

export function createRouter(basename: string) {
  return createBrowserRouter(appRoutes(), { basename: basename.replace(/\/$/, '') || '/' });
}
