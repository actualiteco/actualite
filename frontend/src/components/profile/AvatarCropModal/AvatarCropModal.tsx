import { FC, SyntheticEvent, useEffect, useRef, useState } from 'react';
import ReactCrop from 'react-image-crop';
import type { Crop, PixelCrop } from 'react-image-crop';
import 'react-image-crop/dist/ReactCrop.css';
import { toast } from 'react-toastify';
import Modal from '@components/shared/Modal/Modal';
import Button from '@components/shared/Button/Button';
import { cropToBlob } from '@utils/image.utils';
import styles from './AvatarCropModal.module.scss';

interface AvatarCropModalProps {
    file: File | null;
    onCancel: () => void;
    onConfirm: (blob: Blob) => Promise<void>;
}

interface CropperBodyProps {
    file: File;
    onCancel: () => void;
    onConfirm: (blob: Blob) => Promise<void>;
}

const INITIAL_COVERAGE = 0.9;
const MIN_CROP_SIZE = 48;

const CropperBody: FC<CropperBodyProps> = ({ file, onCancel, onConfirm }) => {
    const imageRef = useRef<HTMLImageElement | null>(null);
    const [url, setUrl] = useState<string | null>(null);
    const [loadFailed, setLoadFailed] = useState(false);
    const [crop, setCrop] = useState<Crop>();
    const [completedCrop, setCompletedCrop] = useState<PixelCrop | null>(null);
    const [isSaving, setIsSaving] = useState(false);

    useEffect(() => {
        const objectUrl = URL.createObjectURL(file);
        setUrl(objectUrl);
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
        if (!image || !completedCrop || completedCrop.width === 0) {
            return;
        }

        const scaleX = image.naturalWidth / image.width;
        const scaleY = image.naturalHeight / image.height;

        setIsSaving(true);

        let blob: Blob;
        try {
            blob = await cropToBlob(
                image,
                {
                    x: completedCrop.x * scaleX,
                    y: completedCrop.y * scaleY,
                    size: Math.min(completedCrop.width * scaleX, completedCrop.height * scaleY),
                },
                file.type,
            );
        } catch {
            toast.error('Не удалось обработать изображение');
            setIsSaving(false);
            return;
        }

        await onConfirm(blob);
        setIsSaving(false);
    };

    if (loadFailed) {
        return (
            <div className={styles.body}>
                <h2 className={styles.title}>Аватар</h2>
                <p className={styles.note}>Не удалось открыть изображение</p>
                <div className={styles.actions}>
                    <Button variant="outlined" onClick={onCancel}>
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
                            alt=""
                            onLoad={handleImageLoad}
                            onError={() => setLoadFailed(true)}
                        />
                    </ReactCrop>
                )}
            </div>

            <div className={styles.actions}>
                <Button variant="outlined" onClick={onCancel}>
                    Отмена
                </Button>
                <Button
                    onClick={handleConfirm}
                    isLoading={isSaving}
                    disabled={!completedCrop || completedCrop.width === 0}
                >
                    Сохранить
                </Button>
            </div>
        </div>
    );
};

const AvatarCropModal: FC<AvatarCropModalProps> = ({ file, onCancel, onConfirm }) => (
    <Modal open={file !== null} onClose={onCancel}>
        {file && <CropperBody file={file} onCancel={onCancel} onConfirm={onConfirm} />}
    </Modal>
);

export default AvatarCropModal;
