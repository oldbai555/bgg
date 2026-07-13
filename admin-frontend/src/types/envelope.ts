export interface Envelope<T = unknown> {
  code: number;
  msg: string;
  data: T;
}

export function isEnvelope(res: unknown): res is Envelope {
  return (
    typeof res === 'object' &&
    res !== null &&
    'code' in res &&
    typeof (res as {code: unknown}).code === 'number'
  );
}
