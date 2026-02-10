"use client";

import {
    createContext,
    useContext,
    useCallback,
    useState,
    type ReactNode,
} from "react";
import { toast } from "sonner";
import { addChapter, ApiError } from "@/lib/api";
import { useQueryClient } from "@tanstack/react-query";

// ── Types ───────────────────────────────────────────────────────────────────

interface Upload {
    id: string;
    bookId: number;
    chapterNumber: number;
    fileName: string;
    status: "uploading" | "done" | "error";
    error?: string;
}

interface UploadContextValue {
    uploads: Upload[];
    uploadChapter: (
        bookId: number,
        chapterNumber: number,
        file: File
    ) => void;
}

// ── Context ─────────────────────────────────────────────────────────────────

const UploadContext = createContext<UploadContextValue | null>(null);

export function useUpload() {
    const ctx = useContext(UploadContext);
    if (!ctx) throw new Error("useUpload must be used within UploadProvider");
    return ctx;
}

// ── Provider ────────────────────────────────────────────────────────────────

let uploadCounter = 0;

export function UploadProvider({ children }: { children: ReactNode }) {
    const [uploads, setUploads] = useState<Upload[]>([]);
    const queryClient = useQueryClient();

    const uploadChapter = useCallback(
        (bookId: number, chapterNumber: number, file: File) => {
            const id = `upload-${++uploadCounter}`;
            const upload: Upload = {
                id,
                bookId,
                chapterNumber,
                fileName: file.name,
                status: "uploading",
            };

            setUploads((prev) => [...prev, upload]);
            toast.info(`Uploading Ch.${chapterNumber}…`, {
                description: file.name,
                duration: 3000,
            });

            addChapter(bookId, chapterNumber, file)
                .then((res) => {
                    setUploads((prev) =>
                        prev.map((u) => (u.id === id ? { ...u, status: "done" } : u))
                    );
                    // Invalidate chapters so any page showing them auto-refreshes
                    queryClient.invalidateQueries({
                        queryKey: ["chapters", bookId],
                    });
                    toast.success(res.message || `Ch.${chapterNumber} uploaded`, {
                        description: file.name,
                    });
                    // Remove from list after a delay
                    setTimeout(() => {
                        setUploads((prev) => prev.filter((u) => u.id !== id));
                    }, 5000);
                })
                .catch((err) => {
                    const message =
                        err instanceof ApiError
                            ? err.businessError
                            : "Upload failed — please retry";
                    setUploads((prev) =>
                        prev.map((u) =>
                            u.id === id ? { ...u, status: "error", error: message } : u
                        )
                    );
                    toast.error(`Ch.${chapterNumber} failed`, {
                        description: message,
                    });
                    // Remove from list after a delay
                    setTimeout(() => {
                        setUploads((prev) => prev.filter((u) => u.id !== id));
                    }, 8000);
                });
        },
        [queryClient]
    );

    return (
        <UploadContext.Provider value={{ uploads, uploadChapter }}>
            {children}
        </UploadContext.Provider>
    );
}
