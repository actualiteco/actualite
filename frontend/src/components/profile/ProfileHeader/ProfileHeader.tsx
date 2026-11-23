import React, { FC } from 'react';
import Avatar from '@components/shared/Avatar/Avatar';
import Button from '@components/shared/Button/Button';
import type { Profile } from '@models/Profile';
import { getDisplayName } from '@models/Profile';
import styles from './ProfileHeader.module.scss';

interface ProfileHeaderProps {
    profile: Profile;
    isOwn: boolean;
    email?: string | null;
    isFollowBusy: boolean;
    onEdit: () => void;
    onToggleFollow: () => void;
}

const ProfileHeader: FC<ProfileHeaderProps> = ({
    profile,
    isOwn,
    email,
    isFollowBusy,
    onEdit,
    onToggleFollow,
}) => (
    <section className={styles.header}>
        <Avatar url={profile.avatarUrl} name={getDisplayName(profile)} size={96} />

        <div className={styles.info}>
            <h2 className={styles.name}>{getDisplayName(profile)}</h2>
            {profile.username && <p className={styles.username}>@{profile.username}</p>}
            {profile.status && <p className={styles.status}>{profile.status}</p>}
            {isOwn && email && <p className={styles.email}>{email}</p>}

            <div className={styles.counters}>
                <span className={styles.counter}>
                    <strong>{profile.followersCount}</strong>подписчиков
                </span>
                <span className={styles.counter}>
                    <strong>{profile.followingCount}</strong>подписок
                </span>
            </div>
        </div>

        <div className={styles.actions}>
            {isOwn ? (
                <Button variant="outlined" onClick={onEdit}>
                    Редактировать
                </Button>
            ) : (
                <Button
                    variant={profile.isFollowing ? 'outlined' : 'contained'}
                    onClick={onToggleFollow}
                    disabled={isFollowBusy}
                >
                    {profile.isFollowing ? 'Отписаться' : 'Подписаться'}
                </Button>
            )}
        </div>
    </section>
);

export default ProfileHeader;
