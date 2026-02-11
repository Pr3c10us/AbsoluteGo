"use client";

import {
    createContext,
    useContext,
    useCallback,
    useState,
    type ReactNode,
} from "react";
import { toast } from "sonner";
import { addChapter, generateScript, generateSplits, ApiError } from "@/lib/api";
import { useQueryClient } from "@tanstack/react-query";

// ── Types ───────────────────────────────────────────────────────────────────

interface UploadTask {
    type: "upload";
    id: string;
    bookId: number;
    chapterNumber: number;
    fileName: string;
    status: "uploading" | "done" | "error";
    error?: string;
}

interface ScriptTask {
    type: "script";
    id: string;
    bookId: number;
    scriptName: string;
    status: "uploading" | "done" | "error";
    error?: string;
}

interface SplitTask {
    type: "split";
    id: string;
    scriptId: number;
    scriptName: string;
    status: "uploading" | "done" | "error";
    error?: string;
}

export type BackgroundTask = UploadTask | ScriptTask | SplitTask;

interface UploadContextValue {
    tasks: BackgroundTask[];
    uploadChapter: (
        bookId: number,
        chapterNumber: number,
        file: File
    ) => void;
    generateScriptTask: (params: {
        bookId: number;
        name: string;
        chapters: number[];
        previousScripts?: number[];
    }) => void;
    generateSplitsTask: (scriptId: number, scriptName: string) => void;
}

// ── Context ─────────────────────────────────────────────────────────────────

const UploadContext = createContext<UploadContextValue | null>(null);

export function useUpload() {
    const ctx = useContext(UploadContext);
    if (!ctx) throw new Error("useUpload must be used within UploadProvider");
    return ctx;
}

// ── Provider ────────────────────────────────────────────────────────────────

let taskCounter = 0;

export function UploadProvider({ children }: { children: ReactNode }) {
    const [tasks, setTasks] = useState<BackgroundTask[]>([]);
    const queryClient = useQueryClient();

    const uploadChapter = useCallback(
        (bookId: number, chapterNumber: number, file: File) => {
            const id = `upload-${++taskCounter}`;
            const task: UploadTask = {
                type: "upload",
                id,
                bookId,
                chapterNumber,
                fileName: file.name,
                status: "uploading",
            };

            setTasks((prev) => [...prev, task]);
            toast.info(`Uploading Ch.${chapterNumber}…`, {
                description: file.name,
                duration: 3000,
            });

            addChapter(bookId, chapterNumber, file)
                .then((res) => {
                    setTasks((prev) =>
                        prev.map((t) => (t.id === id ? { ...t, status: "done" as const } : t))
                    );
                    queryClient.invalidateQueries({
                        queryKey: ["chapters", bookId],
                    });
                    toast.success(res.message || `Ch.${chapterNumber} uploaded`, {
                        description: file.name,
                    });
                    setTimeout(() => {
                        setTasks((prev) => prev.filter((t) => t.id !== id));
                    }, 5000);
                })
                .catch((err) => {
                    const message =
                        err instanceof ApiError
                            ? err.businessError
                            : "Upload failed — please retry";
                    setTasks((prev) =>
                        prev.map((t) =>
                            t.id === id ? { ...t, status: "error" as const, error: message } : t
                        )
                    );
                    toast.error(`Ch.${chapterNumber} failed`, {
                        description: message,
                    });
                    setTimeout(() => {
                        setTasks((prev) => prev.filter((t) => t.id !== id));
                    }, 8000);
                });
        },
        [queryClient]
    );

    const generateScriptTask = useCallback(
        (params: {
            bookId: number;
            name: string;
            chapters: number[];
            previousScripts?: number[];
        }) => {
            const id = `script-${++taskCounter}`;
            const task: ScriptTask = {
                type: "script",
                id,
                bookId: params.bookId,
                scriptName: params.name,
                status: "uploading",
            };

            setTasks((prev) => [...prev, task]);
            toast.info(`Generating "${params.name}"…`, {
                duration: 3000,
            });

            generateScript(params)
                .then(() => {
                    setTasks((prev) =>
                        prev.map((t) => (t.id === id ? { ...t, status: "done" as const } : t))
                    );
                    queryClient.invalidateQueries({
                        queryKey: ["scripts", params.bookId],
                    });
                    toast.success(`"${params.name}" generated`);
                    setTimeout(() => {
                        setTasks((prev) => prev.filter((t) => t.id !== id));
                    }, 5000);
                })
                .catch((err) => {
                    const message =
                        err instanceof ApiError
                            ? err.isValidationError
                                ? err.validationErrors.map((v) => v.message).join(", ")
                                : err.businessError
                            : "Generation failed — please retry";
                    setTasks((prev) =>
                        prev.map((t) =>
                            t.id === id ? { ...t, status: "error" as const, error: message } : t
                        )
                    );
                    toast.error(`"${params.name}" failed`, {
                        description: message,
                    });
                    setTimeout(() => {
                        setTasks((prev) => prev.filter((t) => t.id !== id));
                    }, 8000);
                });
        },
        [queryClient]
    );

    const generateSplitsTask = useCallback(
        (scriptId: number, scriptName: string) => {
            const id = `split-${++taskCounter}`;
            const task: SplitTask = {
                type: "split",
                id,
                scriptId,
                scriptName,
                status: "uploading",
            };

            setTasks((prev) => [...prev, task]);
            toast.info(`Generating splits for "${scriptName}"…`, {
                duration: 3000,
            });

            generateSplits(scriptId)
                .then(() => {
                    setTasks((prev) =>
                        prev.map((t) => (t.id === id ? { ...t, status: "done" as const } : t))
                    );
                    queryClient.invalidateQueries({
                        queryKey: ["splits", scriptId],
                    });
                    toast.success(`Splits for "${scriptName}" generated`);
                    setTimeout(() => {
                        setTasks((prev) => prev.filter((t) => t.id !== id));
                    }, 5000);
                })
                .catch((err) => {
                    const message =
                        err instanceof ApiError
                            ? err.isValidationError
                                ? err.validationErrors.map((v) => v.message).join(", ")
                                : err.businessError
                            : "Split generation failed — please retry";
                    setTasks((prev) =>
                        prev.map((t) =>
                            t.id === id ? { ...t, status: "error" as const, error: message } : t
                        )
                    );
                    toast.error(`Splits for "${scriptName}" failed`, {
                        description: message,
                    });
                    setTimeout(() => {
                        setTasks((prev) => prev.filter((t) => t.id !== id));
                    }, 8000);
                });
        },
        [queryClient]
    );

    return (
        <UploadContext.Provider value={{ tasks, uploadChapter, generateScriptTask, generateSplitsTask }}>
            {children}
        </UploadContext.Provider>
    );
}
