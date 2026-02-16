"use client";

import { memo, useState, useCallback, useRef } from "react";
import Link from "next/link";
import { useParams, useRouter } from "next/navigation";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import {
    fetchBooks,
    fetchChapters,
    fetchScripts,
    deleteChapter,
    addChapter,
    generateScript,
    ApiError,
    type Book,
    type Chapter,
    type Script,
} from "@/lib/api";
import { Button } from "@/components/ui/button";
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
import {
    Dialog,
    DialogContent,
    DialogHeader,
    DialogTitle,
    DialogDescription,
    DialogFooter,
} from "@/components/ui/dialog";

// ── Static SVG icons (hoisted — rendering-hoist-jsx) ────────────────────────

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

const UploadIcon = (
    <svg
        width="20"
        height="20"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        strokeWidth="1.5"
        strokeLinecap="round"
        strokeLinejoin="round"
    >
        <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4" />
        <polyline points="17 8 12 3 7 8" />
        <line x1="12" x2="12" y1="3" y2="15" />
    </svg>
);

const ChapterIcon = (
    <svg
        className="mx-auto mb-3 h-10 w-10 text-neutral-300"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        strokeWidth="1.5"
        strokeLinecap="round"
        strokeLinejoin="round"
    >
        <path d="M2 3h6a4 4 0 0 1 4 4v14a3 3 0 0 0-3-3H2z" />
        <path d="M22 3h-6a4 4 0 0 0-4 4v14a3 3 0 0 1 3-3h7z" />
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

const ScriptIcon = (
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
        <path d="M15 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V7Z" />
        <path d="M14 2v4a2 2 0 0 0 2 2h4" />
        <path d="M10 13H8" />
        <path d="M16 17H8" />
        <path d="M16 13h-2" />
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

// ── Gallery chapter card (full blurURL background, overlaid info) ────────────

const ChapterCard = memo(function ChapterCard({
    chapter,
    bookId,
    onDelete,
    disabled,
}: {
    chapter: Chapter;
    bookId: number;
    onDelete: (chapter: Chapter) => void;
    disabled: boolean;
}) {
    const hasCover = Boolean(chapter.blurURL);

    return (
        <li className="group relative aspect-[2/3] overflow-hidden rounded-[6px_8px_7px_5px] border border-border shadow-[3px_5px_14px_rgba(0,0,0,0.06)] transition-all duration-300 hover:shadow-[4px_7px_24px_rgba(0,0,0,0.15)] hover:-translate-y-0.5 animate-in fade-in-0 zoom-in-95">
            <Link
                href={`/books/${bookId}/chapters/${chapter.id}`}
                className="absolute inset-0 z-10"
                aria-label={`View chapter ${chapter.number}`}
            />
            {/* ── Full background image ── */}
            {hasCover ? (
                <>
                    {/* eslint-disable-next-line @next/next/no-img-element */}
                    <img
                        src={chapter.blurURL}
                        alt={`Chapter ${chapter.number} cover`}
                        className="absolute inset-0 h-full w-full object-cover transition-transform duration-500 group-hover:scale-105"
                    />
                    {/* Bottom scrim for readability */}
                    <div className="absolute inset-0 bg-black/0 transition-colors duration-300 group-hover:bg-black/10" />
                    <div className="absolute inset-x-0 bottom-0 h-2/3 bg-[linear-gradient(to_top,rgba(0,0,0,0.75)_0%,rgba(0,0,0,0.3)_60%,transparent_100%)]" />
                </>
            ) : (
                <div className="absolute inset-0 flex items-center justify-center bg-neutral-50">
                    <span className="text-[5rem] font-black leading-none tracking-tighter text-neutral-100">
                        {String(chapter.number).padStart(2, "0")}
                    </span>
                </div>
            )}

            {/* ── Overlaid content ── */}
            <div className="relative flex h-full flex-col justify-between p-3">
                {/* Top row: ID badge + delete */}
                <div className="flex items-start justify-between">
                    <span
                        className={`rounded-[3px_5px_4px_3px] px-1.5 py-0.5 font-mono text-[10px] backdrop-blur-sm ${hasCover
                            ? "bg-white/20 text-white/80"
                            : "bg-foreground/5 text-muted-foreground"
                            }`}
                    >
                        ID {chapter.id}
                    </span>
                    <Button
                        variant="ghost"
                        size="icon-sm"
                        onClick={(e) => {
                            e.preventDefault();
                            e.stopPropagation();
                            onDelete(chapter);
                        }}
                        disabled={disabled}
                        aria-label={`Delete chapter ${chapter.number}`}
                        className={`relative z-20 h-7 w-7 rounded-full opacity-0 backdrop-blur-sm transition-all group-hover:opacity-100 ${hasCover
                            ? "bg-white/15 text-white hover:bg-white/30"
                            : "hover:bg-neutral-200"
                            }`}
                    >
                        {TrashIcon}
                    </Button>
                </div>

                {/* Bottom row: large chapter number + label */}
                <div>
                    <span
                        className={`block text-[2.5rem] font-black leading-none tracking-tighter sm:text-5xl ${hasCover ? "text-white" : "text-foreground"
                            }`}
                    >
                        {String(chapter.number).padStart(2, "0")}
                    </span>
                    <span
                        className={`mt-1 block text-[10px] font-medium uppercase tracking-[0.25em] ${hasCover ? "text-white/60" : "text-muted-foreground"
                            }`}
                    >
                        Chapter
                    </span>
                </div>
            </div>
        </li>
    );
});

// ── Chapters grid content (rerender-memo) ─────────────────────────────────────

const ChaptersListContent = memo(function ChaptersListContent({
    isLoading,
    fetchError,
    chapters,
    bookId,
    onDelete,
    deleteDisabled,
}: {
    isLoading: boolean;
    fetchError: Error | null;
    chapters: Chapter[];
    bookId: number;
    onDelete: (chapter: Chapter) => void;
    deleteDisabled: boolean;
}) {
    if (isLoading) {
        return (
            <div className="flex items-center justify-center gap-2 py-8 text-sm text-muted-foreground">
                <span className="h-4 w-4 animate-spin rounded-full border-2 border-border border-t-foreground" />
                Loading chapters…
            </div>
        );
    }

    if (fetchError) {
        return (
            <div className="py-8 text-center text-sm font-medium text-foreground">
                Failed to load chapters —{" "}
                {fetchError instanceof ApiError ? fetchError.message : "network error"}
            </div>
        );
    }

    return chapters.length === 0 ? (
        <div className="py-12 text-center text-muted-foreground">
            {ChapterIcon}
            <p className="text-sm font-medium text-foreground">No chapters yet.</p>
            <span className="text-xs">
                Upload a CBR/CBZ file above to add the first chapter.
            </span>
        </div>
    ) : (
        <ul className="grid grid-cols-2 gap-4 sm:grid-cols-3">
            {chapters.map((chapter) => (
                <ChapterCard
                    key={chapter.id}
                    chapter={chapter}
                    bookId={bookId}
                    onDelete={onDelete}
                    disabled={deleteDisabled}
                />
            ))}
        </ul>
    );
});

// ── Main page ───────────────────────────────────────────────────────────────

export default function BookDetailPage() {
    const params = useParams();
    const bookId = Number(params.id);
    const queryClient = useQueryClient();
    const router = useRouter();
    const fileInputRef = useRef<HTMLInputElement>(null);

    // -- state
    const [chapterNumber, setChapterNumber] = useState("");
    const [selectedFile, setSelectedFile] = useState<File | null>(null);
    const [confirmChapter, setConfirmChapter] = useState<Chapter | null>(null);

    // -- script generation shortcut state
    const [scriptModalOpen, setScriptModalOpen] = useState(false);
    const [scriptName, setScriptName] = useState("");
    const [selectedChapters, setSelectedChapters] = useState<number[]>([]);
    const [selectedPrevScripts, setSelectedPrevScripts] = useState<number[]>([]);

    // -- stable callbacks (rerender-functional-setstate)
    const handleOpenChange = useCallback(
        (open: boolean) => { if (!open) setConfirmChapter(null); },
        []
    );
    const handleDeleteClick = useCallback(
        (chapter: Chapter) => setConfirmChapter(chapter),
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

    // -- fetch chapters
    const {
        data: chaptersData,
        isLoading,
        error: fetchError,
    } = useQuery({
        queryKey: ["chapters", bookId],
        queryFn: () => fetchChapters(bookId),
        enabled: !isNaN(bookId) && bookId > 0,
    });

    const chapters = chaptersData?.data?.chapters ?? [];

    // -- fetch scripts (for "previous scripts" in generate modal)
    const { data: scriptsData } = useQuery({
        queryKey: ["scripts", bookId],
        queryFn: () => fetchScripts(bookId),
        enabled: scriptModalOpen && !isNaN(bookId) && bookId > 0,
    });
    const scripts: Script[] = scriptsData?.data?.scripts ?? [];

    // -- mutations (delete stays local; uploads are global via context)
    const deleteMutation = useMutation({
        mutationFn: (chapterId: number) => deleteChapter(chapterId),
        onSuccess: () => {
            setConfirmChapter(null);
            queryClient.invalidateQueries({ queryKey: ["chapters", bookId] });
            toast.success("Chapter deleted");
        },
        onError: (err) => {
            setConfirmChapter(null);
            toast.error(
                err instanceof ApiError ? err.businessError : "Delete failed"
            );
        },
    });

    const handleSubmit = useCallback(
        (e: React.FormEvent) => {
            e.preventDefault();
            const num = parseInt(chapterNumber, 10);
            if (isNaN(num) || num <= 0 || !selectedFile) return;
            // Fire-and-forget — events tracker polls for status
            toast.info(`Adding Ch.${num}…`, { description: selectedFile.name, duration: 3000 });
            addChapter(bookId, num, selectedFile).catch((err) => {
                toast.error(
                    err instanceof ApiError ? err.businessError : "Upload failed — please retry"
                );
            });
            // Immediately reset form
            setChapterNumber("");
            setSelectedFile(null);
            if (fileInputRef.current) fileInputRef.current.value = "";
        },
        [chapterNumber, selectedFile, bookId]
    );

    const handleConfirmDelete = useCallback(() => {
        if (confirmChapter) {
            deleteMutation.mutate(confirmChapter.id);
        }
    }, [confirmChapter, deleteMutation]);

    // -- script generation handlers
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

    const handleGenerateScript = useCallback(
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
            setScriptModalOpen(false);
            setScriptName("");
            setSelectedChapters([]);
            setSelectedPrevScripts([]);
            // Navigate to scripts page so user can see result when ready
            router.push(`/books/${bookId}/scripts`);
        },
        [scriptName, selectedChapters, selectedPrevScripts, bookId, router]
    );

    // -- derived state
    const isSubmitDisabled = !chapterNumber.trim() || !selectedFile;
    const isGenerateDisabled = !scriptName.trim() || selectedChapters.length === 0;

    const bookTitle = book?.title ?? `Book #${bookId}`;

    return (
        <>
            {/* ── Generate Script Modal ── */}
            <Dialog open={scriptModalOpen} onOpenChange={setScriptModalOpen}>
                <DialogContent>
                    <DialogHeader>
                        <DialogTitle>Generate New Script</DialogTitle>
                        <DialogDescription>
                            Select chapters and optionally reference previous scripts for
                            continuity.
                        </DialogDescription>
                    </DialogHeader>

                    <form onSubmit={handleGenerateScript} className="space-y-5">
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
                                Generate Script
                            </Button>
                        </DialogFooter>
                    </form>
                </DialogContent>
            </Dialog>

        <div className="mx-auto max-w-5xl px-6 pb-20 max-sm:px-4">
            {/* ── Delete confirmation ─────────────────────────────────── */}
            <AlertDialog
                open={confirmChapter !== null}
                onOpenChange={handleOpenChange}
            >
                <AlertDialogContent>
                    <AlertDialogHeader>
                        <AlertDialogTitle>Delete chapter?</AlertDialogTitle>
                        <AlertDialogDescription>
                            Chapter {confirmChapter?.number} and all its pages and panels will
                            be permanently removed.
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

            {/* ── Hero ──────────────────────────────────────────────────── */}
            <header className="relative pb-10 pt-20 max-sm:pb-7 max-sm:pt-12">
                <Link
                    href="/"
                    className="mb-6 inline-flex items-center gap-1.5 text-sm text-muted-foreground transition-colors hover:text-foreground"
                >
                    {ArrowLeftIcon}
                    Back to books
                </Link>
                <h1 className="text-[clamp(2.5rem,8vw,4.5rem)] font-black leading-[0.85] tracking-tighter max-sm:text-4xl">
                    {bookTitle}
                </h1>
                {HeroUnderline}
                <span className="mt-4 block text-[11px] font-medium uppercase tracking-[0.3em] text-muted-foreground">
                    CHAPTER MANAGEMENT
                </span>
            </header>

            {/* ── Add Chapter — clean upload area ────────────────────────── */}
            <section className="border-t border-border pb-8 pt-10">
                <h2 className="mb-5 text-2xl font-semibold tracking-tight">
                    Add New Chapter
                </h2>
                <form onSubmit={handleSubmit} className="space-y-4">
                    {/* Chapter number */}
                    <div className="max-w-[8rem]">
                        <label
                            htmlFor="chapter-number-input"
                            className="mb-1.5 block text-xs font-medium uppercase tracking-widest text-muted-foreground"
                        >
                            Chapter №
                        </label>
                        <Input
                            id="chapter-number-input"
                            type="number"
                            min="1"
                            placeholder="1"
                            value={chapterNumber}
                            onChange={(e) => setChapterNumber(e.target.value)}
                        />
                    </div>

                    {/* File drop zone */}
                    <div>
                        <label className="mb-1.5 block text-xs font-medium uppercase tracking-widest text-muted-foreground">
                            Chapter file
                        </label>
                        <input
                            ref={fileInputRef}
                            id="chapter-file-input"
                            type="file"
                            accept=".cbr,.cbz"
                            className="hidden"
                            onChange={(e) => {
                                setSelectedFile(e.target.files?.[0] ?? null);
                            }}
                        />
                        <button
                            type="button"
                            onClick={() => fileInputRef.current?.click()}
                            className={`group/drop flex w-full cursor-pointer flex-col items-center justify-center gap-2 rounded-lg border-2 border-dashed px-6 py-8 text-center transition-colors ${selectedFile
                                ? "border-foreground bg-foreground/[0.03]"
                                : "border-border hover:border-foreground/40 hover:bg-foreground/[0.02]"
                                }`}
                        >
                            <span
                                className={`transition-colors ${selectedFile
                                    ? "text-foreground"
                                    : "text-neutral-300 group-hover/drop:text-neutral-400"
                                    }`}
                            >
                                {UploadIcon}
                            </span>
                            {selectedFile ? (
                                <>
                                    <span className="text-sm font-medium">{selectedFile.name}</span>
                                    <span className="text-xs text-muted-foreground">
                                        {(selectedFile.size / (1024 * 1024)).toFixed(1)} MB — click
                                        to change
                                    </span>
                                </>
                            ) : (
                                <>
                                    <span className="text-sm font-medium text-muted-foreground">
                                        Click to select a file
                                    </span>
                                    <span className="text-xs text-muted-foreground/60">
                                        CBR or CBZ formats accepted
                                    </span>
                                </>
                            )}
                        </button>
                    </div>

                    <Button
                        id="add-chapter-btn"
                        type="submit"
                        disabled={isSubmitDisabled}
                        className="max-sm:w-full"
                    >
                        {"Add Chapter"}
                    </Button>
                </form>
            </section>

            {/* ── Chapters Gallery ──────────────────────────────────────── */}
            <section className="border-t border-border pt-8">
                <div className="mb-6 flex flex-wrap items-center justify-between gap-3">
                    <h2 className="flex items-center gap-2.5 text-2xl font-semibold tracking-tight">
                        Chapters
                        {chapters.length > 0 ? (
                            <span className="rounded-[4px_6px_5px_3px] bg-foreground px-2 py-0.5 text-xs font-medium text-background">
                                {chapters.length}
                            </span>
                        ) : null}
                    </h2>
                    {chapters.length > 0 ? (
                        <Button
                            onClick={() => setScriptModalOpen(true)}
                            className="gap-1.5"
                        >
                            {ScriptIcon}
                            Create Script
                        </Button>
                    ) : null}
                </div>

                <ChaptersListContent
                    isLoading={isLoading}
                    fetchError={fetchError}
                    chapters={chapters}
                    bookId={bookId}
                    onDelete={handleDeleteClick}
                    deleteDisabled={deleteMutation.isPending}
                />
            </section>
        </div>
        </>
    );
}
