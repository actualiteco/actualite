import { useEffect } from 'react';
import { useParams } from 'react-router-dom';
import AppLayout from '@components/layout/AppLayout/AppLayout';
import ProfileHeader from '@components/profile/ProfileHeader/ProfileHeader';
import ProfileSkeleton from '@components/profile/ProfileSkeleton/ProfileSkeleton';
import ProfileError from '@components/profile/ProfileError/ProfileError';
import { useAuthStore } from '@stores/useAuthStore';
import { useProfileStore } from '@stores/useProfileStore';

const ProfilePage = () => {
    const { id } = useParams<{ id: string }>();
    const userId = useAuthStore((state) => state.userId);

    const profile = useProfileStore((state) => state.profile);
    const status = useProfileStore((state) => state.status);
    const error = useProfileStore((state) => state.error);
    const fetchProfile = useProfileStore((state) => state.fetchProfile);
    const reset = useProfileStore((state) => state.reset);

    const isOwn = !id || id === userId;
    const targetId = isOwn ? undefined : id;

    useEffect(() => {
        fetchProfile(targetId);
        return reset;
    }, [targetId, fetchProfile, reset]);

    const renderContent = () => {
        if (status === 'error' && error) {
            return <ProfileError error={error} onRetry={() => fetchProfile(targetId)} />;
        }
        if (status === 'success' && profile) {
            return <ProfileHeader profile={profile} isOwn={isOwn} />;
        }
        return <ProfileSkeleton />;
    };

    return (
        <AppLayout title="Профиль" showBack={Boolean(id)}>
            {renderContent()}
        </AppLayout>
    );
};

export default ProfilePage;
