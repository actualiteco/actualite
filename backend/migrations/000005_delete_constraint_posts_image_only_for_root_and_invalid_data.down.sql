DELETE FROM post_hashtags
WHERE post_id = 'a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d'
  AND hashtag_id IN (
    'e5f6a7b8-c9d0-4e1f-2a3b-4c5d6e7f8a9b',
    'f6a7b8c9-d0e1-4f2a-3b4c-5d6e7f8a9b0c'
  );

DELETE FROM hashtags
WHERE id IN (
  'e5f6a7b8-c9d0-4e1f-2a3b-4c5d6e7f8a9b',
  'f6a7b8c9-d0e1-4f2a-3b4c-5d6e7f8a9b0c'
);

DELETE FROM follows
WHERE follower_id = '8794e589-0ddb-43ce-9f92-16faafcf4ee4'
  AND following_id = '249be7cf-d419-4c54-97f2-d04107806e36';

DELETE FROM posts
WHERE id = 'd4e5f6a7-b8c9-4d0e-1f2a-3b4c5d6e7f8a';

DELETE FROM posts
WHERE id = 'c3d4e5f6-a7b8-4c9d-0e1f-2a3b4c5d6e7f';

DELETE FROM posts
WHERE id IN (
  'a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d', 
  'b2c3d4e5-f6a7-4b8c-9d0e-1f2a3b4c5d6e'
);

INSERT INTO posts(id, user_id, parent_id, content) VALUES
    ('p1', '8794e589-0ddb-43ce-9f92-16faafcf4ee4', NULL, 'Первый пост'),
    ('p2', '249be7cf-d419-4c54-97f2-d04107806e36', NULL, 'Второй пост'),
    ('c1', '249be7cf-d419-4c54-97f2-d04107806e36', 'p1', 'Комментарий к первому посту'),
    ('c2', '8794e589-0ddb-43ce-9f92-16faafcf4ee4', 'c1', 'Ответ на комментарий');

INSERT INTO follows(follower_id, following_id) VALUES
    ('8794e589-0ddb-43ce-9f92-16faafcf4ee4', '249be7cf-d419-4c54-97f2-d04107806e36');

INSERT INTO hashtags(id, name) VALUES
    ('h1', 'postgresql'),
    ('h2', 'dbml');

INSERT INTO post_hashtags(post_id, hashtag_id) VALUES
    ('p1', 'h1'),
    ('p1', 'h2');

ALTER TABLE posts
    ADD CONSTRAINT posts_image_only_for_root
        CHECK (parent_id IS NULL OR image_url IS NULL);
