import React, { useEffect, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { toast } from 'react-toastify';
import AppLayout from '@components/layout/AppLayout/AppLayout';
import PostCard from '@components/post/PostCard/PostCard';
import PostComposer from '@components/post/PostComposer/PostComposer';
import PostList from '@components/post/PostList/PostList';
import Button from '@components/shared/Button/Button';
import EmptyState from '@components/shared/EmptyState/EmptyState';
import { usePaginatedPosts } from '@hooks/usePaginatedPosts';
import { PostService } from '@services/PostService';
import { RoutePath, getPostPath } from '@router/paths';
import { useAuthStore } from '@stores/useAuthStore';
import type { Post } from '@models/api.types';
import { parseApiError } from '@utils/error.utils';

type Status = 'loading' | 'ready' | 'deleted' | 'error';

const PostDetailPage = () => {
    const { id } = useParams<{ id: string }>();
    const navigate = useNavigate();
    const userId = useAuthStore((state) => state.userId);

    const [post, setPost] = useState<Post | null>(null);
    const [status, setStatus] = useState<Status>('loading');
    const [reloadKey, setReloadKey] = useState(0);

    useEffect(() => {
        if (!id) {
            return undefined;
        }
        let cancelled = false;
        setStatus('loading');
        setPost(null);

        PostService.getById(id)
            .then((loaded) => {
                if (cancelled) return;
                if (loaded.deleted_at) {
                    setStatus('deleted');
                    return;
                }
                setPost(loaded);
                setStatus('ready');
            })
            .catch((error) => {
                if (cancelled) return;
                if (parseApiError(error).status === 404) {
                    setStatus('deleted');
                    toast.info('Пост был удалён', { toastId: 'post-deleted' });
                } else {
                    setStatus('error');
                }
            });

        return () => {
            cancelled = true;
        };
    }, [id, reloadKey]);

    const replies = usePaginatedPosts(
        (page, pageSize) => PostService.getReplies(id as string, page, pageSize),
        id,
        status === 'ready',
    );

    const handleDeleted = () => {
        if (post?.parent_id) {
            navigate(getPostPath(post.parent_id), { replace: true });
        } else {
            navigate(RoutePath.feed, { replace: true });
        }
    };

    const title = post?.parent_id ? 'Комментарий' : 'Пост';

    return (
        <AppLayout title={title} showBack>
            {status === 'loading' && <EmptyState title="Загрузка…" />}

            {status === 'deleted' && (
                <EmptyState title="Пост удалён">
                    <p>Этого поста больше нет, но ответы на него остались в профилях авторов.</p>
                </EmptyState>
            )}

            {status === 'error' && (
                <EmptyState title="Не удалось загрузить пост">
                    <Button variant="outlined" onClick={() => setReloadKey((key) => key + 1)}>
                        Повторить
                    </Button>
                </EmptyState>
            )}

            {status === 'ready' && post && (
                <>
                    <PostCard
                        post={post}
                        currentUserId={userId}
                        showParentLink
                        onChange={setPost}
                        onDeleted={handleDeleted}
                    />
                    <PostComposer
                        placeholder="Написать комментарий"
                        submitLabel="Ответить"
                        parentId={post.id}
                        onCreated={replies.prepend}
                    />
                    <PostList
                        posts={replies.posts}
                        isLoading={replies.isLoading}
                        error={replies.error}
                        hasMore={replies.hasMore}
                        loadedOnce={replies.loadedOnce}
                        currentUserId={userId}
                        mode="button"
                        emptyTitle="Комментариев пока нет"
                        emptyText="Станьте первым, кто ответит."
                        onLoadMore={replies.loadMore}
                        onPostChange={replies.update}
                        onPostDeleted={replies.remove}
                    />
                </>
            )}
        </AppLayout>
    );
};

export default PostDetailPage;
