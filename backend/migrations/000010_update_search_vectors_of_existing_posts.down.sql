UPDATE posts
SET search_vector = NULL
WHERE search_vector IS NOT NULL;
