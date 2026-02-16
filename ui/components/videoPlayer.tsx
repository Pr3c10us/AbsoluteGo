import { memo, useEffect, useRef, useState, useCallback } from "react";
import "plyr/dist/plyr.css";
import { useScrollLock } from "@/lib/use-scroll-lock";

const VideoPlayerOverlay = memo(function VideoPlayerOverlay({
    url,
    label,
    onClose,
}: {
    url: string;
    label: string;
    onClose: () => void;
}) {
    useScrollLock();
    // Use any to avoid TypeScript module headaches
    const playerRef = useRef<any>(null);
    const videoContainerRef = useRef<HTMLDivElement>(null);
    const [isVisible, setIsVisible] = useState(false);
    const [isReady, setIsReady] = useState(false);

    useEffect(() => {
        const timer = setTimeout(() => setIsVisible(true), 10);
        return () => clearTimeout(timer);
    }, []);

    useEffect(() => {
        const handleKeyDown = (e: KeyboardEvent) => {
            if (e.key === "Escape") onClose();
        };
        window.addEventListener("keydown", handleKeyDown);
        return () => window.removeEventListener("keydown", handleKeyDown);
    }, [onClose]);

    // Initialize Plyr with dynamic import only
    useEffect(() => {
        let isMounted = true;

        const initPlayer = async () => {
            if (!videoContainerRef.current || !url) return;

            const video = videoContainerRef.current.querySelector("video");
            if (!video) return;

            try {
                // Dynamic import handles both ES modules and CommonJS
                const PlyrModule = await import("plyr");
                const Plyr = PlyrModule.default || PlyrModule;

                if (!isMounted) return;

                // Cleanup previous instance
                if (playerRef.current) {
                    playerRef.current.destroy();
                }

                playerRef.current = new Plyr(video, {
                    controls: [
                        "play-large",
                        "play",
                        "progress",
                        "current-time",
                        "duration",
                        "mute",
                        "volume",
                        "captions",
                        "settings",
                        "pip",
                        "airplay",
                        "fullscreen",
                    ],
                    settings: ["captions", "quality", "speed"],
                    autoplay: true,
                    ratio: "16:9",
                    fullscreen: { enabled: true, iosNative: true },
                    keyboard: { focused: true, global: true },
                    tooltips: { controls: true, seek: true },
                    seekTime: 10,
                });

                playerRef.current.on("ready", () => {
                    if (isMounted) {
                        setIsReady(true);
                        playerRef.current?.play().catch(() => {});
                    }
                });

                playerRef.current.on("error", () => {
                    console.error("Video playback error");
                });
            } catch (err) {
                console.error("Failed to load Plyr:", err);
            }
        };

        initPlayer();

        return () => {
            isMounted = false;
            if (playerRef.current) {
                try {
                    playerRef.current.destroy();
                } catch (e) {
                    // Ignore cleanup errors
                }
                playerRef.current = null;
            }
        };
    }, [url]);

    const handleBackdropClick = useCallback(
        (e: React.MouseEvent<HTMLDivElement>) => {
            if (e.target === e.currentTarget) onClose();
        },
        [onClose],
    );

    if (!url) return null;

    return (
        <div
            className={`fixed inset-0 z-50 flex items-center justify-center bg-black transition-opacity duration-300 ease-out ${
                isVisible ? "opacity-100" : "opacity-0"
            }`}
            onClick={handleBackdropClick}
        >
            <div
                className={`relative w-full max-w-5xl mx-4 transition-all duration-300 ease-out ${
                    isVisible
                        ? "opacity-100 scale-100 translate-y-0"
                        : "opacity-0 scale-95 translate-y-4"
                }`}
            >
                {/* Monochrome Header */}
                <div className="flex items-center justify-between mb-4 px-2">
                    <div className="flex items-center gap-3 min-w-0">
                        <div className="w-1 h-6 bg-white/80 rounded-full shrink-0" />
                        <h2 className="text-lg font-medium text-white truncate pr-4 tracking-tight">
                            {label}
                        </h2>
                    </div>

                    <button
                        onClick={onClose}
                        className="group flex items-center justify-center w-10 h-10 rounded-full bg-white/10 hover:bg-white/20 text-white/70 hover:text-white transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-white/50 shrink-0"
                        aria-label="Close video player"
                    >
                        <svg
                            className="w-5 h-5 transition-transform duration-200 group-hover:rotate-90"
                            fill="none"
                            viewBox="0 0 24 24"
                            stroke="currentColor"
                            strokeWidth={2}
                        >
                            <path
                                strokeLinecap="round"
                                strokeLinejoin="round"
                                d="M6 18L18 6M6 6l12 12"
                            />
                        </svg>
                    </button>
                </div>

                {/* Video Container */}
                <div
                    ref={videoContainerRef}
                    className="relative aspect-video bg-black rounded-lg overflow-hidden shadow-2xl ring-1 ring-white/10"
                >
                    <video
                        className="plyr plyr--full-ui plyr--video w-full h-full"
                        playsInline
                        preload="auto"
                        crossOrigin="anonymous"
                    >
                        <source src={url} type="video/mp4" />
                        Your browser does not support the video element.
                    </video>

                    {/* Loading Spinner */}
                    {!isReady && (
                        <div className="absolute inset-0 flex items-center justify-center bg-black z-10">
                            <div className="w-8 h-8 border-2 border-white/20 border-t-white rounded-full animate-spin" />
                        </div>
                    )}
                </div>

                {/* Monochrome Footer */}
                <div className="mt-4 flex items-center justify-center gap-3 text-sm text-white/40">
                    <span className="hidden sm:inline">Press</span>
                    <kbd className="hidden sm:inline-flex px-2 py-1 rounded bg-white/10 font-mono text-xs border border-white/20 text-white/60">
                        ESC
                    </kbd>
                    <span className="hidden sm:inline">to close</span>
                    <span className="sm:hidden text-xs">
                        Tap outside to close
                    </span>
                </div>
            </div>

            {/* Monochrome Plyr Styles */}
            <style>{`
                .plyr {
                    --plyr-color-main: #ffffff;
                    --plyr-video-background: #000000;
                    --plyr-menu-background: rgba(0,0,0,0.95);
                    --plyr-menu-color: #ffffff;
                    --plyr-control-radius: 4px;
                    --plyr-range-thumb-background: #ffffff;
                    --plyr-range-fill-background: #ffffff;
                    --plyr-range-track-background: rgba(255,255,255,0.2);
                    --plyr-video-control-background-hover: rgba(255,255,255,0.1);
                    --plyr-video-control-color: rgba(255,255,255,0.8);
                    --plyr-video-control-color-hover: #ffffff;
                    --plyr-tooltip-background: rgba(0,0,0,0.9);
                    --plyr-tooltip-color: #ffffff;
                    --plyr-progress-buffered-background: rgba(255,255,255,0.15);
                    --plyr-menu-border-color: rgba(255,255,255,0.1);
                    --plyr-menu-arrow-color: rgba(0,0,0,0.95);
                }
                .plyr--video {
                    background: #000000;
                }
                .plyr__control--overlaid {
                    background: rgba(255,255,255,0.9) !important;
                    color: #000000 !important;
                    padding: 20px !important;
                }
                .plyr__control--overlaid:hover {
                    background: #ffffff !important;
                }
                .plyr__control--overlaid svg {
                    fill: #000000 !important;
                }
                .plyr--loading .plyr__progress__buffer {
                    color: rgba(255,255,255,0.3);
                }
            `}</style>
        </div>
    );
});

export default VideoPlayerOverlay;
