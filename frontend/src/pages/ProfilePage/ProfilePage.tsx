import React, { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { toast } from 'react-toastify';

import AppLayout from '@components/layout/AppLayout/AppLayout';
import PostList from '@components/post/PostList/PostList';
import ProfileEditModal from '@components/profile/EditProfileModal/EditProfileModal';
import ProfileError from '@components/profile/ProfileError/ProfileError';
import ProfileHeader from '@components/profile/ProfileHeader/ProfileHeader';
import ProfileSkeleton from '@components/profile/ProfileSkeleton/ProfileSkeleton';

import { usePaginatedPosts } from '@hooks/usePaginatedPosts';
import { PostService } from '@services/PostService';

import { useAuthStore } from '@stores/useAuthStore';
import { useProfileStore } from '@stores/useProfileStore';

import { getApiErrorMessage, parseApiError } from '@utils/error.utils';

const ProfilePage = () => {
    const { id } = useParams<{ id: string }>();

    const userId = useAuthStore((state) => state.userId);
    const email = useAuthStore((state) => state.email);

    const profile = useProfileStore((state) => state.profile);
    const status = useProfileStore((state) => state.status);
    const error = useProfileStore((state) => state.error);
    const isFollowPending = useProfileStore((state) => state.isFollowPending);

    const fetchProfile = useProfileStore((state) => state.fetchProfile);
    const toggleFollow = useProfileStore((state) => state.toggleFollow);
    const reset = useProfileStore((state) => state.reset);

    const isOwn = !id || id === userId;

    const [reloadKey, setReloadKey] = useState(0);
    const [isEditing, setIsEditing] = useState(false);

    useEffect(() => {
        void fetchProfile(isOwn ? undefined : id);

        return () => {
            reset();
        };
    }, [id, isOwn, reloadKey, fetchProfile, reset]);

    const profileId = profile?.id || undefined;

    const posts = usePaginatedPosts(
        (page, pageSize) =>
            PostService.getUserPosts(profileId as string, page, pageSize),
        profileId,
        Boolean(profileId && status === 'success'),
    );

    const handleToggleFollow = async () => {
        try {
            await toggleFollow();
        } catch (error) {
            toast.error(
                parseApiError(error).status === 429
                    ? 'Достигнут лимит подписок (50 в час). Попробуйте позже'
                    : getApiErrorMessage(
                          error,
                          'Не удалось выполнить действие',
                      ),
            );
        }
    };

    const handleRetry = () => {
        setReloadKey((key) => key + 1);
    };

    return (
        <AppLayout title="Профиль" showBack={Boolean(id) && !isOwn}>
            {(status === 'idle' || status === 'loading') && (
                <ProfileSkeleton />
            )}

            {status === 'error' && error && (
                <ProfileError
                    error={error}
                    onRetry={handleRetry}
                />
            )}

            {status === 'success' && profile && (
                <>
                    <ProfileHeader
                        profile={profile}
                        isOwn={isOwn}
                        email={email}
                        isFollowBusy={isFollowPending}
                        onEdit={() => setIsEditing(true)}
                        onToggleFollow={handleToggleFollow}
                    />

                    <PostList
                        posts={posts.posts}
                        isLoading={posts.isLoading}
                        error={posts.error}
                        hasMore={posts.hasMore}
                        loadedOnce={posts.loadedOnce}
                        currentUserId={userId}
                        emptyTitle="Постов пока нет"
                        showParentLink
                        onLoadMore={posts.loadMore}
                        onPostChange={posts.update}
                        onPostDeleted={posts.remove}
                    />

                    {isOwn && (
                        <ProfileEditModal
                            open={isEditing}
                            profile={profile}
                            onClose={() => setIsEditing(false)}
                            onSaved={() => setIsEditing(false)}
                        />
                    )}
                </>
            )}
        </AppLayout>
    );
};

export default ProfilePage;