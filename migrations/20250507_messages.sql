CREATE TABLE  IF NOT EXISTS  messages (
    id INT AUTO_INCREMENT PRIMARY KEY,
    sender_id INT NOT NULL,       -- ID of the user sending the message
    receiver_id INT NOT NULL,     -- ID of the user receiving the message
    content TEXT NOT NULL,        -- The content of the message
    sent_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Timestamp when the message is sent
    FOREIGN KEY (sender_id) REFERENCES users(id),   -- Assuming a 'users' table exists
    FOREIGN KEY (receiver_id) REFERENCES users(id)  -- Assuming a 'users' table exists
);
