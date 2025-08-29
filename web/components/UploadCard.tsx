"use client";
import { useState, useCallback, useMemo, useRef } from "react";
import Dropzone from "./Dropzone";
import UploadList from "./UploadList";
import { startUpload, completeUpload } from "@/lib/api";
import { putWithProgress } from "@/lib/uploader";
import { makeTraceHeaders } from "@/lib/trace";

type Item = {
  file: File;
  progress: number;
  status: "queued" | "uploading" | "finalizing" | "done" | "error";
  id?: string;
  objectKey?: string;
  errorMsg?: string;
};

export default function UploadCard() {
  const [items, setItems] = useState<Item[]>([]);
  const [uploaded, setUploaded] = useState<string[]>([]);
  const inputRef = useRef<HTMLInputElement | null>(null);

  const onPick = useCallback((files: FileList | null) => {
    if (!files) return;
    const accept = Array.from(files).filter(
      f => /pdf|msword|officedocument|png|jpeg|jpg|gif|ppt|presentation/i.test(f.type)
        || /\.(pdf|docx?|pptx?|png|jpe?g|gif)$/i.test(f.name)
    );
    if (!accept.length) return;
    setItems(prev => [...prev, ...accept.map(f => ({ file: f, progress: 0, status: "queued" as const }))]);
  }, []);

  const removeAt = (i: number) => setItems(prev => prev.filter((_, idx) => idx !== i));
  const patchAt = (i: number, patch: Partial<Item>) =>
    setItems(prev => { const copy = [...prev]; copy[i] = { ...copy[i], ...patch }; return copy; });

  const uploadingLabel = useMemo(() => {
    const u = items.filter(i => i.status !== "done" && i.status !== "error").length;
    const t = items.length;
    return `Uploading – ${u}/${t} file${t === 1 ? "" : "s"}`;
  }, [items]);

  const onUpload = async () => {
    if (!items.length) return;
    const trace = makeTraceHeaders();

    for (let i = 0; i < items.length; i++) {
      const it = items[i];
      if (it.status === "done") continue;
      try {
        patchAt(i, { status: "uploading", progress: 0 });

        const start = await startUpload(it.file.name, it.file.type || "application/octet-stream", trace);
        console.log(start);
        await putWithProgress(start.signed_url, it.file, it.file.type || "application/octet-stream", trace, pct => {
          patchAt(i, { progress: pct });
        });

        patchAt(i, { status: "finalizing", progress: 100 });
        await completeUpload(start.id, trace);

        patchAt(i, { status: "done", id: start.id, objectKey: start.object_key });
        setUploaded(prev => [...prev, it.file.name]);
      } catch (err: any) {
        patchAt(i, { status: "error", errorMsg: err?.message ?? "upload error" });
      }
    }
  };

  return (
    <div className="w-[640px] max-w-[95vw] bg-white rounded-2xl p-8 shadow-2xl">
      <h1 className="text-center text-4xl font-extrabold text-indigo-900 mb-6">Upload Your Resume</h1>

      <Dropzone onDrop={onPick} onClick={() => inputRef.current?.click()} inputRef={inputRef} />

      {items.length > 0 && (
        <UploadList title={uploadingLabel} items={items} onRemove={removeAt} />
      )}

      {uploaded.length > 0 && (
        <UploadList
          title="Uploaded"
          items={uploaded.map(n => ({ file: new File([], n), progress: 100, status: "done" } as Item))}
          readOnly
        />
      )}

      <button
        className="w-full mt-4 py-3 rounded-xl bg-indigo-600 text-white font-bold disabled:opacity-50"
        onClick={onUpload}
        disabled={!items.length}
      >
        UPLOAD FILES
      </button>
    </div>
  );
}
