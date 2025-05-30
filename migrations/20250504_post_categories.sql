CREATE TABLE
    IF NOT EXISTS post_categories (
        post_id INTEGER NOT NULL,
        category_id INTEGER NOT NULL,
        PRIMARY KEY (post_id, category_id),
        FOREIGN KEY (post_id) REFERENCES posts (id) ON DELETE CASCADE,
        FOREIGN KEY (category_id) REFERENCES categories (id) ON DELETE CASCADE
    );