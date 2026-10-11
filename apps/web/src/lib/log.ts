import pino from 'pino';

/** Browser logger. Only warnings and errors are written, to the console. */
export const log = pino({ browser: { asObject: true }, level: 'warn' });
