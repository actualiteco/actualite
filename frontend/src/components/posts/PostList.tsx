import React, { FC, useEffect, useRef } from 'react';
import Button from '@components/shared/Button/Button';
import EmptyState from '@components/shared/EmptyState/EmptyState';
import PostCard from '@components/posts/PostCard';
import type { Post } from '@models/api.types';
import styles from './PostList.module.scss';

interface PostListProps {
    posts: Post[];
    isLoading: boolean;
    error: string | null;
    hasMore: boolean;
    loadedOnce: boolean;
    currentUserId: string | null;
    mode?: 'infinite' | 'button';
    emptyTitle: string;
    emptyText?: string;
    showParentLink?: boolean;
    onLoadMore: () => void;
    onPostChange: (post: Post) => void;
    onPostDeleted: (postId: string) => void;
}

const Sentinel: FC<{ onVisible: () => void; watch: number }> = ({ onVisible, watch }) => {
    const ref = useRef<HTMLDivElement>(null);
    const callback = useRef(onVisible);
    callback.current = onVisible;

    useEffect(() => {
        const element = ref.current;
        if (!element) {
            return undefined;
        }
        const observer = new IntersectionObserver(
            (entries) => {
                if (entries.some((entry) => entry.isIntersecting)) {
                    callback.current();
                }
            },
            { rootMargin: '200px' },
        );
        observer.observe(element);
        return () => observer.disconnect();
    }, [watch]);

    return <div ref={ref} className={styles.sentinel} aria-hidden="true" />;
};

const PostList: FC<PostListProps> = ({
    posts,
    isLoading,
    error,
    hasMore,
    loadedOnce,
    currentUserId,
    mode = 'infinite',
    emptyTitle,
    emptyText,
    showParentLink = false,
    onLoadMore,
    onPostChange,
    onPostDeleted,
}) => {
    const isEmpty = loadedOnce && posts.length === 0;

    return (
        <div>
            {posts.map((post) => (
                <PostCard
                    key={post.id}
                    post={post}
                    currentUserId={currentUserId}
                    showParentLink={showParentLink}
                    onChange={onPostChange}
                    onDeleted={onPostDeleted}
                />
            ))}

            {isEmpty && (
                <EmptyState title={emptyTitle}>{emptyText && <p>{emptyText}</p>}</EmptyState>
            )}

            {isLoading && <div className={styles.status}>Загрузка…</div>}

            {error && (
                <div className={styles.more}>
                    <div>
                        <p className={`${styles.status} ${styles.error}`}>{error}</p>
                        <Button variant="outlined" onClick={onLoadMore}>
                            Повторить
                        </Button>
                    </div>
                </div>
            )}

            {!error && !isLoading && hasMore && mode === 'infinite' && posts.length > 0 && (
                <Sentinel onVisible={onLoadMore} watch={posts.length} />
            )}

            {!error && !isLoading && hasMore && mode === 'button' && posts.length > 0 && (
                <div className={styles.more}>
                    <Button variant="outlined" onClick={onLoadMore}>
                        Показать ещё
                    </Button>
                </div>
            )}
        </div>
    );
};

export default PostList;
