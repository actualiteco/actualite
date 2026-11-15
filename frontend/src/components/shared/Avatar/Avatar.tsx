import { FC, useState } from 'react';
import styles from './Avatar.module.scss';

interface AvatarProps {
    name: string;
    src?: string;
    size?: number;
}

const Avatar: FC<AvatarProps> = ({ name, src, size = 72 }) => {
    const [failed, setFailed] = useState(false);

    const style = { width: size, height: size, fontSize: Math.round(size * 0.4) };
    const letter = Array.from(name.trim())[0]?.toUpperCase() ?? '?';

    if (src && !failed) {
        return (
            <img
                className={styles.avatar}
                style={style}
                src={src}
                alt={name}
                onError={() => setFailed(true)}
            />
        );
    }

    return (
        <div className={`${styles.avatar} ${styles.placeholder}`} style={style} aria-hidden="true">
            {letter}
        </div>
    );
};

export default Avatar;
