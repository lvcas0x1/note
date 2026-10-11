import { z } from 'zod';

// The site's CSP forbids eval. Zod would otherwise probe `new Function` to
// decide whether to compile validators, and the browser reports that probe
// as a CSP violation even though Zod catches it.
z.config({ jitless: true });
