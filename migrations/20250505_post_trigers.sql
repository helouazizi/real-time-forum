CREATE TRIGGER  IF NOT EXISTS  trg_insert_reaction
AFTER INSERT ON post_reactions
BEGIN
    UPDATE posts
    SET 
        total_likes = total_likes + CASE WHEN NEW.reaction = 'like' THEN 1 ELSE 0 END,
        total_dislikes = total_dislikes + CASE WHEN NEW.reaction = 'dislike' THEN 1 ELSE 0 END
    WHERE id = NEW.post_id;
END;


CREATE TRIGGER  IF NOT EXISTS  trg_update_reaction
AFTER UPDATE ON post_reactions
BEGIN
    UPDATE posts
    SET 
        total_likes = total_likes + CASE WHEN NEW.reaction = 'like' THEN 1 ELSE -1 END,
        total_dislikes = total_dislikes + CASE WHEN NEW.reaction = 'dislike' THEN 1 ELSE -1 END
    WHERE id = NEW.post_id;
END;


CREATE TRIGGER  IF NOT EXISTS  trg_delete_reaction
AFTER DELETE ON post_reactions
BEGIN
    UPDATE posts
    SET 
        total_likes = total_likes - CASE WHEN OLD.reaction = 'like' THEN 1 ELSE 0 END,
        total_dislikes = total_dislikes - CASE WHEN OLD.reaction = 'dislike' THEN 1 ELSE 0 END
    WHERE id = OLD.post_id;
END;
