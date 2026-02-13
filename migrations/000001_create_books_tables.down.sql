PRAGMA foreign_keys = OFF;

DROP TRIGGER IF EXISTS trg_panels_updated_at;
DROP TRIGGER IF EXISTS trg_pages_updated_at;

DROP INDEX IF EXISTS idx_videos_audio_id;
DROP INDEX IF EXISTS idx_videos_video_id;
DROP INDEX IF EXISTS idx_slide_shows_video_id;
DROP INDEX IF EXISTS idx_audios_video_id;
DROP INDEX IF EXISTS idx_vab_script_id;
DROP INDEX IF EXISTS idx_script_splits_panel_id;
DROP INDEX IF EXISTS idx_script_splits_page_id;
DROP INDEX IF EXISTS idx_script_splits_chapter_id;
DROP INDEX IF EXISTS idx_script_splits_script_id;
DROP INDEX IF EXISTS idx_scripts_book_id;
DROP INDEX IF EXISTS idx_panels_page_id;
DROP INDEX IF EXISTS idx_pages_chapter_id;
DROP INDEX IF EXISTS idx_chapters_book_id;

DROP TABLE IF EXISTS videos;
DROP TABLE IF EXISTS slide_shows;
DROP TABLE IF EXISTS audios;
DROP TABLE IF EXISTS vab;
DROP TABLE IF EXISTS script_splits;
DROP TABLE IF EXISTS scripts;
DROP TABLE IF EXISTS panels;
DROP TABLE IF EXISTS pages;
DROP TABLE IF EXISTS chapters;
DROP TABLE IF EXISTS books;

PRAGMA foreign_keys = ON;