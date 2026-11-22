import { safeGetItem, safeRemoveItem, safeSetItem } from './storage.utils';

const ACCESS_KEY = 'access_token';
const REFRESH_KEY = 'refresh_token';
const LEGACY_KEY = 'token';

export const getAccessToken = (): string | null => safeGetItem(ACCESS_KEY);
export const getRefreshToken = (): string | null => safeGetItem(REFRESH_KEY);

export const setTokens = (accessToken: string, refreshToken: string): void => {
    const saved = safeSetItem(ACCESS_KEY, accessToken) && safeSetItem(REFRESH_KEY, refreshToken);
    if (!saved) {
        throw new Error('Не удалось сохранить сессию: хранилище браузера недоступно');
    }
};

export const clearTokens = (): void => {
    safeRemoveItem(ACCESS_KEY);
    safeRemoveItem(REFRESH_KEY);
    safeRemoveItem(LEGACY_KEY);
};

export const hasSession = (): boolean => Boolean(getRefreshToken());
