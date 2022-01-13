DROP SCHEMA IF EXISTS tpdb CASCADE;

-- DROP ALL
DROP TABLE IF EXISTS tpdb."User" CASCADE;
DROP TABLE IF EXISTS tpdb."Forum" CASCADE ;
DROP TABLE IF EXISTS tpdb."Thread" CASCADE;
DROP TABLE IF EXISTS tpdb."Post" CASCADE ;
DROP TABLE IF EXISTS tpdb."Vote" CASCADE ;
DROP TABLE IF EXISTS "Service" CASCADE ;

DROP TRIGGER IF EXISTS user_inc ON tpdb."User";
DROP TRIGGER IF EXISTS thread_inc ON tpdb."Thread";
DROP TRIGGER IF EXISTS post_inc ON tpdb."Post";
DROP TRIGGER IF EXISTS forum_inc ON tpdb."Forum";
DROP TRIGGER IF EXISTS post_insert ON tpdb."Post";
DROP TRIGGER IF EXISTS forum_insert ON tpdb."Forum";
DROP TRIGGER IF EXISTS thread_insert ON tpdb."Thread";

DROP FUNCTION IF EXISTS service_user_inc() CASCADE;
DROP FUNCTION IF EXISTS service_thread_inc() CASCADE;
DROP FUNCTION IF EXISTS service_forum_inc() CASCADE;
DROP FUNCTION IF EXISTS service_post_inc() CASCADE;
DROP FUNCTION IF EXISTS post_insert() CASCADE;
DROP FUNCTION IF EXISTS forum_insert() CASCADE;
DROP FUNCTION IF EXISTS thread_insert() CASCADE;
-- --------------------------------

