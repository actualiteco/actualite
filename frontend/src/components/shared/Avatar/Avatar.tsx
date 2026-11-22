import React, { FC, useEffect, useState } from 'react';
import styles from './Avatar.module.scss';

interface AvatarProps {
    url?: string | null;
    name?: string | null;
    size?: number;
}

const Avatar: FC<AvatarProps> = ({ url, name, size = 40 }) => {
    const [failed, setFailed] = useState(false);

    useEffect(() => setFailed(false), [url]);

    const letter = (name ?? '').trim().charAt(0) || '?';

    return (
        <span
            className={styles.avatar}
            style={{ width: size, height: size, fontSize: Math.round(size * 0.4) }}
        >
            {url && !failed ? (
                <img className={styles.image} src={url} alt="" onError={() => setFailed(true)} />
            ) : (
                letter
            )}
        </span>
    );
};

export default Avatar;
