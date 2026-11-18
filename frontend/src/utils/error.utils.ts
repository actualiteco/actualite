import { isAxiosError } from 'axios';

export const ApiErrorCode = {
    InvalidRequest: 'INVALID_REQUEST',
    Internal: 'INTERNAL_ERROR',
    Unauthorized: 'UNAUTHORIZED',
    EmailAlreadyUsed: 'EMAIL_ALREADY_USED',
    Forbidden: 'FORBIDDEN',
    NotFound: 'NOT_FOUND',
    Conflict: 'CONFLICT',
} as const;

export interface ApiErrorInfo {
    status?: number;
    code?: string;
    message: string;
    isNetwork: boolean;
}

interface ApiErrorPayload {
    code?: string;
    message?: string;
}

interface ApiErrorBody extends ApiErrorPayload {
    error?: ApiErrorPayload;
}

const isObject = (value: unknown): value is Record<string, unknown> =>
    typeof value === 'object' && value !== null;

export const parseApiError = (error: unknown): ApiErrorInfo => {
    if (isAxiosError(error)) {
        const data: unknown = error.response?.data;
        const body: ApiErrorBody | undefined = isObject(data) ? (data as ApiErrorBody) : undefined;
        const payload = body?.error ?? body;

        return {
            status: error.response?.status,
            code: payload?.code,
            message: payload?.message || error.message,
            isNetwork: !error.response,
        };
    }
    return { message: error instanceof Error ? error.message : '', isNetwork: false };
};
