import { create } from 'zustand';
import { ProfileService } from '@services/ProfileService';
import type { Profile, ProfileChanges } from '@models/Profile';
import { useAuthStore } from '@stores/useAuthStore';
import { parseApiError } from '@utils/error.utils';
import { ProfileErrorType } from '@components/profile/ProfileError/ProfileError';

export type ProfileLoadStatus = 'idle' | 'loading' | 'success' | 'error';
export type ProfileLoadError = ProfileErrorType;

interface ProfileState {
    profile: Profile | null;
    status: ProfileLoadStatus;
    error: ProfileLoadError | null;
    isFollowPending: boolean;

    fetchProfile: (id?: string) => Promise<void>;
    setProfile: (profile: Profile) => void;
    updateProfile: (changes: ProfileChanges) => Promise<Profile>;
    uploadAvatar: (file: Blob, fileName: string) => Promise<string>;
    toggleFollow: () => Promise<void>;
    reset: () => void;
}

type ProfileData = Pick<
    ProfileState,
    'profile' | 'status' | 'error' | 'isFollowPending'
>;

const initialState: ProfileData = {
    profile: null,
    status: 'idle',
    error: null,
    isFollowPending: false,
};

let lastRequestId = 0;

const toLoadError = (error: unknown): ProfileErrorType => {
    const { status, isNetworkError } = parseApiError(error);

    if (isNetworkError) return 'network';
    if (status === 404 || status === 400) return 'notFound';
    if (status === 403) return 'forbidden';

    return 'unknown';
};

export const useProfileStore = create<ProfileState>()((set, get) => ({
    ...initialState,

    async fetchProfile(id?: string) {
        const requestId = ++lastRequestId;

        set({
            profile: null,
            status: 'loading',
            error: null,
            isFollowPending: false,
        });

        try {
            const profile = id
                ? await ProfileService.getById(id)
                : await ProfileService.getMe();

            if (requestId !== lastRequestId) return;

            set({
                profile,
                status: 'success',
                error: null,
            });

            if (!id) {
                useAuthStore.getState().setProfile(profile);
            }
        } catch (error) {
            if (requestId !== lastRequestId) return;

            set({
                profile: null,
                status: 'error',
                error: toLoadError(error),
            });
        }
    },

    setProfile(profile: Profile) {
        set({
            profile,
            status: 'success',
            error: null,
        });
    },

    async updateProfile(changes: ProfileChanges) {
        const updated = await ProfileService.update(changes);

        set((state) => ({
            profile: state.profile
                ? {
                      ...state.profile,
                      ...updated,
                  }
                : updated,
            status: 'success',
            error: null,
        }));

        const profile = get().profile;

        if (profile && profile.id === updated.id) {
            useAuthStore.getState().setProfile(profile);
            return profile;
        }

        return updated;
    },

    async uploadAvatar(file: Blob, fileName: string) {
        const avatarUrl = await ProfileService.uploadAvatar(file, fileName);

        set((state) => ({
            profile: state.profile
                ? { ...state.profile, avatarUrl }
                : null,
        }));

        const profile = get().profile;

        if (profile && profile.id === useAuthStore.getState().userId) {
            useAuthStore.getState().setProfile(profile);
        }

        return avatarUrl;
    },

    async toggleFollow() {
        const { profile, isFollowPending } = get();

        if (!profile || isFollowPending) return;

        const targetId = profile.id;
        const wasFollowing = profile.isFollowing === true;

        set({ isFollowPending: true });

        const finish = (isFollowing: boolean | null, delta: number) => {
            set((state) => {
                if (!state.profile || state.profile.id !== targetId) {
                    return { isFollowPending: false };
                }

                if (isFollowing === null) {
                    return { isFollowPending: false };
                }

                return {
                    isFollowPending: false,
                    profile: {
                        ...state.profile,
                        isFollowing,
                        followersCount: Math.max(
                            0,
                            state.profile.followersCount + delta,
                        ),
                    },
                };
            });
        };

        try {
            if (wasFollowing) {
                await ProfileService.unfollow(targetId);
            } else {
                await ProfileService.follow(targetId);
            }

            finish(!wasFollowing, wasFollowing ? -1 : 1);
        } catch (error) {
            if (parseApiError(error).status === 400) {
                finish(!wasFollowing, 0);
                return;
            }

            finish(null, 0);
            throw error;
        }
    },

    reset() {
        lastRequestId += 1;
        set({ ...initialState });
    },
}));