DROP SCHEMA IF EXISTS tpdb CASCADE;

-- DROP ALL
DROP TABLE IF EXISTS tpdb."User" CASCADE;
DROP TABLE IF EXISTS tpdb."Forum" CASCADE ;
DROP TABLE IF EXISTS tpdb."Thread" CASCADE;
DROP TABLE IF EXISTS tpdb."Post" CASCADE ;
DROP TABLE IF EXISTS tpdb."Vote" CASCADE ;
DROP TABLE IF EXISTS tpdb."ForumUsers" CASCADE ;

DROP TRIGGER IF EXISTS user_inc ON tpdb."User";
DROP TRIGGER IF EXISTS thread_inc ON tpdb."Thread";
DROP TRIGGER IF EXISTS post_inc ON tpdb."Post";
DROP TRIGGER IF EXISTS forum_inc ON tpdb."Forum";
DROP TRIGGER IF EXISTS post_insert ON tpdb."Post";
DROP TRIGGER IF EXISTS forum_insert ON tpdb."Forum";
DROP TRIGGER IF EXISTS thread_insert ON tpdb."Thread";
DROP TRIGGER IF EXISTS post_insert_forum ON tpdb."Post";
DROP TRIGGER IF EXISTS thread_insert_forum ON tpdb."Thread";
DROP TRIGGER IF EXISTS post_update ON tpdb."Post";

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

DROP INDEX IF EXISTS tpdb.idx_nickname_user;
DROP INDEX IF EXISTS tpdb.idx_email_user;

DROP INDEX IF EXISTS tpdb.idx_slug_forum;

DROP INDEX IF EXISTS tpdb.idx_slug_thread;
DROP INDEX IF EXISTS tpdb.idx_forum_thread;
DROP INDEX IF EXISTS tpdb.idx_forum_thread;
DROP INDEX IF EXISTS tpdb.idx_created_thread;

DROP INDEX IF EXISTS tpdb.idx_thread_post;
DROP INDEX IF EXISTS tpdb.idx_thread_id_post;
DROP INDEX IF EXISTS tpdb.idx_created_post;
DROP INDEX IF EXISTS tpdb.idx_path_post;
DROP INDEX IF EXISTS tpdb.idx_path1_post;
DROP INDEX IF EXISTS tpdb.idx_id_path1_post;
DROP INDEX IF EXISTS tpdb.idx_forum_post;

DROP INDEX IF EXISTS tpdb.idx_user_threadid_vote;

DROP INDEX IF EXISTS tpdb.idx_forum_forumusers;

