export enum AppRoutes {
    AUTH = 'auth',
    FEED = 'feed',
    POST = 'post',
    PROFILE = 'profile',
    OTHERS_PROFILE = 'others_profile',
}

export const RoutePath: Record<AppRoutes, string> = {
    [AppRoutes.AUTH]: '/auth',
    [AppRoutes.FEED]: '/feed',
    [AppRoutes.POST]: '/post/:id',
    [AppRoutes.PROFILE]: '/profile',
    [AppRoutes.OTHERS_PROFILE]: '/profile/:id',
};

export const getPostPath = (postId: string): string => `/post/${postId}`;
export const getProfilePath = (userId: string): string => `/profile/${userId}`;
