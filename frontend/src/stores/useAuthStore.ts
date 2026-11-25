import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import { AuthService } from '@services/AuthService';
import { ProfileService } from '@services/ProfileService';
import { resetSessionExpiredFlag, setUnauthorizedHandler } from '@api/api';
import type { Profile } from '@models/Profile';
import { clearTokens, getAccessToken, hasSession, setTokens } from '@utils/token.utils';

interface AuthState {
    isAuth: boolean;
    userId: string | null;
    email: string | null;
    username: string | null;
    displayName: string | null;
    avatarUrl: string | null;
    isLoading: boolean;

    login: (email: string, password: string) => Promise<void>;
    register: (email: string, username: string, password: string) => Promise<void>;
    logout: () => void;
    clearSession: () => void;
    setProfile: (profile: Profile) => void;
    loadProfile: () => Promise<void>;
}

type AuthData = Pick<
    AuthState,
    'isAuth' | 'userId' | 'email' | 'username' | 'displayName' | 'avatarUrl' | 'isLoading'
>;

const initialState: AuthData = {
    isAuth: false,
    userId: null,
    email: null,
    username: null,
    displayName: null,
    avatarUrl: null,
    isLoading: false,
};

const profileToState = (profile: Profile) => ({
    userId: profile.id || null,
    username: profile.username || null,
    displayName: profile.displayName || null,
    avatarUrl: profile.avatarUrl || null,
});

export const useAuthStore = create<AuthState>()(
    persist(
        (set, get) => ({
            ...initialState,

            async login(email: string, password: string) {
                set({ isLoading: true });
                try {
                    const tokens = await AuthService.login({ email, password });
                    if (!tokens.access_token || !tokens.refresh_token) {
                        throw new Error('Сервер не вернул токены');
                    }
                    setTokens(tokens.access_token, tokens.refresh_token);
                    resetSessionExpiredFlag();

                    let profile: Profile;
                    try {
                        profile = await ProfileService.getMe();
                    } catch (error) {
                        clearTokens();
                        throw error;
                    }
                    set({ isAuth: true, isLoading: false, email, ...profileToState(profile) });
                } catch (e) {
                    set({ isLoading: false });
                    throw e;
                }
            },

            async register(email: string, username: string, password: string) {
                set({ isLoading: true });
                try {
                    await AuthService.register({ email, username, password, role: 'user' });
                } catch (e) {
                    set({ isLoading: false });
                    throw e;
                }
                await get().login(email, password);
            },

            logout() {
                const accessToken = getAccessToken();
                get().clearSession();
                if (accessToken) {
                    AuthService.logout(accessToken).catch(() => undefined);
                }
            },

            clearSession() {
                clearTokens();
                set({ ...initialState });
            },

            setProfile(profile: Profile) {
                set(profileToState(profile));
            },

            async loadProfile() {
                const profile = await ProfileService.getMe();
                get().setProfile(profile);
            },
        }),
        {
            name: 'auth-storage',
            partialize: (state) => ({
                userId: state.userId,
                email: state.email,
                username: state.username,
                displayName: state.displayName,
                avatarUrl: state.avatarUrl,
            }),
            merge: (persisted, current) => {
                if (!hasSession()) {
                    clearTokens();
                    return { ...current, ...initialState };
                }
                return { ...current, ...(persisted as Partial<AuthState>), isAuth: true };
            },
        },
    ),
);

setUnauthorizedHandler(() => useAuthStore.getState().clearSession());
