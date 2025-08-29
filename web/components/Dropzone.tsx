"use client";
import { RefObject, useState } from "react";

export default function Dropzone({
  onDrop,
  onClick,
  inputRef
}: {
  onDrop: (files: FileList | null) => void;
  onClick: () => void;
  inputRef: RefObject<HTMLInputElement>;
}) {
  const [over, setOver] = useState(false);

  return (
    <div
      className={`border-2 border-dashed rounded-xl p-6 text-center cursor-pointer transition ${
        over ? "border-indigo-500 bg-indigo-50 shadow-inner" : "border-indigo-300 bg-indigo-50/50"
      }`}
      onDragOver={(e) => { e.preventDefault(); setOver(true); }}
      onDragLeave={() => setOver(false)}
      onDrop={(e) => { e.preventDefault(); setOver(false); onDrop(e.dataTransfer.files); }}
      onClick={onClick}
    >
      <div className="text-8xl mb-2">☁️</div>
      <div className="text-indigo-900 mb-1">
        <strong>Drag &amp; drop files</strong> or{" "}
        <span className="text-indigo-500 underline">Browse</span>
      </div>
      <div className="text-xs text-indigo-500">
        Supported formats: JPEG, PNG, GIF, MP4, PDF, PSD, AI, Word, PPT
      </div>
      <input
        ref={inputRef}
        type="file"
        multiple
        accept=".pdf,.doc,.docx,.ppt,.pptx,.png,.jpg,.jpeg,.gif"
        onChange={(e) => onDrop(e.target.files)}
        style={{ display: "none" }}
      />
    </div>
  );
}
