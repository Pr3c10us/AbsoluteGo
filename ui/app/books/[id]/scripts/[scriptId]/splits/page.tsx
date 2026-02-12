"use client";

import { memo, useState, useCallback, useMemo } from "react";
import Link from "next/link";
import { useParams } from "next/navigation";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import {
    fetchBooks,
    fetchScripts,
    fetchSplits,
    deleteSplits,
    ApiError,
    type Book,
    type Script,
    type Split,
} from "@/lib/api";
import { useUpload } from "@/lib/upload-context";
import { Button } from "@/components/ui/button";
import Lightbox, { type LightboxItem } from "@/components/lightbox";
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

const SplitEmptyIcon = (
    <svg
        className="mx-auto mb-3 h-10 w-10 text-neutral-300"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        strokeWidth="1.5"
        strokeLinecap="round"
        strokeLinejoin="round"
    >
        <path d="M16 3h5v5" />
        <path d="M8 3H3v5" />
        <path d="M12 22v-8.3a4 4 0 0 0-1.172-2.872L3 3" />
        <path d="m15 9 6-6" />
    </svg>
);

const RefreshIcon = (
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
        <path d="M21 12a9 9 0 0 0-9-9 9.75 9.75 0 0 0-6.74 2.74L3 8" />
        <path d="M3 3v5h5" />
        <path d="M3 12a9 9 0 0 0 9 9 9.75 9.75 0 0 0 6.74-2.74L21 16" />
        <path d="M16 16h5v5" />
    </svg>
);

const SparklesIcon = (
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
        <path d="M9.937 15.5A2 2 0 0 0 8.5 14.063l-6.135-1.582a.5.5 0 0 1 0-.962L8.5 9.936A2 2 0 0 0 9.937 8.5l1.582-6.135a.5.5 0 0 1 .963 0L14.063 8.5A2 2 0 0 0 15.5 9.937l6.135 1.581a.5.5 0 0 1 0 .964L15.5 14.063a2 2 0 0 0-1.437 1.437l-1.582 6.135a.5.5 0 0 1-.963 0z" />
    </svg>
);

const TrashIcon = (
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
        <path d="M3 6h18M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2" />
    </svg>
);

const ExpandIcon = (
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
        <polyline points="15 3 21 3 21 9" />
        <polyline points="9 21 3 21 3 15" />
        <line x1="21" x2="14" y1="3" y2="10" />
        <line x1="3" x2="10" y1="21" y2="14" />
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

// ── Effect label mapping ────────────────────────────────────────────────────

const EFFECT_LABELS: Record<string, string> = {
    panRight: "Pan Right",
    panLeft: "Pan Left",
    panUp: "Pan Up",
    panDown: "Pan Down",
    zoomIn: "Zoom In",
    zoomOut: "Zoom Out",
};

// ── Script viewer overlay ───────────────────────────────────────────────────

const ScriptViewer = memo(function ScriptViewer({
    script,
    onClose,
}: {
    script: Script;
    onClose: () => void;
}) {
    return (
        <div className="fixed inset-0 z-40 flex flex-col bg-white">
            <div className="flex items-center justify-between border-b border-border px-6 py-3">
                <div>
                    <h2 className="text-lg font-bold tracking-tight">
                        {script.name}
                    </h2>
                    <span className="text-xs text-muted-foreground">
                        Full script — Chapters: {script.chapters.join(", ")}
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

// ── Split card ──────────────────────────────────────────────────────────────

const SplitCard = memo(function SplitCard({
    split,
    index,
    onViewImage,
}: {
    split: Split;
    index: number;
    onViewImage: (index: number) => void;
}) {
    return (
        <li className="group overflow-hidden rounded-[6px_8px_7px_5px] border border-border shadow-[3px_5px_14px_rgba(0,0,0,0.06)] transition-all duration-300 hover:shadow-[4px_7px_24px_rgba(0,0,0,0.15)] hover:-translate-y-0.5 animate-in fade-in-0 zoom-in-95">
            <div className="flex gap-4 p-4 max-sm:flex-col">
                {/* ── Panel image ── */}
                {split.panel?.url ? (
                    <div className="group/img relative h-32 w-24 shrink-0 overflow-hidden rounded-[5px_7px_6px_4px] border border-border max-sm:h-40 max-sm:w-full">
                        {/* eslint-disable-next-line @next/next/no-img-element */}
                        <img
                            src={split.panel.url}
                            alt={`Panel ${split.panel.panelNumber}`}
                            className="absolute inset-0 h-full w-full object-cover transition-transform duration-500 group-hover/img:scale-105"
                        />
                        {/* Hover overlay with View Big CTA */}
                        <div className="absolute inset-0 flex items-center justify-center bg-black/0 opacity-0 transition-all duration-300 group-hover/img:bg-black/40 group-hover/img:opacity-100">
                            <Button
                                variant="secondary"
                                size="sm"
                                onClick={() => onViewImage(index)}
                                className="gap-1.5 bg-white text-black shadow-lg hover:bg-neutral-100"
                            >
                                {ExpandIcon}
                                View
                            </Button>
                        </div>
                        <div className="absolute bottom-1 left-1 rounded-[3px_5px_4px_3px] bg-black/70 px-1.5 py-0.5 font-mono text-[10px] font-medium text-white backdrop-blur-sm">
                            P{split.panel.panelNumber}
                        </div>
                    </div>
                ) : null}

                {/* ── Content ── */}
                <div className="flex min-w-0 flex-1 flex-col justify-between">
                    <div>
                        <div className="mb-1.5 flex items-center gap-2">
                            <span className="text-sm font-bold tracking-tight">
                                Split {String(index + 1).padStart(2, "0")}
                            </span>
                            <span className="rounded-[3px_5px_4px_3px] bg-foreground/5 px-1.5 py-0.5 font-mono text-[10px] text-muted-foreground">
                                ID {split.id}
                            </span>
                            <span className="rounded-[3px_5px_4px_3px] bg-foreground px-1.5 py-0.5 text-[10px] font-medium text-background">
                                {EFFECT_LABELS[split.effect] ?? split.effect}
                            </span>
                        </div>
                        <p className="text-sm leading-relaxed text-muted-foreground">
                            {split.content}
                        </p>
                    </div>
                </div>
            </div>
        </li>
    );
});

// ── Splits list content (rerender-memo) ─────────────────────────────────────

const SplitsListContent = memo(function SplitsListContent({
    isLoading,
    fetchError,
    splits,
    onViewImage,
}: {
    isLoading: boolean;
    fetchError: Error | null;
    splits: Split[];
    onViewImage: (index: number) => void;
}) {
    if (isLoading) {
        return (
            <div className="flex items-center justify-center gap-2 py-8 text-sm text-muted-foreground">
                <span className="h-4 w-4 animate-spin rounded-full border-2 border-border border-t-foreground" />
                Loading splits…
            </div>
        );
    }

    if (fetchError) {
        return (
            <div className="py-8 text-center text-sm font-medium text-foreground">
                Failed to load splits —{" "}
                {fetchError instanceof ApiError ? fetchError.message : "network error"}
            </div>
        );
    }

    return splits.length === 0 ? (
        <div className="py-12 text-center text-muted-foreground">
            {SplitEmptyIcon}
            <p className="text-sm font-medium text-foreground">No splits yet.</p>
            <span className="text-xs">
                Generate splits for this script using the button above.
            </span>
        </div>
    ) : (
        <ul className="flex flex-col gap-3">
            {splits.map((split, idx) => (
                <SplitCard
                    key={split.id}
                    split={split}
                    index={idx}
                    onViewImage={onViewImage}
                />
            ))}
        </ul>
    );
});

// ── Main page ───────────────────────────────────────────────────────────────

export default function SplitsPage() {
    const params = useParams();
    const bookId = Number(params.id);
    const scriptId = Number(params.scriptId);
    const queryClient = useQueryClient();
    const { generateSplitsTask, isSplitGenerating } = useUpload();
    const splitInProgress = isSplitGenerating(scriptId);

    // -- state
    const [confirmClear, setConfirmClear] = useState(false);
    const [lightboxIdx, setLightboxIdx] = useState<number | null>(null);
    const [viewingScript, setViewingScript] = useState(false);

    // -- stable callbacks
    const handleClearOpenChange = useCallback(
        (open: boolean) => { if (!open) setConfirmClear(false); },
        []
    );
    const handleViewImage = useCallback(
        (index: number) => setLightboxIdx(index),
        []
    );
    const closeLightbox = useCallback(() => setLightboxIdx(null), []);
    const closeScriptViewer = useCallback(() => setViewingScript(false), []);

    // -- fetch book info
    const { data: booksData } = useQuery({
        queryKey: ["books"],
        queryFn: () => fetchBooks(),
    });
    const book: Book | undefined = booksData?.data?.books?.find(
        (b) => b.id === bookId
    );

    // -- fetch script info
    const { data: scriptsData } = useQuery({
        queryKey: ["scripts", bookId],
        queryFn: () => fetchScripts(bookId),
        enabled: !isNaN(bookId) && bookId > 0,
    });
    const script: Script | undefined = (scriptsData?.data?.scripts ?? []).find(
        (s) => s.id === scriptId
    );

    // -- fetch splits
    const {
        data: splitsData,
        isLoading,
        error: fetchError,
    } = useQuery({
        queryKey: ["splits", scriptId],
        queryFn: () => fetchSplits(scriptId),
        enabled: !isNaN(scriptId) && scriptId > 0,
    });
    const splits: Split[] = splitsData?.data?.splits ?? [];
    const hasSplits = splits.length > 0;

    // -- lightbox items from splits
    const lightboxItems: LightboxItem[] = useMemo(
        () =>
            splits
                .filter((s) => s.panel?.url)
                .map((s, i) => ({
                    url: s.panel.url,
                    label: `Split ${String(i + 1).padStart(2, "0")} — P${s.panel.panelNumber}`,
                })),
        [splits]
    );

    // -- delete splits mutation
    const clearMutation = useMutation({
        mutationFn: () => deleteSplits(scriptId),
        onSuccess: () => {
            setConfirmClear(false);
            queryClient.invalidateQueries({ queryKey: ["splits", scriptId] });
            toast.success("Splits cleared");
        },
        onError: (err) => {
            setConfirmClear(false);
            toast.error(
                err instanceof ApiError ? err.businessError : "Clear failed"
            );
        },
    });

    const handleConfirmClear = useCallback(() => {
        clearMutation.mutate();
    }, [clearMutation]);

    const bookTitle = book?.title ?? `Book #${bookId}`;
    const scriptName = script?.name ?? `Script #${scriptId}`;

    return (
        <>
            {/* ── Panel image lightbox ── */}
            {lightboxIdx !== null ? (
                <Lightbox
                    items={lightboxItems}
                    currentIndex={lightboxIdx}
                    onIndexChange={setLightboxIdx}
                    onClose={closeLightbox}
                />
            ) : null}

            {/* ── Full script viewer ── */}
            {viewingScript && script ? (
                <ScriptViewer script={script} onClose={closeScriptViewer} />
            ) : null}

            {/* ── Clear confirmation ── */}
            <AlertDialog
                open={confirmClear}
                onOpenChange={handleClearOpenChange}
            >
                <AlertDialogContent>
                    <AlertDialogHeader>
                        <AlertDialogTitle>Clear all splits?</AlertDialogTitle>
                        <AlertDialogDescription>
                            All {splits.length} split{splits.length !== 1 ? "s" : ""} for
                            &ldquo;{scriptName}&rdquo; will be permanently removed. You can
                            regenerate them afterwards.
                        </AlertDialogDescription>
                    </AlertDialogHeader>
                    <AlertDialogFooter>
                        <AlertDialogCancel>Cancel</AlertDialogCancel>
                        <AlertDialogAction onClick={handleConfirmClear}>
                            Clear Splits
                        </AlertDialogAction>
                    </AlertDialogFooter>
                </AlertDialogContent>
            </AlertDialog>

            <div className="mx-auto max-w-5xl px-6 pb-20 max-sm:px-4">
                {/* ── Hero ── */}
                <header className="relative pb-10 pt-20 max-sm:pb-7 max-sm:pt-12">
                    <Link
                        href={`/books/${bookId}/scripts`}
                        className="mb-6 inline-flex items-center gap-1.5 text-sm text-muted-foreground transition-colors hover:text-foreground"
                    >
                        {ArrowLeftIcon}
                        {bookTitle} — Scripts
                    </Link>
                    <h1 className="text-[clamp(2.5rem,8vw,4.5rem)] font-black leading-[0.85] tracking-tighter max-sm:text-4xl">
                        {scriptName}
                    </h1>
                    {HeroUnderline}
                    <span className="mt-4 block text-[11px] font-medium uppercase tracking-[0.3em] text-muted-foreground">
                        SCRIPT SPLITS
                    </span>
                </header>

                {/* ── Actions ── */}
                <section className="border-t border-border pb-8 pt-10">
                    <div className="flex items-center justify-between gap-3">
                        <h2 className="text-2xl font-semibold tracking-tight">
                            Splits
                            {hasSplits ? (
                                <span className="ml-2.5 rounded-[4px_6px_5px_3px] bg-foreground px-2 py-0.5 text-xs font-medium text-background">
                                    {splits.length}
                                </span>
                            ) : null}
                        </h2>
                        <div className="flex items-center gap-2">
                            {script ? (
                                <Button
                                    variant="outline"
                                    onClick={() => setViewingScript(true)}
                                    className="gap-1.5"
                                >
                                    {EyeIcon}
                                    View Script
                                </Button>
                            ) : null}
                            {hasSplits ? (
                                <>
                                    <Button
                                        variant="outline"
                                        onClick={() => setConfirmClear(true)}
                                        disabled={clearMutation.isPending}
                                        className="gap-1.5"
                                    >
                                        {TrashIcon}
                                        Clear Splits
                                    </Button>
                                    <Button
                                        onClick={() => generateSplitsTask(scriptId, scriptName)}
                                        disabled={splitInProgress}
                                        className="gap-1.5"
                                    >
                                        {splitInProgress ? (
                                            <span className="h-4 w-4 animate-spin rounded-full border-2 border-background/30 border-t-background" />
                                        ) : (
                                            RefreshIcon
                                        )}
                                        {splitInProgress ? "Generating…" : "Regenerate"}
                                    </Button>
                                </>
                            ) : (
                                <Button
                                    onClick={() => generateSplitsTask(scriptId, scriptName)}
                                    disabled={isLoading || splitInProgress}
                                    className="gap-1.5"
                                >
                                    {splitInProgress ? (
                                        <span className="h-4 w-4 animate-spin rounded-full border-2 border-background/30 border-t-background" />
                                    ) : (
                                        SparklesIcon
                                    )}
                                    {splitInProgress ? "Generating…" : "Generate Splits"}
                                </Button>
                            )}
                        </div>
                    </div>
                </section>

                {/* ── Splits List ── */}
                <section className="border-t border-border pt-8">
                    <SplitsListContent
                        isLoading={isLoading}
                        fetchError={fetchError}
                        splits={splits}
                        onViewImage={handleViewImage}
                    />
                </section>
            </div>
        </>
    );
}
