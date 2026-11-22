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
