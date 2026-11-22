import React, { FC, useEffect, useRef, useState } from 'react';
import Button from '@components/shared/Button/Button';
import Modal from '@components/shared/Modal/Modal';
import { MB } from '@utils/file.utils';
import styles from './AvatarCropModal.module.scss';

export const AVATAR_CROP_SIZE = 400;
export const AVATAR_MAX_BYTES = 2 * MB;
const VIEW_MAX = 420;

interface AvatarCropModalProps {
    file: File | null;
    isSaving: boolean;
    onCancel: () => void;
    onConfirm: (blob: Blob, fileName: string) => void;
}

interface Point {
    x: number;
    y: number;
}

interface DragState {
    pointerX: number;
    pointerY: number;
    origin: Point;
}

const clamp = (value: number, min: number, max: number) => Math.min(Math.max(value, min), max);

const AvatarCropModal: FC<AvatarCropModalProps> = ({ file, isSaving, onCancel, onConfirm }) => {
    const [url, setUrl] = useState<string | null>(null);
    const [size, setSize] = useState<{ width: number; height: number } | null>(null);
    const [pos, setPos] = useState<Point>({ x: 0, y: 0 });
    const [error, setError] = useState<string | null>(null);
    const imageRef = useRef<HTMLImageElement | null>(null);
    const drag = useRef<DragState | null>(null);

    useEffect(() => {
        setSize(null);
        setError(null);
        setUrl(null);
        imageRef.current = null;
        if (!file) {
            return undefined;
        }
        const objectUrl = URL.createObjectURL(file);
        const image = new Image();
        image.onload = () => {
            if (Math.min(image.naturalWidth, image.naturalHeight) < AVATAR_CROP_SIZE) {
                setError(
                    `Изображение слишком маленькое: нужно не менее ${AVATAR_CROP_SIZE}×${AVATAR_CROP_SIZE} px`,
                );
                return;
            }
            imageRef.current = image;
            setSize({ width: image.naturalWidth, height: image.naturalHeight });
            setPos({
                x: Math.round((image.naturalWidth - AVATAR_CROP_SIZE) / 2),
                y: Math.round((image.naturalHeight - AVATAR_CROP_SIZE) / 2),
            });
            setUrl(objectUrl);
        };
        image.onerror = () => setError('Не удалось прочитать изображение');
        image.src = objectUrl;
        return () => URL.revokeObjectURL(objectUrl);
    }, [file]);

    const scale = size ? Math.min(1, VIEW_MAX / Math.max(size.width, size.height)) : 1;

    const onPointerDown = (event: React.PointerEvent<HTMLDivElement>) => {
        event.currentTarget.setPointerCapture(event.pointerId);
        drag.current = { pointerX: event.clientX, pointerY: event.clientY, origin: pos };
    };

    const onPointerMove = (event: React.PointerEvent<HTMLDivElement>) => {
        if (!drag.current || !size) {
            return;
        }
        const { pointerX, pointerY, origin } = drag.current;
        setPos({
            x: clamp(origin.x + (event.clientX - pointerX) / scale, 0, size.width - AVATAR_CROP_SIZE),
            y: clamp(origin.y + (event.clientY - pointerY) / scale, 0, size.height - AVATAR_CROP_SIZE),
        });
    };

    const onPointerUp = () => {
        drag.current = null;
    };

    const handleConfirm = () => {
        const image = imageRef.current;
        if (!image || !file) {
            return;
        }
        const canvas = document.createElement('canvas');
        canvas.width = AVATAR_CROP_SIZE;
        canvas.height = AVATAR_CROP_SIZE;
        const context = canvas.getContext('2d');
        if (!context) {
            setError('Браузер не поддерживает обрезку изображений');
            return;
        }
        context.drawImage(
            image,
            Math.round(pos.x),
            Math.round(pos.y),
            AVATAR_CROP_SIZE,
            AVATAR_CROP_SIZE,
            0,
            0,
            AVATAR_CROP_SIZE,
            AVATAR_CROP_SIZE,
        );
        const isPng = file.type === 'image/png';
        canvas.toBlob(
            (blob) => {
                if (!blob) {
                    setError('Не удалось обрезать изображение');
                } else if (blob.size > AVATAR_MAX_BYTES) {
                    setError('После обрезки файл всё равно больше 2 МБ — выберите другое фото');
                } else {
                    onConfirm(blob, isPng ? 'avatar.png' : 'avatar.jpg');
                }
            },
            isPng ? 'image/png' : 'image/jpeg',
            0.92,
        );
    };

    return (
        <Modal open={Boolean(file)} onClose={onCancel}>
            <div className={styles.content}>
                <h2 className={styles.title}>Аватар</h2>
                <p className={styles.hint}>
                    Перетащите рамку, чтобы выбрать область {AVATAR_CROP_SIZE}×{AVATAR_CROP_SIZE} px.
                </p>

                {url && size && (
                    <div
                        className={styles.stage}
                        style={{ width: size.width * scale, height: size.height * scale }}
                    >
                        <img className={styles.image} src={url} alt="Выбор области аватара" />
                        <div
                            className={styles.box}
                            style={{
                                left: pos.x * scale,
                                top: pos.y * scale,
                                width: AVATAR_CROP_SIZE * scale,
                                height: AVATAR_CROP_SIZE * scale,
                            }}
                            onPointerDown={onPointerDown}
                            onPointerMove={onPointerMove}
                            onPointerUp={onPointerUp}
                            onPointerCancel={onPointerUp}
                        />
                    </div>
                )}

                {error && <p className={styles.error}>{error}</p>}

                <div className={styles.actions}>
                    <Button variant="outlined" onClick={onCancel} disabled={isSaving}>
                        Отмена
                    </Button>
                    <Button onClick={handleConfirm} disabled={!size || Boolean(error)} isLoading={isSaving}>
                        Сохранить
                    </Button>
                </div>
            </div>
        </Modal>
    );
};

export default AvatarCropModal;
