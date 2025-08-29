const API_BASE = process.env.NEXT_PUBLIC_API_BASE!;

export async function startUpload(fileName: string, mimeType: string, trace: Record<string,string>) {
  const r = await fetch(`${API_BASE}/api/v1/cvs/uploads`, {
    method: "POST",
    headers: { "Content-Type": "application/json", ...trace },
    body: JSON.stringify({ file_name: fileName, mime_type: mimeType })
  });
  if (!r.ok) throw new Error(`start upload failed ${r.status}`);
  return r.json();
}

export async function completeUpload(id: string, trace: Record<string,string>) {
  const r = await fetch(`${API_BASE}/api/v1/cvs/${id}/complete`, {
    method: "POST",
    headers: { ...trace }
  });
  if (!r.ok) throw new Error(`complete failed ${r.status}`);
  return r.json();
}
