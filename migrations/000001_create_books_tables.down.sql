PRAGMA foreign_keys = OFF;

DROP TRIGGER IF EXISTS trg_panels_updated_at;
DROP TRIGGER IF EXISTS trg_pages_updated_at;

DROP INDEX IF EXISTS idx_vabs_script_id;
DROP INDEX IF EXISTS idx_script_splits_panel_id;
DROP INDEX IF EXISTS idx_script_splits_page_id;
DROP INDEX IF EXISTS idx_script_splits_chapter_id;
DROP INDEX IF EXISTS idx_script_splits_script_id;
DROP INDEX IF EXISTS idx_scripts_book_id;
DROP INDEX IF EXISTS idx_panels_page_id;
DROP INDEX IF EXISTS idx_pages_chapter_id;
DROP INDEX IF EXISTS idx_chapters_book_id;

DROP TABLE IF EXISTS vabs;
DROP TABLE IF EXISTS splits;
DROP TABLE IF EXISTS scripts;
DROP TABLE IF EXISTS panels;
DROP TABLE IF EXISTS pages;
DROP TABLE IF EXISTS chapters;
DROP TABLE IF EXISTS books;

PRAGMA foreign_keys = ON;