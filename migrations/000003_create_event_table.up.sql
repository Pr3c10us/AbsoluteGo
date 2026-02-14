CREATE TABLE IF NOT EXISTS events
(
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    status      TEXT    NOT NULL,
    operation   TEXT    NOT NULL,
    description TEXT    NOT NULL,
    book_id     INTEGER NOT NULL,
    chapter_id  INTEGER,
    script_id   INTEGER,
    vab_id      INTEGER,
    updated_at  DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (book_id) REFERENCES books (id) ON DELETE CASCADE,
    FOREIGN KEY (chapter_id) REFERENCES chapters (id) ON DELETE CASCADE,
    FOREIGN KEY (script_id) REFERENCES scripts (id) ON DELETE CASCADE,
    FOREIGN KEY (vab_id) REFERENCES vabs (id) ON DELETE CASCADE
);

CREATE TRIGGER IF NOT EXISTS trg_events_updated_at
    AFTER UPDATE
    ON events
    FOR EACH ROW
    WHEN OLD.updated_at = NEW.updated_at OR NEW.updated_at IS NULL
BEGIN
    UPDATE events SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;


CREATE INDEX IF NOT EXISTS idx_events_book_id ON events (book_id);
CREATE INDEX IF NOT EXISTS idx_events_chapter_id ON events (chapter_id);
CREATE INDEX IF NOT EXISTS idx_events_script_id ON events (script_id);
CREATE INDEX IF NOT EXISTS idx_events_vab_id ON events (vab_id);

CREATE INDEX IF NOT EXISTS idx_events_status ON events (status);
CREATE INDEX IF NOT EXISTS idx_events_operation ON events (operation);

CREATE INDEX IF NOT EXISTS idx_events_status_chapter ON events (status, chapter_id);
CREATE INDEX IF NOT EXISTS idx_events_status_script ON events (status, script_id);
CREATE INDEX IF NOT EXISTS idx_events_operation_chapter ON events (operation, chapter_id);