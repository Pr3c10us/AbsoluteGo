DROP INDEX IF EXISTS idx_events_operation_chapter;
DROP INDEX IF EXISTS idx_events_status_script;
DROP INDEX IF EXISTS idx_events_status_chapter;

DROP INDEX IF EXISTS idx_events_operation;
DROP INDEX IF EXISTS idx_events_status;

DROP INDEX IF EXISTS idx_events_vab_id;
DROP INDEX IF EXISTS idx_events_script_id;
DROP INDEX IF EXISTS idx_events_chapter_id;
DROP INDEX IF EXISTS idx_events_book_id;

DROP TRIGGER IF EXISTS trg_events_updated_at;

DROP TABLE IF EXISTS events;