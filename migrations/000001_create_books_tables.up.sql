-- Migration: Initial Schema
-- Created: 2026-02-09

PRAGMA foreign_keys = ON;

-- ============================================
-- Core Book Structure
-- ============================================

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

-- ============================================
-- Scripts
-- ============================================

CREATE TABLE IF NOT EXISTS scripts
(
    id       INTEGER PRIMARY KEY AUTOINCREMENT,
    name     TEXT    NOT NULL,
    content  TEXT,
    book_id  INTEGER NOT NULL,
    chapters TEXT DEFAULT '[]', -- JSON array of integers
    FOREIGN KEY (book_id) REFERENCES books (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS script_splits
(
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    script_id  INTEGER NOT NULL,
    content    TEXT,
    panel_id   INTEGER,
    effect     TEXT,
    FOREIGN KEY (script_id) REFERENCES scripts (id) ON DELETE CASCADE,
    FOREIGN KEY (panel_id) REFERENCES panels (id) ON DELETE CASCADE
);

-- ============================================
-- Video / Audio / Slideshow
-- ============================================

CREATE TABLE IF NOT EXISTS vab
(
    id        INTEGER PRIMARY KEY AUTOINCREMENT,
    name      TEXT    NOT NULL,
    script_id INTEGER NOT NULL,
    url       TEXT,
    music     TEXT DEFAULT '[]', -- JSON array of strings
    FOREIGN KEY (script_id) REFERENCES scripts (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS audios
(
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    video_id    INTEGER NOT NULL,
    page_id     INTEGER,
    voice       TEXT,
    voice_style TEXT,
    url         TEXT,
    FOREIGN KEY (video_id) REFERENCES vab (id) ON DELETE CASCADE,
    FOREIGN KEY (page_id) REFERENCES pages (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS slide_shows
(
    id               INTEGER PRIMARY KEY AUTOINCREMENT,
    video_id         INTEGER NOT NULL,
    script_splits_id INTEGER NOT NULL,
    audio_duration   REAL,
    FOREIGN KEY (video_id) REFERENCES vab (id) ON DELETE CASCADE,
    FOREIGN KEY (script_splits_id) REFERENCES script_splits (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS videos
(
    id       INTEGER PRIMARY KEY AUTOINCREMENT,
    video_id INTEGER NOT NULL,
    page_id  INTEGER,
    url      TEXT,
    audio_id INTEGER,
    FOREIGN KEY (video_id) REFERENCES vab (id) ON DELETE CASCADE,
    FOREIGN KEY (page_id) REFERENCES pages (id) ON DELETE CASCADE,
    FOREIGN KEY (audio_id) REFERENCES audios (id) ON DELETE CASCADE
);


-- ============================================
-- Triggers: Auto-update updated_at
-- ============================================

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

-- ============================================
-- Indexes
-- ============================================

CREATE INDEX IF NOT EXISTS idx_chapters_book_id ON chapters (book_id);
CREATE INDEX IF NOT EXISTS idx_pages_chapter_id ON pages (chapter_id);
CREATE INDEX IF NOT EXISTS idx_panels_page_id ON panels (page_id);
CREATE INDEX IF NOT EXISTS idx_scripts_book_id ON scripts (book_id);
CREATE INDEX IF NOT EXISTS idx_script_splits_script_id ON script_splits (script_id);
CREATE INDEX IF NOT EXISTS idx_script_splits_panel_id ON script_splits (panel_id);
CREATE INDEX IF NOT EXISTS idx_vab_script_id ON vab (script_id);
CREATE INDEX IF NOT EXISTS idx_audios_video_id ON audios (video_id);
CREATE INDEX IF NOT EXISTS idx_slide_shows_video_id ON slide_shows (video_id);
CREATE INDEX IF NOT EXISTS idx_videos_video_id ON videos (video_id);
CREATE INDEX IF NOT EXISTS idx_videos_audio_id ON videos (audio_id);