export interface JwtPayload {
    exp?: number;
    iat?: number;
    [claim: string]: unknown;
}

export const parseJwt = (token: string): JwtPayload | null => {
    try {
        const part = token.split('.')[1];
        if (!part) {
            return null;
        }
        const base64 = part.replace(/-/g, '+').replace(/_/g, '/');
        const padded = base64.padEnd(base64.length + ((4 - (base64.length % 4)) % 4), '=');
        const json = decodeURIComponent(
            atob(padded)
                .split('')
                .map((char) => '%' + ('00' + char.charCodeAt(0).toString(16)).slice(-2))
                .join(''),
        );
        return JSON.parse(json) as JwtPayload;
    } catch {
        return null;
    }
};

export const isTokenExpired = (token: string): boolean => {
    const payload = parseJwt(token);
    if (!payload || typeof payload.exp !== 'number') {
        return false;
    }
    return payload.exp * 1000 <= Date.now();
};

export const truncate = (text: string, max: number): string => {
    if (!text) return '';
    if (text.length <= max) return text;
    return text.slice(0, max).trimEnd() + '…';
};
