import type { Decorator } from '@storybook/react-vite';
import { MemoryRouter, useLocation } from 'react-router';

/** Shows the current router location so play functions can assert navigation. */
export function LocationProbe() {
  const location = useLocation();
  return (
    <output data-testid="location" className="sr-only">
      {location.pathname}
    </output>
  );
}

/** Wraps a story in a memory router that starts at `parameters.route` (default "/"). */
export const withRouter: Decorator = (Story, { parameters }) => {
  const route = typeof parameters.route === 'string' ? parameters.route : '/';
  return (
    <MemoryRouter initialEntries={[route]}>
      <Story />
      <LocationProbe />
    </MemoryRouter>
  );
};
