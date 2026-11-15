import { FC, ReactNode } from 'react';
import Avatar from '@components/shared/Avatar/Avatar';
import type { Profile } from '@models/Profile';
import { getDisplayName } from '@models/Profile';
import styles from './ProfileHeader.module.scss';

interface ProfileHeaderProps {
    profile: Profile;
    isOwn: boolean;
    actions?: ReactNode;
}

const ProfileHeader: FC<ProfileHeaderProps> = ({ profile, isOwn, actions }) => {
    const name = getDisplayName(profile);
    const status = profile.status.trim();

    return (
        <section className={styles.header}>
            <Avatar name={name} src={profile.avatarUrl} />

            <div className={styles.info}>
                <h2 className={styles.name}>{name}</h2>
                <span className={styles.username}>@{profile.username}</span>
                {status && <p className={styles.status}>{status}</p>}
                {isOwn && profile.email && <span className={styles.email}>{profile.email}</span>}

                <div className={styles.stats}>
                    <div className={styles.stat}>
                        <span className={styles.statValue}>{profile.followersCount}</span>
                        <span className={styles.statLabel}>Подписчики</span>
                    </div>
                    <div className={styles.stat}>
                        <span className={styles.statValue}>{profile.followingCount}</span>
                        <span className={styles.statLabel}>Подписки</span>
                    </div>
                </div>
            </div>

            {actions && <div className={styles.actions}>{actions}</div>}
        </section>
    );
};

export default ProfileHeader;
