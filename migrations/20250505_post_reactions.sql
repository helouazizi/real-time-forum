CREATE TABLE IF NOT EXISTS post_reactions (
    user_id INTEGER,
    post_id INTEGER,
    reaction TEXT CHECK(reaction IN ('like', 'dislike')),
    PRIMARY KEY (user_id, post_id)
);