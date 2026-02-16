"use client";

import { memo, useState, useCallback, useMemo } from "react";
import Link from "next/link";
import { useParams } from "next/navigation";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import {
    ArrowLeft,
    GitBranch,
    RefreshCw,
    Sparkles,
    Trash2,
    Maximize2,
    Eye,
    X,
    ChevronDown,
    AudioLines,
    Video,
    Play,
    MoreVertical,
    AlertTriangle,
} from "lucide-react";
import {
    fetchBooks,
    fetchScripts,
    fetchSplits,
    deleteSplits,
    generateSplits,
    generateAllAudios,
    generateAllVideos,
    generateSplitAudio,
    generateSplitVideo,
    ApiError,
    type Book,
    type Script,
    type Split,
} from "@/lib/api";
import { useScrollLock } from "@/lib/use-scroll-lock";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
    Select,
    SelectContent,
    SelectGroup,
    SelectItem,
    SelectLabel,
    SelectTrigger,
    SelectValue,
} from "@/components/ui/select";
import Lightbox, { type LightboxItem } from "@/components/lightbox";
import {
    AlertDialog,
    AlertDialogAction,
    AlertDialogCancel,
    AlertDialogContent,
    AlertDialogDescription,
    AlertDialogFooter,
    AlertDialogHeader,
    AlertDialogTitle,
} from "@/components/ui/alert-dialog";
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
} from "@/components/ui/dialog";
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuSeparator,
    DropdownMenuTrigger,
    DropdownMenuLabel,
} from "@/components/ui/dropdown-menu";
import VideoPlayerOverlayWrapper from "@/components/videoPlayerOverlay";

// ── Static icons (hoisted — rendering-hoist-jsx) ────────────────────────────

const ArrowLeftIcon = <ArrowLeft className="h-4 w-4" />;

const SplitEmptyIcon = (
    <GitBranch
        className="mx-auto mb-3 h-10 w-10 text-neutral-300"
        strokeWidth={1.5}
    />
);

const RefreshIcon = <RefreshCw className="h-4 w-4" />;

const SparklesIcon = <Sparkles className="h-4 w-4" />;

const TrashIcon = <Trash2 className="h-4 w-4" />;

const ExpandIcon = <Maximize2 className="h-3.5 w-3.5" />;

const EyeIcon = <Eye className="h-3.5 w-3.5" />;

const CloseIcon = <X className="h-4 w-4" />;

const HeroUnderline = (
    <svg
        className="mt-2 h-2 w-24 text-foreground"
        viewBox="0 0 120 8"
        fill="none"
    >
        <path
            d="M2 5C25 2 50 7 75 4C100 1 115 6 118 3"
            stroke="currentColor"
            strokeWidth="2"
            strokeLinecap="round"
        />
    </svg>
);

const ChevronDownIcon = <ChevronDown className="h-4 w-4" />;

const AudioIcon = <AudioLines className="h-4 w-4" />;

const VideoIcon = <Video className="h-4 w-4" />;

const PlayIcon = <Play className="h-3.5 w-3.5" />;

const MoreVerticalIcon = <MoreVertical className="h-4 w-4" />;

const AlertTriangleIcon = <AlertTriangle className="h-4 w-4" />;

// ── Effect label mapping ────────────────────────────────────────────────────

const EFFECT_LABELS: Record<string, string> = {
    panRight: "Pan Right",
    panLeft: "Pan Left",
    panUp: "Pan Up",
    panDown: "Pan Down",
    zoomIn: "Zoom In",
    zoomOut: "Zoom Out",
};

// ── Available voices ────────────────────────────────────────────────────────

interface Voice {
    name: string;
    gender: "Male" | "Female";
    desc: string;
}

const VOICES: Voice[] = [
    { name: "Achernar", gender: "Female", desc: "Soft" },
    { name: "Achird", gender: "Male", desc: "Friendly" },
    { name: "Algenib", gender: "Male", desc: "Gravelly" },
    { name: "Algieba", gender: "Male", desc: "Smooth" },
    { name: "Alnilam", gender: "Male", desc: "Firm" },
    { name: "Aoede", gender: "Female", desc: "Breezy" },
    { name: "Autonoe", gender: "Female", desc: "Bright" },
    { name: "Callirrhoe", gender: "Female", desc: "Easy-going" },
    { name: "Charon", gender: "Male", desc: "Informative" },
    { name: "Despina", gender: "Female", desc: "Smooth" },
    { name: "Enceladus", gender: "Male", desc: "Breathy" },
    { name: "Erinome", gender: "Female", desc: "Clear" },
    { name: "Fenrir", gender: "Male", desc: "Excitable" },
    { name: "Gacrux", gender: "Female", desc: "Mature" },
    { name: "Iapetus", gender: "Male", desc: "Clear" },
    { name: "Kore", gender: "Female", desc: "Firm" },
    { name: "Laomedeia", gender: "Female", desc: "Upbeat" },
    { name: "Leda", gender: "Female", desc: "Youthful" },
    { name: "Orus", gender: "Male", desc: "Firm" },
    { name: "Puck", gender: "Male", desc: "Upbeat" },
    { name: "Pulcherrima", gender: "Female", desc: "Forward" },
    { name: "Rasalgethi", gender: "Male", desc: "Informative" },
    { name: "Sadachbia", gender: "Male", desc: "Lively" },
    { name: "Sadaltager", gender: "Male", desc: "Knowledgeable" },
    { name: "Schedar", gender: "Male", desc: "Even" },
    { name: "Sulafat", gender: "Female", desc: "Warm" },
    { name: "Umbriel", gender: "Male", desc: "Easy-going" },
    { name: "Vindemiatrix", gender: "Female", desc: "Gentle" },
    { name: "Zephyr", gender: "Female", desc: "Bright" },
    { name: "Zubenelgenubi", gender: "Male", desc: "Casual" },
];

const MALE_VOICES = VOICES.filter((v) => v.gender === "Male");
const FEMALE_VOICES = VOICES.filter((v) => v.gender === "Female");

// ── Voice dialog for audio generation ───────────────────────────────────────

const VoiceDialog = memo(function VoiceDialog({
    open,
    onOpenChange,
    title,
    description,
    onSubmit,
}: {
    open: boolean;
    onOpenChange: (open: boolean) => void;
    title: string;
    description: string;
    onSubmit: (voice: string, voiceStyle: string) => void;
}) {
    const [voice, setVoice] = useState("Enceladus");
    const [voiceStyle, setVoiceStyle] = useState("");

    const selectedVoice = VOICES.find((v) => v.name === voice);

    const handleSubmit = useCallback(() => {
        if (!voice) {
            toast.error("Please select a voice");
            return;
        }
        onSubmit(voice, voiceStyle.trim());
        onOpenChange(false);
        setVoice("Enceladus");
        setVoiceStyle("");
    }, [voice, voiceStyle, onSubmit, onOpenChange]);

    return (
        <Dialog open={open} onOpenChange={onOpenChange}>
            <DialogContent>
                <DialogHeader>
                    <DialogTitle>{title}</DialogTitle>
                    <DialogDescription>{description}</DialogDescription>
                </DialogHeader>
                <div className="flex flex-col gap-3 py-2">
                    <div>
                        <label className="mb-1.5 block text-xs font-medium uppercase tracking-wider text-muted-foreground">
                            Voice
                        </label>
                        <Select value={voice} onValueChange={setVoice}>
                            <SelectTrigger className="w-full">
                                <SelectValue placeholder="Select a voice">
                                    {selectedVoice
                                        ? `${selectedVoice.name} — ${selectedVoice.desc}`
                                        : "Select a voice"}
                                </SelectValue>
                            </SelectTrigger>
                            <SelectContent
                                position="popper"
                                className="max-h-64"
                            >
                                <SelectGroup>
                                    <SelectLabel>Male</SelectLabel>
                                    {MALE_VOICES.map((v) => (
                                        <SelectItem key={v.name} value={v.name}>
                                            <span className="font-medium">
                                                {v.name}
                                            </span>
                                            <span className="ml-1 text-muted-foreground">
                                                — {v.desc}
                                            </span>
                                        </SelectItem>
                                    ))}
                                </SelectGroup>
                                <SelectGroup>
                                    <SelectLabel>Female</SelectLabel>
                                    {FEMALE_VOICES.map((v) => (
                                        <SelectItem key={v.name} value={v.name}>
                                            <span className="font-medium">
                                                {v.name}
                                            </span>
                                            <span className="ml-1 text-muted-foreground">
                                                — {v.desc}
                                            </span>
                                        </SelectItem>
                                    ))}
                                </SelectGroup>
                            </SelectContent>
                        </Select>
                    </div>
                    <div>
                        <label className="mb-1.5 block text-xs font-medium uppercase tracking-wider text-muted-foreground">
                            Voice Style (optional)
                        </label>
                        <Input
                            value={voiceStyle}
                            onChange={(e) => setVoiceStyle(e.target.value)}
                            placeholder="e.g. calm, narrative tone"
                        />
                    </div>
                </div>
                <DialogFooter>
                    <Button
                        variant="outline"
                        onClick={() => onOpenChange(false)}
                    >
                        Cancel
                    </Button>
                    <Button onClick={handleSubmit} className="gap-1.5">
                        {AudioIcon}
                        Generate
                    </Button>
                </DialogFooter>
            </DialogContent>
        </Dialog>
    );
});

// ── Script viewer overlay ───────────────────────────────────────────────────

const ScriptViewer = memo(function ScriptViewer({
    script,
    onClose,
}: {
    script: Script;
    onClose: () => void;
}) {
    useScrollLock();

    return (
        <div className="fixed inset-0 z-40 flex flex-col bg-white">
            <div className="flex items-center justify-between border-b border-border px-6 py-3">
                <div>
                    <h2 className="text-lg font-bold tracking-tight">
                        {script.name}
                    </h2>
                    <span className="text-xs text-muted-foreground">
                        Full script — Chapters: {script.chapters.join(", ")}
                    </span>
                </div>
                <Button
                    variant="outline"
                    size="icon-sm"
                    onClick={onClose}
                    aria-label="Close script viewer"
                >
                    {CloseIcon}
                </Button>
            </div>
            <div className="flex-1 overflow-y-auto p-6 pb-24 max-sm:p-4 max-sm:pb-24">
                <div className="mx-auto max-w-3xl">
                    <p className="whitespace-pre-wrap text-sm leading-relaxed text-foreground">
                        {script.content}
                    </p>
                </div>
            </div>
        </div>
    );
});

// ── Media status pill ───────────────────────────────────────────────────────

function MediaPill({
    url,
    type,
    onPlay,
}: {
    url: string | null;
    type: "audio" | "video";
    onPlay?: () => void;
}) {
    if (!url) {
        return (
            <span className="inline-flex items-center gap-1 rounded-[3px_5px_4px_3px] bg-neutral-100 px-1.5 py-0.5 text-[10px] font-medium text-neutral-400">
                {type === "audio" ? AudioIcon : VideoIcon}
                None
            </span>
        );
    }

    return (
        <button
            onClick={onPlay}
            className="inline-flex cursor-pointer items-center gap-1 rounded-[3px_5px_4px_3px] bg-foreground px-1.5 py-0.5 text-[10px] font-medium text-background transition-opacity hover:opacity-80"
        >
            {PlayIcon}
            {type === "audio" ? "Audio" : "Video"}
        </button>
    );
}

// ── Audio player overlay ────────────────────────────────────────────────────

const AudioPlayer = memo(function AudioPlayer({
    url,
    label,
    onClose,
}: {
    url: string;
    label: string;
    onClose: () => void;
}) {
    return (
        <div className="fixed inset-x-0 bottom-20 z-40 mx-auto max-w-lg animate-in slide-in-from-bottom-4 duration-200">
            <div className="mx-4 rounded-[6px_8px_7px_5px] border border-border bg-white p-3 shadow-[3px_5px_14px_rgba(0,0,0,0.12)]">
                <div className="mb-2 flex items-center justify-between">
                    <span className="text-xs font-medium">{label}</span>
                    <Button
                        variant="ghost"
                        size="icon-sm"
                        onClick={onClose}
                        aria-label="Close player"
                    >
                        {CloseIcon}
                    </Button>
                </div>
                {/* eslint-disable-next-line jsx-a11y/media-has-caption */}
                <audio controls autoPlay className="w-full" src={url}>
                    Your browser does not support the audio element.
                </audio>
            </div>
        </div>
    );
});

// ── Split card ──────────────────────────────────────────────────────────────

const SplitCard = memo(function SplitCard({
    split,
    index,
    onViewImage,
    onPlayAudio,
    onPlayVideo,
    onGenerateAudio,
    onGenerateVideo,
}: {
    split: Split;
    index: number;
    onViewImage: (index: number) => void;
    onPlayAudio: (split: Split) => void;
    onPlayVideo: (split: Split) => void;
    onGenerateAudio: (split: Split) => void;
    onGenerateVideo: (split: Split) => void;
}) {
    return (
        <li className="group overflow-hidden rounded-[6px_8px_7px_5px] border border-border shadow-[3px_5px_14px_rgba(0,0,0,0.06)] transition-all duration-300 hover:shadow-[4px_7px_24px_rgba(0,0,0,0.15)] hover:-translate-y-0.5 animate-in fade-in-0 zoom-in-95">
            <div className="flex gap-4 p-4 max-sm:flex-col">
                {/* ── Panel image ── */}
                {split.panel?.url ? (
                    <div className="group/img relative h-32 w-24 shrink-0 overflow-hidden rounded-[5px_7px_6px_4px] border border-border max-sm:h-40 max-sm:w-full">
                        {/* eslint-disable-next-line @next/next/no-img-element */}
                        <img
                            src={split.panel.url}
                            alt={`Panel ${split.panel.panelNumber}`}
                            className="absolute inset-0 h-full w-full object-cover transition-transform duration-500 group-hover/img:scale-105"
                        />
                        {/* Hover overlay with View Big CTA */}
                        <div className="absolute inset-0 flex items-center justify-center bg-black/0 opacity-0 transition-all duration-300 group-hover/img:bg-black/40 group-hover/img:opacity-100">
                            <Button
                                variant="secondary"
                                size="sm"
                                onClick={() => onViewImage(index)}
                                className="gap-1.5 bg-white text-black shadow-lg hover:bg-neutral-100"
                            >
                                {ExpandIcon}
                                View
                            </Button>
                        </div>
                        <div className="absolute bottom-1 left-1 rounded-[3px_5px_4px_3px] bg-black/70 px-1.5 py-0.5 font-mono text-[10px] font-medium text-white backdrop-blur-sm">
                            P{split.panel.panelNumber}
                        </div>
                    </div>
                ) : null}

                {/* ── Content ── */}
                <div className="flex min-w-0 flex-1 flex-col justify-between">
                    <div>
                        <div className="mb-1.5 flex flex-wrap items-center gap-2">
                            <span className="text-sm font-bold tracking-tight">
                                Split {String(index + 1).padStart(2, "0")}
                            </span>
                            <span className="rounded-[3px_5px_4px_3px] bg-foreground/5 px-1.5 py-0.5 font-mono text-[10px] text-muted-foreground">
                                ID {split.id}
                            </span>
                            <span className="rounded-[3px_5px_4px_3px] bg-foreground px-1.5 py-0.5 text-[10px] font-medium text-background">
                                {EFFECT_LABELS[split.effect] ?? split.effect}
                            </span>
                        </div>
                        <p className="text-sm leading-relaxed text-muted-foreground">
                            {split.content}
                        </p>
                    </div>

                    {/* ── Media row ── */}
                    <div className="mt-3 flex items-center gap-2">
                        <MediaPill
                            url={split.audioURL}
                            type="audio"
                            onPlay={() => onPlayAudio(split)}
                        />
                        <MediaPill
                            url={split.videoURL}
                            type="video"
                            onPlay={() => onPlayVideo(split)}
                        />

                        {/* Per-split actions */}
                        <DropdownMenu>
                            <DropdownMenuTrigger asChild>
                                <Button
                                    variant="ghost"
                                    size="icon-sm"
                                    className="ml-auto opacity-0 transition-opacity group-hover:opacity-100"
                                    aria-label="Split actions"
                                >
                                    {MoreVerticalIcon}
                                </Button>
                            </DropdownMenuTrigger>
                            <DropdownMenuContent align="end">
                                <DropdownMenuLabel>Generate</DropdownMenuLabel>
                                <DropdownMenuItem
                                    onClick={() => onGenerateAudio(split)}
                                >
                                    {AudioIcon}
                                    {split.audioURL
                                        ? "Regenerate Audio"
                                        : "Generate Audio"}
                                </DropdownMenuItem>
                                <DropdownMenuItem
                                    onClick={() => onGenerateVideo(split)}
                                    disabled={!split.audioURL}
                                >
                                    {VideoIcon}
                                    {split.videoURL
                                        ? "Regenerate Video"
                                        : "Generate Video"}
                                    {!split.audioURL ? (
                                        <span className="ml-auto text-[10px] text-neutral-400">
                                            needs audio
                                        </span>
                                    ) : null}
                                </DropdownMenuItem>
                            </DropdownMenuContent>
                        </DropdownMenu>
                    </div>
                </div>
            </div>
        </li>
    );
});

// ── Splits list content (rerender-memo) ─────────────────────────────────────

const SplitsListContent = memo(function SplitsListContent({
    isLoading,
    fetchError,
    splits,
    onViewImage,
    onPlayAudio,
    onPlayVideo,
    onGenerateAudio,
    onGenerateVideo,
}: {
    isLoading: boolean;
    fetchError: Error | null;
    splits: Split[];
    onViewImage: (index: number) => void;
    onPlayAudio: (split: Split) => void;
    onPlayVideo: (split: Split) => void;
    onGenerateAudio: (split: Split) => void;
    onGenerateVideo: (split: Split) => void;
}) {
    if (isLoading) {
        return (
            <div className="flex items-center justify-center gap-2 py-8 text-sm text-muted-foreground">
                <span className="h-4 w-4 animate-spin rounded-full border-2 border-border border-t-foreground" />
                Loading splits…
            </div>
        );
    }

    if (fetchError) {
        return (
            <div className="py-8 text-center text-sm font-medium text-foreground">
                Failed to load splits —{" "}
                {fetchError instanceof ApiError
                    ? fetchError.message
                    : "network error"}
            </div>
        );
    }

    return splits.length === 0 ? (
        <div className="py-12 text-center text-muted-foreground">
            {SplitEmptyIcon}
            <p className="text-sm font-medium text-foreground">
                No splits yet.
            </p>
            <span className="text-xs">
                Generate splits for this script using the button above.
            </span>
        </div>
    ) : (
        <ul className="flex flex-col gap-3">
            {splits.map((split, idx) => (
                <SplitCard
                    key={split.id}
                    split={split}
                    index={idx}
                    onViewImage={onViewImage}
                    onPlayAudio={onPlayAudio}
                    onPlayVideo={onPlayVideo}
                    onGenerateAudio={onGenerateAudio}
                    onGenerateVideo={onGenerateVideo}
                />
            ))}
        </ul>
    );
});

// ── Main page ───────────────────────────────────────────────────────────────

export default function SplitsPage() {
    const params = useParams();
    const bookId = Number(params.id);
    const scriptId = Number(params.scriptId);
    const queryClient = useQueryClient();

    // -- state
    const [confirmClear, setConfirmClear] = useState(false);
    const [confirmVideos, setConfirmVideos] = useState(false);
    const [lightboxIdx, setLightboxIdx] = useState<number | null>(null);
    const [viewingScript, setViewingScript] = useState(false);
    const [audioPlayerState, setAudioPlayerState] = useState<{
        url: string;
        label: string;
    } | null>(null);
    const [videoPlayerState, setVideoPlayerState] = useState<{
        url: string;
        label: string;
    } | null>(null);

    // Voice dialog state
    const [voiceDialogOpen, setVoiceDialogOpen] = useState(false);
    const [voiceDialogMode, setVoiceDialogMode] = useState<
        { type: "all" } | { type: "split"; splitId: number; splitLabel: string }
    >({ type: "all" });

    // -- stable callbacks
    const handleClearOpenChange = useCallback((open: boolean) => {
        if (!open) setConfirmClear(false);
    }, []);
    const handleConfirmVideosOpenChange = useCallback((open: boolean) => {
        if (!open) setConfirmVideos(false);
    }, []);
    const handleViewImage = useCallback(
        (index: number) => setLightboxIdx(index),
        [],
    );
    const closeLightbox = useCallback(() => setLightboxIdx(null), []);
    const closeScriptViewer = useCallback(() => setViewingScript(false), []);
    const closeAudioPlayer = useCallback(() => setAudioPlayerState(null), []);
    const closeVideoPlayer = useCallback(() => setVideoPlayerState(null), []);

    // -- fetch book info
    const { data: booksData } = useQuery({
        queryKey: ["books"],
        queryFn: () => fetchBooks(),
    });
    const book: Book | undefined = booksData?.data?.books?.find(
        (b) => b.id === bookId,
    );

    // -- fetch script info
    const { data: scriptsData } = useQuery({
        queryKey: ["scripts", bookId],
        queryFn: () => fetchScripts(bookId),
        enabled: !isNaN(bookId) && bookId > 0,
    });
    const script: Script | undefined = (scriptsData?.data?.scripts ?? []).find(
        (s) => s.id === scriptId,
    );

    // -- fetch splits
    const {
        data: splitsData,
        isLoading,
        error: fetchError,
    } = useQuery({
        queryKey: ["splits", scriptId],
        queryFn: () => fetchSplits(scriptId),
        enabled: !isNaN(scriptId) && scriptId > 0,
    });
    const splits: Split[] = splitsData?.data?.splits ?? [];
    const hasSplits = splits.length > 0;

    // -- derived: check if any split has audio
    const hasAnyAudio = useMemo(() => splits.some((s) => s.audioURL), [splits]);

    // -- lightbox items from splits
    const lightboxItems: LightboxItem[] = useMemo(
        () =>
            splits
                .filter((s) => s.panel?.url)
                .map((s, i) => ({
                    url: s.panel.url,
                    label: `Split ${String(i + 1).padStart(2, "0")} — P${s.panel.panelNumber}`,
                    description: s.content,
                    tag: EFFECT_LABELS[s.effect] ?? s.effect,
                    audioURL: s.audioURL,
                    videoURL: s.videoURL,
                    splitId: s.id,
                })),
        [splits],
    );

    // -- lightbox callbacks
    const handleLightboxGenerateAudio = useCallback((splitId: number) => {
        setVoiceDialogMode({
            type: "split",
            splitId,
            splitLabel: `Split ${splitId}`,
        });
        setVoiceDialogOpen(true);
    }, []);

    const handleLightboxGenerateVideo = useCallback((splitId: number) => {
        toast.info("Generating video for split…", { duration: 3000 });
        generateSplitVideo(splitId).catch((err) => {
            toast.error(
                err instanceof ApiError
                    ? err.businessError
                    : "Video generation failed",
            );
        });
    }, []);

    // -- delete splits mutation
    const clearMutation = useMutation({
        mutationFn: () => deleteSplits(scriptId),
        onSuccess: () => {
            setConfirmClear(false);
            queryClient.invalidateQueries({ queryKey: ["splits", scriptId] });
            toast.success("Splits cleared");
        },
        onError: (err) => {
            setConfirmClear(false);
            toast.error(
                err instanceof ApiError ? err.businessError : "Clear failed",
            );
        },
    });

    const handleConfirmClear = useCallback(() => {
        clearMutation.mutate();
    }, [clearMutation]);

    // -- fire-and-forget: generate splits
    const handleGenerateSplits = useCallback(() => {
        toast.info(
            `Generating splits for "${script?.name ?? `Script #${scriptId}`}"…`,
            { duration: 3000 },
        );
        generateSplits(scriptId).catch((err) => {
            toast.error(
                err instanceof ApiError
                    ? err.isValidationError
                        ? err.validationErrors.map((v) => v.message).join(", ")
                        : err.businessError
                    : "Split generation failed — please retry",
            );
        });
    }, [scriptId, script?.name]);

    // -- fire-and-forget: generate all audios
    const handleGenerateAllAudios = useCallback(
        (voice: string, voiceStyle: string) => {
            toast.info("Generating audio for all splits…", { duration: 3000 });
            generateAllAudios({
                scriptId,
                voice,
                voiceStyle: voiceStyle || undefined,
            }).catch((err) => {
                toast.error(
                    err instanceof ApiError
                        ? err.businessError
                        : "Audio generation failed",
                );
            });
        },
        [scriptId],
    );

    // -- fire-and-forget: generate all videos
    const handleGenerateAllVideos = useCallback(() => {
        setConfirmVideos(false);
        toast.info("Generating video for all splits with audio…", {
            duration: 3000,
        });
        generateAllVideos(scriptId).catch((err) => {
            toast.error(
                err instanceof ApiError
                    ? err.businessError
                    : "Video generation failed",
            );
        });
    }, [scriptId]);

    // -- fire-and-forget: generate single split audio
    const handleGenerateSplitAudio = useCallback(
        (voice: string, voiceStyle: string, splitId: number) => {
            toast.info("Generating audio for split…", { duration: 3000 });
            generateSplitAudio({
                splitId,
                voice,
                voiceStyle: voiceStyle || undefined,
            }).catch((err) => {
                toast.error(
                    err instanceof ApiError
                        ? err.businessError
                        : "Audio generation failed",
                );
            });
        },
        [],
    );

    // -- fire-and-forget: generate single split video
    const handleGenerateSplitVideo = useCallback((split: Split) => {
        if (!split.audioURL) {
            toast.error("Generate audio first before creating video");
            return;
        }
        toast.info("Generating video for split…", { duration: 3000 });
        generateSplitVideo(split.id).catch((err) => {
            toast.error(
                err instanceof ApiError
                    ? err.businessError
                    : "Video generation failed",
            );
        });
    }, []);

    // -- media player callbacks
    const handlePlayAudio = useCallback((split: Split) => {
        if (!split.audioURL) return;
        setAudioPlayerState({
            url: split.audioURL,
            label: `Split ${split.id} — Audio`,
        });
    }, []);

    const handlePlayVideo = useCallback((split: Split) => {
        if (!split.videoURL) return;
        setVideoPlayerState({
            url: split.videoURL,
            label: `Split ${split.id} — Video`,
        });
    }, []);

    // -- voice dialog submit handler
    const handleVoiceSubmit = useCallback(
        (voice: string, voiceStyle: string) => {
            if (voiceDialogMode.type === "all") {
                handleGenerateAllAudios(voice, voiceStyle);
            } else {
                handleGenerateSplitAudio(
                    voice,
                    voiceStyle,
                    voiceDialogMode.splitId,
                );
            }
        },
        [voiceDialogMode, handleGenerateAllAudios, handleGenerateSplitAudio],
    );

    // -- open voice dialog for a split
    const handleOpenSplitAudioDialog = useCallback((split: Split) => {
        setVoiceDialogMode({
            type: "split",
            splitId: split.id,
            splitLabel: `Split ${split.id}`,
        });
        setVoiceDialogOpen(true);
    }, []);

    const bookTitle = book?.title ?? `Book #${bookId}`;
    const scriptName = script?.name ?? `Script #${scriptId}`;

    // count splits without audio for warning
    const splitsWithoutAudio = useMemo(
        () => splits.filter((s) => !s.audioURL).length,
        [splits],
    );

    return (
        <>
            {/* ── Panel image lightbox ── */}
            {lightboxIdx !== null ? (
                <Lightbox
                    items={lightboxItems}
                    currentIndex={lightboxIdx}
                    onIndexChange={setLightboxIdx}
                    onClose={closeLightbox}
                    onGenerateAudio={handleLightboxGenerateAudio}
                    onGenerateVideo={handleLightboxGenerateVideo}
                />
            ) : null}

            {/* ── Full script viewer ── */}
            {viewingScript && script ? (
                <ScriptViewer script={script} onClose={closeScriptViewer} />
            ) : null}

            {/* ── Audio player ── */}
            {audioPlayerState ? (
                <AudioPlayer
                    url={audioPlayerState.url}
                    label={audioPlayerState.label}
                    onClose={closeAudioPlayer}
                />
            ) : null}

            {/* ── Video player ── */}
            {videoPlayerState ? (
                <VideoPlayerOverlayWrapper
                    url={videoPlayerState.url}
                    label={videoPlayerState.label}
                    onClose={closeVideoPlayer}
                />
            ) : null}

            {/* ── Voice input dialog ── */}
            <VoiceDialog
                open={voiceDialogOpen}
                onOpenChange={setVoiceDialogOpen}
                title={
                    voiceDialogMode.type === "all"
                        ? "Generate Audio for All Splits"
                        : `Generate Audio — ${voiceDialogMode.splitLabel}`
                }
                description={
                    voiceDialogMode.type === "all"
                        ? "Choose a voice and optional style for all splits."
                        : "Choose a voice and optional style for this split."
                }
                onSubmit={handleVoiceSubmit}
            />

            {/* ── Clear confirmation ── */}
            <AlertDialog
                open={confirmClear}
                onOpenChange={handleClearOpenChange}
            >
                <AlertDialogContent>
                    <AlertDialogHeader>
                        <AlertDialogTitle>Clear all splits?</AlertDialogTitle>
                        <AlertDialogDescription>
                            All {splits.length} split
                            {splits.length !== 1 ? "s" : ""} for &ldquo;
                            {scriptName}&rdquo; will be permanently removed. You
                            can regenerate them afterwards.
                        </AlertDialogDescription>
                    </AlertDialogHeader>
                    <AlertDialogFooter>
                        <AlertDialogCancel>Cancel</AlertDialogCancel>
                        <AlertDialogAction onClick={handleConfirmClear}>
                            Clear Splits
                        </AlertDialogAction>
                    </AlertDialogFooter>
                </AlertDialogContent>
            </AlertDialog>

            {/* ── Generate all videos confirmation with warning ── */}
            <AlertDialog
                open={confirmVideos}
                onOpenChange={handleConfirmVideosOpenChange}
            >
                <AlertDialogContent>
                    <AlertDialogHeader>
                        <AlertDialogTitle className="flex items-center gap-2">
                            {AlertTriangleIcon}
                            Generate Videos for All Splits?
                        </AlertDialogTitle>
                        <AlertDialogDescription asChild>
                            <div className="space-y-2">
                                <p>
                                    This will queue video generation for all
                                    splits that have audio.
                                </p>
                                {splitsWithoutAudio > 0 ? (
                                    <p className="rounded-[4px_6px_5px_3px] border border-neutral-200 bg-neutral-50 px-3 py-2 text-xs font-medium text-foreground">
                                        {splitsWithoutAudio} split
                                        {splitsWithoutAudio !== 1
                                            ? "s"
                                            : ""}{" "}
                                        without audio will be skipped. Generate
                                        audio first for full coverage.
                                    </p>
                                ) : null}
                            </div>
                        </AlertDialogDescription>
                    </AlertDialogHeader>
                    <AlertDialogFooter>
                        <AlertDialogCancel>Cancel</AlertDialogCancel>
                        <AlertDialogAction onClick={handleGenerateAllVideos}>
                            Generate Videos
                        </AlertDialogAction>
                    </AlertDialogFooter>
                </AlertDialogContent>
            </AlertDialog>

            <div className="mx-auto max-w-5xl px-6 pb-20 max-sm:px-4">
                {/* ── Hero ── */}
                <header className="relative pb-10 pt-20 max-sm:pb-7 max-sm:pt-12">
                    <Link
                        href={`/books/${bookId}/scripts`}
                        className="mb-6 inline-flex items-center gap-1.5 text-sm text-muted-foreground transition-colors hover:text-foreground"
                    >
                        {ArrowLeftIcon}
                        {bookTitle} — Scripts
                    </Link>
                    <h1 className="text-[clamp(2.5rem,8vw,4.5rem)] font-black leading-[0.85] tracking-tighter max-sm:text-4xl">
                        {scriptName}
                    </h1>
                    {HeroUnderline}
                    <span className="mt-4 block text-[11px] font-medium uppercase tracking-[0.3em] text-muted-foreground">
                        SCRIPT SPLITS
                    </span>
                </header>

                {/* ── Actions ── */}
                <section className="border-t border-border pb-8 pt-10">
                    <div className="flex items-center justify-between gap-3">
                        <h2 className="text-2xl font-semibold tracking-tight">
                            Splits
                            {hasSplits ? (
                                <span className="ml-2.5 rounded-[4px_6px_5px_3px] bg-foreground px-2 py-0.5 text-xs font-medium text-background">
                                    {splits.length}
                                </span>
                            ) : null}
                        </h2>
                        <div className="flex items-center gap-2">
                            {hasSplits ? (
                                <>
                                    {/* Primary CTA: Generate Audios */}
                                    <Button
                                        onClick={() => {
                                            setVoiceDialogMode({ type: "all" });
                                            setVoiceDialogOpen(true);
                                        }}
                                        variant="outline"
                                        className="gap-1.5"
                                    >
                                        {AudioIcon}
                                        <span className="max-sm:hidden">
                                            Generate Audios
                                        </span>
                                        <span className="sm:hidden">Audio</span>
                                    </Button>

                                    {/* Primary CTA: Generate Videos */}
                                    <Button
                                        onClick={() => setConfirmVideos(true)}
                                        disabled={!hasAnyAudio}
                                        className="gap-1.5"
                                    >
                                        {VideoIcon}
                                        <span className="max-sm:hidden">
                                            Generate Videos
                                        </span>
                                        <span className="sm:hidden">Video</span>
                                    </Button>

                                    {/* Overflow menu for secondary actions */}
                                    <DropdownMenu>
                                        <DropdownMenuTrigger asChild>
                                            <Button
                                                variant="outline"
                                                size="icon-sm"
                                                aria-label="More actions"
                                            >
                                                {ChevronDownIcon}
                                            </Button>
                                        </DropdownMenuTrigger>
                                        <DropdownMenuContent align="end">
                                            {script ? (
                                                <DropdownMenuItem
                                                    onClick={() =>
                                                        setViewingScript(true)
                                                    }
                                                >
                                                    {EyeIcon}
                                                    View Script
                                                </DropdownMenuItem>
                                            ) : null}
                                            <DropdownMenuItem
                                                onClick={handleGenerateSplits}
                                            >
                                                {RefreshIcon}
                                                Regenerate Splits
                                            </DropdownMenuItem>
                                            <DropdownMenuSeparator />
                                            <DropdownMenuItem
                                                onClick={() =>
                                                    setConfirmClear(true)
                                                }
                                                disabled={
                                                    clearMutation.isPending
                                                }
                                            >
                                                {TrashIcon}
                                                Clear Splits
                                            </DropdownMenuItem>
                                        </DropdownMenuContent>
                                    </DropdownMenu>
                                </>
                            ) : (
                                <>
                                    {script ? (
                                        <Button
                                            variant="outline"
                                            onClick={() =>
                                                setViewingScript(true)
                                            }
                                            className="gap-1.5"
                                        >
                                            {EyeIcon}
                                            View Script
                                        </Button>
                                    ) : null}
                                    <Button
                                        onClick={handleGenerateSplits}
                                        disabled={isLoading}
                                        className="gap-1.5"
                                    >
                                        {SparklesIcon}
                                        Generate Splits
                                    </Button>
                                </>
                            )}
                        </div>
                    </div>
                </section>

                {/* ── Splits List ── */}
                <section className="border-t border-border pt-8">
                    <SplitsListContent
                        isLoading={isLoading}
                        fetchError={fetchError}
                        splits={splits}
                        onViewImage={handleViewImage}
                        onPlayAudio={handlePlayAudio}
                        onPlayVideo={handlePlayVideo}
                        onGenerateAudio={handleOpenSplitAudioDialog}
                        onGenerateVideo={handleGenerateSplitVideo}
                    />
                </section>
            </div>
        </>
    );
}
