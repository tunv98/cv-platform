export function putWithProgress(
  url: string,
  file: File,
  contentType: string,
  traceHeaders: Record<string, string>,
  onProgress: (pct: number) => void
): Promise<void> {
  return new Promise((resolve, reject) => {
    const xhr = new XMLHttpRequest();
    xhr.open("PUT", url, true);
    xhr.setRequestHeader("Content-Type", contentType);
    Object.entries(traceHeaders).forEach(([k, v]) => xhr.setRequestHeader(k, v));
    xhr.upload.onprogress = (evt) => {
      if (evt.lengthComputable) onProgress(Math.round((evt.loaded / evt.total) * 100));
    };
    xhr.onload = () => (xhr.status >= 200 && xhr.status < 300) ? resolve() : reject(new Error(`GCS upload failed: ${xhr.status}`));
    xhr.onerror = () => reject(new Error("Network error during upload"));
    xhr.send(file);
  });
}
