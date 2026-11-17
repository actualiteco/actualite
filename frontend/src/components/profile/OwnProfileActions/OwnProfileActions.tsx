import { FC, useState } from 'react';
import Button from '@components/shared/Button/Button';
import AvatarUploader from '@components/profile/AvatarUploader/AvatarUploader';
import EditProfileModal from '@components/profile/EditProfileModal/EditProfileModal';
import type { Profile } from '@models/Profile';
import styles from './OwnProfileActions.module.scss';

interface OwnProfileActionsProps {
    profile: Profile;
}

const OwnProfileActions: FC<OwnProfileActionsProps> = ({ profile }) => {
    const [isEditOpen, setIsEditOpen] = useState(false);

    return (
        <div className={styles.actions}>
            <Button variant="outlined" fullWidth onClick={() => setIsEditOpen(true)}>
                Редактировать профиль
            </Button>
            <AvatarUploader />
            <EditProfileModal
                open={isEditOpen}
                profile={profile}
                onClose={() => setIsEditOpen(false)}
            />
        </div>
    );
};

export default OwnProfileActions;
