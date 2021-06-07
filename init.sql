CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE users (
    nickname CITEXT PRIMARY KEY,
    fullname TEXT NOT NULL,
    email CITEXT UNIQUE,
    about TEXT
);

CREATE TABLE forum (
    slug    CITEXT PRIMARY KEY,
    title   TEXT NOT NULL,
    "user"  CITEXT,
    posts   BIGINT DEFAULT 0,
    threads BIGINT DEFAULT 0,

    FOREIGN KEY ("user") REFERENCES "users" (nickname)
);

CREATE TABLE threads (
    id      SERIAL PRIMARY KEY,
    title   TEXT NOT NULL,
    created TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    author  CITEXT REFERENCES "users" (Nickname),
    forum   CITEXT REFERENCES "forum" (slug),
    message TEXT NOT NULL,
    slug    CITEXT UNIQUE,
    votes   INT DEFAULT 0
);

CREATE TABLE posts (
    id       BIGSERIAL PRIMARY KEY,
    author   CITEXT NOT NULL,
    created  TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    forum    CITEXT,
    message  TEXT NOT NULL,
    parent   BIGINT DEFAULT 0,
    id_thread   INT,
    path     BIGINT[] DEFAULT ARRAY []::INTEGER[],
    isEdited BOOLEAN DEFAULT FALSE,

    FOREIGN KEY (author) REFERENCES "users"  (nickname),
    FOREIGN KEY (id_thread) REFERENCES "threads" (id)
);

CREATE TABLE votes
(
    id        BIGSERIAL PRIMARY KEY,
    author    CITEXT REFERENCES "users" (nickname),
    voice     INT NOT NULL,
    id_thread INT,

    FOREIGN KEY (id_thread) REFERENCES "threads" (id),
    UNIQUE (author, id_thread)
);

CREATE UNLOGGED TABLE forum_users
(
    nickname CITEXT NOT NULL,
    fullname TEXT NOT NULL,
    about    TEXT,
    email    CITEXT,
    slug     CITEXT NOT NULL,

    FOREIGN KEY (nickname) REFERENCES "users" (nickname),
    FOREIGN KEY (slug) REFERENCES "forum" (slug),
    UNIQUE (nickname, slug)
);


CREATE INDEX all_forum_users ON forum_users (nickname, fullname, about, email);
CLUSTER forum_users USING all_forum_users;
CREATE INDEX nickname_forum_users ON forum_users using hash (nickname);
CREATE INDEX forums_users_info ON forum_users (fullname, about, email);

CREATE INDEX if not exists user_nickname ON users using hash (nickname);
CLUSTER users USING user_nickname;
CREATE INDEX if not exists user_email ON users using hash (email);
CREATE INDEX if not exists forum_slug ON forum using hash (slug);

CREATE UNIQUE INDEX if not exists forum_users_unique ON forum_users (slug, nickname);
cluster forum_users using forum_users_unique;

CREATE INDEX if not exists thr_slug ON threads using hash (slug);
CREATE INDEX if not exists thr_date ON threads (created);
-- CREATE INDEX if not exists thr_votes ON threads (id, votes);
CREATE INDEX if not exists thr_forum ON threads using hash (forum);
CREATE INDEX if not exists thr_forum_date ON threads (forum, created);

CREATE INDEX if not exists post_id_path ON posts (id, (path[1]));
CREATE INDEX if not exists post_thread_id_path1_parent ON posts (id_thread, id, (path[1]), parent);
CREATE INDEX if not exists post_thread_path_id ON posts (id_thread, path, id);
CREATE INDEX IF NOT EXISTS post_path1 on posts ((path[1]));
CREATE INDEX IF NOT EXISTS post_thread_id on posts (id_thread, id);
CREATE INDEX IF NOT EXISTS post_thr_id ON posts (id_thread);

CREATE INDEX IF NOT EXISTS post_path1_path_id_desc ON posts ((path[1]) DESC, path, id);
CREATE INDEX IF NOT EXISTS post_path1_path_id_asc ON posts ((path[1]) ASC, path, id);

CREATE UNIQUE INDEX IF NOT EXISTS vote_unique on votes (author, id_thread);


CREATE OR REPLACE FUNCTION updateForumUser() RETURNS TRIGGER AS
$update_forum_user$
DECLARE
    new_fullname CITEXT;
    new_about    CITEXT;
    new_email    CITEXT;
BEGIN
    SELECT fullname, about, email FROM users WHERE nickname = NEW.author INTO new_fullname, new_about, new_email;
    INSERT INTO forum_users (nickname, fullname, about, email, slug)
    VALUES (NEW.author, new_fullname, new_about, new_email, NEW.forum) on conflict do nothing;
    RETURN NEW;
end
$update_forum_user$ LANGUAGE plpgsql;
CREATE TRIGGER thread_update_forum_user
    AFTER INSERT
    ON threads
    FOR EACH ROW
EXECUTE PROCEDURE updateForumUser();
CREATE TRIGGER post_update_forum_user
    AFTER INSERT
    ON posts
    FOR EACH ROW
EXECUTE PROCEDURE updateForumUser();


CREATE OR REPLACE FUNCTION insertVote() RETURNS TRIGGER AS
$insert_vote$
BEGIN
    UPDATE threads SET votes=(votes+NEW.voice) WHERE id=NEW.id_thread;
    RETURN NEW;
end
$insert_vote$ LANGUAGE plpgsql;
CREATE TRIGGER insert_vote
    BEFORE INSERT ON votes
    FOR EACH ROW
EXECUTE PROCEDURE insertVote();


CREATE OR REPLACE FUNCTION updateVotes() RETURNS TRIGGER AS
$update_votes$
BEGIN
    IF OLD.voice <> NEW.voice THEN
        UPDATE threads SET votes=(votes-OLD.voice+NEW.voice) WHERE id=NEW.id_thread;
    END IF;
    return NEW;
end
$update_votes$ LANGUAGE plpgsql;
CREATE TRIGGER update_votes
    BEFORE UPDATE ON votes
    FOR EACH ROW
EXECUTE PROCEDURE updateVotes();


CREATE OR REPLACE FUNCTION insertThread() RETURNS trigger AS
$insert_thread$
BEGIN
    UPDATE forum
    SET threads = threads + 1
    WHERE slug = NEW.forum;
    RETURN NULL;
END;
$insert_thread$ LANGUAGE plpgsql;
CREATE TRIGGER insert_thread
    AFTER INSERT ON threads
    FOR EACH ROW
EXECUTE PROCEDURE insertThread();

CREATE OR REPLACE FUNCTION updatePath() RETURNS TRIGGER AS
$update_path$
DECLARE
    parentPath BIGINT[];
    first_parent_thread INT;
BEGIN
    IF (NEW.parent IS NULL) THEN
        NEW.path := array_append(new.path, new.id);
    ELSE
        SELECT path FROM posts WHERE id = new.parent INTO parentPath;
        SELECT id_thread FROM posts WHERE id = parentPath[1] INTO first_parent_thread;
        IF NOT FOUND OR first_parent_thread != NEW.id_thread THEN
            RAISE EXCEPTION 'parent is from different thread' USING ERRCODE = '00409';
        end if;
        NEW.path := NEW.path || parentPath || new.id;
    end if;
    UPDATE forum SET Posts=Posts + 1 WHERE forum.slug = new.forum;
    RETURN new;
end
$update_path$ LANGUAGE plpgsql;
CREATE TRIGGER update_path
    BEFORE INSERT
    ON posts
    FOR EACH ROW
EXECUTE PROCEDURE updatePath();