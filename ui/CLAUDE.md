# AbsoluteGo UI — Agent Rules

## Background Generation Processes

**All AI generation and long-running API calls MUST run as background tasks.**

When implementing any feature that calls an AI generation endpoint (script generation, split generation, or any future generation endpoint), you MUST:

1. **Use the background task context** (`useUpload` from `@/lib/upload-context`) instead of inline `useMutation` for the generation call.
2. **Fire-and-forget** — close any modal / reset form state immediately after dispatching the task. Do NOT block the UI waiting for the response.
3. **Add a new task method** to the `UploadProvider` in `lib/upload-context.tsx` if one doesn't exist for the operation. Follow the existing pattern:
   - Create a discriminated union member type (e.g. `interface FooTask { type: "foo"; ... }`)
   - Add it to the `BackgroundTask` union
   - Add the fire-and-forget method to the context value interface and implementation
   - Show a `toast.info` on start, `toast.success` on completion, `toast.error` on failure
   - Invalidate the relevant React Query cache key on success
   - Auto-remove the task from the list after a delay (5s success, 8s error)
4. **Update the tracker** (`components/upload-tracker.tsx`) — add the new task type to `getTaskLabel()` and `getTaskSubtitle()` so it renders correctly in the background tasks panel.
5. **Never use `useMutation`** for generation endpoints. `useMutation` is only for quick synchronous operations like delete.

### Rationale

Generation endpoints hit AI models and can take 10–60+ seconds. Blocking the UI with a spinner degrades UX. The background task pill (bottom corner) lets users navigate freely while generation runs, and shows progress/completion status.
