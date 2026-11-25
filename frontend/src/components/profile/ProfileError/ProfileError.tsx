import { FC } from 'react';
import Button from '@components/shared/Button/Button';
import EmptyState from '@components/shared/EmptyState/EmptyState';
import styles from './ProfileError.module.scss';

export type ProfileErrorType =
    | 'notFound'
    | 'forbidden'
    | 'network'
    | 'unknown';

interface ProfileErrorProps {
    error: ProfileErrorType;
    onRetry: () => void;
}

const MESSAGES: Record<ProfileErrorType, { title: string; text: string }> = {
    notFound: {
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

const ProfileError: FC<ProfileErrorProps> = ({ error, onRetry }) => {
    const { title, text } = MESSAGES[error];
    const canRetry = error === 'network' || error === 'unknown';

    return (
        <EmptyState title={title}>
            <p>{text}</p>
            {canRetry && (
                <div className={styles.actions}>
                    <Button onClick={onRetry}>Повторить</Button>
                </div>
            )}
        </EmptyState>
    );
};

export default ProfileError;
