import React, { FC, useEffect, useRef, useState } from 'react';
import { toast } from 'react-toastify';
import Avatar from '@components/shared/Avatar/Avatar';
import Button from '@components/shared/Button/Button';
import Icon from '@components/shared/Icon/Icon';
import { PostService } from '@services/PostService';
import { useAuthStore } from '@stores/useAuthStore';
import type { Post } from '@models/api.types';
import { getApiErrorMessage, parseApiError } from '@utils/error.utils';
import { MB, validateImageFile } from '@utils/file.utils';
import styles from './PostComposer.module.scss';

export const POST_MAX_LENGTH = 280;
const POST_IMAGE_MAX_BYTES = 5 * MB;

interface PostComposerProps {
    placeholder: string;
    submitLabel: string;
    parentId?: string;
    onCreated: (post: Post) => void;
}

const PostComposer: FC<PostComposerProps> = ({ placeholder, submitLabel, parentId, onCreated }) => {
    const { userId, username, displayName, avatarUrl } = useAuthStore();

    const [text, setText] = useState('');
    const [image, setImage] = useState<File | null>(null);
    const [previewUrl, setPreviewUrl] = useState<string | null>(null);
    const [imageError, setImageError] = useState<string | null>(null);
    const [isSubmitting, setIsSubmitting] = useState(false);
    const fileInput = useRef<HTMLInputElement>(null);

    useEffect(() => {
        if (!image) {
            setPreviewUrl(null);
            return undefined;
        }
        const url = URL.createObjectURL(image);
        setPreviewUrl(url);
        return () => URL.revokeObjectURL(url);
    }, [image]);

    const length = [...text].length;
    const trimmed = text.trim();
    const canSubmit = trimmed.length > 0 && length <= POST_MAX_LENGTH && !isSubmitting;

    const handleFile = (event: React.ChangeEvent<HTMLInputElement>) => {
        const file = event.target.files?.[0];
        event.target.value = '';
        if (!file) {
            return;
        }
        const error = validateImageFile(file, POST_IMAGE_MAX_BYTES);
        setImageError(error);
        setImage(error ? null : file);
    };

    const removeImage = () => {
        setImage(null);
        setImageError(null);
    };

    const handleSubmit = async (event: React.FormEvent) => {
        event.preventDefault();
        if (!canSubmit) {
            return;
        }
        setIsSubmitting(true);
        try {
            const imageUrl = image ? await PostService.uploadImage(image) : undefined;
            const created = await PostService.create({
                content: trimmed,
                image_url: imageUrl,
                parent_id: parentId,
            });
            onCreated({
                id: created.id,
                content: created.content ?? trimmed,
                created_at: created.created_at ?? new Date().toISOString(),
                image_url: created.image_url ?? imageUrl,
                likes_count: created.likes_count ?? 0,
                is_liked: false,
                parent_id: created.parent_id ?? parentId,
                user_id: created.user_id ?? userId ?? undefined,
                username: username ?? undefined,
                display_name: displayName ?? username ?? undefined,
                avatar_url: avatarUrl ?? undefined,
            });
            setText('');
            setImage(null);
            setImageError(null);
        } catch (error) {
            if (parseApiError(error).status === 429) {
                toast.error(
                    parentId
                        ? 'Достигнут лимит комментариев (50 в час). Попробуйте позже'
                        : 'Достигнут лимит публикаций (30 в час). Попробуйте позже',
                );
            } else {
                toast.error(getApiErrorMessage(error, 'Не удалось опубликовать'));
            }
        } finally {
            setIsSubmitting(false);
        }
    };

    return (
        <div className={styles.composer}>
            <Avatar url={avatarUrl} name={displayName ?? username} />
            <form className={styles.form} onSubmit={handleSubmit}>
                <textarea
                    className={styles.textarea}
                    value={text}
                    onChange={(event) => setText(event.target.value)}
                    placeholder={placeholder}
                    aria-label={placeholder}
                    rows={3}
                />

                {previewUrl && (
                    <div className={styles.preview}>
                        <img className={styles.previewImage} src={previewUrl} alt="Выбранное фото" />
                        <button
                            type="button"
                            className={styles.removeImage}
                            onClick={removeImage}
                            aria-label="Убрать фото"
                        >
                            <Icon name="close" size={14} />
                        </button>
                    </div>
                )}
                {imageError && <span className={styles.error}>{imageError}</span>}

                <div className={styles.footer}>
                    <div className={styles.tools}>
                        <input
                            ref={fileInput}
                            className={styles.hiddenInput}
                            type="file"
                            accept=".jpg,.jpeg,.png,image/jpeg,image/png"
                            onChange={handleFile}
                        />
                        <Button
                            variant="outlined"
                            onClick={() => fileInput.current?.click()}
                            disabled={isSubmitting}
                        >
                            <Icon name="image" size={16} />
                            Фото
                        </Button>
                        <span
                            className={`${styles.counter} ${
                                length > POST_MAX_LENGTH ? styles.counterOver : ''
                            }`}
                        >
                            {length}/{POST_MAX_LENGTH}
                        </span>
                    </div>
                    <Button type="submit" disabled={!canSubmit} isLoading={isSubmitting}>
                        {submitLabel}
                    </Button>
                </div>
            </form>
        </div>
    );
};

export default PostComposer;
