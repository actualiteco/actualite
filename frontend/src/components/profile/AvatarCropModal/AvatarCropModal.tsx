import React, { FC, SyntheticEvent, useEffect, useRef, useState } from 'react';
import ReactCrop from 'react-image-crop';
import type { Crop, PixelCrop } from 'react-image-crop';
import 'react-image-crop/dist/ReactCrop.css';
import { toast } from 'react-toastify';
import Modal from '@components/shared/Modal/Modal';
import Button from '@components/shared/Button/Button';
import { AVATAR_MAX_BYTES, cropToBlob } from '@utils/image.utils';
import styles from './AvatarCropModal.module.scss';

export { AVATAR_MAX_BYTES } from '@utils/image.utils';

interface AvatarCropModalProps {
    file: File | null;
    isSaving: boolean;
    onCancel: () => void;
    onConfirm: (blob: Blob, fileName: string) => void;
}

const INITIAL_COVERAGE = 0.9;
const MIN_CROP_SIZE = 48;

const CropperBody: FC<Omit<AvatarCropModalProps, 'file'> & { file: File }> = ({
    file,
    isSaving,
    onCancel,
    onConfirm,
}) => {
    const imageRef = useRef<HTMLImageElement | null>(null);
    const [url, setUrl] = useState<string | null>(null);
    const [loadFailed, setLoadFailed] = useState(false);
    const [crop, setCrop] = useState<Crop>();
    const [completedCrop, setCompletedCrop] = useState<PixelCrop | null>(null);

    useEffect(() => {
        const objectUrl = URL.createObjectURL(file);
        setUrl(objectUrl);
        setLoadFailed(false);
        setCrop(undefined);
        setCompletedCrop(null);
        return () => URL.revokeObjectURL(objectUrl);
    }, [file]);

    const handleImageLoad = (event: SyntheticEvent<HTMLImageElement>) => {
        const { width, height } = event.currentTarget;
        const size = Math.min(width, height) * INITIAL_COVERAGE;
        const initial: PixelCrop = {
            unit: 'px',
            x: (width - size) / 2,
            y: (height - size) / 2,
            width: size,
            height: size,
        };
        setCrop(initial);
        setCompletedCrop(initial);
    };

    const handleConfirm = async () => {
        const image = imageRef.current;
        if (!image || !completedCrop || completedCrop.width < MIN_CROP_SIZE) {
            return;
        }

        const scaleX = image.naturalWidth / image.width;
        const scaleY = image.naturalHeight / image.height;
        const cropSize = Math.min(completedCrop.width * scaleX, completedCrop.height * scaleY);

        try {
            const blob = await cropToBlob(
                image,
                {
                    x: completedCrop.x * scaleX,
                    y: completedCrop.y * scaleY,
                    size: cropSize,
                },
                file.type === 'image/png' ? 'image/png' : 'image/jpeg',
            );

            if (blob.size > AVATAR_MAX_BYTES) {
                toast.error('После обрезки файл превышает 2 МБ. Выберите другое фото');
                return;
            }

            onConfirm(blob, file.type === 'image/png' ? 'avatar.png' : 'avatar.jpg');
        } catch {
            toast.error('Не удалось обработать изображение');
        }
    };

    if (loadFailed) {
        return (
            <div className={styles.body}>
                <h2 className={styles.title}>Аватар</h2>
                <p className={styles.note}>Не удалось открыть изображение</p>
                <div className={styles.actions}>
                    <Button variant="outlined" onClick={onCancel} disabled={isSaving}>
                        Закрыть
                    </Button>
                </div>
            </div>
        );
    }

    return (
        <div className={styles.body}>
            <h2 className={styles.title}>Аватар</h2>
            <p className={styles.hint}>Перетащите рамку или потяните за её углы</p>

            <div className={styles.cropArea}>
                {url && (
                    <ReactCrop
                        className={styles.crop}
                        crop={crop}
                        aspect={1}
                        minWidth={MIN_CROP_SIZE}
                        keepSelection
                        onChange={(pixelCrop) => setCrop(pixelCrop)}
                        onComplete={(pixelCrop) => setCompletedCrop(pixelCrop)}
                    >
                        <img
                            ref={imageRef}
                            className={styles.image}
                            src={url}
                            alt="Выбор области аватара"
                            onLoad={handleImageLoad}
                            onError={() => setLoadFailed(true)}
                        />
                    </ReactCrop>
                )}
            </div>

            <div className={styles.actions}>
                <Button variant="outlined" onClick={onCancel} disabled={isSaving}>
                    Отмена
                </Button>
                <Button
                    onClick={handleConfirm}
                    isLoading={isSaving}
                    disabled={!completedCrop || completedCrop.width < MIN_CROP_SIZE || isSaving}
                >
                    Сохранить
                </Button>
            </div>
        </div>
    );
};

const AvatarCropModal: FC<AvatarCropModalProps> = ({ file, isSaving, onCancel, onConfirm }) => (
    <Modal open={file !== null} onClose={onCancel}>
        {file && (
            <CropperBody
                file={file}
                isSaving={isSaving}
                onCancel={onCancel}
                onConfirm={onConfirm}
            />
        )}
    </Modal>
);

export default AvatarCropModal;
