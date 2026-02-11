"use client";

import Link from "next/link";
import { useParams } from "next/navigation";
import { useQuery } from "@tanstack/react-query";
import { fetchBooks, type Book } from "@/lib/api";

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

const HeroUnderline = (
    <svg className="mt-2 h-2 w-24 text-foreground" viewBox="0 0 120 8" fill="none">
        <path
            d="M2 5C25 2 50 7 75 4C100 1 115 6 118 3"
            stroke="currentColor"
            strokeWidth="2"
            strokeLinecap="round"
        />
    </svg>
);

const VideoEmptyIcon = (
    <svg
        className="mx-auto mb-3 h-10 w-10 text-neutral-300"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        strokeWidth="1.5"
        strokeLinecap="round"
        strokeLinejoin="round"
    >
        <path d="m16 13 5.223 3.482a.5.5 0 0 0 .777-.416V7.934a.5.5 0 0 0-.777-.416L16 11" />
        <rect x="2" y="6" width="14" height="12" rx="2" />
    </svg>
);

export default function BookVideosPage() {
    const params = useParams();
    const bookId = Number(params.id);

    const { data: booksData } = useQuery({
        queryKey: ["books"],
        queryFn: () => fetchBooks(),
    });

    const book: Book | undefined = booksData?.data?.books?.find(
        (b) => b.id === bookId
    );

    const bookTitle = book?.title ?? `Book #${bookId}`;

    return (
        <div className="mx-auto max-w-5xl px-6 pb-20 max-sm:px-4">
            <header className="relative pb-10 pt-20 max-sm:pb-7 max-sm:pt-12">
                <Link
                    href={`/books/${bookId}`}
                    className="mb-6 inline-flex items-center gap-1.5 text-sm text-muted-foreground transition-colors hover:text-foreground"
                >
                    {ArrowLeftIcon}
                    {bookTitle}
                </Link>
                <h1 className="text-[clamp(3.5rem,10vw,6rem)] font-bold leading-[0.85] tracking-tighter max-sm:text-5xl">
                    VIDEOS
                </h1>
                {HeroUnderline}
                <span className="mt-4 block text-[11px] font-medium uppercase tracking-[0.3em] text-muted-foreground">
                    VIDEO PRODUCTION
                </span>
            </header>

            <section className="border-t border-border pt-12">
                <div className="py-16 text-center text-muted-foreground">
                    {VideoEmptyIcon}
                    <p className="text-sm font-medium text-foreground">
                        Coming soon.
                    </p>
                    <span className="text-xs">
                        Video management for &ldquo;{bookTitle}&rdquo; is under development.
                    </span>
                </div>
            </section>
        </div>
    );
}
