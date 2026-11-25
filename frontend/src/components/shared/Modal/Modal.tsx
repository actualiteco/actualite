import React, { FC, ReactNode, useEffect } from 'react';
import Icon from '@components/shared/Icon/Icon';
import styles from './Modal.module.scss';

interface ModalProps {
    open: boolean;
    onClose: () => void;
    children: ReactNode;
}

const Modal: FC<ModalProps> = ({ open, onClose, children }) => {
    useEffect(() => {
        if (!open) {
            return undefined;
        }
        const onKeyDown = (event: KeyboardEvent) => {
            if (event.key === 'Escape') {
                onClose();
            }
        };
        document.addEventListener('keydown', onKeyDown);
        return () => document.removeEventListener('keydown', onKeyDown);
    }, [open, onClose]);

    if (!open) return null;

    return (
        <div className={styles.overlay} onClick={onClose}>
            <div
                className={styles.modal}
                role="dialog"
                aria-modal="true"
                onClick={(e) => e.stopPropagation()}
            >
                <button type="button" className={styles.close} onClick={onClose} aria-label="Закрыть">
                    <Icon name="close" size={16} />
                </button>
                {children}
            </div>
        </div>
    );
};

export default Modal;
