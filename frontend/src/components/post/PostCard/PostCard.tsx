import React, { FC, useState } from 'react';
import { Link } from 'react-router-dom';
import { toast } from 'react-toastify';
import PostActionsMenu from '@components/posts/PostActionsMenu';
import Avatar from '@components/shared/Avatar/Avatar';
import Icon from '@components/shared/Icon/Icon';
import { PostService } from '@services/PostService';
import { getPostPath, getProfilePath } from '@router/paths';
import type { Post } from '@models/api.types';
import { getApiErrorMessage } from '@utils/error.utils';
import { formatPostDate } from '@utils/time.utils';
import styles from './PostCard.module.scss';
import { truncate } from '@utils/string.utils';

interface PostCardProps {
    post: Post;
    currentUserId: string | null;
    showParentLink?: boolean;
    onChange: (post: Post) => void;
    onDeleted: (postId: string) => void;
}

const PostCard: FC<PostCardProps> = ({
    post,
    currentUserId,
    showParentLink = false,
    onChange,
    onDeleted,
}) => {
    const [isLikeBusy, setIsLikeBusy] = useState(false);
    const [isDeleting, setIsDeleting] = useState(false);

    if (!post.id) {
        return null;
    }
    const postId = post.id;
    const isOwn = Boolean(currentUserId) && post.user_id === currentUserId;
    const authorName = post.display_name || post.username || '';

    const handleLike = async () => {
        if (isLikeBusy) {
            return;
        }
        setIsLikeBusy(true);
        const wasLiked = Boolean(post.is_liked);
        onChange({
            ...post,
            is_liked: !wasLiked,
            likes_count: Math.max(0, (post.likes_count ?? 0) + (wasLiked ? -1 : 1)),
        });
        try {
            if (wasLiked) {
                await PostService.unlike(postId);
            } else {
                await PostService.like(postId);
            }
        } catch (error) {
            onChange(post);
            toast.error(getApiErrorMessage(error, 'Не удалось изменить лайк'));
        } finally {
            setIsLikeBusy(false);
        }
    };

    // F_Post_5: автор удаляет свой пост
    const handleDelete = async () => {
        setIsDeleting(true);
        try {
            await PostService.remove(postId);
            toast.success('Пост удалён');
            onDeleted(postId);
        } catch (error) {
            setIsDeleting(false);
            toast.error(getApiErrorMessage(error, 'Не удалось удалить пост'));
        }
    };

    const shouldShowParentLink = showParentLink && Boolean(post.parent_id);

    return (
        <article className={styles.card}>
            {post.user_id ? (
                <Link to={getProfilePath(post.user_id)} aria-label={`Профиль ${authorName}`}>
                    <Avatar url={post.avatar_url} name={authorName} />
                </Link>
            ) : (
                <Avatar url={post.avatar_url} name={authorName} />
            )}

            <div className={styles.body}>
                <div className={styles.top}>
                    <div className={styles.head}>
                        {post.user_id ? (
                            <Link to={getProfilePath(post.user_id)} className={styles.author}>
                                {authorName}
                            </Link>
                        ) : (
                            <span className={styles.author}>{authorName}</span>
                        )}
                        {post.username && <span className={styles.meta}>@{post.username}</span>}
                        <span className={styles.meta}>{formatPostDate(post.created_at)}</span>
                    </div>
                    {isOwn && <PostActionsMenu isDeleting={isDeleting} onDelete={handleDelete} />}
                </div>

                {shouldShowParentLink && (
                    <Link to={getPostPath(post.parent_id!)} className={styles.replyTo}>
                        <Icon name="reply" size={12} />
                        В ответ на пост: "{truncate(post.parent_id!, 15)}"
                    </Link>
                )}

                <Link
                    to={getPostPath(postId)}
                    className={`${styles.content} ${
                        shouldShowParentLink ? styles.contentAfterReply : ''
                    }`}
                >
                    {post.content}
                </Link>

                {post.image_url && (
                    <img className={styles.image} src={post.image_url} alt="Изображение к посту" />
                )}

                <div className={styles.actions}>
                    <button
                        type="button"
                        className={`${styles.action} ${post.is_liked ? styles.liked : ''}`}
                        onClick={handleLike}
                        disabled={isLikeBusy}
                        aria-pressed={Boolean(post.is_liked)}
                        aria-label={post.is_liked ? 'Убрать лайк' : 'Поставить лайк'}
                    >
                        <Icon name="heart" size={16} filled={Boolean(post.is_liked)} />
                        <span>{post.likes_count ?? 0}</span>
                    </button>

                    <Link
                        to={getPostPath(postId)}
                        className={styles.action}
                        aria-label="Комментарии"
                    >
                        <Icon name="comment" size={16} />
                        <span>Ответы</span>
                    </Link>

                </div>
            </div>
        </article>
    );
};

export default PostCard;
