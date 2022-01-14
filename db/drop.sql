-- DROP SCHEMA IF EXISTS tpdb CASCADE;

-- DROP ALL
DROP TABLE IF EXISTS "User" CASCADE;
DROP TABLE IF EXISTS Forum CASCADE ;
DROP TABLE IF EXISTS Thread CASCADE;
DROP TABLE IF EXISTS Post CASCADE ;
DROP TABLE IF EXISTS Vote CASCADE ;
DROP TABLE IF EXISTS ForumUsers CASCADE ;

DROP TRIGGER IF EXISTS user_inc ON "User";
DROP TRIGGER IF EXISTS thread_inc ON Thread;
DROP TRIGGER IF EXISTS post_inc ON Post;
DROP TRIGGER IF EXISTS forum_inc ON Forum;
DROP TRIGGER IF EXISTS post_insert ON Post;
DROP TRIGGER IF EXISTS forum_insert ON Forum;
DROP TRIGGER IF EXISTS thread_insert ON Thread;
DROP TRIGGER IF EXISTS post_insert_forum ON Post;
DROP TRIGGER IF EXISTS thread_insert_forum ON Thread;
DROP TRIGGER IF EXISTS post_update ON Post;

DROP FUNCTION IF EXISTS service_user_inc() CASCADE;
DROP FUNCTION IF EXISTS service_thread_inc() CASCADE;
DROP FUNCTION IF EXISTS service_forum_inc() CASCADE;
DROP FUNCTION IF EXISTS service_post_inc() CASCADE;
DROP FUNCTION IF EXISTS post_insert() CASCADE;
DROP FUNCTION IF EXISTS forum_insert() CASCADE;
DROP FUNCTION IF EXISTS thread_insert() CASCADE;
DROP FUNCTION IF EXISTS post_insert_forum() CASCADE;
DROP FUNCTION IF EXISTS thread_insert_forum() CASCADE;
DROP FUNCTION IF EXISTS post_update() CASCADE;
-- --------------------------------

DROP INDEX IF EXISTS idx_nickname_user;
DROP INDEX IF EXISTS idx_email_user;

DROP INDEX IF EXISTS idx_slug_forum;

DROP INDEX IF EXISTS idx_slug_thread;
DROP INDEX IF EXISTS idx_forum_thread;
DROP INDEX IF EXISTS idx_forum_thread;
DROP INDEX IF EXISTS idx_created_thread;

DROP INDEX IF EXISTS idx_thread_post;
DROP INDEX IF EXISTS idx_thread_id_post;
DROP INDEX IF EXISTS idx_created_post;
DROP INDEX IF EXISTS idx_path_post;
DROP INDEX IF EXISTS idx_path1_post;
DROP INDEX IF EXISTS idx_id_path1_post;
DROP INDEX IF EXISTS idx_forum_post;

DROP INDEX IF EXISTS idx_user_threadid_vote;

DROP INDEX IF EXISTS idx_forum_forumusers;

