"use client";

import "./dock.css";

import {
    motion,
    useMotionValue,
    useSpring,
    useTransform,
    AnimatePresence,
} from "motion/react";
import {
    Children,
    cloneElement,
    useEffect,
    useMemo,
    useRef,
    useState,
    type ReactElement,
} from "react";
import Link from "next/link";
import { usePathname } from "next/navigation";

// ── Static SVG icons (hoisted — rendering-hoist-jsx) ────────────────────────

const BookIcon = (
    <svg
        width="20"
        height="20"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        strokeWidth="1.8"
        strokeLinecap="round"
        strokeLinejoin="round"
    >
        <path d="M4 19.5v-15A2.5 2.5 0 0 1 6.5 2H20v20H6.5a2.5 2.5 0 0 1 0-5H20" />
    </svg>
);

const ScriptIcon = (
    <svg
        width="20"
        height="20"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        strokeWidth="1.8"
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

const VideoIcon = (
    <svg
        width="20"
        height="20"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        strokeWidth="1.8"
        strokeLinecap="round"
        strokeLinejoin="round"
    >
        <path d="m16 13 5.223 3.482a.5.5 0 0 0 .777-.416V7.934a.5.5 0 0 0-.777-.416L16 11" />
        <rect x="2" y="6" width="14" height="12" rx="2" />
    </svg>
);

// ── ReactBits Dock internals (adapted) ──────────────────────────────────────

interface DockItemProps {
    children: ReactElement[];
    className?: string;
    onClick?: () => void;
    mouseX: ReturnType<typeof useMotionValue<number>>;
    spring: { mass: number; stiffness: number; damping: number };
    distance: number;
    magnification: number;
    baseItemSize: number;
}

function DockItem({
    children,
    className = "",
    onClick,
    mouseX,
    spring,
    distance,
    magnification,
    baseItemSize,
}: DockItemProps) {
    const ref = useRef<HTMLDivElement>(null);
    const isHovered = useMotionValue(0);

    const mouseDistance = useTransform(mouseX, (val) => {
        const rect = ref.current?.getBoundingClientRect() ?? {
            x: 0,
            width: baseItemSize,
        };
        return val - rect.x - baseItemSize / 2;
    });

    const targetSize = useTransform(
        mouseDistance,
        [-distance, 0, distance],
        [baseItemSize, magnification, baseItemSize]
    );
    const size = useSpring(targetSize, spring);

    return (
        <motion.div
            ref={ref}
            style={{ width: size, height: size }}
            onHoverStart={() => isHovered.set(1)}
            onHoverEnd={() => isHovered.set(0)}
            onFocus={() => isHovered.set(1)}
            onBlur={() => isHovered.set(0)}
            onClick={onClick}
            className={`dock-item ${className}`}
            tabIndex={0}
            role="button"
            aria-haspopup="true"
        >
            {Children.map(children, (child) =>
                cloneElement(child, { isHovered } as Record<string, unknown>)
            )}
        </motion.div>
    );
}

function DockLabel({
    children,
    ...rest
}: {
    children: React.ReactNode;
    isHovered?: ReturnType<typeof useMotionValue<number>>;
}) {
    const { isHovered } = rest;
    const [isVisible, setIsVisible] = useState(false);

    useEffect(() => {
        if (!isHovered) return;
        const unsubscribe = isHovered.on("change", (latest) => {
            setIsVisible(latest === 1);
        });
        return () => unsubscribe();
    }, [isHovered]);

    return (
        <AnimatePresence>
            {isVisible ? (
                <motion.div
                    initial={{ opacity: 0, y: 0 }}
                    animate={{ opacity: 1, y: -10 }}
                    exit={{ opacity: 0, y: 0 }}
                    transition={{ duration: 0.2 }}
                    className="dock-label"
                    role="tooltip"
                    style={{ x: "-50%" }}
                >
                    {children}
                </motion.div>
            ) : null}
        </AnimatePresence>
    );
}

function DockIcon({ children }: { children: React.ReactNode }) {
    return <div className="dock-icon">{children}</div>;
}

// ── Helpers ─────────────────────────────────────────────────────────────────

/** Extract book ID from pathname like /books/42 or /books/42/scripts */
function extractBookId(pathname: string): string | null {
    const match = pathname.match(/^\/books\/(\d+)/);
    return match ? match[1] : null;
}

// ── Main Dock ───────────────────────────────────────────────────────────────

const SPRING = { mass: 0.1, stiffness: 150, damping: 12 };
const MAGNIFICATION = 60;
const DISTANCE = 200;
const PANEL_HEIGHT = 56;
const DOCK_HEIGHT = 200;
const BASE_ITEM_SIZE = 42;

export default function Dock() {
    const pathname = usePathname();
    const mouseX = useMotionValue(Infinity);
    const isHovered = useMotionValue(0);

    const maxHeight = Math.max(
        DOCK_HEIGHT,
        MAGNIFICATION + MAGNIFICATION / 2 + 4
    );
    const heightRow = useTransform(isHovered, [0, 1], [PANEL_HEIGHT, maxHeight]);
    const height = useSpring(heightRow, SPRING);

    // Extract book ID from the current path
    const bookId = extractBookId(pathname);

    // Build context-aware nav items based on the current book
    const navItems = useMemo(() => {
        if (!bookId) return [];
        const base = `/books/${bookId}`;
        return [
            { href: base, label: "Chapters", icon: BookIcon },
            { href: `${base}/scripts`, label: "Scripts", icon: ScriptIcon },
            { href: `${base}/videos`, label: "Videos", icon: VideoIcon },
        ];
    }, [bookId]);

    // Don't show dock on home/books listing page or if no book context
    if (!bookId) return null;

    return (
        <motion.div
            style={{ height, scrollbarWidth: "none" }}
            className="dock-outer"
        >
            <motion.div
                onMouseMove={({ pageX }) => {
                    isHovered.set(1);
                    mouseX.set(pageX);
                }}
                onMouseLeave={() => {
                    isHovered.set(0);
                    mouseX.set(Infinity);
                }}
                className="dock-panel"
                style={{ height: PANEL_HEIGHT }}
                role="toolbar"
                aria-label="Application dock"
            >
                {navItems.map(({ href, label, icon }) => {
                    // Exact match for the book root, startsWith for sub-routes
                    const isActive =
                        href === `/books/${bookId}`
                            ? pathname === href ||
                            pathname === `${href}/chapters` ||
                            pathname.startsWith(`${href}/chapters/`)
                            : pathname === href || pathname.startsWith(href + "/");

                    return (
                        <Link key={href} href={href} className="dock-link">
                            <DockItem
                                mouseX={mouseX}
                                spring={SPRING}
                                distance={DISTANCE}
                                magnification={MAGNIFICATION}
                                baseItemSize={BASE_ITEM_SIZE}
                                className={isActive ? "dock-item--active" : ""}
                            >
                                <DockIcon>{icon}</DockIcon>
                                <DockLabel>{label}</DockLabel>
                            </DockItem>
                        </Link>
                    );
                })}
            </motion.div>
        </motion.div>
    );
}
