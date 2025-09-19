-- CREATE TABLE users (
--     user_id SERIAL PRIMARY KEY
-- );

ALTER TABLE urls
ADD COLUMN user_id INT;

UPDATE urls SET user_id = 1 WHERE user_id IS NULL;

ALTER TABLE urls
ALTER COLUMN user_id SET NOT NULL;

-- ALTER TABLE urls
-- ADD CONSTRAINT fk_urls_user
-- FOREIGN KEY (user_id) REFERENCES users(user_id);
