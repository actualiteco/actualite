import { FC } from 'react';
import Avatar from '@components/shared/Avatar/Avatar';
import type { Post } from '@models/Post';
import { getDisplayName } from '@models/Profile';
import { formatRelativeTime } from '@utils/time.utils';
import styles from './PostCard.module.scss';

interface PostCardProps {
    post: Post;
}

const PostCard: FC<PostCardProps> = ({ post }) => {
    const name = getDisplayName(post.author);

    return (
        <article className={styles.card}>
            <Avatar name={name} src={post.author.avatarUrl} size={40} />

            <div className={styles.body}>
                <div className={styles.meta}>
                    <span className={styles.name}>{name}</span>
                    <span className={styles.username}>@{post.author.username}</span>
                    <span className={styles.time}>{formatRelativeTime(post.createdAt)}</span>
                </div>

                <p className={styles.content}>{post.content}</p>

                {post.imageUrl && (
                    <img className={styles.image} src={post.imageUrl} alt="" loading="lazy" />
                )}

                <div className={`${styles.footer} ${post.isLiked ? styles.liked : ''}`}>
                    <span aria-label="Лайки">♥ {post.likesCount}</span>
                </div>
            </div>
        </article>
    );
};

export default PostCard;
