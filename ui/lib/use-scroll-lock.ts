import { useEffect } from "react";

/**
 * Locks body scroll while the calling component is mounted.
 * Restores the previous overflow value on unmount to avoid
 * clobbering nested overlays (e.g. lightbox inside a viewer).
 */
export function useScrollLock() {
    useEffect(() => {
        const prev = document.body.style.overflow;
        document.body.style.overflow = "hidden";
        return () => {
            document.body.style.overflow = prev;
        };
    }, []);
}
