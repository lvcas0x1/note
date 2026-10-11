import { Link } from 'react-router';
import { Placeholder } from '../components/Placeholder.tsx';

export function NotFoundPage() {
  return (
    <Placeholder title="Page not found">
      <Link className="text-ink underline" to="/">
        Back to the start page
      </Link>
    </Placeholder>
  );
}
