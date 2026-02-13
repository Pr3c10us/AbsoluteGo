CREATE TABLE IF NOT EXISTS events
(
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    status    TEXT,
    operation Text,
    chapter_id   INTEGER,
    script_id   INTEGER,
    vab_id   INTEGER,
    FOREIGN KEY (chapter_id) REFERENCES panels (id) ON DELETE CASCADE,
    FOREIGN KEY (script_id) REFERENCES scripts (id) ON DELETE CASCADE,
    FOREIGN KEY (vab_id) REFERENCES vabs (id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_events_chapter_id ON events (chapter_id);
CREATE INDEX IF NOT EXISTS idx_events_script_id ON events (script_id);
CREATE INDEX IF NOT EXISTS idx_events_vab_id ON events (vab_id);

CREATE INDEX IF NOT EXISTS idx_events_status ON events (status);
CREATE INDEX IF NOT EXISTS idx_events_operation ON events (operation);

CREATE INDEX IF NOT EXISTS idx_events_status_chapter ON events (status, chapter_id);
CREATE INDEX IF NOT EXISTS idx_events_status_script ON events (status, script_id);
CREATE INDEX IF NOT EXISTS idx_events_operation_chapter ON events (operation, chapter_id);