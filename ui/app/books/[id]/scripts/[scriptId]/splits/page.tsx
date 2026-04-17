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
    Film,
    RectangleHorizontal,
    RectangleVertical,
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
    createVAB,
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

const FilmIcon = <Film className="h-4 w-4" />;

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

// ── Resolution presets ────────────────────────────────────────────────────────

interface ResolutionPreset {
    label: string;
    width: number;
    height: number;
    orientation: "horizontal" | "vertical";
}

const RESOLUTION_PRESETS: ResolutionPreset[] = [
    // Horizontal
    { label: "1920 × 1080 (Full HD)", width: 1920, height: 1080, orientation: "horizontal" },
    { label: "2560 × 1440 (2K QHD)", width: 2560, height: 1440, orientation: "horizontal" },
    { label: "3840 × 2160 (4K UHD)", width: 3840, height: 2160, orientation: "horizontal" },
    { label: "1280 × 720 (HD)", width: 1280, height: 720, orientation: "horizontal" },
    // Vertical
    { label: "1080 × 1920 (Full HD)", width: 1080, height: 1920, orientation: "vertical" },
    { label: "1440 × 2560 (2K QHD)", width: 1440, height: 2560, orientation: "vertical" },
    { label: "2160 × 3840 (4K UHD)", width: 2160, height: 3840, orientation: "vertical" },
    { label: "720 × 1280 (HD)", width: 720, height: 1280, orientation: "vertical" },
];

const HORIZONTAL_PRESETS = RESOLUTION_PRESETS.filter((p) => p.orientation === "horizontal");
const VERTICAL_PRESETS = RESOLUTION_PRESETS.filter((p) => p.orientation === "vertical");

const HorizontalIcon = <RectangleHorizontal className="h-3.5 w-3.5" />;
const VerticalIcon = <RectangleVertical className="h-3.5 w-3.5" />;

// ── Video settings dialog for width/height/FPS ──────────────────────────────

const PRESET_CUSTOM = "custom";

const VideoSettingsDialog = memo(function VideoSettingsDialog({
    open,
    onOpenChange,
    title,
    description,
    warningMessage,
    onSubmit,
}: {
    open: boolean;
    onOpenChange: (open: boolean) => void;
    title: string;
    description: string;
    /** Optional warning shown above the inputs (e.g. missing-audio notice). */
    warningMessage?: string | null;
    onSubmit: (settings: {
        width?: number;
        height?: number;
        FPS?: number;
    }) => void;
}) {
    const [preset, setPreset] = useState("");
    const [width, setWidth] = useState("");
    const [height, setHeight] = useState("");
    const [fps, setFps] = useState("");

    const isCustom = preset === PRESET_CUSTOM;

    const handlePresetChange = useCallback((value: string) => {
        setPreset(value);
        if (value === PRESET_CUSTOM) {
            setWidth("");
            setHeight("");
            return;
        }
        const found = RESOLUTION_PRESETS.find(
            (p) => `${p.width}x${p.height}` === value,
        );
        if (found) {
            setWidth(String(found.width));
            setHeight(String(found.height));
        }
    }, []);

    const handleSubmit = useCallback(() => {
        const w = width.trim() ? parseInt(width.trim(), 10) : undefined;
        const h = height.trim() ? parseInt(height.trim(), 10) : undefined;
        const f = fps.trim() ? parseInt(fps.trim(), 10) : undefined;

        if (
            (w !== undefined && (isNaN(w) || w <= 0)) ||
            (h !== undefined && (isNaN(h) || h <= 0)) ||
            (f !== undefined && (isNaN(f) || f <= 0))
        ) {
            toast.error("Values must be positive numbers");
            return;
        }

        onSubmit({ width: w, height: h, FPS: f });
        onOpenChange(false);
        setPreset("");
        setWidth("");
        setHeight("");
        setFps("");
    }, [width, height, fps, onSubmit, onOpenChange]);

    return (
        <Dialog open={open} onOpenChange={onOpenChange}>
            <DialogContent>
                <DialogHeader>
                    <DialogTitle>{title}</DialogTitle>
                    <DialogDescription>{description}</DialogDescription>
                </DialogHeader>
                <div className="flex flex-col gap-3 py-2">
                    {warningMessage ? (
                        <div className="flex items-start gap-2 rounded-[4px_6px_5px_3px] border border-neutral-200 bg-neutral-50 px-3 py-2.5">
                            {AlertTriangleIcon}
                            <p className="text-xs font-medium text-foreground">
                                {warningMessage}
                            </p>
                        </div>
                    ) : null}
                    <p className="text-xs text-muted-foreground">
                        Leave blank to use defaults: 1920 × 1080 @ 30 FPS
                    </p>

                    {/* Resolution preset */}
                    <div>
                        <label className="mb-1.5 block text-xs font-medium uppercase tracking-wider text-muted-foreground">
                            Resolution
                        </label>
                        <Select value={preset} onValueChange={handlePresetChange}>
                            <SelectTrigger className="w-full">
                                <SelectValue placeholder="Default (1920 × 1080)" />
                            </SelectTrigger>
                            <SelectContent
                                position="popper"
                                className="max-h-64"
                            >
                                <SelectGroup>
                                    <SelectLabel className="flex items-center gap-1.5">
                                        {HorizontalIcon}
                                        Horizontal
                                    </SelectLabel>
                                    {HORIZONTAL_PRESETS.map((p) => (
                                        <SelectItem
                                            key={`${p.width}x${p.height}`}
                                            value={`${p.width}x${p.height}`}
                                        >
                                            <span className="inline-flex items-center gap-1.5">
                                                {HorizontalIcon}
                                                {p.label}
                                            </span>
                                        </SelectItem>
                                    ))}
                                </SelectGroup>
                                <SelectGroup>
                                    <SelectLabel className="flex items-center gap-1.5">
                                        {VerticalIcon}
                                        Vertical
                                    </SelectLabel>
                                    {VERTICAL_PRESETS.map((p) => (
                                        <SelectItem
                                            key={`${p.width}x${p.height}`}
                                            value={`${p.width}x${p.height}`}
                                        >
                                            <span className="inline-flex items-center gap-1.5">
                                                {VerticalIcon}
                                                {p.label}
                                            </span>
                                        </SelectItem>
                                    ))}
                                </SelectGroup>
                                <SelectGroup>
                                    <SelectItem value={PRESET_CUSTOM}>
                                        Custom
                                    </SelectItem>
                                </SelectGroup>
                            </SelectContent>
                        </Select>
                    </div>

                    {/* Custom width / height (only when "Custom" selected) */}
                    {isCustom ? (
                        <div className="grid grid-cols-2 gap-3">
                            <div>
                                <label className="mb-1.5 block text-xs font-medium uppercase tracking-wider text-muted-foreground">
                                    Width
                                </label>
                                <Input
                                    type="number"
                                    value={width}
                                    onChange={(e) => setWidth(e.target.value)}
                                    placeholder="1920"
                                    min={1}
                                />
                            </div>
                            <div>
                                <label className="mb-1.5 block text-xs font-medium uppercase tracking-wider text-muted-foreground">
                                    Height
                                </label>
                                <Input
                                    type="number"
                                    value={height}
                                    onChange={(e) => setHeight(e.target.value)}
                                    placeholder="1080"
                                    min={1}
                                />
                            </div>
                        </div>
                    ) : null}

                    {/* FPS (always visible) */}
                    <div>
                        <label className="mb-1.5 block text-xs font-medium uppercase tracking-wider text-muted-foreground">
                            FPS
                        </label>
                        <Input
                            type="number"
                            value={fps}
                            onChange={(e) => setFps(e.target.value)}
                            placeholder="30"
                            min={1}
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
                        {VideoIcon}
                        Generate
                    </Button>
                </DialogFooter>
            </DialogContent>
        </Dialog>
    );
});

// ── VAB name dialog ─────────────────────────────────────────────────────────

const VabNameDialog = memo(function VabNameDialog({
    open,
    onOpenChange,
    onSubmit,
    hasVideoWarning,
    missingCount,
    totalCount,
}: {
    open: boolean;
    onOpenChange: (open: boolean) => void;
    onSubmit: (name: string) => void;
    hasVideoWarning: boolean;
    missingCount: number;
    totalCount: number;
}) {
    const [name, setName] = useState("");

    const handleSubmit = useCallback(() => {
        const trimmed = name.trim();
        if (!trimmed) {
            toast.error("Enter a name for the video audiobook");
            return;
        }
        onSubmit(trimmed);
        onOpenChange(false);
        setName("");
    }, [name, onSubmit, onOpenChange]);

    const handleKeyDown = useCallback(
        (e: React.KeyboardEvent) => {
            if (e.key === "Enter") handleSubmit();
        },
        [handleSubmit],
    );

    return (
        <Dialog open={open} onOpenChange={onOpenChange}>
            <DialogContent>
                <DialogHeader>
                    <DialogTitle className="flex items-center gap-2">
                        {FilmIcon}
                        Create Video Audiobook
                    </DialogTitle>
                    <DialogDescription>
                        All split videos will be merged into a single video audiobook.
                    </DialogDescription>
                </DialogHeader>
                <div className="flex flex-col gap-3 py-2">
                    {hasVideoWarning ? (
                        <div className="flex items-start gap-2 rounded-[4px_6px_5px_3px] border border-neutral-200 bg-neutral-50 px-3 py-2.5">
                            {AlertTriangleIcon}
                            <p className="text-xs font-medium text-foreground">
                                {missingCount} of {totalCount} splits are missing
                                videos. Only splits with videos will be included
                                in the final audiobook.
                            </p>
                        </div>
                    ) : null}
                    <div>
                        <label className="mb-1.5 block text-xs font-medium uppercase tracking-wider text-muted-foreground">
                            Name
                        </label>
                        <Input
                            value={name}
                            onChange={(e) => setName(e.target.value)}
                            onKeyDown={handleKeyDown}
                            placeholder="e.g. Chapter 1 — Final Cut"
                            autoFocus
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
                        {FilmIcon}
                        Create
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
                            {/*<span className="rounded-[3px_5px_4px_3px] bg-foreground px-1.5 py-0.5 text-[10px] font-medium text-background">*/}
                            {/*    {EFFECT_LABELS[split.effect] ?? split.effect}*/}
                            {/*</span>*/}
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

    // VAB dialog state
    const [vabDialogOpen, setVabDialogOpen] = useState(false);

    // Voice dialog state
    const [voiceDialogOpen, setVoiceDialogOpen] = useState(false);
    const [voiceDialogMode, setVoiceDialogMode] = useState<
        { type: "all" } | { type: "split"; splitId: number; splitLabel: string }
    >({ type: "all" });

    // Video settings dialog state
    const [videoDialogOpen, setVideoDialogOpen] = useState(false);
    const [videoDialogMode, setVideoDialogMode] = useState<
        { type: "all" } | { type: "split"; splitId: number }
    >({ type: "all" });

    // -- stable callbacks
    const handleClearOpenChange = useCallback((open: boolean) => {
        if (!open) setConfirmClear(false);
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
    // NOTE: uses "books-all" key to avoid colliding with the useInfiniteQuery
    // on the home page that uses ["books", search].
    const { data: booksData } = useQuery({
        queryKey: ["books-all"],
        queryFn: () => fetchBooks({ page: 1, limit: 500 }),
    });
    const book: Book | undefined = booksData?.data?.books?.find(
        (b) => b.id === bookId,
    );

    // -- fetch script info
    // NOTE: uses "scripts-all" key to avoid colliding with the useInfiniteQuery
    // on the scripts page that uses ["scripts", bookId].
    const { data: scriptsData } = useQuery({
        queryKey: ["scripts-all", bookId],
        queryFn: () => fetchScripts(bookId),
        enabled: !isNaN(bookId) && bookId > 0,
    });
    const script: Script | undefined = (() => {
        const list = scriptsData?.data?.scripts;
        return Array.isArray(list) ? list : [];
    })().find(
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

    // -- derived: audio coverage stats
    const splitsWithAudioCount = useMemo(
        () => splits.filter((s) => s.audioURL).length,
        [splits],
    );
    const allSplitsHaveAudio = hasSplits && splitsWithAudioCount === splits.length;

    // -- derived: video coverage stats
    const splitsWithVideo = useMemo(
        () => splits.filter((s) => s.videoURL).length,
        [splits],
    );
    const allSplitsHaveVideo = hasSplits && splitsWithVideo === splits.length;
    const missingVideoCount = splits.length - splitsWithVideo;

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
        setVideoDialogMode({ type: "split", splitId });
        setVideoDialogOpen(true);
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

    // -- fire-and-forget: create VAB
    const handleCreateVAB = useCallback(
        (name: string) => {
            toast.info(`Queuing video audiobook "${name}"…`, { duration: 3000 });
            createVAB({ scriptId, name }).catch((err) => {
                toast.error(
                    err instanceof ApiError
                        ? err.isValidationError
                            ? err.validationErrors.map((v) => v.message).join(", ")
                            : err.businessError
                        : "Failed to create video audiobook",
                );
            });
        },
        [scriptId],
    );

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
    const handleGenerateAllVideos = useCallback(
        (settings: { width?: number; height?: number; FPS?: number }) => {
            toast.info("Generating video for all splits with audio…", {
                duration: 3000,
            });
            generateAllVideos(scriptId, settings).catch((err) => {
                toast.error(
                    err instanceof ApiError
                        ? err.businessError
                        : "Video generation failed",
                );
            });
        },
        [scriptId],
    );

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

    // -- open video settings dialog for single split
    const handleGenerateSplitVideo = useCallback((split: Split) => {
        if (!split.audioURL) {
            toast.error("Generate audio first before creating video");
            return;
        }
        setVideoDialogMode({ type: "split", splitId: split.id });
        setVideoDialogOpen(true);
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

    // -- video settings dialog submit handler
    const handleVideoSettingsSubmit = useCallback(
        (settings: { width?: number; height?: number; FPS?: number }) => {
            if (videoDialogMode.type === "all") {
                handleGenerateAllVideos(settings);
            } else {
                toast.info("Generating video for split…", { duration: 3000 });
                generateSplitVideo(videoDialogMode.splitId, settings).catch(
                    (err) => {
                        toast.error(
                            err instanceof ApiError
                                ? err.businessError
                                : "Video generation failed",
                        );
                    },
                );
            }
        },
        [videoDialogMode, handleGenerateAllVideos],
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

            {/* ── VAB name dialog ── */}
            <VabNameDialog
                open={vabDialogOpen}
                onOpenChange={setVabDialogOpen}
                onSubmit={handleCreateVAB}
                hasVideoWarning={!allSplitsHaveVideo && hasSplits}
                missingCount={missingVideoCount}
                totalCount={splits.length}
            />

            {/* ── Voice input dialog ── */}
            <VoiceDialog
                open={voiceDialogOpen}
                onOpenChange={setVoiceDialogOpen}
                title={
                    voiceDialogMode.type === "all"
                        ? allSplitsHaveAudio
                            ? "Regenerate Audio for All Splits"
                            : "Generate Audio for All Splits"
                        : `Generate Audio — ${voiceDialogMode.splitLabel}`
                }
                description={
                    voiceDialogMode.type === "all"
                        ? allSplitsHaveAudio
                            ? "All splits already have audio. This will regenerate audio for all splits with the selected voice."
                            : "Choose a voice and optional style for all splits."
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

            {/* ── Video settings dialog ── */}
            <VideoSettingsDialog
                open={videoDialogOpen}
                onOpenChange={setVideoDialogOpen}
                title={
                    videoDialogMode.type === "all"
                        ? allSplitsHaveVideo
                            ? "Regenerate Videos — All Splits"
                            : "Generate Videos — All Splits"
                        : "Generate Video"
                }
                description={
                    videoDialogMode.type === "all"
                        ? allSplitsHaveVideo
                            ? "All splits already have video. This will regenerate video for all splits with the selected settings."
                            : "This will queue video generation for all splits that have audio. Configure resolution and frame rate below."
                        : "Configure resolution and frame rate for this split video."
                }
                warningMessage={
                    videoDialogMode.type === "all" && splitsWithoutAudio > 0
                        ? `${splitsWithoutAudio} split${splitsWithoutAudio !== 1 ? "s" : ""} without audio will be skipped. Generate audio first for full coverage.`
                        : null
                }
                onSubmit={handleVideoSettingsSubmit}
            />

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
                                    {/* Primary CTA: Create Video Audiobook */}
                                    <Button
                                        onClick={() => setVabDialogOpen(true)}
                                        disabled={splitsWithVideo === 0}
                                        className="gap-1.5"
                                    >
                                        {FilmIcon}
                                        <span className="max-sm:hidden">
                                            Video Audiobook
                                        </span>
                                        <span className="sm:hidden">VAB</span>
                                        {!allSplitsHaveVideo && splitsWithVideo > 0 ? (
                                            <span className="ml-0.5 flex h-4 w-4 items-center justify-center rounded-full bg-white/20 text-[9px] font-bold text-white/80">
                                                !
                                            </span>
                                        ) : null}
                                    </Button>

                                    {/* Overflow menu */}
                                    <DropdownMenu>
                                        <DropdownMenuTrigger asChild>
                                            <Button
                                                variant="outline"
                                                className="gap-1.5"
                                                aria-label="More actions"
                                            >
                                                Options
                                                {ChevronDownIcon}
                                            </Button>
                                        </DropdownMenuTrigger>
                                        <DropdownMenuContent align="end">
                                            <DropdownMenuLabel>Generate</DropdownMenuLabel>
                                            <DropdownMenuItem
                                                onClick={() => {
                                                    setVoiceDialogMode({ type: "all" });
                                                    setVoiceDialogOpen(true);
                                                }}
                                            >
                                                {allSplitsHaveAudio ? RefreshIcon : AudioIcon}
                                                {allSplitsHaveAudio ? "Regenerate Audios" : "Generate Audios"}
                                            </DropdownMenuItem>
                                            <DropdownMenuItem
                                                onClick={() => {
                                                    setVideoDialogMode({ type: "all" });
                                                    setVideoDialogOpen(true);
                                                }}
                                                disabled={!hasAnyAudio}
                                            >
                                                {allSplitsHaveVideo ? RefreshIcon : VideoIcon}
                                                {allSplitsHaveVideo ? "Regenerate Videos" : "Generate Videos"}
                                                {!hasAnyAudio ? (
                                                    <span className="ml-auto text-[10px] text-neutral-400">
                                                        needs audio
                                                    </span>
                                                ) : null}
                                            </DropdownMenuItem>
                                            <DropdownMenuSeparator />
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
                                            <DropdownMenuItem asChild>
                                                <Link
                                                    href={`/books/${bookId}/videos`}
                                                >
                                                    {FilmIcon}
                                                    View Videos
                                                </Link>
                                            </DropdownMenuItem>
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
                    {/* Video coverage warning */}
                    {hasSplits && !allSplitsHaveVideo && splitsWithVideo > 0 ? (
                        <div className="mb-6 flex items-start gap-2.5 rounded-[4px_6px_5px_3px] border border-neutral-200 bg-neutral-50 px-4 py-3">
                            {AlertTriangleIcon}
                            <div className="min-w-0">
                                <p className="text-sm font-medium text-foreground">
                                    {missingVideoCount} split
                                    {missingVideoCount !== 1 ? "s" : ""} missing
                                    videos
                                </p>
                                <p className="mt-0.5 text-xs text-muted-foreground">
                                    Generate videos for all splits before creating a
                                    video audiobook for full coverage.{" "}
                                    {splitsWithVideo} of {splits.length} ready.
                                </p>
                            </div>
                        </div>
                    ) : null}
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
