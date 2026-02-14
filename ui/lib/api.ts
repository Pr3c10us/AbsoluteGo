// ---------------------------------------------------------------------------
// Absolute API — centralised fetch utilities
// ---------------------------------------------------------------------------

const API_BASE_URL =
  process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:5000";

// ---- Types ----------------------------------------------------------------

export interface Book {
  id: number;
  title: string;
}

export interface GetBooksResponse {
  message: string;
  data: { books: Book[] };
}

export interface MutationResponse {
  message: string;
}

export interface ValidationField {
  field: string;
  tag: string;
  message: string;
}

export interface ValidationErrorBody {
  error: ValidationField[];
  message: string;
}

export interface BusinessErrorBody {
  error: string;
  message: string;
}

// ---- Error class ----------------------------------------------------------

export class ApiError extends Error {
  status: number;
  body: ValidationErrorBody | BusinessErrorBody;

  constructor(status: number, body: ValidationErrorBody | BusinessErrorBody) {
    super(body.message ?? "API Error");
    this.status = status;
    this.body = body;
  }

  get isValidationError(): boolean {
    return this.status === 422;
  }

  get validationErrors(): ValidationField[] {
    if (this.isValidationError && Array.isArray(this.body.error)) {
      return this.body.error as ValidationField[];
    }
    return [];
  }

  get businessError(): string {
    if (!this.isValidationError && typeof this.body.error === "string") {
      return this.body.error;
    }
    return this.message;
  }
}

// ---- Fetch wrapper --------------------------------------------------------

export async function apiFetch<T>(
  path: string,
  options?: RequestInit
): Promise<T> {
  const res = await fetch(`${API_BASE_URL}${path}`, {
    headers: { "Content-Type": "application/json", ...options?.headers },
    ...options,
  });

  const json = await res.json();

  if (!res.ok) {
    throw new ApiError(res.status, json);
  }

  return json as T;
}

// ---- Book API calls -------------------------------------------------------

export function fetchBooks(title?: string): Promise<GetBooksResponse> {
  const query = title ? `?title=${encodeURIComponent(title)}` : "";
  return apiFetch<GetBooksResponse>(`/api/v1/book${query}`);
}

export function addBook(title: string): Promise<MutationResponse> {
  return apiFetch<MutationResponse>("/api/v1/book", {
    method: "POST",
    body: JSON.stringify({ title }),
  });
}

export function deleteBook(bookId: number): Promise<MutationResponse> {
  return apiFetch<MutationResponse>(`/api/v1/book/${bookId}`, {
    method: "DELETE",
  });
}

// ---- Chapter types --------------------------------------------------------

export interface Chapter {
  id: number;
  number: number;
  bookId: number;
  blurURL: string;
}

export interface GetChaptersResponse {
  message: string;
  data: { chapters: Chapter[] };
}

// ---- Chapter API calls ----------------------------------------------------

export function fetchChapters(
  bookId: number,
  number?: number
): Promise<GetChaptersResponse> {
  const params = new URLSearchParams({ bookId: String(bookId) });
  if (number !== undefined) params.set("number", String(number));
  return apiFetch<GetChaptersResponse>(
    `/api/v1/book/chapter?${params.toString()}`
  );
}

export async function addChapter(
  bookId: number,
  chapterNumber: number,
  file: File
): Promise<MutationResponse> {
  const formData = new FormData();
  formData.append("bookId", String(bookId));
  formData.append("chapter", String(chapterNumber));
  formData.append("book", file);

  const res = await fetch(`${API_BASE_URL}/api/v1/book/chapter`, {
    method: "POST",
    body: formData,
  });

  const json = await res.json();
  if (!res.ok) throw new ApiError(res.status, json);
  return json as MutationResponse;
}

export function deleteChapter(
  chapterId: number
): Promise<MutationResponse> {
  return apiFetch<MutationResponse>(`/api/v1/book/chapter/${chapterId}`, {
    method: "DELETE",
  });
}

// ---- Page types -----------------------------------------------------------

export interface Page {
  id: number;
  chapterId: number;
  url: string;
  llmurl: string;
  mime: string;
  pageNumber: number;
  updatedAt: string;
}

export interface GetPagesResponse {
  message: string;
  data: { pages: Page[] };
}

// ---- Page API calls -------------------------------------------------------

export function fetchPages(chapterId: number): Promise<GetPagesResponse> {
  return apiFetch<GetPagesResponse>(
    `/api/v1/book/page?chapterId=${chapterId}`
  );
}

// ---- Panel types ----------------------------------------------------------

export interface Panel {
  id: number;
  pageId: number;
  url: string;
  panelNumber: number;
  updatedAt: string;
}

export interface GetPanelsResponse {
  message: string;
  data: { panels: Panel[] };
}

// ---- Panel API calls ------------------------------------------------------

export function fetchPanels(pageId: number): Promise<GetPanelsResponse> {
  return apiFetch<GetPanelsResponse>(
    `/api/v1/book/panel?pageId=${pageId}`
  );
}

// ---- Script types ---------------------------------------------------------

export interface Script {
  id: number;
  name: string;
  content: string;
  bookId: number;
  chapters: number[];
}

export interface GetScriptsResponse {
  message: string;
  data: { scripts: Script[] | null };
}

export interface GenerateScriptResponse {
  message: string;
  data: { script: string; scriptId: number };
}

// ---- Script API calls -----------------------------------------------------

export function fetchScripts(
  bookId: number,
  opts?: { name?: string; ids?: number[] }
): Promise<GetScriptsResponse> {
  const params = new URLSearchParams({ bookId: String(bookId) });
  if (opts?.name) params.set("name", opts.name);
  if (opts?.ids) {
    for (const id of opts.ids) params.append("id", String(id));
  }
  return apiFetch<GetScriptsResponse>(`/api/v1/script?${params.toString()}`);
}

export function generateScript(body: {
  bookId: number;
  name: string;
  chapters: number[];
  previousScripts?: number[];
}): Promise<GenerateScriptResponse> {
  return apiFetch<GenerateScriptResponse>("/api/v1/script", {
    method: "POST",
    body: JSON.stringify(body),
  });
}

export function deleteScript(scriptId: number): Promise<MutationResponse> {
  return apiFetch<MutationResponse>(`/api/v1/script/${scriptId}`, {
    method: "DELETE",
  });
}

// ---- Split types ----------------------------------------------------------

export interface SplitPanel {
  id: number;
  pageId: number;
  url: string;
  panelNumber: number;
  updatedAt: string;
}

export interface Split {
  id: number;
  scriptId: number;
  content: string;
  panelId: number;
  effect: string;
  panel: SplitPanel;
}

export interface GetSplitsResponse {
  message: string;
  data: { splits: Split[] | null };
}

// ---- Split API calls ------------------------------------------------------

export function fetchSplits(scriptId: number): Promise<GetSplitsResponse> {
  return apiFetch<GetSplitsResponse>(`/api/v1/script/split/${scriptId}`);
}

export function generateSplits(
  scriptId: number
): Promise<MutationResponse> {
  return apiFetch<MutationResponse>(`/api/v1/script/split/${scriptId}`, {
    method: "POST",
  });
}

export function deleteSplits(scriptId: number): Promise<MutationResponse> {
  return apiFetch<MutationResponse>(`/api/v1/script/split/${scriptId}`, {
    method: "DELETE",
  });
}

// ---- Event types ----------------------------------------------------------

export type EventStatus =
  | "enqueue"
  | "processing"
  | "failed"
  | "successful"
  | "retry";

export type EventOperation =
  | "add_chapter"
  | "gen_script"
  | "gen_script_split"
  | "gen_audio"
  | "gen_video"
  | "merge_video";

export interface EventItem {
  Id: number;
  Status: EventStatus;
  Operation: EventOperation;
  Description: string;
  BookId: number;
  ChapterId: number;
  ScriptId: number;
  VabId: number;
  UpdatedAt: string;
}

export interface GetEventsResponse {
  message: string;
  data: { events: EventItem[] | null };
}

// ---- Event API calls ------------------------------------------------------

export function fetchEvents(params?: {
  page?: number;
  limit?: number;
  status?: EventStatus;
  operation?: EventOperation;
}): Promise<GetEventsResponse> {
  const qs = new URLSearchParams();
  if (params?.page) qs.set("page", String(params.page));
  if (params?.limit) qs.set("limit", String(params.limit));
  if (params?.status) qs.set("status", params.status);
  if (params?.operation) qs.set("operation", params.operation);
  const query = qs.toString();
  return apiFetch<GetEventsResponse>(
    `/api/v1/event${query ? `?${query}` : ""}`
  );
}
