import React, { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { toast } from 'react-toastify';
import AppLayout from '@components/layout/AppLayout/AppLayout';
import PostList from '@components/posts/PostList';
import ProfileEditModal from '@components/profile/EditProfileModal/EditProfileModal';
import ProfileError, { ProfileErrorType } from '@components/profile/ProfileError/ProfileError';
import ProfileHeader from '@components/profile/ProfileHeader/ProfileHeader';
import ProfileSkeleton from '@components/profile/ProfileSkeleton/ProfileSkeleton';
import { usePaginatedPosts } from '@hooks/usePaginatedPosts';
import { PostService } from '@services/PostService';
import { ProfileService } from '@services/ProfileService';
import { useAuthStore } from '@stores/useAuthStore';
import type { Profile } from '@models/Profile';
import { getApiErrorMessage, parseApiError } from '@utils/error.utils';

type Status = 'loading' | 'ready' | ProfileErrorType;

const ProfilePage = () => {
    const { id } = useParams<{ id: string }>();
    const { userId, email, setProfile: saveMyProfile } = useAuthStore();
    const isOwn = !id || id === userId;

    const [profile, setProfile] = useState<Profile | null>(null);
    const [status, setStatus] = useState<Status>('loading');
    const [reloadKey, setReloadKey] = useState(0);
    const [isEditing, setIsEditing] = useState(false);
    const [isFollowBusy, setIsFollowBusy] = useState(false);

    useEffect(() => {
        let cancelled = false;

        setStatus('loading');
        setProfile(null);

        const request = isOwn
            ? ProfileService.getMe()
            : ProfileService.getById(id as string);

        request
            .then((loaded) => {
                if (cancelled) return;

                setProfile(loaded);
                setStatus('ready');

                if (isOwn) {
                    saveMyProfile(loaded);
                }
            })
            .catch((error: unknown) => {
                if (cancelled) return;

                const { status: code, isNetworkError } = parseApiError(error);

                if (code === 404 || code === 400) {
                    setStatus('notFound');
                } else if (code === 403) {
                    setStatus('forbidden');
                } else if (isNetworkError) {
                    setStatus('network');
                } else {
                    setStatus('unknown');
                }
            });

        return () => {
            cancelled = true;
        };
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [id, isOwn, reloadKey]);

    const profileId = profile?.id || undefined;
    const posts = usePaginatedPosts(
        (page, pageSize) =>
            PostService.getUserPosts(profileId as string, page, pageSize),
        profileId,
        Boolean(profileId),
    );

    const handleToggleFollow = async () => {
        if (!profile?.id || isFollowBusy) {
            return;
        }

        setIsFollowBusy(true);

        const before = profile;
        const wasFollowing = before.isFollowing === true;

        setProfile({
            ...before,
            isFollowing: !wasFollowing,
            followersCount: Math.max(
                0,
                before.followersCount + (wasFollowing ? -1 : 1),
            ),
        });

        try {
            if (wasFollowing) {
                await ProfileService.unfollow(before.id);
            } else {
                await ProfileService.follow(before.id);
            }
        } catch (error) {
            setProfile((current) =>
                current?.id === before.id ? before : current,
            );

            toast.error(
                parseApiError(error).status === 429
                    ? 'Достигнут лимит подписок (50 в час). Попробуйте позже'
                    : getApiErrorMessage(error, 'Не удалось выполнить действие'),
            );
        } finally {
            setIsFollowBusy(false);
        }
    };

    const handleSaved = (updated: Profile) => {
        setProfile(updated);

        if (isOwn) {
            saveMyProfile(updated);
        }
    };

    const handleRetry = () => {
        setReloadKey((key) => key + 1);
    };

    return (
        <AppLayout title="Профиль" showBack={Boolean(id) && !isOwn}>
            {status === 'loading' && <ProfileSkeleton />}

            {status !== 'loading' && status !== 'ready' && (
                <ProfileError
                    error={status}
                    onRetry={handleRetry}
                />
            )}

            {status === 'ready' && profile && (
                <>
                    <ProfileHeader
                        profile={profile}
                        isOwn={isOwn}
                        email={email}
                        isFollowBusy={isFollowBusy}
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
                            onSaved={handleSaved}
                        />
                    )}
                </>
            )}
        </AppLayout>
    );
};

export default ProfilePage;
