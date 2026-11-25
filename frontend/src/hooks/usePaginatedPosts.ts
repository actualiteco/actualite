import { useCallback, useEffect, useRef, useState } from 'react';
import type { Post, PostsPage } from '@models/api.types';
import { getApiErrorMessage } from '@utils/error.utils';

export const PAGE_SIZE = 20;

type PageFetcher = (page: number, pageSize: number) => Promise<PostsPage>;

export const usePaginatedPosts = (
    fetchPage: PageFetcher,
    resetKey: string | undefined,
    enabled = true,
) => {
    const [posts, setPosts] = useState<Post[]>([]);
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [hasMore, setHasMore] = useState(true);
    const [loadedOnce, setLoadedOnce] = useState(false);

    const fetchRef = useRef(fetchPage);
    fetchRef.current = fetchPage;

    const generation = useRef(0);
    const loadingRef = useRef(false);
    const hasMoreRef = useRef(true);
    const pageRef = useRef(0);

    const loadMore = useCallback(async () => {
        if (loadingRef.current || !hasMoreRef.current) {
            return;
        }
        loadingRef.current = true;
        setIsLoading(true);
        setError(null);
        const gen = generation.current;
        try {
            const next = pageRef.current + 1;
            const data = await fetchRef.current(next, PAGE_SIZE);
            if (gen !== generation.current) {
                return;
            }
            const received = data.posts ?? [];
            const totalPages = data.pagination?.total_pages;
            const more = totalPages !== undefined ? next < totalPages : received.length >= PAGE_SIZE;

            pageRef.current = next;
            hasMoreRef.current = more;
            setHasMore(more);
            setLoadedOnce(true);
            setPosts((prev) => {
                const known = new Set(prev.map((post) => post.id));
                return [...prev, ...received.filter((post) => !known.has(post.id))];
            });
        } catch (e) {
            if (gen === generation.current) {
                setError(getApiErrorMessage(e, 'Не удалось загрузить посты'));
            }
        } finally {
            if (gen === generation.current) {
                loadingRef.current = false;
                setIsLoading(false);
            }
        }
    }, []);

    const reset = useCallback(() => {
        generation.current += 1;
        loadingRef.current = false;
        hasMoreRef.current = true;
        pageRef.current = 0;
        setPosts([]);
        setHasMore(true);
        setLoadedOnce(false);
        setError(null);
        setIsLoading(false);
    }, []);

    useEffect(() => {
        reset();

        if (enabled) {
            void loadMore();
        }

        return () => {
            generation.current += 1;
            loadingRef.current = false;
        };
    }, [resetKey, enabled, reset, loadMore]);

    const prepend = useCallback((post: Post) => {
        setPosts((prev) => [post, ...prev.filter((item) => item.id !== post.id)]);
    }, []);

    const update = useCallback((post: Post) => {
        setPosts((prev) => prev.map((item) => (item.id === post.id ? post : item)));
    }, []);

    const remove = useCallback((postId: string) => {
        setPosts((prev) => prev.filter((item) => item.id !== postId));
    }, []);

    return { posts, isLoading, error, hasMore, loadedOnce, loadMore, prepend, update, remove };
};
