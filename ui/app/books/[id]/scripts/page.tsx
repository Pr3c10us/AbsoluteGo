"use client";

import { memo, useState, useCallback, useEffect, useRef } from "react";
import Link from "next/link";
import { useParams, useSearchParams } from "next/navigation";
import { useInfiniteQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { useQuery } from "@tanstack/react-query";
import { toast } from "sonner";
import { ArrowLeft, Trash2, FileText, Eye, GitBranch, Plus, Check, X } from "lucide-react";
import {
    fetchBooks,
    fetchChapters,
    fetchScripts,
    deleteScript,
    generateScript,
    ApiError,
    type Book,
    type Chapter,
    type Script,
} from "@/lib/api";
import { useScrollLock } from "@/lib/use-scroll-lock";
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

const PAGE_LIMIT = 20;

// ── Static icons (hoisted — rendering-hoist-jsx) ────────────────────────────

const ArrowLeftIcon = <ArrowLeft className="h-4 w-4" />;

const TrashIcon = <Trash2 className="h-3.5 w-3.5" />;

const ScriptEmptyIcon = <FileText className="mx-auto mb-3 h-10 w-10 text-neutral-300" strokeWidth={1.5} />;

const EyeIcon = <Eye className="h-3.5 w-3.5" />;

const SplitIcon = <GitBranch className="h-3.5 w-3.5" />;

const PlusIcon = <Plus className="h-4 w-4" />;

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

const CheckIcon = <Check className="h-3.5 w-3.5" strokeWidth={2.5} />;

// ── Script content viewer overlay ───────────────────────────────────────────

const ScriptViewer = memo(function ScriptViewer({
    script,
    onClose,
}: {
    script: Script;
    onClose: () => void;
}) {
    useScrollLock();

    const CloseIcon = <X className="h-4 w-4" />;

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
    sentinelRef,
    isFetchingNextPage,
    hasNextPage,
}: {
    isLoading: boolean;
    fetchError: Error | null;
    scripts: Script[];
    bookId: number;
    onView: (script: Script) => void;
    onDelete: (script: Script) => void;
    deleteDisabled: boolean;
    sentinelRef: React.RefObject<HTMLDivElement | null>;
    isFetchingNextPage: boolean;
    hasNextPage: boolean;
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
        <>
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
            {/* Infinite scroll sentinel */}
            <div ref={sentinelRef} className="h-1" />
            {isFetchingNextPage ? (
                <div className="flex items-center justify-center gap-2 py-6 text-sm text-muted-foreground">
                    <span className="h-4 w-4 animate-spin rounded-full border-2 border-border border-t-foreground" />
                    Loading more…
                </div>
            ) : !hasNextPage && scripts.length >= PAGE_LIMIT ? (
                <p className="py-6 text-center text-xs text-muted-foreground">
                    All scripts loaded
                </p>
            ) : null}
        </>
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
    const searchParams = useSearchParams();

    // -- state
    const [modalOpen, setModalOpen] = useState(false);
    const [scriptName, setScriptName] = useState("");
    const [selectedChapters, setSelectedChapters] = useState<number[]>([]);
    const [selectedPrevScripts, setSelectedPrevScripts] = useState<number[]>([]);
    const [viewingScript, setViewingScript] = useState<Script | null>(null);
    const [confirmScript, setConfirmScript] = useState<Script | null>(null);

    // -- infinite scroll sentinel
    const sentinelRef = useRef<HTMLDivElement | null>(null);

    // -- stable callbacks
    const handleView = useCallback((script: Script) => setViewingScript(script), []);
    const closeViewer = useCallback(() => setViewingScript(null), []);
    const handleDeleteClick = useCallback((script: Script) => setConfirmScript(script), []);
    const handleConfirmOpenChange = useCallback(
        (open: boolean) => { if (!open) setConfirmScript(null); },
        []
    );

    // -- fetch book info (high limit to ensure the current book is found)
    const { data: booksData } = useQuery({
        queryKey: ["books-all"],
        queryFn: () => fetchBooks({ page: 1, limit: 500 }),
    });
    const book: Book | undefined = booksData?.data?.books?.find(
        (b) => b.id === bookId
    );

    // -- fetch chapters (for modal — no pagination needed, fetch all)
    const { data: chaptersData } = useQuery({
        queryKey: ["chapters-all", bookId],
        queryFn: () => fetchChapters(bookId),
        enabled: !isNaN(bookId) && bookId > 0,
    });
    const chapters: Chapter[] = chaptersData?.data?.chapters ?? [];

    // -- fetch scripts (infinite scroll)
    const {
        data: scriptsData,
        isLoading,
        error: fetchError,
        fetchNextPage,
        hasNextPage,
        isFetchingNextPage,
    } = useInfiniteQuery({
        queryKey: ["scripts", bookId],
        queryFn: ({ pageParam }) =>
            fetchScripts(bookId, { page: pageParam, limit: PAGE_LIMIT }),
        initialPageParam: 1,
        getNextPageParam: (lastPage, allPages) => {
            const fetched = lastPage.data?.scripts?.length ?? 0;
            return fetched < PAGE_LIMIT ? undefined : allPages.length + 1;
        },
        enabled: !isNaN(bookId) && bookId > 0,
    });
    const scripts: Script[] = scriptsData?.pages.flatMap((p) => p.data?.scripts ?? []) ?? [];

    // -- intersection observer for sentinel
    useEffect(() => {
        const el = sentinelRef.current;
        if (!el) return;
        const observer = new IntersectionObserver(
            (entries) => {
                if (entries[0].isIntersecting && hasNextPage && !isFetchingNextPage) {
                    fetchNextPage();
                }
            },
            { rootMargin: "200px" }
        );
        observer.observe(el);
        return () => observer.disconnect();
    }, [hasNextPage, isFetchingNextPage, fetchNextPage]);

    // -- auto-open script viewer from ?scriptId= query param (event tracker deep link)
    useEffect(() => {
        const sid = searchParams.get("scriptId");
        if (!sid || scripts.length === 0) return;
        const target = scripts.find((s) => s.id === Number(sid));
        if (target) setViewingScript(target);
    }, [searchParams, scripts]);

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
            const name = scriptName.trim();
            // Fire-and-forget — events tracker polls for status
            toast.info(`Generating "${name}"…`, { duration: 3000 });
            generateScript({
                bookId,
                name,
                chapters: selectedChapters,
                previousScripts: selectedPrevScripts.length > 0 ? selectedPrevScripts : undefined,
            }).catch((err) => {
                toast.error(
                    err instanceof ApiError
                        ? err.isValidationError
                            ? err.validationErrors.map((v) => v.message).join(", ")
                            : err.businessError
                        : "Generation failed — please retry"
                );
            });
            setModalOpen(false);
            setScriptName("");
            setSelectedChapters([]);
            setSelectedPrevScripts([]);
        },
        [scriptName, selectedChapters, selectedPrevScripts, bookId]
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
                        sentinelRef={sentinelRef}
                        isFetchingNextPage={isFetchingNextPage}
                        hasNextPage={hasNextPage}
                    />
                </section>
            </div>
        </>
    );
}
