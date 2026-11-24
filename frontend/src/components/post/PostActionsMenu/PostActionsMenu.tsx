import React, { FC, useEffect, useRef, useState } from 'react';
import Icon from '@components/shared/Icon/Icon';
import styles from './PostActionsMenu.module.scss';

interface PostActionsMenuProps {
    isDeleting: boolean;
    onDelete: () => void;
}

const PostActionsMenu: FC<PostActionsMenuProps> = ({ isDeleting, onDelete }) => {
    const [isOpen, setIsOpen] = useState(false);
    const [isConfirming, setIsConfirming] = useState(false);
    const wrapperRef = useRef<HTMLDivElement>(null);

    const close = () => {
        setIsOpen(false);
        setIsConfirming(false);
    };

    useEffect(() => {
        if (!isOpen) {
            return undefined;
        }
        const onPointerDown = (event: MouseEvent) => {
            if (wrapperRef.current && !wrapperRef.current.contains(event.target as Node)) {
                close();
            }
        };
        const onKeyDown = (event: KeyboardEvent) => {
            if (event.key === 'Escape') {
                close();
            }
        };
        document.addEventListener('mousedown', onPointerDown);
        document.addEventListener('keydown', onKeyDown);
        return () => {
            document.removeEventListener('mousedown', onPointerDown);
            document.removeEventListener('keydown', onKeyDown);
        };
    }, [isOpen]);

    const handleDelete = () => {
        if (!isConfirming) {
            setIsConfirming(true);
            return;
        }
        onDelete();
    };

    return (
        <div className={styles.wrapper} ref={wrapperRef}>
            <button
                type="button"
                className={styles.trigger}
                aria-label="Действия с постом"
                aria-haspopup="menu"
                aria-expanded={isOpen}
                onClick={() => (isOpen ? close() : setIsOpen(true))}
            >
                <Icon name="more" size={18} />
            </button>

            {isOpen && (
                <div className={styles.menu} role="menu">
                    <button
                        type="button"
                        role="menuitem"
                        className={`${styles.item} ${styles.danger}`}
                        onClick={handleDelete}
                        disabled={isDeleting}
                    >
                        <Icon name="trash" size={16} />
                        {isConfirming ? 'Точно удалить?' : 'Удалить'}
                    </button>
                </div>
            )}
        </div>
    );
};

export default PostActionsMenu;
