CREATE TABLE
    IF NOT EXISTS posts (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_id INTEGER NOT NULL,
        title TEXT NOT NULL,
        content TEXT NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        total_likes INTEGER DEFAULT 0 CHECK (total_likes >= 0),
        total_dislikes INTEGER DEFAULT 0 CHECK (total_dislikes >= 0),
        total_comments INTEGER DEFAULT 0 CHECK (total_comments >= 0),
        FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
    );

-- -- Increment on insert
-- CREATE TRIGGER IF NOT EXISTS increment_like_count AFTER INSERT ON post_likes BEGIN
-- UPDATE posts
-- SET total_likes = total_likes + 1
-- WHERE id = NEW.post_id;
-- END;
-- -- Decrement on delete
-- CREATE TRIGGER IF NOT EXISTS decrement_like_count AFTER DELETE ON post_likes BEGIN
-- UPDATE posts
-- SET total_likes = total_likes - 1
-- WHERE id = OLD.post_id;
-- END;
