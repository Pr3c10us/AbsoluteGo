"use client";

import dynamic from "next/dynamic";
import { memo } from "react";

const VideoPlayerOverlay = dynamic(() => import("./videoPlayer"), {
    ssr: false,
    loading: () => (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/95 backdrop-blur-sm">
            <div className="flex flex-col items-center gap-3">
                <div className="h-12 w-12 animate-spin rounded-full border-4 border-white/20 border-t-white"></div>
                <p className="text-sm text-white/60">Loading player...</p>
            </div>
        </div>
    ),
});

export default memo(VideoPlayerOverlay);
