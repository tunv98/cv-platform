function randomHex(bytes: number) {
  const arr = new Uint8Array(bytes);
  crypto.getRandomValues(arr);
  return Array.from(arr).map(b => b.toString(16).padStart(2, "0")).join("");
}

export function makeTraceHeaders() {
  const traceId = randomHex(16);
  const spanId = randomHex(8);
  const traceparent = `00-${traceId}-${spanId}-01`;
  const gcp = `${traceId}/${parseInt(spanId, 16)};o=1`;
  return { traceparent, "X-Cloud-Trace-Context": gcp };
}
