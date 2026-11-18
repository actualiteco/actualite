import * as Yup from 'yup';

// F_Prof_1, F_Prof_2: никнейм 4-15 символов из [a-zA-Z0-9_]
export const USERNAME_REGEX = /^[a-zA-Z0-9_]+$/;
export const USERNAME_MIN = 4;
export const USERNAME_MAX = 15;

// F_Prof_4, NF_Data_4: отображаемое имя 1-30 символов
export const DISPLAY_NAME_MIN = 1;
export const DISPLAY_NAME_MAX = 30;

// F_Prof_5, NF_Data_5: статус до 100 символов
export const STATUS_MAX = 100;

/**
 * Длина строки в символах (как runes на бэкенде): эмодзи = 1.
 * `str.length` считает UTF-16 единицы и даёт 2 для эмодзи.
 */
export const charCount = (value: string): number => Array.from(value).length;

/** Обрезка до `max` символов без разрыва суррогатных пар. */
export const truncateChars = (value: string, max: number): string => {
    const chars = Array.from(value);
    return chars.length <= max ? value : chars.slice(0, max).join('');
};

/** Оставляет в никнейме только допустимые символы и обрезает до максимума */
export const sanitizeUsername = (value: string): string =>
    truncateChars(value.replace(/[^a-zA-Z0-9_]/g, ''), USERNAME_MAX);

export const isBlank = (value: string | undefined | null): boolean => !value || value.trim() === '';

export const usernameField = () =>
    Yup.string()
        .required('Введите никнейм')
        .min(USERNAME_MIN, `Никнейм: от ${USERNAME_MIN} до ${USERNAME_MAX} символов`)
        .max(USERNAME_MAX, `Никнейм: от ${USERNAME_MIN} до ${USERNAME_MAX} символов`)
        .matches(USERNAME_REGEX, 'Никнейм может содержать только a-z, A-Z, 0-9 и _');

export const displayNameField = () =>
    Yup.string()
        .test('not-blank', 'Введите имя', (value) => !isBlank(value))
        .test(
            'max-chars',
            `Имя: от ${DISPLAY_NAME_MIN} до ${DISPLAY_NAME_MAX} символов`,
            (value) => !value || charCount(value) <= DISPLAY_NAME_MAX,
        );

export const statusField = () =>
    Yup.string().test(
        'max-chars',
        `Статус: не более ${STATUS_MAX} символов`,
        (value) => !value || charCount(value) <= STATUS_MAX,
    );
