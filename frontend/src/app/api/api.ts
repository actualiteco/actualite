import axios, { AxiosError, InternalAxiosRequestConfig, isAxiosError } from 'axios';
import { toast } from 'react-toastify';
import type { TokensResponse } from '@models/api.types';
import { clearTokens, getAccessToken, getRefreshToken, setTokens } from '@utils/token.utils';

export const API_URL: string = __API_URL__;

const AUTH_ENDPOINTS = ['/auth/login', '/auth/register', '/auth/refresh', '/auth/logout'];

type RetriableConfig = InternalAxiosRequestConfig & { _retry?: boolean };

let isSessionExpiredHandled = false;

let onUnauthorized: (() => void) | null = null;

export const setUnauthorizedHandler = (handler: () => void): void => {
    onUnauthorized = handler;
};

export const resetSessionExpiredFlag = (): void => {
    isSessionExpiredHandled = false;
};

const expireSession = (): void => {
    if (isSessionExpiredHandled) {
        return;
    }
    isSessionExpiredHandled = true;
    clearTokens();
    onUnauthorized?.();
    toast.error('Сессия истекла. Войдите снова', { toastId: 'session-expired' });
};

const $api = axios.create({
    baseURL: API_URL,
    timeout: 10000,
    withCredentials: false,
});

const refreshClient = axios.create({
    baseURL: API_URL,
    timeout: 10000,
    headers: { 'Content-Type': 'application/json' },
});

let refreshPromise: Promise<string> | null = null;

const refreshAccessToken = (): Promise<string> => {
    if (refreshPromise) {
        return refreshPromise;
    }
    const refreshToken = getRefreshToken();
    if (!refreshToken) {
        return Promise.reject(new Error('Нет refresh-токена'));
    }
    const promise: Promise<string> = refreshClient
        .post<TokensResponse>('/auth/refresh', { refresh_token: refreshToken })
        .then(({ data }) => {
            if (!data.access_token || !data.refresh_token) {
                throw new Error('Сервер не вернул токены');
            }
            setTokens(data.access_token, data.refresh_token);
            return data.access_token;
        })
        .finally(() => {
            refreshPromise = null;
        });
    refreshPromise = promise;
    return promise;
};

$api.interceptors.request.use((config) => {
    const token = getAccessToken();
    if (token) {
        config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
});

$api.interceptors.response.use(
    (response) => response,
    async (error: AxiosError) => {
        if (error.code === 'ERR_NETWORK' || error.message === 'Network Error') {
            toast.error('Не удаётся подключиться к серверу', {
                autoClose: 8000,
                toastId: 'network-error',
            });
        }

        const config = error.config as RetriableConfig | undefined;
        const url = config?.url ?? '';
        const isAuthEndpoint = AUTH_ENDPOINTS.some((endpoint) => url.endsWith(endpoint));

        if (error.response?.status === 401 && config && !isAuthEndpoint) {
            if (!config._retry && getRefreshToken()) {
                config._retry = true;

                const sentToken = String(config.headers.Authorization ?? '').replace('Bearer ', '');
                const currentToken = getAccessToken();
                if (currentToken && sentToken && sentToken !== currentToken) {
                    config.headers.Authorization = `Bearer ${currentToken}`;
                    return $api(config);
                }

                try {
                    const token = await refreshAccessToken();
                    config.headers.Authorization = `Bearer ${token}`;
                    return $api(config);
                } catch (refreshError) {
                    if (isAxiosError(refreshError) && !refreshError.response) {
                        return Promise.reject(error);
                    }
                }
            }
            expireSession();
        }
        return Promise.reject(error);
    },
);

export default $api;
