import { ChangeEvent, useRef, useState } from 'react';
import { toast } from 'react-toastify';
import Button from '@components/shared/Button/Button';
import AvatarCropModal from '@components/profile/AvatarCropModal/AvatarCropModal';
import { useProfileStore } from '@stores/useProfileStore';
import { parseApiError } from '@utils/error.utils';
import { validateAvatarFile } from '@utils/image.utils';

const AvatarUploader = () => {
    const uploadAvatar = useProfileStore((state) => state.uploadAvatar);
    const inputRef = useRef<HTMLInputElement | null>(null);
    const [file, setFile] = useState<File | null>(null);

    const handleSelect = (event: ChangeEvent<HTMLInputElement>) => {
        const selected = event.target.files?.[0];
        event.target.value = '';

        if (!selected) {
            return;
        }

        const error = validateAvatarFile(selected);
        if (error) {
            toast.error(error);
            return;
        }
        setFile(selected);
    };

    const handleConfirm = async (blob: Blob) => {
        try {
            await uploadAvatar(blob);
            toast.success('Аватар обновлён');
            setFile(null);
        } catch (error) {
            const { status, isNetwork } = parseApiError(error);
            if (isNetwork) {
                return;
            }
            toast.error(
                status === 400
                    ? 'Не удалось загрузить аватар: проверьте файл'
                    : 'Не удалось загрузить аватар. Попробуйте ещё раз',
            );
        }
    };

    return (
        <>
            <Button variant="outlined" fullWidth onClick={() => inputRef.current?.click()}>
                Сменить фото
            </Button>
            <input
                ref={inputRef}
                type="file"
                accept="image/png,image/jpeg"
                hidden
                onChange={handleSelect}
            />
            <AvatarCropModal file={file} onCancel={() => setFile(null)} onConfirm={handleConfirm} />
        </>
    );
};

export default AvatarUploader;
