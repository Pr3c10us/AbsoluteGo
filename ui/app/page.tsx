import Link from "next/link";

// ── Static SVG icons (hoisted — rendering-hoist-jsx) ────────────────────────

const BookIcon = (
  <svg
    width="32"
    height="32"
    viewBox="0 0 24 24"
    fill="none"
    stroke="currentColor"
    strokeWidth="1.5"
    strokeLinecap="round"
    strokeLinejoin="round"
  >
    <path d="M4 19.5v-15A2.5 2.5 0 0 1 6.5 2H20v20H6.5a2.5 2.5 0 0 1 0-5H20" />
  </svg>
);

const ScriptIcon = (
  <svg
    width="32"
    height="32"
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

const VideoIcon = (
  <svg
    width="32"
    height="32"
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

const HeroUnderline = (
  <svg className="mt-3 h-2 w-32 text-foreground" viewBox="0 0 120 8" fill="none">
    <path
      d="M2 5C25 2 50 7 75 4C100 1 115 6 118 3"
      stroke="currentColor"
      strokeWidth="2"
      strokeLinecap="round"
    />
  </svg>
);

const ArrowRightIcon = (
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
    <path d="M5 12h14" />
    <path d="m12 5 7 7-7 7" />
  </svg>
);

// ── Nav cards data ──────────────────────────────────────────────────────────

const NAV_CARDS = [
  {
    href: "/books",
    label: "BOOK",
    description: "Manage your comic book library — add chapters, view pages & panels",
    icon: BookIcon,
  },
  {
    href: "/scripts",
    label: "SCRIPT",
    description: "Write and organize scripts for your comic series",
    icon: ScriptIcon,
  },
  {
    href: "/videos",
    label: "VIDEO",
    description: "Create and manage video content from your comics",
    icon: VideoIcon,
  },
] as const;

// ── Home page ───────────────────────────────────────────────────────────────

export default function HomePage() {
  return (
    <div className="mx-auto flex min-h-screen max-w-5xl flex-col px-6 max-sm:px-4">
      {/* ── Hero ── */}
      <header className="flex flex-col items-center pt-32 pb-16 text-center max-sm:pt-20 max-sm:pb-10">
        <span className="mb-6 inline-block rounded-[4px_6px_5px_3px] bg-foreground px-3 py-1 text-[10px] font-bold uppercase tracking-[0.35em] text-background">
          MANAGEMENT SUITE
        </span>
        <h1 className="text-[clamp(4rem,14vw,8rem)] font-black leading-[0.85] tracking-tighter">
          ABSOLUTE
        </h1>
        {HeroUnderline}
        <p className="mt-6 max-w-md text-sm leading-relaxed text-muted-foreground">
          Your end-to-end workspace for comic books — from raw pages to
          AI-segmented panels, scripts, and video content.
        </p>
      </header>

      {/* ── Navigation cards ── */}
      <section className="grid gap-4 pb-20 sm:grid-cols-3">
        {NAV_CARDS.map(({ href, label, description, icon }) => (
          <Link
            key={href}
            href={href}
            className="group relative flex flex-col justify-between overflow-hidden rounded-[8px_10px_9px_7px] border border-border p-6 shadow-[3px_5px_14px_rgba(0,0,0,0.06)] transition-all duration-300 hover:shadow-[4px_7px_24px_rgba(0,0,0,0.15)] hover:-translate-y-1"
          >
            <div>
              <span className="mb-4 inline-flex items-center justify-center rounded-[5px_7px_6px_4px] border border-border p-2.5 text-foreground transition-transform duration-200 group-hover:scale-110">
                {icon}
              </span>
              <h2 className="mt-3 text-2xl font-black tracking-tighter">
                {label}
              </h2>
              <p className="mt-2 text-xs leading-relaxed text-muted-foreground">
                {description}
              </p>
            </div>
            <div className="mt-6 flex items-center gap-1.5 text-xs font-semibold text-foreground opacity-0 transition-opacity duration-200 group-hover:opacity-100">
              Enter
              {ArrowRightIcon}
            </div>
          </Link>
        ))}
      </section>
    </div>
  );
}
