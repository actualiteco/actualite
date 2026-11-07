ALTER TABLE posts DROP CONSTRAINT IF EXISTS posts_image_only_for_root;

DELETE FROM post_hashtags
WHERE (post_id, hashtag_id) IN (
    ('p1', 'h1'),
    ('p1', 'h2')
);

DELETE FROM hashtags
WHERE id IN ('h1', 'h2');

DELETE FROM follows
WHERE follower_id  = '8794e589-0ddb-43ce-9f92-16faafcf4ee4'
  AND following_id = '249be7cf-d419-4c54-97f2-d04107806e36';

DELETE FROM posts
WHERE id IN ('c2', 'c1');

DELETE FROM posts
WHERE id IN ('p1', 'p2');

INSERT INTO posts(id, user_id, parent_id, content) VALUES
    ('a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d', '8794e589-0ddb-43ce-9f92-16faafcf4ee4', NULL, 'Первый пост'),
    ('b2c3d4e5-f6a7-4b8c-9d0e-1f2a3b4c5d6e', '249be7cf-d419-4c54-97f2-d04107806e36', NULL, 'Второй пост'),
    ('c3d4e5f6-a7b8-4c9d-0e1f-2a3b4c5d6e7f', '249be7cf-d419-4c54-97f2-d04107806e36', 'a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d', 'Комментарий к первому посту'),
    ('d4e5f6a7-b8c9-4d0e-1f2a-3b4c5d6e7f8a', '8794e589-0ddb-43ce-9f92-16faafcf4ee4', 'c3d4e5f6-a7b8-4c9d-0e1f-2a3b4c5d6e7f', 'Ответ на комментарий');

INSERT INTO follows(follower_id, following_id) VALUES
    ('8794e589-0ddb-43ce-9f92-16faafcf4ee4', '249be7cf-d419-4c54-97f2-d04107806e36');

INSERT INTO hashtags(id, name) VALUES
    ('e5f6a7b8-c9d0-4e1f-2a3b-4c5d6e7f8a9b', 'postgresql'),
    ('f6a7b8c9-d0e1-4f2a-3b4c-5d6e7f8a9b0c', 'dbml');

INSERT INTO post_hashtags(post_id, hashtag_id) VALUES
    ('a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d', 'e5f6a7b8-c9d0-4e1f-2a3b-4c5d6e7f8a9b'),
    ('a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d', 'f6a7b8c9-d0e1-4f2a-3b4c-5d6e7f8a9b0c');
