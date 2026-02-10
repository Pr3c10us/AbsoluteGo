"use client";

import "./upload-tracker.css";

import {
    memo,
    useState,
    useRef,
    useCallback,
    useEffect,
    type PointerEvent as ReactPointerEvent,
} from "react";
import { useUpload } from "@/lib/upload-context";

// ── Static SVG icons (hoisted — rendering-hoist-jsx) ────────────────────────

const UploadingSpinner = (
    <svg className="ut-spinner" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5">
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

const UploadIcon = (
    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
        <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4" />
        <polyline points="17 8 12 3 7 8" />
        <line x1="12" y1="3" x2="12" y2="15" />
    </svg>
);

// ── Corner snapping ─────────────────────────────────────────────────────────

type Corner = "tl" | "tr" | "bl" | "br";

const MARGIN = 20;
const PILL_SIZE = 44;
const DOCK_CLEARANCE = 20;

function getPillPosition(corner: Corner) {
    const vw = typeof window !== "undefined" ? window.innerWidth : 1200;
    const vh = typeof window !== "undefined" ? window.innerHeight : 800;

    switch (corner) {
        case "tl": return { x: MARGIN, y: MARGIN };
        case "tr": return { x: vw - PILL_SIZE - MARGIN, y: MARGIN };
        case "bl": return { x: MARGIN, y: vh - PILL_SIZE - DOCK_CLEARANCE };
        case "br": return { x: vw - PILL_SIZE - MARGIN, y: vh - PILL_SIZE - DOCK_CLEARANCE };
    }
}

function findClosestCorner(cx: number, cy: number): Corner {
    const vw = typeof window !== "undefined" ? window.innerWidth : 1200;
    const vh = typeof window !== "undefined" ? window.innerHeight : 800;
    return cx < vw / 2
        ? cy < vh / 2 ? "tl" : "bl"
        : cy < vh / 2 ? "tr" : "br";
}

// ── Upload Tracker ──────────────────────────────────────────────────────────

const UploadTracker = memo(function UploadTracker() {
    const { uploads } = useUpload();
    const [corner, setCorner] = useState<Corner>("br");
    const [expanded, setExpanded] = useState(false);
    const [dragging, setDragging] = useState(false);
    const [dragPos, setDragPos] = useState<{ x: number; y: number } | null>(null);
    const dragOffsetRef = useRef({ x: 0, y: 0 });
    const didDragRef = useRef(false);
    const hideTimerRef = useRef<ReturnType<typeof setTimeout> | null>(null);

    const snappedPos = getPillPosition(corner);

    // ── Hover show/hide with delay ──
    const cancelHide = useCallback(() => {
        if (hideTimerRef.current) {
            clearTimeout(hideTimerRef.current);
            hideTimerRef.current = null;
        }
    }, []);

    const scheduleHide = useCallback(() => {
        cancelHide();
        hideTimerRef.current = setTimeout(() => setExpanded(false), 300);
    }, [cancelHide]);

    const showPanel = useCallback(() => {
        cancelHide();
        setExpanded(true);
    }, [cancelHide]);

    // ── Drag handlers (pill) ──
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

    // Re-snap on resize
    useEffect(() => {
        const h = () => setCorner((c) => c);
        window.addEventListener("resize", h);
        return () => window.removeEventListener("resize", h);
    }, []);

    // Cleanup timer
    useEffect(() => () => { if (hideTimerRef.current) clearTimeout(hideTimerRef.current); }, []);

    if (uploads.length === 0) return null;

    const activeCount = uploads.filter((u) => u.status === "uploading").length;
    const position = dragPos ?? snappedPos;
    const isTopCorner = corner === "tl" || corner === "tr";
    const isRightCorner = corner === "tr" || corner === "br";

    // Panel position: anchored to pill
    const panelStyle: React.CSSProperties = {
        left: isRightCorner ? position.x + PILL_SIZE - 260 : position.x,
        ...(isTopCorner
            ? { top: position.y + PILL_SIZE + 8 }
            : { top: position.y - 8, transform: "translateY(-100%)" }),
    };

    return (
        <>
            {/* ── Persistent pill ── */}
            <button
                className={`ut-pill ${dragging ? "ut-pill--dragging" : ""} ${activeCount > 0 ? "ut-pill--active" : ""}`}
                style={{
                    left: position.x,
                    top: position.y,
                    transition: dragging
                        ? "none"
                        : "left 0.35s cubic-bezier(0.4,0,0.2,1), top 0.35s cubic-bezier(0.4,0,0.2,1)",
                }}
                onClick={() => {
                    if (!didDragRef.current) setExpanded((v) => !v);
                }}
                onMouseEnter={showPanel}
                onMouseLeave={scheduleHide}
                onPointerDown={onPointerDown}
                onPointerMove={onPointerMove}
                onPointerUp={onPointerUp}
                aria-label="Upload tracker"
            >
                {activeCount > 0 ? UploadingSpinner : UploadIcon}
                {activeCount > 0 ? (
                    <span className="ut-pill-badge">{activeCount}</span>
                ) : null}
            </button>

            {/* ── Expanded panel ── */}
            {expanded ? (
                <div
                    className="ut-panel"
                    style={panelStyle}
                    onMouseEnter={showPanel}
                    onMouseLeave={scheduleHide}
                >
                    <div className="ut-panel-header">
                        <span className="ut-panel-title">
                            {activeCount > 0
                                ? `Processing ${activeCount} upload${activeCount > 1 ? "s" : ""}…`
                                : "Uploads"}
                        </span>
                    </div>
                    <ul className="ut-panel-list">
                        {uploads.map((upload) => (
                            <li key={upload.id} className="ut-panel-item">
                                <span className="ut-panel-item-icon">
                                    {upload.status === "uploading"
                                        ? UploadingSpinner
                                        : upload.status === "done"
                                            ? CheckIcon
                                            : ErrorIcon}
                                </span>
                                <div className="ut-panel-item-info">
                                    <span className="ut-panel-item-name">
                                        Ch.{upload.chapterNumber}
                                    </span>
                                    <span className="ut-panel-item-file">
                                        {upload.fileName}
                                    </span>
                                </div>
                                <span className={`ut-panel-badge ut-panel-badge--${upload.status}`}>
                                    {upload.status === "uploading"
                                        ? "Processing"
                                        : upload.status === "done"
                                            ? "Done"
                                            : "Failed"}
                                </span>
                            </li>
                        ))}
                    </ul>
                </div>
            ) : null}
        </>
    );
});

export default UploadTracker;
