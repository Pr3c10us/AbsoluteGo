"use client";

import { memo, useState, useCallback, useMemo, useRef, useEffect } from "react";
import Link from "next/link";
import { useParams } from "next/navigation";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import { ArrowLeft, Film, Play, Trash2, AlertTriangle, Search, X } from "lucide-react";
import {
    fetchBooks,
    fetchScripts,
    fetchVABs,
    deleteVAB,
    ApiError,
    type Book,
    type Script,
    type VAB,
} from "@/lib/api";
import { Input } from "@/components/ui/input";
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
import VideoPlayerOverlayWrapper from "@/components/videoPlayerOverlay";

// ── Static icons ─────────────────────────────────────────────────────────────

const ArrowLeftIcon = <ArrowLeft className="h-4 w-4" />;

const AlertTriangleIcon = <AlertTriangle className="h-4 w-4" />;

const SearchIcon = <Search className="h-3.5 w-3.5" />;

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

// ── Debounce hook ─────────────────────────────────────────────────────────────

function useDebounce<T>(value: T, delay: number): T {
    const [debounced, setDebounced] = useState(value);
    useEffect(() => {
        const timer = setTimeout(() => setDebounced(value), delay);
        return () => clearTimeout(timer);
    }, [value, delay]);
    return debounced;
}

// ── Video thumbnail card ─────────────────────────────────────────────────────

const VABThumbnailCard = memo(function VABThumbnailCard({
    vab,
    scriptName,
    bookId,
    onPlay,
    onDelete,
    deleteDisabled,
}: {
    vab: VAB;
    scriptName: string | undefined;
    bookId: number;
    onPlay: (vab: VAB) => void;
    onDelete: (vab: VAB) => void;
    deleteDisabled: boolean;
}) {
    const videoRef = useRef<HTMLVideoElement>(null);
    const [hovered, setHovered] = useState(false);

    const handleMouseEnter = useCallback(() => {
        setHovered(true);
        // Short preview on hover — silent, muted
        videoRef.current?.play().catch(() => {});
    }, []);

    const handleMouseLeave = useCallback(() => {
        setHovered(false);
        if (videoRef.current) {
            videoRef.current.pause();
            videoRef.current.currentTime = 0;
        }
    }, []);

    return (
        <li className="group animate-in fade-in-0 zoom-in-95">
            {/* ── Thumbnail ── */}
            <div
                className="relative mb-3 aspect-video w-full cursor-pointer overflow-hidden rounded-[6px_8px_7px_5px] bg-black"
                onClick={() => vab.Url && onPlay(vab)}
                onMouseEnter={vab.Url ? handleMouseEnter : undefined}
                onMouseLeave={vab.Url ? handleMouseLeave : undefined}
            >
                {vab.Url ? (
                    <>
                        {/* Native video element — first frame = thumbnail */}
                        <video
                            ref={videoRef}
                            src={vab.Url}
                            muted
                            playsInline
                            preload="metadata"
                            className={`h-full w-full object-cover transition-all duration-500 ${hovered ? "scale-105" : "scale-100"}`}
                        />

                        {/* Play button overlay */}
                        <div
                            className={`absolute inset-0 flex items-center justify-center bg-black/30 transition-opacity duration-200 ${hovered ? "opacity-100" : "opacity-0"}`}
                        >
                            <div className="flex h-14 w-14 items-center justify-center rounded-full bg-white shadow-[0_4px_24px_rgba(0,0,0,0.4)]">
                                <Play
                                    className="h-6 w-6 translate-x-0.5 text-black"
                                    fill="currentColor"
                                    strokeWidth={0}
                                />
                            </div>
                        </div>

                        {/* Badge — bottom-right corner */}
                        <div className="absolute bottom-2 right-2 rounded-[3px_5px_4px_3px] bg-black/80 px-1.5 py-0.5 font-mono text-[10px] font-medium text-white backdrop-blur-sm">
                            VAB
                        </div>
                    </>
                ) : (
                    /* Processing state */
                    <div className="flex h-full w-full flex-col items-center justify-center gap-2">
                        <div className="h-6 w-6 animate-spin rounded-full border-2 border-white/20 border-t-white/60" />
                        <span className="text-[11px] font-medium text-white/40">
                            Processing…
                        </span>
                    </div>
                )}
            </div>

            {/* ── Metadata ── */}
            <div className="flex items-start justify-between gap-2 px-0.5 pt-2">
                <div className="min-w-0">
                    <h3
                        className="truncate text-base font-black tracking-tight text-foreground"
                        title={vab.Name}
                    >
                        {vab.Name}
                    </h3>
                    <div className="mt-1.5 flex flex-wrap items-center gap-1.5">
                        {scriptName ? (
                            <Link
                                href={`/books/${bookId}/scripts/${vab.ScriptId}/splits`}
                                onClick={(e) => e.stopPropagation()}
                                className="rounded-[3px_5px_4px_3px] bg-foreground px-1.5 py-0.5 text-[10px] font-medium uppercase tracking-wider text-background transition-opacity hover:opacity-70"
                            >
                                {scriptName}
                            </Link>
                        ) : null}
                        <span className="font-mono text-[10px] text-neutral-400">
                            #{vab.Id}
                        </span>
                    </div>
                </div>
                <button
                    onClick={(e) => {
                        e.stopPropagation();
                        onDelete(vab);
                    }}
                    disabled={deleteDisabled}
                    aria-label={`Delete ${vab.Name}`}
                    className="mt-0.5 shrink-0 cursor-pointer text-neutral-300 opacity-0 transition-all duration-150 hover:text-foreground group-hover:opacity-100 disabled:pointer-events-none"
                >
                    <Trash2 className="h-3.5 w-3.5" />
                </button>
            </div>
        </li>
    );
});

// ── VABs grid ────────────────────────────────────────────────────────────────

const VABsGrid = memo(function VABsGrid({
    isLoading,
    fetchError,
    vabs,
    scripts,
    bookId,
    onPlay,
    onDelete,
    deleteDisabled,
    searchActive,
}: {
    isLoading: boolean;
    fetchError: Error | null;
    vabs: VAB[];
    scripts: Script[];
    bookId: number;
    onPlay: (vab: VAB) => void;
    onDelete: (vab: VAB) => void;
    deleteDisabled: boolean;
    searchActive: boolean;
}) {
    if (isLoading) {
        return (
            <div className="flex items-center justify-center gap-2 py-12 text-sm text-muted-foreground">
                <span className="h-4 w-4 animate-spin rounded-full border-2 border-border border-t-foreground" />
                {searchActive ? "Searching…" : "Loading video audiobooks…"}
            </div>
        );
    }

    if (fetchError) {
        return (
            <div className="py-8 text-center text-sm font-medium text-foreground">
                Failed to load video audiobooks —{" "}
                {fetchError instanceof ApiError
                    ? fetchError.message
                    : "network error"}
            </div>
        );
    }

    if (vabs.length === 0) {
        return searchActive ? (
            <div className="py-16 text-center text-muted-foreground">
                <Search
                    className="mx-auto mb-3 h-10 w-10 text-neutral-300"
                    strokeWidth={1.5}
                />
                <p className="text-sm font-medium text-foreground">
                    No audiobooks match your search.
                </p>
                <span className="text-xs">Try a different name.</span>
            </div>
        ) : (
            <div className="py-16 text-center text-muted-foreground">
                <Film
                    className="mx-auto mb-3 h-10 w-10 text-neutral-300"
                    strokeWidth={1.5}
                />
                <p className="text-sm font-medium text-foreground">
                    No video audiobooks yet.
                </p>
                <span className="text-xs">
                    Go to a script&apos;s splits page and click{" "}
                    <strong>Video Audiobook</strong> to create one.
                </span>
            </div>
        );
    }

    return (
        <ul className="grid grid-cols-1 gap-x-6 gap-y-10 sm:grid-cols-2">
            {vabs.map((vab) => {
                const scriptName = scripts.find(
                    (s) => s.id === vab.ScriptId,
                )?.name;
                return (
                    <VABThumbnailCard
                        key={vab.Id}
                        vab={vab}
                        scriptName={scriptName}
                        bookId={bookId}
                        onPlay={onPlay}
                        onDelete={onDelete}
                        deleteDisabled={deleteDisabled}
                    />
                );
            })}
        </ul>
    );
});

// ── Main page ────────────────────────────────────────────────────────────────

export default function BookVideosPage() {
    const params = useParams();
    const bookId = Number(params.id);
    const queryClient = useQueryClient();

    const [confirmVAB, setConfirmVAB] = useState<VAB | null>(null);
    const [videoPlayerState, setVideoPlayerState] = useState<{
        url: string;
        label: string;
    } | null>(null);

    // Search — input value is immediate, debouncedSearch is sent to the backend
    const [search, setSearch] = useState("");
    const debouncedSearch = useDebounce(search, 350);

    const handleConfirmOpenChange = useCallback((open: boolean) => {
        if (!open) setConfirmVAB(null);
    }, []);
    const closeVideoPlayer = useCallback(() => setVideoPlayerState(null), []);
    const handleDeleteClick = useCallback((vab: VAB) => setConfirmVAB(vab), []);
    const clearSearch = useCallback(() => setSearch(""), []);
    const handlePlay = useCallback((vab: VAB) => {
        if (!vab.Url) return;
        setVideoPlayerState({ url: vab.Url, label: vab.Name });
    }, []);

    const { data: booksData } = useQuery({
        queryKey: ["books"],
        queryFn: () => fetchBooks(),
    });
    const book: Book | undefined = booksData?.data?.books?.find(
        (b) => b.id === bookId,
    );

    const { data: scriptsData } = useQuery({
        queryKey: ["scripts", bookId],
        queryFn: () => fetchScripts(bookId),
        enabled: !isNaN(bookId) && bookId > 0,
    });
    const scripts: Script[] = useMemo(
        () => scriptsData?.data?.scripts ?? [],
        [scriptsData],
    );

    // VABs — queryKey includes debouncedSearch so React Query re-fetches when it changes
    const {
        data: vabsData,
        isLoading,
        isFetching,
        error: fetchError,
    } = useQuery({
        queryKey: ["vabs", bookId, debouncedSearch],
        queryFn: () =>
            fetchVABs({
                bookId,
                name: debouncedSearch.trim() || undefined,
            }),
        enabled: !isNaN(bookId) && bookId > 0,
    });
    const vabs: VAB[] = useMemo(
        () => vabsData?.data?.vabs ?? [],
        [vabsData],
    );

    const processingCount = useMemo(
        () => vabs.filter((v) => !v.Url).length,
        [vabs],
    );

    // Whether we're still waiting for the debounced query to resolve
    const searchActive = debouncedSearch.trim().length > 0;
    const isSearching = search !== debouncedSearch || (searchActive && isFetching);

    const deleteMutation = useMutation({
        mutationFn: (vabId: number) => deleteVAB(vabId),
        onSuccess: () => {
            setConfirmVAB(null);
            queryClient.invalidateQueries({ queryKey: ["vabs", bookId] });
            toast.success("Video audiobook deleted");
        },
        onError: (err) => {
            setConfirmVAB(null);
            toast.error(
                err instanceof ApiError ? err.businessError : "Delete failed",
            );
        },
    });

    const handleConfirmDelete = useCallback(() => {
        if (confirmVAB) deleteMutation.mutate(confirmVAB.Id);
    }, [confirmVAB, deleteMutation]);

    const bookTitle = book?.title ?? `Book #${bookId}`;

    return (
        <>
            {videoPlayerState ? (
                <VideoPlayerOverlayWrapper
                    url={videoPlayerState.url}
                    label={videoPlayerState.label}
                    onClose={closeVideoPlayer}
                />
            ) : null}

            <AlertDialog
                open={confirmVAB !== null}
                onOpenChange={handleConfirmOpenChange}
            >
                <AlertDialogContent>
                    <AlertDialogHeader>
                        <AlertDialogTitle>
                            Delete video audiobook?
                        </AlertDialogTitle>
                        <AlertDialogDescription>
                            &ldquo;{confirmVAB?.Name}&rdquo; will be permanently
                            removed. This cannot be undone.
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
                        Video Audiobooks
                    </h1>
                    {HeroUnderline}
                    <span className="mt-4 block text-[11px] font-medium uppercase tracking-[0.3em] text-muted-foreground">
                        GENERATED PRODUCTIONS
                    </span>
                </header>

                {/* ── Section header ── */}
                <section className="border-t border-border pt-10">
                    <div className="mb-8 flex flex-wrap items-center justify-between gap-3">
                        <h2 className="flex items-center gap-2.5 text-2xl font-semibold tracking-tight">
                            Audiobooks
                            {vabs.length > 0 ? (
                                <span className="rounded-[4px_6px_5px_3px] bg-foreground px-2 py-0.5 text-xs font-medium text-background">
                                    {vabs.length}
                                </span>
                            ) : null}
                        </h2>

                        {/* Search — always shown so users can search even before knowing count */}
                        <div className="relative w-full sm:w-56">
                            <span className="pointer-events-none absolute left-2.5 top-1/2 -translate-y-1/2 text-muted-foreground">
                                {isSearching ? (
                                    <span className="h-3.5 w-3.5 animate-spin rounded-full border-2 border-border border-t-foreground inline-block" />
                                ) : (
                                    SearchIcon
                                )}
                            </span>
                            <Input
                                value={search}
                                onChange={(e) => setSearch(e.target.value)}
                                placeholder="Search audiobooks…"
                                className="h-8 pl-7 pr-7 text-sm"
                            />
                            {search ? (
                                <button
                                    onClick={clearSearch}
                                    aria-label="Clear search"
                                    className="absolute right-2 top-1/2 -translate-y-1/2 cursor-pointer text-muted-foreground transition-colors hover:text-foreground"
                                >
                                    <X className="h-3.5 w-3.5" />
                                </button>
                            ) : null}
                        </div>
                    </div>

                    {/* Processing notice */}
                    {processingCount > 0 ? (
                        <div className="mb-8 flex items-start gap-2.5 rounded-[4px_6px_5px_3px] border border-neutral-200 bg-neutral-50 px-4 py-3">
                            {AlertTriangleIcon}
                            <div className="min-w-0">
                                <p className="text-sm font-medium text-foreground">
                                    {processingCount} audiobook
                                    {processingCount !== 1 ? "s" : ""}{" "}
                                    processing
                                </p>
                                <p className="mt-0.5 text-xs text-muted-foreground">
                                    Video merging is in progress. Check the
                                    event tracker for live status updates.
                                </p>
                            </div>
                        </div>
                    ) : null}

                    <VABsGrid
                        isLoading={isLoading}
                        fetchError={fetchError as Error | null}
                        vabs={vabs}
                        scripts={scripts}
                        bookId={bookId}
                        onPlay={handlePlay}
                        onDelete={handleDeleteClick}
                        deleteDisabled={deleteMutation.isPending}
                        searchActive={searchActive}
                    />
                </section>
            </div>
        </>
    );
}
