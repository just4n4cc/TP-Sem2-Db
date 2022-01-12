DROP SCHEMA IF EXISTS tpdb CASCADE;
CREATE EXTENSION IF NOT EXISTS citext;
CREATE SCHEMA tpdb;

-- DROP ALL
DROP TABLE IF EXISTS tpdb."User" CASCADE;
DROP TABLE IF EXISTS tpdb."Forum" CASCADE ;
DROP TABLE IF EXISTS tpdb."Thread" CASCADE;
DROP TABLE IF EXISTS tpdb."Post" CASCADE ;
DROP TABLE IF EXISTS "Service" CASCADE ;

DROP TRIGGER IF EXISTS user_inc ON tpdb."User";
DROP TRIGGER IF EXISTS thread_inc ON tpdb."Thread";
DROP TRIGGER IF EXISTS post_inc ON tpdb."Post";
DROP TRIGGER IF EXISTS forum_inc ON tpdb."Forum";

DROP FUNCTION IF EXISTS service_user_inc() CASCADE;
DROP FUNCTION IF EXISTS service_thread_inc() CASCADE;
DROP FUNCTION IF EXISTS service_forum_inc() CASCADE;
DROP FUNCTION IF EXISTS service_post_inc() CASCADE;
-- --------------------------------

CREATE UNLOGGED TABLE "Service"
(
    id SERIAL PRIMARY KEY,
    forums INT DEFAULT 0,
    posts INT DEFAULT 0,
    threads INT DEFAULT 0,
    users INT DEFAULT 0
);
INSERT INTO "Service" DEFAULT VALUES;

CREATE UNLOGGED TABLE tpdb."User"
(
    id SERIAL PRIMARY KEY,
    nickname CITEXT UNIQUE NOT NULL,
    fullname TEXT NOT NULL,
    about TEXT,
    email CITEXT UNIQUE NOT NULL
);

CREATE UNLOGGED TABLE tpdb."Forum"
(
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    "user" CITEXT REFERENCES tpdb."User"(nickname) NOT NULL,
    slug CITEXT UNIQUE NOT NULL,
    posts INT DEFAULT 0,
    threads INT DEFAULT 0
);

CREATE UNLOGGED TABLE tpdb."Thread"
(
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    author CITEXT REFERENCES tpdb."User"(nickname) NOT NULL,
    forum CITEXT REFERENCES tpdb."Forum"(slug) NOT NULL,
    message TEXT NOT NULL,
    votes INT DEFAULT 0,
    slug CITEXT UNIQUE NOT NULL,
    created TIMESTAMP WITH TIME ZONE
);

CREATE UNLOGGED TABLE tpdb."Post"
(
    id SERIAL PRIMARY KEY,
    parent INT DEFAULT 0,
    author CITEXT REFERENCES tpdb."User"(nickname) NOT NULL,
    message TEXT NOT NULL,
    isEdited bool NOT NULL DEFAULT FALSE,
    forum CITEXT REFERENCES tpdb."Forum"(slug) NOT NULL,
    thread INT REFERENCES tpdb."Thread"(id) NOT NULL,
    created  TIMESTAMP WITH TIME ZONE
);


-- SERVICE USERS
CREATE FUNCTION service_user_inc() RETURNS TRIGGER AS $user_inc$
BEGIN
    UPDATE "Service" SET users = users + 1 WHERE id = 1;
    RETURN NULL;
END;
$user_inc$ LANGUAGE plpgsql;
CREATE TRIGGER user_inc AFTER INSERT ON tpdb."User" EXECUTE PROCEDURE service_user_inc();

-- SERVICE THREADS
CREATE FUNCTION service_thread_inc() RETURNS TRIGGER AS $thread_inc$
BEGIN
    UPDATE "Service" SET threads = threads + 1 WHERE id = 1;
    RETURN NULL;
END;
$thread_inc$ LANGUAGE plpgsql;
CREATE TRIGGER thread_inc AFTER INSERT ON tpdb."Thread" EXECUTE PROCEDURE service_thread_inc();

-- SERVICE POSTS
CREATE FUNCTION service_post_inc() RETURNS TRIGGER AS $post_inc$
DECLARE
    record_count integer;
BEGIN
    SELECT COUNT(*) FROM newtbl INTO record_count;
    UPDATE "Service" SET posts = posts + record_count WHERE id = 1;
    RETURN NULL;
END;
$post_inc$ LANGUAGE plpgsql;
CREATE TRIGGER post_inc AFTER INSERT ON tpdb."Post"
    REFERENCING NEW TABLE as newtbl
    FOR EACH STATEMENT EXECUTE PROCEDURE service_post_inc();

-- SERVICE FORUMS
CREATE FUNCTION service_forum_inc() RETURNS TRIGGER AS $forum_inc$
BEGIN
    UPDATE "Service" SET forums = forums + 1 WHERE id = 1;
    RETURN NULL;
END;
$forum_inc$ LANGUAGE plpgsql;
CREATE TRIGGER forum_inc AFTER INSERT ON tpdb."Forum" EXECUTE PROCEDURE service_forum_inc();


-- SELECT * FROM tpdb."User";
-- SELECT * FROM tpdb."Forum";
-- SELECT * FROM tpdb."Thread";
-- SELECT * FROM tpdb."Post";
-- SELECT * FROM "Service";

-- INSERT INTO tpdb."User" (nickname, fullname, about, email) VALUES ('lala', 'lala', 'lala', 'lala@mail.ru');
-- INSERT INTO tpdb."User" (nickname, fullname, about, email) VALUES ('sasha', 'lala', 'lala', 'sasha@mail.ru');
-- INSERT INTO tpdb."User" (nickname, fullname, about, email) VALUES ('ksyasha', 'lala', 'lala', 'lalala@mail.ru');
--
-- INSERT INTO tpdb."Forum" (title, "user", slug, posts, threads) VALUES ('forum', 'lala', 'someforum', 0, 0);
-- INSERT INTO tpdb."Forum" (title, "user", slug, posts, threads) VALUES ('lala', 'ksyasha', 'lalaforum', 0, 0);
--
-- INSERT INTO tpdb."Thread" (title, author, forum, message, votes, slug, created)
--     VALUES ('thread1', 'lala', 'someforum', 'idk why but', 0, 'thread1', '2022-01-10T15:31:49.659Z');
-- INSERT INTO tpdb."Thread" (title, author, forum, message, votes, slug, created)
--     VALUES ('thread2', 'lala', 'someforum', 'idk why but again', 0, 'thread2', '2021-01-10T15:31:49.659Z');
-- INSERT INTO tpdb."Thread" (title, author, forum, message, votes, slug, created)
--     VALUES ('thread3', 'ksyasha', 'lalaforum', 'u dont know?', 0, 'thread3', '2021-12-10T15:31:49.659Z');
-- INSERT INTO tpdb."Thread" (title, author, forum, message, votes, slug, created)
--     VALUES ('thread4', 'ksyasha', 'lalaforum', 'idk either', 0, 'thread4', '2020-01-10T15:31:49.659Z');
--
-- INSERT INTO tpdb."Post" (parent, author, message, isedited, forum, thread, created)
--     VALUES (0, 'lala', 'helpppp', false, 'someforum', 1, '2023-01-10T15:31:49.659Z');
-- INSERT INTO tpdb."Post" (parent, author, message, isedited, forum, thread, created)
--     VALUES (0, 'sasha', 'helpppp her faster', false, 'someforum', 1, '2023-02-10T15:31:49.659Z');
-- INSERT INTO tpdb."Post" (parent, author, message, isedited, forum, thread, created)
--     VALUES (0, 'ksyasha', 'haha haaa', false, 'someforum', 1, '2023-03-10T15:31:49.659Z');
-- INSERT INTO tpdb."Post" (parent, author, message, isedited, forum, thread, created)
--     VALUES (0, 'lala', 'not funny at all...', false, 'someforum', 1, '2023-04-10T15:31:49.659Z');
--
-- INSERT INTO tpdb."Post" (parent, author, message, isedited, forum, thread, created)
--     VALUES (0, 'lala', 'helpppp', false, 'lalaforum', 3, '2023-01-10T15:31:49.659Z');
-- INSERT INTO tpdb."Post" (parent, author, message, isedited, forum, thread, created)
--     VALUES (0, 'lala', 'helpppp!!', false, 'lalaforum', 3, '2023-02-10T15:31:49.659Z');
-- INSERT INTO tpdb."Post" (parent, author, message, isedited, forum, thread, created)
--     VALUES (0, 'lala', 'pls', false, 'lalaforum', 3, '2023-03-10T15:31:49.659Z');
-- INSERT INTO tpdb."Post" (parent, author, message, isedited, forum, thread, created)
--     VALUES (0, 'lala', '((', false, 'lalaforum', 3, '2023-04-10T15:31:49.659Z');

-- insert into tpdb."Post"
-- (parent, author, message, forum, thread)
-- values (0, 'lala', 'hahahahah', 'someforum', 1), (0, 'lala', 'ahhahaha', 'someforum', 1) returning *;

-- update tpdb."Post" set created = '2023-04-10T15:31:49.659Z' where created is null;