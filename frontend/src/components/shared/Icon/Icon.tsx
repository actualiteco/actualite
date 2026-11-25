import React, { FC, ReactNode } from 'react';

export type IconName =
    | 'home'
    | 'search'
    | 'bell'
    | 'user'
    | 'logout'
    | 'comment'
    | 'heart'
    | 'arrowLeft'
    | 'close'
    | 'more'
    | 'reply'
    | 'trash'
    | 'image';

const ICONS: Record<IconName, ReactNode> = {
    home: <path d="M3.5 10.5 12 3.5l8.5 7V20a.5.5 0 0 1-.5.5h-5v-6h-6v6H4a.5.5 0 0 1-.5-.5z" />,
    search: (
        <>
            <circle cx="11" cy="11" r="6.5" />
            <path d="m16 16 4.5 4.5" />
        </>
    ),
    bell: (
        <>
            <path d="M6 9.5a6 6 0 1 1 12 0c0 6 2.5 7.5 2.5 7.5h-17S6 15.5 6 9.5z" />
            <path d="M10 20.5a2.2 2.2 0 0 0 4 0" />
        </>
    ),
    user: (
        <>
            <circle cx="12" cy="8" r="4" />
            <path d="M4.5 20.5a7.5 7.5 0 0 1 15 0" />
        </>
    ),
    logout: (
        <>
            <path d="M9.5 20.5H5a.5.5 0 0 1-.5-.5V4a.5.5 0 0 1 .5-.5h4.5" />
            <path d="m16 8 4 4-4 4" />
            <path d="M20 12H9.5" />
        </>
    ),
    comment: <path d="M20.5 11.5a8 8 0 0 1-11.7 7.1L3.5 20l1.4-4.9a8 8 0 1 1 15.6-3.6z" />,
    heart: (
        <path d="M12 20.5s-8.5-5-8.5-11.2A4.6 4.6 0 0 1 12 6.7a4.6 4.6 0 0 1 8.5 2.6C20.5 15.5 12 20.5 12 20.5z" />
    ),
    arrowLeft: (
        <>
            <path d="M19.5 12h-15" />
            <path d="m10.5 6-6 6 6 6" />
        </>
    ),
    close: <path d="m6 6 12 12M18 6 6 18" />,
    more: (
        <>
            <circle cx="5.5" cy="12" r="1.6" fill="currentColor" stroke="none" />
            <circle cx="12" cy="12" r="1.6" fill="currentColor" stroke="none" />
            <circle cx="18.5" cy="12" r="1.6" fill="currentColor" stroke="none" />
        </>
    ),
    reply: (
        <>
            <path d="M5 5v5.5a3 3 0 0 0 3 3h11" />
            <path d="m15 9.5 4 4-4 4" />
        </>
    ),
    trash: (
        <>
            <path d="M4 7h16" />
            <path d="M9.5 7V4.5h5V7" />
            <path d="m6 7 .9 12.5a1 1 0 0 0 1 .9h8.2a1 1 0 0 0 1-.9L18 7" />
            <path d="M10 11v5.5M14 11v5.5" />
        </>
    ),
    image: (
        <>
            <rect x="3.5" y="4.5" width="17" height="15" rx="1.5" />
            <circle cx="9" cy="10" r="1.5" />
            <path d="m20.5 16-5-5-8 8.5" />
        </>
    ),
};

interface IconProps {
    name: IconName;
    size?: number;
    /** Заливка (например, активный лайк) */
    filled?: boolean;
    className?: string;
}

const Icon: FC<IconProps> = ({ name, size = 18, filled = false, className }) => (
    <svg
        width={size}
        height={size}
        viewBox="0 0 24 24"
        fill={filled ? 'currentColor' : 'none'}
        stroke="currentColor"
        strokeWidth={1.5}
        strokeLinecap="round"
        strokeLinejoin="round"
        aria-hidden="true"
        focusable="false"
        className={className}
        style={{ display: 'block', flexShrink: 0 }}
    >
        {ICONS[name]}
    </svg>
);

export default Icon;
