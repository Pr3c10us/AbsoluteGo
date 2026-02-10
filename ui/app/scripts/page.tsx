import Link from "next/link";

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

const ScriptIcon = (
    <svg
        className="mx-auto mb-3 h-10 w-10 text-neutral-300"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        strokeWidth="1.5"
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

export default function ScriptsPage() {
    return (
        <div className="mx-auto min-h-screen max-w-5xl px-6 pb-20 max-sm:px-4">
            <header className="relative pb-10 pt-20 max-sm:pb-7 max-sm:pt-12">
                <Link
                    href="/"
                    className="mb-6 inline-flex items-center gap-1.5 text-sm text-muted-foreground transition-colors hover:text-foreground"
                >
                    {ArrowLeftIcon}
                    Home
                </Link>
                <h1 className="text-[clamp(3.5rem,10vw,6rem)] font-bold leading-[0.85] tracking-tighter max-sm:text-5xl">
                    SCRIPTS
                </h1>
                {HeroUnderline}
                <span className="mt-4 block text-[11px] font-medium uppercase tracking-[0.3em] text-muted-foreground">
                    WRITING WORKSPACE
                </span>
            </header>

            <section className="border-t border-border pt-12">
                <div className="py-16 text-center text-muted-foreground">
                    {ScriptIcon}
                    <p className="text-sm font-medium text-foreground">
                        Coming soon.
                    </p>
                    <span className="text-xs">
                        Script management is under development.
                    </span>
                </div>
            </section>
        </div>
    );
}
