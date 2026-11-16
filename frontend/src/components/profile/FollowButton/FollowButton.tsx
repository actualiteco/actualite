import { toast } from 'react-toastify';
import Button from '@components/shared/Button/Button';
import { useProfileStore } from '@stores/useProfileStore';
import { parseApiError } from '@utils/error.utils';

const FollowButton = () => {
    const isFollowing = useProfileStore((state) => state.profile?.isFollowing === true);
    const isPending = useProfileStore((state) => state.isFollowPending);
    const toggleFollow = useProfileStore((state) => state.toggleFollow);

    const handleClick = async () => {
        try {
            await toggleFollow();
        } catch (error) {
            const { status, isNetwork } = parseApiError(error);
            if (isNetwork) {
                return;
            }
            toast.error(
                status === 429
                    ? 'Слишком много подписок. Попробуйте позже'
                    : 'Не удалось выполнить действие. Попробуйте ещё раз',
            );
        }
    };

    return (
        <Button
            variant={isFollowing ? 'outlined' : 'contained'}
            isLoading={isPending}
            onClick={handleClick}
        >
            {isFollowing ? 'Отписаться' : 'Подписаться'}
        </Button>
    );
};

export default FollowButton;
