import styles from './ProfileSkeleton.module.scss';

const ProfileSkeleton = () => (
    <div className={styles.skeleton} role="status" aria-label="Загрузка профиля">
        <div className={`${styles.block} ${styles.avatar}`} />
        <div className={styles.lines}>
            <div className={`${styles.block} ${styles.name}`} />
            <div className={`${styles.block} ${styles.username}`} />
            <div className={`${styles.block} ${styles.status}`} />
            <div className={`${styles.block} ${styles.stats}`} />
        </div>
    </div>
);

export default ProfileSkeleton;
