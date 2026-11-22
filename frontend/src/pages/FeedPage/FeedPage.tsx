import React from 'react';
import AppLayout from '@components/layout/AppLayout/AppLayout';
import PostComposer from '@components/posts/PostComposer';
import PostList from '@components/posts/PostList';
import { usePaginatedPosts } from '@hooks/usePaginatedPosts';
import { PostService } from '@services/PostService';
import { useAuthStore } from '@stores/useAuthStore';

const FeedPage = () => {
    const userId = useAuthStore((state) => state.userId);
    const feed = usePaginatedPosts(PostService.getFeed, []);

    return (
        <AppLayout title="Главная">
            <PostComposer
                placeholder="Что нового?"
                submitLabel="Опубликовать"
                onCreated={feed.prepend}
            />
            <PostList
                posts={feed.posts}
                isLoading={feed.isLoading}
                error={feed.error}
                hasMore={feed.hasMore}
                loadedOnce={feed.loadedOnce}
                currentUserId={userId}
                emptyTitle="Пока здесь пусто"
                emptyText="Подпишитесь на авторов или опубликуйте первый пост."
                showParentLink
                onLoadMore={feed.loadMore}
                onPostChange={feed.update}
                onPostDeleted={feed.remove}
            />
        </AppLayout>
    );
};

export default FeedPage;
