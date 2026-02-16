"use client";

import { memo, useState, useCallback, useMemo } from "react";
import Link from "next/link";
import { useParams } from "next/navigation";
import { useQuery, useQueries } from "@tanstack/react-query";
import { ArrowLeft, Maximize2, LayoutGrid, X, FileImage, Layers, File } from "lucide-react";
import {
    fetchBooks,
    fetchChapters,
    fetchPages,
    fetchPanels,
    ApiError,
    type Book,
    type Chapter,
    type Page,
    type Panel,
} from "@/lib/api";
import { Button } from "@/components/ui/button";
import { useScrollLock } from "@/lib/use-scroll-lock";
import Lightbox, { type LightboxItem } from "@/components/lightbox";

// ── Static icons (hoisted — rendering-hoist-jsx) ────────────────────────────

const ArrowLeftIcon = <ArrowLeft className="h-4 w-4" />;

const ExpandIcon = <Maximize2 className="h-3.5 w-3.5" />;

const GridIcon = <LayoutGrid className="h-3.5 w-3.5" />;

const CloseIcon = <X className="h-4 w-4" />;

const PagesEmptyIcon = <FileImage className="mx-auto mb-3 h-10 w-10 text-neutral-300" strokeWidth={1.5} />;

const LayersIcon = <Layers className="h-3.5 w-3.5" />;

const FileIcon = <File className="h-3.5 w-3.5" />;

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

const PanelsEmptyIcon = <Layers className="mx-auto mb-3 h-10 w-10 text-neutral-300" strokeWidth={1.5} />;

// ── Panel viewer overlay (with its own lightbox) ────────────────────────────

const PanelViewer = memo(function PanelViewer({
    page,
    onClose,
}: {
    page: Page;
    onClose: () => void;
}) {
    useScrollLock();

    const { data, isLoading, error } = useQuery({
        queryKey: ["panels", page.id],
        queryFn: () => fetchPanels(page.id),
    });

    const panels: Panel[] = data?.data?.panels ?? [];

    // Lightbox for panels (rerender-functional-setstate)
    const [panelLightboxIdx, setPanelLightboxIdx] = useState<number | null>(null);
    const closePanelLightbox = useCallback(() => setPanelLightboxIdx(null), []);

    const panelLightboxItems: LightboxItem[] = useMemo(
        () =>
            panels.map((p) => ({
                url: p.url,
                label: `Panel ${p.panelNumber}`,
            })),
        [panels]
    );

    return (
        <>
            {panelLightboxIdx !== null ? (
                <Lightbox
                    items={panelLightboxItems}
                    currentIndex={panelLightboxIdx}
                    onIndexChange={setPanelLightboxIdx}
                    onClose={closePanelLightbox}
                />
            ) : null}

            <div className="fixed inset-0 z-40 flex flex-col bg-white">
                {/* ── Header ── */}
                <div className="flex items-center justify-between border-b border-border px-6 py-3">
                    <div>
                        <h2 className="text-lg font-bold tracking-tight">
                            Page {page.pageNumber} — Panels
                        </h2>
                        <span className="text-xs text-muted-foreground">
                            {panels.length > 0
                                ? `${panels.length} panel${panels.length > 1 ? "s" : ""} segmented`
                                : isLoading
                                    ? "Loading…"
                                    : "No panels"}
                        </span>
                    </div>
                    <Button
                        variant="outline"
                        size="icon-sm"
                        onClick={onClose}
                        aria-label="Close panel viewer"
                    >
                        {CloseIcon}
                    </Button>
                </div>

                {/* ── Content ── */}
                <div className="flex-1 overflow-y-auto p-6 max-sm:p-4">
                    {isLoading ? (
                        <div className="flex items-center justify-center gap-2 py-16 text-sm text-muted-foreground">
                            <span className="h-4 w-4 animate-spin rounded-full border-2 border-border border-t-foreground" />
                            Loading panels…
                        </div>
                    ) : error ? (
                        <div className="py-16 text-center text-sm font-medium text-foreground">
                            Failed to load panels —{" "}
                            {error instanceof ApiError ? error.message : "network error"}
                        </div>
                    ) : panels.length === 0 ? (
                        <div className="py-16 text-center text-muted-foreground">
                            <p className="text-sm font-medium text-foreground">
                                No panels found.
                            </p>
                            <span className="text-xs">
                                This page has not been segmented yet.
                            </span>
                        </div>
                    ) : (
                        <div className="mx-auto grid max-w-5xl grid-cols-2 gap-4 sm:grid-cols-3 lg:grid-cols-4">
                            {panels.map((panel, idx) => (
                                <div
                                    key={panel.id}
                                    className="group relative overflow-hidden rounded-[5px_7px_6px_4px] border border-border shadow-[2px_4px_12px_rgba(0,0,0,0.04)] transition-all duration-300 hover:shadow-[3px_5px_18px_rgba(0,0,0,0.1)] hover:-translate-y-0.5"
                                >
                                    {/* eslint-disable-next-line @next/next/no-img-element */}
                                    <img
                                        src={panel.url}
                                        alt={`Panel ${panel.panelNumber}`}
                                        className="block w-full transition-transform duration-500 group-hover:scale-[1.02]"
                                    />
                                    {/* Hover overlay with View Big CTA */}
                                    <div className="absolute inset-0 flex items-center justify-center bg-black/0 opacity-0 transition-all duration-300 group-hover:bg-black/40 group-hover:opacity-100">
                                        <Button
                                            variant="secondary"
                                            size="sm"
                                            onClick={() => setPanelLightboxIdx(idx)}
                                            className="gap-1.5 bg-white text-black shadow-lg hover:bg-neutral-100"
                                        >
                                            {ExpandIcon}
                                            View Big
                                        </Button>
                                    </div>
                                    {/* Panel number badge */}
                                    <div className="absolute bottom-2 left-2 rounded-[3px_5px_4px_3px] bg-black/70 px-2 py-0.5 font-mono text-[10px] font-medium text-white backdrop-blur-sm">
                                        P{panel.panelNumber}
                                    </div>
                                </div>
                            ))}
                        </div>
                    )}
                </div>
            </div>
        </>
    );
});

// ── Page card ───────────────────────────────────────────────────────────────

const PageCard = memo(function PageCard({
    page,
    index,
    onViewPanels,
    onViewBig,
}: {
    page: Page;
    index: number;
    onViewPanels: (page: Page) => void;
    onViewBig: (index: number) => void;
}) {
    return (
        <li className="group relative overflow-hidden rounded-[6px_8px_7px_5px] border border-border shadow-[3px_5px_14px_rgba(0,0,0,0.06)] transition-all duration-300 hover:shadow-[4px_7px_24px_rgba(0,0,0,0.15)] hover:-translate-y-0.5 animate-in fade-in-0 zoom-in-95">
            {/* ── Page image ── */}
            <div className="relative aspect-[2/3] overflow-hidden bg-neutral-50">
                {/* eslint-disable-next-line @next/next/no-img-element */}
                <img
                    src={page.url}
                    alt={`Page ${page.pageNumber}`}
                    className="absolute inset-0 h-full w-full object-cover transition-transform duration-500 group-hover:scale-105"
                />
                {/* Hover overlay with CTAs */}
                <div className="absolute inset-0 flex flex-col items-center justify-center gap-2 bg-black/0 opacity-0 transition-all duration-300 group-hover:bg-black/40 group-hover:opacity-100">
                    <Button
                        variant="secondary"
                        size="sm"
                        onClick={() => onViewBig(index)}
                        className="gap-1.5 bg-white text-black shadow-lg hover:bg-neutral-100"
                    >
                        {ExpandIcon}
                        View Big
                    </Button>
                    <Button
                        variant="secondary"
                        size="sm"
                        onClick={() => onViewPanels(page)}
                        className="gap-1.5 bg-white text-black shadow-lg hover:bg-neutral-100"
                    >
                        {GridIcon}
                        View Panels
                    </Button>
                </div>
            </div>

            {/* ── Page info footer ── */}
            <div className="flex items-center justify-between px-3 py-2.5">
                <div>
                    <span className="block text-xl font-black leading-none tracking-tighter">
                        {String(page.pageNumber).padStart(2, "0")}
                    </span>
                    <span className="block text-[9px] font-medium uppercase tracking-[0.25em] text-muted-foreground">
                        Page
                    </span>
                </div>
                <span className="rounded-[3px_5px_4px_3px] bg-foreground/5 px-1.5 py-0.5 font-mono text-[10px] text-muted-foreground">
                    ID {page.id}
                </span>
            </div>
        </li>
    );
});

// ── Pages list content (rerender-memo) ──────────────────────────────────────

const PagesListContent = memo(function PagesListContent({
    isLoading,
    fetchError,
    pages,
    onViewPanels,
    onViewBig,
}: {
    isLoading: boolean;
    fetchError: Error | null;
    pages: Page[];
    onViewPanels: (page: Page) => void;
    onViewBig: (index: number) => void;
}) {
    if (isLoading) {
        return (
            <div className="flex items-center justify-center gap-2 py-8 text-sm text-muted-foreground">
                <span className="h-4 w-4 animate-spin rounded-full border-2 border-border border-t-foreground" />
                Loading pages…
            </div>
        );
    }

    if (fetchError) {
        return (
            <div className="py-8 text-center text-sm font-medium text-foreground">
                Failed to load pages —{" "}
                {fetchError instanceof ApiError ? fetchError.message : "network error"}
            </div>
        );
    }

    return pages.length === 0 ? (
        <div className="py-12 text-center text-muted-foreground">
            {PagesEmptyIcon}
            <p className="text-sm font-medium text-foreground">No pages found.</p>
            <span className="text-xs">
                This chapter doesn&apos;t have any pages yet.
            </span>
        </div>
    ) : (
        <ul className="grid grid-cols-2 gap-4 sm:grid-cols-3">
            {pages.map((page, idx) => (
                <PageCard
                    key={page.id}
                    page={page}
                    index={idx}
                    onViewPanels={onViewPanels}
                    onViewBig={onViewBig}
                />
            ))}
        </ul>
    );
});

// ── Panel card (for all-panels view) ─────────────────────────────────────────

interface PanelWithPage extends Panel {
    pageNumber: number;
}

const AllPanelCard = memo(function AllPanelCard({
    panel,
    index,
    onViewBig,
}: {
    panel: PanelWithPage;
    index: number;
    onViewBig: (index: number) => void;
}) {
    return (
        <li className="group relative overflow-hidden rounded-[6px_8px_7px_5px] border border-border shadow-[3px_5px_14px_rgba(0,0,0,0.06)] transition-all duration-300 hover:shadow-[4px_7px_24px_rgba(0,0,0,0.15)] hover:-translate-y-0.5 animate-in fade-in-0 zoom-in-95">
            <div className="relative overflow-hidden bg-neutral-50">
                {/* eslint-disable-next-line @next/next/no-img-element */}
                <img
                    src={panel.url}
                    alt={`Panel ${panel.panelNumber}`}
                    className="block w-full transition-transform duration-500 group-hover:scale-[1.02]"
                />
                {/* Hover overlay */}
                <div className="absolute inset-0 flex items-center justify-center bg-black/0 opacity-0 transition-all duration-300 group-hover:bg-black/40 group-hover:opacity-100">
                    <Button
                        variant="secondary"
                        size="sm"
                        onClick={() => onViewBig(index)}
                        className="gap-1.5 bg-white text-black shadow-lg hover:bg-neutral-100"
                    >
                        {ExpandIcon}
                        View
                    </Button>
                </div>
            </div>
            {/* Info footer */}
            <div className="flex items-center justify-between px-3 py-2">
                <div>
                    <span className="block text-sm font-bold leading-none tracking-tight">
                        P{panel.panelNumber}
                    </span>
                    <span className="block text-[9px] font-medium uppercase tracking-[0.25em] text-muted-foreground">
                        Page {String(panel.pageNumber).padStart(2, "0")}
                    </span>
                </div>
                <span className="rounded-[3px_5px_4px_3px] bg-foreground/5 px-1.5 py-0.5 font-mono text-[10px] text-muted-foreground">
                    ID {panel.id}
                </span>
            </div>
        </li>
    );
});

// ── All panels list content ──────────────────────────────────────────────────

const AllPanelsListContent = memo(function AllPanelsListContent({
    pages,
    pagesLoading,
    onViewBig,
}: {
    pages: Page[];
    pagesLoading: boolean;
    onViewBig: (index: number) => void;
}) {
    // Fetch panels for every page in parallel
    const panelQueries = useQueries({
        queries: pages.map((p) => ({
            queryKey: ["panels", p.id] as const,
            queryFn: () => fetchPanels(p.id),
        })),
    });

    const isLoading = pagesLoading || panelQueries.some((q) => q.isLoading);
    const hasError = panelQueries.some((q) => q.error);

    // Aggregate all panels with page context
    // Derive a stable key from query data to avoid variable-length deps
    const panelDataKey = panelQueries.map((q) => q.dataUpdatedAt).join(",");
    const allPanels: PanelWithPage[] = useMemo(
        () =>
            pages.flatMap((page, i) => {
                const panels: Panel[] = panelQueries[i]?.data?.data?.panels ?? [];
                return panels.map((p) => ({ ...p, pageNumber: page.pageNumber }));
            }),
        // eslint-disable-next-line react-hooks/exhaustive-deps
        [pages, panelDataKey]
    );

    if (isLoading) {
        return (
            <div className="flex items-center justify-center gap-2 py-8 text-sm text-muted-foreground">
                <span className="h-4 w-4 animate-spin rounded-full border-2 border-border border-t-foreground" />
                Loading panels…
            </div>
        );
    }

    if (hasError) {
        return (
            <div className="py-8 text-center text-sm font-medium text-foreground">
                Failed to load some panels — please try again.
            </div>
        );
    }

    return allPanels.length === 0 ? (
        <div className="py-12 text-center text-muted-foreground">
            {PanelsEmptyIcon}
            <p className="text-sm font-medium text-foreground">No panels found.</p>
            <span className="text-xs">
                Pages haven&apos;t been segmented into panels yet.
            </span>
        </div>
    ) : (
        <ul className="grid grid-cols-2 gap-4 sm:grid-cols-3 lg:grid-cols-4">
            {allPanels.map((panel, idx) => (
                <AllPanelCard
                    key={panel.id}
                    panel={panel}
                    index={idx}
                    onViewBig={onViewBig}
                />
            ))}
        </ul>
    );
});

// ── Main page ───────────────────────────────────────────────────────────────

export default function ChapterDetailPage() {
    const params = useParams();
    const bookId = Number(params.id);
    const chapterId = Number(params.chapterId);

    // -- view toggle: "pages" | "panels"
    type ViewMode = "pages" | "panels";
    const [viewMode, setViewMode] = useState<ViewMode>("pages");

    // -- panel viewer state (page-specific overlay)
    const [viewingPage, setViewingPage] = useState<Page | null>(null);
    const handleViewPanels = useCallback(
        (page: Page) => setViewingPage(page),
        []
    );
    const closePanelViewer = useCallback(() => setViewingPage(null), []);

    // -- page lightbox state
    const [pageLightboxIdx, setPageLightboxIdx] = useState<number | null>(null);
    const closePageLightbox = useCallback(() => setPageLightboxIdx(null), []);

    // -- all-panels lightbox state
    const [panelLightboxIdx, setPanelLightboxIdx] = useState<number | null>(null);
    const closePanelLightbox = useCallback(() => setPanelLightboxIdx(null), []);

    // -- fetch book info
    const { data: booksData } = useQuery({
        queryKey: ["books"],
        queryFn: () => fetchBooks(),
    });
    const book: Book | undefined = booksData?.data?.books?.find(
        (b) => b.id === bookId
    );

    // -- fetch chapter info
    const { data: chaptersData } = useQuery({
        queryKey: ["chapters", bookId],
        queryFn: () => fetchChapters(bookId),
        enabled: !isNaN(bookId) && bookId > 0,
    });
    const chapter: Chapter | undefined = chaptersData?.data?.chapters?.find(
        (c) => c.id === chapterId
    );

    // -- fetch pages
    const {
        data: pagesData,
        isLoading,
        error: fetchError,
    } = useQuery({
        queryKey: ["pages", chapterId],
        queryFn: () => fetchPages(chapterId),
        enabled: !isNaN(chapterId) && chapterId > 0,
    });

    const pages: Page[] = pagesData?.data?.pages ?? [];
    const bookTitle = book?.title ?? `Book #${bookId}`;
    const chapterNumber = chapter?.number ?? "?";

    // -- lightbox items for pages
    const pageLightboxItems: LightboxItem[] = useMemo(
        () =>
            pages.map((p) => ({
                url: p.url,
                label: `Page ${p.pageNumber}`,
            })),
        [pages]
    );

    // -- fetch all panels for all-panels lightbox (parallel queries)
    const allPanelQueries = useQueries({
        queries: pages.map((p) => ({
            queryKey: ["panels", p.id] as const,
            queryFn: () => fetchPanels(p.id),
            enabled: viewMode === "panels",
        })),
    });

    const allPanelDataKey = allPanelQueries.map((q) => q.dataUpdatedAt).join(",");
    const allPanelLightboxItems: LightboxItem[] = useMemo(
        () =>
            pages.flatMap((page, i) => {
                const panels: Panel[] = allPanelQueries[i]?.data?.data?.panels ?? [];
                return panels.map((p) => ({
                    url: p.url,
                    label: `P${p.panelNumber} — Page ${page.pageNumber}`,
                }));
            }),
        // eslint-disable-next-line react-hooks/exhaustive-deps
        [pages, allPanelDataKey]
    );

    return (
        <>
            {/* ── Page lightbox ── */}
            {pageLightboxIdx !== null ? (
                <Lightbox
                    items={pageLightboxItems}
                    currentIndex={pageLightboxIdx}
                    onIndexChange={setPageLightboxIdx}
                    onClose={closePageLightbox}
                />
            ) : null}

            {/* ── All-panels lightbox ── */}
            {panelLightboxIdx !== null ? (
                <Lightbox
                    items={allPanelLightboxItems}
                    currentIndex={panelLightboxIdx}
                    onIndexChange={setPanelLightboxIdx}
                    onClose={closePanelLightbox}
                />
            ) : null}

            {/* ── Panel viewer overlay (single page) ── */}
            {viewingPage !== null ? (
                <PanelViewer page={viewingPage} onClose={closePanelViewer} />
            ) : null}

            <div className="mx-auto max-w-5xl px-6 pb-20 max-sm:px-4">
                {/* ── Hero ──────────────────────────────────────────────────── */}
                <header className="relative pb-10 pt-20 max-sm:pb-7 max-sm:pt-12">
                    <Link
                        href={`/books/${bookId}`}
                        className="mb-6 inline-flex items-center gap-1.5 text-sm text-muted-foreground transition-colors hover:text-foreground"
                    >
                        {ArrowLeftIcon}
                        {bookTitle}
                    </Link>
                    <h1 className="text-[clamp(2.5rem,8vw,4.5rem)] font-black leading-[0.85] tracking-tighter max-sm:text-4xl">
                        Chapter {String(chapterNumber).padStart(2, "0")}
                    </h1>
                    {HeroUnderline}
                    <span className="mt-4 block text-[11px] font-medium uppercase tracking-[0.3em] text-muted-foreground">
                        {viewMode === "pages" ? "PAGE VIEWER" : "PANEL VIEWER"}
                    </span>
                </header>

                {/* ── Gallery section ───────────────────────────────────────── */}
                <section className="border-t border-border pt-8">
                    <div className="mb-6 flex flex-wrap items-center justify-between gap-3">
                        <h2 className="flex items-center gap-2.5 text-2xl font-semibold tracking-tight">
                            {viewMode === "pages" ? "Pages" : "Panels"}
                            {pages.length > 0 ? (
                                <span className="rounded-[4px_6px_5px_3px] bg-foreground px-2 py-0.5 text-xs font-medium text-background">
                                    {pages.length}
                                </span>
                            ) : null}
                        </h2>

                        {/* ── Segmented toggle ── */}
                        <div className="inline-flex items-center rounded-[5px_7px_6px_4px] border border-border bg-neutral-50 p-0.5">
                            <button
                                onClick={() => setViewMode("pages")}
                                className={`inline-flex cursor-pointer items-center gap-1.5 rounded-[4px_6px_5px_3px] px-3 py-1.5 text-xs font-medium transition-all ${
                                    viewMode === "pages"
                                        ? "bg-foreground text-background shadow-sm"
                                        : "text-muted-foreground hover:text-foreground"
                                }`}
                            >
                                {FileIcon}
                                Pages
                            </button>
                            <button
                                onClick={() => setViewMode("panels")}
                                className={`inline-flex cursor-pointer items-center gap-1.5 rounded-[4px_6px_5px_3px] px-3 py-1.5 text-xs font-medium transition-all ${
                                    viewMode === "panels"
                                        ? "bg-foreground text-background shadow-sm"
                                        : "text-muted-foreground hover:text-foreground"
                                }`}
                            >
                                {LayersIcon}
                                Panels
                            </button>
                        </div>
                    </div>

                    {viewMode === "pages" ? (
                        <PagesListContent
                            isLoading={isLoading}
                            fetchError={fetchError}
                            pages={pages}
                            onViewPanels={handleViewPanels}
                            onViewBig={setPageLightboxIdx}
                        />
                    ) : (
                        <AllPanelsListContent
                            pages={pages}
                            pagesLoading={isLoading}
                            onViewBig={setPanelLightboxIdx}
                        />
                    )}
                </section>
            </div>
        </>
    );
}
