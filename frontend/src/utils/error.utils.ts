import { isAxiosError } from 'axios';
import { ApiErrorCode } from '@models/apiErrorCodes';

export interface ApiErrorInfo {
    status?: number;
    code?: string;
    message: string;
    isNetworkError: boolean;
}

interface ApiErrorDetails {
    code?: string;
    message?: string;
}

interface ApiErrorBody extends ApiErrorDetails {
    error?: ApiErrorDetails;
}

export const parseApiError = (error: unknown): ApiErrorInfo => {
    if (isAxiosError(error)) {
        const body = error.response?.data as ApiErrorBody | undefined;
        const details = body?.error ?? body;
        return {
            status: error.response?.status,
            code: details?.code,
            message: details?.message || error.message,
            isNetworkError: !error.response,
        };
    }
    return { message: error instanceof Error ? error.message : '', isNetworkError: false };
};

const CODE_MESSAGES: Partial<Record<string, string>> = {
    [ApiErrorCode.InvalidRequest]: 'Что-то пошло не так',
    [ApiErrorCode.Unauthorized]: 'Нужно войти в аккаунт',
    [ApiErrorCode.Forbidden]: 'Недостаточно прав',
    [ApiErrorCode.NotFound]: 'Не найдено',
    [ApiErrorCode.Conflict]: 'Конфликт данных',
    [ApiErrorCode.EmailAlreadyUsed]: 'Email уже используется',
    [ApiErrorCode.UsernameAlreadyUsed]: 'Никнейм занят',
    [ApiErrorCode.TokenExpired]: 'Сессия истекла. Войдите снова',
    [ApiErrorCode.SessionRevokedOrExpired]: 'Сессия истекла. Войдите снова',
    [ApiErrorCode.RequestTooLarge]: 'Файл слишком большой',
    [ApiErrorCode.UnsupportedImageType]: 'Неподходящий формат файла',
    [ApiErrorCode.Internal]: 'Ошибка сервера. Попробуйте позже',
};

const STATUS_MESSAGES: Record<number, string> = {
    401: 'Нужно войти в аккаунт',
    403: 'Недостаточно прав',
    404: 'Не найдено',
    409: 'Конфликт данных',
    413: 'Файл слишком большой',
    415: 'Неподходящий формат файла',
    429: 'Слишком много запросов. Попробуйте позже',
};

export const getApiErrorMessage = (error: unknown, fallback: string): string => {
    const { status, code, isNetworkError, message } = parseApiError(error);
    if (isNetworkError) {
        return 'Нет соединения с сервером';
    }
    if (code && CODE_MESSAGES[code]) {
        return CODE_MESSAGES[code];
    }
    if (status && STATUS_MESSAGES[status]) {
        return STATUS_MESSAGES[status];
    }
    if (status && status >= 500) {
        return 'Ошибка сервера. Попробуйте позже';
    }
    return !status && message ? message : fallback;
};
