"use client";

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
import { BookOpen, FileText, Video } from "lucide-react";

// ── Static icons (hoisted — rendering-hoist-jsx) ────────────────────────────

const BookIcon = <BookOpen className="h-5 w-5" strokeWidth={1.8} />;

const ScriptIcon = <FileText className="h-5 w-5" strokeWidth={1.8} />;

const VideoIcon = <Video className="h-5 w-5" strokeWidth={1.8} />;

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
            className={`relative inline-flex items-center justify-center rounded-[10px_12px_11px_9px] bg-black border border-white/[0.08] cursor-pointer outline-none text-white/60 transition-colors hover:text-white/95 ${className}`}
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
                    className="absolute -top-9 left-1/2 w-fit whitespace-pre rounded-[5px_7px_6px_4px] border border-white/10 bg-black px-2.5 py-2.5 text-[0.7rem] font-medium tracking-[0.04em] text-white"
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
    return <div className="flex items-center justify-center">{children}</div>;
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
            className="fixed bottom-0 left-1/2 -translate-x-1/2 z-40 flex max-w-full items-center mx-2"
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
                className="absolute bottom-3 left-1/2 -translate-x-1/2 flex items-end w-fit gap-3 rounded-[14px_16px_15px_13px] bg-black border border-white/10 px-2 pb-2 shadow-[0_8px_32px_rgba(0,0,0,0.35)]"
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
                        <Link key={href} href={href} className="no-underline leading-none">
                            <DockItem
                                mouseX={mouseX}
                                spring={SPRING}
                                distance={DISTANCE}
                                magnification={MAGNIFICATION}
                                baseItemSize={BASE_ITEM_SIZE}
                                className={isActive ? "!text-white bg-white/[0.12] shadow-[0_0_12px_rgba(255,255,255,0.1)]" : ""}
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
