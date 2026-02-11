"use client";

import { memo, useState, useCallback } from "react";
import Link from "next/link";
import { useParams } from "next/navigation";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import {
    fetchBooks,
    fetchChapters,
    fetchScripts,
    deleteScript,
    ApiError,
    type Book,
    type Chapter,
    type Script,
} from "@/lib/api";
import { useUpload } from "@/lib/upload-context";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
    Dialog,
    DialogContent,
    DialogHeader,
    DialogTitle,
    DialogDescription,
    DialogFooter,
} from "@/components/ui/dialog";
import {
    AlertDialog,
    AlertDialogAction,
    AlertDialogCancel,
    AlertDialogContent,
    AlertDialogDescription,
    AlertDialogFooter,
    AlertDialogHeader,
    AlertDialogTitle,
} from "@/components/ui/alert-dialog";

// ── Static SVG icons (hoisted — rendering-hoist-jsx) ────────────────────────

const ArrowLeftIcon = (
    <svg
        width="16"
        height="16"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        strokeWidth="2"
        strokeLinecap="round"
        strokeLinejoin="round"
    >
        <path d="m12 19-7-7 7-7" />
        <path d="M19 12H5" />
    </svg>
);

const TrashIcon = (
    <svg
        width="14"
        height="14"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        strokeWidth="2"
        strokeLinecap="round"
        strokeLinejoin="round"
    >
        <path d="M3 6h18M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2" />
    </svg>
);

const ScriptEmptyIcon = (
    <svg
        className="mx-auto mb-3 h-10 w-10 text-neutral-300"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        strokeWidth="1.5"
        strokeLinecap="round"
        strokeLinejoin="round"
    >
        <path d="M15 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V7Z" />
        <path d="M14 2v4a2 2 0 0 0 2 2h4" />
        <path d="M10 13H8" />
        <path d="M16 17H8" />
        <path d="M16 13h-2" />
    </svg>
);

const EyeIcon = (
    <svg
        width="14"
        height="14"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        strokeWidth="2"
        strokeLinecap="round"
        strokeLinejoin="round"
    >
        <path d="M2.062 12.348a1 1 0 0 1 0-.696 10.75 10.75 0 0 1 19.876 0 1 1 0 0 1 0 .696 10.75 10.75 0 0 1-19.876 0" />
        <circle cx="12" cy="12" r="3" />
    </svg>
);

const SplitIcon = (
    <svg
        width="14"
        height="14"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        strokeWidth="2"
        strokeLinecap="round"
        strokeLinejoin="round"
    >
        <path d="M16 3h5v5" />
        <path d="M8 3H3v5" />
        <path d="M12 22v-8.3a4 4 0 0 0-1.172-2.872L3 3" />
        <path d="m15 9 6-6" />
    </svg>
);

const PlusIcon = (
    <svg
        width="16"
        height="16"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        strokeWidth="2"
        strokeLinecap="round"
        strokeLinejoin="round"
    >
        <path d="M5 12h14" />
        <path d="M12 5v14" />
    </svg>
);

const HeroUnderline = (
    <svg
        className="mt-2 h-2 w-24 text-foreground"
        viewBox="0 0 120 8"
        fill="none"
    >
        <path
            d="M2 5C25 2 50 7 75 4C100 1 115 6 118 3"
            stroke="currentColor"
            strokeWidth="2"
            strokeLinecap="round"
        />
    </svg>
);

const CheckIcon = (
    <svg
        width="14"
        height="14"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        strokeWidth="2.5"
        strokeLinecap="round"
        strokeLinejoin="round"
    >
        <path d="M20 6 9 17l-5-5" />
    </svg>
);

// ── Script content viewer overlay ───────────────────────────────────────────

const ScriptViewer = memo(function ScriptViewer({
    script,
    onClose,
}: {
    script: Script;
    onClose: () => void;
}) {
    const CloseIcon = (
        <svg
            width="16"
            height="16"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            strokeWidth="2"
            strokeLinecap="round"
            strokeLinejoin="round"
        >
            <path d="M18 6 6 18" />
            <path d="m6 6 12 12" />
        </svg>
    );

    return (
        <div className="fixed inset-0 z-40 flex flex-col bg-white">
            {/* ── Header ── */}
            <div className="flex items-center justify-between border-b border-border px-6 py-3">
                <div>
                    <h2 className="text-lg font-bold tracking-tight">
                        {script.name}
                    </h2>
                    <span className="text-xs text-muted-foreground">
                        Script #{script.id} — Chapters: {script.chapters.join(", ")}
                    </span>
                </div>
                <Button
                    variant="outline"
                    size="icon-sm"
                    onClick={onClose}
                    aria-label="Close script viewer"
                >
                    {CloseIcon}
                </Button>
            </div>

            {/* ── Content ── */}
            <div className="flex-1 overflow-y-auto p-6 pb-24 max-sm:p-4 max-sm:pb-24">
                <div className="mx-auto max-w-3xl">
                    <p className="whitespace-pre-wrap text-sm leading-relaxed text-foreground">
                        {script.content}
                    </p>
                </div>
            </div>
        </div>
    );
});

// ── Script card ─────────────────────────────────────────────────────────────

const ScriptCard = memo(function ScriptCard({
    script,
    bookId,
    onView,
    onDelete,
    deleteDisabled,
}: {
    script: Script;
    bookId: number;
    onView: (script: Script) => void;
    onDelete: (script: Script) => void;
    deleteDisabled: boolean;
}) {
    return (
        <li className="group relative overflow-hidden rounded-[6px_8px_7px_5px] border border-border shadow-[3px_5px_14px_rgba(0,0,0,0.06)] transition-all duration-300 hover:shadow-[4px_7px_24px_rgba(0,0,0,0.15)] hover:-translate-y-0.5 animate-in fade-in-0 zoom-in-95">
            <div className="p-4">
                {/* ── Top row: name + delete ── */}
                <div className="mb-3 flex items-start justify-between gap-2">
                    <div className="min-w-0">
                        <h3 className="truncate text-lg font-bold tracking-tight">
                            {script.name}
                        </h3>
                        <span className="text-[10px] font-medium uppercase tracking-[0.25em] text-muted-foreground">
                            Chapters: {script.chapters.join(", ")}
                        </span>
                    </div>
                    <div className="flex items-center gap-1">
                        <span className="rounded-[3px_5px_4px_3px] bg-foreground/5 px-1.5 py-0.5 font-mono text-[10px] text-muted-foreground">
                            ID {script.id}
                        </span>
                        <Button
                            variant="ghost"
                            size="icon-sm"
                            onClick={() => onDelete(script)}
                            disabled={deleteDisabled}
                            aria-label={`Delete ${script.name}`}
                            className="h-7 w-7 opacity-0 transition-opacity group-hover:opacity-100"
                        >
                            {TrashIcon}
                        </Button>
                    </div>
                </div>

                {/* ── Preview ── */}
                <p className="mb-4 line-clamp-3 text-sm leading-relaxed text-muted-foreground">
                    {script.content}
                </p>

                {/* ── CTAs ── */}
                <div className="flex items-center gap-2">
                    <Button
                        variant="secondary"
                        size="sm"
                        onClick={() => onView(script)}
                        className="gap-1.5"
                    >
                        {EyeIcon}
                        View Script
                    </Button>
                    <Button
                        variant="outline"
                        size="sm"
                        asChild
                        className="gap-1.5"
                    >
                        <Link href={`/books/${bookId}/scripts/${script.id}/splits`}>
                            {SplitIcon}
                            View Splits
                        </Link>
                    </Button>
                </div>
            </div>
        </li>
    );
});

// ── Scripts list content (rerender-memo) ────────────────────────────────────

const ScriptsListContent = memo(function ScriptsListContent({
    isLoading,
    fetchError,
    scripts,
    bookId,
    onView,
    onDelete,
    deleteDisabled,
}: {
    isLoading: boolean;
    fetchError: Error | null;
    scripts: Script[];
    bookId: number;
    onView: (script: Script) => void;
    onDelete: (script: Script) => void;
    deleteDisabled: boolean;
}) {
    if (isLoading) {
        return (
            <div className="flex items-center justify-center gap-2 py-8 text-sm text-muted-foreground">
                <span className="h-4 w-4 animate-spin rounded-full border-2 border-border border-t-foreground" />
                Loading scripts…
            </div>
        );
    }

    if (fetchError) {
        return (
            <div className="py-8 text-center text-sm font-medium text-foreground">
                Failed to load scripts —{" "}
                {fetchError instanceof ApiError ? fetchError.message : "network error"}
            </div>
        );
    }

    return scripts.length === 0 ? (
        <div className="py-12 text-center text-muted-foreground">
            {ScriptEmptyIcon}
            <p className="text-sm font-medium text-foreground">No scripts yet.</p>
            <span className="text-xs">
                Generate your first script using the button above.
            </span>
        </div>
    ) : (
        <ul className="grid grid-cols-1 gap-4 sm:grid-cols-2">
            {scripts.map((script) => (
                <ScriptCard
                    key={script.id}
                    script={script}
                    bookId={bookId}
                    onView={onView}
                    onDelete={onDelete}
                    deleteDisabled={deleteDisabled}
                />
            ))}
        </ul>
    );
});

// ── Multi-select checkbox item ──────────────────────────────────────────────

function CheckboxItem({
    checked,
    onChange,
    label,
}: {
    checked: boolean;
    onChange: (checked: boolean) => void;
    label: string;
}) {
    return (
        <button
            type="button"
            onClick={() => onChange(!checked)}
            className={`flex items-center gap-2 rounded-md border px-3 py-2 text-sm transition-colors ${checked
                ? "border-foreground bg-foreground text-background"
                : "border-border bg-white text-foreground hover:bg-neutral-50"
                }`}
        >
            <span
                className={`flex h-4 w-4 shrink-0 items-center justify-center rounded-[3px] border transition-colors ${checked
                    ? "border-background bg-background"
                    : "border-neutral-300"
                    }`}
            >
                {checked ? (
                    <span className="text-foreground">{CheckIcon}</span>
                ) : null}
            </span>
            {label}
        </button>
    );
}

// ── Main page ───────────────────────────────────────────────────────────────

export default function ScriptsPage() {
    const params = useParams();
    const bookId = Number(params.id);
    const queryClient = useQueryClient();
    const { generateScriptTask } = useUpload();

    // -- state
    const [modalOpen, setModalOpen] = useState(false);
    const [scriptName, setScriptName] = useState("");
    const [selectedChapters, setSelectedChapters] = useState<number[]>([]);
    const [selectedPrevScripts, setSelectedPrevScripts] = useState<number[]>([]);
    const [viewingScript, setViewingScript] = useState<Script | null>(null);
    const [confirmScript, setConfirmScript] = useState<Script | null>(null);

    // -- stable callbacks
    const handleView = useCallback((script: Script) => setViewingScript(script), []);
    const closeViewer = useCallback(() => setViewingScript(null), []);
    const handleDeleteClick = useCallback((script: Script) => setConfirmScript(script), []);
    const handleConfirmOpenChange = useCallback(
        (open: boolean) => { if (!open) setConfirmScript(null); },
        []
    );

    // -- fetch book info
    const { data: booksData } = useQuery({
        queryKey: ["books"],
        queryFn: () => fetchBooks(),
    });
    const book: Book | undefined = booksData?.data?.books?.find(
        (b) => b.id === bookId
    );

    // -- fetch chapters (for the modal dropdown)
    const { data: chaptersData } = useQuery({
        queryKey: ["chapters", bookId],
        queryFn: () => fetchChapters(bookId),
        enabled: !isNaN(bookId) && bookId > 0,
    });
    const chapters: Chapter[] = chaptersData?.data?.chapters ?? [];

    // -- fetch scripts
    const {
        data: scriptsData,
        isLoading,
        error: fetchError,
    } = useQuery({
        queryKey: ["scripts", bookId],
        queryFn: () => fetchScripts(bookId),
        enabled: !isNaN(bookId) && bookId > 0,
    });
    const scripts: Script[] = scriptsData?.data?.scripts ?? [];

    // -- delete script mutation
    const deleteMutation = useMutation({
        mutationFn: (scriptId: number) => deleteScript(scriptId),
        onSuccess: () => {
            setConfirmScript(null);
            queryClient.invalidateQueries({ queryKey: ["scripts", bookId] });
            toast.success("Script deleted");
        },
        onError: (err) => {
            setConfirmScript(null);
            toast.error(
                err instanceof ApiError ? err.businessError : "Delete failed"
            );
        },
    });

    const handleConfirmDelete = useCallback(() => {
        if (confirmScript) {
            deleteMutation.mutate(confirmScript.id);
        }
    }, [confirmScript, deleteMutation]);

    const handleGenerate = useCallback(
        (e: React.FormEvent) => {
            e.preventDefault();
            if (!scriptName.trim() || selectedChapters.length === 0) return;
            generateScriptTask({
                bookId,
                name: scriptName.trim(),
                chapters: selectedChapters,
                previousScripts: selectedPrevScripts.length > 0 ? selectedPrevScripts : undefined,
            });
            setModalOpen(false);
            setScriptName("");
            setSelectedChapters([]);
            setSelectedPrevScripts([]);
        },
        [scriptName, selectedChapters, selectedPrevScripts, bookId, generateScriptTask]
    );

    const toggleChapter = useCallback((chapterNum: number) => {
        setSelectedChapters((prev) =>
            prev.includes(chapterNum)
                ? prev.filter((n) => n !== chapterNum)
                : [...prev, chapterNum]
        );
    }, []);

    const togglePrevScript = useCallback((scriptId: number) => {
        setSelectedPrevScripts((prev) =>
            prev.includes(scriptId)
                ? prev.filter((id) => id !== scriptId)
                : [...prev, scriptId]
        );
    }, []);

    const bookTitle = book?.title ?? `Book #${bookId}`;
    const isGenerateDisabled =
        !scriptName.trim() ||
        selectedChapters.length === 0;

    return (
        <>
            {/* ── Script viewer overlay ── */}
            {viewingScript !== null ? (
                <ScriptViewer script={viewingScript} onClose={closeViewer} />
            ) : null}

            {/* ── Delete confirmation ── */}
            <AlertDialog
                open={confirmScript !== null}
                onOpenChange={handleConfirmOpenChange}
            >
                <AlertDialogContent>
                    <AlertDialogHeader>
                        <AlertDialogTitle>Delete script?</AlertDialogTitle>
                        <AlertDialogDescription>
                            &ldquo;{confirmScript?.name}&rdquo; and all its splits will be
                            permanently removed.
                        </AlertDialogDescription>
                    </AlertDialogHeader>
                    <AlertDialogFooter>
                        <AlertDialogCancel>Cancel</AlertDialogCancel>
                        <AlertDialogAction onClick={handleConfirmDelete}>
                            Delete
                        </AlertDialogAction>
                    </AlertDialogFooter>
                </AlertDialogContent>
            </AlertDialog>

            {/* ── Generate Script Modal ── */}
            <Dialog open={modalOpen} onOpenChange={setModalOpen}>
                <DialogContent>
                    <DialogHeader>
                        <DialogTitle>Generate New Script</DialogTitle>
                        <DialogDescription>
                            Select chapters and optionally reference previous scripts for
                            continuity.
                        </DialogDescription>
                    </DialogHeader>

                    <form onSubmit={handleGenerate} className="space-y-5">
                        {/* Script name */}
                        <div>
                            <label
                                htmlFor="script-name"
                                className="mb-1.5 block text-xs font-medium uppercase tracking-widest text-muted-foreground"
                            >
                                Script Name
                            </label>
                            <Input
                                id="script-name"
                                placeholder="e.g. Chapter 1"
                                value={scriptName}
                                onChange={(e) => setScriptName(e.target.value)}
                            />
                        </div>

                        {/* Chapter selection */}
                        <div>
                            <label className="mb-1.5 block text-xs font-medium uppercase tracking-widest text-muted-foreground">
                                Chapters
                            </label>
                            {chapters.length === 0 ? (
                                <p className="text-sm text-muted-foreground">
                                    No chapters available. Upload chapters first.
                                </p>
                            ) : (
                                <div className="flex flex-wrap gap-2">
                                    {chapters.map((ch) => (
                                        <CheckboxItem
                                            key={ch.id}
                                            checked={selectedChapters.includes(ch.number)}
                                            onChange={() => toggleChapter(ch.number)}
                                            label={`Ch. ${ch.number}`}
                                        />
                                    ))}
                                </div>
                            )}
                        </div>

                        {/* Previous scripts (optional) */}
                        {scripts.length > 0 ? (
                            <div>
                                <label className="mb-1.5 block text-xs font-medium uppercase tracking-widest text-muted-foreground">
                                    Previous Scripts (optional)
                                </label>
                                <div className="flex flex-wrap gap-2">
                                    {scripts.map((s) => (
                                        <CheckboxItem
                                            key={s.id}
                                            checked={selectedPrevScripts.includes(s.id)}
                                            onChange={() => togglePrevScript(s.id)}
                                            label={s.name}
                                        />
                                    ))}
                                </div>
                            </div>
                        ) : null}

                        <DialogFooter>
                            <Button
                                type="submit"
                                disabled={isGenerateDisabled}
                            >
                                {"Generate Script"}
                            </Button>
                        </DialogFooter>
                    </form>
                </DialogContent>
            </Dialog>

            <div className="mx-auto max-w-5xl px-6 pb-20 max-sm:px-4">
                {/* ── Hero ── */}
                <header className="relative pb-10 pt-20 max-sm:pb-7 max-sm:pt-12">
                    <Link
                        href={`/books/${bookId}`}
                        className="mb-6 inline-flex items-center gap-1.5 text-sm text-muted-foreground transition-colors hover:text-foreground"
                    >
                        {ArrowLeftIcon}
                        {bookTitle}
                    </Link>
                    <h1 className="text-[clamp(2.5rem,8vw,4.5rem)] font-black leading-[0.85] tracking-tighter max-sm:text-4xl">
                        Scripts
                    </h1>
                    {HeroUnderline}
                    <span className="mt-4 block text-[11px] font-medium uppercase tracking-[0.3em] text-muted-foreground">
                        SCRIPT MANAGEMENT
                    </span>
                </header>

                {/* ── Generate CTA ── */}
                <section className="border-t border-border pb-8 pt-10">
                    <div className="flex items-center justify-between gap-3">
                        <h2 className="text-2xl font-semibold tracking-tight">
                            Generate
                        </h2>
                        <Button
                            onClick={() => setModalOpen(true)}
                            className="gap-1.5"
                        >
                            {PlusIcon}
                            New Script
                        </Button>
                    </div>
                </section>

                {/* ── Scripts List ── */}
                <section className="border-t border-border pt-8">
                    <div className="mb-6 flex flex-wrap items-center justify-between gap-3">
                        <h2 className="flex items-center gap-2.5 text-2xl font-semibold tracking-tight">
                            Scripts
                            {scripts.length > 0 ? (
                                <span className="rounded-[4px_6px_5px_3px] bg-foreground px-2 py-0.5 text-xs font-medium text-background">
                                    {scripts.length}
                                </span>
                            ) : null}
                        </h2>
                    </div>

                    <ScriptsListContent
                        isLoading={isLoading}
                        fetchError={fetchError}
                        scripts={scripts}
                        bookId={bookId}
                        onView={handleView}
                        onDelete={handleDeleteClick}
                        deleteDisabled={deleteMutation.isPending}
                    />
                </section>
            </div>
        </>
    );
}
