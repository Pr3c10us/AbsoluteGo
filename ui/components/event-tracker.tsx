"use client";

import {
    memo,
    useState,
    useRef,
    useCallback,
    useEffect,
    type PointerEvent as ReactPointerEvent,
} from "react";
import { useRouter } from "next/navigation";
import { useQuery } from "@tanstack/react-query";
import {
    fetchEvents,
    type EventItem,
    type EventStatus,
    type EventOperation,
} from "@/lib/api";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { ScrollArea } from "@/components/ui/scroll-area";
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from "@/components/ui/select";

// ── Static SVG icons (hoisted) ──────────────────────────────────────────────

const SpinnerIcon = (
    <svg className="animate-spin" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5">
        <path d="M21 12a9 9 0 1 1-6.219-8.56" />
    </svg>
);

const CheckIcon = (
    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round">
        <path d="M20 6 9 17l-5-5" />
    </svg>
);

const ErrorIcon = (
    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round">
        <circle cx="12" cy="12" r="10" />
        <path d="m15 9-6 6" />
        <path d="m9 9 6 6" />
    </svg>
);

const QueueIcon = (
    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
        <circle cx="12" cy="12" r="10" />
        <polyline points="12 6 12 12 16 14" />
    </svg>
);

const RetryIcon = (
    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
        <path d="M21 12a9 9 0 0 0-9-9 9.75 9.75 0 0 0-6.74 2.74L3 8" />
        <path d="M3 3v5h5" />
        <path d="M3 12a9 9 0 0 0 9 9 9.75 9.75 0 0 0 6.74-2.74L21 16" />
        <path d="M16 16h5v5" />
    </svg>
);

const EventsIcon = (
    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
        <path d="M16 3h5v5" />
        <path d="M8 3H3v5" />
        <path d="M12 22v-8.3a4 4 0 0 0-1.172-2.872L3 3" />
        <path d="m15 9 6-6" />
    </svg>
);

// ── Status helpers ──────────────────────────────────────────────────────────

function getStatusIcon(status: EventStatus) {
    switch (status) {
        case "enqueue": return QueueIcon;
        case "processing": return SpinnerIcon;
        case "failed": return ErrorIcon;
        case "successful": return CheckIcon;
        case "retry": return RetryIcon;
    }
}

function getStatusLabel(status: EventStatus): string {
    switch (status) {
        case "enqueue": return "Queued";
        case "processing": return "Processing";
        case "failed": return "Failed";
        case "successful": return "Done";
        case "retry": return "Retrying";
    }
}

function getStatusBadgeVariant(status: EventStatus): "default" | "secondary" | "destructive" | "outline" {
    switch (status) {
        case "successful": return "default";
        case "failed": return "destructive";
        case "processing":
        case "retry": return "secondary";
        case "enqueue": return "outline";
    }
}

function getOperationLabel(op: EventOperation): string {
    switch (op) {
        case "add_chapter": return "Add Chapter";
        case "gen_script": return "Generate Script";
        case "gen_script_split": return "Generate Splits";
        case "gen_audio": return "Generate Audio";
        case "gen_video": return "Generate Video";
        case "merge_video": return "Merge Video";
    }
}

// ── Relative time ───────────────────────────────────────────────────────────

function timeAgo(dateStr: string): string {
    const diff = Date.now() - new Date(dateStr).getTime();
    const secs = Math.floor(diff / 1000);
    if (secs < 60) return "just now";
    const mins = Math.floor(secs / 60);
    if (mins < 60) return `${mins}m ago`;
    const hrs = Math.floor(mins / 60);
    if (hrs < 24) return `${hrs}h ago`;
    const days = Math.floor(hrs / 24);
    return `${days}d ago`;
}

// ── Corner snapping ─────────────────────────────────────────────────────────

type Corner = "tl" | "tr" | "bl" | "br";

const MARGIN = 20;
const PILL_SIZE = 44;
const DOCK_CLEARANCE = 20;
const PAGE_SIZE = 10;

function getPillPosition(corner: Corner) {
    const vw = window.innerWidth;
    const vh = window.innerHeight;
    switch (corner) {
        case "tl": return { x: MARGIN, y: MARGIN };
        case "tr": return { x: vw - PILL_SIZE - MARGIN, y: MARGIN };
        case "bl": return { x: MARGIN, y: vh - PILL_SIZE - DOCK_CLEARANCE };
        case "br": return { x: vw - PILL_SIZE - MARGIN, y: vh - PILL_SIZE - DOCK_CLEARANCE };
    }
}

function findClosestCorner(cx: number, cy: number): Corner {
    const vw = window.innerWidth;
    const vh = window.innerHeight;
    return cx < vw / 2
        ? cy < vh / 2 ? "tl" : "bl"
        : cy < vh / 2 ? "tr" : "br";
}

// ── Event click → navigation ────────────────────────────────────────────────

function getEventHref(event: EventItem): string | null {
    const bookId = event.BookId;
    if (!bookId) return null;

    switch (event.Operation) {
        case "add_chapter":
            return event.ChapterId
                ? `/books/${bookId}/chapters/${event.ChapterId}`
                : `/books/${bookId}`;
        case "gen_script":
            return event.ScriptId
                ? `/books/${bookId}/scripts?scriptId=${event.ScriptId}`
                : `/books/${bookId}/scripts`;
        case "gen_script_split":
        case "gen_audio":
        case "gen_video":
            return event.ScriptId
                ? `/books/${bookId}/scripts/${event.ScriptId}/splits`
                : `/books/${bookId}/scripts`;
        case "merge_video":
            return `/books/${bookId}/videos`;
        default:
            return `/books/${bookId}/scripts`;
    }
}

// ── Filter options ──────────────────────────────────────────────────────────

const STATUS_OPTIONS: { value: string; label: string }[] = [
    { value: "all", label: "All Status" },
    { value: "enqueue", label: "Queued" },
    { value: "processing", label: "Processing" },
    { value: "failed", label: "Failed" },
    { value: "successful", label: "Done" },
    { value: "retry", label: "Retrying" },
];

const OPERATION_OPTIONS: { value: string; label: string }[] = [
    { value: "all", label: "All Operations" },
    { value: "add_chapter", label: "Add Chapter" },
    { value: "gen_script", label: "Gen Script" },
    { value: "gen_script_split", label: "Gen Splits" },
    { value: "gen_audio", label: "Gen Audio" },
    { value: "gen_video", label: "Gen Video" },
    { value: "merge_video", label: "Merge Video" },
];

// ── Event Tracker ───────────────────────────────────────────────────────────

const EventTracker = memo(function EventTracker() {
    const router = useRouter();
    const [mounted, setMounted] = useState(false);
    const [corner, setCorner] = useState<Corner>("br");
    const [expanded, setExpanded] = useState(false);
    const [pinned, setPinned] = useState(false);
    const [dragging, setDragging] = useState(false);
    const [dragPos, setDragPos] = useState<{ x: number; y: number } | null>(null);
    const dragOffsetRef = useRef({ x: 0, y: 0 });
    const didDragRef = useRef(false);
    const hideTimerRef = useRef<ReturnType<typeof setTimeout> | null>(null);
    const panelRef = useRef<HTMLDivElement>(null);

    // ── Filters & pagination state ──
    const [filterStatus, setFilterStatus] = useState("all");
    const [filterOperation, setFilterOperation] = useState("all");
    const [page, setPage] = useState(1);

    useEffect(() => setMounted(true), []);
    useEffect(() => setPage(1), [filterStatus, filterOperation]);

    // Poll events every 5 seconds
    const { data } = useQuery({
        queryKey: ["events", filterStatus, filterOperation, page],
        queryFn: () =>
            fetchEvents({
                limit: PAGE_SIZE,
                page,
                ...(filterStatus !== "all" ? { status: filterStatus as EventStatus } : {}),
                ...(filterOperation !== "all" ? { operation: filterOperation as EventOperation } : {}),
            }),
        refetchInterval: 5000,
    });

    const events: EventItem[] = data?.data?.events ?? [];
    const hasMore = events.length === PAGE_SIZE;

    // Active count (unfiltered) for pill badge
    const { data: allData } = useQuery({
        queryKey: ["events", "active-count"],
        queryFn: () => fetchEvents({ limit: 100 }),
        refetchInterval: 5000,
    });

    const allEvents: EventItem[] = allData?.data?.events ?? [];
    const activeCount = allEvents.filter(
        (e) => e.Status === "processing" || e.Status === "enqueue" || e.Status === "retry"
    ).length;

    // ── Hover show/hide with delay ──
    const cancelHide = useCallback(() => {
        if (hideTimerRef.current) {
            clearTimeout(hideTimerRef.current);
            hideTimerRef.current = null;
        }
    }, []);

    const scheduleHide = useCallback(() => {
        if (pinned) return;
        // Don't hide while a Radix Select dropdown is open
        if (document.querySelector("[data-radix-select-viewport]")) return;
        cancelHide();
        hideTimerRef.current = setTimeout(() => setExpanded(false), 300);
    }, [cancelHide, pinned]);

    const showPanel = useCallback(() => {
        cancelHide();
        setExpanded(true);
    }, [cancelHide]);

    // ── Close on click outside when pinned ──
    useEffect(() => {
        if (!pinned) return;
        const handleClick = (e: MouseEvent) => {
            const target = e.target as HTMLElement;
            // Ignore clicks inside any Radix portal (Select, Popover, etc.)
            if (
                target.closest?.("[data-radix-select-viewport]") ||
                target.closest?.("[data-radix-select-content]") ||
                target.closest?.("[data-radix-popper-content-wrapper]") ||
                target.closest?.("[role='listbox']")
            ) return;
            if (panelRef.current && !panelRef.current.contains(target)) {
                setPinned(false);
                setExpanded(false);
            }
        };
        const id = setTimeout(() => document.addEventListener("mousedown", handleClick), 0);
        return () => {
            clearTimeout(id);
            document.removeEventListener("mousedown", handleClick);
        };
    }, [pinned]);

    // ── Drag handlers ──
    const onPointerDown = useCallback(
        (e: ReactPointerEvent<HTMLButtonElement>) => {
            didDragRef.current = false;
            const rect = e.currentTarget.getBoundingClientRect();
            dragOffsetRef.current = { x: e.clientX - rect.left, y: e.clientY - rect.top };
            setDragging(true);
            setDragPos({ x: rect.left, y: rect.top });
            e.currentTarget.setPointerCapture(e.pointerId);
        },
        []
    );

    const onPointerMove = useCallback(
        (e: ReactPointerEvent<HTMLButtonElement>) => {
            if (!dragging) return;
            didDragRef.current = true;
            setDragPos({
                x: e.clientX - dragOffsetRef.current.x,
                y: e.clientY - dragOffsetRef.current.y,
            });
        },
        [dragging]
    );

    const onPointerUp = useCallback(
        (e: ReactPointerEvent<HTMLButtonElement>) => {
            if (!dragging) return;
            setDragging(false);
            const cx = e.clientX - dragOffsetRef.current.x + PILL_SIZE / 2;
            const cy = e.clientY - dragOffsetRef.current.y + PILL_SIZE / 2;
            setCorner(findClosestCorner(cx, cy));
            setDragPos(null);
        },
        [dragging]
    );

    useEffect(() => {
        const h = () => setCorner((c) => c);
        window.addEventListener("resize", h);
        return () => window.removeEventListener("resize", h);
    }, []);

    useEffect(() => () => { if (hideTimerRef.current) clearTimeout(hideTimerRef.current); }, []);

    const handlePillClick = useCallback(() => {
        if (didDragRef.current) return;
        if (pinned) { setPinned(false); setExpanded(false); }
        else { setPinned(true); setExpanded(true); }
    }, [pinned]);

    const handleEventClick = useCallback(
        (event: EventItem) => {
            const href = getEventHref(event);
            if (href) {
                router.push(href);
                setPinned(false);
                setExpanded(false);
            }
        },
        [router]
    );

    if (!mounted) return null;

    const snappedPos = getPillPosition(corner);
    const position = dragPos ?? snappedPos;
    const isTopCorner = corner === "tl" || corner === "tr";
    const isRightCorner = corner === "tr" || corner === "br";

    const panelStyle: React.CSSProperties = {
        left: isRightCorner ? position.x + PILL_SIZE - 340 : position.x,
        ...(isTopCorner
            ? { top: position.y + PILL_SIZE + 8 }
            : { top: position.y - 8, transform: "translateY(-100%)" }),
    };

    return (
        <>
            {/* ── Persistent pill ── */}
            <button
                className={`fixed z-40 pointer-events-auto w-11 h-11 rounded-full border border-white/[0.12] bg-black text-white flex items-center justify-center select-none touch-none p-0 overflow-visible ${dragging ? "cursor-grabbing shadow-[0_8px_32px_rgba(0,0,0,0.5)]" : "cursor-grab shadow-[0_4px_20px_rgba(0,0,0,0.3)] hover:shadow-[0_4px_24px_rgba(0,0,0,0.45)]"} ${activeCount > 0 ? "shadow-[0_0_16px_rgba(0,0,0,0.3),0_0_0_2px_rgba(255,255,255,0.08)]" : ""}`}
                style={{
                    left: position.x,
                    top: position.y,
                    transition: dragging
                        ? "none"
                        : "left 0.35s cubic-bezier(0.4,0,0.2,1), top 0.35s cubic-bezier(0.4,0,0.2,1)",
                }}
                onClick={handlePillClick}
                onMouseEnter={() => { if (!pinned) showPanel(); }}
                onMouseLeave={() => { if (!pinned) scheduleHide(); }}
                onPointerDown={onPointerDown}
                onPointerMove={onPointerMove}
                onPointerUp={onPointerUp}
                aria-label="Events"
            >
                {activeCount > 0 ? SpinnerIcon : EventsIcon}
                {activeCount > 0 ? (
                    <span className="absolute -top-1 -right-1 min-w-[18px] h-[18px] rounded-full bg-white text-black text-[0.6rem] font-extrabold flex items-center justify-center leading-none pointer-events-none">{activeCount}</span>
                ) : null}
            </button>

            {/* ── Expanded panel ── */}
            {expanded ? (
                <div
                    ref={panelRef}
                    className="fixed z-40 pointer-events-auto w-[340px] rounded-xl border border-black/[0.08] bg-white shadow-[0_8px_32px_rgba(0,0,0,0.12)] font-sans animate-in fade-in-0 zoom-in-[0.96] duration-150 overflow-hidden"
                    style={panelStyle}
                    onMouseEnter={() => { if (!pinned) showPanel(); }}
                    onMouseLeave={() => { if (!pinned) scheduleHide(); }}
                >
                    {/* Header */}
                    <div className="flex items-center px-3 py-2.5 border-b border-black/[0.06]">
                        <span className="text-[0.7rem] font-bold tracking-[0.02em] uppercase text-black">
                            {activeCount > 0
                                ? `Processing ${activeCount} event${activeCount > 1 ? "s" : ""}…`
                                : "Events"}
                        </span>
                    </div>

                    {/* Filters */}
                    <div className="flex gap-1.5 px-3 py-2 border-b border-black/[0.06]">
                        <Select value={filterStatus} onValueChange={setFilterStatus}>
                            <SelectTrigger size="sm" className="h-7 text-xs flex-1">
                                <SelectValue />
                            </SelectTrigger>
                            <SelectContent>
                                {STATUS_OPTIONS.map((opt) => (
                                    <SelectItem key={opt.value} value={opt.value}>
                                        {opt.label}
                                    </SelectItem>
                                ))}
                            </SelectContent>
                        </Select>

                        <Select value={filterOperation} onValueChange={setFilterOperation}>
                            <SelectTrigger size="sm" className="h-7 text-xs flex-1">
                                <SelectValue />
                            </SelectTrigger>
                            <SelectContent>
                                {OPERATION_OPTIONS.map((opt) => (
                                    <SelectItem key={opt.value} value={opt.value}>
                                        {opt.label}
                                    </SelectItem>
                                ))}
                            </SelectContent>
                        </Select>
                    </div>

                    {/* Event list */}
                    {events.length === 0 ? (
                        <div className="py-6 px-3 text-center text-[0.7rem] text-neutral-400">No events found</div>
                    ) : (
                        <ScrollArea className="max-h-[260px]">
                            <ul className="list-none m-0 py-1.5">
                                {events.map((event) => {
                                    const href = getEventHref(event);
                                    const clickable = href !== null;
                                    return (
                                        <li
                                            key={event.Id}
                                            className={`flex items-center gap-2 px-3 py-2 rounded transition-colors duration-150 ${clickable ? "cursor-pointer hover:bg-black/[0.03]" : "cursor-default opacity-60"}`}
                                            onClick={() => {
                                                if (clickable) handleEventClick(event);
                                            }}
                                        >
                                            <span className="flex shrink-0 text-black">
                                                {getStatusIcon(event.Status)}
                                            </span>
                                            <div className="flex flex-col flex-1 min-w-0">
                                                <span className="text-xs font-semibold text-black truncate">
                                                    {getOperationLabel(event.Operation)}
                                                </span>
                                                <span className="text-[0.6rem] text-neutral-500 truncate max-w-[200px]">
                                                    {event.Description}
                                                </span>
                                                <span className="text-[0.55rem] text-neutral-400 mt-px">
                                                    {timeAgo(event.UpdatedAt)}
                                                </span>
                                            </div>
                                            <Badge variant={getStatusBadgeVariant(event.Status)}>
                                                {getStatusLabel(event.Status)}
                                            </Badge>
                                        </li>
                                    );
                                })}
                            </ul>
                        </ScrollArea>
                    )}

                    {/* Pagination */}
                    {(page > 1 || hasMore) && (
                        <div className="flex items-center justify-between px-3 py-1.5 pb-2 border-t border-black/[0.06]">
                            <Button
                                variant="ghost"
                                size="sm"
                                className="h-7 text-xs px-2"
                                disabled={page <= 1}
                                onClick={() => setPage((p) => Math.max(1, p - 1))}
                            >
                                ← Prev
                            </Button>
                            <span className="text-[0.6rem] text-neutral-500 font-medium">Page {page}</span>
                            <Button
                                variant="ghost"
                                size="sm"
                                className="h-7 text-xs px-2"
                                disabled={!hasMore}
                                onClick={() => setPage((p) => p + 1)}
                            >
                                Next →
                            </Button>
                        </div>
                    )}
                </div>
            ) : null}
        </>
    );
});

export default EventTracker;
