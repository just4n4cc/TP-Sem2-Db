-- INIT ALL
CREATE EXTENSION IF NOT EXISTS citext;

-- TABLES ----------------------------------
CREATE UNLOGGED TABLE "User"
(
    id SERIAL PRIMARY KEY,
    nickname CITEXT COLLATE "C" UNIQUE NOT NULL,
    fullname TEXT NOT NULL,
    about TEXT,
    email CITEXT UNIQUE NOT NULL
);

CREATE UNLOGGED TABLE Forum
(
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    "user" CITEXT REFERENCES "User"(nickname) NOT NULL,
    slug CITEXT UNIQUE NOT NULL,
    posts INT DEFAULT 0,
    threads INT DEFAULT 0
);

CREATE UNLOGGED TABLE Thread
(
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    author CITEXT REFERENCES "User"(nickname) NOT NULL,
    forum CITEXT REFERENCES Forum(slug) NOT NULL,
    message TEXT NOT NULL,
    votes INT DEFAULT 0,
    slug CITEXT,
    created TIMESTAMP WITH TIME ZONE
);

CREATE UNLOGGED TABLE Post
(
    id SERIAL PRIMARY KEY,
    parent INT DEFAULT 0,
    author CITEXT REFERENCES "User"(nickname) NOT NULL,
    message TEXT NOT NULL,
    isEdited bool NOT NULL DEFAULT FALSE,
    forum CITEXT REFERENCES Forum(slug) NOT NULL,
    thread INT REFERENCES Thread(id) NOT NULL,
    created  TIMESTAMP WITH TIME ZONE,
    path INT[] NOT NULL
);

CREATE UNLOGGED TABLE Vote
(
    id SERIAL PRIMARY KEY,
    threadid INT REFERENCES Thread(id) NOT NULL,
    "user" CITEXT REFERENCES "User"(nickname) NOT NULL,
    vote INT NOT NULL,
    UNIQUE (threadid, "user")
);

CREATE UNLOGGED TABLE ForumUsers
(
    id SERIAL PRIMARY KEY,
    "user" CITEXT COLLATE "C" REFERENCES "User"(nickname) NOT NULL,
    forum CITEXT REFERENCES Forum(slug) NOT NULL,
    UNIQUE (forum, "user")
);


-- TRIGGERS ----------------------------------
-- POST INSERT
CREATE FUNCTION post_insert() RETURNS TRIGGER AS $post_insert$
DECLARE
    prevpath INT[];
BEGIN
    IF new.parent != 0 THEN
        SELECT path, thread FROM Post WHERE id = new.parent INTO prevpath;
        new.path := array_append(prevpath, new.id);
    ELSE
        new.path[1] := new.id;
    END IF;
    RETURN new;
END
$post_insert$ LANGUAGE plpgsql;
CREATE TRIGGER post_insert BEFORE INSERT ON Post FOR EACH ROW EXECUTE PROCEDURE post_insert();

-- POST INSERT FORUM
CREATE FUNCTION post_insert_forum() RETURNS TRIGGER AS $post_insert_forum$
BEGIN
    UPDATE Forum SET posts = posts + 1 WHERE slug = new.forum;
    INSERT INTO ForumUsers ("user", forum) VALUES (new.author, new.forum) ON CONFLICT DO NOTHING;
    RETURN new;
END
$post_insert_forum$ LANGUAGE plpgsql;
CREATE TRIGGER post_insert_forum AFTER INSERT ON Post FOR EACH ROW EXECUTE PROCEDURE post_insert_forum();

-- THREAD INSERT FORUM
CREATE FUNCTION thread_insert_forum() RETURNS TRIGGER AS $thread_insert_forum$
BEGIN
    UPDATE Forum SET threads = threads + 1 WHERE slug = new.forum;
    INSERT INTO ForumUsers ("user", forum) VALUES (new.author, new.forum) ON CONFLICT DO NOTHING;
    RETURN new;
END
$thread_insert_forum$ LANGUAGE plpgsql;
CREATE TRIGGER thread_insert_forum AFTER INSERT ON Thread FOR EACH ROW EXECUTE PROCEDURE thread_insert_forum();

-- FORUM INSTERT
CREATE FUNCTION forum_insert() RETURNS TRIGGER AS $forum_insert$
DECLARE
    nick text;
BEGIN
    SELECT nickname FROM "User" WHERE nickname = new."user" INTO nick;
    IF nick != '' THEN
        new."user" := nick;
    END IF;
    RETURN new;
END
$forum_insert$ LANGUAGE plpgsql;
CREATE TRIGGER forum_insert BEFORE INSERT ON Forum FOR EACH ROW EXECUTE PROCEDURE forum_insert();

-- THREAD INSTERT
CREATE FUNCTION thread_insert() RETURNS TRIGGER AS $thread_insert$
DECLARE
    nick text;
    forum text;
BEGIN
    SELECT nickname FROM "User" WHERE nickname = new."author" INTO nick;
    IF nick != '' THEN
        new."author" := nick;
    END IF;
    SELECT slug FROM Forum WHERE slug = new."forum" INTO forum;
    IF forum != '' THEN
        new."forum" := forum;
    END IF;
    RETURN new;
END
$thread_insert$ LANGUAGE plpgsql;
CREATE TRIGGER thread_insert BEFORE INSERT ON Thread FOR EACH ROW EXECUTE PROCEDURE thread_insert();

-- POST UPDATE
CREATE FUNCTION post_update() RETURNS TRIGGER AS $post_update$
BEGIN
    IF old.message != new.message THEN
        new.isedited := true;
    END IF;
    RETURN new;
END
$post_update$ LANGUAGE plpgsql;
CREATE TRIGGER post_update BEFORE UPDATE ON Post FOR EACH ROW EXECUTE PROCEDURE post_update();


-- INDEXES -----------------
CREATE INDEX IF NOT EXISTS idx_nickname_user ON "User" USING hash(nickname);
CREATE INDEX IF NOT EXISTS idx_email_user ON "User" USING hash(email);

CREATE INDEX IF NOT EXISTS idx_slug_forum ON Forum USING hash(slug);

CREATE UNIQUE INDEX IF NOT EXISTS idx_slug_thread ON Thread USING btree(slug) WHERE slug <> '';
CREATE INDEX IF NOT EXISTS idx_forum_thread ON Thread USING hash(forum);
CREATE INDEX IF NOT EXISTS idx_created_thread ON Thread USING btree(created);
CREATE INDEX IF NOT EXISTS idx_forum_created_thread ON Thread USING btree(forum, created);

CREATE INDEX IF NOT EXISTS idx_thread_post ON Post USING btree(thread);
CREATE INDEX IF NOT EXISTS idx_thread_id_post ON Post USING btree(thread, id);
CREATE INDEX IF NOT EXISTS idx_created_post ON Post USING btree(created);
CREATE INDEX IF NOT EXISTS idx_path1_post ON Post USING btree((path[1]));
CREATE INDEX IF NOT EXISTS idx_path_post ON Post USING btree(path);
CREATE INDEX IF NOT EXISTS idx_parent_post ON Post USING btree(parent);
CREATE INDEX IF NOT EXISTS idx_forum_post ON Post USING hash(forum);

CREATE INDEX IF NOT EXISTS idx_user_threadid_vote ON Vote USING btree("user", threadid);

CREATE INDEX IF NOT EXISTS idx_forum_forumusers ON ForumUsers USING hash(forum);
CREATE INDEX IF NOT EXISTS idx_forum_forumusers ON ForumUsers USING hash("user");
CREATE INDEX IF NOT EXISTS idx_forum_user_forumusers ON ForumUsers USING btree(forum, "user");

VACUUM;
VACUUM ANALYZE;
