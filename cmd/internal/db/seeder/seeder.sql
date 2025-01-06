-- RAW SQL to seed the database with some initial data

-- ===========================
-- 1. USERS
-- account numbers: 
-- 1234567891234567
-- 3998442793406687
-- 7180218105191773
-- 6129856725721562
-- ===========================

INSERT INTO users (account_number) 
VALUES
  ('$2a$12$Sghxkjm.nrzwR0ym.U82qOT6P.EfW16.L.h6pQBBAreGhjaEUVzWS'),
  ('$2a$12$S7uiGoEMGKJFLbZHWFJjUO.lao3MeI7EfECrg41/KY/n3VhXFDY16'),
  ('$2a$12$agN83W3sLUZDqGjTtYAHAeR.4iRrcNgcy9LpMXXzNNze8Qi3FGpCa'),
  ('$2a$12$IcQn6xuJp2ezw00NeNQ/c.iuZMZrKSBJy.mIdEDdbX6R7JfjhY9ZK');

-- ===========================
-- 2. POSTS
-- ===========================
-- We'll omit 'created_at' so it uses the default CURRENT_TIMESTAMP
INSERT INTO posts (content, user_id)
VALUES
  ('Hello from user 1!', 1),
  ('User 2 shares this post', 2),
  ('User 3 says hi', 3),
  ('User 4 is here too', 4),
  ('User 1 posts again!', 1);

-- ===========================
-- 3. COMMENTS
-- ===========================
INSERT INTO comments (content, post_id, user_id)
VALUES
  ('User 2 commenting on post #1', 1, 2),
  ('User 3 commenting on post #1', 1, 3),
  ('User 1 commenting on post #2', 2, 1),
  ('User 4 responding to post #3', 3, 4),
  ('User 3 responding to post #4', 4, 3);

-- ===========================
-- 4. POSTS_LIKES
-- ===========================
INSERT INTO posts_likes (post_id, user_id)
VALUES
  (1, 1),  -- user 1 likes post #1
  (1, 2),  -- user 2 likes post #1
  (1, 3),  -- user 3 likes post #1
  (2, 1),  -- user 1 likes post #2
  (2, 3),  -- user 3 likes post #2
  (2, 4),  -- user 4 likes post #2
  (3, 4);  -- user 4 likes post #3
