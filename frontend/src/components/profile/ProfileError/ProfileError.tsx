import { FC } from 'react';
import Button from '@components/shared/Button/Button';
import EmptyState from '@components/shared/EmptyState/EmptyState';
import type { ProfileLoadError } from '@stores/useProfileStore';
import styles from './ProfileError.module.scss';

interface ProfileErrorProps {
    error: ProfileLoadError;
    onRetry: () => void;
}

const MESSAGES: Record<ProfileLoadError, { title: string; text: string }> = {
    not_found: {
        title: 'Пользователь не найден',
        text: 'Возможно, профиль был удалён или ссылка неверна.',
    },
    forbidden: {
        title: 'Нет доступа',
        text: 'У вас нет прав для просмотра этого профиля.',
    },
    network: {
        title: 'Нет соединения',
        text: 'Не удаётся подключиться к серверу. Проверьте интернет и попробуйте ещё раз.',
    },
    unknown: {
        title: 'Что-то пошло не так',
        text: 'Не удалось загрузить профиль. Попробуйте ещё раз.',
    },
};

const RETRYABLE: ProfileLoadError[] = ['network', 'unknown'];

const ProfileError: FC<ProfileErrorProps> = ({ error, onRetry }) => {
    const { title, text } = MESSAGES[error];

    return (
        <EmptyState title={title}>
            <p>{text}</p>
            {RETRYABLE.includes(error) && (
                <div className={styles.actions}>
                    <Button onClick={onRetry}>Повторить</Button>
                </div>
            )}
        </EmptyState>
    );
};

export default ProfileError;
