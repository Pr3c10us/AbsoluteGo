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
