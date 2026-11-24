const MINUTE = 60 * 1000;
const HOUR = 60 * MINUTE;
const DAY = 24 * HOUR;
const WEEK = 7 * DAY;

export const formatRelativeTime = (
    iso: string,
    now: number = Date.now(),
): string => {
    const timestamp = Date.parse(iso);

    if (Number.isNaN(timestamp)) {
        return '';
    }

    const diff = Math.max(0, now - timestamp);

    if (diff < MINUTE) {
        return 'только что';
    }

    if (diff < HOUR) {
        return `${Math.floor(diff / MINUTE)} мин назад`;
    }

    if (diff < DAY) {
        return `${Math.floor(diff / HOUR)} ч назад`;
    }

    if (diff < WEEK) {
        return `${Math.floor(diff / DAY)} дн назад`;
    }

    return new Date(timestamp).toLocaleDateString('ru-RU', {
        day: 'numeric',
        month: 'long',
        year: 'numeric',
    });
};

const formatter = new Intl.DateTimeFormat('ru', {
    day: 'numeric',
    month: 'short',
    hour: '2-digit',
    minute: '2-digit',
});

const formatterWithYear = new Intl.DateTimeFormat('ru', {
    day: 'numeric',
    month: 'short',
    year: 'numeric',
});

export const formatPostDate = (iso?: string): string => {
    if (!iso) {
        return '';
    }
    const date = new Date(iso);
    if (Number.isNaN(date.getTime())) {
        return '';
    }
    return date.getFullYear() === new Date().getFullYear()
        ? formatter.format(date)
        : formatterWithYear.format(date);
};
