"use client";

import {
    memo,
    useCallback,
    useEffect,
    useRef,
    useState,
    type PointerEvent as ReactPointerEvent,
} from "react";
import { X, ChevronLeft, ChevronRight, AlignLeft, AudioLines, Video, Play, Sparkles } from "lucide-react";
import { Button } from "@/components/ui/button";
import { useScrollLock } from "@/lib/use-scroll-lock";

// ── Types ───────────────────────────────────────────────────────────────────

export interface LightboxItem {
    url: string;
    label: string;
    /** Optional text displayed alongside the image (e.g. split content) */
    description?: string;
    /** Optional small tag shown above the description (e.g. effect name) */
    tag?: string;
    /** Optional audio URL for this item */
    audioURL?: string | null;
    /** Optional video URL for this item */
    videoURL?: string | null;
    /** Split ID for action callbacks */
    splitId?: number;
}

interface LightboxProps {
    items: LightboxItem[];
    currentIndex: number;
    onIndexChange: (index: number) => void;
    onClose: () => void;
    /** Callback to generate audio for a split (opens voice dialog) */
    onGenerateAudio?: (splitId: number) => void;
    /** Callback to generate video for a split */
    onGenerateVideo?: (splitId: number) => void;
}

// ── Static icons (hoisted — rendering-hoist-jsx) ────────────────────────────

const CloseIcon = <X className="h-5 w-5" />;

const ChevronLeftIcon = <ChevronLeft className="h-6 w-6" />;

const ChevronRightIcon = <ChevronRight className="h-6 w-6" />;

const TextIcon = <AlignLeft className="h-4 w-4" />;

const LbAudioIcon = <AudioLines className="h-3.5 w-3.5" />;

const LbVideoIcon = <Video className="h-3.5 w-3.5" />;

const LbPlayIcon = <Play className="h-3.5 w-3.5" />;

const LbSparklesIcon = <Sparkles className="h-3.5 w-3.5" />;

// ── Component ───────────────────────────────────────────────────────────────

const Lightbox = memo(function Lightbox({
    items,
    currentIndex,
    onIndexChange,
    onClose,
    onGenerateAudio,
    onGenerateVideo,
}: LightboxProps) {
    const [playingAudio, setPlayingAudio] = useState(false);
    const hasPrev = currentIndex > 0;
    const hasNext = currentIndex < items.length - 1;
    const hasAnyDescription = items.some((item) => item.description);

    const [showText, setShowText] = useState(hasAnyDescription);
    const toggleText = useCallback(() => setShowText((v) => !v), []);

    useScrollLock();

    // -- refs for stable keyboard handler (advanced-event-handler-refs,
    //    rerender-use-ref-transient-values) ──────────────────────────────
    const indexRef = useRef(currentIndex);
    const lengthRef = useRef(items.length);
    const onIndexChangeRef = useRef(onIndexChange);
    const onCloseRef = useRef(onClose);
    indexRef.current = currentIndex;
    lengthRef.current = items.length;
    onIndexChangeRef.current = onIndexChange;
    onCloseRef.current = onClose;

    // Stable callbacks that never change identity (rerender-functional-setstate)
    const goPrev = useCallback(() => {
        const idx = indexRef.current;
        if (idx > 0) onIndexChangeRef.current(idx - 1);
    }, []);

    const goNext = useCallback(() => {
        const idx = indexRef.current;
        if (idx < lengthRef.current - 1) onIndexChangeRef.current(idx + 1);
    }, []);

    // Keyboard navigation — effect registers once, never re-subscribes
    // (client-event-listeners)
    useEffect(() => {
        function handleKey(e: KeyboardEvent) {
            if (e.key === "ArrowLeft") {
                e.preventDefault();
                const idx = indexRef.current;
                if (idx > 0) onIndexChangeRef.current(idx - 1);
            } else if (e.key === "ArrowRight") {
                e.preventDefault();
                const idx = indexRef.current;
                if (idx < lengthRef.current - 1) onIndexChangeRef.current(idx + 1);
            } else if (e.key === "Escape") {
                e.preventDefault();
                onCloseRef.current();
            }
        }

        window.addEventListener("keydown", handleKey);
        return () => {
            window.removeEventListener("keydown", handleKey);
        };
    }, []);

    // -- Drag-to-scroll for thumbnail strip --
    const thumbRef = useRef<HTMLDivElement>(null);
    const [thumbDragging, setThumbDragging] = useState(false);
    const thumbStartRef = useRef({ x: 0, scrollLeft: 0 });

    const onThumbPointerDown = useCallback(
        (e: ReactPointerEvent<HTMLDivElement>) => {
            const el = thumbRef.current;
            if (!el) return;
            setThumbDragging(true);
            thumbStartRef.current = { x: e.clientX, scrollLeft: el.scrollLeft };
            el.setPointerCapture(e.pointerId);
        },
        []
    );

    const onThumbPointerMove = useCallback(
        (e: ReactPointerEvent<HTMLDivElement>) => {
            if (!thumbDragging || !thumbRef.current) return;
            const dx = e.clientX - thumbStartRef.current.x;
            thumbRef.current.scrollLeft = thumbStartRef.current.scrollLeft - dx;
        },
        [thumbDragging]
    );

    const onThumbPointerUp = useCallback(() => {
        setThumbDragging(false);
    }, []);

    // Scroll active thumbnail into view on index change + reset audio player
    useEffect(() => {
        setPlayingAudio(false);
        const el = thumbRef.current;
        if (!el) return;
        const active = el.children[currentIndex] as HTMLElement | undefined;
        if (active) {
            active.scrollIntoView({ behavior: "smooth", inline: "center", block: "nearest" });
        }
    }, [currentIndex]);

    // js-early-exit
    if (items.length === 0) return null;
    const current = items[currentIndex];

    return (
        <div className="fixed inset-0 z-50 flex flex-col bg-black/95 overflow-hidden">
            {/* ── Top bar ── */}
            <div className="flex items-center justify-between px-4 py-3 sm:px-6">
                <span className="text-sm font-medium text-white/80">
                    {current.label}
                </span>
                <div className="flex items-center gap-3">
                    <span className="font-mono text-xs text-white/50">
                        {currentIndex + 1} / {items.length}
                    </span>
                    {hasAnyDescription ? (
                        <Button
                            variant="ghost"
                            size="icon-sm"
                            onClick={toggleText}
                            aria-label={showText ? "Hide text" : "Show text"}
                            className={`cursor-pointer transition-colors hover:bg-white/10 hover:text-white ${showText ? "text-white bg-white/10" : "text-white/50"}`}
                        >
                            {TextIcon}
                        </Button>
                    ) : null}
                    <Button
                        variant="ghost"
                        size="icon-sm"
                        onClick={onClose}
                        aria-label="Close viewer"
                        className="cursor-pointer text-white/70 hover:bg-white/10 hover:text-white"
                    >
                        {CloseIcon}
                    </Button>
                </div>
            </div>

            {/* ── Main content area ── */}
            <div className="relative flex flex-1 overflow-hidden max-sm:flex-col">
                {/* ── Image area (click outside image → close) ── */}
                <div
                    className="relative flex flex-1 items-center justify-center overflow-hidden px-14 max-sm:px-2 cursor-pointer"
                    onClick={(e) => {
                        if (e.target === e.currentTarget) onClose();
                    }}
                >
                    {hasPrev ? (
                        <button
                            onClick={goPrev}
                            aria-label="Previous"
                            className="cursor-pointer absolute left-2 top-1/2 z-10 flex h-10 w-10 -translate-y-1/2 items-center justify-center rounded-full bg-white/10 text-white/70 backdrop-blur-sm transition-colors hover:bg-white/20 hover:text-white sm:left-4 sm:h-12 sm:w-12"
                        >
                            {ChevronLeftIcon}
                        </button>
                    ) : null}

                    {/* eslint-disable-next-line @next/next/no-img-element */}
                    <img
                        key={current.url}
                        src={current.url}
                        alt={current.label}
                        className="max-h-[calc(100vh-10rem)] max-w-full object-contain animate-in fade-in-0 zoom-in-95 duration-200 cursor-default"
                        onClick={(e) => e.stopPropagation()}
                    />

                    {hasNext ? (
                        <button
                            onClick={goNext}
                            aria-label="Next"
                            className="cursor-pointer absolute right-2 top-1/2 z-10 flex h-10 w-10 -translate-y-1/2 items-center justify-center rounded-full bg-white/10 text-white/70 backdrop-blur-sm transition-colors hover:bg-white/20 hover:text-white sm:right-4 sm:h-12 sm:w-12"
                        >
                            {ChevronRightIcon}
                        </button>
                    ) : null}
                </div>

                {/* ── Description panel ── */}
                {showText && current.description ? (
                    <div className="sm:w-80 sm:shrink-0 sm:border-l sm:border-white/10 max-sm:max-h-[35vh] max-sm:border-t max-sm:border-white/10 overflow-y-auto animate-in slide-in-from-right-4 sm:slide-in-from-right-4 max-sm:slide-in-from-bottom-4 duration-200">
                        <div className="p-4 sm:p-5">
                            {current.tag ? (
                                <span className="mb-3 inline-block rounded-[3px_5px_4px_3px] bg-white/15 px-2 py-0.5 text-[10px] font-semibold uppercase tracking-widest text-white/70">
                                    {current.tag}
                                </span>
                            ) : null}
                            <p className="whitespace-pre-wrap text-[13px] leading-relaxed text-white/85">
                                {current.description}
                            </p>

                            {/* ── Media section ── */}
                            {(current.audioURL || current.videoURL || current.splitId) ? (
                                <div className="mt-4 space-y-3 border-t border-white/10 pt-4">
                                    <span className="block text-[10px] font-semibold uppercase tracking-widest text-white/50">
                                        Media
                                    </span>

                                    {/* Audio player */}
                                    {current.audioURL ? (
                                        <div>
                                            <div className="mb-1.5 flex items-center gap-1.5 text-[11px] font-medium text-white/70">
                                                {LbAudioIcon}
                                                Audio
                                            </div>
                                            {playingAudio ? (
                                                <div>
                                                    {/* eslint-disable-next-line jsx-a11y/media-has-caption */}
                                                    <audio
                                                        controls
                                                        autoPlay
                                                        className="w-full h-8"
                                                        src={current.audioURL}
                                                        onEnded={() => setPlayingAudio(false)}
                                                    />
                                                </div>
                                            ) : (
                                                <button
                                                    onClick={() => setPlayingAudio(true)}
                                                    className="inline-flex cursor-pointer items-center gap-1.5 rounded-[3px_5px_4px_3px] bg-white/15 px-2.5 py-1 text-[11px] font-medium text-white/80 transition-colors hover:bg-white/25"
                                                >
                                                    {LbPlayIcon}
                                                    Play Audio
                                                </button>
                                            )}
                                        </div>
                                    ) : current.splitId && onGenerateAudio ? (
                                        <button
                                            onClick={() => onGenerateAudio(current.splitId!)}
                                            className="inline-flex cursor-pointer items-center gap-1.5 rounded-[3px_5px_4px_3px] bg-white/10 px-2.5 py-1 text-[11px] font-medium text-white/60 transition-colors hover:bg-white/20 hover:text-white/80"
                                        >
                                            {LbSparklesIcon}
                                            Generate Audio
                                        </button>
                                    ) : null}

                                    {/* Video link / generate */}
                                    {current.videoURL ? (
                                        <div>
                                            <div className="mb-1.5 flex items-center gap-1.5 text-[11px] font-medium text-white/70">
                                                {LbVideoIcon}
                                                Video
                                            </div>
                                            <a
                                                href={current.videoURL}
                                                target="_blank"
                                                rel="noopener noreferrer"
                                                className="inline-flex cursor-pointer items-center gap-1.5 rounded-[3px_5px_4px_3px] bg-white/15 px-2.5 py-1 text-[11px] font-medium text-white/80 transition-colors hover:bg-white/25"
                                            >
                                                {LbPlayIcon}
                                                Play Video
                                            </a>
                                        </div>
                                    ) : current.splitId && current.audioURL && onGenerateVideo ? (
                                        <button
                                            onClick={() => onGenerateVideo(current.splitId!)}
                                            className="inline-flex cursor-pointer items-center gap-1.5 rounded-[3px_5px_4px_3px] bg-white/10 px-2.5 py-1 text-[11px] font-medium text-white/60 transition-colors hover:bg-white/20 hover:text-white/80"
                                        >
                                            {LbSparklesIcon}
                                            Generate Video
                                        </button>
                                    ) : current.splitId && !current.audioURL ? (
                                        <span className="block text-[11px] text-white/30">
                                            Video needs audio first
                                        </span>
                                    ) : null}
                                </div>
                            ) : null}
                        </div>
                    </div>
                ) : null}
            </div>

            {/* ── Bottom thumbnail strip (drag-to-scroll, no overflow scrollbar) ── */}
            {items.length > 1 ? (
                <div
                    ref={thumbRef}
                    className="flex items-center gap-1.5 overflow-x-hidden px-4 py-3 sm:gap-2 sm:px-6"
                    style={{
                        cursor: thumbDragging ? "grabbing" : "grab",
                        scrollBehavior: thumbDragging ? "auto" : "smooth",
                    }}
                    onPointerDown={onThumbPointerDown}
                    onPointerMove={onThumbPointerMove}
                    onPointerUp={onThumbPointerUp}
                    onPointerCancel={onThumbPointerUp}
                >
                    {items.map((item, i) => (
                        <button
                            key={item.url}
                            onClick={() => onIndexChange(i)}
                            aria-label={item.label}
                            className={`cursor-pointer h-12 w-9 shrink-0 overflow-hidden rounded-[3px_4px_3px_4px] border-2 transition-all sm:h-14 sm:w-10 ${i === currentIndex
                                ? "border-white shadow-[0_0_8px_rgba(255,255,255,0.3)]"
                                : "border-transparent opacity-50 hover:opacity-80"
                                }`}
                        >
                            {/* eslint-disable-next-line @next/next/no-img-element */}
                            <img
                                src={item.url}
                                alt={item.label}
                                className="h-full w-full object-cover pointer-events-none"
                            />
                        </button>
                    ))}
                </div>
            ) : null}
        </div>
    );
});

export default Lightbox;
