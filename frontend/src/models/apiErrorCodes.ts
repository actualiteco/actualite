export const ApiErrorCode = {
    InvalidRequest: 'INVALID_REQUEST',
    Internal: 'INTERNAL_ERROR',
    Unauthorized: 'UNAUTHORIZED',
    EmailAlreadyUsed: 'EMAIL_ALREADY_USED',
    Forbidden: 'FORBIDDEN',
    NotFound: 'NOT_FOUND',
    Conflict: 'CONFLICT',
    UsernameAlreadyUsed: 'USERNAME_ALREADY_USED',
    TokenExpired: 'TOKEN_EXPIRED',
    SessionRevokedOrExpired: 'SESSION_REVOKED_OR_EXPIRED',
    // todo: после мержа добавить нижние плдчеркивания
    RequestTooLarge: 'REQUEST IS TOO LARGE',
    UnsupportedImageType: 'UNSUPPORTED IMAGE TYPE',
} as const;

export type ApiErrorCode = (typeof ApiErrorCode)[keyof typeof ApiErrorCode];
