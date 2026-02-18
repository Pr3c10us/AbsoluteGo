PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS books
(
    id    INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS chapters
(
    id       INTEGER PRIMARY KEY AUTOINCREMENT,
    number   INTEGER NOT NULL,
    book_id  INTEGER NOT NULL,
    blur_url TEXT    NOT NULL,
    FOREIGN KEY (book_id) REFERENCES books (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS pages
(
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    chapter_id  INTEGER NOT NULL,
    url         TEXT,
    llm_url     TEXT,
    mime        TEXT,
    page_number INTEGER NOT NULL,
    updated_at  DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (chapter_id) REFERENCES chapters (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS panels
(
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    page_id      INTEGER NOT NULL,
    url          TEXT,
    panel_number INTEGER NOT NULL,
    updated_at   DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (page_id) REFERENCES pages (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS scripts
(
    id       INTEGER PRIMARY KEY AUTOINCREMENT,
    name     TEXT    NOT NULL,
    content  TEXT,
    book_id  INTEGER NOT NULL,
    chapters TEXT DEFAULT '[]',
    FOREIGN KEY (book_id) REFERENCES books (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS splits
(
    id               INTEGER PRIMARY KEY AUTOINCREMENT,
    script_id        INTEGER NOT NULL,
    content          TEXT,
    previous_content TEXT,
    panel_id         INTEGER,
    effect           TEXT,
    audio_url        TEXT,
    audio_duration   REAL,
    video_url        TEXT,
    FOREIGN KEY (script_id) REFERENCES scripts (id) ON DELETE CASCADE,
    FOREIGN KEY (panel_id) REFERENCES panels (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS vabs
(
    id        INTEGER PRIMARY KEY AUTOINCREMENT,
    name      TEXT    NOT NULL,
    url       TEXT,
    script_id INTEGER NOT NULL,
    book_id INTEGER NOT NULL,
    FOREIGN KEY (script_id) REFERENCES scripts (id) ON DELETE CASCADE,
    FOREIGN KEY (book_id) REFERENCES books (id) ON DELETE CASCADE
);

CREATE TRIGGER IF NOT EXISTS trg_pages_updated_at
    AFTER UPDATE
    ON pages
    FOR EACH ROW
    WHEN OLD.updated_at = NEW.updated_at OR NEW.updated_at IS NULL
BEGIN
    UPDATE pages SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

CREATE TRIGGER IF NOT EXISTS trg_panels_updated_at
    AFTER UPDATE
    ON panels
    FOR EACH ROW
    WHEN OLD.updated_at = NEW.updated_at OR NEW.updated_at IS NULL
BEGIN
    UPDATE panels SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

CREATE INDEX IF NOT EXISTS idx_chapters_book_id ON chapters (book_id);
CREATE INDEX IF NOT EXISTS idx_pages_chapter_id ON pages (chapter_id);
CREATE INDEX IF NOT EXISTS idx_panels_page_id ON panels (page_id);
CREATE INDEX IF NOT EXISTS idx_scripts_book_id ON scripts (book_id);
CREATE INDEX IF NOT EXISTS idx_splits_script_id ON splits (script_id);
CREATE INDEX IF NOT EXISTS idx_splits_panel_id ON splits (panel_id);
CREATE INDEX IF NOT EXISTS idx_vabs_script_id ON vabs (script_id);
CREATE INDEX IF NOT EXISTS idx_vabs_book_id ON vabs (book_id);