import { FC, useEffect } from 'react';
import Button from '@components/shared/Button/Button';
import EmptyState from '@components/shared/EmptyState/EmptyState';
import PostCard from '@components/post/PostCard/PostCard';
import { useInfiniteScroll } from '@hooks/useInfiniteScroll';
import { useUserPostsStore } from '@stores/useUserPostsStore';
import styles from './UserPosts.module.scss';

interface UserPostsProps {
    userId: string;
    isOwn: boolean;
}

const SKELETON_ITEMS = [0, 1, 2];

const UserPosts: FC<UserPostsProps> = ({ userId, isOwn }) => {
    const posts = useUserPostsStore((state) => state.posts);
    const status = useUserPostsStore((state) => state.status);
    const error = useUserPostsStore((state) => state.error);
    const hasMore = useUserPostsStore((state) => state.hasMore);
    const isLoadingMore = useUserPostsStore((state) => state.isLoadingMore);
    const loadMoreFailed = useUserPostsStore((state) => state.loadMoreFailed);
    const load = useUserPostsStore((state) => state.load);
    const loadMore = useUserPostsStore((state) => state.loadMore);
    const reset = useUserPostsStore((state) => state.reset);

    useEffect(() => {
        load(userId);
        return reset;
    }, [userId, load, reset]);

    const sentinelRef = useInfiniteScroll({
        onReach: loadMore,
        enabled: status === 'success' && hasMore && !isLoadingMore && !loadMoreFailed,
    });

    if (status === 'error') {
        const text =
            error === 'network'
                ? 'Не удаётся подключиться к серверу. Проверьте интернет и попробуйте ещё раз.'
                : 'Не удалось загрузить посты. Попробуйте ещё раз.';

        return (
            <div className={styles.section}>
                <EmptyState title="Не удалось загрузить посты">
                    <p>{text}</p>
                    <div className={styles.actions}>
                        <Button onClick={() => load(userId)}>Повторить</Button>
                    </div>
                </EmptyState>
            </div>
        );
    }

    if (status !== 'success') {
        return (
            <div className={styles.section} role="status" aria-label="Загрузка постов">
                {SKELETON_ITEMS.map((item) => (
                    <div key={item} className={styles.skeletonItem} />
                ))}
            </div>
        );
    }

    if (posts.length === 0) {
        return (
            <div className={styles.section}>
                <EmptyState title={isOwn ? 'У вас пока нет постов' : 'Пока нет постов'} />
            </div>
        );
    }

    return (
        <div className={styles.section}>
            {posts.map((post) => (
                <PostCard key={post.id} post={post} />
            ))}

            {isLoadingMore && <p className={styles.note}>Загрузка…</p>}

            {loadMoreFailed && (
                <div className={styles.footer}>
                    <p className={styles.note}>Не удалось загрузить ещё посты</p>
                    <Button variant="outlined" onClick={loadMore}>
                        Повторить
                    </Button>
                </div>
            )}

            <div ref={sentinelRef} className={styles.sentinel} />
        </div>
    );
};

export default UserPosts;
