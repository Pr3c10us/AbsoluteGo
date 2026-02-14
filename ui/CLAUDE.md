# AbsoluteGo UI — Agent Rules

## Background Generation Processes

**All AI generation and long-running API calls are fire-and-forget — the backend queues them as events.**

When implementing any feature that calls an AI generation endpoint (chapter upload, script generation, split generation, audio/video generation, or any future generation endpoint), you MUST:

1. **Call the API function directly** (e.g. `addChapter()`, `generateScript()`, `generateSplits()` from `@/lib/api`) — there is no background task context.
2. **Fire-and-forget** — close any modal / reset form state immediately after dispatching the call. Do NOT block the UI waiting for the response.
3. **Show a `toast.info()`** when dispatching the call so the user knows it was sent.
4. **Catch errors** with `.catch()` and show `toast.error()` — the only client-side error handling needed.
5. **Never use `useMutation`** for generation endpoints. `useMutation` is only for quick synchronous operations like delete.

### Events Tracker

The `EventTracker` component (`components/event-tracker.tsx`) polls `GET /api/v1/event` every 5 seconds using `useQuery` and displays a draggable pill + panel at the bottom corner. It shows all server-side events with their status (queued, processing, failed, successful, retrying). Users can hover to preview events or click to pin the panel open.

**Do NOT** use `UploadProvider`, `useUpload`, or any client-side background task tracking. Those have been removed.

### Rationale

Generation endpoints hit AI models and can take 10–60+ seconds. The backend queues these as events and tracks their status server-side. The `EventTracker` polls for updates, keeping the UI responsive while showing real-time progress.
