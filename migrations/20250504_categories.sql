CREATE TABLE
    IF NOT EXISTS categories (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        category_name TEXT NOT NULL UNIQUE
    );

INSERT
OR IGNORE INTO categories (category_name)
VALUES
    ('Technology'),
    ('Science'),
    ('Health'),
    ('Lifestyle'),
    ('Education'),
    ('Gaming'),
    ('Business');