"use client";
type Item = {
  file: File;
  progress: number;
  status: "queued" | "uploading" | "finalizing" | "done" | "error";
  errorMsg?: string;
};

export default function UploadList({
  title,
  items,
  onRemove,
  readOnly = false
}: {
  title: string;
  items: Item[];
  onRemove?: (index: number) => void;
  readOnly?: boolean;
}) {
  return (
    <div className="mt-5">
      <div className="text-indigo-700 font-semibold mb-2">{title}</div>
      {items.map((it, idx) => (
        <div
          key={idx}
          className={`flex items-center justify-between border rounded-lg p-2 mb-2 ${
            it.status === "done"
              ? "border-green-400 bg-green-50"
              : it.status === "error"
              ? "border-red-300 bg-red-50"
              : "border-indigo-200 bg-white"
          }`}
        >
          <div className="truncate text-indigo-900 mr-3">{it.file.name}</div>
          <div className="flex items-center gap-2 min-w-[160px]">
            {it.status !== "done" && it.status !== "error" && (
              <div className="w-[140px] h-1.5 bg-indigo-100 rounded-full overflow-hidden">
                <div
                  className="h-full bg-indigo-500"
                  style={{ width: `${it.progress}%` }}
                />
              </div>
            )}
            {it.status === "done" && <span className="text-green-600 font-bold">✓</span>}
            {it.status === "error" && <span className="text-red-600 font-bold" title={it.errorMsg}>!</span>}
            {!readOnly && (
              <button
                className="text-indigo-400 hover:text-red-500"
                onClick={() => onRemove?.(idx)}
              >
                ✖
              </button>
            )}
          </div>
        </div>
      ))}
    </div>
  );
}
