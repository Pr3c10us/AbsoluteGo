"use client";

import { memo, useState, useCallback } from "react";
import Link from "next/link";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import { Trash2, BookOpen } from "lucide-react";
import {
  fetchBooks,
  addBook,
  deleteBook,
  ApiError,
  type Book,
  type ValidationField,
} from "@/lib/api";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
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

// ── Static icons (hoisted — rendering-hoist-jsx) ────────────────────────────

const TrashIcon = <Trash2 className="h-4 w-4" />;

const BookIcon = <BookOpen className="mx-auto mb-3 h-10 w-10 text-neutral-300" />;

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

// ── Book item row (memoised — rerender-memo) ────────────────────────────────

const BookItem = memo(function BookItem({
  book,
  onDelete,
  disabled,
}: {
  book: Book;
  onDelete: (book: Book) => void;
  disabled: boolean;
}) {
  return (
    <li className="flex items-center justify-between gap-3 rounded-lg border border-border px-4 py-3.5 shadow-[2px_4px_12px_rgba(0,0,0,0.04)] transition-shadow hover:shadow-[2px_4px_16px_rgba(0,0,0,0.08)] animate-in fade-in-0 slide-in-from-bottom-1">
      <Link
        href={`/books/${book.id}`}
        className="flex items-center gap-2.5 min-w-0 flex-1 group"
      >
        <span className="shrink-0 font-mono text-xs text-muted-foreground">
          #{book.id}
        </span>
        <span className="truncate text-sm font-medium group-hover:underline">
          {book.title}
        </span>
      </Link>
      <Button
        variant="outline"
        size="icon-sm"
        onClick={(e) => {
          e.preventDefault();
          onDelete(book);
        }}
        disabled={disabled}
        aria-label={`Delete ${book.title}`}
        className="shrink-0"
      >
        {TrashIcon}
      </Button>
    </li>
  );
});

// ── List content (rerender-memo, rendering-conditional-render) ─────────────

const ListContent = memo(function ListContent({
  isLoading,
  fetchError,
  books,
  onDelete,
  deleteDisabled,
}: {
  isLoading: boolean;
  fetchError: Error | null;
  books: Book[];
  onDelete: (book: Book) => void;
  deleteDisabled: boolean;
}) {
  if (isLoading) {
    return (
      <div className="flex items-center justify-center gap-2 py-8 text-sm text-muted-foreground">
        <span className="h-4 w-4 animate-spin rounded-full border-2 border-border border-t-foreground" />
        Loading…
      </div>
    );
  }

  if (fetchError) {
    return (
      <div className="py-8 text-center text-sm font-medium text-foreground">
        Failed to load books —{" "}
        {fetchError instanceof ApiError
          ? fetchError.message
          : "network error"}
      </div>
    );
  }

  return books.length === 0 ? (
    <div className="py-12 text-center text-muted-foreground">
      {BookIcon}
      <p className="text-sm font-medium text-foreground">No books yet.</p>
      <span className="text-xs">Add your first book above to get started.</span>
    </div>
  ) : (
    <ul className="flex flex-col gap-1.5">
      {books.map((book) => (
        <BookItem
          key={book.id}
          book={book}
          onDelete={onDelete}
          disabled={deleteDisabled}
        />
      ))}
    </ul>
  );
});

// ── Main page ───────────────────────────────────────────────────────────────

export default function HomePage() {
  const queryClient = useQueryClient();

  // -- state
  const [search, setSearch] = useState("");
  const [newTitle, setNewTitle] = useState("");
  const [validationErrors, setValidationErrors] = useState<ValidationField[]>(
    []
  );
  const [confirmBook, setConfirmBook] = useState<Book | null>(null);

  // -- stable callbacks (rerender-functional-setstate)
  const clearConfirm = useCallback(() => setConfirmBook(null), []);
  const handleOpenChange = useCallback(
    (open: boolean) => { if (!open) setConfirmBook(null); },
    []
  );
  const handleDeleteClick = useCallback(
    (book: Book) => setConfirmBook(book),
    []
  );

  // -- queries
  const {
    data,
    isLoading,
    error: fetchError,
  } = useQuery({
    queryKey: ["books", search],
    queryFn: () => fetchBooks(search || undefined),
  });

  const books = data?.data?.books ?? [];

  // -- mutations
  const addMutation = useMutation({
    mutationFn: (title: string) => addBook(title),
    onSuccess: (res) => {
      setNewTitle("");
      setValidationErrors([]);
      queryClient.invalidateQueries({ queryKey: ["books"] });
      toast.success(res.message || "Book created");
    },
    onError: (err) => {
      if (!(err instanceof ApiError)) {
        toast.error("Network error — please retry");
        return;
      }
      if (err.isValidationError) {
        setValidationErrors(err.validationErrors);
      } else {
        toast.error(err.businessError);
      }
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (bookId: number) => deleteBook(bookId),
    onSuccess: (res) => {
      clearConfirm();
      queryClient.invalidateQueries({ queryKey: ["books"] });
      toast.success(res.message || "Book deleted");
    },
    onError: (err) => {
      clearConfirm();
      toast.error(
        err instanceof ApiError ? err.businessError : "Delete failed"
      );
    },
  });

  const handleSubmit = useCallback(
    (e: React.FormEvent) => {
      e.preventDefault();
      setValidationErrors([]);
      const trimmed = newTitle.trim();
      if (!trimmed) return;
      addMutation.mutate(trimmed);
    },
    [newTitle, addMutation]
  );

  const handleConfirmDelete = useCallback(() => {
    if (confirmBook) {
      deleteMutation.mutate(confirmBook.id);
    }
  }, [confirmBook, deleteMutation]);

  // -- derived state (rerender-derived-state)
  const titleError = validationErrors.find(
    (v) => v.field.toLowerCase() === "title"
  );
  const isSubmitDisabled = addMutation.isPending || !newTitle.trim();

  return (
    <div className="mx-auto max-w-5xl px-6 pb-20 max-sm:px-4">
      {/* ── Delete confirmation (shadcn AlertDialog) ────────────── */}
      <AlertDialog
        open={confirmBook !== null}
        onOpenChange={handleOpenChange}
      >
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Delete book?</AlertDialogTitle>
            <AlertDialogDescription>
              &ldquo;{confirmBook?.title}&rdquo; and all its chapters, pages,
              and panels will be permanently removed.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <AlertDialogAction onClick={handleConfirmDelete}>
              Delete
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>

      {/* ── Hero ──────────────────────────────────────────────────── */}
      <header className="relative pb-10 pt-20 max-sm:pb-7 max-sm:pt-12">
        <h1 className="text-[clamp(3.5rem,10vw,6rem)] font-bold leading-[0.85] tracking-tighter max-sm:text-5xl">
          BOOKS
        </h1>
        {HeroUnderline}
        <span className="mt-4 block text-[11px] font-medium uppercase tracking-[0.3em] text-muted-foreground">
          LIBRARY MANAGEMENT
        </span>
      </header>

      {/* ── Add Book ──────────────────────────────────────────────── */}
      <section className="border-t border-border pb-8 pt-10">
        <h2 className="mb-5 flex items-center gap-2.5 text-2xl font-semibold tracking-tight">
          Add New Book
        </h2>
        <form
          onSubmit={handleSubmit}
          className="flex items-start gap-2.5 max-sm:flex-col"
        >
          <div className="min-w-0 flex-1 max-sm:w-full">
            <Input
              id="book-title-input"
              placeholder="Enter book title…"
              value={newTitle}
              onChange={(e) => {
                setNewTitle(e.target.value);
                setValidationErrors([]);
              }}
              disabled={addMutation.isPending}
              className={titleError ? "border-foreground" : ""}
            />
            {titleError !== undefined ? (
              <span className="mt-1.5 block text-xs font-medium text-foreground">
                {titleError.message}
              </span>
            ) : null}
          </div>
          <Button
            id="add-book-btn"
            type="submit"
            disabled={isSubmitDisabled}
            className="shrink-0 max-sm:w-full"
          >
            {addMutation.isPending ? "Adding…" : "Add Book"}
          </Button>
        </form>
      </section>

      {/* ── Books List ────────────────────────────────────────────── */}
      <section className="border-t border-border pt-8">
        <div className="mb-6 flex flex-wrap items-center justify-between gap-3 max-sm:flex-col max-sm:items-stretch">
          <h2 className="flex items-center gap-2.5 text-2xl font-semibold tracking-tight">
            Library
            {books.length > 0 ? (
              <span className="rounded-[4px_6px_5px_3px] bg-foreground px-2 py-0.5 text-xs font-medium text-background">
                {books.length}
              </span>
            ) : null}
          </h2>
          <Input
            id="search-input"
            placeholder="Filter by title…"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            className="max-w-60 max-sm:max-w-full"
          />
        </div>

        <ListContent
          isLoading={isLoading}
          fetchError={fetchError}
          books={books}
          onDelete={handleDeleteClick}
          deleteDisabled={deleteMutation.isPending}
        />
      </section>
    </div>
  );
}
